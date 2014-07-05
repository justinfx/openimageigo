#include <OpenImageIO/imageio.h>

#include <string>

#include "oiio.h"


extern OIIO::TypeDesc fromTypeDesc(TypeDesc fmt);
extern TypeDesc toTypeDesc(OIIO::TypeDesc fmt);



extern "C" {

#include "_cgo_export.h"


void deleteImageOutput(ImageOutput *out) {
	delete static_cast<OIIO::ImageOutput*>(out);
}

ImageOutput* ImageOutput_Create() {
	return (ImageOutput*) OIIO::ImageOutput();
}

ImageOutput* ImageOutput_Create_filename(const char* filename, const char* plugin_searchpath) {
	std::string s_filename(filename);
	std::string s_path(plugin_searchpath);
	return (ImageOutput*) OIIO::ImageOutput::create(s_filename, s_path);
}

bool ImageOutput_open(ImageOutput *out, const char* name, ImageSpec* newspec, OpenMode mode) {
	std::string s_name(name);
	return static_cast<OIIO::ImageOutput*>(out)->open(s_name, *(static_cast<OIIO::ImageSpec*>(newspec)), mode);
}

// bool ImageOutput_open_subimages(ImageOutput* out, const char* name, int subimages, const ImageSpec* newspec);

const char* ImageOutput_geterror(ImageOutput *out) {
	return static_cast<OIIO::ImageOutput*>(out)->geterror().c_str();
}

const char* ImageOutput_format_name(ImageOutput *out) {
	return static_cast<OIIO::ImageOutput*>(out)->format_name();
}

const ImageSpec* ImageOutput_spec(ImageOutput *out) {
	const OIIO::ImageSpec *spec = &(static_cast<OIIO::ImageOutput*>(out)->spec());
	return (ImageSpec*) spec;
}

bool ImageOutput_supports(ImageOutput *out, const char* feature) {
	std::string s_feature(feature);
	return static_cast<OIIO::ImageOutput*>(out)->supports(s_feature);
}

bool ImageOutput_close(ImageOutput *out) {
	return static_cast<OIIO::ImageOutput*>(out)->close();
}

// bool ImageOutput_write_scanline(ImageOutput* out, int y, int z, TypeDesc format,
//                              	const void *data, stride_t xstride);

// bool ImageOutput_write_scanlines(ImageOutput* out, int ybegin, int yend, int z,
// 	                              TypeDesc format, const void *data,
// 	                              stride_t xstride, stride_t ystride);

// bool ImageOutput_write_tile(ImageOutput* out, int x, int y, int z, TypeDesc format,
// 	                         const void *data, stride_t xstride,
// 	                         stride_t ystride, stride_t zstride);

// bool ImageOutput_write_tiles(ImageOutput* out, int xbegin, int xend, int ybegin, int yend,
// 	                          int zbegin, int zend, TypeDesc format,
// 	                          const void *data, stride_t xstride,
// 	                          stride_t ystride, stride_t zstride);

// bool ImageOutput_write_rectangle(ImageOutput* out, int xbegin, int xend, int ybegin, int yend,
// 	                              int zbegin, int zend, TypeDesc format,
// 	                              const void *data, stride_t xstride,
// 	                              stride_t ystride, stride_t zstride);

bool ImageOutput_write_image(ImageOutput* out, TypeDesc format, const void *data, void *cbk_data)
{
	ProgressCallback cbk = &image_progress_callback;

	return static_cast<OIIO::ImageOutput*>(out)->write_image(
												fromTypeDesc(format),
												data,
												OIIO::AutoStride,
												OIIO::AutoStride,
												OIIO::AutoStride,
												cbk,
												cbk_data);
}

// bool ImageOutput_write_deep_scanlines(ImageOutput* out, int ybegin, int yend, int z,
//                                    		const DeepData* deepdata);

// bool ImageOutput_write_deep_tiles(ImageOutput* out, int xbegin, int xend, int ybegin, int yend,
// 	                               int zbegin, int zend, const DeepData* deepdata);

// bool ImageOutput_write_deep_image(ImageOutput* out, const DeepData* deepdata);

// bool ImageOutput_copy_image(ImageOutput* out, ImageInput *in);

// int ImageOutput_send_to_output(ImageOutput* out, const char *format, ...);
// int ImageOutput_send_to_client(ImageOutput* out, const char *format, ...);













// bool ImageOutput_read_image_floats(ImageOutput *in, float* data) {
// 	return static_cast<OIIO::ImageOutput*>(in)->read_image(data);
// }

// bool ImageOutput_read_image_format(ImageOutput *in, TypeDesc format, void* data, void* cbk_data)
// {
// 	ProgressCallback cbk = &image_progress_callback;

// 	return static_cast<OIIO::ImageOutput*>(in)->read_image(
// 												fromTypeDesc(format),
// 												data,
// 												OIIO::AutoStride,
// 												OIIO::AutoStride,
// 												OIIO::AutoStride,
// 												cbk,
// 												cbk_data);
// }

// bool ImageOutput_read_scanline_floats(ImageOutput *in, int y, int z, float* data) {
// 	return static_cast<OIIO::ImageOutput*>(in)->read_scanline(y, z, data);
// }



} // extern "C"


