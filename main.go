package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"talker/config"
	"talker/session"
	"text/template"
	"time"
)

const cookieName = "session_talker"

var (
	GitTag string = "dev"

	//go:embed assets/*
	assets embed.FS

	//go:embed templates/*
	templates embed.FS

	tmpl *template.Template
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := session.SC.Get(r)
	if !ok {
		sid, sd = session.SC.Create()
	}

	// renew session
	session.SC.Save(w, sid, sd)

	//////////////////////////

	err := tmpl.ExecuteTemplate(w, "main.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

	// http.Redirect(w, r, "/payments", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		index, err := assets.ReadFile("assets/login.html")
		if err != nil {
			log.Fatal(err)
		}
		t, err := template.New("login.html").Parse(string(index))
		if err != nil {
			log.Fatal(err)
		}

		// exec template
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// login logic

	// create session
	sid, sd := session.SC.Create()

	// save session
	session.SC.Save(w, sid, sd)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sid, _, ok := session.SC.Get(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// remove session
	session.SC.Delete(w, sid)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// templates FS
	fstemplates, err := fs.Sub(templates, ".")
	if err != nil {
		log.Fatal(err)
	}

	// print

	tmpl = template.Must(template.ParseFS(fstemplates, "templates/*.gohtml"))

	err = config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	session.Create(cookieName)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			session.SC.RemoveExpired()
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           config.CFG.Listen,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on port %s\n", config.CFG.Listen)
	log.Fatal(s.ListenAndServe())
}
