package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func httpError(w http.ResponseWriter, err error) {
	log.Printf("Error: %v", err)
	w.WriteHeader(500)
	if newErr := errorPage.ExecuteTemplate(w, "base", nil); newErr != nil {
		log.Printf("Template Error: %v", newErr)
	}
}

func upload(dir string, maxBytes int64, parseBytes int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limit size of request
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

		// parse with upto `parseBytes` in memory, otherwise use disk
		err := r.ParseMultipartForm(parseBytes)
		if err != nil {
			httpError(w, err)
			return
		}

		file, header, err := r.FormFile("uploaded_file")
		if err != nil {
			httpError(w, err)
			return
		}
		defer file.Close()

		log.Printf("Uploaded File: %v %v %v", header.Filename, header.Size, header.Header)

		newFile, err := os.OpenFile(filepath.Join(dir, header.Filename), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			httpError(w, err)
			return
		}
		defer newFile.Close()

		io.Copy(newFile, file)
		err = successPage.ExecuteTemplate(w, "base", nil)
		if err != nil {
			httpError(w, err)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if err := indexPage.ExecuteTemplate(w, "base", nil); err != nil {
		httpError(w, err)
	}
}

func main() {
	var (
		addr        = flag.String("address", ":8080", "listen address")
		uploadDir   = flag.String("dir", "/tmp", "directory to upload to")
		timeoutSecs = flag.Int("timeout", 15, "number of seconds to wait before timing out")
		maxBytes    = flag.Int("max-bytes", 10<<20, "maximum request size in bytes, default is 10MB")
		parseBytes  = flag.Int("parse-bytes", 10<<20, "use upto this number of bytes to parse, otherwise use disk for overflow, default is 10MB")
	)
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/upload/", upload(*uploadDir, int64(*maxBytes), int64(*parseBytes)))
	mux.HandleFunc("/", index)

	timeout := time.Duration(*timeoutSecs) * time.Second
	s := http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}

	log.Printf("Listening on: %s", *addr)
	log.Fatal(s.ListenAndServe())
}
