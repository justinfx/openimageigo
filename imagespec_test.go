package oiio

import (
	"testing"
)

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

	spec := in.Spec()

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

}

func TestImageSpecStringAttribute(t *testing.T) {
	spec, err := getTestImageSpec()
	if err != nil {
		t.Error(err.Error())
	}

	e := "Expected value %s for attr FOO_TEST; got %s"

	expected := ""
	actual := spec.AttributeString("FOO_TEST")
	if actual != expected {
		t.Errorf(e, expected, actual)
	}

	expected = "TEST_DEFAULT"
	actual = spec.AttributeString("FOO_TEST", expected)
	if actual != expected {
		t.Errorf(e, expected, actual)
	}

	expected = "foo_value"
	spec.SetAttribute("FOO_TEST", expected)
	actual = spec.AttributeString("FOO_TEST")
	if actual != expected {
		t.Errorf(e, expected, actual)
	}
}

func TestImageSpecFloatAttribute(t *testing.T) {
	spec, err := getTestImageSpec()
	if err != nil {
		t.Error(err.Error())
	}

	e := "Expected value %s for attr FOO_TEST; got %s"

	var expected float32
	actual := spec.AttributeFloat("FOO_TEST")
	if actual != expected {
		t.Errorf(e, expected, actual)
	}

	expected = -1.5
	actual = spec.AttributeFloat("FOO_TEST", expected)
	if actual != expected {
		t.Errorf(e, expected, actual)
	}

	expected = 123.456
	spec.SetAttribute("FOO_TEST", expected)
	actual = spec.AttributeFloat("FOO_TEST")
	if actual != expected {
		t.Errorf(e, expected, actual)
	}
}

func TestImageSpecIntAttribute(t *testing.T) {
	spec, err := getTestImageSpec()
	if err != nil {
		t.Error(err.Error())
	}

	const name = `FOO_TEST`

	e := "Expected value %s for attr FOO_TEST; got %s"

	var expected int
	actual := spec.AttributeInt(name)
	if actual != expected {
		t.Errorf(e, expected, actual)
	}

	expected = -1
	actual = spec.AttributeInt(name, expected)
	if actual != expected {
		t.Errorf(e, expected, actual)
	}

	expected = 123
	spec.SetAttribute(name, expected)
	actual = spec.AttributeInt(name)
	if actual != expected {
		t.Errorf(e, expected, actual)
	}
}

func TestImageSpecEraseAttribute(t *testing.T) {
	spec, err := getTestImageSpec()
	if err != nil {
		t.Error(err.Error())
	}

	var number int = 123
	var str string = "test"

	spec.SetAttribute("INT_VALUE", number)
	spec.SetAttribute("STR_VALUE", str)

	if actual := spec.AttributeInt("INT_VALUE"); actual != number {
		t.Fatalf("INT_VALUE attribute not set properly. Expected %d, got %d", number, actual)
	}

	if actual := spec.AttributeString("STR_VALUE"); actual != str {
		t.Fatalf("STR_VALUE attribute not set properly. Expected %s, got %s", str, actual)
	}

	spec.EraseAttribute("StR_VaLUe", true) // case-sensitive

	if spec.AttributeString("STR_VALUE") == "" {
		t.Fatal("STR_VALUE was not expected to have been erased using case-sensitive option")
	}

	spec.EraseAttribute("StR_VaLUe", false) // case-insensitive

	if spec.AttributeString("STR_VALUE") != "" {
		t.Error("Expected STR_VALUE attribute to have been erased")
	}

	spec.EraseAttributeType("INT_VALUE", TypeDouble, false)

	if spec.AttributeInt("INT_VALUE") == 0 {
		t.Fatal("INT_VALUE was not expected to have been erased non-matching search type")
	}

	spec.EraseAttributeType("INT_VALUE", TypeInt, false)

	if spec.AttributeInt("INT_VALUE") != 0 {
		t.Error("Expected INT_VALUE attribute to have been erased")
	}
}

func getTestImageSpec() (*ImageSpec, error) {
	in, err := OpenImageInput(TEST_IMAGE)
	if err != nil {
		return nil, err
	}

	spec := in.Spec()
	return spec, nil
}
