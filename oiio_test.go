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
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	png.Encode(tmpfile, m)
	tmpfile.Close()

	os.Rename(tmpfile.Name(), TEST_IMAGE)

}

func TestNewImageInput(t *testing.T) {
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual := in.FormatName()
	if actual != "png" {
		t.Errorf("Expected FormatName 'png' but got %q", actual)
	}

	if err = in.Close(); err != nil {
		t.Fatal(err.Error())
	}
}
