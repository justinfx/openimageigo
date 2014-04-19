#include <OpenImageIO/imagebuf.h>

#include "oiio.h"

OIIO::ImageBuf::IBStorage fromIBStorage(IBStorage s) {
	switch (s) {
	case IBSTORAGE_LOCALBUFFER: return OIIO::ImageBuf::LOCALBUFFER;
	case IBSTORAGE_APPBUFFER: 	return OIIO::ImageBuf::APPBUFFER;
	case IBSTORAGE_IMAGECACHE: 	return OIIO::ImageBuf::IMAGECACHE;	
	}
	return OIIO::ImageBuf::UNINITIALIZED;
}

IBStorage toIBStorage(OIIO::ImageBuf::IBStorage s) {
	if (s == OIIO::ImageBuf::LOCALBUFFER) return IBSTORAGE_LOCALBUFFER;
	if (s == OIIO::ImageBuf::APPBUFFER)   return IBSTORAGE_APPBUFFER;
	if (s == OIIO::ImageBuf::IMAGECACHE)  return IBSTORAGE_IMAGECACHE;
	return IBSTORAGE_UNINITIALIZED;
}

extern "C" {

const char* ImageBuf_geterror(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->geterror().c_str();
}

void deleteImageBuf(ImageBuf *buf) {
	delete static_cast<OIIO::ImageBuf*>(buf);
}

ImageBuf* ImageBuf_New() {
	return (ImageBuf*) new OIIO::ImageBuf();
}

ImageBuf* ImageBuf_New_WithCache(const char* name, ImageCache *imagecache) {
	std::string s_name(name);
	return (ImageBuf*) new OIIO::ImageBuf(s_name, static_cast<OIIO::ImageCache*>(imagecache));
}

ImageBuf* ImageBuf_New_Spec(const ImageSpec* spec, void* buffer) {
	return (ImageBuf*) new OIIO::ImageBuf( *(static_cast<const OIIO::ImageSpec*>(spec)), buffer );
}

ImageBuf* ImageBuf_New_WithBuffer(const char* name, const ImageSpec* spec, void *buffer) {
	std::string s_name(name);
	return (ImageBuf*) new OIIO::ImageBuf(s_name, *(static_cast<const OIIO::ImageSpec*>(spec)), buffer);
}

ImageBuf* ImageBuf_New_SubImage(const char* name, int subimage, int miplevel, ImageCache* imagecache) {
	std::string s_name(name);
	return (ImageBuf*) new OIIO::ImageBuf(s_name, subimage, miplevel, 
											static_cast<OIIO::ImageCache*>(imagecache));
}


void ImageBuf_clear(ImageBuf* buf) {
	static_cast<OIIO::ImageBuf*>(buf)->clear();
}

void ImageBuf_reset_subimage(ImageBuf* buf, const char* name, int subimage, int miplevel, 
							 ImageCache *imagecache) {
	std::string s_name(name);
	static_cast<OIIO::ImageBuf*>(buf)->reset(s_name, subimage, miplevel, 
												static_cast<OIIO::ImageCache*>(imagecache));
}

void ImageBuf_reset_name_cache(ImageBuf* buf, const char* name, ImageCache *imagecache) {
	std::string s_name(name);
	static_cast<OIIO::ImageBuf*>(buf)->reset(s_name, static_cast<OIIO::ImageCache*>(imagecache));
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


bool ImageBuf_read(ImageBuf* buf, int subimage, int miplevel, bool force, TypeDesc convert, void *cbk_data) {
	ProgressCallback cbk = &read_image_format_callback;

	return static_cast<OIIO::ImageBuf*>(buf)->read(
												subimage,
												miplevel,
												force,
												fromTypeDesc(convert), 
												cbk,
												cbk_data);
}



const char* ImageBuf_name(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->name().c_str();
}

const char* ImageBuf_file_format_name(ImageBuf* buf) {
	return static_cast<OIIO::ImageBuf*>(buf)->file_format_name().c_str();
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



} // extern "C"


