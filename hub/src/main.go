package main

import (
	"HTTPHandler"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dlintw/goconf"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
 	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/alexjlockwood/gcm"
)

type Serv struct {
	Name string
	Key  string
	IP   string
}

var _cfile *goconf.ConfigFile     // A file-handler required to hold the config file
var _confMap map[string]string    // A map to hold the config file params, we don't read the file again and again
var _store map[string]interface{} // Store info about the server
var _servers map[string]interface{}

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
	_store = make(map[string]interface{})
	_servers = make(map[string]interface {})
	// Channels to manage data
	httpChannel := make(chan map[string]interface{})                         // {action : {object}}
	poller := time.NewTicker(time.Duration(int64(5) * int64(time.Second))) // Actual time variable that is a counter
	alertChannel := make(chan map[string]string)
	pollerChannel :=  make(chan map[string]interface {})
	//serverChannel := make(chan []Serv)
	HTTPHandler.Init()

	// Load the contents of the servers.json to the _servers
	file, e := ioutil.ReadFile("./cockpit/servers.json")
	if e != nil {
		log.Fatal("File Error")
	}
	log.Print(file)
	if string(file) != "" {
		err = json.Unmarshal(file,&_servers)
		if err != nil{
			log.Fatal("Error reading servers.json")
		}
	}



	// Initiate HTTP Listener
	// Listener for the incoming HTTP request
	go func(store map[string]interface {},servers map[string]interface {} , ch chan map[string]interface{}) {
		// Define the http listener
		router := mux.NewRouter()
		// ROUTER
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"status\":\"alive\"}"))
		})
		router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
			HTTPHandler.Auth(w, r)
		})
		router.HandleFunc("/add/{action}", func(w http.ResponseWriter, r *http.Request) {
			HTTPHandler.Add(w, r, ch)
		})
		router.HandleFunc("/server/{action}", func(w http.ResponseWriter, r *http.Request) {
			HTTPHandler.Server(w, r, servers,store)
		})
		router.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
				// TODO : Put a conditional key here, we need that key to verify that message isn't bogus
				// If a message is recieved on this func, this needs to be sent as an alert ASAP
				log.Print("Message recieved")
				serv := r.FormValue("alert")
				err := r.FormValue("error")
				res := make(map[string]string)
				res["server"] = serv
				res["error"] = err
				alertChannel <- res
				w.Write([]byte("{\"status\":\"posted\"}"))
			})
		http.Handle("/", router)
		http.ListenAndServe(":63000", nil) // TODO add this to conf
	}(_store, _servers, httpChannel)


	for {
		select {
		case item := <-httpChannel:
			// Poll on Request
			for k, _ := range item {
				if k == "server" {
					// Do soemthing
					// Add this info to server map
					ser := make(map[string]interface{})
					ser["name"] = item[k].(map[string]interface{})["name"].(string)
					ser["key"] = item[k].(map[string]interface{})["key"].(string)
					ser["ip"] = item[k].(map[string]interface{})["ip"].(string)
					if _,ok := _servers[item[k].(map[string]interface{})["name"].(string)]; !ok{
						_servers[item[k].(map[string]interface{})["name"].(string)] = ser
					}
					fi, err := os.OpenFile("./cockpit/servers.json", os.O_RDWR, 0660);
					if err != nil { panic(err) }
					// close fi on exit and check for its returned error
					defer func() {
						if err := fi.Close(); err != nil {
							panic(err)
						}
					}()
					js, err := json.Marshal(_servers)
					if err != nil{
						log.Panic("Error Marshalling json")
					}else{
						_,err := fi.Write(js); if err != nil{
							log.Panic("Unable to write to json file")
						}
					}
				}
			}
		case <-poller.C:
			srv := make(map [string]interface {})
			for k, v := range _servers {
				srv[k] = v
			}
			go func(srv map[string]interface {}){
				r := make (map[string]interface {})
				// Poller
				// Now loop over servers to check the status of the server
				// Dump the contents of _server as json file
				for _ ,val := range srv{
					ser := val.(map[string]interface {})
					resp, err := http.Get("http://" + ser["ip"].(string) + ":63001"+"/") // TODO: Add it to conf
					if err != nil {
						log.Print("Error in fetching value for the IP " + ser["ip"].(string))
						resp = nil
					}
					//log.Print(resp)
					if resp != nil {
						var f interface{}
						v, er := ioutil.ReadAll(resp.Body)
						if er != nil {
							log.Print("Error reading value from the server")
						} else {
							er := json.Unmarshal([]byte(v), &f)
							if er != nil {
								log.Print("Error in Unmarshaling data")
								var temp interface{}
								r[ser["name"].(string)] = temp
							} else {
								r[ser["name"].(string)] = f
							}
						}
					}else{
						t := make(map[string]interface {})
						t["status"] = "dead"
						r[ser["name"].(string)] = t
					}
					pollerChannel <- r
				}
			}(srv)



		case alert := <-alertChannel:
			// as of now just print this shit
			data := map[string]interface{}{"server": alert["server"],
				"error": alert["error"]}
			regIDs := []string{"APA91bFLfsr-5lsuunNLXs_qRxCjPW6c9spZIUAnuIXUlwQ0pSEi1U9Ln0ctLjSeD4xT5McsIv2MJSBlQE-v3Y8jm6JtNbzKlWXnK90goeEweQ36t-sGkTf_nVZ-jiB-TfxWS3mQHJxkd1M7PXf5-PWLeWfkOpusESD0z_Q0NMIy2tqAp9X_K7s"}
			msg := gcm.NewMessage(data, regIDs...)
			sender := &gcm.Sender{ApiKey: "AIzaSyCNJeYio3_-9tRF8fjuUep7BeV0c0rE8TU"}
			_, err := sender.Send(msg, 2)
			if err != nil {
				log.Print("Failed to send message:", err)
				return
			}
			log.Print(alert)
		case item := <-pollerChannel:
			// Copy to _server
			for k, v := range item {
				_store[k]= v
			}
		}
	}
}

