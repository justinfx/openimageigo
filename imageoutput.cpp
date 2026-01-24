#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"
#include "handles.h"

using oiio_go::UniqueHandle;

extern "C" {

void deleteImageOutput(ImageOutput *out) {
	delete static_cast<UniqueHandle<OIIO::ImageOutput>*>(out);
}

ImageOutput* ImageOutput_Create(const char* filename, const char* plugin_searchpath) {
	auto out = OIIO::ImageOutput::create(filename, plugin_searchpath);
	if (!out) return nullptr;
	return (ImageOutput*) new UniqueHandle<OIIO::ImageOutput>(std::move(out));
}

char* ImageOutput_geterror(ImageOutput *out) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageOutput>*>(out);
	std::string sstring = handle->get()->geterror();
	if (sstring.empty()){
		return NULL;
	}
	return strdup(sstring.c_str());
}

const char* ImageOutput_format_name(ImageOutput *out) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageOutput>*>(out);
	return handle->get()->format_name();
}

const ImageSpec* ImageOutput_spec(ImageOutput *out) {
	auto* handle = static_cast<UniqueHandle<OIIO::ImageOutput>*>(out);
	const OIIO::ImageSpec *spec = &(handle->get()->spec());
	return (ImageSpec*) spec;
}

bool ImageOutput_supports(ImageOutput *out, const char* feature){
	auto* handle = static_cast<UniqueHandle<OIIO::ImageOutput>*>(out);
	return handle->get()->supports(feature);
}


} // extern "C"
