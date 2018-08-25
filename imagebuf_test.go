package oiio

import (
	"fmt"
	"os"
	"testing"
)

func TestNewImageBuf(t *testing.T) {
	buf := NewImageBuf()
	if buf.Initialized() {
		t.Fatal("ImageBuf should not be considered initialized")
	}
}

func TestNewImageBufInitSpec(t *testing.T) {
	buf := NewImageBuf()
	if buf.Initialized() {
		t.Fatal("Expected ImageBuf not to be initialized")
	}

	err := buf.InitSpec(TEST_IMAGE, 0, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	spec := buf.Spec()
	width, height := spec.Width(), spec.Height()
	if width != 128 || height != 64 {
		t.Errorf("Expected ImageSpec width==128, height==64, got %v, %v", width, height)
	}
}

func TestImageBufReadImage(t *testing.T) {
	// Open New
	cache := CreateImageCache(true)

	buf, err := NewImageBufPathCache(TEST_IMAGE, cache)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !buf.Initialized() {
		t.Fatal("ImageBuf was not initialized")
	}

	err = buf.Read(true)
	if err != nil {
		t.Fatal(err.Error())
	}

	storage := buf.Storage()
	if storage == IBStorageUninitialized {
		t.Error("Got IBStorage value: Uninitialized")
	}

	if buf.Name() != TEST_IMAGE {
		t.Errorf("Expected name %v; got %v", TEST_IMAGE, buf.Name())
	}

	if buf.FileFormatName() != "png" {
		t.Errorf("Expected format png; got %v", buf.FileFormatName())
	}

	curr, total := buf.SubImage(), buf.NumSubImages()
	if curr != 0 && total != 1 {
		t.Errorf("Expected SubImage 0 and total of 1; got %v, %v", curr, total)
	}
	curr, total = buf.MipLevel(), buf.NumMipLevels()
	if curr != 0 && total != 1 {
		t.Errorf("Expected MipLevel 0 and total of 1; got %v, %v", curr, total)
	}

	if buf.NumChannels() != 3 {
		t.Errorf("Expected 3 channels in image, got %v", buf.NumChannels())
	}

	buf.Orientation()

	width, height := buf.OrientedWidth(), buf.OrientedHeight()
	if width != 128 || height != 64 {
		t.Errorf("Expected width==128, height==64, got %v, %v", width, height)
	}

	width, height = buf.OrientedFullWidth(), buf.OrientedFullHeight()
	if width != 128 || height != 64 {
		t.Errorf("Expected width==128, height==64, got %v, %v", width, height)
	}

	x, y := buf.OrientedX(), buf.OrientedY()
	if x != 0 || y != 0 {
		t.Errorf("Expected orientation (0,0), got (%v,%v)", x, y)
	}
	x, y = buf.OrientedFullX(), buf.OrientedFullY()
	if x != 0 || y != 0 {
		t.Errorf("Expected orientation (0,0), got (%v,%v)", x, y)
	}
	x, y = buf.XMin(), buf.YMin()
	if x != 0 || y != 0 {
		t.Errorf("Expected origin (0,0), got (%v,%v)", x, y)
	}
	x, y = buf.XMax(), buf.YMax()
	if x != width-1 || y != height-1 {
		t.Errorf("Expected x/y max (%v,%v), got (%v,%v)", width-1, height-1, x, y)
	}
}

func TestImageBufSpec(t *testing.T) {
	// Open New
	cache := CreateImageCache(true)

	buf, err := NewImageBufPathCache(TEST_IMAGE, cache)
	if err != nil {
		t.Fatal(err.Error())
	}

	spec := buf.Spec()
	width, height := spec.Width(), spec.Height()
	if width != 128 || height != 64 {
		t.Errorf("Expected ImageSpec width==128, height==64, got %v, %v", width, height)
	}

	specmod := buf.SpecMod()
	width, height = specmod.Width(), specmod.Height()
	if width != 128 || height != 64 {
		t.Errorf("Expected ImageSpec width==128, height==64, got %v, %v", width, height)
	}

	if spec.Format() != buf.PixelType() {
		t.Errorf("Expected TypeDesc %v, got %v", spec.Format(), buf.PixelType())
	}

	if !buf.PixelsValid() {
		t.Error("Expected ImageBuf.PixelsValid() == true")
	}

	if !buf.CachedPixels() {
		t.Error("Expected ImageBuf.CachedPixels() == true")
	}

	if cache.ptr != buf.ImageCache().ptr {
		t.Error("Expected ImageBuf.ImageCache() to return the same ImageCache as originally used")
	}

	names := specmod.ChannelNames()
	expected := []string{"B", "G", "A"}
	specmod.SetChannelNames(expected)
	specmod = buf.SpecMod()
	names = specmod.ChannelNames()
	if fmt.Sprintf("%v", names) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expected %v names; got %v", expected, names)
	}

	roi := buf.ROI()
	width, height = roi.Width(), roi.Height()
	if width != 128 || height != 64 {
		t.Errorf("Expected ROI width==128, height==64, got %v, %v", width, height)
	}

	roi = buf.ROIFull()
	width, height = roi.Width(), roi.Height()
	if width != 128 || height != 64 {
		t.Errorf("Expected ROI width==128, height==64, got %v, %v", width, height)
	}

	roi.SetXEnd(64)
	roi.SetYEnd(32)
	buf.SetROIFull(roi)
	roi = buf.ROIFull()
	width, height = roi.Width(), roi.Height()
	if width != 64 || height != 32 {
		t.Errorf("Expected ROI width==64, height==32, got %v, %v", width, height)
	}

	buf.SetROIFull(NewROIRegion2D(0, 128, 0, 128))
	otherROI := NewROIRegion2D(0, 1024, 0, 1024)
	if buf.ContainsROI(otherROI) {
		// https://github.com/OpenImageIO/oiio/pull/1996
		// ImageBuf::contains_roi() is broken and always returns true
		//t.Errorf("Expected ImageBuf %v to not contain %v", buf.ROI(), otherROI)
	}
	otherROI = NewROIRegion2D(10, 20, 10, 20)
	otherROI.SetChannelsEnd(3)
	if !buf.ContainsROI(otherROI) {
		t.Errorf("Expected ImageBuf %v to contain %v", buf.ROI(), otherROI)
	}

	if buf.Deep() {
		t.Error("Expected ImageBuf.Deep() == false")
	}
}

