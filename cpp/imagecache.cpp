#include <OpenImageIO/imagebuf.h>

#include "oiio.h"

extern "C" {

ImageCache* ImageCache_Create(bool shared) {
	return (ImageCache*) OIIO::ImageCache::create(shared);
}

void ImageCache_Destroy(ImageCache *x, bool teardown) {
	OIIO::ImageCache::destroy(static_cast<OIIO::ImageCache*>(x), teardown);
}

const char* ImageCache_geterror(ImageCache* x) {
	std::string sstring = static_cast<OIIO::ImageCache*>(x)->geterror();
	if (sstring.empty()) {
		return NULL;
	}
	return sstring.c_str();
}

const char* ImageCache_getstats(ImageCache *x, int level) {
	return static_cast<OIIO::ImageCache*>(x)->getstats(level).c_str();
}

void ImageCache_reset_stats(ImageCache *x) {
	static_cast<OIIO::ImageCache*>(x)->reset_stats();
}

void ImageCache_invalidate(ImageCache *x, const char *filename) {
	OIIO::ustring s(filename);
	static_cast<OIIO::ImageCache*>(x)->invalidate(s);
}

void ImageCache_invalidate_all(ImageCache *x, bool force) {
	static_cast<OIIO::ImageCache*>(x)->invalidate_all(force);
}


} // extern "C"


