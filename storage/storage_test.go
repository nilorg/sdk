package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultStorage_Upload_Download_Exist_Remove(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "storage_test_crud")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	ds := NewDefaultStorage()
	ds.BasePath = tempDir

	filename := "test_file.txt"
	content := "Hello Storage!"
	ctx := context.Background()

	// 1. Test Upload
	reader := bytes.NewBufferString(content)
	fullName, err := ds.Upload(ctx, reader, filename)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}
	if fullName != filename {
		t.Errorf("Upload returned fullName %s, want %s", fullName, filename)
	}

	// Verify file exists on disk
	diskPath := filepath.Join(tempDir, filename)
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		t.Errorf("File not created on disk at %s", diskPath)
	}

	// 2. Test Exist
	exists, err := ds.Exist(ctx, filename)
	if err != nil {
		t.Fatalf("Exist failed: %v", err)
	}
	if !exists {
		t.Error("Exist returned false, want true")
	}

	// Test Exist for non-existent file
	exists, err = ds.Exist(ctx, "non_existent.txt")
	if err != nil {
		t.Fatalf("Exist failed for non-existent file: %v", err)
	}
	if exists {
		t.Error("Exist returned true for non-existent file, want false")
	}

	// 3. Test Download
	var buf bytes.Buffer
	info, err := ds.Download(ctx, &buf, filename)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	if buf.String() != content {
		t.Errorf("Download content = %q, want %q", buf.String(), content)
	}
	if info.Size() != int64(len(content)) {
		t.Errorf("Download info size = %d, want %d", info.Size(), len(content))
	}

	// 4. Test Remove
	err = ds.Remove(ctx, filename)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// Verify file removed
	if _, err := os.Stat(diskPath); !os.IsNotExist(err) {
		t.Errorf("File still exists on disk after Remove")
	}

	// Verify Exist returns false
	exists, _ = ds.Exist(ctx, filename)
	if exists {
		t.Error("Exist returned true after Remove, want false")
	}
}

func TestDefaultStorage_Upload_Mkdir(t *testing.T) {
	// Test that Upload creates directories if they don't exist
	tempDir, err := os.MkdirTemp("", "storage_test_mkdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	ds := NewDefaultStorage()
	ds.BasePath = tempDir

	filename := "subdir/nested/test.txt"
	content := "Nested file"
	ctx := context.Background()

	_, err = ds.Upload(ctx, bytes.NewBufferString(content), filename)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	diskPath := filepath.Join(tempDir, filename)
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		t.Errorf("Nested file not created at %s", diskPath)
	}
}

func TestDefaultStorage_Upload_Rename(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "storage_test_rename")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	ds := NewDefaultStorage()
	ds.BasePath = tempDir

	filename := "original.txt"
	newFilename := "renamed.txt"
	content := "Renamed content"

	ctx := context.Background()
	// Add Rename to context
	ctx = NewRenameContext(ctx, func(name string) string {
		if name == filename {
			return newFilename
		}
		return name
	})

	fullName, err := ds.Upload(ctx, bytes.NewBufferString(content), filename)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// Upload returns the *original* filename as fullName in current implementation?
	// Let's check implementation.
	// storage.go: fullName = filename (line 83)
	// Wait, line 57: if rename, ok := FromRenameContext(ctx); ok { filename = rename(filename) }
	// So filename is updated.
	// Then line 83: fullName = filename.
	// So it should return the NEW filename.

	if fullName != newFilename {
		t.Errorf("Upload returned fullName %s, want %s", fullName, newFilename)
	}

	// Check file exists with new name
	if _, err := os.Stat(filepath.Join(tempDir, newFilename)); os.IsNotExist(err) {
		t.Errorf("File not found with new name %s", newFilename)
	}
	// Check file does NOT exist with old name
	if _, err := os.Stat(filepath.Join(tempDir, filename)); !os.IsNotExist(err) {
		t.Errorf("File found with old name %s, should not exist", filename)
	}
}

func TestDefaultStorage_Upload_Error_Cleanup(t *testing.T) {
	// Test the cleanup logic we added earlier
	tempDir, err := os.MkdirTemp("", "storage_test_cleanup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	ds := NewDefaultStorage()
	ds.BasePath = tempDir

	filename := "fail_upload.txt"

	// Create a reader that fails
	failReader := &FailReader{}

	ctx := context.Background()
	_, err = ds.Upload(ctx, failReader, filename)
	if err == nil {
		t.Fatal("Expected Upload to fail, but it succeeded")
	}

	// Verify file does NOT exist (cleaned up)
	diskPath := filepath.Join(tempDir, filename)
	if _, err := os.Stat(diskPath); !os.IsNotExist(err) {
		t.Errorf("File %s should have been cleaned up after failed upload", diskPath)
	}
}

type FailReader struct{}

func (r *FailReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
