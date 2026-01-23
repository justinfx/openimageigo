package oiio

/*
#include "stdlib.h"

#include "oiio.h"

*/
import "C"

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// ROI helper struct describing a region of interest in an image.
// The region is [xbegin,xend) x [begin,yend) x [zbegin,zend), with the "end" designators
// signifying one past the last pixel, a la C++ STL style.
type ROI struct {
	ptr *unsafe.Pointer
}

var roi_all *ROI

func init() {
	roi_all = NewROI()
}

func newROI(i unsafe.Pointer) *ROI {
	holder := &i
	in := &ROI{ptr: holder}
	runtime.AddCleanup(in, cleanupROI, holder)
	return in
}

func cleanupROI(ptr *unsafe.Pointer) {
	p := atomic.SwapPointer(ptr, nil)
	if p != nil {
		C.deleteROI(p)
	}
}

// NewROI default constructor is an undefined region.
func NewROI() *ROI {
	return newROI(C.ROI_New())
}

// NewROIRegion2D constructor with an explicitly defined region, where you are.
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

// NewROIRegion3D constructor with an explicitly defined region.
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

// Destroy the object immediately instead of waiting for GC.
func (r *ROI) Destroy() {
	cleanupROI(r.ptr)
	runtime.KeepAlive(r)
}

func (r *ROI) p() unsafe.Pointer {
	return atomic.LoadPointer(r.ptr)
}

func (r *ROI) validOrAllPtr() unsafe.Pointer {
	if r == nil || r.ptr == nil {
		return atomic.LoadPointer(roi_all.ptr)
	}
	return r.p()
}

// Copy returns a new copy of the ROI that can be freely modified..
func (r *ROI) Copy() *ROI {
	rc := C.ROI_Copy(r.p())
	runtime.KeepAlive(r)
	return newROI(rc)
}

// String returns a printable string representation
// of the ROI, containing just the origin (X,Y) and Width,Height.
func (r *ROI) String() string {
	return fmt.Sprintf("ROI:{X: %d, Y: %d, W: %d, H: %d, ...}",
		r.XBegin(), r.YBegin(), r.Width(), r.Height())
}

// Defined returns true if a region defined?
func (r *ROI) Defined() bool {
	ret := bool(C.ROI_defined(r.p()))
	runtime.KeepAlive(r)
	return ret
}

// Width of the region (X)
func (r *ROI) Width() int {
	ret := int(C.ROI_width(r.p()))
	runtime.KeepAlive(r)
	return ret
}

// Height of the region (Y)
func (r *ROI) Height() int {
	ret := int(C.ROI_height(r.p()))
	runtime.KeepAlive(r)
	return ret
}

// Depth of the region (Z)
func (r *ROI) Depth() int {
	ret := int(C.ROI_depth(r.p()))
	runtime.KeepAlive(r)
	return ret
}

// NumChannels number of channels in the region.
func (r *ROI) NumChannels() int {
	ret := int(C.ROI_nchannels(r.p()))
	runtime.KeepAlive(r)
	return ret
}

// NumPixels number of total pixels in the region.
// This is Width * Height * Depth
func (r *ROI) NumPixels() int {
	ret := int(C.ROI_npixels(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) XBegin() int {
	ret := int(C.ROI_xbegin(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetXBegin(x int) {
	C.ROI_set_xbegin(r.p(), C.int(x))
	runtime.KeepAlive(r)
}

func (r *ROI) XEnd() int {
	ret := int(C.ROI_xend(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetXEnd(x int) {
	C.ROI_set_xend(r.p(), C.int(x))
	runtime.KeepAlive(r)
}

func (r *ROI) YBegin() int {
	ret := int(C.ROI_ybegin(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetYBegin(y int) {
	C.ROI_set_ybegin(r.p(), C.int(y))
	runtime.KeepAlive(r)
}

func (r *ROI) YEnd() int {
	ret := int(C.ROI_yend(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetYEnd(y int) {
	C.ROI_set_yend(r.p(), C.int(y))
	runtime.KeepAlive(r)
}

func (r *ROI) ZBegin() int {
	ret := int(C.ROI_zbegin(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetZBegin(z int) {
	C.ROI_set_zbegin(r.p(), C.int(z))
	runtime.KeepAlive(r)
}

func (r *ROI) ZEnd() int {
	ret := int(C.ROI_zend(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetZEnd(z int) {
	C.ROI_set_zend(r.p(), C.int(z))
	runtime.KeepAlive(r)
}

func (r *ROI) ChannelsBegin() int {
	ret := int(C.ROI_chbegin(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetChannelsBegin(ch int) {
	C.ROI_set_chbegin(r.p(), C.int(ch))
	runtime.KeepAlive(r)
}

func (r *ROI) ChannelsEnd() int {
	ret := int(C.ROI_chend(r.p()))
	runtime.KeepAlive(r)
	return ret
}

func (r *ROI) SetChannelsEnd(ch int) {
	C.ROI_set_chend(r.p(), C.int(ch))
	runtime.KeepAlive(r)
}
