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

// ImageSpec describes the data format of an image – dimensions, layout, 
// number and meanings of image channels.
type ImageSpec struct {
	ptr unsafe.Pointer
}

func newImageSpec(i unsafe.Pointer) *ImageSpec {
	spec := new(ImageSpec)
	spec.ptr = i
	runtime.SetFinalizer(spec, deleteImageSpec)
	return spec
}

func deleteImageSpec(i *ImageSpec) {
	if i.ptr != nil {
		C.free(i.ptr)
		i.ptr = nil
	}
}

// given just the data format, set the default quantize and set all other channels to something reasonable.
func NewImageSpec(format TypeDesc) *ImageSpec {
	spec := C.ImageSpec_New((C.TypeDesc)(format))
	return newImageSpec(spec)
}

// for simple 2D scanline image with nothing special. If fmt is not supplied, default to unsigned 8-bit data.
func NewImageSpecSize(x, y, chans int, format TypeDesc) *ImageSpec {
	spec := C.ImageSpec_New_Size(C.int(x), C.int(y), C.int(chans), (C.TypeDesc)(format))
	return newImageSpec(spec)
}

// Set the data format, and as a side effect set quantize to good defaults for that format
func (s *ImageSpec) SetFormat(format TypeDesc) {
	C.ImageSpec_set_format(s.ptr, (C.TypeDesc)(format))
}

// Set the channelnames to reasonable defaults ("R", "G", "B", "A"),
// and alpha_channel, based on the number of channels.
func (s *ImageSpec) DefaultChannelNames() {
	C.ImageSpec_default_channel_names(s.ptr)
}

// Return the number of bytes for each channel datum, assuming they are
// all stored using the data format given by format.
func (s *ImageSpec) ChannelBytes() int {
	return int(C.ImageSpec_channel_bytes(s.ptr))
}

// Return the number of bytes needed for the single specified channel.
// If native is false (default), compute the size of one channel of format,
// but if native is true, compute the size of the channel in terms of the "native"
// data format of that channel as stored in the file.
func (s *ImageSpec) ChannelBytesChan(chanNum int, native bool) int {
	return int(C.ImageSpec_channel_bytes_chan(s.ptr, C.int(chanNum), C.bool(native)))
}

// Return the number of bytes for each pixel (counting all channels). If native is
// false, assume all channels are in format, but if native is true,
// compute the size of a pixel in the "native" data format of the file (these may
// differ in the case of per-channel formats). This will return a max value
// in the event of an overflow where it's not representable in a int.
func (s *ImageSpec) PixelBytes(native bool) int {
	return int(C.ImageSpec_pixel_bytes(s.ptr, C.bool(native)))
}

// Return the number of bytes for just the subset of channels in each pixel described
// by [chanBegin, chanEnd). If native is false, assume all channels are in format,
// but if native is true, compute the size of a pixel in the "native" data format of
// the file (these may differ in the case of per-channel formats). This will return
// a max value in the event of an overflow where it's not representable in a int.
func (s *ImageSpec) PixelBytesChans(chanBegin, chanEnd int, native bool) int {
	return int(C.ImageSpec_pixel_bytes_chans(s.ptr, C.int(chanBegin), C.int(chanEnd), C.bool(native)))
}

// Return the number of bytes for each scanline.
// This will return a max value in the event of an overflow where it's not
// representable in an int. If native is false, assume all channels are in
// format, but if native is true, compute the size of a pixel in the "native"
// data format of the file (these may differ in the case of per-channel formats).
func (s *ImageSpec) ScanlineBytes(native bool) int {
	return int(C.ImageSpec_scanline_bytes(s.ptr, C.bool(native)))
}

// Return the number of pixels for a tile. This will return a max value in the event
// of an overflow where it's not representable in an int.
func (s *ImageSpec) TilePixels() int {
	return int(C.ImageSpec_tile_pixels(s.ptr))
}

// Return the number of bytes for each a tile of the image.
// This will return a max value in the event of an overflow where it's not
// representable in an imagesize_t. If native is false, assume all channels are
// in format, but if native is true, compute the size of a pixel in the "native"
// data format of the file (these may differ in the case of per-channel formats).
func (s *ImageSpec) TileBytes(native bool) int {
	return int(C.ImageSpec_tile_bytes(s.ptr, C.bool(native)))
}

// Return the number of pixels for an entire image.
// This will return a max value in the event of an overflow where
// it's not representable in an int.
func (s *ImageSpec) ImagePixels() int {
	return int(C.ImageSpec_image_pixels(s.ptr))
}

// Return the number of bytes for an entire image.
// This will return a max value in the event of an overflow where it's not
// representable in an int. If native is false, assume all channels are in
// format, but if native is true, compute the size of a pixel in the "native"
// data format of the file (these may differ in the case of per-channel formats).
func (s *ImageSpec) ImageBytes(native bool) int {
	return int(C.ImageSpec_image_bytes(s.ptr, C.bool(native)))
}

