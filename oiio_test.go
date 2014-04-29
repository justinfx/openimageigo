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

	defer tmpfile.Close()

	TEST_IMAGE = fmt.Sprintf("%s.png", tmpfile.Name())

	m := image.NewRGBA(image.Rect(0, 0, 128, 64))
	blue := color.RGBA{0, 0, 255, 255}
	gray := color.RGBA{128, 128, 128, 255}

	draw.Draw(m, m.Bounds(), &image.Uniform{gray}, image.ZP, draw.Src)

	r := m.Bounds().Inset(16)
	draw.Draw(m, r, &image.Uniform{blue}, image.ZP, draw.Over)

	png.Encode(tmpfile, m)

	os.Rename(tmpfile.Name(), TEST_IMAGE)
}

func createOutputFile() string {
	tmpfile, err := ioutil.TempFile("", "oiio_unittest_output_")
	if err != nil {
		panic(err.Error())
	}

	defer tmpfile.Close()

	name := fmt.Sprintf("%s.png", tmpfile.Name())
	os.Rename(tmpfile.Name(), name)
	return name
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}

func checkFataError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}
