package oiio

import (
	"testing"
)

func TestNewROIEmpty(t *testing.T) {
	roi := NewROI()
	if roi.Defined() {
		t.Error("Expected empty ROI to returned Defined=false")
	}
}

func TestNewROIRegion2D(t *testing.T) {
	roi := NewROIRegion2D(100, 500, 50, 600)
	if !roi.Defined() {
		t.Error("Expected ROI to have a defined region")
	}

	actual := roi.Width()
	if actual != 400 {
		t.Errorf("Expected width 400, got %v", actual)
	}

	actual = roi.Height()
	if actual != 550 {
		t.Errorf("Expected height 550, got %v", actual)
	}

	actual = roi.Depth()
	if actual != 1 {
		t.Errorf("Expected depth 1, got %v", actual)
	}

	actual = roi.NumChannels()
	if actual != 1000 {
		t.Errorf("Expected 1000 channels, got %v", actual)
	}

	actual = roi.NumPixels()
	expected := (500 - 100) * (600 - 50) * 1
	if actual != expected {
		t.Errorf("Expected %v pixels, got %v", expected, actual)
	}
}

func TestNewROIRegion3D(t *testing.T) {
	roi := NewROIRegion3D(100, 500, 50, 600, 10, 200, 3, 10)
	if !roi.Defined() {
		t.Error("Expected ROI to have a defined region")
	}

	actual := roi.Width()
	if actual != 400 {
		t.Errorf("Expected width 400, got %v", actual)
	}

	actual = roi.Height()
	if actual != 550 {
		t.Errorf("Expected height 550, got %v", actual)
	}

	actual = roi.Depth()
	if actual != 190 {
		t.Errorf("Expected depth 1, got %v", actual)
	}

	actual = roi.NumChannels()
	if actual != 7 {
		t.Errorf("Expected 1000 channels, got %v", actual)
	}

	actual = roi.NumPixels()
	expected := (500 - 100) * (600 - 50) * (200 - 10)
	if actual != expected {
		t.Errorf("Expected %v pixels, got %v", expected, actual)
	}
}

func TestROIProperties(t *testing.T) {
	roi := NewROIRegion3D(100, 500, 50, 600, 10, 200, 3, 10)

	if roi.XBegin() != 100 {
		t.Errorf("Expected XBegin 100, got %v", roi.XBegin())
	}
	roi.SetXBegin(200)
	if roi.XBegin() != 200 {
		t.Errorf("Expected XBegin 200, got %v", roi.XBegin())
	}
	if roi.XEnd() != 500 {
		t.Errorf("Expected XEnd 500, got %v", roi.XEnd())
	}
	roi.SetXEnd(1000)
	if roi.XEnd() != 1000 {
		t.Errorf("Expected XEnd 1000, got %v", roi.XEnd())
	}
	if roi.YBegin() != 50 {
		t.Errorf("Expected YBegin 50, got %v", roi.YBegin())
	}
	roi.SetYBegin(100)
	if roi.YBegin() != 100 {
		t.Errorf("Expected YBegin 100, got %v", roi.YBegin())
	}
	if roi.YEnd() != 600 {
		t.Errorf("Expected YEnd 600, got %v", roi.YEnd())
	}
	roi.SetYEnd(1200)
	if roi.YEnd() != 1200 {
		t.Errorf("Expected YEnd 1200, got %v", roi.YEnd())
	}
	if roi.ZBegin() != 10 {
		t.Errorf("Expected ZBegin 10, got %v", roi.ZBegin())
	}
	roi.SetZBegin(20)
	if roi.ZBegin() != 20 {
		t.Errorf("Expected ZBegin 20, got %v", roi.ZBegin())
	}
	if roi.ZEnd() != 200 {
		t.Errorf("Expected ZEnd 200, got %v", roi.ZEnd())
	}
	roi.SetZEnd(400)
	if roi.ZEnd() != 400 {
		t.Errorf("Expected ZEnd 400, got %v", roi.ZEnd())
	}
	if roi.ChannelsBegin() != 3 {
		t.Errorf("Expected ChannelsBegin 3, got %v", roi.ChannelsBegin())
	}
	roi.SetChannelsBegin(6)
	if roi.ChannelsBegin() != 6 {
		t.Errorf("Expected ChannelsBegin 6, got %v", roi.ChannelsBegin())
	}
	if roi.ChannelsEnd() != 10 {
		t.Errorf("Expected ChannelsEnd 10, got %v", roi.ChannelsEnd())
	}
	roi.SetChannelsEnd(20)
	if roi.ChannelsEnd() != 20 {
		t.Errorf("Expected ChannelsEnd 20, got %v", roi.ChannelsEnd())
	}
}