// Verify that on this platform, a size_t is big enough to hold the number of bytes
// (and pixels) in a scanline, a tile, and the whole image. If this returns false,
// the image is much too big to allocate and read all at once, so client apps beware
// and check these routines for overflows!
func (s *ImageSpec) SizeSafe() bool {
	return bool(C.ImageSpec_size_safe(s.ptr))
}

func (s *ImageSpec) ChannelFormat(chanNum int) TypeDesc {
	return (TypeDesc)(C.ImageSpec_channelformat(s.ptr, C.int(chanNum)))
}

// Properties
func (s *ImageSpec) X() int {
	return int(C.ImageSpec_x(s.ptr))
}

func (s *ImageSpec) Y() int {
	return int(C.ImageSpec_y(s.ptr))
}

// origin (upper left corner) of pixel data
func (s *ImageSpec) Z() int {
	return int(C.ImageSpec_z(s.ptr))
}

// width of the pixel data window
func (s *ImageSpec) Width() int {
	return int(C.ImageSpec_width(s.ptr))
}

// height of the pixel data window
func (s *ImageSpec) Height() int {
	return int(C.ImageSpec_height(s.ptr))
}

// depth of pixel data, >1 indicates a "volume"
func (s *ImageSpec) Depth() int {
	return int(C.ImageSpec_depth(s.ptr))
}

// origin of the full (display) window
func (s *ImageSpec) FullX() int {
	return int(C.ImageSpec_full_x(s.ptr))
}

// origin of the full (display) window
func (s *ImageSpec) FullY() int {
	return int(C.ImageSpec_full_y(s.ptr))
}

// origin of the full (display) window
func (s *ImageSpec) FullZ() int {
	return int(C.ImageSpec_full_z(s.ptr))
}

// width of the full (display) window
func (s *ImageSpec) FullWidth() int {
	return int(C.ImageSpec_full_width(s.ptr))
}

// height of the full (display) window
func (s *ImageSpec) FullHeight() int {
	return int(C.ImageSpec_full_height(s.ptr))
}

// depth of the full (display) window
func (s *ImageSpec) FullDepth() int {
	return int(C.ImageSpec_full_depth(s.ptr))
}

// tile width (0 for a non-tiled image)
func (s *ImageSpec) TileWidth() int {
	return int(C.ImageSpec_tile_width(s.ptr))
}

// tile height (0 for a non-tiled image)
func (s *ImageSpec) TileHeight() int {
	return int(C.ImageSpec_tile_height(s.ptr))
}

func (s *ImageSpec) TileDepth() int {
	return int(C.ImageSpec_tile_depth(s.ptr))
}

// number of image channels, e.g., 4 for RGBA
func (s *ImageSpec) NumChannels() int {
	return int(C.ImageSpec_nchannels(s.ptr))
}

// data format of the channels
func (s *ImageSpec) Format() TypeDesc {
	return (TypeDesc)(C.ImageSpec_format(s.ptr))
}

// Optional per-channel formats.
func (s *ImageSpec) ChannelFormats() []TypeDesc {
	formats := make([]TypeDesc, s.NumChannels())
	formats_ptr := (*C.TypeDesc)(unsafe.Pointer(&formats[0]))
	C.ImageSpec_channelformats(s.ptr, formats_ptr)
	return formats
}

func (s *ImageSpec) ChannelNames() []string {
	names := make([]string, s.NumChannels())
	c_names := make([]*C.char, s.NumChannels())
	c_names_ptr := (**C.char)(unsafe.Pointer(&c_names[0]))
	C.ImageSpec_channelnames(s.ptr, c_names_ptr)
	for i, c := range c_names {
		names[i] = C.GoString(c)
	}
	return names
}

// Index of alpha channel, or -1 if not known.
func (s *ImageSpec) AlphaChannel() int {
	return int(C.ImageSpec_alpha_channel(s.ptr))
}

// Index of depth channel, or -1 if not known.
func (s *ImageSpec) ZChannel() int {
	return int(C.ImageSpec_z_channel(s.ptr))
}

// Contains deep data.
func (s *ImageSpec) Deep() bool {
	return bool(C.ImageSpec_deep(s.ptr))
}

// quantization of black (0.0) level
func (s *ImageSpec) QuantBlack() int {
	return int(C.ImageSpec_quant_black(s.ptr))
}

// quantization of white (1.0) level
func (s *ImageSpec) QuantWhite() int {
	return int(C.ImageSpec_quant_white(s.ptr))
}

// quantization minimum clamp value
func (s *ImageSpec) QuantMin() int {
	return int(C.ImageSpec_quant_min(s.ptr))
}

// quantization maximum clamp value
func (s *ImageSpec) QuantMax() int {
	return int(C.ImageSpec_quant_max(s.ptr))
}
