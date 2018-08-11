#include <OpenImageIO/imagebuf.h>

#include "oiio.h"

extern "C" {

void deleteROI(ROI *roi) {
	delete static_cast<OIIO::ROI*>(roi);
}

ROI* ROI_New() {
	ROI* roi =  new OIIO::ROI();
	return (ROI*) roi;
}

ROI* ROI_NewOptions(int xbegin, int xend, int ybegin, int yend, int zbegin, int zend, int chbegin, int chend) {
	ROI* roi =  new OIIO::ROI(xbegin, xend, ybegin, yend, zbegin, zend, chbegin, chend);
	return (ROI*) roi;
}

ROI* ROI_Copy(const ROI *roi) {
	ROI* rc =  new OIIO::ROI(*(static_cast<const OIIO::ROI*>(roi)));
	return (ROI*) rc;
}

bool ROI_defined(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->defined();
}

int ROI_width(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->width();
}

int ROI_height(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->height();
}

int ROI_depth(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->depth();
}

int ROI_nchannels(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->nchannels();
}

imagesize_t ROI_npixels(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->npixels();
}

int ROI_xbegin(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->xbegin;
}

void ROI_set_xbegin(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->xbegin = val;
} 

int ROI_xend(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->xend;
}

void ROI_set_xend(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->xend = val;
} 

int ROI_ybegin(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->ybegin;
}

void ROI_set_ybegin(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->ybegin = val;
} 

int ROI_yend(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->yend;
}

void ROI_set_yend(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->yend = val;
} 

int ROI_zbegin(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->zbegin;
}

void ROI_set_zbegin(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->zbegin = val;
} 

int ROI_zend(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->zend;
}

void ROI_set_zend(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->zend = val;
} 

int ROI_chbegin(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->chbegin;
}

void ROI_set_chbegin(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->chbegin = val;
} 

int ROI_chend(ROI* roi) {
	return static_cast<OIIO::ROI*>(roi)->chend;
}

void ROI_set_chend(ROI* roi, int val) {
	static_cast<OIIO::ROI*>(roi)->chend = val;
}


} // extern "C"