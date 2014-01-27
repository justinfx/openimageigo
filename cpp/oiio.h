#ifndef _OPENIMAGEIGO_OIIO_H_
#define _OPENIMAGEIGO_OIIO_H_

#include <stdbool.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef ptrdiff_t stride_t;
typedef unsigned long long imagesize_t;

typedef void ImageSpec;
typedef void ImageInput;
typedef void ImageOutput;
typedef void DeepData;

typedef bool(* ProgressCallback)(void *opaque_data, float portion_done);

// TypeDesc
// 

typedef enum TypeDesc {
	TYPE_UNKNOWN 	= -1,
	TYPE_UINT8 		= 0,
	TYPE_INT8 		= 1,
	TYPE_UINT16 	= 2,
	TYPE_INT16 		= 3,
	TYPE_UINT 		= 4,
	TYPE_INT 		= 5,
	TYPE_UINT64 	= 6,
	TYPE_INT64 		= 7,
	TYPE_HALF 		= 8,
	TYPE_FLOAT 		= 9,
	TYPE_DOUBLE 	= 10
} TypeDesc;

// ImageInput
// 

void deleteImageInput(ImageInput *in);

ImageInput* ImageInput_Open(const char* filename, const ImageSpec *config);
ImageInput* ImageInput_Create(const char* filename, const char* plugin_searchpath);

const char* ImageInput_format_name(ImageInput *in);
bool ImageInput_valid_file(ImageInput *in, const char* filename);
bool ImageInput_open(ImageInput *in, const char* name, ImageSpec* newspec);
const ImageSpec* ImageInput_spec(ImageInput *in);
bool ImageInput_supports(ImageInput *in, const char* feature);
bool ImageInput_close(ImageInput *in);

// int ImageInput_current_subimage(ImageInput *in);
// int ImageInput_current_miplevel(ImageInput *in);
// bool ImageInput_seek_subimage(ImageInput *in, int subimage, ImageSpec* newspec);
// bool ImageInput_seek_subimage_mip(ImageInput *in, int subimage, int miplevel, ImageSpec* newspec);
bool ImageInput_read_scanline_floats(ImageInput *in, int y, int z, float* data);
// bool ImageInput_read_scanline_format(ImageInput *in, int y, int z, TypeDesc format, void* data, stride_t xstride);
// bool ImageInput_read_tile(ImageInput *in, int x, int y, int z, float* data);
// bool ImageInput_read_tile_format(ImageInput *in, int x, int y, int z, TypeDesc format, void* data, 
// 									stride_t xstride, stride_t ystride, stride_t zstride);
bool ImageInput_read_image_floats(ImageInput *in, float* data);

// // TODO: Progress Callback?
bool ImageInput_read_image_format(ImageInput *in, TypeDesc format, void* data, void* cbk_data);

// bool ImageInput_read_native_scanline(ImageInput *in, int y, int z, void *data);
// bool ImageInput_read_native_tile(ImageInput *in, int x, int y, int z, void *data);
// bool ImageInput_read_native_tiles(ImageInput *in, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend, void *data);
// bool ImageInput_read_native_deep_scanlines(ImageInput *in, int ybegin, int yend, int z, int chbegin, int chend, DeepData* deepdata);
// bool ImageInput_read_native_deep_tiles(ImageInput *in, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend, 
// 											int chbegin, int chend, DeepData &deepdata);
// bool ImageInput_read_native_deep_image(ImageInput *in, DeepData* deepdata);
// int ImageInput_send_to_input(ImageInput *in, const char *format,...);
// int ImageInput_send_to_client(ImageInput *in, const char *format,...);
 
const char* ImageInput_geterror(ImageInput *in);

 
// ImageSpec
// 

void deleteImageSpec(ImageSpec *spec);

ImageSpec* ImageSpec_New(TypeDesc fmt);
ImageSpec* ImageSpec_New_Size(int xres, int yres, int nchans, TypeDesc fmt);

