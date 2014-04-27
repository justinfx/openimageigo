package oiio

/*
#include "stdlib.h"

#include "cpp/oiio.h"

*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Helper struct describing a region of interest in an image.
// The region is [xbegin,xend) x [begin,yend) x [zbegin,zend), with the "end" designators
// signifying one past the last pixel, a la C++ STL style.
type ROI struct {
	ptr unsafe.Pointer
}

func newROI(i unsafe.Pointer) *ROI {
	in := &ROI{i}
	runtime.SetFinalizer(in, deleteROI)
	return in
}

func deleteROI(i *ROI) {
	if i.ptr != nil {
		C.free(i.ptr)
		i.ptr = nil
	}
}

// Default constructor is an undefined region.
func NewROI() *ROI {
	return newROI(C.ROI_New())
}

// Constructor with an explicitly defined region, where you are
// concerned with just the X/Y region, and not the Z or the channels
func NewROIRegion2D(xbegin, xend, ybegin, yend int) *ROI {
	ptr := C.ROI_NewOptions(
		C.int(xbegin),
		C.int(xend),
		C.int(ybegin),
		C.int(yend),
		C.int(0),
		C.int(1),
		C.int(0),
		C.int(1000),
	)
	return newROI(ptr)
}

// Constructor with an explicitly defined region.
// Reasonable default values are:
//   zbegin  = 0
//   zend    = 1
//   chbegin = 0
//   chend   = 1000
func NewROIRegion3D(xbegin, xend, ybegin, yend, zbegin, zend, chbegin, chend int) *ROI {
	ptr := C.ROI_NewOptions(
		C.int(xbegin),
		C.int(xend),
		C.int(ybegin),
		C.int(yend),
		C.int(zbegin),
		C.int(zend),
		C.int(chbegin),
		C.int(chend),
	)
	return newROI(ptr)
}

// Is a region defined?
func (r *ROI) Defined() bool {
	return bool(C.ROI_defined(r.ptr))
}

// Width of the region (X)
func (r *ROI) Width() int {
	return int(C.ROI_width(r.ptr))
}

// Height of the region (Y)
func (r *ROI) Height() int {
	return int(C.ROI_height(r.ptr))
}

// Depth of the region (Z)
func (r *ROI) Depth() int {
	return int(C.ROI_depth(r.ptr))
}

// Number of channels in the region
func (r *ROI) NumChannels() int {
	return int(C.ROI_nchannels(r.ptr))
}

// Number of total pixels in the region
// This is Width * Height * Depth
func (r *ROI) NumPixels() int {
	return int(C.ROI_npixels(r.ptr))
}

func (r *ROI) XBegin() int {
	return int(C.ROI_xbegin(r.ptr))
}

func (r *ROI) SetXBegin(x int) {
	C.ROI_set_xbegin(r.ptr, C.int(x))
}

func (r *ROI) XEnd() int {
	return int(C.ROI_xend(r.ptr))
}

func (r *ROI) SetXEnd(x int) {
	C.ROI_set_xend(r.ptr, C.int(x))
}

func (r *ROI) YBegin() int {
	return int(C.ROI_ybegin(r.ptr))
}

func (r *ROI) SetYBegin(y int) {
	C.ROI_set_ybegin(r.ptr, C.int(y))
}

func (r *ROI) YEnd() int {
	return int(C.ROI_yend(r.ptr))
}

func (r *ROI) SetYEnd(y int) {
	C.ROI_set_yend(r.ptr, C.int(y))
}

func (r *ROI) ZBegin() int {
	return int(C.ROI_zbegin(r.ptr))
}

func (r *ROI) SetZBegin(z int) {
	C.ROI_set_zbegin(r.ptr, C.int(z))
}

func (r *ROI) ZEnd() int {
	return int(C.ROI_zend(r.ptr))
}

func (r *ROI) SetZEnd(z int) {
	C.ROI_set_zend(r.ptr, C.int(z))
}

func (r *ROI) ChannelsBegin() int {
	return int(C.ROI_chbegin(r.ptr))
}

func (r *ROI) SetChannelsBegin(ch int) {
	C.ROI_set_chbegin(r.ptr, C.int(ch))
}

func (r *ROI) ChannelsEnd() int {
	return int(C.ROI_chend(r.ptr))
}

func (r *ROI) SetChannelsEnd(ch int) {
	C.ROI_set_chend(r.ptr, C.int(ch))
}
