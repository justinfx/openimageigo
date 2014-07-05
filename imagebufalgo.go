package oiio

/*
#cgo LDFLAGS: -lboost_thread -lboost_system

#include "stdlib.h"

#include "cpp/imagebufalgo.h"

*/
import "C"

import (
	"errors"
	"unsafe"
)

const (
	//
	// Let OIIO choose the best filter
	FilterDefault = ""
	// Let OIIO choose the best filter and filter width
	FilterDefaultWidth = 0.0
	// Use the global OIIO-determined thread count
	GlobalThreads = 0
)

func checkBufAndROI(dst *ImageBuf, roi *ROI) error {
	if dst.Initialized() || (roi != nil && roi.ptr != nil && roi.Defined()) {
		return nil
	}
	return errors.New("ImageBuf and ROI cannot both be undefined." +
		"ImageBufAlgo without any guess about region of interest")
}

// Zero out (set to 0, black) the image region.
// Only the pixels (and channels) in dst that are specified by roi will be altered;
// the default roi is to alter all the pixels in dst.
// If dst is uninitialized, it will be resized to be a float ImageBuf large enough to
// hold the region specified by roi. It is an error to pass both an uninitialied dst and an undefined roi.
// The nthreads parameter specifies how many threads (potentially) may be used, but it's not a guarantee.
// If nthreads == 0, it will use the global OIIO attribute "nthreads". If nthreads == 1, it guarantees
// that it will not launch any new threads.
// Works on all pixel data types.
func Zero(dst *ImageBuf, roi *ROI, nthreads int) error {
	err := checkBufAndROI(dst, roi)
	if err != nil {
		return err
	}

	ok := bool(C.zero(dst.ptr, validOrAllROIPtr(roi), C.int(nthreads)))
	if !ok {
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
func ColorConvert(dst, src *ImageBuf, from, to string, unpremult bool, roi *ROI, nthreads int) error {
	c_from := C.CString(from)
	defer C.free(unsafe.Pointer(c_from))

	c_to := C.CString(to)
	defer C.free(unsafe.Pointer(c_to))

	ok := C.colorconvert(dst.ptr, src.ptr, c_from, c_to, C.bool(unpremult), validOrAllROIPtr(roi), C.int(nthreads))
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
func ColorConvertProcessor(dst, src *ImageBuf, cp *ColorProcessor, unpremult bool, roi *ROI, nthreads int) error {
	ok := C.colorconvert_processor(dst.ptr, src.ptr, cp.ptr, C.bool(unpremult), validOrAllROIPtr(roi), C.int(nthreads))
	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}

// Set dst, over the region of interest, to be a resized version of the corresponding portion of src
// (mapping such that the "full" image window of each correspond to each other, regardless of resolution).
// Will choose a reasonable default high-quality default filter (blackman-harris when upsizing, lanczos3 when downsizing)
// The nthreads parameter specifies how many threads (potentially) may be used, but it's not a guarantee.
// If nthreads == 0, it will use the global OIIO attribute "nthreads". If nthreads == 1, it guarantees that it will
// not launch any new threads.
// Works on all pixel data types.
func Resize(dst, src *ImageBuf, roi *ROI, nthreads int) error {
	c_filtname := C.CString("")
	defer C.free(unsafe.Pointer(c_filtname))

	ok := C.resize(dst.ptr, src.ptr, c_filtname, C.float(0.0), validOrAllROIPtr(roi), C.int(nthreads))
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
// The nthreads parameter specifies how many threads (potentially) may be used, but it's not a guarantee.
// If nthreads == 0, it will use the global OIIO attribute "nthreads". If nthreads == 1, it guarantees that it will
// not launch any new threads.
// Works on all pixel data types.
func Resample(dst, src *ImageBuf, interpolate bool, roi *ROI, nthreads int) error {
	ok := C.resample(dst.ptr, src.ptr, C.bool(interpolate), validOrAllROIPtr(roi), C.int(nthreads))
	if !bool(ok) {
		return dst.LastError()
	}
	return nil
}
