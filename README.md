Dump
====

A simple file upload server.

I built this utility server to solve the problem of extracting files out of
a host where I couldn't SCP or S3 upload from (or copy via external drive).

For instance, imagine physically debugging a remote client's host and needing
to copy the debug logs back to your machine which you left at the office and
for *reasons* you can't use external drives, SCP or S3. With this, you instead
upload the files to a server you control that's running this utility
(with the client's approval of course).

The current version doesn't have any security considerations so if you plan
on serving this on the public internet, defend it with a proxy that offers
authentication e.g. nginx with basic auth and let's encrypt tls certs or
nginx with client cert authentication.

Usage:
------

1. [Setup Go](https://golang.org/doc/install)
2. Clone this repo and cd into the repo directory.
3. Run `go build`. This will make a `dump` binary in the same directory.
   You can also build for other platforms by using `GOOS` and `GOARCH`
   environment variables. [Learn more](https://github.com/golang/go/wiki/WindowsCrossCompiling).
4. Run `./dump -help` to see the list of options provided.
   If you run `./dump` without any arguments, it will use the defaults as listed
   in the options.
5. Assuming you ran it with defaults, browse to `http://localhost:8080` and upload a file.
