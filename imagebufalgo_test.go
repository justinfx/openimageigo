package oiio

import (
	"reflect"
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

func TestAlgoFill(t *testing.T) {
	spec := NewImageSpecSize(32, 32, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := []float32{1.0, 0.7, 0.7}
	if err = Fill(buf, expected, nil, GlobalThreads); err != nil {
		t.Fatal(err.Error())
	}

	roi := NewROIRegion2D(0, 1, 0, 1)
	roi.SetChannelsEnd(3)
	iface, err := buf.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual := iface.([]float32)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected pixels %v; Got %v", expected, actual)
	}
}

func TestAlgoChannels(t *testing.T) {
	// Create a source image
	src, err := NewImageBufSpec(NewImageSpecSize(32, 32, 4, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}

	fill := []float32{0.25, 0.5, 0.75, 1.0}
	if err = Fill(src, fill, nil, 0); err != nil {
		t.Fatal(err.Error())
	}

	// Strip the alpha channel
	dst := NewImageBuf()

	if err = Channels(dst, src, 3, nil); err != nil {
		t.Fatal(err.Error())
	}

	// Has right number of channels
	if dst.NumChannels() != 3 {
		t.Errorf("Expected 3 channels; got %d", dst.NumChannels())
	}

	rgb := dst

	// Check that the pixels for the region are correct
	roi := NewROIRegion2D(0, 1, 0, 1)
	roi.SetChannelsEnd(3)

	iface, err := dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels := iface.([]float32)
	if !reflect.DeepEqual(pixels, fill[0:3]) {
		t.Errorf("Expected pixels %v;  got %v", fill[0:3], pixels)
	}

	// Create a one channel alpha image
	dst = NewImageBuf()
	opts := ChannelOpts{Order: []int32{3}}

	if err = Channels(dst, src, 1, &opts); err != nil {
		t.Fatal(err.Error())
	}

	// Has right number of channels
	if dst.NumChannels() != 1 {
		t.Errorf("Expected 1 channels; got %d", dst.NumChannels())
	}

	// Test reordering channels
	dst = NewImageBuf()
	opts = ChannelOpts{Order: []int32{2, 1, 0}}
	if err = Channels(dst, src, 3, &opts); err != nil {
		t.Fatal(err.Error())
	}

	iface, err = dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels = iface.([]float32)
	expected := []float32{fill[2], fill[1], fill[0]}
	if !reflect.DeepEqual(pixels, expected) {
		t.Errorf("Expected pixels %v; got %v", expected, pixels)
	}

	// Set explicit values on the channels
	dst = NewImageBuf()
	opts = ChannelOpts{
		Order:  []int32{0, 1, -1, -1},
		Values: []float32{0, 0, 1.0, 0.5},
	}
	if err = Channels(dst, src, 4, &opts); err != nil {
		t.Fatal(err.Error())
	}

	roi.SetChannelsEnd(4)
	iface, err = dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected = []float32{0.25, 0.5, 1.0, 0.5}
	pixels = iface.([]float32)
	if !reflect.DeepEqual(pixels, expected) {
		t.Errorf("Expected pixels %v; got %v", expected, pixels)
	}

	// Rename channels
	dst = NewImageBuf()
	opts = ChannelOpts{
		Order:    []int32{0, 1, 2, -1},
		Values:   []float32{0, 0, 0, 1.0},
		NewNames: []string{"", "", "", "A"},
	}
	if err = Channels(dst, rgb, 4, &opts); err != nil {
		t.Fatal(err.Error())
	}
	if dst.NumChannels() != 4 {
		t.Errorf("Expected 4 channels; got %d", dst.NumChannels())
	}
	expected_names := []string{"R", "G", "B", "A"}
	actual := dst.Spec().ChannelNames()
	if !reflect.DeepEqual(actual, expected_names) {
		t.Errorf("Expected names %v; got %v", expected_names, actual)
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

func TestAlgoPaste2D(t *testing.T) {
	srcSpec := NewImageSpecSize(16, 16, 3, TypeFloat)
	src, err := NewImageBufSpec(srcSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	topExpected := []float32{1.0, 0.7, 0.7}
	if err = Fill(src, topExpected, nil, GlobalThreads); err != nil {
		t.Fatal(err.Error())
	}

	// Destination is bigger than source
	dstSpec := NewImageSpecSize(32, 32, 3, TypeFloat)
	dst, err := NewImageBufSpec(dstSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	bottomExpected := []float32{0.5, 0.2, 0.2}
	if err = Fill(dst, bottomExpected, nil, GlobalThreads); err != nil {
		t.Fatal(err.Error())
	}

	// Paste source into destination at top left
	if err = Paste2D(dst, src, 0, 0, nil, GlobalThreads); err != nil {
		t.Fatal(err.Error())
	}

	// Top left pixel
	roi := NewROIRegion2D(0, 1, 0, 1)
	roi.SetChannelsEnd(3)
	topIface, err := dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Bottom right pixel
	roi = NewROIRegion2D(31, 32, 31, 32)
	roi.SetChannelsEnd(3)
	bottomIface, err := dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	topActual := topIface.([]float32)
	bottomActual := bottomIface.([]float32)

	if !reflect.DeepEqual(topExpected, topActual) {
		t.Fatalf("Expected pixels %v; Got %v", topExpected, topActual)
	}

	if !reflect.DeepEqual(bottomExpected, bottomActual) {
		t.Fatalf("Expected pixels %v; Got %v", bottomExpected, bottomActual)
	}
}
