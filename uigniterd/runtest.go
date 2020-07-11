package main

import (
	"log"
	"time"
)

func runtest() {
	opt := &Options{
		1,
		128,
		DefaultKernel,
		ImageRoot + "hello.raw",
		true,
		"--bootchart hello",
	}

	for i := 0; i < 5; i++ {
		err := runVM(opt)
		if err != nil {
			log.Print(err)
		}
		time.Sleep(time.Second)
	}
}
