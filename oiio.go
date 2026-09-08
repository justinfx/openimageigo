/*
Package oiio provides OpenImageIO bindings for Go.

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

// TypeDesc various representation formats for image data.
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

// ProgressCallback for image processing functions that accept a callback to monitor progress.
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

// SetAttribute sets a global OpenImageIO attribute such as "threads",
// "exr_threads", or "openexr:core". Supported value types are int, float32, and string.
// Returns true if the attribute was successfully set.
func SetAttribute(name string, val interface{}) bool {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var ok bool
	switch t := val.(type) {
	case string:
		c_val := C.CString(t)
		defer C.free(unsafe.Pointer(c_val))
		ok = bool(C.OIIO_attribute_string(c_name, c_val))
	case int:
		ok = bool(C.OIIO_attribute_int(c_name, C.int(t)))
	case float32:
		ok = bool(C.OIIO_attribute_float(c_name, C.float(t)))
	default:
		ok = false
	}
	return ok
}

// AttributeInt returns a global OpenImageIO int attribute, or defaultVal (or 0
// if not specified) if the attribute is not found.
func AttributeInt(name string, defaultVal ...int) int {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var val C.int
	ok := bool(C.OIIO_getattribute_int(c_name, &val))

	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return int(val)
}

// AttributeFloat returns a global OpenImageIO float attribute, or defaultVal (or 0
// if not specified) if the attribute is not found.
func AttributeFloat(name string, defaultVal ...float32) float32 {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var val C.float
	ok := bool(C.OIIO_getattribute_float(c_name, &val))

	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return float32(val)
}

// AttributeString returns a global OpenImageIO string attribute, or defaultVal (or ""
// if not specified) if the attribute is not found.
func AttributeString(name string, defaultVal ...string) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var c_val *C.char
	ok := bool(C.OIIO_getattribute_string(c_name, &c_val))

	if !ok || c_val == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	val := C.GoString(c_val)
	C.free(unsafe.Pointer(c_val))
	return val
}