func TestImageBufReadImageCallbacks(t *testing.T) {
	// Open New
	buf, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	// With success callback
	//
	buf.Clear()
	var progress ProgressCallback = func(done float32) bool {
		// no cancel
		return false
	}

	err = buf.ReadCallback(true, &progress)
	if err != nil {
		t.Fatal(err.Error())
	}

	// With stop callback
	//
	buf.Clear()
	progress = func(done float32) bool {
		// cancel
		return true
	}

	err = buf.ReadCallback(true, &progress)
	if err != nil {
		t.Fatal(err.Error())
	}

}

func TestImageBufWriteFile(t *testing.T) {
	// Open New
	buf, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	outfile := createOutputFile()
	defer os.Remove(outfile)

	buf.SetWriteTiles(0, 0, 0)
	buf.SetWriteFormat(TypeUint8)

	checkFatalError(t, buf.WriteFile(outfile, ""))

	info, _ := os.Stat(outfile)
	if info.Size() == 0 {
		t.Fatal("Image file was 0 bytes. Write failed")
	}

	checkFatalError(t, buf.WriteFile(outfile, "png"))

	var progress ProgressCallback = func(done float32) bool {
		// no cancel
		return false
	}

	checkFatalError(t, buf.WriteFileProgress(outfile, "", &progress))
}

func TestImageBufCopySwap(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	other := NewImageBuf()
	checkFatalError(t, other.Copy(src))
	checkFatalError(t, other.CopyMetadata(src))
	checkFatalError(t, other.CopyPixels(src))
	checkFatalError(t, other.Swap(src))
}

func TestImageBufGetPixels(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	var pixel_iface interface{}

	pixel_iface, err = src.GetPixels(TypeFloat)
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

	expected := src.OrientedWidth() * src.OrientedHeight() * src.Spec().Depth() * src.NumChannels()
	actual := len(float_pixels)
	if expected != actual {
		t.Fatalf("Expected to get a total of %d pixels; Got %d", expected, actual)
	}
}

func TestImageBufGetFloatPixels(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	pixels, err := src.GetFloatPixels()
	if err != nil {
		t.Fatal(err.Error())
	}

	if pixels[0] == 0 {
		t.Fatal("First pixel of test image was 0")
	}

	expected := src.OrientedWidth() * src.OrientedHeight() * src.Spec().Depth() * src.NumChannels()
	actual := len(pixels)
	if expected != actual {
		t.Fatalf("Expected to get a total of %d pixels; Got %d", expected, actual)
	}
}

func TestImageBufGetPixelRegion(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	all_pixels, err := src.GetFloatPixels()
	if err != nil {
		t.Fatal(err.Error())
	}

	roi := src.ROI()
	roi.SetXEnd(10)
	roi.SetYEnd(10)
	pixel_iface, err := src.GetPixelRegion(roi, TypeFloat)
	if err != nil {
		t.Fatal(err.Error())
	}

	some_pixels := pixel_iface.([]float32)

	spec := src.Spec()
	expected := 10 * 10 * spec.Depth() * spec.NumChannels()
	actual := len(some_pixels)
	if expected != actual {
		t.Fatalf("Expected to get a total of %d pixels; Got %d", expected, actual)
	}

	for i := 0; i < len(some_pixels); i++ {
		if some_pixels[i] != all_pixels[i] {
			t.Fatalf("Pixel values do not match; Expected %v. got %v", all_pixels[i], some_pixels[i])
		}
	}
}
