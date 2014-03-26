/*
	Cockpit - UI

	Description:
		The purpose of this package is to be the UI back-end of the system.
		This is an interface between Cockpit - HUB and Web-Browser/PhoneGap app.
		This package's primary purpose is to serve static contents and some business logic
		It is also responsible for maintaining session/user details
*/

package main

import (
	"github.com/dlintw/goconf"
	"github.com/gorilla/mux"
	"handler"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var _cfile *goconf.ConfigFile  // A file-handler required to hold the config file
var _confMap map[string]string // A map to hold the config file params, we don't read the file again and again

func main() {
	/* -------------------------------------  */
	/* 			Boilerplate					  */
	/* -------------------------------------  */
	/*
		Setting up Logging, this is standard practice for all go programs
	*/
	pwd, err := os.Getwd() // get present working directory
	if err != nil {
		log.Fatal(err)
	}
	fo, err := os.Create(pwd + string(filepath.Separator) + "server.log") // create log file in current dir
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.Writer(fo))
	// 		Done setting up Logging
	// Initialize the config variables
	// This is triggered only once and loads up the conf file.
	if func() bool {
		_confMap = make(map[string]string)
		pwd, err := os.Getwd() // get present working directory
		if err != nil {
			log.Fatal(err)
		}
		path := pwd + string(filepath.Separator) + "ui.conf"
		_, err = os.Stat(path)
		if err != nil {
			log.Fatal("Config file not found.")
		}
		_cfile, err = goconf.ReadConfigFile(path)
		if err != nil {
			log.Fatal("Error reading config file." + err.Error())
		}
		return true
	}() {
		log.Print("Loaded Logging, starting up networking interface ...")
		handler.Init()
		go func() {
			router := mux.NewRouter()
			/*
				Routes
			*/
			router.HandleFunc("/login", handler.Login).Methods("GET")
			/* Static content*/
			router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
			// Lets start the listener
			http.Handle("/", router)
			http.ListenAndServe(":63000", nil)
		}()
		select {}
	} else {
		log.Fatal("Error loading conf file ...")
	}

}
