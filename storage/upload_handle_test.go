package storage

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadHandle(t *testing.T) {
	// Setup temporary directory
	tempDir, err := os.MkdirTemp("", "storage_upload_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a buffer to write our multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add a file field
	filename := "test_upload.txt"
	fileContent := "This is a test upload file."
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.WriteString(part, fileContent)
	if err != nil {
		t.Fatal(err)
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create request
	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create UploadStorager
	us := &DefaultUploadStorage{
		DefaultStorage: &DefaultStorage{
			BasePath: tempDir,
		},
	}

	// Test UploadHandle
	ctx := context.Background()
	infos, err := UploadHandle(ctx, req, us, "file")
	if err != nil {
		t.Fatalf("UploadHandle failed: %v", err)
	}

	if len(infos) != 1 {
		t.Fatalf("Expected 1 uploaded file info, got %d", len(infos))
	}

	info := infos[0]
	if info.Filename() != filename {
		t.Errorf("Expected filename %s, got %s", filename, info.Filename())
	}

	// Verify file content on disk
	uploadedFilePath := filepath.Join(tempDir, filename)
	content, err := os.ReadFile(uploadedFilePath)
	if err != nil {
		t.Fatalf("Failed to read uploaded file: %v", err)
	}

	if string(content) != fileContent {
		t.Errorf("Expected content %q, got %q", fileContent, string(content))
	}
}

func TestUploadFileInfo(t *testing.T) {
	info := &uploadFileInfo{
		fullName: "/path/to/file.txt",
		filename: "file.txt",
		size:     1024,
		header:   make(textproto.MIMEHeader),
	}
	info.header.Set("Content-Type", "text/plain")

	if info.FullName() != "/path/to/file.txt" {
		t.Errorf("FullName() = %v, want %v", info.FullName(), "/path/to/file.txt")
	}
	if info.Filename() != "file.txt" {
		t.Errorf("Filename() = %v, want %v", info.Filename(), "file.txt")
	}
	if info.Size() != 1024 {
		t.Errorf("Size() = %v, want %v", info.Size(), 1024)
	}
	if info.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("Header() Content-Type = %v, want %v", info.Header().Get("Content-Type"), "text/plain")
	}
}
