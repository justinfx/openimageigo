package oiio

import (
	"testing"
)

func TestCreateImageCache(t *testing.T) {
	private := CreateImageCache(false)
	private.Destroy(true)

	CreateImageCache(true)

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

	cache.Close("test")
	cache.CloseAll()
}
