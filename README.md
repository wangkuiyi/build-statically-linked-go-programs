# How to Build Statically Linked Go Programs


## Conclusion

We can build our Go programs, as long as they don't depend on `cgo`,
into a fully statically linked executable binary file.

This is simplifies the development of servers and distributed systems,
since we don't need to deploy dependent shared libraries onto the
cluster.


## Background

It has been very confusing to me that how can I make sure that I can
build my Go program into a fully statically linked executable binary
file.

[This post](http://blog.hashbangbash.com/2014/04/linking-golang-statically/) explains
that if my program uses `cgo`, then it cannot be built into fully
statically linked executable binary file.

[This post](http://matthewkwilliams.com/index.php/2014/09/28/go-executables-are-statically-linked-except-when-they-are-not/) claims
that the standard Go package "`net/http` do not have a pure go
implementation at this point".  This matches with my recent experience
that if my program depends on and depends only on `net/http`, `go
build` generates a executable file depending on `libc.so`.  However,
is it possible that I build my program that depends on `net/http` into
a fully statically linked file?

It seems
from
[this explanation](https://github.com/golang/go/issues/9344#issuecomment-149442382) that
I'd have to build the Go compiler with `CGO_ENABLED=0` specified,
before I can use the compiler to build my program into a statically
linked file.  However, the following experiments show that I don't
have to rebuilt the Go compiler, but only need to use `CGO_ENABLED=0`
with `go build my-program.go`.


## Experiments

I
created
[this `Vagrantfile`](https://github.com/wangkuiyi/build-statically-linked-go-programs/blob/master/Vagrantfile),
so that I can create three VMs conveniently by:

```
git clone https://github.com/wangkuiyi/build-statically-linked-go-programs
cd build-statically-linked-go-programs
vagrant up
```

These 3 VMs are:

1. `release`, which has pre-built Go 1.7.1 downloaded from the official site and installed into `/usr/loca/go/bin`,
1. `built`, which has a Go 1.7.1 built from source code with the command `CGO_ENABLED=0 ./all.bash`, and
1. `docker`, which doesn't have Go compiler but with Docker installed, so I can run `docker golang`.


On these 3 VMs, we build the same
simple
[Go program](https://github.com/wangkuiyi/build-statically-linked-go-programs/blob/master/hello/hello.go) that
depends on `net/http` with and without `CGO_ENABLED=0`.

```
package main

import (
    "fmt"
    "net/http"
    "log"
)


func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	    fmt.Fprintf(w, "Hello!")
    })

    if e := http.ListenAndServe(":9090", nil); e != nil {
        log.Fatal("ListenAndServe: ", e)
    }
}
```

1. on the `release` VM:

   ```
vagrant@vagrant-ubuntu-trusty-64:~$ go build /vagrant/hello/hello.go
vagrant@vagrant-ubuntu-trusty-64:~$ ldd hello 
	linux-vdso.so.1 =>  (0x00007ffdcce8e000)
	libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f3437ad7000)
	libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f3437712000)
	/lib64/ld-linux-x86-64.so.2 (0x00007f3437cf5000)
vagrant@vagrant-ubuntu-trusty-64:~$ CGO_ENABLED=0 go build /vagrant/hello/hello.go
vagrant@vagrant-ubuntu-trusty-64:~$ ldd hello 
	not a dynamic executable
   ```

1. on the `built` VM

   ```
vagrant@vagrant-ubuntu-trusty-64:~$ go build /vagrant/hello/hello.go
vagrant@vagrant-ubuntu-trusty-64:~$ ldd hello 
	linux-vdso.so.1 =>  (0x00007ffee4cdb000)
	libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007fd84e60b000)
	libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd84e246000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fd84e829000)
vagrant@vagrant-ubuntu-trusty-64:~$ CGO_ENABLED=0 go build /vagrant/hello/hello.go
vagrant@vagrant-ubuntu-trusty-64:~$ ldd hello 
	not a dynamic executable
   ```

1. on the `docker` VM

   ```
vagrant@vagrant-ubuntu-trusty-64:~$ sudo docker run -it --rm -v /vagrant:/vagrant golang:alpine /bin/ash
/go # apk add --update ldd
/go # go build /vagrant/hello/hello.go
/go # ldd hello 
	/lib/ld-musl-x86_64.so.1 (0x7f55934aa000)
	libc.musl-x86_64.so.1 => /lib/ld-musl-x86_64.so.1 (0x7f55934aa000)
/go # CGO_ENABLED=0 go build /vagrant/hello/hello.go
/go # ldd hello 
ldd: hello: Not a valid dynamic program
   ```


We can see that in all experiments, as long as we use `CGO_ENABLED=1`,
we get fully statically linked executable file.

<!--  LocalWords:  cgo http libc issuecomment Vagrantfile VMs cd pre
 -->
<!--  LocalWords:  golang fmt func HandleFunc ResponseWriter Fprintf
 -->
<!--  LocalWords:  ListenAndServe VM ldd sudo rm apk
 -->
