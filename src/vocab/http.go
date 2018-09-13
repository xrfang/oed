package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

var LastModified time.Time

func init() {
	bt, _ := strconv.Atoi(_BUILT_)
	LastModified = time.Unix(int64(bt), 0)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func getCookie(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func setCookie(w http.ResponseWriter, name, value string, age int) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		MaxAge: age,
		Secure: false,
	})
}

func modified(w http.ResponseWriter, r *http.Request) bool {
	ims, err := time.Parse(time.RFC1123, r.Header.Get("If-Modified-Since"))
	if err != nil || LastModified.UnixNano() > ims.UnixNano() {
		return true
	}
	http.Error(w, http.StatusText(304), 304)
	return false
}

func setContentType(w http.ResponseWriter, filePath string) {
	switch strings.ToLower(path.Ext(filePath)) {
	case ".css":
		w.Header().Add("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Add("Content-Type", "application/javascript")
	case ".jpg", ".jpeg":
		w.Header().Add("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Add("Content-Type", "image/png")
	case ".txt", ".log":
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	default:
		w.Header().Add("Content-Type", "application/octet-stream")
	}
}

func sendAsset(w http.ResponseWriter, ident interface{}) {
	var name string
	switch ident.(type) {
	case string:
		name = ident.(string)
	case *http.Request:
		r := ident.(*http.Request)
		if !modified(w, r) {
			return
		}
		name = r.URL.Path
	default:
		panic(fmt.Errorf("invalid ident type %T", ident))
	}
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	data, err := Asset(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setContentType(w, name)
	w.Header().Add("Last-Modified", _BUILT_)
	w.Write(data)
}

func renderTemplate(w http.ResponseWriter, tpl string, args interface{}) {
	var buf bytes.Buffer
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
	}()
	helper := template.FuncMap{
		"ver": func() string {
			return "V" + _G_REVS + "." + _G_HASH
		},
	}
	body, err := Asset("templates/" + tpl)
	assert(err)
	t, err := template.New("body").Funcs(helper).Parse(string(body))
	assert(err)
	shared, err := AssetDir("templates/shared")
	if err == nil {
		for _, name := range shared {
			s, err := Asset("templates/shared/" + name)
			assert(err)
			t, err = t.Parse(string(s))
			assert(err)
		}
	}
	assert(t.Execute(&buf, args))
}
