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

bool colorconvert(ImageBuf *dst, const ImageBuf *src, const char *from, const char *to,
				   bool unpremult, ROI* roi, int nthreads) {

	return OIIO::ImageBufAlgo::colorconvert(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			from,
			to,
			unpremult,
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

} // extern "C"


