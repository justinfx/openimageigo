#include <OpenImageIO/imagecache.h>

#include "oiio.h"
#include "handles.h"
#include <string>

using oiio_go::Handle;

extern "C" {

ImageCache* ImageCache_Create(bool shared) {
	auto cache = OIIO::ImageCache::create(shared);
	return (ImageCache*) new Handle<OIIO::ImageCache>(cache);
}

void ImageCache_Destroy(ImageCache *x, bool teardown) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	OIIO::ImageCache::destroy(handle->ptr, teardown);
}

void deleteImageCache(ImageCache *x) {
	delete static_cast<Handle<OIIO::ImageCache>*>(x);
}

void ImageCache_clear(ImageCache *x) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	// clear() removed in OIIO 3.x, use invalidate_all(true) instead
	handle->get()->invalidate_all(true);
}

char* ImageCache_geterror(ImageCache* x) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	std::string sstring = handle->get()->geterror();
	if (sstring.empty()) {
		return NULL;
	}
	return strdup(sstring.c_str());
}

char* ImageCache_getstats(ImageCache *x, int level) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	std::string str = handle->get()->getstats(level);
	return strdup(str.c_str());
}

void ImageCache_reset_stats(ImageCache *x) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	handle->get()->reset_stats();
}

void ImageCache_invalidate(ImageCache *x, const char *filename) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	OIIO::ustring s(filename);
	handle->get()->invalidate(s);
}

void ImageCache_invalidate_all(ImageCache *x, bool force) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	handle->get()->invalidate_all(force);
}

// Attribute setters
bool ImageCache_attribute_int(ImageCache *x, const char *name, int val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->attribute(name, val);
}

bool ImageCache_attribute_float(ImageCache *x, const char *name, float val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->attribute(name, val);
}

bool ImageCache_attribute_double(ImageCache *x, const char *name, double val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->attribute(name, val);
}

bool ImageCache_attribute_string(ImageCache *x, const char *name, const char *val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->attribute(name, val);
}

// Attribute getters
bool ImageCache_getattribute_int(ImageCache *x, const char *name, int *val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->getattribute(name, *val);
}

bool ImageCache_getattribute_float(ImageCache *x, const char *name, float *val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->getattribute(name, *val);
}

bool ImageCache_getattribute_double(ImageCache *x, const char *name, double *val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	return handle->get()->getattribute(name, *val);
}

bool ImageCache_getattribute_string(ImageCache *x, const char *name, char **val) {
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(x);
	std::string str;
	bool ok = handle->get()->getattribute(name, str);
	if (ok && val) {
		*val = strdup(str.c_str());
	}
	return ok;
}

} // extern "C"
