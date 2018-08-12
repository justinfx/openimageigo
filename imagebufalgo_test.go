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
	checkError(t, Zero(buf, AlgoOpts{ROI: roi}))

	if buf.XBegin() != 0 || buf.XEnd() != 100 || buf.YBegin() != 0 || buf.YEnd() != 200 {
		t.Error("buf was not initialized to ROI")
	}

	checkError(t, Zero(buf, AlgoOpts{Threads: 1}))
}

func TestAlgoFill(t *testing.T) {
	spec := NewImageSpecSize(32, 32, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := []float32{1.0, 0.7, 0.7}
	checkFatalError(t, Fill(buf, expected))

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

	checkFatalError(t, Checker2D(buf, 4, 4, dark, light, 0, 0))

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
	checkFatalError(t, Fill(src, fill))

	// Strip the alpha channel
	dst := NewImageBuf()

	checkFatalError(t, Channels(dst, src, 3, nil))

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

	checkFatalError(t, Channels(dst, src, 1, &opts))

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
	checkFatalError(t, Channels(dst, src, 4, &opts))

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
	checkFatalError(t, Channels(dst, rgb, 4, &opts))
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
	checkFatalError(t, Channels(dst, rgb, 4, &opts))
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

	// Has right number of channels
	if rgb.NumChannels() != 3 {
		t.Errorf("Expected 3 channels; got %d", rgb.NumChannels())
	}

	expected := []string{"R", "G", "B"}
	actual := rgb.Spec().ChannelNames()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected names %v; got %v", expected, actual)
	}

	spec := NewImageSpecSize(32, 32, 1, TypeFloat)
	spec.SetChannelNames([]string{"A"})
	a, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	fill = []float32{1}
	if err = Fill(a, fill); err != nil {
		t.Fatal(err.Error())
	}

	// Has right number of channels
	if a.NumChannels() != 1 {
		t.Errorf("Expected 1 channels; got %d", a.NumChannels())
	}

	expected = []string{"A"}
	actual = a.Spec().ChannelNames()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected names %v; got %v", expected, actual)
	}

	dst := NewImageBuf()

	checkFatalError(t, ChannelAppend(dst, rgb, a))

	// Has right number of channels
	if dst.NumChannels() != 4 {
		t.Errorf("Expected 4 channels; got %d", dst.NumChannels())
	}

	expected = []string{"R", "G", "B", "A"}
	actual = dst.Spec().ChannelNames()
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

	checkFatalError(t, Crop(dst, src, AlgoOpts{ROI: roi}))

	expect_x := roi.XBegin()
	expect_y := roi.YBegin()
	expect_w := roi.Width()
	expect_h := roi.Height()
	actual_w := dst.Spec().Width()
	actual_h := dst.Spec().Height()
	actual_x := dst.Spec().X()
	actual_y := dst.Spec().Y()

	if actual_x != expect_x || actual_y != expect_y {
		t.Errorf("Expected origin to be (%d,%d), got (%d,%d)",
			expect_x, expect_y, actual_x, actual_y)
	}
	if expect_w != actual_w {
		t.Errorf("Expected width %d, got %d", expect_w, actual_w)
	}
	if expect_h != actual_h {
		t.Errorf("Expected width %d, got %d", expect_h, actual_h)
	}

}

