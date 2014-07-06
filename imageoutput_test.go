package oiio

import (
	"os"
	"testing"
)

func TestImageOutputCreate(t *testing.T) {
	// Open test image
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	outfile := createOutputFile()
	defer os.Remove(outfile)

	out, err := CreateImageOutput(outfile)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Aribitrary feature name, just to test
	out.Supports("xyz")

	actual := out.FormatName()
	if actual != "png" {
		t.Errorf("Expected FormatName 'png' but got %q", actual)
	}

	if err = out.Close(); err != nil {
		t.Fatal(err.Error())
	}

	// Re-open
	spec, _ := in.Spec()
	if err = out.Open(outfile, spec, OpenModeCreate); err != nil {
		t.Fatal(err.Error())
	}

	if actual = out.FormatName(); actual != "png" {
		t.Errorf("Expected FormatName 'png' but got %q", actual)
	}

}

func TestImageOutputWritePixels(t *testing.T) {
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err.Error())
	}

	pixels, err := src.GetPixels(TypeUint8)
	if err != nil {
		t.Fatal(err.Error())
	}

	outfile := createOutputFile()
	defer os.Remove(outfile)

	out, err := CreateImageOutput(outfile)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("Writing to output file: %q", outfile)
	spec := src.Spec()
	if err = out.Open(outfile, spec, OpenModeCreate); err != nil {
		t.Fatal(err.Error())
	}

	if err = out.WriteImageFormat(pixels, TypeUint8, nil); err != nil {
		t.Fatal(err.Error())
	}

	if err = out.Close(); err != nil {
		t.Fatal("Error closing file:", err.Error())
	}

}
