#ifndef _OPENIMAGEIGO_COLOR_H_
#define _OPENIMAGEIGO_COLOR_H_

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void ColorProcessor;
typedef void ColorConfig;

ColorConfig* New_ColorConfig();
ColorConfig* New_ColorConfigPath(const char * filename);

void deleteColorConfig(ColorConfig* c);

bool ColorConfig_error(ColorConfig* c);

const char* ColorConfig_geterror(ColorConfig* c);

int ColorConfig_getNumColorSpaces(ColorConfig* c);

const char * ColorConfig_getColorSpaceNameByIndex(ColorConfig* c, int index);

const char * ColorConfig_getColorSpaceNameByRole(ColorConfig* c, const char *role);

int ColorConfig_getNumLooks(ColorConfig* c);

const char * ColorConfig_getLookNameByIndex(ColorConfig* c, int index);

ColorProcessor* ColorConfig_createColorProcessor(ColorConfig* c, const char * inputColorSpace,
                                     				const char * outputColorSpace);

// ColorProcessor* createLookTransform (ColorConfig* c, const char * looks,
//                                      const char * inputColorSpace,
//                                      const char * outputColorSpace,
//                                      bool inverse,
//                                      const char *context_key,
//                                      const char *context_value);

int ColorConfig_getNumDisplays(ColorConfig* c);

const char * ColorConfig_getDisplayNameByIndex(ColorConfig* c, int index);

int ColorConfig_getNumViews(ColorConfig* c, const char * display);

const char * ColorConfig_getViewNameByIndex(ColorConfig* c, const char * display, int index);

// const char * getDefaultDisplayName(ColorConfig* c);

// const char * getDefaultViewName(ColorConfig* c, const char * display);

// ColorProcessor* createDisplayTransform (ColorConfig* c, const char * display,
//                                         const char * view,
//                                         const char * inputColorSpace,
//                                         const char * looks,
//                                         const char * context_key,
//                                         const char * context_value);

void deleteColorProcessor(ColorProcessor* processor);

bool supportsOpenColorIO();

#ifdef __cplusplus
}
#endif
#endif