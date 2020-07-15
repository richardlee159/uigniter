package main

import (
    "fmt"
    "bytes"
    "time"
    "encoding/json"
    "net/http"
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

func sayHello(w http.ResponseWriter,r *http.Request){
    fmt.Fprintf(w,"Hello!")
    ch <- 0
}

func main(){
    http.HandleFunc("/",sayHello)
    go http.ListenAndServe(":8888",nil)

    opt:=&Options{
        DiskPath:"curl",
        ReadOnly:true,
        CommandLine:"curl 149.28.211.182:8888",
    }
    json, _ := json.Marshal(opt)
    body := bytes.NewBuffer(json)

    start := time.Now().UnixNano()
    fmt.Print(http.Post("http://127.0.0.1:6666/vm/run","application/json",body))
    <- ch
    end := time.Now().UnixNano()
    fmt.Print(end-start)
}