/*
	This function val is for reading the content of _cfile, in case the value isn't present in
	_confMap. It also inserts the value in _confMap, if it isn't present
*/
func val(sec string, opt string) string {
	if sec == "" || opt == "" {
		log.Fatal("Section or Option value sent to config cannot be empty")
	}
	if _confMap[sec+opt] != "" {
		return _confMap[sec+opt]
	} else {
		// cfile is file object from the init.go
		val, err := _cfile.GetRawString(sec, opt)
		if err != nil {
			// Throw exception but let the call proceed
			log.Fatal("Error while getting value for the [" + sec + "/" + opt + "] : (Error) " + err.Error())
			return "err"
		}
		_confMap[sec+opt] = val
		return val
	}
}

func mailer(subject string, content string) bool {
	type Euser struct {
		User   string
		Pass   string
		Server string
		Port   int
	}
	port, _ := strconv.Atoi(val("email", "port"))
	emailUser := &Euser{val("email", "user"), val("email", "password"), val("email", "server"), port}
	auth := smtp.PlainAuth("",
		emailUser.User,
		emailUser.Pass,
		emailUser.Server,
	)
	ml := strings.Split(val("email", "to"), ",")
	for _, val := range ml {
		tos := strings.Split(val, ":")
		from := mail.Address{"Unbxd Monitoring Tool", emailUser.User}
		to := mail.Address{strings.TrimSpace(tos[0]), strings.TrimSpace(tos[1])}
		title := subject
		body := content
		header := make(map[string]string)
		header["From"] = from.String()
		header["To"] = to.String()
		header["Subject"] = func(String string) string {
			// use mail's rfc2047 to encode any string
			addr := mail.Address{String, ""}
			return strings.Trim(addr.String(), " <>")
		}(title)
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
		header["Content-Transfer-Encoding"] = "base64"

		message := ""
		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
		err := smtp.SendMail(
			emailUser.Server+":"+strconv.Itoa(emailUser.Port),
			auth,
			from.Address,
			[]string{to.Address},
			[]byte(message),
		)
		if err != nil {
			log.Print("Unable to send mail" + err.Error())
		}
	}
	return true
}