func TestAlgoCut(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()
	roi := NewROIRegion2D(10, 60, 20, 40)

	checkFatalError(t, Cut(dst, src, AlgoOpts{ROI: roi}))

	expect_w := 50
	expect_h := 20
	actual_w := dst.Spec().Width()
	actual_h := dst.Spec().Height()
	actual_x := dst.Spec().X()
	actual_y := dst.Spec().Y()

	if actual_x != 0 || actual_y != 0 {
		t.Errorf("Expected origin to be (0,0), got (%d,%d)", actual_x, actual_y)
	}
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

	checkFatalError(t, Checker2D(buf, 4, 4, red, blue, 0, 0))

	dst := NewImageBuf()

	// Flip
	checkFatalError(t, Flip(dst, buf))

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
	checkFatalError(t, Flop(dst, buf))

	iface, err = dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels = iface.([]float32)
	if !reflect.DeepEqual(pixels, expected) {
		t.Errorf("Expected pixels %v; got %v", expected, pixels)
	}

	// Rotate180
	checkFatalError(t, Rotate180(dst, buf))

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

func TestAlgoRotate(t *testing.T) {
	spec := NewImageSpecSize(16, 16, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	red := []float32{1, 0, 0}
	blue := []float32{0, 0, 1}

	checkFatalError(t, Checker2D(buf, 4, 4, red, blue, 0, 0))

	dst := NewImageBuf()

	checkFatalError(t, Rotate90(dst, buf))
	checkFatalError(t, Rotate180(dst, buf))
	checkFatalError(t, Rotate270(dst, buf))
	checkFatalError(t, Rotate(dst, buf, 60, "", 0, true))
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
	checkFatalError(t, Transpose(dst, buf))

	pixels, _ := dst.GetFloatPixels()
	table := []float32{0, .75, .25, 1}
	for i, expected := range table {
		actual := pixels[i]
		if actual != expected {
			t.Errorf("Expected value %f at index %d; got %f", expected, i, actual)
		}
	}
}

func TestAlgoColorAdd(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(1, 1, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}
	reset := func() { Fill(buf, []float32{0, .25, .75}) }
	dst := NewImageBuf()

	reset()
	checkFatalError(t, AddValue(dst, buf, .25))
	actual, _ := dst.GetFloatPixels()
	expected := []float32{.25, .5, 1}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	reset()
	checkFatalError(t, AddValues(dst, buf, []float32{.1, .2, .15}))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{.1, .45, .9}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	reset()
	b, _ := NewImageBufSpec(NewImageSpecSize(1, 1, 3, TypeFloat))
	Fill(b, []float32{.2, .2, .2})
	checkFatalError(t, Add(dst, buf, b))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{.2, .45, .95}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}
}

func TestAlgoColorSub(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(1, 1, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}
	reset := func() { Fill(buf, []float32{.5, .75, 1}) }
	dst := NewImageBuf()

	reset()
	checkFatalError(t, SubValue(dst, buf, .25))
	actual, _ := dst.GetFloatPixels()
	expected := []float32{.25, .5, .75}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	reset()
	checkFatalError(t, SubValues(dst, buf, []float32{.1, .2, .15}))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{.4, .55, .85}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	reset()
	b, _ := NewImageBufSpec(NewImageSpecSize(1, 1, 3, TypeFloat))
	Fill(b, []float32{.2, .2, .2})
	checkFatalError(t, Sub(dst, buf, b))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{.3, .55, .8}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}
}

func TestAlgoColorMul(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(1, 1, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}
	reset := func() { Fill(buf, []float32{.1, .2, .3}) }
	dst := NewImageBuf()

	reset()
	checkFatalError(t, MulValue(dst, buf, 2))
	actual, _ := dst.GetFloatPixels()
	expected := []float32{.2, .4, .6}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	reset()
	checkFatalError(t, MulValues(dst, buf, []float32{2, 3, .5}))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{.2, .6, .15}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	reset()
	b, _ := NewImageBufSpec(NewImageSpecSize(1, 1, 3, TypeFloat))
	Fill(b, []float32{.5, .25, 1})
	checkFatalError(t, Mul(dst, buf, b))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{.05, .05, .3}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}
}

func TestAlgoColorConvert(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Error(err.Error())
	}

	dst := NewImageBuf()

	checkError(t, ColorConvert(dst, src, "lnf", "srgb8", false))

	cfg, err := NewColorConfig()
	if err != nil {
		t.Fatal(err)
	}

	cp, err := cfg.CreateColorProcessor("lnf", "srgb8")
	if err != nil {
		t.Fatal(err.Error())
	}

	dst = NewImageBuf()
	checkError(t, ColorConvertProcessor(dst, src, cp, false))
}

