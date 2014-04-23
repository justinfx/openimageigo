package oiio

/*
#include "stdlib.h"

#include "cpp/oiio.h"

*/
import "C"

import (
	"errors"
	"reflect"
	"runtime"
	"unsafe"
)

// ImageInput abstracts the reading of an image file in a file format-agnostic manner.
type ImageInput struct {
	ptr unsafe.Pointer
}

func newImageInput(i unsafe.Pointer) *ImageInput {
	in := &ImageInput{i}
	runtime.SetFinalizer(in, deleteImageInput)
	return in
}

func deleteImageInput(i *ImageInput) {
	if i.ptr != nil {
		C.free(i.ptr)
		i.ptr = nil
	}
}

// Create an ImageInput subclass instance that is able to read the given file and open it,
// returning the opened ImageInput if successful. If it fails, return error.
func OpenImageInput(filename string) (*ImageInput, error) {
	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	cfg := unsafe.Pointer(nil)
	ptr := C.ImageInput_Open(c_str, cfg)

	in := newImageInput(ptr)

	return in, in.LastError()
}

// Return the last error generated by API calls.
// An nil error will be returned if no error has occured.
func (i *ImageInput) LastError() error {
	err := C.GoString(C.ImageInput_geterror(i.ptr))
	if err == "" {
		return nil
	}
	return errors.New(err)
}

// Open file with given name. Return true if the file was found and opened okay.
func (i *ImageInput) Open(filename string) error {
	i.Close()

	deleteImageInput(i)

	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	cfg := unsafe.Pointer(nil)
	ptr := C.ImageInput_Open(c_str, cfg)
	i.ptr = ptr

	return i.LastError()
}

// Close an image that we are totally done with.
func (i *ImageInput) Close() error {
	if !bool(C.ImageInput_close(i.ptr)) {
		return i.LastError()
	}
	return nil
}

// Return the name of the format implemented by this image.
func (i *ImageInput) FormatName() string {
	return C.GoString(C.ImageInput_format_name(i.ptr))
}

// Return true if the named file is file of the type for this ImageInput.
// The implementation will try to determine this as efficiently as possible,
// in most cases much less expensively than doing a full Open().
// Note that a file can appear to be of the right type (i.e., ValidFIle() returning true)
// but still fail a subsequent call to Open(), such as if the contents of the file are
// truncated, nonsensical, or otherwise corrupted.
func (i *ImageInput) ValidFile(filename string) bool {
	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))
	return bool(C.ImageInput_valid_file(i.ptr, c_str))
}

// Given the name of a 'feature', return whether this ImageInput supports input of images
// with the given properties. Feature names that ImageIO plugins are expected to recognize
// include: none at this time.
// Note that main advantage of this approach, versus having separate individual supports_foo()
// methods, is that this allows future expansion of the set of possible queries without changing
// the API, adding new entry points, or breaking linkage compatibility.
func (i *ImageInput) Supports(feature string) bool {
	c_str := C.CString(feature)
	defer C.free(unsafe.Pointer(c_str))
	return bool(C.ImageInput_supports(i.ptr, c_str))
}

// Return a reference to the image format specification of the current subimage/MIPlevel.
// Note that the contents of the spec are invalid before Open() or after Close(), and may
// change with a call to SeekSubImage().
func (i *ImageInput) Spec() (*ImageSpec, error) {
	ptr := C.ImageInput_spec(i.ptr)
	err := i.LastError()
	if err != nil {
		return nil, err
	}
	return &ImageSpec{ptr}, nil
}

// Read the entire image of width * height * depth * channels into contiguous float32 pixels.
// Read tiles or scanlines automatically.
func (i *ImageInput) ReadImage() ([]float32, error) {
	spec, err := i.Spec()
	if err != nil {
		return nil, err
	}

	size := spec.Width() * spec.Height() * spec.Depth() * spec.NumChannels()
	pixels := make([]float32, size)
	pixels_ptr := (*C.float)(unsafe.Pointer(&pixels[0]))
	C.ImageInput_read_image_floats(i.ptr, pixels_ptr)

	return pixels, i.LastError()
}

// Read the entire image of width * height * depth * channels into contiguous pixels.
// Read tiles or scanlines automatically.
//
// This call supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
// The underlying type of data is determined by the given TypeDesc.
// Returned interface{} will be:
//     TypeUint8   => []uint8
//     TypeInt8    => []int8
//     TypeUint16  => []uint16
//     TypeInt16   => []int16
//     TypeUint    => []uint
//     TypeInt     => []int
//     TypeUint64  => []uint64
//     TypeInt64   => []int64
//     TypeHalf    => []float32
//     TypeFloat   => []float32
//     TypeDouble  => []float64
//
// Example:
//
//     // Without a callback
//     val, err := in.ReadImageFormat(TypeFloat, nil)
//     if err != nil {
//         panic(err.Error())
//     }
//     floatPixels := val.([]float32)
//
//     // With a callback
//     var cbk ProgressCallback = func(done float32) bool {
//         fmt.Printf("Progress: %0.2f\n", done)
//         // Keep processing (return true to abort)
//         return false
//     }
//     val, _ = in.ReadImageFormat(TypeFloat, &cbk)
//     floatPixels = val.([]float32)
//
func (i *ImageInput) ReadImageFormat(format TypeDesc, progress *ProgressCallback) (interface{}, error) {
	spec, err := i.Spec()
	if err != nil {
		return nil, err
	}

	size := spec.Width() * spec.Height() * spec.Depth() * spec.NumChannels()

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
		return nil, errors.New("TypeDesc is not valid for this operation")

	}

	var cbk unsafe.Pointer
	if progress != nil {
		cbk = unsafe.Pointer(progress)
	}

	C.ImageInput_read_image_format(i.ptr, (C.TypeDesc)(format), ptr, cbk)

	return pixel_iface, i.LastError()
}

// Read the scanline that includes pixels (*,y,z) into data, converting if necessary
// from the native data format of the file into contiguous float32 pixels (z==0 for non-volume images).
// The size of the slice is: width * depth * channels
func (i *ImageInput) ReadScanline(y, z int) ([]float32, error) {
	spec, err := i.Spec()
	if err != nil {
		return nil, err
	}

	size := spec.Width() * spec.Depth() * spec.NumChannels()
	pixels := make([]float32, size)
	pixels_ptr := (*C.float)(unsafe.Pointer(&pixels[0]))
	C.ImageInput_read_scanline_floats(i.ptr, C.int(y), C.int(z), pixels_ptr)

	return pixels, i.LastError()
}