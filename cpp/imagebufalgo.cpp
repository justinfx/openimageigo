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

bool channels(ImageBuf *dst, const ImageBuf *src, int nchannels, const int32_t *channelorder,
			   const float *channelvalues, const char **newchannelnames,
			   bool shuffle_channel_names)
{
	std::string *str_names = NULL;

	if (nchannels > 0 && newchannelnames != NULL) {
		str_names = new(std::string[nchannels]);
		const char** temp = newchannelnames;
		for (size_t i=0; i < nchannels; i++) {
			str_names[i] = std::string(*temp++);
		}
	}

	bool ok = OIIO::ImageBufAlgo::channels(*(static_cast<OIIO::ImageBuf*>(dst)),
											*(static_cast<const OIIO::ImageBuf*>(src)),
											nchannels, channelorder, channelvalues,
											str_names, shuffle_channel_names );
	if (str_names != NULL)
		delete [] str_names;
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
				   			bool unpremult, ROI* roi, int nthreads) {

	return OIIO::ImageBufAlgo::colorconvert(
			*(static_cast<OIIO::ImageBuf*>(dst)),
			*(static_cast<const OIIO::ImageBuf*>(src)),
			static_cast<const OIIO::ColorProcessor*>(processor),
			unpremult,
			*(static_cast<OIIO::ROI*>(roi)),
			nthreads);
}

bool resize(ImageBuf *dst, const ImageBuf *src, const char *filtername,
			 float filterwidth, ROI* roi, int nthreads) {

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