func TestAlgoPremult(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(1, 1, 4, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}
	Fill(buf, []float32{1, .6, .4, .5})
	dst := NewImageBuf()

	checkFatalError(t, Premult(dst, buf))
	actual, _ := dst.GetFloatPixels()
	expected := []float32{.5, .3, .2, .5}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}

	dst, src := NewImageBuf(), dst
	checkFatalError(t, Unpremult(dst, src))
	actual, _ = dst.GetFloatPixels()
	expected = []float32{1, .6, .4, .5}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pixels %v, got %v", expected, actual)
	}
}

func TestAlgoConstantColor(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(16, 16, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}
	constants := []float32{.25, .5, 1}
	Fill(buf, constants)

	if !IsConstantColor(buf) {
		t.Error("Expected constant color. Got false.")
	}

	colors := ConstantColors(buf)
	if !reflect.DeepEqual(colors, constants) {
		t.Errorf("Expected %v pixels. Got %v", constants, colors)
	}

	roi := buf.ROI()
	roi.SetChannelsEnd(2)
	colors = ConstantColors(buf, AlgoOpts{ROI: roi})
	if !reflect.DeepEqual(colors, constants[:2]) {
		t.Errorf("Expected %v pixels. Got %v", constants, colors)
	}
}

func TestAlgoIsConstantChannel(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(16, 16, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}

	// Fill all channels with a constant color
	constants := []float32{.25, .5, .8}
	Fill(buf, constants)

	for i, val := range constants {
		if !IsConstantChannel(buf, i, val) {
			t.Errorf("Expected constant value %0.2f for %d channel. Got false.", val, i)
		}
	}

	// Fill part of the 3rd channel with another color
	roi := buf.ROI().Copy()
	roi.SetChannelsBegin(2)
	roi.SetChannelsEnd(3)
	roi.SetXEnd(8)
	roi.SetYEnd(8)
	constants2 := []float32{.1, .2, .3}
	Fill(buf, constants2, AlgoOpts{ROI: roi})

	// First 2 channels should be constant, still
	for i, val := range constants[:2] {
		if !IsConstantChannel(buf, i, val) {
			t.Errorf("Expected constant value %0.2f for %d channel. Got false.", val, i)
		}
	}

	// 3rd channel should not be constant anymore.
	// Check both the old and new fill values
	for _, v := range []float32{constants[2], constants2[2]} {
		if IsConstantChannel(buf, 2, v) {
			t.Error("Expected channel 2 to *not* be a constant value")
		}
	}
}

func TestAlgoIsMonochrome(t *testing.T) {
	buf, err := NewImageBufSpec(NewImageSpecSize(16, 16, 3, TypeFloat))
	if err != nil {
		t.Fatal(err.Error())
	}

	var val float32 = 0
	height := buf.OrientedHeight()
	roi := buf.ROI().Copy()

	for i := 1; i <= 10; i++ {
		val = float32(i) / float32(height)
		roi.SetYEnd(height / i)
		Fill(buf, []float32{val, val, val}, AlgoOpts{ROI: roi})
	}

	if !IsMonochrome(buf) {
		t.Error("Expected image to be monochrome. Got false.")
	}
}

func TestAlgoComputeHash(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	aHash := ComputePixelHashSHA1(src, "extra", 4)
	err = src.LastError()
	if err != nil {
		t.Fatal(err.Error())
	}

	aHash2 := ComputePixelHashSHA1(src, "extra", 4)
	err = src.LastError()
	if err != nil {
		t.Fatal(err.Error())
	}

	if aHash == "" || aHash2 == "" || aHash != aHash2 {
		t.Fatalf("Computed pixel hash %q != %q", aHash, aHash2)
	}
}

func TestAlgoResize(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	dst := NewImageBuf()
	roi := NewROIRegion2D(0, 64, 0, 32)

	checkError(t, Resize(dst, src, AlgoOpts{ROI: roi}))

	if dst.OrientedWidth() != 64 || dst.OrientedHeight() != 32 {
		t.Logf("Expected width/height == 64/32, but got %v/%v", dst.OrientedWidth(), dst.OrientedHeight())
	}

	checkError(t, ResizeFilter(dst, src, "lanczos3", 1.0, AlgoOpts{ROI: roi}))

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
	checkError(t, Resample(dst, src, true, AlgoOpts{ROI: roi}))

	if dst.OrientedWidth() != 64 || dst.OrientedHeight() != 32 {
		t.Logf("Expected width/height == 64/32, but got %v/%v", dst.OrientedWidth(), dst.OrientedHeight())
	}
}

