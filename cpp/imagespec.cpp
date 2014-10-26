#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"


OIIO::TypeDesc fromTypeDesc(TypeDesc fmt) {
	switch (fmt) {
	case TYPE_UINT8: 	return OIIO::TypeDesc::UINT8;
	case TYPE_INT8: 	return OIIO::TypeDesc::INT8;
	case TYPE_UINT16: 	return OIIO::TypeDesc::UINT16;
	case TYPE_INT16: 	return OIIO::TypeDesc::INT16;	
	case TYPE_UINT: 	return OIIO::TypeDesc::UINT; 	
	case TYPE_INT: 		return OIIO::TypeDesc::INT; 	
	case TYPE_UINT64: 	return OIIO::TypeDesc::UINT64; 	
	case TYPE_INT64: 	return OIIO::TypeDesc::INT64; 	
	case TYPE_HALF: 	return OIIO::TypeDesc::HALF; 
	case TYPE_FLOAT: 	return OIIO::TypeDesc::FLOAT; 
	case TYPE_DOUBLE: 	return OIIO::TypeDesc::DOUBLE; 
	case TYPE_UNKNOWN: 	return OIIO::TypeDesc::UNKNOWN; 
	}
	return OIIO::TypeDesc::UNKNOWN;
}

TypeDesc toTypeDesc(OIIO::TypeDesc fmt) {
	if (fmt == OIIO::TypeDesc::UINT8) 	return TYPE_UINT8;
	if (fmt == OIIO::TypeDesc::INT8) 	return TYPE_INT8;
	if (fmt == OIIO::TypeDesc::UINT16) 	return TYPE_UINT16;
	if (fmt == OIIO::TypeDesc::INT16) 	return TYPE_INT16;
	if (fmt == OIIO::TypeDesc::UINT)	return TYPE_UINT;
	if (fmt == OIIO::TypeDesc::INT)		return TYPE_INT;
	if (fmt == OIIO::TypeDesc::UINT64) 	return TYPE_UINT64;
	if (fmt == OIIO::TypeDesc::INT64) 	return TYPE_INT64;
	if (fmt == OIIO::TypeDesc::HALF) 	return TYPE_HALF;
	if (fmt == OIIO::TypeDesc::FLOAT) 	return TYPE_FLOAT;
	if (fmt == OIIO::TypeDesc::DOUBLE) 	return TYPE_DOUBLE;
	return TYPE_UNKNOWN;
}

