# Uigniter

Light weight orchestration tools for OSv unikernels, backed by Firecracker.

- Uigniterd: the uigniter server which boots and manages OSv/firecracker instances, controlled via a RESTful API
- Uigniterctl: the command line tool to communicate with uigniterd (not developed yet)

## Prerequisites

- A Linux machine with KVM enabled
- [Firecracker](https://firecracker-microvm.github.io/) installed (in env PATH)

## Getting started

### Build from source

```shell
git clone https://github.com/richardlee159/uigniter.git
cd uigniter/uigniterd
go build
```

### Run  uigniter server

```shell
sudo ./uigniterd
```

Now the server is listening on 127.0.0.1:6666.

Note that Uigniter loads OSv kernels and disk images from a root repository folder: `/var/lib/uigniter` . For now, it's created automatically but managed manually, so you need to copy your OSv kernels into `/var/lib/uigniter/kernel` and images into `/var/lib/uigniter/image` .

### Usage

- create new instance

  POST `http://127.0.0.1:6666/vm/create`

  The request body is json format. For example: 

  ```json
  {
      "image_name": "hello",
      "cmdline": "hello",
      "read_only":true
  }
  ```

  will use the image `/var/lib/uigniter/image/hello.raw` and the kernel `/var/lib/uigniter/kernel/kernel.elf` .

  Get status code 201 if success, and the response body will tell you the id and ipv4 address of the running instance (in json format, of course) .

- stop an instance

  POST `http://127.0.0.1:6666/vm/{id}/stop`

  Get status code 201 if success.



