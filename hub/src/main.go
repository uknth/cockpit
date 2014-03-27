package main

import (
	"HTTPHandler"
	"github.com/dlintw/goconf"
	"github.com/gorilla/mux"
	"io"
	"log"
	"memcached"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"encoding/json"
	"io/ioutil"
)

type Server struct{
	Name string
	IP string
}


var _cfile *goconf.ConfigFile  // A file-handler required to hold the config file
var _confMap map[string]string // A map to hold the config file params, we don't read the file again and again
var _store map[string]interface{}  // Store info about the server
var _servers []Server



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

	// Initialize memecached
	memcached.Init()

	/* -------------------------------------  */
	/*			Coding Begins				  */
	/* -------------------------------------  */
	//Variable stores everything
	_store = make(map[string]interface{})
	httpChannel := make(chan map[string]interface{})// {action : {object}}
	pollCounter := 300
	poller := time.NewTicker(time.Duration(int64(pollCounter) * int64(time.Second))) // Actual time variable that is a counter


	// Initiate HTTP Listener
	// Listener for the incoming HTTP request
	go func(store map[string]interface{}, ch chan map[string]interface{}) {
		// Define the http listener
		router := mux.NewRouter()
		// ROUTERS
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("{\"status\":\"alive\"}"))
			})
		router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
			HTTPHandler.Auth(w, r)
		})
		router.HandleFunc("/add/{action}", func(w http.ResponseWriter, r *http.Request) {
			HTTPHandler.Add(w, r, ch)
		})
		http.Handle("/", router)
		http.ListenAndServe(":63000", nil) // TODO add this to conf
	}(_store, httpChannel)




	for {
		select {
		case item := <-httpChannel:
			// Poll on Request
			for k , _ := range item{
				if k == "server"{
					// Do soemthing
					// Add this info to server map
					ser := Server{item[k].(map[string]interface {})["name"].(string),
						item[k].(map[string]interface {})["ip"].(string)}
					flag := false;
					for k,_ := range _servers {
						if _servers[k].Name == ser.Name{
							flag = true
						}
					}
					if !flag{
						_servers = append(_servers,ser)
					}

					log.Print(_servers)
				}
			}
		case <-poller.C:
			// Poller
			// Now loop over servers to check the status of the server
			for _ , val := range _servers {
				// Make a get request to each of these and add that info to the _store
				resp , err := http.Get(val.IP)
				if err != nil{
					log.Print("Error in fetching value for the IP " + val.IP)
					resp = nil
				}
				if resp != nil {
					var f interface {}
					v, er := ioutil.ReadAll(resp.Body)
					if er != nil{
						log.Print("Error reading value from the server")
					}else{
						er := json.Unmarshal([]byte(v) , &f)
						if er != nil{
							log.Print("Error in Unmarshaling data")
							var temp interface {}
							_store[val.Name] = temp
						}else{
							_store[val.Name] = f
						}
					}
				}
			}

		}
	}
}
