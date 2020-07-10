package main

import (
	"log"
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
	err := runVM(opt)
	if err != nil {
		log.Fatalln(err)
	}
	// for i := 0; i < 5; i++ {
	// 	runVM(opt)
	// 	time.Sleep(time.Second)
	// }
	// time.Sleep(time.Minute)
}
