package oiio

import (
	"testing"
)

func TestOpenImageOutput(t *testing.T) {
	out, err := OpenImageOutput("png")
	if err != nil {
		t.Fatal(err.Error())
	}

	out.Supports("xyz")

	actual := out.FormatName()
	if actual != "png" {
		t.Errorf("Expected FormatName 'png' but got actual %q", actual)
	}

}

func TestImageOutputSupports(t *testing.T) {

	// Test using .exr format as it has the widest support feature
	out, err := OpenImageOutput("exr")

	if err != nil {
		t.Fatal(err.Error())
	}

	expectSupports := []string{
		"tiles",
		"mipmap",
		"alpha",
		"nchannels",
		"channelformats",
		"displaywindow",
		"origin",
		"negativeorigin",
		"arbitrary_metadata",
		"exif",
		"iptc",
		"multiimage",
		"deepdata",
	}

	for _, expect := range expectSupports {
		if !out.Supports(expect) {
			t.Errorf("Expected support for feature %q (format %q)", expect, out.FormatName())
		}
	}

	if out.Supports("invalidfeature") {
		t.Error("Supports() returned true for a feature we expected to report false")
	}

}
