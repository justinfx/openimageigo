#ifndef _OPENIMAGEIGO_OIIO_H_
#define _OPENIMAGEIGO_OIIO_H_

#include <stdbool.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif


static char** makeCharArray(int size) {
        return (char**)calloc(sizeof(char*), size);
}

static void setArrayString(char **a, char *s, int n) {
        a[n] = s;
}

static void freeCharArray(char **a, int size) {
        int i;
        for (i = 0; i < size; i++)
                free(a[i]);
        free(a);
}


typedef ptrdiff_t stride_t;
typedef unsigned long long imagesize_t;

typedef void ImageSpec;
typedef void ImageInput;
typedef void ImageOutput;
typedef void ImageCache;
typedef void DeepData;
typedef void ImageBuf;
typedef void ROI;
typedef void Tile;

typedef bool(* ProgressCallback)(void *opaque_data, float portion_done);



// Enums
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


typedef enum IBStorage {
	IBSTORAGE_UNINITIALIZED,
	IBSTORAGE_LOCALBUFFER,
	IBSTORAGE_APPBUFFER,
	IBSTORAGE_IMAGECACHE,
} IBStorage;


typedef enum WrapMode {
	WrapDefault,
	WrapBlack,
	WrapClamp,
	WrapPeriodic,
	WrapMirror,
	_WrapLast
} WrapMode;


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

int ImageInput_current_subimage(ImageInput *in);
int ImageInput_current_miplevel(ImageInput *in);
bool ImageInput_seek_subimage(ImageInput *in, int subimage, ImageSpec* newspec);
bool ImageInput_seek_subimage_miplevel(ImageInput *in, int subimage, int miplevel, ImageSpec* newspec);
bool ImageInput_read_scanline_floats(ImageInput *in, int y, int z, float* data);
// bool ImageInput_read_scanline_format(ImageInput *in, int y, int z, TypeDesc format, void* data, stride_t xstride);
bool ImageInput_read_tile_floats(ImageInput *in, int x, int y, int z, float* data);
// bool ImageInput_read_tile_format(ImageInput *in, int x, int y, int z, TypeDesc format, void* data,
// 									stride_t xstride, stride_t ystride, stride_t zstride);
bool ImageInput_read_image_floats(ImageInput *in, float* data);
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

// std::string metadata_val(const ImageIOParameter &p, bool human=false);
const char* ImageSpec_to_xml(ImageSpec *spec);
// void from_xml(const char *xml)
// bool valid_tile_range(int xbegin, int xend, int ybegin, int yend, int zbegin, int zend)

TypeDesc ImageSpec_channelformat(ImageSpec *spec, int chan);
// void ImageSpec_get_channelformats(ImageSpec *spec, std::vector< TypeDesc > &formats);

// Properties
int ImageSpec_x(ImageSpec *spec);
void ImageSpec_set_x(ImageSpec *spec, int val);
int ImageSpec_y(ImageSpec *spec);
void ImageSpec_set_y(ImageSpec *spec, int val);
int ImageSpec_z(ImageSpec *spec);
void ImageSpec_set_z(ImageSpec *spec, int val);
int ImageSpec_width(ImageSpec *spec);
void ImageSpec_set_width(ImageSpec *spec, int val);
int ImageSpec_height(ImageSpec *spec);
void ImageSpec_set_height(ImageSpec *spec, int val);
int ImageSpec_depth(ImageSpec *spec);
void ImageSpec_set_depth(ImageSpec *spec, int val);
int ImageSpec_full_x(ImageSpec *spec);
void ImageSpec_set_full_x(ImageSpec *spec, int val);
int ImageSpec_full_y(ImageSpec *spec);
void ImageSpec_set_full_y(ImageSpec *spec, int val);
int ImageSpec_full_z(ImageSpec *spec);
void ImageSpec_set_full_z(ImageSpec *spec, int val);
int ImageSpec_full_width(ImageSpec *spec);
void ImageSpec_set_full_width(ImageSpec *spec, int val);
int ImageSpec_full_height(ImageSpec *spec);
void ImageSpec_set_full_height(ImageSpec *spec, int val);
int ImageSpec_full_depth(ImageSpec *spec);
void ImageSpec_set_full_depth(ImageSpec *spec, int val);
int ImageSpec_tile_width(ImageSpec *spec);
void ImageSpec_set_tile_width(ImageSpec *spec, int val);
int ImageSpec_tile_height(ImageSpec *spec);
void ImageSpec_set_tile_height(ImageSpec *spec, int val);
int ImageSpec_tile_depth(ImageSpec *spec);
void ImageSpec_set_tile_depth(ImageSpec *spec, int val);
int ImageSpec_nchannels(ImageSpec *spec);
void ImageSpec_set_nchannels(ImageSpec *spec, int val);
TypeDesc ImageSpec_format(ImageSpec *spec);
void ImageSpec_set_format(ImageSpec *spec, TypeDesc format);
void ImageSpec_channelformats(ImageSpec *spec, TypeDesc *out);
void ImageSpec_set_channelformats(ImageSpec *spec, TypeDesc *formats);
void ImageSpec_channelnames(ImageSpec *spec, char** out);
void ImageSpec_set_channelnames(ImageSpec *spec, char** names);
int ImageSpec_alpha_channel(ImageSpec *spec);
void ImageSpec_set_alpha_channel(ImageSpec *spec, int val);
int ImageSpec_z_channel(ImageSpec *spec);
void ImageSpec_set_z_channel(ImageSpec *spec, int val);
bool ImageSpec_deep(ImageSpec *spec);
void ImageSpec_set_deep(ImageSpec *spec, bool val);