extern "C" {

void deleteImageSpec(ImageSpec *spec) {
	delete static_cast<OIIO::ImageSpec*>(spec);
}

ImageSpec* ImageSpec_New(TypeDesc fmt) {
	return (ImageSpec*) new OIIO::ImageSpec(fromTypeDesc(fmt));
}

ImageSpec* ImageSpec_New_Size(int xres, int yres, int nchans, TypeDesc fmt) {
	return (ImageSpec*) new OIIO::ImageSpec(xres, yres, nchans, fromTypeDesc(fmt));
}

void ImageSpec_default_channel_names(ImageSpec *spec) {
	static_cast<OIIO::ImageSpec*>(spec)->default_channel_names();
}

size_t ImageSpec_channel_bytes(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->channel_bytes();
}

size_t ImageSpec_channel_bytes_chan(ImageSpec *spec, int chan, bool native) {
	return static_cast<OIIO::ImageSpec*>(spec)->channel_bytes(chan, native);
}

size_t ImageSpec_pixel_bytes(ImageSpec *spec, bool native) {
	return static_cast<OIIO::ImageSpec*>(spec)->pixel_bytes(native);
}

size_t ImageSpec_pixel_bytes_chans(ImageSpec *spec, int chbegin, int chend, bool native) {
	return static_cast<OIIO::ImageSpec*>(spec)->pixel_bytes(chbegin, chend, native);
}

imagesize_t	ImageSpec_scanline_bytes(ImageSpec *spec, bool native) {
	return static_cast<OIIO::ImageSpec*>(spec)->scanline_bytes(native);
}

imagesize_t	ImageSpec_tile_pixels(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->tile_pixels();
}

imagesize_t	ImageSpec_tile_bytes(ImageSpec *spec, bool native) {
	return static_cast<OIIO::ImageSpec*>(spec)->tile_bytes(native);
}

imagesize_t	ImageSpec_image_pixels(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->image_pixels();
}

imagesize_t	ImageSpec_image_bytes(ImageSpec *spec, bool native) {
	return static_cast<OIIO::ImageSpec*>(spec)->image_bytes(native);
}

bool ImageSpec_size_safe (ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->size_t_safe();
}

TypeDesc ImageSpec_channelformat(ImageSpec *spec, int chan) {
	OIIO::TypeDesc c_spec = static_cast<OIIO::ImageSpec*>(spec)->channelformat(chan);
	return toTypeDesc(c_spec);
}

// Properties
int ImageSpec_x(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->x;
}

void ImageSpec_set_x(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->x = val;
}

int ImageSpec_y(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->y;
}

void ImageSpec_set_y(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->y = val;
}

int ImageSpec_z(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->z;
}

void ImageSpec_set_z(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->z = val;
}

int ImageSpec_width(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->width;
}

void ImageSpec_set_width(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->width = val;
}

int ImageSpec_height(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->height;
}

void ImageSpec_set_height(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->height = val;
}

int ImageSpec_depth(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->depth;
}

void ImageSpec_set_depth(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->depth = val;
}

int ImageSpec_full_x(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->full_x;
}

void ImageSpec_set_full_x(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->full_x = val;
}

int ImageSpec_full_y(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->full_y;
}

void ImageSpec_set_full_y(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->full_y = val;
}

int ImageSpec_full_z(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->full_z;
}

void ImageSpec_set_full_z(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->full_z = val;
}

int ImageSpec_full_width(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->full_width;
}

void ImageSpec_set_full_width(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->full_width = val;
}

int ImageSpec_full_height(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->full_height;
}

void ImageSpec_set_full_height(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->full_height = val;
}

int ImageSpec_full_depth(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->full_depth;
}

void ImageSpec_set_full_depth(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->full_depth = val;
}

int ImageSpec_tile_width(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->tile_width;
}

void ImageSpec_set_tile_width(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->tile_width = val;
}

int ImageSpec_tile_height(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->tile_height;
}

void ImageSpec_set_tile_height(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->tile_height = val;
}

int ImageSpec_tile_depth(ImageSpec *spec ){
	return static_cast<OIIO::ImageSpec*>(spec)->tile_depth;
}

void ImageSpec_set_tile_depth(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->tile_depth = val;
}

int ImageSpec_nchannels(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->nchannels;
}

void ImageSpec_set_nchannels(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->nchannels = val;
}

TypeDesc ImageSpec_format(ImageSpec *spec){
	OIIO::TypeDesc c_typ = static_cast<OIIO::ImageSpec*>(spec)->format;
	return toTypeDesc(c_typ);
}

void ImageSpec_set_format(ImageSpec *spec, TypeDesc fmt) {
	static_cast<OIIO::ImageSpec*>(spec)->set_format(fromTypeDesc(fmt));
}

void ImageSpec_channelformats(ImageSpec *spec, TypeDesc* out) {
	std::vector<OIIO::TypeDesc> vec = static_cast<OIIO::ImageSpec*>(spec)->channelformats;
	for (std::vector<OIIO::TypeDesc>::size_type i = 0; i != vec.size(); i++) {
		out[i] = toTypeDesc(vec[i]);
	}
}

void ImageSpec_set_channelformats(ImageSpec *spec, TypeDesc* formats){
	OIIO::ImageSpec *ptr = static_cast<OIIO::ImageSpec*>(spec);
	std::vector<OIIO::TypeDesc> vec = ptr->channelformats;
	for (std::vector<std::string>::size_type i = 0; i != vec.size(); i++) {
		vec[i] = fromTypeDesc(formats[i]);
	}
	ptr->channelformats = vec;
}

void ImageSpec_channelnames(ImageSpec *spec, char** out) {
	std::vector<std::string> vec = static_cast<OIIO::ImageSpec*>(spec)->channelnames;
	for (std::vector<std::string>::size_type i = 0; i != vec.size(); i++) {
		out[i] = (char*)vec[i].c_str();
	}
}

void ImageSpec_set_channelnames(ImageSpec *spec, char** names) {
	OIIO::ImageSpec *ptr = static_cast<OIIO::ImageSpec*>(spec);
	std::vector<std::string> vec = ptr->channelnames;
	for (std::vector<std::string>::size_type i = 0; i != vec.size(); i++) {
		vec[i] = std::string(names[i]);
	}
	ptr->channelnames = vec;
}

int ImageSpec_alpha_channel(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->alpha_channel;
}

void ImageSpec_set_alpha_channel(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->alpha_channel = val;
}

int ImageSpec_z_channel(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->z_channel;
}

void ImageSpec_set_z_channel(ImageSpec *spec, int val) {
	static_cast<OIIO::ImageSpec*>(spec)->z_channel = val;
}

bool ImageSpec_deep(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->deep;
}

void ImageSpec_set_deep(ImageSpec *spec, bool val) {
	static_cast<OIIO::ImageSpec*>(spec)->deep = val;
}

const char* ImageSpec_to_xml(ImageSpec *spec) {
	return static_cast<OIIO::ImageSpec*>(spec)->to_xml().c_str();
}

void ImageSpec_attribute_type_data(ImageSpec *spec, const char* name, TypeDesc type, const void *value) {
	static_cast<OIIO::ImageSpec*>(spec)->attribute(name, fromTypeDesc(type), value);
}

void ImageSpec_attribute_type_char(ImageSpec *spec, const char* name, TypeDesc type, const char* value) {
	static_cast<OIIO::ImageSpec*>(spec)->attribute(name, fromTypeDesc(type), value);
}

void ImageSpec_attribute_uint(ImageSpec *spec, const char* name, unsigned int value) {
	static_cast<OIIO::ImageSpec*>(spec)->attribute(name, value);
}

void ImageSpec_attribute_int(ImageSpec *spec, const char* name, int value) {
	static_cast<OIIO::ImageSpec*>(spec)->attribute(name, value);
}

void ImageSpec_attribute_float(ImageSpec *spec, const char* name, float value) {
	static_cast<OIIO::ImageSpec*>(spec)->attribute(name, value);
}
void ImageSpec_attribute_char(ImageSpec *spec, const char* name, const char* value) {
	static_cast<OIIO::ImageSpec*>(spec)->attribute(name, value);	
}

int ImageSpec_get_int_attribute(ImageSpec *spec, const char* name, int defaultval) {
	return static_cast<OIIO::ImageSpec*>(spec)->get_int_attribute(name, defaultval);	
}

float ImageSpec_get_float_attribute(ImageSpec *spec, const char* name, float defaultval) {
	return static_cast<OIIO::ImageSpec*>(spec)->get_float_attribute(name, defaultval);	
}

const char* ImageSpec_get_string_attribute(ImageSpec *spec, const char* name, const char* defaultval) {
	return static_cast<OIIO::ImageSpec*>(spec)->get_string_attribute(name, defaultval).c_str();	
}

} // extern "C"


