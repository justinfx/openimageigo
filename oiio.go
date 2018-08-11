/*

OpenImageIO bindings

https://sites.google.com/site/openimageio/home

OpenImageIO is a library for reading and writing images, and a bunch of related classes,
utilities, and applications.  There is a particular emphasis on formats and functionality
used in professional, large-scale animation and visual effects work for film.
OpenImageIO is used extensively in animation and VFX studios all over the world, and is
also incorporated into several commercial products.

*/
package oiio

/*
#cgo CPPFLAGS: -I/usr/local/include
#cgo LDFLAGS: -lstdc++

#include "stdlib.h"

#include "oiio.h"

*/
import "C"

import (
	"unsafe"
)

// Various representation formats for image data
type TypeDesc int

const (
	TypeUnknown TypeDesc = C.TYPE_UNKNOWN
	TypeUint8   TypeDesc = C.TYPE_UINT8
	TypeInt8    TypeDesc = C.TYPE_INT8
	TypeUint16  TypeDesc = C.TYPE_UINT16
	TypeInt16   TypeDesc = C.TYPE_INT16
	TypeUint    TypeDesc = C.TYPE_UINT
	TypeInt     TypeDesc = C.TYPE_INT
	TypeUint64  TypeDesc = C.TYPE_UINT64
	TypeInt64   TypeDesc = C.TYPE_INT64
	TypeHalf    TypeDesc = C.TYPE_HALF
	TypeFloat   TypeDesc = C.TYPE_FLOAT
	TypeDouble  TypeDesc = C.TYPE_DOUBLE
)

// For image processing functions that accept a callback to monitor progress.
// A function that will be passed a float value indicating the progress
// percentage of the current operation. If the functon returns true, then
// the process should be aborted. Return false to allow processing to continue.
type ProgressCallback func(done float32) bool

//export image_progress_callback
func image_progress_callback(goCallback unsafe.Pointer, done C.float) C.bool {
	if goCallback == nil {
		return C.bool(false)
	}

	fn := *(*ProgressCallback)(goCallback)
	cancel := fn(float32(done))
	return C.bool(cancel)
}
