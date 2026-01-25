package oiio

import (
	"testing"
)

// TestWorkflowReadModifyWrite tests the most common workflow:
// Read an image, modify pixels, and write it back out
func TestWorkflowReadModifyWrite(t *testing.T) {
	tmpDir := t.TempDir()

	// Read existing image
	buf, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err)
	}
	if !buf.Initialized() {
		t.Fatal("ImageBuf not initialized after read")
	}

	spec := buf.Spec()
	if spec.Width() != 128 || spec.Height() != 64 {
		t.Errorf("Expected 128x64, got %dx%d", spec.Width(), spec.Height())
	}

	// Modify image - fill with color
	err = Fill(buf, []float32{0.5, 0.5, 0.5})
	if err != nil {
		t.Fatalf("Fill failed: %s", err)
	}

	// Write to new file
	output := tmpDir + "/workflow_modified.png"

	err = buf.WriteFile(output, "png")
	if err != nil {
		t.Fatalf("WriteFile failed: %s", err)
	}

	// Verify output file exists and can be read
	bufCheck, err := NewImageBufPath(output)
	if err != nil {
		t.Fatalf("Failed to read output: %s", err)
	}
	if !bufCheck.Initialized() {
		t.Fatal("Output image not initialized")
	}
}

// TestWorkflowImageCacheMultipleReads tests ImageCache with multiple images
func TestWorkflowImageCacheMultipleReads(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Configure cache for performance
	cache.SetAttribute("max_memory_MB", 256.0)
	cache.SetAttribute("autotile", 64)

	// Read same image multiple times through cache (should be cached)
	for i := 0; i < 5; i++ {
		buf, err := NewImageBufPathCache(TEST_IMAGE, cache)
		if err != nil {
			t.Fatalf("Iteration %d: %s", i, err)
		}
		if !buf.Initialized() {
			t.Fatalf("Iteration %d: ImageBuf not initialized", i)
		}
	}

	// Check cache stats
	stats := cache.GetStats(1)
	if stats == "" {
		t.Error("Cache stats empty")
	}
}

// TestWorkflowCopyAndConvert tests copying with format conversion
func TestWorkflowCopyAndConvert(t *testing.T) {
	// Read source image
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err)
	}

	// Create destination with different channel count
	srcSpec := src.Spec()
	dstSpec := NewImageSpecSize(
		srcSpec.Width(),
		srcSpec.Height(),
		4, // RGBA instead of RGB
		TypeFloat,
	)

	dst, err := NewImageBufSpec(dstSpec)
	if err != nil {
		t.Fatal(err)
	}

	// Copy and add alpha channel
	chanOpts := &ChannelOpts{
		Order:  []int32{0, 1, 2, -1},
		Values: []float32{0, 0, 0, 1.0},
	}
	err = Channels(dst, src, 4, chanOpts)
	if err != nil {
		t.Fatalf("Channels conversion failed: %s", err)
	}

	// Verify output
	outSpec := dst.Spec()
	if outSpec.NumChannels() != 4 {
		t.Errorf("Expected 4 channels, got %d", outSpec.NumChannels())
	}
}

// TestWorkflowImagePipeline tests a real-world image processing pipeline:
// Read → Resize → Color Correct → Write
func TestWorkflowImagePipeline(t *testing.T) {
	// Step 1: Read source image
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err)
	}

	srcSpec := src.Spec()

	// Step 2: Resize to 256x128
	resizedSpec := NewImageSpecSize(256, 128, srcSpec.NumChannels(), TypeFloat)
	resized, err := NewImageBufSpec(resizedSpec)
	if err != nil {
		t.Fatal(err)
	}

	roi := NewROIRegion2D(0, 256, 0, 128)
	err = Resize(resized, src, AlgoOpts{ROI: roi})
	if err != nil {
		t.Fatalf("Resize failed: %s", err)
	}

	// Step 3: Adjust brightness (multiply by 1.2)
	brightened, err := NewImageBufSpec(resizedSpec)
	if err != nil {
		t.Fatal(err)
	}

	err = MulValues(brightened, resized, []float32{1.2, 1.2, 1.2})
	if err != nil {
		t.Fatalf("MulValues failed: %s", err)
	}

	// Step 4: Write final result
	output := t.TempDir() + "/workflow_pipeline.png"

	err = brightened.WriteFile(output, "png")
	if err != nil {
		t.Fatalf("WriteFile failed: %s", err)
	}

	// Verify output
	result, err := NewImageBufPath(output)
	if err != nil {
		t.Fatal(err)
	}

	resultSpec := result.Spec()
	if resultSpec.Width() != 256 || resultSpec.Height() != 128 {
		t.Errorf("Expected 256x128, got %dx%d", resultSpec.Width(), resultSpec.Height())
	}
}

