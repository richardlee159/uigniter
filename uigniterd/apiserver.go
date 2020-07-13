package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Options struct {
	VcpuCount   int
	MemSize     int
	KernelPath  string
	DiskPath    string
	ReadOnly    bool
	CommandLine string
}

func runAPIServer() {
	router := mux.NewRouter()
	vmRouter := router.PathPrefix("/vm").Subrouter()

	vmRouter.HandleFunc("/create", CreateVMHandler).Methods("POST")
	vmRouter.HandleFunc("/{id}/stop", StopVMHandler).Methods("POST")

	log.Println("API server listening on port 6666")
	log.Fatal(http.ListenAndServe(":6666", router))
}

func CreateVMHandler(w http.ResponseWriter, r *http.Request) {

	var req map[string]string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	image, ok := req["image_name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cmdline, ok := req["cmdline"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	opt := &Options{
		1,
		128,
		DefaultKernel,
		ImageRoot + image + ".raw",
		true,
		cmdline,
	}

	vm, err := RunVM(opt)
	if err != nil {
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