void ImageSpec_set_format(ImageSpec *spec, TypeDesc fmt);
void ImageSpec_default_channel_names(ImageSpec *spec);
size_t ImageSpec_channel_bytes(ImageSpec *spec);
size_t ImageSpec_channel_bytes_chan(ImageSpec *spec, int chan, bool native);
size_t ImageSpec_pixel_bytes(ImageSpec *spec, bool native);
size_t ImageSpec_pixel_bytes_chans(ImageSpec *spec, int chbegin, int chend, bool native);
imagesize_t	ImageSpec_scanline_bytes(ImageSpec *spec, bool native);
imagesize_t	ImageSpec_tile_pixels(ImageSpec *spec);
imagesize_t	ImageSpec_tile_bytes(ImageSpec *spec, bool native);
imagesize_t	ImageSpec_image_pixels(ImageSpec *spec);
imagesize_t	ImageSpec_image_bytes(ImageSpec *spec, bool native);
bool ImageSpec_size_safe(ImageSpec *spec);

// void 	attribute (const std::string &name, TypeDesc type, const void *value)
// void 	attribute (const std::string &name, TypeDesc type, const std::string &value)
// void 	attribute (const std::string &name, unsigned int value)
// void 	attribute (const std::string &name, int value)
// void 	attribute (const std::string &name, float value)
// void 	attribute (const std::string &name, const char *value)
// void 	attribute (const std::string &name, const std::string &value)
// void 	erase_attribute (const std::string &name, TypeDesc searchtype=TypeDesc::UNKNOWN, bool casesensitive=false)
// ImageIOParameter * 	find_attribute (const std::string &name, TypeDesc searchtype=TypeDesc::UNKNOWN, bool casesensitive=false)
// const ImageIOParameter * 	find_attribute (const std::string &name, TypeDesc searchtype=TypeDesc::UNKNOWN, bool casesensitive=false) const
// int 	get_int_attribute (const std::string &name, int defaultval=0) const
// float 	get_float_attribute (const std::string &name, float defaultval=0) const
// std::string 	get_string_attribute (const std::string &name, const std::string &defaultval=std::string()) const
// std::string 	metadata_val (const ImageIOParameter &p, bool human=false) const
// std::string 	to_xml () const
// void 	from_xml (const char *xml)
// bool 	valid_tile_range (int xbegin, int xend, int ybegin, int yend, int zbegin, int zend)

TypeDesc ImageSpec_channelformat(ImageSpec *spec, int chan);
// void ImageSpec_get_channelformats(ImageSpec *spec, std::vector< TypeDesc > &formats);

// Properties
int ImageSpec_x(ImageSpec *spec);
int ImageSpec_y(ImageSpec *spec);
int ImageSpec_z(ImageSpec *spec);
int ImageSpec_width(ImageSpec *spec);
int ImageSpec_height(ImageSpec *spec);
int ImageSpec_depth(ImageSpec *spec);
int ImageSpec_full_x(ImageSpec *spec);
int ImageSpec_full_y(ImageSpec *spec);
int ImageSpec_full_z(ImageSpec *spec);
int ImageSpec_full_width(ImageSpec *spec);
int ImageSpec_full_height(ImageSpec *spec);
int ImageSpec_full_depth(ImageSpec *spec);
int ImageSpec_tile_width(ImageSpec *spec);
int ImageSpec_tile_height(ImageSpec *spec);
int ImageSpec_tile_depth(ImageSpec *spec);
int ImageSpec_nchannels(ImageSpec *spec);
TypeDesc ImageSpec_format(ImageSpec *spec);
void ImageSpec_channelformats(ImageSpec *spec, TypeDesc *out);
void ImageSpec_channelnames(ImageSpec *spec, char** out);
int ImageSpec_alpha_channel(ImageSpec *spec);
int ImageSpec_z_channel(ImageSpec *spec);
bool ImageSpec_deep(ImageSpec *spec);
int ImageSpec_quant_black(ImageSpec *spec);
int ImageSpec_quant_white(ImageSpec *spec);
int ImageSpec_quant_min(ImageSpec *spec);
int ImageSpec_quant_max(ImageSpec *spec);
// extra_attribs?

#ifdef __cplusplus
}
#endif
#endif