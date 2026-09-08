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

func checkFatalError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGlobalAttributeInt(t *testing.T) {
	if !SetAttribute("threads", 2) {
		t.Error("SetAttribute(threads, 2) returned false")
	}
	val := AttributeInt("threads")
	if val != 2 {
		t.Errorf("AttributeInt(threads) = %d, want 2", val)
	}

	// restore the hardware default so other tests are unaffected
	SetAttribute("threads", 0)
}

func TestGlobalAttributeString(t *testing.T) {
	if !SetAttribute("plugin_searchpath", "/tmp/oiio_plugins") {
		t.Error("SetAttribute(plugin_searchpath, /tmp/oiio_plugins) returned false")
	}
	val := AttributeString("plugin_searchpath")
	if val != "/tmp/oiio_plugins" {
		t.Errorf("AttributeString(plugin_searchpath) = %q, want \"/tmp/oiio_plugins\"", val)
	}
}

func TestGlobalAttributeDefault(t *testing.T) {
	intVal := AttributeInt("nonexistent_int", 42)
	if intVal != 42 {
		t.Errorf("AttributeInt with default = %d, want 42", intVal)
	}

	floatVal := AttributeFloat("nonexistent_float", 3.14)
	if floatVal != 3.14 {
		t.Errorf("AttributeFloat with default = %f, want 3.14", floatVal)
	}

	stringVal := AttributeString("nonexistent_string", "default")
	if stringVal != "default" {
		t.Errorf("AttributeString with default = %q, want \"default\"", stringVal)
	}
}

func TestGlobalSetAttributeInvalidType(t *testing.T) {
	if SetAttribute("threads", []int{1, 2, 3}) {
		t.Error("SetAttribute with invalid type (slice) should return false")
	}
}
