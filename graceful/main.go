package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	server   *http.Server
	listener net.Listener
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world233333"))
}

func main() {
	go signalHandler()

	graceful := flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
	flag.Parse()

	http.HandleFunc("/hello", handler)
	server = &http.Server{Addr: ":9999"}

	log.Printf("args: %v, graceful: %v", os.Args, *graceful)
	var err error
	if *graceful {
		log.Print("main: Listening to existing file descriptor 3.")
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
	} else {
		log.Print("main: Listening on a new file descriptor.")
		listener, err = net.Listen("tcp", server.Addr)
	}

	if err != nil {
		log.Fatalf("listener error: %v", err)
	}

	err = server.Serve(listener)
	log.Printf("server.Serve err: %v\n", err)
}

func reload() error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		log.Printf("listener is not tcp listener")
		return errors.New("listener is not tcp listener")
	}
	f, err := tl.File()
	if err != nil {
		log.Printf("get listener file error: %v", err)
		return err
	}

	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}

func signalHandler() {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-ch
		log.Printf("signal: %v", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// stop
			signal.Stop(ch)
			server.Shutdown(context.Background())
		case syscall.SIGUSR2:
			// reload
			err := reload()
			if err != nil {
				log.Fatalf("graceful restart error: %v", err)
			}
			server.Shutdown(context.Background())
		}
	}
}