void ImageSpec_attribute_type_data(ImageSpec *spec, const char* name, TypeDesc type, const void *value);
void ImageSpec_attribute_type_char(ImageSpec *spec, const char* name, TypeDesc type, const char* value);
void ImageSpec_attribute_uint(ImageSpec *spec, const char* name, unsigned int value);
void ImageSpec_attribute_int(ImageSpec *spec, const char* name, int value);
void ImageSpec_attribute_float(ImageSpec *spec, const char* name, float value);
void ImageSpec_attribute_char(ImageSpec *spec, const char* name, const char* value);
int ImageSpec_get_int_attribute(ImageSpec *spec, const char* name, int defaultval);
float ImageSpec_get_float_attribute(ImageSpec *spec, const char* name, float defaultval);
const char* ImageSpec_get_string_attribute(ImageSpec *spec, const char* name, const char* defaultval);
// void erase_attribute(const char* name, TypeDesc searchtype=TypeDesc::UNKNOWN, bool casesensitive=false)
// ImageIOParameter * find_attribute(const char* name, TypeDesc searchtype=TypeDesc::UNKNOWN, bool casesensitive=false)
// const ImageIOParameter * find_attribute(const char* name, TypeDesc searchtype=TypeDesc::UNKNOWN, bool casesensitive=false);



// ImageBuf
//

ImageBuf* ImageBuf_New();
ImageBuf* ImageBuf_New_WithCache(const char* name, ImageCache *imagecache);
ImageBuf* ImageBuf_New_WithBuffer(const char* name, const ImageSpec* spec, void *buffer);
ImageBuf* ImageBuf_New_SubImage(const char* name, int subimage, int miplevel, ImageCache* imagecache);
ImageBuf* ImageBuf_New_Spec(const ImageSpec* spec);

void ImageBuf_clear(ImageBuf* buf);
void ImageBuf_reset_subimage(ImageBuf* buf, const char* name, int subimage, int miplevel, ImageCache *imagecache);
void ImageBuf_reset_name_cache(ImageBuf* buf, const char* name, ImageCache *imagecache);
void ImageBuf_reset_spec(ImageBuf* buf, ImageSpec* spec);
void ImageBuf_reset_name_spec(ImageBuf* buf, const char* name, const ImageSpec* spec);

