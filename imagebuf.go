package oiio

/*
#include "stdlib.h"

#include "oiio.h"

*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// Description of where the pixels live for this ImageBuf
type IBStorage int

const (
	// Derive the file format from the file path name (empty string)
	FileFormatAuto = ""

	IBStorageLocalBuffer   IBStorage = C.IBSTORAGE_LOCALBUFFER
	IBStorageAppBuffer     IBStorage = C.IBSTORAGE_APPBUFFER
	IBStorageImageCache    IBStorage = C.IBSTORAGE_IMAGECACHE
	IBStorageUninitialized IBStorage = C.IBSTORAGE_UNINITIALIZED
)

// An ImageBuf is a simple in-memory representation of a 2D image.
// It uses ImageInput and ImageOutput underneath for its file I/O, and has simple
// routines for setting and getting individual pixels, that hides most of the details
// of memory layout and data representation (translating to/from float automatically).
type ImageBuf struct {
	ptr unsafe.Pointer
}

func newImageBuf(i unsafe.Pointer) *ImageBuf {
	spec := new(ImageBuf)
	spec.ptr = i
	runtime.SetFinalizer(spec, deleteImageBuf)
	return spec
}

func deleteImageBuf(i *ImageBuf) {
	if i.ptr != nil {
		C.ImageBuf_clear(i.ptr)
		C.deleteImageBuf(i.ptr)
		i.ptr = nil
	}
	runtime.KeepAlive(i)
}

// Return the last error generated by API calls.
// An nil error will be returned if no error has occured.
func (i *ImageBuf) LastError() error {
	c_str := C.ImageBuf_geterror(i.ptr)
	if c_str == nil {
		return nil
	}
	runtime.KeepAlive(i)
	err := C.GoString(c_str)
	C.free(unsafe.Pointer(c_str))
	if err == "" {
		return nil
	}
	return errors.New(err)
}

// Construct an empty/uninitialized ImageBuf. This is relatively useless until you call reset().
func NewImageBuf() *ImageBuf {
	buf := C.ImageBuf_New()
	return newImageBuf(buf)
}

// Construct an ImageBuf to read the named image - but don't actually read it yet!
// The image will actually be read when other methods need to access the spec and/or pixels,
// or when an explicit call to init_spec() or read() is made, whichever comes first.
//
// Uses the global/shared ImageCache.
func NewImageBufPath(path string) (*ImageBuf, error) {
	return NewImageBufPathCache(path, nil)
}

// Construct an ImageBuf to read the named image - but don't actually read it yet!
// The image will actually be read when other methods need to access the spec and/or pixels,
// or when an explicit call to init_spec() or read() is made, whichever comes first.
//
// Uses an explicitly passed ImageCache
func NewImageBufPathCache(path string, cache *ImageCache) (*ImageBuf, error) {
	c_str := C.CString(path)
	defer C.free(unsafe.Pointer(c_str))

	var ptr unsafe.Pointer = nil
	if cache != nil {
		ptr = cache.ptr
	}

	buf := newImageBuf(C.ImageBuf_New_WithCache(c_str, ptr))
	runtime.KeepAlive(cache)
	err := buf.LastError()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Construct an Imagebuf given a proposed spec describing the image size and type,
// and allocate storage for the pixels of the image (whose values will be uninitialized).
func NewImageBufSpec(spec *ImageSpec) (*ImageBuf, error) {
	buf := newImageBuf(C.ImageBuf_New_Spec(spec.ptr))
	runtime.KeepAlive(spec)
	err := buf.LastError()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Is this ImageBuf object initialized?
func (i *ImageBuf) Initialized() bool {
	if i.ptr == nil {
		return false
	}
	ret := bool(C.ImageBuf_initialized(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Destroy the object immediately instead of waiting for GC.
func (i *ImageBuf) Destroy() {
	runtime.SetFinalizer(i, nil)
	deleteImageBuf(i)
}

// Restore the ImageBuf to an uninitialized state.
func (i *ImageBuf) Clear() {
	if i.ptr == nil {
		return
	}
	C.ImageBuf_clear(i.ptr)
	runtime.KeepAlive(i)
}

func (i *ImageBuf) Storage() IBStorage {
	ret := IBStorage(C.ImageBuf_storage(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Initialize this ImageBuf with the named image file, and read its header to
// fill out the spec correctly. Return true if this succeeded, false if the file
// could not be read. But don't allocate or read the pixels.
func (i *ImageBuf) InitSpec(filename string, subimage, miplevel int) error {
	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	ok := bool(C.ImageBuf_init_spec(i.ptr, c_str, C.int(subimage), C.int(miplevel)))
	if !ok {
		return i.LastError()
	}
	runtime.KeepAlive(i)
	return nil
}

// Reset forgets all previous info, and resets this ImageBuf to a new image
// that is uninitialized (no pixel values, no size or spec).
// 'cache' may be nil, or can be a specific cache to use instead of the global.
func (i *ImageBuf) Reset(filename string, cache *ImageCache) error {
	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	var cachePtr unsafe.Pointer = nil
	if cache != nil {
		cachePtr = cache.ptr
	}

	C.ImageBuf_reset_name_cache(i.ptr, c_str, cachePtr)
	err := i.LastError()
	runtime.KeepAlive(i)
	return err
}

// ResetSubImage forgets all previous info, and resets this ImageBuf to a new image
// that is uninitialized (no pixel values, no size or spec).
// If 'config' is not nil, it points to an ImageSpec giving requests
// or special instructions to be passed on to the eventual
// ImageInput.Open() call.
// 'cache' may be nil, or can be a specific cache to use instead of the global.
func (i *ImageBuf) ResetSubImage(filename string, subimage, miplevel int, cache *ImageCache, config *ImageSpec) error {
	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	var cachePtr unsafe.Pointer = nil
	if cache != nil {
		cachePtr = cache.ptr
	}

	var specPtr unsafe.Pointer = nil
	if config != nil {
		specPtr = config.ptr
	}

	C.ImageBuf_reset_subimage(i.ptr, c_str, C.int(subimage), C.int(miplevel), cachePtr, specPtr)
	err := i.LastError()
	runtime.KeepAlive(i)
	return err
}

// ResetSpec forgets all previous info, and reset this ImageBuf to a blank
// image of the given dimensions.
func (i *ImageBuf) ResetSpec(spec *ImageSpec) error {
	if spec == nil || spec.ptr == nil {
		return errors.New("nil ImageSpec")
	}

	var specPtr unsafe.Pointer = spec.ptr

	C.ImageBuf_reset_spec(i.ptr, specPtr)
	err := i.LastError()
	runtime.KeepAlive(i)
	return err
}

// ResetNameSpec forgets all previous info, and reset this ImageBuf to a blank
// image of the given name and dimensions.
func (i *ImageBuf) ResetNameSpec(filename string, spec *ImageSpec) error {
	if spec == nil || spec.ptr == nil {
		return errors.New("nil ImageSpec")
	}

	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	var specPtr unsafe.Pointer = spec.ptr

	C.ImageBuf_reset_name_spec(i.ptr, c_str, specPtr)
	err := i.LastError()
	runtime.KeepAlive(i)
	return err
}

// Read the file from disk. Generally will skip the read if we've already got a current
// version of the image in memory, unless force==true.
// This uses ImageInput underneath, so will read any file format for which an appropriate
// imageio plugin can be found.
func (i *ImageBuf) Read(force bool) error {
	ret := i.ReadFormatCallback(force, TypeUnknown, nil)
	runtime.KeepAlive(i)
	return ret
}

// Read the file from disk. Generally will skip the read if we've already got a current
// version of the image in memory, unless force==true.
// This uses ImageInput underneath, so will read any file format for which an appropriate
// imageio plugin can be found.
func (i *ImageBuf) ReadSubImage(subimage, miplevel int, force bool) error {
	ret := i.ReadSubImageFormatCallback(subimage, miplevel, force, TypeUnknown, nil)
	runtime.KeepAlive(i)
	return ret
}

// Read the file from disk. Generally will skip the read if we've already got a current
// version of the image in memory, unless force==true.
// This uses ImageInput underneath, so will read any file format for which an appropriate
// imageio plugin can be found.
//
// This call optionally supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
func (i *ImageBuf) ReadCallback(force bool, progress *ProgressCallback) error {
	ret := i.ReadFormatCallback(force, TypeUnknown, progress)
	runtime.KeepAlive(i)
	runtime.KeepAlive(progress)
	return ret
}

// Read the file from disk. Generally will skip the read if we've already got a current
// version of the image in memory, unless force==true.
// This uses ImageInput underneath, so will read any file format for which an appropriate
// imageio plugin can be found.
//
// This call optionally supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
func (i *ImageBuf) ReadSubImageCallback(subimage, miplevel int, force bool, progress *ProgressCallback) error {
	ret := i.ReadSubImageFormatCallback(subimage, miplevel, force, TypeUnknown, progress)
	runtime.KeepAlive(i)
	runtime.KeepAlive(progress)
	return ret
}

// Read the file from disk. Generally will skip the read if we've already got a current
// version of the image in memory, unless force==true.
// This uses ImageInput underneath, so will read any file format for which an appropriate
// imageio plugin can be found.
//
// Specify a specific conversion format or TypeUnknown for automatic handling.
//
// This call optionally supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
func (i *ImageBuf) ReadFormatCallback(force bool, convert TypeDesc, progress *ProgressCallback) error {
	var cbk unsafe.Pointer
	if progress != nil {
		cbk = unsafe.Pointer(progress)
	}

	ok := C.ImageBuf_read(i.ptr, 0, 0, C.bool(force), C.TypeDesc(convert), cbk)
	if !bool(ok) {
		return i.LastError()
	}
	runtime.KeepAlive(i)
	runtime.KeepAlive(progress)

	return nil
}

// Read the file from disk. Generally will skip the read if we've already got a current
// version of the image in memory, unless force==true.
// This uses ImageInput underneath, so will read any file format for which an appropriate
// imageio plugin can be found.
//
// Specify a specific conversion format or TypeUnknown for automatic handling.
//
// This call optionally supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
func (i *ImageBuf) ReadSubImageFormatCallback(subimage, miplevel int, force bool, convert TypeDesc, progress *ProgressCallback) error {
	var cbk unsafe.Pointer
	if progress != nil {
		cbk = unsafe.Pointer(progress)
	}

	ok := C.ImageBuf_read(i.ptr, C.int(subimage), C.int(miplevel), C.bool(force), C.TypeDesc(convert), cbk)
	if !bool(ok) {
		return i.LastError()
	}
	runtime.KeepAlive(i)
	runtime.KeepAlive(progress)

	return nil
}

// Write the image to the named file and file format
// (fileformat=="" means to infer the type from the filename extension).
func (i *ImageBuf) WriteFile(filepath, fileformat string) error {
	ret := i.WriteFileProgress(filepath, fileformat, nil)
	runtime.KeepAlive(i)
	return ret
}

// Write the image to the named file and file format
// (fileformat=="" means to infer the type from the filename extension).
//
// This call optionally supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
func (i *ImageBuf) WriteFileProgress(filepath, fileformat string, progress *ProgressCallback) error {
	var cbk unsafe.Pointer
	if progress != nil {
		cbk = unsafe.Pointer(progress)
	}

	c_path := C.CString(filepath)
	defer C.free(unsafe.Pointer(c_path))

	c_fmt := C.CString(fileformat)
	defer C.free(unsafe.Pointer(c_fmt))

	ok := C.ImageBuf_write_file(i.ptr, c_path, c_fmt, cbk)
	if !bool(ok) {
		return i.LastError()
	}
	runtime.KeepAlive(i)
	runtime.KeepAlive(progress)

	return nil
}

// Write the image to the open ImageOutput 'out'. Return true if all went ok, false if there were errors writing.
// It does NOT close the file when it's done (and so may be called in a loop to write a multi-image file).
func (i *ImageBuf) WriteImageOutput(output *ImageOutput) error {
	ret := i.WriteImageOutputProgress(output, nil)
	runtime.KeepAlive(i)
	runtime.KeepAlive(output)
	return ret
}

// Write the image to the open ImageOutput 'out'. Return true if all went ok, false if there were errors writing.
// It does NOT close the file when it's done (and so may be called in a loop to write a multi-image file).
//
// This call optionally supports passing a callback pointer to both track the progress,
// and to optionally abort the processing. The callback function will receive
// a float32 value indicating the percentage done of the processing, and should
// return true if the process should abort, and false if it should continue.
//
func (i *ImageBuf) WriteImageOutputProgress(output *ImageOutput, progress *ProgressCallback) error {
	var cbk unsafe.Pointer
	if progress != nil {
		cbk = unsafe.Pointer(progress)
	}

	ok := C.ImageBuf_write_output(i.ptr, output.ptr, cbk)
	if !bool(ok) {
		return i.LastError()
	}
	runtime.KeepAlive(i)
	runtime.KeepAlive(output)
	runtime.KeepAlive(progress)

	return nil
}

// Inform the ImageBuf what data format you'd like for any subsequent write().
func (i *ImageBuf) SetWriteFormat(format TypeDesc) {
	C.ImageBuf_set_write_format(i.ptr, C.TypeDesc(format))
	runtime.KeepAlive(i)
}

// Inform the ImageBuf what tile size (or no tiling, for 0) for any subsequent Write*()
func (i *ImageBuf) SetWriteTiles(width, height, depth int) {
	C.ImageBuf_set_write_tiles(i.ptr, C.int(width), C.int(height), C.int(depth))
	runtime.KeepAlive(i)
}

// Copy all the metadata from src to this (except for pixel data resolution,
// channel information, and data format).
func (i *ImageBuf) CopyMetadata(src *ImageBuf) error {
	C.ImageBuf_copy_metadata(i.ptr, src.ptr)
	runtime.KeepAlive(i)
	runtime.KeepAlive(src)
	return i.LastError()
}

// Copy the pixel data from src to this, automatically converting to the existing data
// format of this. It only copies pixels in the overlap regions (and channels) of the two
// images; pixel data in this that do exist in src will be set to 0, and pixel data in src
// that do not exist in this will not be copied.
func (i *ImageBuf) CopyPixels(src *ImageBuf) error {
	ok := bool(C.ImageBuf_copy_pixels(i.ptr, src.ptr))
	runtime.KeepAlive(i)
	runtime.KeepAlive(src)
	if !ok {
		return i.LastError()
	}
	return nil
}

// Try to copy the pixels and metadata from src to this, returning true upon success
// and false upon error/failure.
//
// If the previous state of this was uninitialized, owning its own local pixel memory,
// or referring to a read-only image backed by ImageCache, then local pixel memory will
// be allocated to hold the new pixels and the call always succeeds unless the memory cannot be allocated.
//
// If this previously referred to an app-owned memory buffer, the memory cannot be re-allocated,
// so the call will only succeed if the app-owned buffer is already the correct resolution and
// number of channels. The data type of the pixels will be converted automatically to the data
// type of the app buffer.
func (i *ImageBuf) Copy(src *ImageBuf) error {
	ok := bool(C.ImageBuf_copy(i.ptr, src.ptr))
	runtime.KeepAlive(i)
	runtime.KeepAlive(src)
	if !ok {
		return i.LastError()
	}
	return nil
}

// Swap with another ImageBuf.
func (i *ImageBuf) Swap(other *ImageBuf) error {
	C.ImageBuf_swap(i.ptr, other.ptr)
	runtime.KeepAlive(other)
	runtime.KeepAlive(i)
	return i.LastError()
}

// Return a reference to the image spec that describes the buffer.
func (i *ImageBuf) Spec() *ImageSpec {
	ret := &ImageSpec{C.ImageBuf_spec(i.ptr)}
	runtime.KeepAlive(i)
	return ret
}

// Return a reference to the "native" image spec (that describes the file, which may be slightly
// different than the spec of the ImageBuf, particularly if the IB is backed by an ImageCache
// that is imposing some particular data format or tile size).
func (i *ImageBuf) NativeSpec() *ImageSpec {
	ret := &ImageSpec{C.ImageBuf_nativespec(i.ptr)}
	runtime.KeepAlive(i)
	return ret
}

// Return a writable reference to the image spec that describes the buffer.
// Use with extreme caution! If you use this for anything other than adding
// attribute metadata, you are really taking your chances!
func (i *ImageBuf) SpecMod() *ImageSpec {
	ret := &ImageSpec{C.ImageBuf_specmod(i.ptr)}
	runtime.KeepAlive(i)
	return ret
}

// Return the name of this image.
func (i *ImageBuf) Name() string {
	ret := C.GoString(C.ImageBuf_name(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the name of the image file format of the disk file we read into this image.
// Returns an empty string if this image was not the result of a Read().
func (i *ImageBuf) FileFormatName() string {
	ret := C.GoString(C.ImageBuf_file_format_name(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the index of the subimage are we currently viewing
func (i *ImageBuf) SubImage() int {
	ret := int(C.ImageBuf_subimage(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the number of subimages in the file.
func (i *ImageBuf) NumSubImages() int {
	ret := int(C.ImageBuf_nsubimages(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the index of the miplevel are we currently viewing
func (i *ImageBuf) MipLevel() int {
	ret := int(C.ImageBuf_miplevel(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the number of miplevels of the current subimage.
func (i *ImageBuf) NumMipLevels() int {
	ret := int(C.ImageBuf_nmiplevels(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the number of color channels in the image.
func (i *ImageBuf) NumChannels() int {
	ret := int(C.ImageBuf_nchannels(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// GetFloatPixels retrieves the rectangle of pixels spanning the entire image,
// at the current subimage and MIP-map level, storing the float pixel values in
// a []float32
func (i *ImageBuf) GetFloatPixels() ([]float32, error) {
	spec := i.Spec()
	size := spec.Width() * spec.Height() * spec.Depth() * spec.NumChannels()
	pixels := make([]float32, size)
	ptr := unsafe.Pointer(&pixels[0])

	roi := i.ROI()

	ok := bool(C.ImageBuf_get_pixel_channels(
		i.ptr,
		C.int(roi.XBegin()), C.int(roi.XEnd()),
		C.int(roi.YBegin()), C.int(roi.YEnd()),
		C.int(roi.ZBegin()), C.int(roi.ZEnd()),
		C.int(roi.ChannelsBegin()), C.int(roi.ChannelsEnd()),
		(C.TypeDesc)(TypeFloat), ptr),
	)

	if !ok {
		return nil, i.LastError()
	}
	runtime.KeepAlive(i)

	return pixels, nil
}

// GetPixels retrieves the rectangle of pixels spanning the entire image,
// at the current subimage and MIP-map level, storing the pixel values in
// a slice that has been casted to an interface, and in the pixel format type
// described by 'format'.
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
//     val, err := buf.GetPixels(TypeFloat)
//     if err != nil {
//         panic(err.Error())
//     }
//     floatPixels := val.([]float32)
//
func (i *ImageBuf) GetPixels(format TypeDesc) (interface{}, error) {
	pixel_iface, ptr, err := allocatePixelBuffer(i.Spec(), format)
	if err != nil {
		return nil, err
	}

	roi := i.ROI()

	ok := bool(C.ImageBuf_get_pixel_channels(
		i.ptr,
		C.int(roi.XBegin()), C.int(roi.XEnd()),
		C.int(roi.YBegin()), C.int(roi.YEnd()),
		C.int(roi.ZBegin()), C.int(roi.ZEnd()),
		C.int(roi.ChannelsBegin()), C.int(roi.ChannelsEnd()),
		(C.TypeDesc)(format), ptr),
	)

	if !ok {
		return nil, i.LastError()
	}
	runtime.KeepAlive(i)

	return pixel_iface, nil
}

// GetPixelRegion retrieves the rectangle of pixels defined by an ROI,
// at the current subimage and MIP-map level, storing the pixel values in
// a slice that has been casted to an interface, and in the pixel format type
// described by 'format'.
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
//     val, err := buf.GetPixelRegion(roi, TypeFloat)
//     if err != nil {
//         panic(err.Error())
//     }
//     floatPixels := val.([]float32)
//
func (i *ImageBuf) GetPixelRegion(roi *ROI, format TypeDesc) (interface{}, error) {
	pixel_iface, ptr, err := allocatePixelBufferSize(roi.NumPixels()*roi.NumChannels(), format)
	if err != nil {
		return nil, err
	}

	ok := bool(C.ImageBuf_get_pixel_channels(
		i.ptr,
		C.int(roi.XBegin()), C.int(roi.XEnd()),
		C.int(roi.YBegin()), C.int(roi.YEnd()),
		C.int(roi.ZBegin()), C.int(roi.ZEnd()),
		C.int(roi.ChannelsBegin()), C.int(roi.ChannelsEnd()),
		(C.TypeDesc)(format), ptr),
	)
	runtime.KeepAlive(i)
	runtime.KeepAlive(roi)

	if !ok {
		return nil, i.LastError()
	}

	return pixel_iface, nil
}

// By default, image pixels are ordered from the top of the display to the bottom, and within
// each scanline, from left to right (i.e., the same ordering as English text and scan progression
// on a CRT). But the "Orientation" field can suggest that it should be displayed
// with a different orientation, according to the TIFF/EXIF conventions:
//   1 normal (top to bottom, left to right)
//   2 flipped horizontally (top to botom, right to left)
//   3 rotate 180 (bottom to top, right to left)
//   4 flipped vertically (bottom to top, left to right)
//   5 transposed (left to right, top to bottom)
//   6 rotated 90 clockwise (right to left, top to bottom)
//   7 transverse (right to left, bottom to top)
//   8 rotated 90 counter-clockwise (left to right, bottom to top)
func (i *ImageBuf) Orientation() int {
	ret := int(C.ImageBuf_orientation(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedWidth() int {
	ret := int(C.ImageBuf_oriented_width(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedHeight() int {
	ret := int(C.ImageBuf_oriented_height(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedX() int {
	ret := int(C.ImageBuf_oriented_x(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedY() int {
	ret := int(C.ImageBuf_oriented_y(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedFullWidth() int {
	ret := int(C.ImageBuf_oriented_full_width(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedFullHeight() int {
	ret := int(C.ImageBuf_oriented_full_height(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedFullX() int {
	ret := int(C.ImageBuf_oriented_full_x(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) OrientedFullY() int {
	ret := int(C.ImageBuf_oriented_full_y(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the beginning (minimum) x coordinate of the defined image.
func (i *ImageBuf) XBegin() int {
	ret := int(C.ImageBuf_xbegin(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the end (one past maximum) x coordinate of the defined image.
func (i *ImageBuf) XEnd() int {
	ret := int(C.ImageBuf_xend(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the beginning (minimum) y coordinate of the defined image
func (i *ImageBuf) YBegin() int {
	ret := int(C.ImageBuf_ybegin(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the end (one past maximum) y coordinate of the defined image.
func (i *ImageBuf) YEnd() int {
	ret := int(C.ImageBuf_yend(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the beginning (minimum) z coordinate of the defined image.
func (i *ImageBuf) ZBegin() int {
	ret := int(C.ImageBuf_zbegin(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the end (one past maximum) z coordinate of the defined image.
func (i *ImageBuf) ZEnd() int {
	ret := int(C.ImageBuf_zend(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the end (one past maximum) z coordinate of the defined image.
func (i *ImageBuf) XMin() int {
	ret := int(C.ImageBuf_xmin(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the maximum x coordinate of the defined image.
func (i *ImageBuf) XMax() int {
	ret := int(C.ImageBuf_xmax(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the minimum y coordinate of the defined image.
func (i *ImageBuf) YMin() int {
	ret := int(C.ImageBuf_ymin(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the maximum y coordinate of the defined image.
func (i *ImageBuf) YMax() int {
	ret := int(C.ImageBuf_ymax(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the minimum z coordinate of the defined image.
func (i *ImageBuf) ZMin() int {
	ret := int(C.ImageBuf_zmin(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the maximum z coordinate of the defined image.
func (i *ImageBuf) ZMax() int {
	ret := int(C.ImageBuf_zmax(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Set the "full" (a.k.a. display) window to [xbegin,xend) x [ybegin,yend) x [zbegin,zend).
func (i *ImageBuf) SetFull(xbegin, xend, ybegin, yend, zbegin, zend int) {
	C.ImageBuf_set_full(
		i.ptr,
		C.int(xbegin), C.int(xend),
		C.int(ybegin), C.int(yend),
		C.int(zbegin), C.int(zend))
	runtime.KeepAlive(i)
}

// Return pixel data window for this ImageBuf as a ROI.
func (i *ImageBuf) ROI() *ROI {
	ret := newROI(C.ImageBuf_roi(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return full/display window for this ImageBuf as a ROI.
func (i *ImageBuf) ROIFull() *ROI {
	ret := newROI(C.ImageBuf_roi_full(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Set full/display window for this ImageBuf to a ROI. Does NOT change the channels of the spec,
// regardless of newroi.
func (i *ImageBuf) SetROIFull(roi *ROI) error {
	C.ImageBuf_set_roi_full(i.ptr, roi.ptr)
	runtime.KeepAlive(i)
	runtime.KeepAlive(roi)
	return i.LastError()
}

func (i *ImageBuf) PixelsValid() bool {
	ret := bool(C.ImageBuf_pixels_valid(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) PixelType() TypeDesc {
	ret := TypeDesc(C.ImageBuf_pixeltype(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

func (i *ImageBuf) CachedPixels() bool {
	ret := bool(C.ImageBuf_cachedpixels(i.ptr))
	runtime.KeepAlive(i)
	return ret
}

// Return the ImageCache in which backs this ImageBuf
func (i *ImageBuf) ImageCache() *ImageCache {
	ret := &ImageCache{C.ImageBuf_imagecache(i.ptr)}
	runtime.KeepAlive(i)
	return ret
}

// Does this ImageBuf store deep data?
func (i *ImageBuf) Deep() bool {
	ret := bool(C.ImageBuf_deep(i.ptr))
	runtime.KeepAlive(i)
	return ret
}
