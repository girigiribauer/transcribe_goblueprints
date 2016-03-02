package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	//"github.com/girigiribauer/transcribe_goblueprints/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

// templ は1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈します
	gomniauth.SetSecurityKey("__my_security_key_is_here__")
	gomniauth.WithProviders(
		facebook.New(
			os.Getenv("CHATAPP_FACEBOOK_CLIENTID"),
			os.Getenv("CHATAPP_FACEBOOK_SECRETKEY"),
			"http://localhost:8080/auth/callback/facebook",
		),
		github.New(
			os.Getenv("CHATAPP_GITHUB_CLIENTID"),
			os.Getenv("CHATAPP_GITHUB_SECRETKEY"),
			"http://localhost:8080/auth/callback/github",
		),
		google.New(
			os.Getenv("CHATAPP_GOOGLE_CLIENTID"),
			os.Getenv("CHATAPP_GOOGLE_SECRETKEY"),
			"http://localhost:8080/auth/callback/google",
		),
	)
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()
	// Webサーバを開始します
	log.Println("Webサーバを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
