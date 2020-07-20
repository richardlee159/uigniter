package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	resp, _ := http.Get("http://149.28.211.182:8888")
	fmt.Print(resp)
	time.Sleep(time.Minute * 2)
}
