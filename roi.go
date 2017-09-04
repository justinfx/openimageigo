package oiio

/*
#include "stdlib.h"

#include "cpp/oiio.h"

*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Helper struct describing a region of interest in an image.
// The region is [xbegin,xend) x [begin,yend) x [zbegin,zend), with the "end" designators
// signifying one past the last pixel, a la C++ STL style.
type ROI struct {
	ptr unsafe.Pointer
}

var roi_all *ROI

func init() {
	roi_all = NewROI()
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
	runtime.KeepAlive(i)
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

func (r *ROI) validOrAllPtr() unsafe.Pointer {
	if r == nil || r.ptr == nil {
		return roi_all.ptr
	}
	return r.ptr
}

// Return a new copy of the ROI that can be freely modified.
func (r *ROI) Copy() *ROI {
	rc := C.ROI_Copy(r.ptr)
	runtime.KeepAlive(r)
	return newROI(rc)
}

// String returns a printable string representation
// of the ROI, containing just the origin (X,Y) and Width,Height.
func (r *ROI) String() string {
	return fmt.Sprintf("ROI:{X: %d, Y: %d, W: %d, H: %d, ...}",
		r.XBegin(), r.YBegin(), r.Width(), r.Height())
}

// Is a region defined?
func (r *ROI) Defined() bool {
	ret := bool(C.ROI_defined(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

// Width of the region (X)
func (r *ROI) Width() int {
	ret := int(C.ROI_width(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

// Height of the region (Y)
func (r *ROI) Height() int {
	ret := int(C.ROI_height(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

// Depth of the region (Z)
func (r *ROI) Depth() int {
	ret := int(C.ROI_depth(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

// Number of channels in the region
func (r *ROI) NumChannels() int {
	ret := int(C.ROI_nchannels(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

// Number of total pixels in the region
// This is Width * Height * Depth
func (r *ROI) NumPixels() int {
	ret := int(C.ROI_npixels(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) XBegin() int {
	ret := int(C.ROI_xbegin(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetXBegin(x int) {
	C.ROI_set_xbegin(r.ptr, C.int(x))
	runtime.KeepAlive(r)
}

func (r *ROI) XEnd() int {
	ret := int(C.ROI_xend(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetXEnd(x int) {
	C.ROI_set_xend(r.ptr, C.int(x))
	runtime.KeepAlive(r)
}

func (r *ROI) YBegin() int {
	ret := int(C.ROI_ybegin(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetYBegin(y int) {
	C.ROI_set_ybegin(r.ptr, C.int(y))
	runtime.KeepAlive(r)
}

func (r *ROI) YEnd() int {
	ret := int(C.ROI_yend(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetYEnd(y int) {
	C.ROI_set_yend(r.ptr, C.int(y))
	runtime.KeepAlive(r)
}

func (r *ROI) ZBegin() int {
	ret := int(C.ROI_zbegin(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetZBegin(z int) {
	C.ROI_set_zbegin(r.ptr, C.int(z))
	runtime.KeepAlive(r)
}

func (r *ROI) ZEnd() int {
	ret := int(C.ROI_zend(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetZEnd(z int) {
	C.ROI_set_zend(r.ptr, C.int(z))
	runtime.KeepAlive(r)
}

func (r *ROI) ChannelsBegin() int {
	ret := int(C.ROI_chbegin(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetChannelsBegin(ch int) {
	C.ROI_set_chbegin(r.ptr, C.int(ch))
	runtime.KeepAlive(r)
}

func (r *ROI) ChannelsEnd() int {
	ret := int(C.ROI_chend(r.ptr))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetChannelsEnd(ch int) {
	C.ROI_set_chend(r.ptr, C.int(ch))
	runtime.KeepAlive(r)
}
