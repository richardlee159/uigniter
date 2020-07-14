package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/satori/uuid"
)

type Options struct {
	VcpuCount   int    `json:"-"`
	MemSize     int    `json:"-"`
	KernelPath  string `json:"-"`
	DiskPath    string `json:"image_name"`
	ReadOnly    bool   `json:"read_only"`
	CommandLine string `json:"cmdline"`
}

func runAPIServer() {
	router := mux.NewRouter()
	vmRouter := router.PathPrefix("/vm").Subrouter()

	vmRouter.HandleFunc("/run", RunVMHandler).Methods("POST")
	vmRouter.HandleFunc("/{id}/stop", StopVMHandler).Methods("POST")

	log.Println("API server listening on port 6666")
	log.Fatal(http.ListenAndServe(":6666", router))
}

func RunVMHandler(w http.ResponseWriter, r *http.Request) {

	opt := &Options{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, opt)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	opt.VcpuCount = 1
	opt.MemSize = 128
	opt.KernelPath = DefaultKernel

	if opt.ReadOnly {
		opt.DiskPath = ImageRoot + opt.DiskPath + ".raw"
	} else {
		imageDiskPath := ImageRoot + opt.DiskPath + ".raw"
		vmDiskPath := VMRoot + opt.DiskPath + uuid.NewV4().String()[:8] + ".raw"
		err := copy(imageDiskPath, vmDiskPath)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		opt.DiskPath = vmDiskPath
	}

	vm, err := RunVM(opt)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	b, _ := json.Marshal(map[string]string{
		"id":         vm.uuid,
		"ip_address": vm.ipAddr,
	})
	w.Write(b)
}

func StopVMHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := StopVM(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func copy(src, dst string) error {
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
