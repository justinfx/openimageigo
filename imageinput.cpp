#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"
#include "imagespec.h"


extern "C" {

void deleteImageInput(ImageInput *in) {
	delete static_cast<OIIO::ImageInput*>(in);
}

ImageInput* ImageInput_Open(const char* filename, const ImageSpec *config) {
	std::string s_filename(filename);
	return (ImageInput*) OIIO::ImageInput::open(s_filename, static_cast<const OIIO::ImageSpec*>(config));
}

ImageInput* ImageInput_Create(const char* filename, const char* plugin_searchpath) {
	std::string s_filename(filename);
	std::string s_path(plugin_searchpath);
	return (ImageInput*) OIIO::ImageInput::create(s_filename, s_path);

}

const char* ImageInput_geterror(ImageInput *in) {
	std::string sstring = static_cast<OIIO::ImageInput*>(in)->geterror();
	if (sstring.empty()) {
		return NULL;
	}
	return sstring.c_str();
}

const char* ImageInput_format_name(ImageInput *in) {
	return static_cast<OIIO::ImageInput*>(in)->format_name();
}

bool ImageInput_valid_file(ImageInput *in, const char* filename) {
	std::string s_name(filename);
	return static_cast<OIIO::ImageInput*>(in)->valid_file(s_name);	
}

bool ImageInput_open(ImageInput *in, const char* name, ImageSpec* newspec) {
	std::string s_name(name);
	return static_cast<OIIO::ImageInput*>(in)->open(s_name, *(static_cast<OIIO::ImageSpec*>(newspec)));	
}

const ImageSpec* ImageInput_spec(ImageInput *in) {
	const OIIO::ImageSpec *spec = &(static_cast<OIIO::ImageInput*>(in)->spec());	
	return (ImageSpec*) spec;
}

bool ImageInput_supports(ImageInput *in, const char* feature) {
	std::string s_feature(feature);
	return static_cast<OIIO::ImageInput*>(in)->supports(s_feature);	
}

bool ImageInput_close(ImageInput *in) {
	return static_cast<OIIO::ImageInput*>(in)->close();	
}

int ImageInput_current_subimage(ImageInput *in) {
	return static_cast<OIIO::ImageInput*>(in)->current_subimage();	
}

bool ImageInput_seek_subimage(ImageInput *in, int subimage, ImageSpec* newspec) {
	return static_cast<OIIO::ImageInput*>(in)->seek_subimage(
													subimage, 
													*(static_cast<OIIO::ImageSpec*>(newspec)));	
}

int ImageInput_current_miplevel(ImageInput *in) {
	return static_cast<OIIO::ImageInput*>(in)->current_miplevel();	
}

bool ImageInput_seek_subimage_miplevel(ImageInput *in, int subimage, int miplevel, ImageSpec* newspec) {
	return static_cast<OIIO::ImageInput*>(in)->seek_subimage(
													subimage, 
													miplevel,
													*(static_cast<OIIO::ImageSpec*>(newspec)));	
}

bool ImageInput_read_image_floats(ImageInput *in, float* data) {
	return static_cast<OIIO::ImageInput*>(in)->read_image(data);	
}

bool ImageInput_read_image_format(ImageInput *in, TypeDesc format, void* data, void* cbk_data)
{	
	ProgressCallback cbk = NULL;
	if (cbk_data != NULL) {
		cbk = &image_progress_callback;
	}

	return static_cast<OIIO::ImageInput*>(in)->read_image(
												fromTypeDesc(format), 
												data,
												OIIO::AutoStride,
												OIIO::AutoStride,
												OIIO::AutoStride,
												cbk,
												cbk_data);
}

bool ImageInput_read_scanline_floats(ImageInput *in, int y, int z, float* data) {
	return static_cast<OIIO::ImageInput*>(in)->read_scanline(y, z, data);	
}

bool ImageInput_read_tile_floats(ImageInput *in, int x, int y, int z, float* data) {
	return static_cast<OIIO::ImageInput*>(in)->read_tile(x, y, z, data);	
}


} // extern "C"


