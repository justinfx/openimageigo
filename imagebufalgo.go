package oiio

/*
#include "stdlib.h"

#include "imagebufalgo.h"

*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

const (
	// Let OIIO choose the best filter
	FilterDefault = ""
	// Let OIIO choose the best filter and filter width
	FilterDefaultWidth = 0.0
	// Use a default system font
	FontNameDefault = ""
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

	runtime.KeepAlive(dst)

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

	runtime.KeepAlive(dst)
	runtime.KeepAlive(values)

	if !ok {
		return dst.LastError()
	}

	return nil
}

// Checker2D sets the pixels in the destination image within the specified region to a checkerboard pattern
// with origin given by the offset values, checker size given by the width and height values, and
// alternating between color1 and color2. The colors must contain enough values for all channels in the image.
//
// Example:
//     // Create a new 640x480 RGB image, fill it with a two-toned gray
//     // checkerboard, the checkers being 64x64 pixels each.
//     spec := NewImageSpecSize(640, 480, 3, TypeFloat)
//     dark := []float32{0.1, 0.1, 0.1}
//     light := []float32{0.4, 0.4, 0.4}
//     Checker(spec, 64, 64, dark, light, 0, 0)
//
func Checker2D(dst *ImageBuf, width, height int, color1, color2 []float32,
	xoffset, yoffset int, opts ...AlgoOpts) error {
	return Checker(dst, width, height, 1, color1, color2, xoffset, yoffset, 0, opts...)
}

// Checker sets the pixels in the destination image within the specified region to a checkerboard pattern
// with origin given by the offset values, checker size given by the width, height,
// depth values, and alternating between color1 and color2. The colors must contain enough values
// for all channels in the image.
func Checker(dst *ImageBuf, width, height, depth int, color1, color2 []float32,
	xoffset, yoffset, zoffset int, opts ...AlgoOpts) error {

	opt := flatAlgoOpts(opts)
	err := checkBufAndROI(dst, opt.ROI)
	if err != nil {
		return err
	}

	c1_ptr := (*C.float)(unsafe.Pointer(&color1[0]))
	c2_ptr := (*C.float)(unsafe.Pointer(&color2[0]))

	ok := bool(C.checker(dst.ptr, C.int(width), C.int(height), C.int(depth),
		c1_ptr, c2_ptr, C.int(xoffset), C.int(yoffset), C.int(zoffset),
		opt.ROI.validOrAllPtr(), C.int(opt.Threads)),
	)
	runtime.KeepAlive(dst)
	runtime.KeepAlive(color1)
	runtime.KeepAlive(color2)

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
// ing that it's only renaming channels, not reordering them.
// If ChannelOpts.NewNames is not nil, it points to a list of new channel names. Channels
// for which ChannelOpts.NewNames[i] is the empty string (or all channels, if ChannelOpts.NewNames
// == nil) will be named as follows: If ChannelOpts.ShuffleNames is false, the result-
// ing dst image will have default channel names in the usual order ("R", "G", etc.), but
// if ChannelOpts.ShuffleNames is true, the names will be taken from the corresponding
// channels of the source image - be careful with this, shuffling both channel ordering and
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

	runtime.KeepAlive(src)
	runtime.KeepAlive(dst)
	runtime.KeepAlive(opts)

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// ChannelAppend appends the channels of A and B together into dst over the region of interest.
// If a roi is passed, it will be interpreted as being the union of the pixel windows of
// A and B (and all channels of both images). If dst is not already initialized to a size, it
// will be resized to be big enough for the region.
func ChannelAppend(dst, a, b *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.channel_append(dst.ptr, a.ptr, b.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(opts)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// TODO: Flatten does not seem to work as expected
//
// Flatten copies pixels from deep image src into non-deep dst, compositing the depth samples
// within each pixel to yield a single "flat" value per pixel. If src is not deep, it just copies
// the pixels without alteration.
// func Flatten(dst, src *ImageBuf, opts ...AlgoOpts) error {
// 	opt := flatAlgoOpts(opts)
//
// 	ok := C.flatten(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))
// 	if !bool(ok) {
// 		return dst.LastError()
// 	}
//
// 	return nil
// }

// Crop resets dst to be the specified region of src.
// Note that the crop operation does not actually move the pixels on the image plane or
// adjust the full/display window; it merely restricts which pixels are copied from src to
// dst. (Note the difference compared to Cut()).
func Crop(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.crop(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Set dst to the flattened composite of deep image src. That is, it
// converts a deep image to a simple flat image by front-to-back
// compositing the samples within each pixel.  If src is already a non-
// deep/flat image, it will just copy pixel values from src to dst. If dst
// is not already an initialized ImageBuf, it will be sized to match src
// (but made non-deep).
//
// 'roi' specifies the region of dst's pixels which will be computed;
// existing pixels outside this range will not be altered.  If not
// specified, the default ROI value will be the pixel data window of src.
func Flatten(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.flatten(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Cut assigns to dst the designated region of src, but shifted to be at the
// (0,0) origin, and with the full/display resolution set to be identical
// to the data region.
//
// The nthreads AlgoOpts specifies how many threads (potentially) may
// be used, but it's not a guarantee.  If nthreads == 0, it will use
// the global OIIO attribute "nthreads".  If nthreads == 1, it
// guarantees that it will not launch any new threads.
func Cut(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.cut(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

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

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Copy into dst, beginning at (xbegin, ybegin, zbegin, chbegin), the pixels of src described by roi.
// If roi is nil, the entirety of src will be used.
func Paste(dst, src *ImageBuf, xbegin, ybegin, zbegin, chbegin int, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.paste(dst.ptr, C.int(xbegin), C.int(ybegin), C.int(zbegin), C.int(chbegin),
		src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Flip copies src (or a subregion of src) to the corresponding pixels of dst, but with the scanlines
// exchanged vertically.
func Flip(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.flip(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Flop copies src (or a subregion of src) to the corresponding pixels of dst, but with the columns
// exchanged horizontally.
func Flop(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.flop(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Flipflop copies src (or a subregion of src to the corresponding pixels of dst, but with both the
// rows exchanged vertically and the columns exchanged horizontally (this is equivalent to
// a 180 degree rotation).
func Flipflop(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.flipflop(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Transpose copies src (or a subregion of src to the corresponding transposed (x$y) pixels
func Transpose(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.transpose(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels and channels within the designated region, set dst to the sum of image A
// and image B. All of the images must have the same number of channels.
func Add(dst, a, b *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.add(dst.ptr, a.ptr, b.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels and channels within the designated region, set dst to the sum of image src
// and float value (added to all channels). All of the images must have the same number of channels.
func AddValue(dst, src *ImageBuf, value float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.add_value(dst.ptr, src.ptr, C.float(value), opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels and channels within the designated region, set dst to the sum of image src
// and per-channel float slice values. All of the images must have the same number of channels.
func AddValues(dst, src *ImageBuf, values []float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_ptr := (*C.float)(unsafe.Pointer(&values[0]))

	ok := C.add_values(dst.ptr, src.ptr, c_ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(values)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels and channels within the designated region, subtract image B from
// image A. All of the images must have the same number of channels.
func Sub(dst, a, b *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.sub(dst.ptr, a.ptr, b.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels and channels within the designated region, subtract float value (subtracted from
// all channels) from image src. All of the images must have the same number of channels.
func SubValue(dst, src *ImageBuf, value float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.sub_value(dst.ptr, src.ptr, C.float(value), opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels and channels within the designated region, subtract per-channel float slice B
// from image src. All of the images must have the same number of channels.
func SubValues(dst, src *ImageBuf, values []float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_ptr := (*C.float)(unsafe.Pointer(&values[0]))

	ok := C.sub_values(dst.ptr, src.ptr, c_ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(values)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels within the designated region, multiply the pixel values of image A by image B
// (channel by channel), putting the product in dst. All of the images must have the same number
// of channels.
func Mul(dst, a, b *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.mul(dst.ptr, a.ptr, b.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels within the designated region, multiply the pixel values of image A by float
// value B (applied to all channels), putting the product in dst. All of the images must have
// the same number of channels.
func MulValue(dst, src *ImageBuf, value float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.mul_value(dst.ptr, src.ptr, C.float(value), opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// For all pixels within the designated region, multiply the pixel values of image A by
// per-channel float array, putting the product in dst. All of the images must have the
// same number of channels.
func MulValues(dst, src *ImageBuf, values []float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_ptr := (*C.float)(unsafe.Pointer(&values[0]))

	ok := C.mul_values(dst.ptr, src.ptr, c_ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(values)
	runtime.KeepAlive(opt)

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

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

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

	runtime.KeepAlive(dst)
	runtime.KeepAlive(cp)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// Premult copies pixels from src to dst, and in the process multiply all color channels (those not
// alpha or z) by the alpha value, to "premultiply" them. This presumes that the image starts
// of as "unassociated alpha" a.k.a. "non-premultipled." The alterations are restricted to the
// pixels and channels of the supplied ROI (which defaults to all of src). This is just a copy
// if there is no identified alpha channel (and a no-op if dst and src are the same image).
func Premult(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.premult(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Unpremult copies pixels from src to dst, and in the process divide all color channels (those not alpha
// or z) by the alpha value, to "un-premultiply" them. This presumes that the image starts
// of as "associated alpha" a.k.a. "premultipled." The alterations are restricted to the pixels
// and channels of the supplied ROI (which defaults to all of src). Pixels in which the alpha
// channel is 0 will not be modified (since the operation is undefined in that case). This is
// just a copy if there is no identified alpha channel (and a no-op if dst and src are the same
// image).
func Unpremult(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	ok := C.unpremult(dst.ptr, src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// IsConstantColor returns true if all pixels of src within the ROI have the same values
// (for the subset of channels described by roi)
func IsConstantColor(src *ImageBuf, opts ...AlgoOpts) bool {
	opt := flatAlgoOpts(opts)
	ok := C.is_constant_color(src.ptr, nil, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	return bool(ok)
}

// ConstantColors returns a slice of the constant pixel values of all channels within
// either the src image, or the specific channel range of the given ROI.
// If the image does not have constant colors in the given channel range, a nil value
// is returned.
func ConstantColors(src *ImageBuf, opts ...AlgoOpts) []float32 {
	opt := flatAlgoOpts(opts)

	// Determine how big of a slice to allocation,
	// depending on whether they have passed an ROI
	// or not.
	num := src.NumChannels()
	if opt.ROI != nil {
		roi_num := opt.ROI.NumChannels()
		if roi_num < num {
			num = roi_num
		}
	}
	values := make([]float32, num)
	c_ptr := (*C.float)(unsafe.Pointer(&values[0]))

	ok := C.is_constant_color(src.ptr, c_ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(src)
	runtime.KeepAlive(values)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return nil
	}
	return values
}

// IsConstantChannel returns true if all pixels of src within the
// ROI have the given channel value val
func IsConstantChannel(src *ImageBuf, channel int, val float32, opts ...AlgoOpts) bool {
	opt := flatAlgoOpts(opts)
	ok := C.is_constant_channel(src.ptr, C.int(channel), C.float(val),
		opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	return bool(ok)
}

// IsMonochrome returns true if the image is monochrome within the ROI,
// that is, for all pixels within the region, do all channels
// [roi.chbegin, roi.chend) have the same value? If roi is not defined
// (the default), it will be understood to be all of the defined pixels
// and channels of source.
func IsMonochrome(src *ImageBuf, opts ...AlgoOpts) bool {
	opt := flatAlgoOpts(opts)
	ok := C.is_monochrome(src.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	return bool(ok)
}

// ComputePixelHashSHA1 computes the SHA-1 byte hash for all the pixels in the specifed region of the image. If
// blocksize > 0, the function will compute separate SHA-1 hashes of each blocksize
// batch of scanlines, then return a hash of the individual hashes. This is just as strong a
// hash, but will NOT match a single hash of the entire image (blocksize == 0). But by
// breaking up the hash into independent blocks, we can parallelize across multiple threads,
// given by nthreads. The extrainfo provides additional text that will be incorporated
// into the hash.
func ComputePixelHashSHA1(src *ImageBuf, extraInfo string, blockSize int, opts ...AlgoOpts) string {
	opt := flatAlgoOpts(opts)

	cExtraInfo := C.CString(extraInfo)
	defer C.free(unsafe.Pointer(cExtraInfo))

	if blockSize < 0 {
		blockSize = 0
	}

	cStr := C.computePixelHashSHA1(src.ptr, cExtraInfo,
		opt.ROI.validOrAllPtr(), C.int(blockSize), C.int(opt.Threads))

	str := C.GoString(cStr)
	C.free(unsafe.Pointer(cStr))

	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	return str
}

// Set dst, over the region of interest, to be a resized version of the
// corresponding portion of src (mapping such that the "full" image
// window of each correspond to each other, regardless of resolution).
//
// Default high-quality filters are used: blackman-harris when upsizing,
// lanczos3 when downsizing
//
// The nthreads AlgoOpts specifies how many threads (potentially) may
// be used, but it's not a guarantee.  If nthreads == 0, it will use
// the global OIIO attribute "nthreads".  If nthreads == 1, it
// guarantees that it will not launch any new threads.
func Resize(dst, src *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_filtname := C.CString("")
	defer C.free(unsafe.Pointer(c_filtname))

	ok := C.resize(dst.ptr, src.ptr, c_filtname, C.float(0.0), opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}

	return nil
}

// Set dst, over the region of interest, to be a resized version of the
// corresponding portion of src (mapping such that the "full" image
// window of each correspond to each other, regardless of resolution).
//
// The filter is used to weight the src pixels falling underneath it for
// each dst pixel.  The caller may specify a reconstruction filter by name
// and width (expressed  in pixels unts of the dst image), or ResizeFilter() will
// choose a reasonable default high-quality default filter (blackman-harris
// when upsizing, lanczos3 when downsizing) if the empty string is passed
// or if filterWidth is 0.
//
// The nthreads AlgoOpts specifies how many threads (potentially) may
// be used, but it's not a guarantee.  If nthreads == 0, it will use
// the global OIIO attribute "nthreads".  If nthreads == 1, it
// guarantees that it will not launch any new threads.
func ResizeFilter(dst, src *ImageBuf, filter string, filterWidth float32, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)

	c_filtname := C.CString(filter)
	defer C.free(unsafe.Pointer(c_filtname))

	ok := C.resize(dst.ptr, src.ptr, c_filtname, C.float(filterWidth), opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

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

	runtime.KeepAlive(dst)
	runtime.KeepAlive(src)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// Over sets dst to the composite of A over B using the Porter/Duff definition
// of "over", returning true upon success and false for any of a
// variety of failures (as described below).
//
// A and B (and dst, if already defined/allocated) must have valid alpha
// channels identified by their ImageSpec alpha_channel field.  If A or
// B do not have alpha channels (as determined by those rules) or if
// the number of non-alpha channels do not match between A and B,
// Over() will fail.
//
// If dst is not already an initialized ImageBuf, it will be sized to
// encompass the minimal rectangular pixel region containing the union
// of the defined pixels of A and B, and with a number of channels
// equal to the number of non-alpha channels of A and B, plus an alpha
// channel.  However, if dst is already initialized, it will not be
// resized, and the "over" operation will apply to its existing pixel
// data window.  In this case, dst must have an alpha channel designated
// and must have the same number of non-alpha channels as A and B,
// otherwise it will fail.
//
// 'roi' AlgoOpts specifies the region of dst's pixels which will be computed;
// existing pixels outside this range will not be altered.  If not
// specified, the default ROI value will be interpreted as a request to
// apply "A over B" to the entire region of dst's pixel data.
//
// A, B, and dst need not perfectly overlap in their pixel data windows;
// pixel values of A or B that are outside their respective pixel data
// window will be treated as having "zero" (0,0,0...) value.
//
// The nthreads AlgoOpts specifies how many threads (potentially) may
// be used, but it's not a guarantee.  If nthreads == 0, it will use
// the global OIIO attribute "nthreads".  If nthreads == 1, it
// guarantees that it will not launch any new threads.
func Over(dst, a, b *ImageBuf, opts ...AlgoOpts) error {
	opt := flatAlgoOpts(opts)
	ok := C.over(dst.ptr, a.ptr, b.ptr, opt.ROI.validOrAllPtr(), C.int(opt.Threads))

	runtime.KeepAlive(dst)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(opt)

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// RenderTextColor renders a text string into image dst, essentially doing an "over" of
// the character into the existing pixel data.
// The baseline of the first character will start at position (x,y), and with a
// nominal height of fontSize (in pixels).
// The font is given by fontName as either a full pathname to the font file,
// or a font basename which can be found in the standard system font locations.
// If an empty string is provided or font is not found, it defaults to some
// reasonable system font if not supplied at all.
// The characters will be drawn in opaque white (1.0,1.0,...) in all channels,
// unless color is supplied (and is expected to point to a float slice of length at
// least equal to R.Spec().NumChannels).
func RenderTextColor(dst *ImageBuf, x, y int, text string, fontSize int, fontName string, color []float32) error {
	c_text := C.CString(text)
	defer C.free(unsafe.Pointer(c_text))

	c_fontName := C.CString(fontName)
	defer C.free(unsafe.Pointer(c_fontName))

	var color_ptr *C.float
	if color != nil && len(color) > 0 {
		color_ptr = (*C.float)(unsafe.Pointer(&color[0]))
	}

	ok := C.render_text(dst.ptr, C.int(x), C.int(y), c_text, C.int(fontSize), c_fontName, color_ptr)

	runtime.KeepAlive(dst)
	runtime.KeepAlive(color)

	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}
