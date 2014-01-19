#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"


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
	return static_cast<OIIO::ImageInput*>(in)->geterror().c_str();
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

bool ImageInput_read_image(ImageInput *in, float* data) {
	return static_cast<OIIO::ImageInput*>(in)->read_image(data);	
}

} // extern "C"

