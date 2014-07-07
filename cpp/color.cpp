#include <OpenImageIO/color.h>


extern "C" {

ColorConfig* New_ColorConfig() {
	return static_cast<ColorConfig*>(new OIIO::ColorConfig());
}

ColorConfig* New_ColorConfigPath(const char * filename) {
	return static_cast<ColorConfig*>(new OIIO::ColorConfig(filename));
}

bool supportsOpenColorIO() {
	return OIIO::ColorConfig::supportsOpenColorIO();
}

bool ColorConfig_error(ColorConfig* c) {
	return static_cast<OIIO::ColorConfig*>(c)->error();
}

const char* ColorConfig_geterror(ColorConfig* c) {
	if (!ColorConfig_error(c)) {
		return NULL;
	}
	return static_cast<OIIO::ImageOutput*>(c)->geterror().c_str();
}

int ColorConfig_getNumColorSpaces(ColorConfig* c) {
	return static_cast<OIIO::ColorConfig*>(c)->getNumColorSpaces();
}

const char * ColorConfig_getColorSpaceNameByIndex(ColorConfig* c, int index) {
	return static_cast<OIIO::ColorConfig*>(c)->getColorSpaceNameByIndex(index);
}

int ColorConfig_getNumLooks(ColorConfig* c) {
	return static_cast<OIIO::ColorConfig*>(c)->getNumLooks();
}

const char * ColorConfig_getLookNameByIndex(ColorConfig* c, int index) {
	return static_cast<OIIO::ColorConfig*>(c)->getLookNameByIndex(index);
}

int ColorConfig_getNumDisplays(ColorConfig* c) {
	return static_cast<OIIO::ColorConfig*>(c)->getNumDisplays();	
}

const char * ColorConfig_getDisplayNameByIndex(ColorConfig* c, int index) {
	return static_cast<OIIO::ColorConfig*>(c)->getDisplayNameByIndex(index);
}

int ColorConfig_getNumViews(ColorConfig* c, const char * display) {
	return static_cast<OIIO::ColorConfig*>(c)->getNumViews(display);	
}

const char * ColorConfig_getViewNameByIndex(ColorConfig* c, const char * display, int index) {
	return static_cast<OIIO::ColorConfig*>(c)->getViewNameByIndex(display, index);
}

const char * ColorConfig_getColorSpaceNameByRole(ColorConfig* c, const char *role) {
	return static_cast<OIIO::ColorConfig*>(c)->getColorSpaceNameByRole(role);	
}

void deleteColorProcessor(ColorProcessor* processor) {
	OIIO::ColorConfig::deleteColorProcessor(static_cast<OIIO::ColorProcessor*>(processor));
}

ColorProcessor* ColorConfig_createColorProcessor(ColorConfig* c, const char * inputColorSpace,
                                     				             const char * outputColorSpace) {
	OIIO::ColorProcessor *cp;
	cp = static_cast<OIIO::ColorConfig*>(c)->createColorProcessor(inputColorSpace, outputColorSpace);	
	return static_cast<ColorProcessor*>(cp);
}

} // extern "C"


