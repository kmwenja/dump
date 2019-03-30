package main

import "html/template"

var base = template.Must(template.New("base").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>{{ block "title" . -}}{{- end }} - Dump</title>
	<style>
	    body {
		    max-width: 960px;
			margin: 0 auto;
			padding: 5px 10px;
		}
	</style>
  </head>
  <body>
	{{ block "content" . }}{{ end }}
  </body>
</html>
`))

var indexPage = template.Must(template.Must(base.Clone()).New("index").Parse(`
{{ define "title" }}Upload File{{end}}
{{ define "content" }}
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
  <input type="submit" value="Upload" />
</form>
{{end}}
`))

var errorPage = template.Must(template.Must(base.Clone()).New("error").Parse(`
{{ define "title" }} Error {{ end }}
{{ define "content" }}
<p>Internal error. <a href="/">Upload again.</a></p>
{{ end }}
`))

var successPage = template.Must(template.Must(base.Clone()).New("success").Parse(`
{{ define "title" }} Success {{ end }}
{{ define "content" }}
<p>Upload done. <a href="/">Upload again.</a></p>
{{ end }}
`))
