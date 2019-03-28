package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, errorPage)
	log.Printf("Error: %v", err)
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
		fmt.Fprintf(w, successPage)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexPage))
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

var indexPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Uploads</title>
	<style>
	    body {
		    max-width: 960px;
		}
	</style>
  </head>
  <body>
    <h1>Upload</h1>
    <form
      enctype="multipart/form-data"
      action="/upload/"
      method="post"
    >
	  <div>
          <input type="file" name="uploaded_file" />
	  </div>
	  <hr/>
      <input type="submit" value="upload" />
    </form>
  </body>
</html>
`

var errorPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Upload Error</title>
	<style>
	    body {
		    max-width: 960px;
		}
	</style>
  </head>
  <body>
	<p>Internal error. <a href="/">Upload again.</a></p>
  </body>
</html>
`
var successPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Upload Error</title>
	<style>
	    body {
		    max-width: 960px;
		}
	</style>
  </head>
  <body>
	<p>Upload done. <a href="/">Upload again.</a></p>
  </body>
</html>
`
