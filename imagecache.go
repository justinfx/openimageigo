package oiio

/*
#include "stdlib.h"

#include "cpp/oiio.h"

*/
import "C"

import (
	"errors"
	"log"
	"runtime"
	"unsafe"
)

// Define an API to an abstract class that manages image files, caches of open
// file handles as well as tiles of pixels so that truly huge amounts of image
// data may be accessed by an application with low memory footprint.
type ImageCache struct {
	ptr unsafe.Pointer
}

func newImageCache(i unsafe.Pointer) *ImageCache {
	return &ImageCache{i}
}

// Create an ImageCache. *This should be freed by calling ImageCache.Destroy()*
//
// If shared==true, it's intended to be shared with other like-minded owners
// in the same process who also ask for a shared cache.
//
// If false, a private image cache will be created.
func CreateImageCache(shared bool) *ImageCache {
	ptr := C.ImageCache_Create(C.bool(false))
	return newImageCache(ptr)
}

// Destroy a ImageCache that was created using CreateImageCache().
// When 'teardown' parameter is set to true, it will fully destroy even a "shared" ImageCache.
func (i *ImageCache) Destroy(teardown bool) {
	if i.ptr != nil {
		C.ImageCache_Destroy(i.ptr, C.bool(teardown))
		i.ptr = nil
	}
	runtime.KeepAlive(i)
}

// Return the last error generated by API calls.
// An nil error will be returned if no error has occured.
func (i *ImageCache) LastError() error {
	c_str := C.ImageCache_geterror(i.ptr)
	runtime.KeepAlive(i)
	if c_str == nil {
		return nil
	}
	err := C.GoString(c_str)
	C.free(unsafe.Pointer(c_str))
	if err == "" {
		return nil
	}
	return errors.New(err)
}

// Close everything, free resources, start from scratch.
// Deprecated [oiio 1.7]
func (i *ImageCache) Clear() {
	log.Println("Deprecated ImageCache.Clear() [oiio 1.7]")
}

// Return the statistics output as a huge string.
// Suitable default for level == 1
func (i *ImageCache) GetStats(level int) string {
	c_stats := C.ImageCache_getstats(i.ptr, C.int(level))
	runtime.KeepAlive(i)
	stats := C.GoString(c_stats)
	C.free(unsafe.Pointer(c_stats))
	return stats
}

// Reset most statistics to be as they were with a fresh ImageCache.
// Caveat emptor: this does not flush the cache itelf, so the resulting
// statistics from the next set of texture requests will not match the number
// of tile reads, etc., that would have resulted from a new ImageCache.
func (i *ImageCache) ResetStats() {
	C.ImageCache_reset_stats(i.ptr)
	runtime.KeepAlive(i)
}

// Invalidate any loaded tiles or open file handles associated with the filename,
// so that any subsequent queries will be forced to re-open the file or re-load
// any tiles (even those that were previously loaded and would ordinarily be reused).
// A client might do this if, for example, they are aware that an image being held
// in the cache has been updated on disk. This is safe to do even if other procedures
// are currently holding reference-counted tile pointers from the named image,
// but those procedures will not get updated pixels until they release the tiles they
// are holding.
func (i *ImageCache) Invalidate(filename string) {
	c_str := C.CString(filename)
	defer C.free(unsafe.Pointer(c_str))

	C.ImageCache_invalidate(i.ptr, c_str)
	runtime.KeepAlive(i)
}

// Invalidate all loaded tiles and open file handles. This is safe to do even if other
// procedures are currently holding reference-counted tile pointers from the named
// image, but those procedures will not get updated pixels until they release the tiles
// they are holding.
// If force is true, everything will be invalidated, no matter how wasteful it is, but
// if force is false, in actuality files will only be invalidated if their modification
// times have been changed since they were first opened.
func (i *ImageCache) InvalidateAll(force bool) {
	C.ImageCache_invalidate_all(i.ptr, C.bool(force))
	runtime.KeepAlive(i)
}
