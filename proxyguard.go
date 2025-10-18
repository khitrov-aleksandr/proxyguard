package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	/*r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "welcome"})
	})

	r.Any("/*", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "welcome"})
	})

	http.ListenAndServe(":7070", r)
	*/

	url, _ := url.Parse("https://mail.ru")
	rProxy := httputil.NewSingleHostReverseProxy(url)

	/*rProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}*/

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pp" {
			w.Write([]byte("Hello from Go HTTP Server PP!"))
		} else {
			rProxy.ServeHTTP(w, r)
		}
	}

	http.HandleFunc("/", handler)

	http.ListenAndServe(":7070", nil)
}
