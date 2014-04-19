package oiio

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

var TEST_IMAGE string

func init() {
	tmpfile, err := ioutil.TempFile("", "oiio_unittest_")
	if err != nil {
		panic(err.Error())
	}

	TEST_IMAGE = fmt.Sprintf("%s.png", tmpfile.Name())

	m := image.NewRGBA(image.Rect(0, 0, 128, 64))
	blue := color.RGBA{0, 0, 255, 255}
	gray := color.RGBA{128, 128, 128, 255}

	draw.Draw(m, m.Bounds(), &image.Uniform{gray}, image.ZP, draw.Src)

	r := m.Bounds().Inset(16)
	draw.Draw(m, r, &image.Uniform{blue}, image.ZP, draw.Over)

	png.Encode(tmpfile, m)
	tmpfile.Close()

	os.Rename(tmpfile.Name(), TEST_IMAGE)
}

// ImageInput
//
func TestOpenImageInput(t *testing.T) {
	// Open New
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !in.ValidFile(TEST_IMAGE) {
		t.Errorf("Test image %q should have been a valid file", TEST_IMAGE)
	}

	// Aribitrary feature name, just to test
	in.Supports("xyz")

	actual := in.FormatName()
	if actual != "png" {
		t.Errorf("Expected FormatName 'png' but got %q", actual)
	}

	if err = in.Close(); err != nil {
		t.Fatal(err.Error())
	}

	// Re-open
	if err = in.Open(TEST_IMAGE); err != nil {
		t.Fatal(err.Error())
	}

	if actual = in.FormatName(); actual != "png" {
		t.Errorf("Expected FormatName 'png' but got %q", actual)
	}

}

func TestImageInputReadImage(t *testing.T) {
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Simple read
	var pixels []float32
	pixels, err = in.ReadImage()
	if err != nil {
		t.Fatal(err.Error())
	}
	if pixels[0] == 0 {
		t.Fatal("First pixel of test image was 0")
	}
}

func TestImageInputReadImageFormat(t *testing.T) {
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	var pixel_iface interface{}

	// nil Callback read
	//
	pixel_iface, err = in.ReadImageFormat(TypeFloat, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	float_pixels, ok := pixel_iface.([]float32)
	if !ok {
		t.Fatal("Interface could not be converted to a []float21")
	}

	if float_pixels[0] == 0 {
		t.Fatal("First pixel of test image was 0")
	}

	// With callback
	//
	var progress ProgressCallback = func(done float32) bool {
		// no cancel
		return false
	}

	pixel_iface, err = in.ReadImageFormat(TypeFloat, &progress)
	if err != nil {
		t.Fatal(err.Error())
	}

	float_pixels, _ = pixel_iface.([]float32)
	if float_pixels[0] == 0 {
		t.Fatal("First pixel of test image was 0")
	}

	// With callback
	//
	progress = func(done float32) bool {
		// cancel
		return true
	}

	pixel_iface, err = in.ReadImageFormat(TypeFloat, &progress)
	if err != nil {
		t.Fatal(err.Error())
	}

	float_pixels, _ = pixel_iface.([]float32)
	if float_pixels[0] != 0 {
		t.Fatal("First pixel of test image should be 0, since callback issued a cancel")
	}

}

func TestImageInputReadScanline(t *testing.T) {
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	var pixels []float32
	pixels, err = in.ReadScanline(0, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	if pixels[0] == 0 {
		t.Fatal("First pixel of test image was 0")
	}

}

// ImageSpec
//
func TestNewImageSpec(t *testing.T) {
	spec := NewImageSpec(TypeFloat)
	spec = NewImageSpecSize(512, 512, 3, TypeDouble)

	spec.SetFormat(TypeHalf)
	spec.DefaultChannelNames()

	expected := 2
	bytes := spec.ChannelBytes()
	if bytes != expected {
		t.Errorf("Expected 2, got %v", bytes)
	}

	bytes = spec.ChannelBytesChan(0, false)
	if bytes != expected {
		t.Errorf("Expected 2, got %v", bytes)
	}

	bytes = spec.ChannelBytesChan(1, true)
	if bytes != expected {
		t.Errorf("Expected 2, got %v", bytes)
	}

	expected = 6
	bytes = spec.PixelBytes(false)
	if bytes != expected {
		t.Errorf("Expected 6, got %v", bytes)
	}

	bytes = spec.PixelBytes(true)
	if bytes != expected {
		t.Errorf("Expected 6, got %v", bytes)
	}

	if bytes != expected {
		t.Errorf("Expected 6, got %v", bytes)
	}

	bytes = spec.PixelBytesChans(0, 3, false)
	if bytes != expected {
		t.Errorf("Expected 6, got %v", bytes)
	}

	bytes = spec.PixelBytesChans(0, 3, true)
	if bytes != expected {
		t.Errorf("Expected 6, got %v", bytes)
	}

	spec.TileBytes(true)
	spec.TilePixels()
	spec.ScanlineBytes(true)
	spec.ImageBytes(true)
	spec.ImagePixels()
	spec.SizeSafe()

	format := spec.ChannelFormat(0)
	if format != TypeHalf {
		t.Errorf("Expected TypeHalf (8), got %v", format)
	}
}

func TestImageSpecProperties(t *testing.T) {
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	spec, err := in.Spec()
	if err != nil {
		t.Fatal(err.Error())
	}
	if spec.Width() != 128 {
		t.Errorf("Expected 128;  got %v", spec.Width())
	}
	if spec.Height() != 64 {
		t.Errorf("Expected 64;  got %v", spec.Height())
	}
	if spec.NumChannels() != 3 {
		t.Errorf("Expected 3;  got %v", spec.NumChannels())
	}
	if spec.AlphaChannel() != -1 {
		t.Errorf("Expected alpha index to be -1;  got %v", spec.AlphaChannel())
	}
	if spec.Format() != TypeUint8 {
		t.Errorf("Expected data format to be TypeUint8; got %v", spec.Format())
	}

	actual := spec.ChannelNames()
	if len(actual) != 3 || actual[0] != "R" || actual[1] != "G" || actual[2] != "B" {
		t.Errorf("Expected channel nanes R,G,B; got %v", actual)
	}

	spec.ChannelFormats()
	spec.X()
	spec.Y()
	spec.Z()
	spec.Depth()
	spec.FullX()
	spec.FullY()
	spec.FullZ()
	spec.FullWidth()
	spec.FullHeight()
	spec.FullDepth()
	spec.TileWidth()
	spec.TileHeight()
	spec.TileDepth()
	spec.ZChannel()
	spec.Deep()
	spec.QuantBlack()
	spec.QuantWhite()
	spec.QuantMin()
	spec.QuantMax()

}
