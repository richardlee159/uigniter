package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Options struct {
	VcpuCount   int    `json:"-"`
	MemSize     int    `json:"-"`
	KernelPath  string `json:"-"`
	DiskPath    string `json:"image_name"`
	ReadOnly    bool   `json:"read_only"`
	CommandLine string `json:"cmdline"`
}

var ch = make(chan int)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
	ch <- 0
}

func main() {
	http.HandleFunc("/", sayHello)
	go http.ListenAndServe(":8888", nil)

	opt := &Options{
		DiskPath:    "test-client",
		ReadOnly:    true,
		CommandLine: "",
	}
	json, _ := json.Marshal(opt)

	for i := 0; i < 100; i++ {
		body := bytes.NewBuffer(json)
		start := time.Now().UnixNano()
		fmt.Print(http.Post("http://127.0.0.1:6666/vm/create", "application/json", body))
		<-ch
		end := time.Now().UnixNano()
		fmt.Print((end - start) / 1000000)
	}
}