// TestWorkflowMultipleFormats tests reading and writing different formats
func TestWorkflowMultipleFormats(t *testing.T) {
	// Read source
	src, err := NewImageBufPath(TEST_IMAGE)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	formats := []string{"tif", "jpg"}
	for _, format := range formats {
		output := tmpDir + "/workflow_format." + format

		err = src.WriteFile(output, format)
		if err != nil {
			t.Fatalf("WriteFile(%s) failed: %s", format, err)
		}

		// Verify can read back
		check, err := NewImageBufPath(output)
		if err != nil {
			t.Fatalf("Failed to read %s: %s", format, err)
		}
		if !check.Initialized() {
			t.Fatalf("Failed to initialize %s image", format)
		}
	}
}

// TestWorkflowImageCacheInvalidation tests cache invalidation patterns
func TestWorkflowImageCacheInvalidation(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	// Load image through cache
	buf1, err := NewImageBufPathCache(TEST_IMAGE, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !buf1.Initialized() {
		t.Fatal("ImageBuf not initialized")
	}

	// Invalidate specific file
	cache.Invalidate(TEST_IMAGE)

	// Load again (should re-read from disk)
	buf2, err := NewImageBufPathCache(TEST_IMAGE, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !buf2.Initialized() {
		t.Fatal("ImageBuf not initialized after invalidation")
	}

	// Invalidate all
	cache.InvalidateAll(true)

	// Load again
	buf3, err := NewImageBufPathCache(TEST_IMAGE, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !buf3.Initialized() {
		t.Fatal("ImageBuf not initialized after InvalidateAll")
	}
}

// TestWorkflowPixelAccess tests reading and modifying individual pixels
func TestWorkflowPixelAccess(t *testing.T) {
	// Create small test image
	spec := NewImageSpecSize(4, 4, 3, TypeFloat)
	buf, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err)
	}

	// Fill with known pattern
	err = Fill(buf, []float32{0.5, 0.5, 0.5})
	if err != nil {
		t.Fatal(err)
	}

	// Get pixels from entire image
	pixels, err := buf.GetFloatPixels()
	if err != nil {
		t.Fatalf("GetFloatPixels failed: %s", err)
	}

	// Verify pixel values
	for i := 0; i < len(pixels); i++ {
		if pixels[i] != 0.5 {
			t.Errorf("Pixel %d: expected 0.5, got %f", i, pixels[i])
			break
		}
	}
}

// TestWorkflowComposite tests compositing two images
func TestWorkflowComposite(t *testing.T) {
	// Create two RGBA images (Over requires alpha)
	spec := NewImageSpecSize(128, 64, 4, TypeFloat)

	bg, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err)
	}
	Fill(bg, []float32{0.2, 0.2, 0.2, 1.0})

	fg, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err)
	}
	Fill(fg, []float32{0.8, 0.8, 0.8, 0.5}) // 50% alpha

	// Composite using Over operation
	result, err := NewImageBufSpec(spec)
	if err != nil {
		t.Fatal(err)
	}

	err = Over(result, fg, bg)
	if err != nil {
		t.Fatalf("Over failed: %s", err)
	}

	// Write result
	output := t.TempDir() + "/workflow_composite.png"

	err = result.WriteFile(output, "png")
	if err != nil {
		t.Fatal(err)
	}
}

// TestWorkflowBatchProcessing tests processing multiple images efficiently
func TestWorkflowBatchProcessing(t *testing.T) {
	cache := CreateImageCache(false)
	defer cache.Destroy(true)

	cache.SetAttribute("max_memory_MB", 512.0)

	tmpDir := t.TempDir()

	// Simulate batch processing of the same image multiple times
	// (in real-world, these would be different images)
	outputs := make([]string, 3)
	for i := 0; i < 3; i++ {
		// Load through cache
		src, err := NewImageBufPathCache(TEST_IMAGE, cache)
		if err != nil {
			t.Fatalf("Batch %d load failed: %s", i, err)
		}

		// Process: resize to 64x32
		spec := NewImageSpecSize(64, 32, src.Spec().NumChannels(), TypeFloat)
		dst, err := NewImageBufSpec(spec)
		if err != nil {
			t.Fatal(err)
		}

		roi := NewROIRegion2D(0, 64, 0, 32)
		err = Resize(dst, src, AlgoOpts{ROI: roi})
		if err != nil {
			t.Fatalf("Batch %d resize failed: %s", i, err)
		}

		// Write output
		outputs[i] = tmpDir + "/workflow_batch_" + string(rune('0'+i)) + ".png"

		err = dst.WriteFile(outputs[i], "png")
		if err != nil {
			t.Fatalf("Batch %d write failed: %s", i, err)
		}
	}

	// Verify all outputs
	for i, output := range outputs {
		check, err := NewImageBufPath(output)
		if err != nil {
			t.Fatalf("Batch %d verify failed: %s", i, err)
		}
		checkSpec := check.Spec()
		if checkSpec.Width() != 64 || checkSpec.Height() != 32 {
			t.Errorf("Batch %d: expected 64x32, got %dx%d", i, checkSpec.Width(), checkSpec.Height())
		}
	}
}
