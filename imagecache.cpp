#include <OpenImageIO/imagebuf.h>

#include "oiio.h"
#include <string>

extern "C" {

ImageCache* ImageCache_Create(bool shared) {
	return (ImageCache*) OIIO::ImageCache::create(shared);
}

void ImageCache_Destroy(ImageCache *x, bool teardown) {
	OIIO::ImageCache::destroy(static_cast<OIIO::ImageCache*>(x), teardown);
}

void ImageCache_close(ImageCache *x, const char *filename) {
	static_cast<OIIO::ImageCache*>(x)->close(OIIO::ustring(filename));
}

void ImageCache_close_all(ImageCache *x) {
	static_cast<OIIO::ImageCache*>(x)->close_all();
}

char* ImageCache_geterror(ImageCache* x) {
	std::string sstring = static_cast<OIIO::ImageCache*>(x)->geterror();
	if (sstring.empty()) {
		return NULL;
	}
	return strdup(sstring.c_str());
}

char* ImageCache_getstats(ImageCache *x, int level) {
	std::string str = static_cast<OIIO::ImageCache*>(x)->getstats(level);
	return strdup(str.c_str());
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


