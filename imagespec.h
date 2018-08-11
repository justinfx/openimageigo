#ifndef _OPENIMAGEIGO_IMAGESPEC_H_
#define _OPENIMAGEIGO_IMAGESPEC_H_

#include <OpenImageIO/imagebuf.h>

OIIO::TypeDesc fromTypeDesc(TypeDesc fmt);
TypeDesc toTypeDesc(OIIO::TypeDesc fmt);

#endif