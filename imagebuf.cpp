#include <OpenImageIO/imagebuf.h>

#include "oiio.h"
#include "imagespec.h"
#include "handles.h"

using oiio_go::Handle;
using oiio_go::UniqueHandle;

OIIO::ImageBuf::IBStorage fromIBStorage(IBStorage s) {
	switch (s) {
	case IBSTORAGE_LOCALBUFFER: return OIIO::ImageBuf::LOCALBUFFER;
	case IBSTORAGE_APPBUFFER: 	return OIIO::ImageBuf::APPBUFFER;
	case IBSTORAGE_IMAGECACHE: 	return OIIO::ImageBuf::IMAGECACHE;
	case IBSTORAGE_UNINITIALIZED: 	return OIIO::ImageBuf::UNINITIALIZED;
	default: return OIIO::ImageBuf::UNINITIALIZED;
	}
}

IBStorage toIBStorage(OIIO::ImageBuf::IBStorage s) {
	if (s == OIIO::ImageBuf::LOCALBUFFER) return IBSTORAGE_LOCALBUFFER;
	if (s == OIIO::ImageBuf::APPBUFFER)   return IBSTORAGE_APPBUFFER;
	if (s == OIIO::ImageBuf::IMAGECACHE)  return IBSTORAGE_IMAGECACHE;
	return IBSTORAGE_UNINITIALIZED;
}


extern "C" {

char* ImageBuf_geterror(ImageBuf* buf) {
	if (!static_cast<OIIO::ImageBuf*>(buf)->has_error()) {
		return NULL;
	}
	std::string err = static_cast<OIIO::ImageBuf*>(buf)->geterror();
	return strdup(err.c_str());
}

void deleteImageBuf(ImageBuf *buf) {
	delete static_cast<OIIO::ImageBuf*>(buf);
}

ImageBuf* ImageBuf_New() {
	return (ImageBuf*) new OIIO::ImageBuf();
}

ImageBuf* ImageBuf_New_WithCache(const char* name, ImageCache *imagecache) {
	std::string s_name(name);
	if (imagecache == nullptr) {
		return (ImageBuf*) new OIIO::ImageBuf(s_name, 0, 0, nullptr);
	}
	// OIIO 2.x uses raw ImageCache* pointers, not shared_ptr
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(imagecache);
	return (ImageBuf*) new OIIO::ImageBuf(s_name, 0, 0, handle->ptr.get());
}

ImageBuf* ImageBuf_New_Spec(const ImageSpec* spec) {
	return (ImageBuf*) new OIIO::ImageBuf(*(static_cast<const OIIO::ImageSpec*>(spec)));
}

ImageBuf* ImageBuf_New_WithBuffer(const char* name, const ImageSpec* spec, void *buffer) {
	// OIIO 2.x constructor takes name parameter
	std::string s_name(name);
	return (ImageBuf*) new OIIO::ImageBuf(s_name, *(static_cast<const OIIO::ImageSpec*>(spec)), buffer);
}

ImageBuf* ImageBuf_New_SubImage(const char* name, int subimage, int miplevel, ImageCache* imagecache) {
	std::string s_name(name);
	if (imagecache == nullptr) {
		return (ImageBuf*) new OIIO::ImageBuf(s_name, subimage, miplevel, nullptr);
	}
	// OIIO 2.x uses raw ImageCache* pointers, not shared_ptr
	auto* handle = static_cast<Handle<OIIO::ImageCache>*>(imagecache);
	return (ImageBuf*) new OIIO::ImageBuf(s_name, subimage, miplevel, handle->ptr.get());
}


void ImageBuf_clear(ImageBuf* buf) {
	static_cast<OIIO::ImageBuf*>(buf)->clear();
}

void ImageBuf_reset_subimage(ImageBuf* buf, const char* name, int subimage, int miplevel,
							 ImageCache *imagecache, const ImageSpec* spec) {
	std::string s_name(name);
	// OIIO 2.x uses raw ImageCache* pointers, not shared_ptr
	OIIO::ImageCache* cache_ptr = nullptr;
	if (imagecache != nullptr) {
		auto* handle = static_cast<Handle<OIIO::ImageCache>*>(imagecache);
		cache_ptr = handle->ptr.get();
	}
	static_cast<OIIO::ImageBuf*>(buf)->reset(s_name, subimage, miplevel,
												cache_ptr,
												static_cast<const OIIO::ImageSpec*>(spec));
}

void ImageBuf_reset_name_cache(ImageBuf* buf, const char* name, ImageCache *imagecache) {
	std::string s_name(name);
	// OIIO 2.x uses raw ImageCache* pointers, not shared_ptr
	OIIO::ImageCache* cache_ptr = nullptr;
	if (imagecache != nullptr) {
		auto* handle = static_cast<Handle<OIIO::ImageCache>*>(imagecache);
		cache_ptr = handle->ptr.get();
	}
	static_cast<OIIO::ImageBuf*>(buf)->reset(s_name, 0, 0, cache_ptr);
}

void ImageBuf_reset_name_spec(ImageBuf* buf, const char* name, const ImageSpec* spec) {
	std::string s_name(name);
	static_cast<OIIO::ImageBuf*>(buf)->reset(s_name, *(static_cast<const OIIO::ImageSpec*>(spec)));
}

void ImageBuf_reset_spec(ImageBuf* buf, ImageSpec* spec) {
	static_cast<OIIO::ImageBuf*>(buf)->reset(*(static_cast<const OIIO::ImageSpec*>(spec)));
}


IBStorage ImageBuf_storage(ImageBuf* buf) {
	return toIBStorage(static_cast<OIIO::ImageBuf*>(buf)->storage());
}

bool ImageBuf_initialized(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->initialized();
}

bool ImageBuf_init_spec(ImageBuf* buf, const char* filename, int subimage, int miplevel) {
	return static_cast<OIIO::ImageBuf*>(buf)->init_spec(filename, subimage, miplevel);
}

bool ImageBuf_read(ImageBuf* buf, int subimage, int miplevel, bool force, TypeDesc convert, void *cbk_data) {
	ProgressCallback cbk = NULL;
	if (cbk_data != NULL) {
		cbk = &image_progress_callback;
	}
	return static_cast<OIIO::ImageBuf*>(buf)->read(subimage,
												   miplevel,
												   force,
												   fromTypeDesc(convert), 
												   cbk,
												   cbk_data);
}


bool ImageBuf_write_file(ImageBuf* buf, const char* filename, const char* fileformat, void *cbk_data) {
	ProgressCallback cbk = NULL;
	if (cbk_data != NULL) {
		cbk = &image_progress_callback;
	}
	return static_cast<OIIO::ImageBuf*>(buf)->write(filename, OIIO::TypeUnknown, fileformat, cbk, cbk_data);
}

bool ImageBuf_write_output(ImageBuf* buf, ImageOutput *out, void *cbk_data) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageOutput>*>(out);
	ProgressCallback cbk = NULL;
	if (cbk_data != NULL) {
		cbk = &image_progress_callback;
	}
	return static_cast<OIIO::ImageBuf*>(buf)->write(handle->get(), cbk, cbk_data);
}

