package oiio

import (
	"errors"
	"reflect"
	"unsafe"
)

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
