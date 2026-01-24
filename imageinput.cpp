#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"
#include "imagespec.h"
#include "handles.h"

using oiio_go::UniqueHandle;


extern "C" {

void deleteImageInput(ImageInput *in) {
	delete static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
}

ImageInput* ImageInput_Open(const char* filename, const ImageSpec *config) {
	auto in = OIIO::ImageInput::open(filename, static_cast<const OIIO::ImageSpec*>(config));
	if (!in) return nullptr;
	return (ImageInput*) new UniqueHandle<OIIO::ImageInput>(std::move(in));
}

ImageInput* ImageInput_Create(const char* filename, const char* plugin_searchpath) {
	// OIIO 3.x signature: create(filename, do_open, config, ioproxy, plugin_searchpath)
	auto in = OIIO::ImageInput::create(filename, false, nullptr, nullptr, plugin_searchpath);
	if (!in) return nullptr;
	return (ImageInput*) new UniqueHandle<OIIO::ImageInput>(std::move(in));
}

char* ImageInput_geterror(ImageInput *in) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	std::string sstring = handle->get()->geterror();
	if (sstring.empty()) {
		return NULL;
	}
	return strdup(sstring.c_str());
}

const char* ImageInput_format_name(ImageInput *in) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->format_name();
}

bool ImageInput_valid_file(ImageInput *in, const char* filename) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->valid_file(filename);
}

bool ImageInput_open(ImageInput *in, const char* name, ImageSpec* newspec) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->open(name, *(static_cast<OIIO::ImageSpec*>(newspec)));
}

const ImageSpec* ImageInput_spec(ImageInput *in) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	const OIIO::ImageSpec *spec = &(handle->get()->spec());
	return (ImageSpec*) spec;
}

bool ImageInput_supports(ImageInput *in, const char* feature) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->supports(feature);
}

bool ImageInput_close(ImageInput *in) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->close();
}

int ImageInput_current_subimage(ImageInput *in) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->current_subimage();
}

bool ImageInput_seek_subimage(ImageInput *in, int subimage, ImageSpec* newspec) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	// In OIIO 3.x, seek_subimage(int, ImageSpec&) is deprecated
	// Use seek_subimage(int, int) instead and copy spec afterward
	bool ok = handle->get()->seek_subimage(subimage, 0);
	if (ok && newspec) {
		*(static_cast<OIIO::ImageSpec*>(newspec)) = handle->get()->spec();
	}
	return ok;
}

int ImageInput_current_miplevel(ImageInput *in) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->current_miplevel();
}

bool ImageInput_seek_subimage_miplevel(ImageInput *in, int subimage, int miplevel, ImageSpec* newspec) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	// In OIIO 3.x, seek_subimage(int, int, ImageSpec&) is deprecated
	// Use seek_subimage(int, int) instead and copy spec afterward
	bool ok = handle->get()->seek_subimage(subimage, miplevel);
	if (ok && newspec) {
		*(static_cast<OIIO::ImageSpec*>(newspec)) = handle->get()->spec();
	}
	return ok;
}

bool ImageInput_read_image_floats(ImageInput *in, float* data) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	// OIIO 3.x requires chbegin, chend parameters and stride parameters
	return handle->get()->read_image(0, 0, 0, -1, OIIO::TypeDesc::FLOAT, data,
									 OIIO::AutoStride, OIIO::AutoStride, OIIO::AutoStride);
}

bool ImageInput_read_image_format(ImageInput *in, TypeDesc format, void* data, void* cbk_data)
{
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	ProgressCallback cbk = NULL;
	if (cbk_data != NULL) {
		cbk = &image_progress_callback;
	}

	return handle->get()->read_image(
		0, 0,  // subimage, miplevel
		0, -1,  // chbegin, chend (all channels)
		fromTypeDesc(format),
		data,
		OIIO::AutoStride,
		OIIO::AutoStride,
		OIIO::AutoStride,
		cbk,
		cbk_data);
}

bool ImageInput_read_scanline_floats(ImageInput *in, int y, int z, float* data) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->read_scanline(y, z, OIIO::TypeDesc::FLOAT, data);
}

bool ImageInput_read_tile_floats(ImageInput *in, int x, int y, int z, float* data) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageInput>*>(in);
	return handle->get()->read_tile(x, y, z, OIIO::TypeDesc::FLOAT, data);
}


} // extern "C"
