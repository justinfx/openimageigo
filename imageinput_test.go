package oiio

import (
	"testing"
)

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