void ImageBuf_set_write_format(ImageBuf* buf, TypeDesc format) {
	static_cast<OIIO::ImageBuf*>(buf)->set_write_format(fromTypeDesc(format));
}

void ImageBuf_set_write_tiles(ImageBuf* buf, int width, int height, int depth) {
	static_cast<OIIO::ImageBuf*>(buf)->set_write_tiles(width, height, depth);
}

void ImageBuf_copy_metadata(ImageBuf* dst, const ImageBuf* src) {
	const OIIO::ImageBuf *src_ptr = static_cast<const OIIO::ImageBuf*>(src);
	static_cast<OIIO::ImageBuf*>(dst)->copy_metadata(*src_ptr);
}

bool ImageBuf_copy_pixels(ImageBuf* dst, const ImageBuf* src) {
	const OIIO::ImageBuf *src_ptr = static_cast<const OIIO::ImageBuf*>(src);
	return static_cast<OIIO::ImageBuf*>(dst)->copy_pixels(*src_ptr);
}

bool ImageBuf_copy(ImageBuf* dst, const ImageBuf* src) {
	const OIIO::ImageBuf *src_ptr = static_cast<const OIIO::ImageBuf*>(src);
	return static_cast<OIIO::ImageBuf*>(dst)->copy(*src_ptr);
}

void ImageBuf_swap(ImageBuf* buf, ImageBuf* other) {
	OIIO::ImageBuf *other_ptr = static_cast<OIIO::ImageBuf*>(other);
	static_cast<OIIO::ImageBuf*>(buf)->swap(*other_ptr);
}

const ImageSpec* ImageBuf_spec(ImageBuf* buf) {
	OIIO::ImageSpec *spec = new OIIO::ImageSpec(static_cast<OIIO::ImageBuf*>(buf)->spec());
	return static_cast<const ImageSpec*>(spec);
}

ImageSpec* ImageBuf_specmod(ImageBuf* buf) {
	OIIO::ImageSpec *spec = &(static_cast<OIIO::ImageBuf*>(buf)->specmod());
	return static_cast<ImageSpec*>(spec);
}

const ImageSpec* ImageBuf_nativespec(ImageBuf* buf) {
	OIIO::ImageSpec *spec = new OIIO::ImageSpec(static_cast<OIIO::ImageBuf*>(buf)->nativespec());
	return static_cast<const ImageSpec*>(spec);
}

const char* ImageBuf_name(ImageBuf* buf) {
	// Return pointer to internal string data (valid until ImageBuf is modified)
	return static_cast<OIIO::ImageBuf*>(buf)->name().data();
}

