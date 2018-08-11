#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"

extern "C" {

void deleteImageOutput(ImageOutput *out) {
	delete static_cast<OIIO::ImageOutput*>(out);
}

ImageOutput* ImageOutput_Create(const char* filename, const char* plugin_searchpath) {
	std::string s_filename(filename);
	std::string s_path(plugin_searchpath);
	return (ImageOutput*) OIIO::ImageOutput::create(s_filename, s_path);
}

const char* ImageOutput_geterror(ImageOutput *out) {
	std::string sstring = static_cast<OIIO::ImageOutput*>(out)->geterror();
	if (sstring.empty()){
		return NULL;
	}
	return sstring.c_str();
}

const char* ImageOutput_format_name(ImageOutput *out) {
	return static_cast<OIIO::ImageOutput*>(out)->format_name();
}

const ImageSpec* ImageOutput_spec(ImageOutput *out) {
	const OIIO::ImageSpec *spec = &(static_cast<OIIO::ImageOutput*>(out)->spec());
	return (ImageSpec*) spec;
}

bool ImageOutput_supports(ImageOutput *out, const char* feature){
	std::string s_feature(feature);
	return static_cast<OIIO::ImageOutput*>(out)->supports(s_feature);
}


} // extern "C"
