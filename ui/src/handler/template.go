package handler

import (
	"html/template"
	"log"
	"path/filepath"
)

func loadTempl(typ string) *template.Template {
	if _, ok := _templateMap[typ]; !ok {
		path := "template"+ string(filepath.Separator)
		templ, err := template.ParseFiles(
			path+"index.html",
			path+"_header.html",
			path+"_body.html",
			path+"_footer.html",
			path+"_extra.html",
			path+typ+string(filepath.Separator)+"ctrl_header.html",
			path+typ+string(filepath.Separator)+"ctrl_body.html",
			path+typ+string(filepath.Separator)+"ctrl_footer.html",
			path+typ+string(filepath.Separator)+"ctrl_extra.html",
		)
		_templateMap[typ] = templ
		if err != nil {
			log.Panic("Error loading template " + err.Error())
			return nil
		}
	}
	return _templateMap[typ]
}


func contextMap() map[string]interface {}{
	mp := map[string]interface{}{
		"client": "Cockpit",
		"logo":   "http://beta.dashboard.unbxd.com/users/images/logo-unbxd-black.png",
	}
	return (mp)
}