IBStorage ImageBuf_storage(ImageBuf* buf);
bool ImageBuf_initialized(ImageBuf* buf);
bool ImageBuf_read(ImageBuf* buf, int subimage, int miplevel, bool force, TypeDesc convert, void *cbk_data);
bool ImageBuf_init_spec(ImageBuf* buf, const char* filename, int subimage, int miplevel);
bool ImageBuf_write_file(ImageBuf* buf, const char* filename, const char* fileformat, void *cbk_data);
bool ImageBuf_write_output(ImageBuf* buf, ImageOutput *out, void *cbk_data);
void ImageBuf_set_write_format(ImageBuf* buf, TypeDesc format);
void ImageBuf_set_write_tiles(ImageBuf* buf, int width, int height, int depth);
void ImageBuf_copy_metadata(ImageBuf* dst, const ImageBuf* src);
bool ImageBuf_copy_pixels(ImageBuf* dst, const ImageBuf* src);
bool ImageBuf_copy(ImageBuf* dst, const ImageBuf* src);
void ImageBuf_swap(ImageBuf* buf, ImageBuf* other);
const char* ImageBuf_geterror(ImageBuf* buf);
const ImageSpec* ImageBuf_spec(ImageBuf* buf);
ImageSpec* ImageBuf_specmod(ImageBuf* buf);
const ImageSpec* ImageBuf_nativespec(ImageBuf* buf);
const char* ImageBuf_name(ImageBuf* buf);
const char* ImageBuf_file_format_name(ImageBuf* buf);
int ImageBuf_subimage(ImageBuf* buf);
int ImageBuf_nsubimages(ImageBuf* buf);
int ImageBuf_miplevel(ImageBuf* buf);
int ImageBuf_nmiplevels(ImageBuf* buf);
int ImageBuf_nchannels(ImageBuf* buf);
// float ImageBuf_getchannel(ImageBuf* buf, int x, int y, int z, int c, WrapMode wrap);
// void ImageBuf_getpixel(ImageBuf* buf, int x, int y, float *pixel, int maxchannels);
// void ImageBuf_getpixel_xyz(ImageBuf* buf, int x, int y, int z, float *pixel, int maxchannels, WrapMode wrap);
// void ImageBuf_interppixel(ImageBuf* buf, float x, float y, float *pixel, WrapMode wrap);
// void ImageBuf_interppixel_NDC(ImageBuf* buf, float s, float t, float *pixel, WrapMode wrap);
// void ImageBuf_interppixel_NDC_full(ImageBuf* buf, float s, float t, float *pixel, WrapMode wrap);
// void ImageBuf_setpixel(ImageBuf* buf, int x, int y, const float *pixel, int maxchannels);
// void ImageBuf_setpixel_xyz(ImageBuf* buf, int x, int y, int z, const float *pixel, int maxchannels);
// void ImageBuf_setpixel_index(ImageBuf* buf, int i, const float *pixel, int maxchannels);
bool ImageBuf_get_pixel_channels(ImageBuf* buf, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend, int chbegin, int chend, TypeDesc format, void *result);
// bool ImageBuf_get_pixels(ImageBuf* buf, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend, TypeDesc format, void *result);

int ImageBuf_orientation(ImageBuf* buf);
int ImageBuf_oriented_width(ImageBuf* buf);
int ImageBuf_oriented_height(ImageBuf* buf);
int ImageBuf_oriented_x(ImageBuf* buf);
int ImageBuf_oriented_y(ImageBuf* buf);
int ImageBuf_oriented_full_width(ImageBuf* buf);
int ImageBuf_oriented_full_height(ImageBuf* buf);
int ImageBuf_oriented_full_x(ImageBuf* buf);
int ImageBuf_oriented_full_y(ImageBuf* buf);

int ImageBuf_xbegin(ImageBuf* buf);
int ImageBuf_xend(ImageBuf* buf);
int ImageBuf_ybegin(ImageBuf* buf);
int ImageBuf_yend(ImageBuf* buf);
int ImageBuf_zbegin(ImageBuf* buf);
int ImageBuf_zend(ImageBuf* buf);
int ImageBuf_xmin(ImageBuf* buf);
int ImageBuf_xmax(ImageBuf* buf);
int ImageBuf_ymin(ImageBuf* buf);
int ImageBuf_ymax(ImageBuf* buf);
int ImageBuf_zmin(ImageBuf* buf);
int ImageBuf_zmax(ImageBuf* buf);

void ImageBuf_set_full(ImageBuf* buf, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend);
// void ImageBuf_set_full_border(ImageBuf* buf, int xbegin, int xend, int ybegin, int yend, int zbegin, int zend, const float *bordercolor);

ROI* ImageBuf_roi(ImageBuf* buf);
ROI* ImageBuf_roi_full(ImageBuf* buf);
void ImageBuf_set_roi_full(ImageBuf* buf, ROI* newroi);

bool ImageBuf_pixels_valid(ImageBuf* buf);
TypeDesc ImageBuf_pixeltype(ImageBuf* buf);
// void* ImageBuf_localpixels(ImageBuf* buf);
// const void* ImageBuf_localpixels(ImageBuf* buf);
bool ImageBuf_cachedpixels(ImageBuf* buf);
ImageCache* ImageBuf_imagecache(ImageBuf* buf);
// void* ImageBuf_pixeladdr(ImageBuf* buf, int x, int y);
// void* ImageBuf_pixeladdr_z(ImageBuf* buf, int x, int y, int z);
bool ImageBuf_deep(ImageBuf* buf);
// int ImageBuf_deep_samples(ImageBuf* buf, int x, int y, int z);
// const void* ImageBuf_deep_pixel_ptr(ImageBuf* buf, int x, int y, int z, int c);
// float ImageBuf_deep_value(ImageBuf* buf, int x, int y, int z, int c, int s);
// DeepData* ImageBuf_deepdata(ImageBuf* buf);

