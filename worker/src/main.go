package main

import (
	"github.com/dlintw/goconf"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"strings"
	"net/url"
)

type Command struct {
	Inst string
	Arg string
	Oth string
}

var _cfile *goconf.ConfigFile  // A file-handler required to hold the config file
var _confMap map[string]string // A map to hold the config file params, we don't read the file again and again
var alert chan map[string]interface{}
var data chan map[string]interface{}
var _store chan map[string]interface {}

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
	fo, err := os.Create(pwd + string(filepath.Separator) + "worker.log") // create log file in current dir
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.Writer(fo))
	// 		Done setting up Logging
	if func() bool {
		_confMap = make(map[string]string)
		pwd, err := os.Getwd() // get present working directory
		if err != nil {
			log.Fatal(err)
		}
		path := pwd + string(filepath.Separator) + "cockpit" + string(filepath.Separator) + "workerconf"
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

	poller := time.NewTicker(time.Duration(int64(5) * int64(time.Second))) // Actual time variable that is a counter]

	// Make a listener.
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"status\":\"alive\"}"))
		})
		router.HandleFunc("/detail" ,func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("{\"status\":\"detail\"}"))
		})
		router.HandleFunc("/status" ,func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("{\"status\":\"status\"}"))
		})
		http.Handle("/", router)
		http.ListenAndServe(":63001", nil) // TODO add this to conf
	}()
	for {
		select {
		case item := <-alert:
			// This handles alert, and is recieved at push
			// Dial port 6300
			log.Print(item)
		case item := <-data:
			log.Print(item)
		case <-poller.C:
			// Poll for values at certain period of time
			go func(){
				// Checks the commands and raises alerts
				// Read the command.json file
				file, e := ioutil.ReadFile("./cockpit/command.json")
				if e != nil {
					log.Fatal("File Error")
				}
				var cjson interface {}
				err := json.Unmarshal(file,&cjson)
				if err != nil{
					log.Fatal("Error reading command.json")
				}
				for _, val := range cjson.(map[string]interface {})["commands"].(map[string]interface{}){
					command :=  val.(map[string]interface {})["command"].(string)
					argument := strings.Fields(val.(map[string]interface {})["argument"].(string))
					_, err := exec.Command(command ,argument...).Output()
					if err != nil {
						log.Print("Unable to run program")
						// This means the execution failed.
						//We raise the alert request here.
						_ ,err := http.PostForm("http://localhost:63000/alert",
							url.Values{"alert": {"TestServer1"}, "error": {"Error Running Program"}})
						if err != nil {
							// og
							log.Print("Error sending data to server")
						}
					}
				}
			}()
		}
	}

}
