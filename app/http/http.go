// Package http provides a convenient way to impliment http servers
package http

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/pottava/docker-webui/app/config"
	"github.com/pottava/docker-webui/app/logs"
	"github.com/pottava/docker-webui/app/misc"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

// RequestGetParam retrives a request parameter
func RequestGetParam(r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	return value, (value != "")
}

// RequestGetParamS retrives a request parameter as string
func RequestGetParamS(r *http.Request, key, def string) string {
	value, found := RequestGetParam(r, key)
	if !found {
		return def
	}
	return value
}

// RequestGetParamI retrives a request parameter as int
func RequestGetParamI(r *http.Request, key string, def int) int {
	value, found := RequestGetParam(r, key)
	if !found {
		return def
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

// SplittedUpperStrings split word to array and change those words to UpperCase
func SplittedUpperStrings(value string) []string {
	splitted := strings.FieldsFunc(value, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	words := make([]string, len(splitted))
	for i, val := range splitted {
		words[i] = strings.ToUpper(val)
	}
	return words
}

// RequestPostParam retrives a POST request parameter
func RequestPostParam(r *http.Request, key string) (string, bool) {
	value := r.PostFormValue(key)
	return value, (value != "")
}

// RequestPostParamS retrives a request parameter as string
func RequestPostParamS(r *http.Request, key, def string) string {
	value, found := RequestPostParam(r, key)
	if !found {
		return def
	}
	return value
}

// Chain enables middleware chaining
func Chain(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return chain(true, false, true, f)
}

// AssetsChain enables middleware chaining
func AssetsChain(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return chain(false, true, false, f)
}

// RenderText write data as a simple text
func RenderText(w http.ResponseWriter, data string, err error) {
	if IsInvalid(w, err, "@RenderText") {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(html.EscapeString(data)))
}

// RenderHTML write data as a HTML text with template
func RenderHTML(w http.ResponseWriter, templatePath []string, data interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relatives := append([]string{"base.tmpl"}, templatePath...)
	templates := make([]string, len(relatives))
	for idx, template := range relatives {
		templates[idx] = path.Join(cfg.StaticFilePath, "views", template)
	}

	tmpl, err := template.ParseFiles(templates...)
	if IsInvalid(w, err, "@RenderHTML") {
		return
	}
	if err := tmpl.Execute(w, struct {
		AppName        string
		StaticFileHost string
		Mode           string
		Data           interface{}
	}{cfg.Name, cfg.StaticFileHost, cfg.Mode, data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error.Printf("ERROR: @RenderHTML %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
}

// RenderJSON write data as a json
func RenderJSON(w http.ResponseWriter, data interface{}, err error) {
	if IsInvalid(w, err, "@RenderJSON") {
		return
	}
	js, err := json.Marshal(data)
	if IsInvalid(w, err, "@RenderJSON") {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(js)
}

// IsInvalid checks if the second argument represents a real error
func IsInvalid(w http.ResponseWriter, err error, caption string) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error.Printf("ERROR: %s %s", caption, err.Error())
		return true
	}
	return false
}

type customResponseWriter struct {
	io.Writer
	http.ResponseWriter
	status int
}

func (r *customResponseWriter) Write(b []byte) (int, error) {
	if r.Header().Get("Content-Type") == "" {
		r.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return r.Writer.Write(b)
}

func (r *customResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}

func chain(log, cors, validate bool, f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(custom(log, cors, validate, f))
}

func custom(log, cors, validate bool, f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// compress settings
		ioWriter := w.(io.Writer)
		for _, val := range misc.ParseCsvLine(r.Header.Get("Accept-Encoding")) {
			if val == "gzip" {
				w.Header().Set("Content-Encoding", "gzip")
				g := gzip.NewWriter(w)
				defer g.Close()
				ioWriter = g
				break
			}
			if val == "deflate" {
				w.Header().Set("Content-Encoding", "deflate")
				z := zlib.NewWriter(w)
				defer z.Close()
				ioWriter = z
				break
			}
		}
		writer := &customResponseWriter{Writer: ioWriter, ResponseWriter: w, status: 200}

		// route to the controllers
		f(writer, r)
	}
}

func header(r *http.Request, key string) (string, bool) {
	if r.Header == nil {
		return "", false
	}
	if candidate := r.Header[key]; !misc.ZeroOrNil(candidate) {
		return candidate[0], true
	}
	return "", false
}
