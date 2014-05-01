#include <OpenImageIO/imagebufalgo.h>
#include <OpenImageIO/color.h>

#include <string>

#include "oiio.h"
#include "color.h"

extern "C" {

bool zero(ImageBuf *dst, ROI* roi, int nthreads) {
	return OIIO::ImageBufAlgo::zero(*(static_cast<OIIO::ImageBuf*>(dst)),
									*(static_cast<OIIO::ROI*>(roi)),
									nthreads);
}

// bool fill(ImageBuf *dst, const float *values, ROI* roi, int nthreads) {
// 	return OIIO::ImageBufAlgo::fill(*(static_cast<OIIO::ImageBuf*>(dst)),
// 									values,
// 									*(static_cast<OIIO::ROI*>(roi)),
// 									nthreads);
// }

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


