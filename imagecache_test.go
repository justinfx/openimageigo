package oiio

import (
	"testing"
)

func TestCreateImageCache(t *testing.T) {
	private := CreateImageCache(false)
	private.Clear()
	private.Destroy(true)

	shared := CreateImageCache(true)
	shared.Clear()
}

func TestImageCacheStats(t *testing.T) {
	cache := CreateImageCache(true)
	stats := cache.GetStats(1)
	if stats == "" {
		t.Error("GetStats() returned an empty string")
	}

	cache.ResetStats()
	cache.Invalidate("test")
	cache.InvalidateAll(true)
}

func TestImageCacheAttributeInt(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Test setting and getting int attributes
	if !cache.SetAttribute("max_open_files", 50) {
		t.Error("SetAttribute(max_open_files, 50) returned false")
	}
	val := cache.AttributeInt("max_open_files")
	if val != 50 {
		t.Errorf("AttributeInt(max_open_files) = %d, want 50", val)
	}

	// Test autotile attribute
	if !cache.SetAttribute("autotile", 64) {
		t.Error("SetAttribute(autotile, 64) returned false")
	}
	val = cache.AttributeInt("autotile")
	if val != 64 {
		t.Errorf("AttributeInt(autotile) = %d, want 64", val)
	}
}

func TestImageCacheAttributeFloat(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Test setting and getting float attributes
	if !cache.SetAttribute("max_memory_MB", float32(512.0)) {
		t.Error("SetAttribute(max_memory_MB, 512.0) returned false")
	}
	val := cache.AttributeFloat("max_memory_MB")
	if val != 512.0 {
		t.Errorf("AttributeFloat(max_memory_MB) = %f, want 512.0", val)
	}
}

func TestImageCacheAttributeString(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Test setting and getting string attributes
	if !cache.SetAttribute("searchpath", "/tmp/images") {
		t.Error("SetAttribute(searchpath, /tmp/images) returned false")
	}
	val := cache.AttributeString("searchpath")
	if val != "/tmp/images" {
		t.Errorf("AttributeString(searchpath) = %q, want \"/tmp/images\"", val)
	}
}

func TestImageCacheAttributeDefaults(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Test default values for non-existent attributes
	intVal := cache.AttributeInt("nonexistent_int", 42)
	if intVal != 42 {
		t.Errorf("AttributeInt with default = %d, want 42", intVal)
	}

	floatVal := cache.AttributeFloat("nonexistent_float", 3.14)
	if floatVal != 3.14 {
		t.Errorf("AttributeFloat with default = %f, want 3.14", floatVal)
	}

	stringVal := cache.AttributeString("nonexistent_string", "default")
	if stringVal != "default" {
		t.Errorf("AttributeString with default = %q, want \"default\"", stringVal)
	}
}

func TestImageCacheSetAttributeInvalidType(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Test that invalid types return false
	if cache.SetAttribute("max_memory_MB", []int{1, 2, 3}) {
		t.Error("SetAttribute with invalid type (slice) should return false")
	}
}
