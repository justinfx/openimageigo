#ifndef _OPENIMAGEIGO_IMAGEBUFALGO_H_
#define _OPENIMAGEIGO_IMAGEBUFALGO_H_

#include <stdint.h>

#include "oiio.h"
#include "color.h"

#ifdef __cplusplus
extern "C" {
#endif


bool zero(ImageBuf *dst, ROI* roi, int nthreads);

bool fill(ImageBuf *dst, const float *values, ROI* roi, int nthreads);

bool checker(ImageBuf *dst, int width, int height, int depth, const float *color1, const float *color2,
			  int xoffset, int yoffset, int zoffset, ROI* roi, int nthreads);

bool channels(ImageBuf *dst, const ImageBuf *src, int nchannels, const int32_t *channelorder,
			   const float *channelvalues, const char **newchannelnames, bool shuffle_channel_names);

bool channel_append(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads);

bool flatten(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool cut (ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool crop(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool paste(ImageBuf *dst, int xbegin, int ybegin, int zbegin, int chbegin,
			const ImageBuf *src, ROI* srcroi, int nthreads);

bool flip(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool flop(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool transpose(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

// bool circular_shift(ImageBuf *dst, const ImageBuf *src, int xshift, int yshift,
// 					 int zshift=0, ROI* roi, int nthreads);

// bool clamp(ImageBuf *dst, const ImageBuf *src, const float *min=NULL, const float *max=NULL,
// 			bool clampalpha01=false, ROI* roi, int nthreads);

// bool clamp(ImageBuf *dst, const ImageBuf *src, float min=-std::numeric_limits< float >::max(),
// 			float max=std::numeric_limits< float >::max(), bool clampalpha01=false, ROI* roi, int nthreads);

bool add(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads);

bool add_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads);

bool add_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads);

bool sub(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads);

bool sub_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads);

bool sub_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads);

bool mul(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads);

bool mul_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads);

bool mul_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads);

/* OIIO >= 1.5.x
// bool pow_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads);

// bool pow_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads);
*/

// bool channel_sum(ImageBuf *dst, const ImageBuf *src, const float *weights=NULL, ROI* roi, int nthreads);

// bool rangecompress(ImageBuf *dst, const ImageBuf *src, bool useluma=false, ROI* roi, int nthreads);

// bool rangeexpand(ImageBuf *dst, const ImageBuf *src, bool useluma=false, ROI* roi, int nthreads);

bool colorconvert(ImageBuf* dst, const ImageBuf* src,
	                const char* from, const char* to,
	                bool unpremult,
	                const char* context_key,
	                const char* context_value,
	                ColorConfig *colorconfig,
	                ROI* roi, int nthreads);

// bool ociolook(ImageBuf *dst, const ImageBuf *src, const char *looks, const char *from,
// 			   const char *to, bool unpremult=false, bool inverse=false, const char *key=NULL,
// 			   const char *value=NULL, ROI* roi, int nthreads);

// bool ociodisplay(ImageBuf *dst, const ImageBuf *src, const char *display, const char *view,
// 				  const char *from=NULL, const char *looks=NULL, bool unpremult=false,
// 				  const char *key=NULL, const char *value=NULL, ROI* roi, int nthreads);

bool colorconvert_processor(ImageBuf *dst, const ImageBuf *src, const ColorProcessor *processor,
				   			bool unpremult, ROI* roi, int nthreads);

// bool colorconvert(float *color, int nchannels, const ColorProcessor *processor, bool unpremult);

bool unpremult(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool premult(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

// bool computePixelStats(PixelStats *stats, const ImageBuf *src, ROI* roi, int nthreads);

// bool compare(const ImageBuf *A, const ImageBuf *B, float failthresh, float warnthresh,
// 			  CompareResults *result, ROI* roi, int nthreads);

// int compare_Yee(const ImageBuf *A, const ImageBuf *B, CompareResults *result, float luminance=100,
// 				 float fov=45, ROI* roi, int nthreads)
bool is_constant_color(const ImageBuf *src, float *color, ROI* roi, int nthreads);

bool is_constant_channel(const ImageBuf *src, int channel, float val, ROI* roi, int nthreads);

bool is_monochrome(const ImageBuf *src, ROI* roi, int nthreads);

// bool color_count(const ImageBuf *src, imagesize_t *count, int ncolors, const float *color,
// 				  const float *eps=NULL, ROI* roi, int nthreads);

// bool color_range_check(const ImageBuf *src, imagesize_t *lowcount, imagesize_t *highcount,
// 						imagesize_t *inrangecount, const float *low, const float *high, ROI* roi, int nthreads);

// ROI* nonzero_region(const ImageBuf *src, ROI* roi, int nthreads);

char* computePixelHashSHA1(const ImageBuf *src, const char *extrainfo,
						   ROI* roi, int blocksize, int nthreads);

// bool warp (ImageBuf *dst, const ImageBuf *src,
//                     const Imath::M33f &M,
//                     const char* filtername,
//                     float filterwidth,
//                     bool recompute_roi,
//                     ImageBuf::WrapMode wrap = ImageBuf::WrapDefault,
//                     ROI* roi, int nthreads);

// bool warp (ImageBuf *dst, const ImageBuf *src,
//                     const Imath::M33f &M,
//                     const Filter2D *filter,
//                     bool recompute_roi,
//                     ImageBuf::WrapMode wrap = ImageBuf::WrapDefault,
//                     ROI* roi, int nthreads);

// bool reorient(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool rotate90(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool rotate180(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool rotate270(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

bool rotate (ImageBuf *dst, const ImageBuf *src, float angle, const char* filtername, 
				float filterwidth, bool recompute_roi, ROI *roi, int nthreads);

// bool rotate (ImageBuf *dst, const ImageBuf *src, float angle, Filter2D *filter,
//              bool recompute_roi, ROI *rot, int nthreads);

// bool rotate (ImageBuf *dst, const ImageBuf *src, float angle, 
//              float center_x, float center_y, const char* filtername, 
//              float filterwidth, bool recompute_roi, ROI *roi, int nthreads);

// bool rotate (ImageBuf *dst, const ImageBuf *src, float angle, float center_x, 
// 			 float center_y, Filter2D *filter, bool recompute_roi = false, 
// 			 ROI *rot, int nthreads);

bool resize(ImageBuf *dst, const ImageBuf *src, const char *filtername,
			 float filterwidth, ROI* roi, int nthreads);

// bool resize(ImageBuf *dst, const ImageBuf *src, Filter2D *filter, ROI* roi, int nthreads);

bool resample(ImageBuf *dst, const ImageBuf *src, bool interpolate, ROI* roi, int nthreads);

// bool convolve(ImageBuf *dst, const ImageBuf *src, const ImageBuf *kernel, bool normalize=true,
// 			   ROI* roi, int nthreads);

// bool make_kernel(ImageBuf *dst, const char *name, float width, float height, float depth=1.0f, bool normalize=true);

// bool unsharp_mask(ImageBuf *dst, const ImageBuf *src, const char *kernel="gaussian", float width=3.0f,
// 				   float contrast=1.0f, float threshold=0.0f, ROI* roi, int nthreads);

// bool fft(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

// bool ifft(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

// bool fixNonFinite(ImageBuf *dst, const ImageBuf *src, NonFiniteFixMode mode=NONFINITE_BOX3,
// 				   int *pixelsFixed=NULL, ROI* roi, int nthreads);

// bool fillholes_pushpull(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads);

// bool from_IplImage(ImageBuf *dst, const IplImage *ipl, TypeDesc convert=TypeDesc::UNKNOWN);

// IplImage* to_IplImage(const ImageBuf *src);

// bool capture_image(ImageBuf *dst, int cameranum=0, TypeDesc convert=TypeDesc::UNKNOWN);

bool over(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads);

// bool zover(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, bool z_zeroisinf=false, ROI* roi, int nthreads);

bool render_text(ImageBuf *dst, int x, int y, const char *text, int fontsize,
				  const char *fontname, const float *textcolor);

bool render_box(ImageBuf *dst, int x1, int y1, int x2, int y2,
                const float* color, size_t ncolors, bool fill,
                ROI* roi, int nthreads);

bool render_line(ImageBuf *dst, int x1, int y1, int x2, int y2,
                 const float* color, size_t ncolors, bool skip_first_point,
                 ROI* roi, int nthreads);

bool render_point(ImageBuf *dst, int x, int y,
                  const float* color, size_t ncolors, ROI* roi, int nthreads);

// bool histogram(const ImageBuf *src, int channel, std::vector< imagesize_t > *histogram, int bins=256,
// 				float min=0, float max=1, imagesize_t *submin=NULL, imagesize_t *supermax=NULL, ROI* roi);

// bool histogram_draw(ImageBuf *dst, const std::vector< imagesize_t > *histogram);

// bool make_texture(MakeTextureMode mode, const ImageBuf *input, const char *outputfilename,
// 				   const ImageSpec *config, std::ostream *outstream=NULL);

// bool make_texture(MakeTextureMode mode, const char *filename, const char *outputfilename,
// 				   const ImageSpec *config, std::ostream *outstream=NULL);

// bool make_texture(MakeTextureMode mode, const std::vector< char > *filenames, const char *outputfilename,
// 				   const ImageSpec *config, std::ostream *outstream=NULL);

#ifdef __cplusplus
}
#endif
#endif