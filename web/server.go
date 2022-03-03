package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/hectorcorrea/texto/textdb"
)

var blogRouter Router
var db textdb.TextDb

func init() {
	blogRouter.Add("POST", "/blog/new", blogNew)
	blogRouter.Add("GET", "/blog", blogViewAll)
	blogRouter.Add("GET", "/blog/:id", blogOne)
}

func blogViewAll(s session, values map[string]string) {
	vm := db.ListAll()
	renderTemplate(s, "views/blogList.html", vm)
}

func blogOne(s session, values map[string]string) {
	vm := ""
	renderTemplate(s, "views/blogOne.html", vm)
}

func blogNew(s session, values map[string]string) {
	err := db.CreateNewEntry()
	if err != nil {
		// TODO: render error page
		panic("error creating new entry")
	}
	// TODO: redirect
	vm := db.ListAll()
	renderTemplate(s, "views/blogList.html", vm)
}

func homePage(resp http.ResponseWriter, req *http.Request) {
	vm := ""
	session := newSession(resp, req)
	renderTemplate(session, "views/home.html", vm)
}

func blogPages(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)
	found, route := blogRouter.FindRoute(req.Method, req.URL.Path)
	if found {
		values := route.UrlValues(req.URL.Path)
		route.handler(session, values)
	} else {
		log.Printf("not found")
	}
}

func loadTemplate(s session, viewName string) (*template.Template, error) {
	t, err := template.New("layout").ParseFiles("views/layout.html", viewName)
	if err != nil {
		log.Printf("Error loading template %s (%s)", viewName, s.req.URL.Path)
		return nil, err
	} else {
		log.Printf("Loaded template %s (%s)", viewName, s.req.URL.Path)
		return t, nil
	}
}

func renderTemplate(s session, viewName string, viewModel interface{}) {
	t, err := loadTemplate(s, viewName)
	if err != nil {
		log.Printf("Error loading: %s, %s ", viewName, err)
	} else {
		err = t.Execute(s.resp, viewModel)
		if err != nil {
			log.Printf("Error rendering: %s, %s ", viewName, err)
		}
	}
}

func StartWebServer(settings Settings) {
	log.Printf("Listening for requests at %s\n", "http://"+settings.Address)
	db = textdb.InitTextDb(settings.DataFolder)
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/favicon.ico", fs)
	http.Handle("/robots.txt", fs)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/blog/", blogPages)
	http.HandleFunc("/", homePage)

	err := http.ListenAndServe(settings.Address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}
