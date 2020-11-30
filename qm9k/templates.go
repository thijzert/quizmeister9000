package qm9k

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

func init() {
	var mustLoad = []string{
		"headers",
		"chat",
		"login",
		"home",
	}

	for _, n := range mustLoad {
		_, err := getTemplate(n)
		if err != nil {
			log.Fatalf("error parsing template '%s': %s", n, err)
		}
	}
}

var parsedTemplates map[string]*template.Template

func getTemplate(name string) (*template.Template, error) {
	if parsedTemplates == nil {
		parsedTemplates = make(map[string]*template.Template)
	}

	if assetsEmbedded {
		if tp, ok := parsedTemplates[name]; ok {
			return tp, nil
		}
	}

	var tp *template.Template

	if name == "headers" {
		b, err := getAsset(path.Join("templates", name+".html"))
		if err != nil {
			return nil, err
		}

		funcs := template.FuncMap{
			"ErrorHeadline": UserHeadline,
			"ErrorMessage":  UserMessage,
		}

		tp, err = template.New("headers").Funcs(funcs).Parse(string(b))
		if err != nil {
			return nil, err
		}
	} else {
		headers, err := getTemplate("headers")
		if err != nil {
			return nil, err
		}

		tp, err = headers.Clone()
		if err != nil {
			return nil, err
		}

		b, err := getAsset(path.Join("templates", name+".html"))
		if err != nil {
			return nil, err
		}

		_, err = tp.Parse(string(b))
		if err != nil {
			return nil, err
		}
	}

	parsedTemplates[name] = tp
	return tp, nil
}

// appRoot finds the relative path to the application root
func (s *Server) appRoot(r *http.Request) string {
	// Find the relative path for the application root by counting the number of slashes in the relative URL
	c := strings.Count(r.URL.Path, "/") - 1
	if c == 0 {
		return "./"
	}
	return strings.Repeat("../", c)
}

func (s *Server) executeTemplate(name string, pageData interface{}, w http.ResponseWriter, r *http.Request) {
	tpl, err := getTemplate(name)
	if err != nil {
		s.errorHandler(err, w, r)
	}

	w.Header()["Content-Type"] = []string{"text/html; charset=UTF-8"}

	tpData := struct {
		AppRoot       string
		AssetLocation string
		PageCSS       string
		Session       *Session
		PageData      interface{}
	}{
		AppRoot:       s.appRoot(r),
		AssetLocation: s.appRoot(r) + "assets",
		Session:       s.MaybeSession(r),
		PageData:      pageData,
	}

	if _, err := getAsset(path.Join("dist", "css", "pages", name+".css")); err == nil {
		tpData.PageCSS = name
	}

	var b bytes.Buffer
	err = tpl.ExecuteTemplate(&b, "headers", tpData)
	if err != nil {
		s.errorHandler(err, w, r)
		return
	}
	io.Copy(w, &b)
}
