extern "C" {
	#include "_cgo_export.h"
}

#include <OpenImageIO/imageio.h>

#include <string>

extern "C" {

bool OIIO_attribute_int(const char* name, int val) {
	return OIIO::attribute(name, val);
}

bool OIIO_attribute_float(const char* name, float val) {
	return OIIO::attribute(name, val);
}

bool OIIO_attribute_string(const char* name, const char* val) {
	return OIIO::attribute(name, OIIO::string_view(val));
}

bool OIIO_getattribute_int(const char* name, int* val) {
	return OIIO::getattribute(name, *val);
}

bool OIIO_getattribute_float(const char* name, float* val) {
	return OIIO::getattribute(name, *val);
}

bool OIIO_getattribute_string(const char* name, char** val) {
	std::string s;
	bool ok = OIIO::getattribute(name, s);
	if (ok) {
		*val = strdup(s.c_str());
	}
	return ok;
}

} // extern "C"
