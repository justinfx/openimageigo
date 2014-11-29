package oiio

/*
#include "stdlib.h"

#include "cpp/imagebufalgo.h"

*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	// Let OIIO choose the best filter
	FilterDefault = ""
	// Let OIIO choose the best filter and filter width
	FilterDefaultWidth = 0.0
	// Use the global OIIO-determined thread count
	GlobalThreads = 0
)

// AlgoOpts allows common arguments to be passed to algorithm functions.
type AlgoOpts struct {
	// The Region of interest to use when applying the operation.
	// Some algorithms require either an ROI or to fallback on the
	// spec of the ImageBuf being given.
	ROI *ROI
	// The Threads parameter specifies how many threads (potentially) may
	// be used, but it's not a guarantee.
	// If Threads == 0, it will use the global OIIO attribute "threads".
	// If Threads == 1, it guarantees that it will
	// not launch any new threads.
	Threads int
}

func checkBufAndROI(dst *ImageBuf, roi *ROI) error {
	if dst.Initialized() || (roi != nil && roi.ptr != nil && roi.Defined()) {
		return nil
	}
	return errors.New("ImageBuf and ROI cannot both be undefined." +
		"ImageBufAlgo without any guess about region of interest")
}

func flatAlgoOpts(opts []AlgoOpts) AlgoOpts {
	var opt AlgoOpts
	for _, o := range opts {
		opt.Threads = o.Threads
		if o.ROI != nil {
			opt.ROI = o.ROI
		}
	}
	return opt
}

// Zero out (set to 0, black) the image region.
// Only the pixels (and channels) in dst that are specified by roi will be altered;
// the default roi is to alter all the pixels in dst.
// If dst is uninitialized, it will be resized to be a float ImageBuf large enough to
// hold the region specified by roi. It is an error to pass both an uninitialied dst and an undefined roi.
// Works on all pixel data types.
func Zero(dst *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)
	err := checkBufAndROI(dst, opt.ROI)
	if err != nil {
		return err
	}

	ok := bool(C.zero(dst.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads)))
	if !ok {
		return dst.LastError()
	}

	return nil
}

// Fill sets the pixels in the destination image within the specified region to the values in []values.
// The values slice must point to at least chend values, or the number of channels in the
// image, whichever is smaller.
func Fill(dst *ImageBuf, values []float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)
	err := checkBufAndROI(dst, opt.ROI)
	if err != nil {
		return err
	}

	c_ptr := (*C.float)(unsafe.Pointer(&values[0]))
	ok := bool(C.fill(dst.ptr, c_ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads)))
	if !ok {
		return dst.LastError()
	}

	return nil
}

// ChannelOpts are options that can be passed to the Channels() function
//
// For any channel in which ChannelOpts.Order[i] < 0, it will just make dst channel i be a constant value
// set to ChannelOpts.Values[i] (if ChannelOpts.Values is not nil) or 0.0 (if ChannelOpts.Values is
// nil).
// If ChannelOpts.Order is nil, it will be interpreted as {0, 1, ..., nchannels-1}, mean-
// ing that it’s only renaming channels, not reordering them.
// If ChannelOpts.NewNames is not nil, it points to a list of new channel names. Channels
// for which ChannelOpts.NewNames[i] is the empty string (or all channels, if ChannelOpts.NewNames
// == nil) will be named as follows: If ChannelOpts.ShuffleNames is false, the result-
// ing dst image will have default channel names in the usual order ("R", "G", etc.), but
// if ChannelOpts.ShuffleNames is true, the names will be taken from the corresponding
// channels of the source image – be careful with this, shuffling both channel ordering and
// their names could result in no semantic change at all, if you catch the drift.
type ChannelOpts struct {
	Order        []int32
	Values       []float32
	NewNames     []string
	ShuffleNames bool
}

// Generic channel shuffling: copy src to dst, but with channels in the order specified by
// ChannelOpts.Order[0..nchannels-1]. Does not support in-place operation.
//
// See ChannelOpts docs for more details on the options.
func Channels(dst, src *ImageBuf, nchannels int, opts ...*ChannelOpts) error {
	var (
		order    *C.int32_t
		values   *C.float
		newNames **C.char
		shuffle  C.bool = C.bool(false)
	)

	var opt *ChannelOpts
	if len(opts) > 0 {
		opt = opts[len(opts)-1]
	}

	if opt != nil {
		shuffle = C.bool(opt.ShuffleNames)

		if opt.Order != nil {
			if len(opt.Order) < nchannels {
				return fmt.Errorf("ChannelOpts.Order length %d is less than nchannels %d",
					len(opt.Order), nchannels)
			}
			order = (*C.int32_t)(unsafe.Pointer(&opt.Order[0]))
		}

		if opt.Values != nil {
			if len(opt.Values) < nchannels {
				return fmt.Errorf("ChannelOpts.Values length %d is less than nchannels %d",
					len(opt.Values), nchannels)
			}
			values = (*C.float)(unsafe.Pointer(&opt.Values[0]))
		}

		if opt.NewNames != nil {
			if len(opt.NewNames) < nchannels {
				return fmt.Errorf("ChannelOpts.NewNames length %d is less than nchannels %d",
					len(opt.NewNames), nchannels)
			}
			nameSize := len(opt.NewNames)
			newNames = C.makeCharArray(C.int(nameSize))
			defer C.freeCharArray(newNames, C.int(nameSize))
			for i, s := range opt.NewNames {
				C.setArrayString(newNames, C.CString(s), C.int(i))
			}
		}
	}

	ok := C.channels(dst.ptr, src.ptr, C.int(nchannels), order, values, newNames, shuffle)
	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// Copy into dst, beginning at (xbegin, ybegin), the pixels of src described by roi.
// If roi is nil, the entirety of src will be used.
func Paste2D(dst, src *ImageBuf, xbegin, ybegin int, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.paste(dst.ptr, C.int(xbegin), C.int(ybegin), C.int(0), C.int(0),
		src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Copy pixels within the ROI from src to dst, applying a color transform.
// If dst is not yet initialized, it will be allocated to the same size as specified by roi. If roi is not
// defined it will be all of dst, if dst is defined, or all of src, if dst is not yet defined.
// In-place operations (dst == src) are supported.
// If unpremult is true, unpremultiply before color conversion, then premultiply after the color conversion.
// You may want to use this flag if your image contains an alpha channel.
// Works with all data types.
func ColorConvert(dst, src *ImageBuf, from, to string, unpremult bool, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_from := C.CString(from)
	defer C.free(unsafe.Pointer(c_from))

	c_to := C.CString(to)
	defer C.free(unsafe.Pointer(c_to))

	ok := C.colorconvert(dst.ptr, src.ptr, c_from, c_to, C.bool(unpremult),
		opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// Copy pixels within the ROI from src to dst, applying a color transform.
// If dst is not yet initialized, it will be allocated to the same size as specified by roi.
// If roi is not defined it will be all of dst, if dst is defined, or all of src, if dst is not yet defined.
// In-place operations (dst == src) are supported.
// If unpremult is true, unpremultiply before color conversion, then premultiply after the color conversion.
// You may want to use this flag if your image contains an alpha channel.
// Works with all data types.
func ColorConvertProcessor(dst, src *ImageBuf, cp *ColorProcessor, unpremult bool, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.colorconvert_processor(dst.ptr, src.ptr, cp.ptr, C.bool(unpremult),
		opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// Set dst, over the region of interest, to be a resized version of the corresponding portion of src
// (mapping such that the "full" image window of each correspond to each other, regardless of resolution).
// Will choose a reasonable default high-quality default filter (blackman-harris when upsizing, lanczos3 when downsizing)
// Works on all pixel data types.
func Resize(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_filtname := C.CString("")
	defer C.free(unsafe.Pointer(c_filtname))

	ok := C.resize(dst.ptr, src.ptr, c_filtname, C.float(0.0), opt.ROI.validOrAllPtr(), C.int(opt.Threads))
	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Set dst, over the region of interest, to be a resampled version of the corresponding portion of src
// (mapping such that the "full" image window of each correspond to each other, regardless of resolution).
// Unlike Resize(), Resample does not take a filter; it just samples either with a bilinear
// interpolation (if interpolate is true, the default) or uses the single "closest" pixel (if interpolate is false).
// This makes it a lot faster than a proper Resize(), though obviously with lower quality (aliasing when downsizing,
// pixel replication when upsizing).
// Works on all pixel data types.
func Resample(dst, src *ImageBuf, interpolate bool, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.resample(dst.ptr, src.ptr, C.bool(interpolate), opt.ROI.validOrAllPtr(), C.int(opt.Threads))
	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}
