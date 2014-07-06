package oiio

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var supportedSliceKinds = []reflect.Kind{
	reflect.Uint8,
	reflect.Int8,
	reflect.Uint16,
	reflect.Int16,
	reflect.Int,
	reflect.Uint,
	reflect.Uint64,
	reflect.Int64,
	reflect.Float32,
	reflect.Float64,
}

// Given an ImageSpec, a slice is allocated to the size that
// is able to contain the pixels of the ImageSpec dimensions.
// The TypeDesc determines what format the pixels will be stored in.
// Returns the slice, casted to an interface.
func allocatePixelBuffer(spec *ImageSpec, format TypeDesc) (interface{}, unsafe.Pointer, error) {
	size := spec.Width() * spec.Height() * spec.Depth() * spec.NumChannels()
	return allocatePixelBufferSize(size, format)
}

// Given a specific size, a slice is allocated that
// is able to contain the number of pixels.
// The TypeDesc determines what format the pixels will be stored in.
// Returns the slice, casted to an interface.
func allocatePixelBufferSize(size int, format TypeDesc) (interface{}, unsafe.Pointer, error) {
	var (
		pixel_iface interface{}
		ptr         unsafe.Pointer
	)

	switch format {

	case TypeUint8:
		pixels := make([]uint8, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeInt8:
		pixels := make([]int8, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeUint16:
		pixels := make([]uint16, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeInt16:
		pixels := make([]int16, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeUint:
		pixels := make([]uint, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeInt:
		pixels := make([]int, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeUint64:
		pixels := make([]uint64, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeInt64:
		pixels := make([]int64, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeFloat, TypeHalf:
		pixels := make([]float32, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	case TypeDouble:
		pixels := make([]float64, size)
		pixel_iface = reflect.ValueOf(pixels).Interface()
		ptr = unsafe.Pointer(&pixels[0])

	default:
		return nil, nil, errors.New("TypeDesc is not valid for this operation")

	}

	return pixel_iface, ptr, nil
}

// Takes a slice, passed in generically as an interface, and returns
// the underlying pointer to the data.
// The slice type must be one of the supported slice types to contain pixel data.
func pixelsToPtr(slice interface{}) (unsafe.Pointer, error) {
	slice_t := reflect.TypeOf(slice)
	if slice_t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Not a slice. Received type %s", slice_t.Kind().String())
	}

	elem := slice_t.Elem().Kind()
	ok := false
	for _, okElem := range supportedSliceKinds {
		ok = elem == okElem
		if ok {
			break
		}
	}

	if !ok {
		return nil, fmt.Errorf("Slice type is not one of the supported pixel data types %q", supportedSliceKinds)
	}

	val := reflect.ValueOf(slice)
	ptr := val.Pointer()

	if ptr == 0 {
		return unsafe.Pointer(nil), nil
	}

	return unsafe.Pointer(ptr), nil
}
