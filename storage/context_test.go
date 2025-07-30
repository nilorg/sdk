package storage

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFromRangeRequestContext(t *testing.T) {
	tests := []struct {
		name      string
		start     int64
		end       int64
		expectOk  bool
		expectNil bool
	}{
		{
			name:     "valid range request",
			start:    0,
			end:      100,
			expectOk: true,
		},
		{
			name:     "valid range to end",
			start:    100,
			end:      0,
			expectOk: true,
		},
		{
			name:     "valid range negative end",
			start:    0,
			end:      -1,
			expectOk: true,
		},
		{
			name:      "invalid range - negative start",
			start:     -1,
			end:       100,
			expectOk:  false,
			expectNil: true,
		},
		{
			name:      "invalid range - end less than start",
			start:     100,
			end:       50,
			expectOk:  false,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// 创建包含Range请求的context
			ctxWithRange := NewRangeRequestContext(ctx, tt.start, tt.end)

			// 从context中获取Range请求
			rangeReq, ok := FromRangeRequestContext(ctxWithRange)

			if tt.expectNil {
				// 对于无效的Range请求，NewRangeRequestContext应该返回原context
				// 所以FromRangeRequestContext应该返回nil, false
				if ok {
					t.Errorf("expected ok=false for invalid range, got ok=true")
				}
				if rangeReq != nil {
					t.Errorf("expected rangeReq=nil for invalid range, got %+v", rangeReq)
				}
				return
			}

			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectOk, ok)
			}

			if !ok {
				return
			}

			if rangeReq == nil {
				t.Fatal("expected non-nil rangeReq")
			}

			if rangeReq.Start != tt.start {
				t.Errorf("expected start=%d, got start=%d", tt.start, rangeReq.Start)
			}

			if rangeReq.End != tt.end {
				t.Errorf("expected end=%d, got end=%d", tt.end, rangeReq.End)
			}
		})
	}
}

func TestFromRangeRequestContext_EmptyContext(t *testing.T) {
	ctx := context.Background()

	// 从空context中获取Range请求
	rangeReq, ok := FromRangeRequestContext(ctx)

	if ok {
		t.Errorf("expected ok=false for empty context, got ok=true")
	}

	if rangeReq != nil {
		t.Errorf("expected rangeReq=nil for empty context, got %+v", rangeReq)
	}
}

func TestRangeRequest_Methods(t *testing.T) {
	tests := []struct {
		name         string
		start        int64
		end          int64
		expectValid  bool
		expectToEnd  bool
		expectLength int64
	}{
		{
			name:         "valid range 0-99",
			start:        0,
			end:          99,
			expectValid:  true,
			expectToEnd:  false,
			expectLength: 100,
		},
		{
			name:         "valid range to end",
			start:        100,
			end:          0,
			expectValid:  true,
			expectToEnd:  true,
			expectLength: -1,
		},
		{
			name:         "valid range negative end",
			start:        0,
			end:          -1,
			expectValid:  true,
			expectToEnd:  true,
			expectLength: -1,
		},
		{
			name:         "invalid negative start",
			start:        -1,
			end:          100,
			expectValid:  false,
			expectToEnd:  false,
			expectLength: 102,
		},
		{
			name:         "invalid end less than start",
			start:        100,
			end:          50,
			expectValid:  false,
			expectToEnd:  false,
			expectLength: -49,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rangeReq := &RangeRequest{
				Start: tt.start,
				End:   tt.end,
			}

			if valid := rangeReq.IsValid(); valid != tt.expectValid {
				t.Errorf("expected IsValid()=%v, got %v", tt.expectValid, valid)
			}

			if toEnd := rangeReq.IsToEnd(); toEnd != tt.expectToEnd {
				t.Errorf("expected IsToEnd()=%v, got %v", tt.expectToEnd, toEnd)
			}

			if length := rangeReq.Length(); length != tt.expectLength {
				t.Errorf("expected Length()=%d, got %d", tt.expectLength, length)
			}
		})
	}
}

