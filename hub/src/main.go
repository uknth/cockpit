package main

import (
	"github.com/dlintw/goconf"
	"os"
	"log"
	"io"
	"path/filepath"
	"github.com/gorilla/mux"
	"HTTPHandler"
	"net/http"

)

var _cfile *goconf.ConfigFile  // A file-handler required to hold the config file
var _confMap map[string]string // A map to hold the config file params, we don't read the file again and again
var _store  map[string]interface{}


func main() {
	// Begin here.
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
		path := pwd + string(filepath.Separator) + "cockpit" + string(filepath.Separator) + "conf"
		_, err = os.Stat(path)
		if err != nil {
			log.Fatal("Config file not found. Config file is located in same directory as executable in a folder .arwen/config")
		}
		_cfile, err = goconf.ReadConfigFile(path)
		if err != nil {
			log.Fatal("Error reading config file." + err.Error())
		}
		return true
	}() {
		log.Print("Loaded Logging, starting up server ...")
	} else {
		log.Fatal("Error loading conf file ...")
	}
	/* -------------------------------------  */
	/*			Coding Begins				  */
	/* -------------------------------------  */
	//Variable stores everything

	_store  = make (map[string]interface{})
	ch := make (chan map[string]interface{}) // {action : {object}}
	// Initiate HTTP Listener
	go func(store map[string]interface {}, ch chan map[string]interface{}){
		// Define the http listener
		router := mux.NewRouter() // TODO add this in conf
		// ROUTERS
		router.HandleFunc("/auth", func(w http.ResponseWriter,r *http.Request){
				HTTPHandler.Auth(w,r)
			})
		router.HandleFunc("/add/{action}/",func(w http.ResponseWriter,r *http.Request){
				HTTPHandler.Add(w,r,ch)
			} )
		http.Handle("/", router)
		http.ListenAndServe(":63000", nil) // TODO add this to conf
	}(_store,ch)
	for {
		select {
		case item := <- ch:
			// Poll on Request
			log.Print(item)
//		case <-ticker.C:
			// Poller
		}
	}
}
