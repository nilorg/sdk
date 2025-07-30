package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nilorg/sdk/storage"
)

func main() {
	// 创建存储实例
	storageInstance := storage.NewDefaultDownloadStorage("attachment")

	// 设置HTTP处理器
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径获取文件名
		filename := r.URL.Path[len("/download/"):]
		if filename == "" {
			http.Error(w, "文件名不能为空", http.StatusBadRequest)
			return
		}

		// 使用支持范围请求的下载处理器
		ctx := context.Background()
		err := storage.DownloadHandleWithRange(ctx, r, w, storageInstance, filename)
		if err != nil {
			log.Printf("下载文件失败: %v", err)
			http.Error(w, "文件下载失败", http.StatusInternalServerError)
			return
		}
	})

	// 设置一个简单的上传处理器用于测试
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "获取上传文件失败", http.StatusBadRequest)
			return
		}
		defer file.Close()

		ctx := context.Background()
		uploadStorage := storage.NewDefaultUploadStorage()
		fullName, err := uploadStorage.Upload(ctx, file, header.Filename)
		if err != nil {
			log.Printf("上传文件失败: %v", err)
			http.Error(w, "文件上传失败", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "文件上传成功: %s", fullName)
	})

	// 示例：如何在代码中使用FromRangeRequestContext
	demonstrateRangeContext()

	// 启动HTTP服务器
	fmt.Println("服务器启动在 :8080")
	fmt.Println("上传: POST /upload")
	fmt.Println("下载: GET /download/{filename}")
	fmt.Println("范围下载示例: curl -H 'Range: bytes=0-99' http://localhost:8080/download/test.txt")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func demonstrateRangeContext() {
	fmt.Println("=== FromRangeRequestContext 使用示例 ===")

	// 示例1：创建范围请求context
	ctx := context.Background()
	rangeCtx := storage.NewRangeRequestContext(ctx, 0, 99) // 读取前100个字节

	// 从context中获取范围请求
	rangeReq, ok := storage.FromRangeRequestContext(rangeCtx)
	if ok {
		fmt.Printf("范围请求: 开始=%d, 结束=%d\n", rangeReq.Start, rangeReq.End)
		fmt.Printf("是否到文件末尾: %v\n", rangeReq.IsToEnd())
		fmt.Printf("请求长度: %d\n", rangeReq.Length())
		fmt.Printf("请求是否有效: %v\n", rangeReq.IsValid())
	}

	// 示例2：创建到文件末尾的范围请求
	rangeToEndCtx := storage.NewRangeRequestContext(ctx, 100, 0) // 从第101个字节读取到末尾

	rangeReq2, ok := storage.FromRangeRequestContext(rangeToEndCtx)
	if ok {
		fmt.Printf("\n范围请求(到末尾): 开始=%d, 结束=%d\n", rangeReq2.Start, rangeReq2.End)
		fmt.Printf("是否到文件末尾: %v\n", rangeReq2.IsToEnd())
		fmt.Printf("请求长度: %d\n", rangeReq2.Length())
	}

	// 示例3：检查空context
	emptyRangeReq, ok := storage.FromRangeRequestContext(ctx)
	if !ok {
		fmt.Println("\n空context中没有范围请求")
	}
	if emptyRangeReq == nil {
		fmt.Println("空context返回nil范围请求")
	}

	fmt.Println("======================================")
}