const char* ImageBuf_file_format_name(ImageBuf* buf) {
	// Return pointer to internal string data (valid until ImageBuf is modified)
	return static_cast<OIIO::ImageBuf*>(buf)->file_format_name().data();
}

int ImageBuf_subimage(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->subimage();
}

int ImageBuf_nsubimages(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->nsubimages();
}

int ImageBuf_miplevel(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->miplevel();
}

int ImageBuf_nmiplevels(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->nmiplevels();
}

int ImageBuf_nchannels(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->nchannels();
}

bool ImageBuf_get_pixel_channels(ImageBuf* buf,
								 int xbegin, int xend,
								 int ybegin, int yend,
								 int zbegin, int zend,
								 int chbegin, int chend,
								 TypeDesc format,
								 void *result)
{
	// In OIIO 3.x, get_pixel_channels was removed, use get_pixels with ROI
	OIIO::ROI roi(xbegin, xend, ybegin, yend, zbegin, zend, chbegin, chend);
	return static_cast<OIIO::ImageBuf*>(buf)->get_pixels(roi, fromTypeDesc(format), result);
}

int ImageBuf_orientation(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->orientation();
}

int ImageBuf_oriented_width(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_width();
}

int ImageBuf_oriented_height(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_height();
}

int ImageBuf_oriented_x(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_x();
}

int ImageBuf_oriented_y(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_y();
}

int ImageBuf_oriented_full_width(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_full_width();
}

int ImageBuf_oriented_full_height(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_full_height();
}

int ImageBuf_oriented_full_x(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_full_x();
}

int ImageBuf_oriented_full_y(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->oriented_full_y();
}

int ImageBuf_xbegin(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->xbegin();
}

int ImageBuf_xend(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->xend();
}

int ImageBuf_ybegin(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->ybegin();
}

int ImageBuf_yend(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->yend();
}

int ImageBuf_zbegin(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->zbegin();
}

int ImageBuf_zend(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->zend();
}

int ImageBuf_xmin(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->xmin();
}

int ImageBuf_xmax(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->xmax();
}

int ImageBuf_ymin(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->ymin();
}

int ImageBuf_ymax(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->ymax();
}

int ImageBuf_zmin(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->zmin();
}

int ImageBuf_zmax(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->zmax();
}

void ImageBuf_set_full(ImageBuf* buf, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend) {
	static_cast<OIIO::ImageBuf*>(buf)->set_full(xbegin, xend, ybegin, yend, zbegin, zend);
}

ROI* ImageBuf_roi(ImageBuf* buf) {
	OIIO::ROI roi(static_cast<OIIO::ImageBuf*>(buf)->roi());
	OIIO::ROI *ptr = new OIIO::ROI();
	*ptr = roi;
	return static_cast<ROI*>(ptr);
}

ROI* ImageBuf_roi_full(ImageBuf* buf) {
	OIIO::ROI roi(static_cast<OIIO::ImageBuf*>(buf)->roi_full());
	OIIO::ROI *ptr = new OIIO::ROI();
	*ptr = roi;
	return static_cast<ROI*>(ptr);
}

void ImageBuf_set_roi_full(ImageBuf* buf, ROI* newroi) {
	static_cast<OIIO::ImageBuf*>(buf)->set_roi_full(*(static_cast<OIIO::ROI*>(newroi)));
}

bool ImageBuf_pixels_valid(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->pixels_valid();
}

TypeDesc ImageBuf_pixeltype(ImageBuf* buf) {
	return toTypeDesc(static_cast<OIIO::ImageBuf*>(buf)->pixeltype());
}

// void* ImageBuf_localpixels(ImageBuf* buf);

// const void* ImageBuf_localpixels(ImageBuf* buf);

bool ImageBuf_cachedpixels(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->cachedpixels();
}

ImageCache* ImageBuf_imagecache(ImageBuf* buf) {
	// OIIO 2.x returns raw ImageCache* pointer
	auto* cache = static_cast<OIIO::ImageBuf*>(buf)->imagecache();
	if (!cache) return nullptr;
	// Wrap raw pointer in a non-owning shared_ptr (doesn't delete on destruction)
	std::shared_ptr<OIIO::ImageCache> shared_cache(cache, [](OIIO::ImageCache*){});
	return (ImageCache*) new Handle<OIIO::ImageCache>(shared_cache);
}

// void* ImageBuf_pixeladdr(ImageBuf* buf, int x, int y);

// void* ImageBuf_pixeladdr_z(ImageBuf* buf, int x, int y, int z);

bool ImageBuf_deep(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->deep();
}

// int ImageBuf_deep_samples(ImageBuf* buf, int x, int y, int z);

// const void* ImageBuf_deep_pixel_ptr(ImageBuf* buf, int x, int y, int z, int c);

// float ImageBuf_deep_value(ImageBuf* buf, int x, int y, int z, int c, int s);

// DeepData* ImageBuf_deepdata(ImageBuf* buf);

} // extern "C"


