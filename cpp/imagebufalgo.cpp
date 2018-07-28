#include <OpenImageIO/imagebufalgo.h>
#include <OpenImageIO/color.h>

#include <stdint.h>
#include <string>

#include "oiio.h"
#include "color.h"

extern "C" {

bool zero(ImageBuf *dst, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::zero(*(static_cast<OIIO::ImageBuf*>(dst)),
									*(static_cast<OIIO::ROI*>(roi)),
									nthreads);
}

bool fill(ImageBuf *dst, const float *values, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::fill(*(static_cast<OIIO::ImageBuf*>(dst)),
									values,
									*(static_cast<OIIO::ROI*>(roi)),
									nthreads);
}

bool checker(ImageBuf *dst, int width, int height, int depth, const float *color1, const float *color2,
			  int xoffset, int yoffset, int zoffset, ROI* roi, int nthreads) 
{
	return OIIO::ImageBufAlgo::checker(
		*(static_cast<OIIO::ImageBuf*>(dst)),
		width, height, depth, 
		color1, color2,
		xoffset, yoffset, zoffset,
		*(static_cast<OIIO::ROI*>(roi)),
		nthreads);	
}

bool channels(ImageBuf *dst, const ImageBuf *src, int nchannels, const int32_t *channelorder,
			   const float *channelvalues, const char **newchannelnames,
			   bool shuffle_channel_names)
{
	std::vector<std::string> vec_names;

	if (nchannels > 0 && newchannelnames != NULL) {
		vec_names.assign(newchannelnames, newchannelnames+nchannels);
	}

	bool ok = OIIO::ImageBufAlgo::channels(*(static_cast<OIIO::ImageBuf*>(dst)),
											*(static_cast<const OIIO::ImageBuf*>(src)),
											nchannels, channelorder, channelvalues,
											&vec_names[0], shuffle_channel_names );
	return ok;
}

bool channel_append(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, 
	ROI* roi, int nthreads) 
{
	return OIIO::ImageBufAlgo::channel_append(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			*(static_cast<const OIIO::ImageBuf*>(B)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);		
}

bool flatten(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::flatten(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool crop(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::crop(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool cut (ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::cut(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool paste(ImageBuf *dst, int xbegin, int ybegin, int zbegin, int chbegin,
			const ImageBuf *src, ROI* srcroi, int nthreads) 
{
	return OIIO::ImageBufAlgo::paste(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			xbegin,
			ybegin,
			zbegin,
			chbegin,
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(srcroi)),
			nthreads);
}

bool flip(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::flip(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool flop(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::flop(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool transpose(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::transpose(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool add(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::add(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			*(static_cast<const OIIO::ImageBuf*>(B)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool add_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::add(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			B,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);		
}

bool add_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::add(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			B,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool sub(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::sub(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			*(static_cast<const OIIO::ImageBuf*>(B)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool sub_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::sub(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			B,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool sub_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::sub(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			B,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool mul(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::mul(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			*(static_cast<const OIIO::ImageBuf*>(B)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool mul_values(ImageBuf *dst, const ImageBuf *A, const float *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::mul(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			B,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool mul_value(ImageBuf *dst, const ImageBuf *A, float B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::mul(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			B,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool colorconvert(ImageBuf* dst, const ImageBuf* src,
	                const char* from, 
	                const char* to,
	                bool unpremult,
	                const char* context_key,
	                const char* context_value,
	                ColorConfig *colorconfig,
	                ROI* roi, int nthreads) 
{
	return OIIO::ImageBufAlgo::colorconvert(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			from, 
			to,
			unpremult,
			context_key,
			context_value,
			static_cast<OIIO::ColorConfig*>(colorconfig),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);

}

bool colorconvert_processor(ImageBuf *dst, const ImageBuf *src, const ColorProcessor *processor,
				   			bool unpremult, ROI* roi, int nthreads) 
{
	return OIIO::ImageBufAlgo::colorconvert(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			static_cast<const OIIO::ColorProcessor*>(processor),
			unpremult,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool unpremult(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::unpremult(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool premult(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::premult(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool is_constant_color(const ImageBuf *src, float *color, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::isConstantColor(
			*(static_cast<const OIIO::ImageBuf*>(src)),
			color,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool is_constant_channel(const ImageBuf *src, int channel, float val, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::isConstantChannel(
				*(static_cast<const OIIO::ImageBuf*>(src)),
				channel,
				val, 
				*(static_cast<OIIO::ROI*>(roi)),
				nthreads);		
}

bool is_monochrome(const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::isMonochrome(
				*(static_cast<const OIIO::ImageBuf*>(src)),
				*(static_cast<OIIO::ROI*>(roi)),
				nthreads);	
}

char* computePixelHashSHA1(const ImageBuf *src, const char *extrainfo,
						   ROI* roi, int blocksize, int nthreads) {
	std::string aHash = OIIO::ImageBufAlgo::computePixelHashSHA1(
							*(static_cast<const OIIO::ImageBuf*>(src)),
							extrainfo,
							*(static_cast<OIIO::ROI*>(roi)),
							blocksize,
							nthreads);

    char* c = (char*)(malloc(aHash.length()));
    strcpy(c, aHash.c_str());
    return c;
}

bool rotate90(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::rotate90(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool rotate180(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::rotate180(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool rotate270(ImageBuf *dst, const ImageBuf *src, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::rotate270(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool rotate(ImageBuf *dst, const ImageBuf *src, float angle, const char* filtername, 
				float filterwidth, bool recompute_roi, ROI *roi, int nthreads) 
{
	return OIIO::ImageBufAlgo::rotate(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			angle,
			filtername,
			filterwidth,
			recompute_roi,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool resize(ImageBuf *dst, const ImageBuf *src, const char *filtername,
			 float filterwidth, ROI* roi, int nthreads) 
{
	return OIIO::ImageBufAlgo::resize(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			filtername,
			filterwidth,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);

}

bool resample(ImageBuf *dst, const ImageBuf *src, bool interpolate, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::resample(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			interpolate,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool over(ImageBuf *dst, const ImageBuf *A, const ImageBuf *B, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::over(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(A)),
			*(static_cast<const OIIO::ImageBuf*>(B)),
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);	
}

bool render_text(ImageBuf *dst, int x, int y, const char *text, int fontsize,
				  const char *fontname, const float *textcolor) {

	if (fontsize <= 0) fontsize = 16;

	return OIIO::ImageBufAlgo::render_text(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			x, y, 
			OIIO::string_view(text), 
			fontsize, 
			OIIO::string_view(fontname),
			textcolor);
}

} // extern "C"