// ROI
//
void deleteROI(ROI* roi);

ROI* ROI_New();
ROI* ROI_NewOptions(int xbeing, int xend, int ybegin, int yend, int zbegin, int zend, int chbegin, int chend);
ROI* ROI_Copy(const ROI *roi); 

bool ROI_defined(ROI* roi);
int ROI_width(ROI* roi);
int ROI_height(ROI* roi);
int ROI_depth(ROI* roi);
int ROI_nchannels(ROI* roi);
imagesize_t ROI_npixels(ROI* roi);

// Properties
int ROI_xbegin(ROI* roi);
void ROI_set_xbegin(ROI* roi, int val);
int ROI_xend(ROI* roi);
void ROI_set_xend(ROI* roi, int val);
int ROI_ybegin(ROI* roi);
void ROI_set_ybegin(ROI* roi, int val);
int ROI_yend(ROI* roi);
void ROI_set_yend(ROI* roi, int val);
int ROI_zbegin(ROI* roi);
void ROI_set_zbegin(ROI* roi, int val);
int ROI_zend(ROI* roi);
void ROI_set_zend(ROI* roi, int val);
int ROI_chbegin(ROI* roi);
void ROI_set_chbegin(ROI* roi, int val);
int ROI_chend(ROI* roi);
void ROI_set_chend(ROI* roi, int val);


// ImageCache
//
ImageCache* ImageCache_Create(bool shared);
void ImageCache_Destroy(ImageCache *x, bool teardown);

void ImageCache_clear(ImageCache *x);

// bool ImageCache_attribute(ImageCache *x, const char *name, TypeDesc type, const void *val);

// bool ImageCache_attribute_int(ImageCache *x, const char *name, int val);
// bool ImageCache_attribute_float(ImageCache *x, const char *name, float val);
// bool ImageCache_attribute_double(ImageCache *x, const char *name, double val);
// bool ImageCache_attribute_char(ImageCache *x, const char *name, const char **val);

// bool ImageCache_getattribute(ImageCache *x, const char *name, TypeDesc type, void *val);
// bool ImageCache_getattribute_int(ImageCache *x, const char *name, int *val);
// bool ImageCache_getattribute_float(ImageCache *x, const char *name, float *val);
// bool ImageCache_getattribute_double(ImageCache *x, const char *name, double *val);
// bool ImageCache_getattribute_char(ImageCache *x, const char *name, char **val);

// char* ImageCache_resolve_filename(ImageCache *x, const char *filename);

// bool ImageCache_get_image_info(ImageCache *x, char *filename, int subimage, int miplevel,
//                      			char *dataname, TypeDesc datatype, void *data);

// bool ImageCache_get_imagespec(ImageCache *x, char *filename, ImageSpec &spec,
//                             int subimage=0, int miplevel=0,
//                             bool native=false);

// const ImageSpec* ImageCache_imagespec(ImageCache *x, char *filename, int subimage, int miplevel, bool native);

// bool get_pixels(ImageCache *x, char *filename, int subimage, int miplevel,
//                          int xbegin, int xend, int ybegin, int yend,
//                          int zbegin, int zend,
//                          TypeDesc format, void *result);

// bool get_pixels(ImageCache *x, char *filename,
//                 int subimage, int miplevel, int xbegin, int xend,
//                 int ybegin, int yend, int zbegin, int zend,
//                 int chbegin, int chend, TypeDesc format, void *result,
//                 stride_t xstride=AutoStride, stride_t ystride=AutoStride,
//                 stride_t zstride=AutoStride);

// Tile* ImageCache_get_tile(ImageCache *x, char *filename, int subimage, int miplevel,
//                             int x, int y, int z);

// void ImageCache_release_tile(ImageCache *x, Tile *tile);
// const void* ImageCache_tile_pixels(ImageCache *x, Tile *tile, TypeDesc *format);
// bool ImageCache_add_file(ImageCache *x, char *filename, ImageInput::Creator creator);
// bool ImageCache_add_tile(ImageCache *x, char *filename, int subimage, int miplevel,
// 		                 int x, int y, int z, TypeDesc format, const void *buffer,
// 		                 stride_t xstride=AutoStride, stride_t ystride=AutoStride,
// 		                 stride_t zstride=AutoStride);
const char* ImageCache_geterror(ImageCache *x);
const char* ImageCache_getstats(ImageCache *x, int level);
void ImageCache_reset_stats(ImageCache *x);
void ImageCache_invalidate(ImageCache *x, const char *filename);
void ImageCache_invalidate_all(ImageCache *x, bool force);


#ifdef __cplusplus
}
#endif
#endif