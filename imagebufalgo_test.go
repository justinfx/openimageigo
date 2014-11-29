package oiio

import (
	"reflect"
	"testing"
)

func TestAlgoZero(t *testing.T) {
	buf := NewImageBuf()
	err := Zero(buf)
	if err == nil {
		t.Error("Expected passing nil buffer and nil ROI to raise error")
	}

	roi := NewROIRegion2D(0, 100, 0, 200)
	if err = Zero(buf, AlgoOpts{ROI: roi}); err != nil {
		t.Error(err.Error())
	}

	if buf.XBegin() != 0 || buf.XEnd() != 100 || buf.YBegin() != 0 || buf.YEnd() != 200 {
		t.Error("buf was not initialized to ROI")
	}

	err = Zero(buf, AlgoOpts{Threads: 1})
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
	if err = Fill(buf, expected); err != nil {
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

func TestAlgoChecker(t *testing.T) {
	spec := NewImageSpecSize(16, 16, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	dark := []float32{.1, .1, .1}
	light := []float32{.4, .4, .4}

	if err = Checker2D(buf, 4, 4, dark, light, 0, 0); err != nil {
		t.Fatal(err.Error())
	}

	roi := NewROIRegion2D(0, 1, 0, 1)
	roi.SetChannelsEnd(3)
	iface, err := buf.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual := iface.([]float32)
	if !reflect.DeepEqual(dark, actual) {
		t.Fatalf("Expected pixels %v; Got %v", dark, actual)
	}

	roi = NewROIRegion2D(14, 15, 0, 1)
	roi.SetChannelsEnd(3)
	iface, err = buf.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual = iface.([]float32)
	if !reflect.DeepEqual(light, actual) {
		t.Fatalf("Expected pixels %v; Got %v", light, actual)
	}
}

func TestAlgoChannels(t *testing.T) {
	// Create a source image
	src, err := NewImageBufSpec(NewImageSpecSize(32, 32, 4, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}

	fill := []float32{0.25, 0.5, 0.75, 1.0}
	if err = Fill(src, fill); err != nil {
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

	dst = NewImageBuf()
	expected_names = []string{"A2", "B2", "G2", "R2"}
	opts = ChannelOpts{NewNames: expected_names}
	if err = Channels(dst, rgb, 4, &opts); err != nil {
		t.Fatal(err.Error())
	}
	actual = dst.Spec().ChannelNames()
	if !reflect.DeepEqual(actual, expected_names) {
		t.Errorf("Expected names %v; got %v", expected_names, actual)
	}
}

func TestAlgoChannelAppend(t *testing.T) {
	// Create a source image
	rgb, err := NewImageBufSpec(NewImageSpecSize(32, 32, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}

	fill := []float32{1, 1, 1}
	if err = Fill(rgb, fill); err != nil {
		t.Fatal(err.Error())
	}

	a, err := NewImageBufSpec(NewImageSpecSize(32, 32, 1, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}

	fill = []float32{1}
	if err = Fill(a, fill); err != nil {
		t.Fatal(err.Error())
	}

	dst := NewImageBuf()

	if err = ChannelAppend(dst, rgb, a); err != nil {
		t.Fatal(err.Error())
	}

	// Has right number of channels
	if dst.NumChannels() != 4 {
		t.Errorf("Expected 4 channels; got %d", dst.NumChannels())
	}

	expected := []string{"R", "G", "B", "A"}
	actual := dst.Spec().ChannelNames()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected names %v; got %v", expected, actual)
	}
}

// TODO: Flatten does not seem to work as expected
//
// func TestAlgoFlatten(t *testing.T) {
// 	// Create a source image
// 	src, err := NewImageBufSpec(NewImageSpecSize(32, 32, 1, TypeFloat))
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
//
// 	fill := []float32{0.5}
// 	if err = Fill(src, fill); err != nil {
// 		t.Fatal(err.Error())
// 	}
//
// 	dst, _ := NewImageBufSpec(src.Spec())
// 	if err = Flatten(dst, src); err != nil {
// 		t.Fatal(err.Error())
// 	}
//
// 	// Check that the pixels for the region are correct
// 	firstPixel := NewROIRegion2D(0, 1, 0, 1)
// 	firstPixel.SetChannelsEnd(1)
// 	iface, err := dst.GetPixelRegion(firstPixel, TypeFloat)
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	pixels := iface.([]float32)
// 	if !reflect.DeepEqual(pixels, fill[:1]) {
// 		t.Errorf("Expected pixels %v;  got %v", fill[:1], pixels)
// 	}
// }

func TestAlgoCrop(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()
	roi := NewROIRegion2D(10, 60, 20, 40)

	if err = Crop(dst, src, AlgoOpts{ROI: roi}); err != nil {
		t.Fatal(err.Error())
	}

	expect_w := 50
	expect_h := 20
	actual_w := dst.Spec().Width()
	actual_h := dst.Spec().Height()

	if expect_w != actual_w {
		t.Errorf("Expected width %d, got %d", expect_w, actual_w)
	}
	if expect_h != actual_h {
		t.Errorf("Expected width %d, got %d", expect_h, actual_h)
	}

}

func TestAlgoFlipFlop(t *testing.T) {
	spec := NewImageSpecSize(16, 16, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	red := []float32{1, 0, 0}
	blue := []float32{0, 0, 1}

	if err = Checker2D(buf, 4, 4, red, blue, 0, 0); err != nil {
		t.Fatal(err.Error())
	}

	dst := NewImageBuf()

	// Flip
	if err = Flip(dst, buf); err != nil {
		t.Fatal(err.Error())
	}

	roi := NewROIRegion2D(0, 1, 0, 1)
	roi.SetChannelsEnd(3)
	iface, err := dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := blue
	pixels := iface.([]float32)
	if !reflect.DeepEqual(pixels, expected) {
		t.Errorf("Expected pixels %v; got %v", expected, pixels)
	}

	// Flop
	if err = Flop(dst, buf); err != nil {
		t.Fatal(err.Error())
	}

	iface, err = dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels = iface.([]float32)
	if !reflect.DeepEqual(pixels, expected) {
		t.Errorf("Expected pixels %v; got %v", expected, pixels)
	}

	// Flopflop
	if err = Flipflop(dst, buf); err != nil {
		t.Fatal(err.Error())
	}

	iface, err = dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected = red
	pixels = iface.([]float32)
	if !reflect.DeepEqual(pixels, expected) {
		t.Errorf("Expected pixels %v; got %v", expected, pixels)
	}
}

func TestAlgoTranspose(t *testing.T) {
	spec := NewImageSpecSize(2, 2, 1, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}
	// top-left
	Fill(buf, []float32{0.0}, AlgoOpts{ROI: NewROIRegion2D(0, 1, 0, 1)})
	// top-right
	Fill(buf, []float32{.25}, AlgoOpts{ROI: NewROIRegion2D(1, 2, 0, 1)})
	// bottom-left
	Fill(buf, []float32{.75}, AlgoOpts{ROI: NewROIRegion2D(0, 1, 1, 2)})
	// bottom-right
	Fill(buf, []float32{1.0}, AlgoOpts{ROI: NewROIRegion2D(1, 2, 1, 2)})

	dst := NewImageBuf()
	if err = Transpose(dst, buf); err != nil {
		t.Fatal(err.Error())
	}

	pixels, _ := dst.GetFloatPixels()
	table := []float32{0, .75, .25, 1}
	for i, expected := range table {
		actual := pixels[i]
		if actual != expected {
			t.Errorf("Expected value %f at index %d; got %f", expected, i, actual)
		}
	}
}

func TestAlgoColorConvert(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()

	err = ColorConvert(dst, src, "lnf", "srgb8", false)
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
	err = ColorConvertProcessor(dst, src, cp, false)
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
	err = Resize(dst, src, AlgoOpts{ROI: roi})
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
	err = Resample(dst, src, true, AlgoOpts{ROI: roi})
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
	if err = Fill(src, topExpected); err != nil {
		t.Fatal(err.Error())
	}

	// Destination is bigger than source
	dstSpec := NewImageSpecSize(32, 32, 3, TypeFloat)
	dst, err := NewImageBufSpec(dstSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	bottomExpected := []float32{0.5, 0.2, 0.2}
	if err = Fill(dst, bottomExpected); err != nil {
		t.Fatal(err.Error())
	}

	// Paste source into destination at top left
	if err = Paste2D(dst, src, 0, 0); err != nil {
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
