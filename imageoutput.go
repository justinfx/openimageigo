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

// ImageOutput abstracts the writing of an image file in a file format-agnostic manner.
type ImageOutput struct {
	ptr unsafe.Pointer
}

func newImageOutput(i unsafe.Pointer) *ImageOutput {
	in := &ImageOutput{i}
	runtime.SetFinalizer(in, deleteImageOutput)
	return in
}

func deleteImageOutput(i *ImageOutput) {
	if i.ptr != nil {
		C.free(i.ptr)
		i.ptr = nil
	}
}

