package handler

/*
	All handlers are defined here
 */

import(
	"net/http"
	"html/template"
	"log"
)

var _templateMap map[string] *template.Template

func Init(){
	_templateMap = make(map[string]*template.Template)
}

func Login(w http.ResponseWriter,r *http.Request){
	templ := loadTempl("Login")
	mp := contextMap()
	err := templ.ExecuteTemplate(w, "html", mp)
	if err != nil{
		log.Panic("Error executing template." + err.Error())
	}
}

