package oiio

import (
	"testing"
)

func TestAlgoZero(t *testing.T) {
	buf := NewImageBuf()
	roi := NewROI()
	err := Zero(buf, roi, GlobalThreads)
	if err == nil {
		t.Error("Expected passing nil buffer and nil ROI to raise error")
	}

	roi = NewROIRegion2D(0, 100, 0, 200)
	if err = Zero(buf, roi, GlobalThreads); err != nil {
		t.Error(err.Error())
	}

	if buf.XBegin() != 0 || buf.XEnd() != 100 || buf.YBegin() != 0 || buf.YEnd() != 200 {
		t.Error("buf was not initialized to ROI")
	}

	err = Zero(buf, nil, 1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlgoColorConvert(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()

	err = ColorConvert(dst, src, "lnf", "srgb8", false, nil, GlobalThreads)
	if err != nil {
		t.Error(err.Error())
	}

	cfg, err := NewColorConfig()
	if err != nil {
		t.Fatal(err)
	}

	cp, err := cfg.CreateColorProcessor("lnf", "srgb8")
	if err != nil {
		t.Fatal(err.Error())
	}

	dst = NewImageBuf()
	err = ColorConvertProcessor(dst, src, cp, false, nil, GlobalThreads)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlgoResize(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()
	roi := NewROIRegion2D(0, 64, 0, 32)
	err = Resize(dst, src, roi, GlobalThreads)
	if err != nil {
		t.Error(err.Error())
	}

	if dst.OrientedWidth() != 64 || dst.OrientedHeight() != 32 {
		t.Logf("Expected width/height == 64/32, but got %v/%v", dst.OrientedWidth(), dst.OrientedHeight())
	}
}

func TestAlgoResample(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()
	roi := NewROIRegion2D(0, 64, 0, 32)
	err = Resample(dst, src, true, roi, GlobalThreads)
	if err != nil {
		t.Error(err.Error())
	}

	if dst.OrientedWidth() != 64 || dst.OrientedHeight() != 32 {
		t.Logf("Expected width/height == 64/32, but got %v/%v", dst.OrientedWidth(), dst.OrientedHeight())
	}
}
