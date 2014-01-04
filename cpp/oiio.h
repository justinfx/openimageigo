#ifndef _OPENIMAGEIGO_OIIO_H_
#define _OPENIMAGEIGO_OIIO_H_

#include <stdbool.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef ptrdiff_t stride_t;

typedef void TypeDesc;
typedef void ImageSpec;
typedef void ImageInput;
typedef void ImageOutput;
typedef void DeepData;

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
// bool ImageInput_read_scanline(ImageInput *in, int y, int z, float* data);
// bool ImageInput_read_scanline_format(ImageInput *in, int y, int z, TypeDesc format, void* data, stride_t xstride);
// bool ImageInput_read_tile(ImageInput *in, int x, int y, int z, float* data);
// bool ImageInput_read_tile_format(ImageInput *in, int x, int y, int z, TypeDesc format, void* data, 
// 									stride_t xstride, stride_t ystride, stride_t zstride);
// bool ImageInput_read_image(ImageInput *in, float* data);

// // TODO: Progress Callback?
// bool ImageInput_read_image_format(ImageInput *in, TypeDesc format, void* data, 
// 									stride_t xstride, stride_t ystride, stride_t zstride);

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


#ifdef __cplusplus
}
#endif
#endif