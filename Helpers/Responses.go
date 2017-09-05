package Helpers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

type Page struct {
	Title string
}

func loadPage(title string) (*Page, error) {
	return &Page{Title: title}, nil
}

func ConvertToJSON(status string, user_data map[string]interface{}) string {
	data := map[string]interface{}{
		"status": status,
		"data":   user_data,
	}
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("bad request")
	}
	return string(json_data)
}

/*
	to render html pages
*/
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string) {
	t, err := template.ParseFiles("static/templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.Header().Set("X-CSRF-Token", csrf.Token(r))
	//fmt.Println(csrf.Token(r))
	err = t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
	to serve static css and js files
*/
func ServeStaticCSS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, _ := loadPage("css file")
	filename := ps.ByName("filename")
	t, _ := template.ParseFiles("static/css/" + filename)
	w.Header().Set("Cache-Control", "no-cache, must-revalidate") //HTTP 1.1
	w.Header().Set("Pragma", "no-cache")                         //HTTP 1.0
	w.Header().Set("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
	w.Header().Set("Content-Type", "text/css")
	t.Execute(w, p)
}

func ServeStaticJS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, _ := loadPage("js file")
	filename := ps.ByName("filename")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate") //HTTP 1.1
	w.Header().Set("Pragma", "no-cache")                         //HTTP 1.0
	w.Header().Set("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
	w.Header().Set("Content-Type", "text/javascript")
	t, _ := template.ParseFiles("static/js/" + filename)
	t.Execute(w, p)
}
