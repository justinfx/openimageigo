package oiio

/*
#include "stdlib.h"

#include "cpp/color.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// The ColorProcessor encapsulates a baked color transformation, suitable for
// application to raw pixels, or ImageBuf(s). These are generated using
// ColorConfig.CreateColorProcessor, and referenced in ImageBufAlgo
// (amongst other places)
type ColorProcessor struct {
	ptr unsafe.Pointer
}

func newColorProcessor(i unsafe.Pointer) *ColorProcessor {
	in := &ColorProcessor{i}
	runtime.SetFinalizer(in, deleteColorProcessor)
	return in
}

func deleteColorProcessor(i *ColorProcessor) {
	if i.ptr != nil {
		C.deleteColorProcessor(i.ptr)
		i.ptr = nil
	}
	runtime.KeepAlive(i)
}

// Represents the set of all color transformations that are allowed.
// If OpenColorIO is enabled at build time, this configuration is loaded
// at runtime, allowing the user to have complete control of all color
// transformation math. ($OCIO)  (See opencolorio.org for details).
// If OpenColorIO is not enabled at build time, a generic color configuration
// is provided for minimal color support.
//
// NOTE: ColorConfig(s) and ColorProcessor(s) are potentially heavy-weight.
// Their construction / destruction should be kept to a minimum.
type ColorConfig struct {
	ptr unsafe.Pointer
}

func newColorConfig(i unsafe.Pointer) *ColorConfig {
	in := &ColorConfig{i}
	runtime.SetFinalizer(in, deleteColorConfig)
	return in
}

func deleteColorConfig(i *ColorConfig) {
	if i.ptr != nil {
		C.free(i.ptr)
		i.ptr = nil
	}
	runtime.KeepAlive(i)
}

// Return if OpenImageIO was built with OCIO support
func SupportsOpenColorIO() bool {
	return bool(C.supportsOpenColorIO())
}

// If OpenColorIO is enabled at build time, initialize with the current
// color configuration. ($OCIO)
// If OpenColorIO is not enabled, this does nothing.
//
// Multiple calls to this are inexpensive.
func NewColorConfig() (*ColorConfig, error) {
	c := newColorConfig(C.New_ColorConfig())
	return c, c.error()
}

// If OpenColorIO is enabled at build time, initialize with the
// specified color configuration (.ocio) file
// If OpenColorIO is not enabled, this will result in an error.
//
// Multiple calls to this are potentially expensive.
func NewColorConfigPath(path string) (*ColorConfig, error) {
	c_str := C.CString(path)
	defer C.free(unsafe.Pointer(c_str))
	c := newColorConfig(C.New_ColorConfigPath(c_str))
	return c, c.error()
}

// Get the number of ColorSpace(s) defined in this configuration
func (c *ColorConfig) NumColorSpaces() int {
	ret := int(C.ColorConfig_getNumColorSpaces(c.ptr))
	runtime.KeepAlive(c)
	return ret
}

// Return the name of the colorspace at a given index
func (c *ColorConfig) ColorSpaceNameByIndex(index int) string {
	ret := C.GoString(C.ColorConfig_getColorSpaceNameByIndex(c.ptr, C.int(index)))
	runtime.KeepAlive(c)
	return ret
}

// Get the number of Looks defined in this configuration
func (c *ColorConfig) NumLooks() int {
	ret := int(C.ColorConfig_getNumLooks(c.ptr))
	runtime.KeepAlive(c)
	return ret
}

// Return the name of the look at a given index
func (c *ColorConfig) LookNameByIndex(index int) string {
	ret := C.GoString(C.ColorConfig_getLookNameByIndex(c.ptr, C.int(index)))
	runtime.KeepAlive(c)
	return ret
}

// Get the number of displays defined in this configuration
func (c *ColorConfig) NumDisplays() int {
	ret := int(C.ColorConfig_getNumDisplays(c.ptr))
	runtime.KeepAlive(c)
	return ret
}

// Return the name of the display at a given index
func (c *ColorConfig) DisplayNameByIndex(index int) string {
	ret := C.GoString(C.ColorConfig_getDisplayNameByIndex(c.ptr, C.int(index)))
	runtime.KeepAlive(c)
	return ret
}

// Get the number of displays defined in this configuration
func (c *ColorConfig) NumViews(displayName string) int {
	c_str := C.CString(displayName)
	defer C.free(unsafe.Pointer(c_str))
	ret := int(C.ColorConfig_getNumViews(c.ptr, c_str))
	runtime.KeepAlive(c)
	return ret
}

// Get the name of a view at a specific index of a display
func (c *ColorConfig) ViewNameByIndex(displayName string, index int) string {
	c_str := C.CString(displayName)
	defer C.free(unsafe.Pointer(c_str))
	ret := C.GoString(C.ColorConfig_getViewNameByIndex(c.ptr, c_str, C.int(index)))
	runtime.KeepAlive(c)
	return ret
}

// Get the name of the color space representing the named role,
// or empty string if none could be identified.
func (c *ColorConfig) ColorSpaceNameByRole(role string) string {
	c_str := C.CString(role)
	defer C.free(unsafe.Pointer(c_str))
	ret := C.GoString(C.ColorConfig_getColorSpaceNameByRole(c.ptr, c_str))
	runtime.KeepAlive(c)
	return ret
}

// Given the specified input and output ColorSpace, construct the
// processor.  It is possible that this will return nil and an error, if the
// inputColorSpace doesnt exist, the outputColorSpace doesn't
// exist, or if the specified transformation is illegal (for
// example, it may require the inversion of a 3D-LUT, etc).  When
// the user is finished with a ColorProcess, ColorProcess.Destroy()
// should be called.  ColorProcessor(s) remain valid even if the
// ColorConfig that created them no longer exists.
//
// Multiple calls to this are potentially expensive, so you should
// call once to create a ColorProcessor to use on an entire image
// (or multiple images), NOT for every scanline or pixel
// separately!
func (c *ColorConfig) CreateColorProcessor(inColorSpace, outColorSpace string) (*ColorProcessor, error) {
	c_in := C.CString(inColorSpace)
	defer C.free(unsafe.Pointer(c_in))

	c_out := C.CString(outColorSpace)
	defer C.free(unsafe.Pointer(c_out))

	ptr := C.ColorConfig_createColorProcessor(c.ptr, c_in, c_out)
	err := c.error()
	if err != nil {
		return nil, err
	}

	runtime.KeepAlive(c)

	return newColorProcessor(ptr), nil
}

// This routine will return the error string (and clear any error
// flags).  If no error has occurred since the last time GetError()
// was called, it will return an empty string.
func (c *ColorConfig) error() error {
	isError := C.ColorConfig_error(c.ptr)
	if C.bool(isError) {
		return errors.New(C.GoString(C.ColorConfig_geterror(c.ptr)))
	}
	runtime.KeepAlive(c)
	return nil
}