func TestDownloadWithRangeRequest(t *testing.T) {
	// 创建临时目录和测试文件
	tmpDir, err := os.MkdirTemp("", "storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testContent := "0123456789abcdefghijklmnopqrstuvwxyz"
	testFileName := "test.txt"
	testFilePath := filepath.Join(tmpDir, testFileName)
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建存储实例
	storage := &DefaultStorage{BasePath: tmpDir}

	tests := []struct {
		name            string
		start           int64
		end             int64
		expectedData    string
		expectedSize    int64
		shouldHaveRange bool
	}{
		{
			name:            "no range request",
			shouldHaveRange: false,
			expectedData:    testContent,
			expectedSize:    int64(len(testContent)),
		},
		{
			name:            "range request first 10 bytes",
			start:           0,
			end:             9,
			shouldHaveRange: true,
			expectedData:    "0123456789",
			expectedSize:    10,
		},
		{
			name:            "range request middle bytes",
			start:           10,
			end:             19,
			shouldHaveRange: true,
			expectedData:    "abcdefghij",
			expectedSize:    10,
		},
		{
			name:            "range request to end",
			start:           30,
			end:             0, // 0 表示到文件末尾
			shouldHaveRange: true,
			expectedData:    "uvwxyz",
			expectedSize:    6,
		},
		{
			name:            "range request with negative end",
			start:           30,
			end:             -1, // -1 表示到文件末尾
			shouldHaveRange: true,
			expectedData:    "uvwxyz",
			expectedSize:    6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.shouldHaveRange {
				ctx = NewRangeRequestContext(ctx, tt.start, tt.end)
			}

			var buf bytes.Buffer
			info, err := storage.Download(ctx, &buf, testFileName)
			if err != nil {
				t.Fatalf("Download failed: %v", err)
			}

			if info == nil {
				t.Fatal("Expected non-nil info")
			}

			actualData := buf.String()
			if actualData != tt.expectedData {
				t.Errorf("Expected data %q, got %q", tt.expectedData, actualData)
			}

			if info.Size() != tt.expectedSize {
				t.Errorf("Expected size %d, got %d", tt.expectedSize, info.Size())
			}
		})
	}
}

func TestParseRangeHeader(t *testing.T) {
	tests := []struct {
		name        string
		rangeHeader string
		expectStart int64
		expectEnd   int64
		expectError bool
	}{
		{
			name:        "bytes range with start and end",
			rangeHeader: "bytes=0-1023",
			expectStart: 0,
			expectEnd:   1023,
			expectError: false,
		},
		{
			name:        "bytes range to end",
			rangeHeader: "bytes=1024-",
			expectStart: 1024,
			expectEnd:   0,
			expectError: false,
		},
		{
			name:        "invalid format - missing bytes prefix",
			rangeHeader: "0-1023",
			expectError: true,
		},
		{
			name:        "invalid format - missing dash",
			rangeHeader: "bytes=0",
			expectError: true,
		},
		{
			name:        "invalid start position",
			rangeHeader: "bytes=abc-1023",
			expectError: true,
		},
		{
			name:        "invalid end position",
			rangeHeader: "bytes=0-xyz",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := parseRangeHeader(tt.rangeHeader)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if start != tt.expectStart {
				t.Errorf("Expected start=%d, got start=%d", tt.expectStart, start)
			}

			if end != tt.expectEnd {
				t.Errorf("Expected end=%d, got end=%d", tt.expectEnd, end)
			}
		})
	}
}

func TestDownloadHandleWithRange(t *testing.T) {
	// 创建临时目录和测试文件
	tmpDir, err := os.MkdirTemp("", "storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testContent := "0123456789abcdefghijklmnopqrstuvwxyz"
	testFileName := "test.txt"
	testFilePath := filepath.Join(tmpDir, testFileName)
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建存储实例
	storage := &DefaultDownloadStorage{
		DefaultStorage:  &DefaultStorage{BasePath: tmpDir},
		dispositionType: "attachment",
	}

	tests := []struct {
		name           string
		rangeHeader    string
		expectedStatus int
		expectedData   string
	}{
		{
			name:           "no range header",
			rangeHeader:    "",
			expectedStatus: http.StatusOK,
			expectedData:   testContent,
		},
		{
			name:           "range request first 10 bytes",
			rangeHeader:    "bytes=0-9",
			expectedStatus: http.StatusPartialContent,
			expectedData:   "0123456789",
		},
		{
			name:           "range request to end",
			rangeHeader:    "bytes=30-",
			expectedStatus: http.StatusPartialContent,
			expectedData:   "uvwxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/download/"+testFileName, nil)
			if tt.rangeHeader != "" {
				req.Header.Set("Range", tt.rangeHeader)
			}

			w := httptest.NewRecorder()
			ctx := context.Background()

			err := DownloadHandleWithRange(ctx, req, w, storage, testFileName)
			if err != nil {
				t.Fatalf("DownloadHandleWithRange failed: %v", err)
			}

			// 检查状态码
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// 检查响应数据
			actualData := w.Body.String()
			if actualData != tt.expectedData {
				t.Errorf("Expected data %q, got %q", tt.expectedData, actualData)
			}

			// 检查Accept-Ranges头
			if acceptRanges := w.Header().Get("Accept-Ranges"); acceptRanges != "bytes" {
				t.Errorf("Expected Accept-Ranges header to be 'bytes', got %q", acceptRanges)
			}

			// 对于范围请求，检查Content-Range头
			if tt.rangeHeader != "" {
				contentRange := w.Header().Get("Content-Range")
				if contentRange == "" {
					t.Error("Expected Content-Range header for range request")
				}
			}
		})
	}
}
