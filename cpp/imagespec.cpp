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
	return (ImageSpec*) new OIIO::ImageSpec(fromTypeDesc(fmt));;
}

ImageSpec* ImageSpec_New_Size(int xres, int yres, int nchans, TypeDesc fmt) {
	return (ImageSpec*) new OIIO::ImageSpec(xres, yres, nchans, fromTypeDesc(fmt));
}

void ImageSpec_set_format(ImageSpec *spec, TypeDesc fmt) {
	static_cast<OIIO::ImageSpec*>(spec)->set_format(fromTypeDesc(fmt));
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
int ImageSpec_x(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->x;
}

int ImageSpec_y(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->y;
}

int ImageSpec_z(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->z;
}

int ImageSpec_width(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->width;
}

int ImageSpec_height(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->height;
}

int ImageSpec_depth(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->depth;
}

int ImageSpec_full_x(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->full_x;
}

int ImageSpec_full_y(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->full_y;
}

int ImageSpec_full_z(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->full_z;
}

int ImageSpec_full_width(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->full_width;
}

int ImageSpec_full_height(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->full_height;
}

int ImageSpec_full_depth(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->full_depth;
}

int ImageSpec_tile_width(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->tile_width;
}

int ImageSpec_tile_height(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->tile_height;
}

int ImageSpec_tile_depth(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->tile_depth;
}

int ImageSpec_nchannels(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->nchannels;
}

TypeDesc ImageSpec_format(ImageSpec *spec){
	OIIO::TypeDesc c_typ = static_cast<OIIO::ImageSpec*>(spec)->format;
	return toTypeDesc(c_typ);
}

void ImageSpec_channelformats(ImageSpec *spec, TypeDesc* out) {
	std::vector<OIIO::TypeDesc> vec = static_cast<OIIO::ImageSpec*>(spec)->channelformats;
	for (std::vector<OIIO::TypeDesc>::size_type i = 0; i != vec.size(); i++) {
		out[i] = toTypeDesc(vec[i]);
	}
}

void ImageSpec_channelnames(ImageSpec *spec, char** out) {
	std::vector<std::string> vec = static_cast<OIIO::ImageSpec*>(spec)->channelnames;
	for (std::vector<std::string>::size_type i = 0; i != vec.size(); i++) {
		out[i] = (char*)vec[i].c_str();
	}
}

int ImageSpec_alpha_channel(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->alpha_channel;
}

int ImageSpec_z_channel(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->z_channel;
}

bool ImageSpec_deep(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->deep;
}

int ImageSpec_quant_black(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->quant_black;
}

int ImageSpec_quant_white(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->quant_white;
}

int ImageSpec_quant_min(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->quant_min;
}

int ImageSpec_quant_max(ImageSpec *spec){
	return static_cast<OIIO::ImageSpec*>(spec)->quant_max;
}

// extra_attribs?

} // extern "C"


