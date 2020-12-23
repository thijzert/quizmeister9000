package qm9k

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
	"github.com/thijzert/quizmeister9000/qm9k/weberrors"
)

type htmlHandler struct {
	Server         *Server
	TemplateName   string
	RequestDecoder handlers.RequestDecoder
	Handler        handlers.RequestHandler
}

// HTMLFunc creates a HTTP handler that outputs HTML
func (s *Server) HTMLFunc(handler handlers.RequestHandler, decoder handlers.RequestDecoder, templateName string) http.Handler {
	return htmlHandler{
		Server:         s,
		TemplateName:   templateName,
		RequestDecoder: decoder,
		Handler:        handler,
	}
}

var csp string

func init() {
}

func (h htmlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.RequestDecoder(r)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	tpl, err := h.Server.getTemplate(h.TemplateName)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	state := h.Server.getState(r)
	newState, resp, err := h.Handler(state, req)
	if err != nil {
		if _, ok := err.(weberrors.Redirector); ok {
			// HACK: save the state for redirect events
			saveErr := h.Server.setState(newState)
			if saveErr != nil {
				h.Error(w, r, saveErr)
				return
			}
		}
		h.Error(w, r, err)
		return
	}

	err = h.Server.setState(newState)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	w.Header()["Content-Type"] = []string{"text/html; charset=UTF-8"}

	csp := ""
	csp += "default-src 'self' blob: data: ; "
	csp += "script-src 'self' blob: ; "
	csp += "style-src 'self' data: 'unsafe-inline'; "
	csp += "img-src 'self' blob: data: ; "
	csp += "connect-src 'self' blob: data: ; "
	csp += "frame-src 'none' ; "
	csp += "frame-ancestors 'none'; "
	csp += "form-action 'self'; "
	w.Header()["Content-Security-Policy"] = []string{csp}
	w.Header()["X-Frame-Options"] = []string{"deny"}
	w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	w.Header()["Referrer-Policy"] = []string{"strict-origin-when-cross-origin"}
	w.Header()["X-Content-Type-Options"] = []string{"nosniff"}

	tpData := struct {
		AppRoot       string
		AssetLocation string
		PageCSS       string
		Session       *Session
		CurrentUser   handlers.User
		Request       handlers.Request
		Response      handlers.Response
	}{
		AppRoot:       h.appRoot(r),
		AssetLocation: h.appRoot(r) + "assets",
		Session:       h.Server.MaybeSession(r),
		Request:       req,
		Response:      resp,
	}

	cssName := h.TemplateName
	if len(cssName) > 5 && cssName[0:5] == "full/" {
		cssName = cssName[5:]
	}

	if _, err := getAsset(path.Join("dist", "css", "pages", cssName+".css")); err == nil {
		tpData.PageCSS = cssName
	}

	if !tpData.Session.UserID.Empty() {
		u, ok := h.Server.getUser(tpData.Session.UserID)
		if ok {
			tpData.CurrentUser = u
		}
	}

	var b bytes.Buffer
	err = tpl.ExecuteTemplate(&b, "headers", tpData)
	if err != nil {
		h.Error(w, r, err)
		return
	}
	io.Copy(w, &b)
}

func (s *Server) getTemplate(name string) (*template.Template, error) {
	if s.parsedTemplates == nil {
		s.parsedTemplates = make(map[string]*template.Template)
	}

	if assetsEmbedded {
		if tp, ok := s.parsedTemplates[name]; ok {
			return tp, nil
		}
	}

	var tp *template.Template

	b, err := getAsset(path.Join("templates", name+".html"))
	if err != nil {
		return nil, err
	}

	funcs := template.FuncMap{}

	if name == "full/headers" {
		tp, err = template.New("headers").Funcs(funcs).Parse(string(b))
		if err != nil {
			return nil, err
		}
	} else if len(name) > 5 && name[:5] == "full/" {
		headers, err := s.getTemplate("full/headers")
		if err != nil {
			return nil, err
		}

		tp, err = headers.Clone()
		if err != nil {
			return nil, err
		}

		_, err = tp.Parse(string(b))
		if err != nil {
			return nil, err
		}
	} else {
		tp, err = template.New("headers").Funcs(funcs).Parse(string(b))
		if err != nil {
			return nil, err
		}
	}

	s.parsedTemplates[name] = tp
	return tp, nil
}

// appRoot finds the relative path to the application root
func (h htmlHandler) appRoot(r *http.Request) string {
	return h.Server.appRoot(r)
}

func (h htmlHandler) Error(w http.ResponseWriter, r *http.Request, err error) {
	if redir, ok := err.(weberrors.Redirector); ok {
		st := 302
		if statuser, okk := redir.(weberrors.HTTPError); okk {
			st = statuser.HTTPStatus()
		}
		h.Server.redirect(w, r, h.appRoot(r)+redir.RedirectLocation(), st)
		return
	}

	// TODO: we may need to set a different status entirely
	w.WriteHeader(500)
	fmt.Fprintf(w, "Error: %s", err)
}
