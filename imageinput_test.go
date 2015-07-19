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

	checkFatalError(t, in.Close())

	// Re-open
	checkFatalError(t, in.Open(TEST_IMAGE))

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

func TestImageInputSubimage(t *testing.T) {
	filepath := `testdata/subimages.exr`
	in, err := OpenImageInput(filepath)
	checkFatalError(t, err)

	expected := 0
	actual := in.CurrentSubimage()
	if actual != expected {
		t.Fatalf("Expected subimage %d, got %d", expected, actual)
	}

	expected = 1
	ok := in.SeekSubimage(expected, nil)
	if !ok {
		t.Fatalf("Failed while seeking to subimage %d", expected)
	}
	actual = in.CurrentSubimage()
	if actual != expected {
		t.Fatalf("Expected subimage %d, got %d", expected, actual)
	}

	expected = 2
	newSpec := NewImageSpec(TypeUnknown)
	ok = in.SeekSubimage(expected, newSpec)
	if !ok {
		t.Fatalf("Failed while seeking to subimage %d", expected)
	}
	actual = in.CurrentSubimage()
	if actual != expected {
		t.Fatalf("Expected subimage %d, got %d", expected, actual)
	}
	nchans := newSpec.NumChannels()
	if nchans != 1 {
		t.Fatalf("Expected to find 1 channel in subimage %d, got %d", expected, nchans)
	}
	if newSpec.Format() != in.Spec().Format() {
		t.Fatalf("Expected subimage format %v, got %v", in.Spec().Format(), newSpec.Format())
	}
}

func TestImageInputMipLevel(t *testing.T) {
	filepath := `testdata/checker_mip.tx`
	in, err := OpenImageInput(filepath)
	checkFatalError(t, err)

	expected := 0
	actual := in.CurrentMipLevel()
	if actual != expected {
		t.Fatalf("Expected mip level %d, got %d", expected, actual)
	}

	expected = 1
	ok := in.SeekMipLevel(0, expected, nil)
	if !ok {
		t.Fatalf("Failed while seeking to Mip level %d", expected)
	}
	actual = in.CurrentMipLevel()
	if actual != expected {
		t.Fatalf("Expected mip level %d, got %d", expected, actual)
	}

	expected = 2
	newSpec := NewImageSpec(TypeUnknown)
	ok = in.SeekMipLevel(0, expected, newSpec)
	if !ok {
		t.Fatalf("Failed while seeking to mip level %d", expected)
	}
	actual = in.CurrentMipLevel()
	if actual != expected {
		t.Fatalf("Expected mip level %d, got %d", expected, actual)
	}

	expected = 8
	actual = 0
	for i := 0; in.SeekMipLevel(0, i, newSpec); i++ {
		actual++
	}
	if actual != expected {
		t.Fatalf("Expected total number of mip levels to be %d, got %d", expected, actual)
	}

}
