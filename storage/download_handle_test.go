package storage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadHandle_Comprehensive(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "storage_download_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filename := "test_download.txt"
	content := "Download me!"
	fullPath := filepath.Join(tempDir, filename)
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	ds := NewDefaultDownloadStorage("attachment")
	ds.BasePath = tempDir

	// Test Case 1: Basic Download
	t.Run("Basic Download", func(t *testing.T) {
		rr := httptest.NewRecorder()

		err := DownloadHandle(context.Background(), rr, ds, filename)
		if err != nil {
			t.Fatalf("DownloadHandle failed: %v", err)
		}

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if rr.Body.String() != content {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), content)
		}

		// Check Content-Disposition
		expectedCD := "attachment; filename=\"test_download.txt\""
		if cd := rr.Header().Get("Content-Disposition"); cd != expectedCD {
			t.Errorf("Content-Disposition = %v, want %v", cd, expectedCD)
		}
	})

	// Test Case 2: Download with Range (using DownloadHandleWithRange)
	t.Run("Download With Range", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/download", nil)
		req.Header.Set("Range", "bytes=0-7") // "Download"
		rr := httptest.NewRecorder()

		err := DownloadHandleWithRange(context.Background(), req, rr, ds, filename)
		if err != nil {
			t.Fatalf("DownloadHandleWithRange failed: %v", err)
		}

		if status := rr.Code; status != http.StatusPartialContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPartialContent)
		}

		expectedBody := "Download"
		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
		}

		// Check Content-Range
		// Total length is 12 ("Download me!")
		// Range 0-7 is 8 bytes.
		expectedCR := "bytes 0-7/*"
		if cr := rr.Header().Get("Content-Range"); cr != expectedCR {
			t.Errorf("Content-Range = %v, want %v", cr, expectedCR)
		}
	})

	// Test Case 3: Download Suffix Range
	t.Run("Download Suffix Range", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/download", nil)
		req.Header.Set("Range", "bytes=-3") // "me!"
		rr := httptest.NewRecorder()

		err := DownloadHandleWithRange(context.Background(), req, rr, ds, filename)
		if err != nil {
			t.Fatalf("DownloadHandleWithRange failed: %v", err)
		}

		if status := rr.Code; status != http.StatusPartialContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPartialContent)
		}

		expectedBody := "me!"
		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
		}
	})

	// Test Case 4: Download Single Byte Range (0-0)
	t.Run("Download Single Byte Range", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/download", nil)
		req.Header.Set("Range", "bytes=0-0") // "D"
		rr := httptest.NewRecorder()

		err := DownloadHandleWithRange(context.Background(), req, rr, ds, filename)
		if err != nil {
			t.Fatalf("DownloadHandleWithRange failed: %v", err)
		}

		if status := rr.Code; status != http.StatusPartialContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPartialContent)
		}

		expectedBody := "D"
		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
		}

		// Check Content-Range
		expectedCR := "bytes 0-0/*"
		if cr := rr.Header().Get("Content-Range"); cr != expectedCR {
			t.Errorf("Content-Range = %v, want %v", cr, expectedCR)
		}
	})
}

func TestDefaultDownloadStorage_DispositionType(t *testing.T) {
	ds := NewDefaultDownloadStorage("inline")
	if ds.DispositionType() != "inline" {
		t.Errorf("DispositionType() = %v, want %v", ds.DispositionType(), "inline")
	}
}
