package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}

	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

func worker(wg *sync.WaitGroup, id int, url string) {
	resp, err := http.Get(url)
	if err != nil {
		print(err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("     ", resp.StatusCode, "           ", url)
	}

	defer wg.Done()
}

func main() {
	var wg sync.WaitGroup
	i := 0
	file, err := os.Open("dict.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	log.Print("Scan started")
	println("=============================================")
	println("TIME                  STATUS CODE         URL")
	println("=============================================")

	for scanner.Scan() {
		i++
		wg.Add(1)
		SetConsoleTitle(os.Args[1] + scanner.Text())
		go worker(&wg, i, os.Args[1]+scanner.Text())
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}
