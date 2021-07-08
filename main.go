package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

func main() {

	dirFlag := flag.String("d", "./", "Directory for cache. Default is CWD")
	httpPort := flag.Int("p", 80, "Port to serve on")
	memeoryFS := flag.Bool("m", false, "Use memory storage")

	dir := *dirFlag

	var filesystem webdav.FileSystem
	if *memeoryFS == true {
		filesystem = webdav.Dir(dir)
	} else {
		filesystem = webdav.NewMemFS()
	}

	srv := &webdav.Handler{
		FileSystem: filesystem,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("%s %s %s ERROR: %s\n", r.RemoteAddr, r.Method, r.URL, err)
			} else {
				log.Printf("%s %s %s \n", r.RemoteAddr, r.Method, r.URL)
			}
		},
	}

	http.Handle("/", srv)

	fmt.Printf("Start server on %d ...\n", *httpPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
