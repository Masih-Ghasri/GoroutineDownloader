package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// تابع دانلود فایل
func downloadFile(url string, filepath string, wg *sync.WaitGroup, progress chan<- string) {
	defer wg.Done() // وقتی کار تموم شد، به WaitGroup اطلاع بده

	// ایجاد فایل جدید
	output, err := os.Create(filepath)
	if err != nil {
		progress <- fmt.Sprintf("خطا در ایجاد فایل %s: %v", filepath, err)
		return
	}
	defer output.Close()

	// دریافت فایل از اینترنت
	response, err := http.Get(url)
	if err != nil {
		progress <- fmt.Sprintf("خطا در دانلود %s: %v", url, err)
		return
	}
	defer response.Body.Close()

	// کپی کردن داده‌ها به فایل
	_, err = io.Copy(output, response.Body)
	if err != nil {
		progress <- fmt.Sprintf("خطا در ذخیره فایل %s: %v", filepath, err)
		return
	}

	progress <- fmt.Sprintf("دانلود %s با موفقیت انجام شد!", filepath)
}

func main() {
	startTime := time.Now() // زمان شروع اجرا

	// لیست فایل‌ها برای دانلود
	fileUrls := map[string]string{
		"file1.jpg": "https://example.com/file1.jpg",
		"file2.jpg": "https://example.com/file2.jpg",
		"file3.jpg": "https://example.com/file3.jpg",
	}

	var wg sync.WaitGroup          // WaitGroup برای هماهنگی Goroutine‌ها
	progress := make(chan string) // Channel برای گزارش پیشرفت

	// شروع دانلود هر فایل در یک Goroutine جداگانه
	for filepath, url := range fileUrls {
		wg.Add(1)
		go downloadFile(url, filepath, &wg, progress)
	}

	// Goroutine برای چاپ پیشرفت دانلود
	go func() {
		for msg := range progress {
			fmt.Println(msg)
		}
	}()

	wg.Wait() // منتظر ماندن تا تمام Goroutine‌ها کارشون تموم شه
	close(progress) // بستن Channel

	fmt.Printf("تمام فایل‌ها با موفقیت دانلود شدند! زمان اجرا: %v \n", time.Since(startTime))
}