func TestAlgoOver(t *testing.T) {
	srcSpec := NewImageSpecSize(16, 16, 4, TypeFloat)
	srcSpec.SetAlphaChannel(3)

	srcA, err := NewImageBufSpec(srcSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	srcB, err := NewImageBufSpec(srcSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	Fill(srcA, []float32{1.0, 0.0, 0.0, 0.5})
	Fill(srcB, []float32{0.0, 0.0, 1.0, 1.0})

	Premult(srcA, srcA)

	dst := NewImageBuf()
	err = Over(dst, srcA, srcB)
	if err != nil {
		t.Fatalf("Failed while performing Over operation: %s", err)
	}

	expected := []float32{.5, 0, .5}

	// Top left pixel
	roi := NewROIRegion2D(0, 1, 0, 1)
	roi.SetChannelsEnd(3)
	topIface, err := dst.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual := topIface.([]float32)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected pixels %#v; Got %#v", expected, actual)
	}
}

func ExampleOver() {
	srcSpec := NewImageSpecSize(16, 16, 4, TypeFloat)
	srcSpec.SetAlphaChannel(3)

	srcA, _ := NewImageBufSpec(srcSpec)
	srcB, _ := NewImageBufSpec(srcSpec)

	// Red with 50% alpha
	Fill(srcA, []float32{1.0, 0.0, 0.0, 0.5})
	// Blue with solid alpha
	Fill(srcB, []float32{0.0, 0.0, 1.0, 1.0})

	// Make sure the foreground image is premult
	Premult(srcA, srcA)

	dst := NewImageBuf()
	Over(dst, srcA, srcB)
}

func TestAlgoPaste2D(t *testing.T) {
	srcSpec := NewImageSpecSize(16, 16, 3, TypeFloat)
	src, err := NewImageBufSpec(srcSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	topExpected := []float32{1.0, 0.7, 0.7}
	checkFatalError(t, Fill(src, topExpected))

	// Destination is bigger than source
	dstSpec := NewImageSpecSize(32, 32, 3, TypeFloat)
	dst, err := NewImageBufSpec(dstSpec)
	if err != nil {
		t.Fatal(err.Error())
	}

	bottomExpected := []float32{0.5, 0.2, 0.2}
	checkFatalError(t, Fill(dst, bottomExpected))

	// Paste source into destination at top left
	checkFatalError(t, Paste2D(dst, src, 0, 0))

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

func TestAlgoRenderText(t *testing.T) {
	spec := NewImageSpecSize(256, 256, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	Fill(buf, []float32{0.3, 0.3, 0.3})

	err = RenderTextColor(buf, 0, spec.Height()/2, "UNITTEST", 18,
		FontNameDefault, []float32{1.0, .1, .2})

	if err != nil {
		t.Fatalf("Failed to render text: %s", err)
	}
	buf.WriteFile("/tmp/text.png", "png")
}

func TestAlgoRenderX(t *testing.T) {
	spec := NewImageSpecSize(256, 256, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err.Error())
	}

	Fill(buf, []float32{.9, .9, .9})

	err = RenderBox(buf, 50, 50, 100, 100, []float32{1.0, .1, .2}, true)
	if err != nil {
		t.Fatalf("Failed to render box: %s", err)
	}

	err = RenderLine(buf, 0, 50, 256, 50, []float32{.1, 1.0, .2}, false)
	if err != nil {
		t.Fatalf("Failed to render line: %s", err)
	}

	for i := 0; i < 10; i++ {
		err = RenderPoint(buf, 128+i, 128, []float32{.2, .1, 1.0})
		if err != nil {
			t.Fatalf("Failed to render point: %s", err)
		}
	}

	buf.WriteFile("/tmp/out.png", "png")
}
