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

	if err = buf.WriteFile(outfile, ""); err != nil {
		t.Fatal(err.Error())
	}

	info, _ := os.Stat(outfile)
	if info.Size() == 0 {
		t.Fatal("Image file was 0 bytes. Write failed")
	}

	if err = buf.WriteFile(outfile, "png"); err != nil {
		t.Fatal(err.Error())
	}

	var progress ProgressCallback = func(done float32) bool {
		// no cancel
		return false
	}

	if err = buf.WriteFileProgress(outfile, "", &progress); err != nil {
		t.Fatal(err.Error())
	}

}

func TestImageBufCopySwap(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	other := NewImageBuf()
	if err = other.Copy(src); err != nil {
		t.Fatal(err.Error())
	}

	if err = other.CopyMetadata(src); err != nil {
		t.Fatal(err.Error())
	}

	if err = other.CopyPixels(src); err != nil {
		t.Fatal(err.Error())
	}

	if err = other.Swap(src); err != nil {
		t.Fatal(err.Error())
	}
}
