package site

import (
	"html/template"
	"log"
	"net/http"
)

func ShowCommentsPage(w http.ResponseWriter, comments []Comment) {
	w.WriteHeader(http.StatusOK)
	terr := commentsTemplate.Execute(w, map[string]interface{}{
		"Comments": comments,
	})
	if terr != nil {
		log.Println(terr)
	}
}

func ShowErrorPage(w http.ResponseWriter, statuscode int, title string, err error) {
	w.WriteHeader(statuscode)
	terr := errorTemplate.Execute(w, map[string]interface{}{
		"Title": title,
		"Error": err,
	})
	if terr != nil {
		log.Println(terr)
	}
}

var commentsTemplate = template.Must(template.New(``).Parse(`
<!DOCTYPE html>
<html><body>
	<h1>Comment Site</h1>
	<form action="/comment" method="POST">
		<input type="text" name="user"/>
		<input type="text" name="comment"/>
		<input type="submit" />
	</form>
	<h2>Comments</h2>
	<div>
	{{ range .Comments }}
		<div class="comment">
			{{ .User }} - {{ .Text }}
		</div>
	{{ end }}
	</div>
</body></html>
`))

var errorTemplate = template.Must(template.New(``).Parse(`
<!DOCTYPE html>
<html>
<head>
	{{if .Redirect}}<meta http-equiv="refresh" content="0; url={{.Redirect}}">{{end}}
</head>
<body>
	<h1>{{.Title}}</h1>
	{{if .Error}}<p>{{.Error}}</p>{{end}}
</body>
</html>
`))
