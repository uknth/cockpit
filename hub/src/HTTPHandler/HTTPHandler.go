package HTTPHandler

import(
	"net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"errors"
)


func Auth(w http.ResponseWriter,r *http.Request){
	user := r.FormValue("USER_ID")
	token := r.FormValue("TOKEN")
	// now put it in memcached
	if token != "" && user != "" {
		// TODO: We can add multiple server, but i guess we don't need them
		err := MemSet(user,token)
		if err != nil {
			log.Panic("Error writing data to Memcached." + err.Error())
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte("{\"status\":\"success\"}"))
	}else{
		http.Error(w, "{\"status\":\"Invalid Data\"}", 400)
	}
}


func Add (w http.ResponseWriter,r *http.Request , ch chan map[string]interface{}){
	err := validateRequest(w,r)
	if err != nil{
		return
	}
	vars := mux.Vars(r)
	action := vars["action"]
	var f interface{}
	// We get the value as formValue, with a prefix of 'products'
	err = json.Unmarshal([]byte(r.FormValue("input")), &f)
	if err != nil {
		log.Panic("Invalid Json data sent" + err.Error())
		http.Error(w,"{\"status\":\"Invalid Data\"}" , 400)
	}
	res := make(map[string]interface{})
	res[action] = f
	ch <- res
	w.Write([]byte("{\"status\":\"success\"}"))
}


func Server(w http.ResponseWriter,r *http.Request ,servers map[string]interface {},store map[string]interface {}){
	err := validateRequest(w,r)
	if err != nil{
		return
	}
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "list" {
		// Return the list here
		type Ser struct {
			Name string
			Key string
			Ip string
			Status string
		}
		type Res struct{
			Servers []Ser
		}
		var servs []Ser
		for _,v := range servers{
			name := v.(map[string]interface {})["name"].(string)
			key := v.(map[string]interface {})["key"].(string)
			ip := v.(map[string]interface {})["ip"].(string)
			var status string
			if _,ok := store[name] ; !ok{
				status = "unknown"
			}else{

				status = store[name].(map[string]interface {})["status"].(string)
			}
			ser := Ser{name,key,ip, status}
			servs = append(servs,ser)
		}
		resp := Res{servs}
		jsn,err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("{\"status\":\"error\"}"))
			return
		}
		w.Write(jsn)
		return
	}
}


func validateRequest(w http.ResponseWriter,r *http.Request) error{
	// Make a request to memcached and get the token based on key
	user := r.Header.Get("USER_ID")
	token := r.Header.Get("TOKEN")
	if user != "" && token != "" {
		// Check it the token exists in memcached
		val, err := MemGet(user)
		if err != nil{
			// Throw panic
			log.Panic()
		}else{
			if val != token{
				// Valid request
				// Break the request
				http.Error(w,"{\"status\":\"Logged Out\"}" , 415)
				return  errors.New("Logged Out")
			}
		}
	}else{
		log.Print("U:"+user + "  T:" + token)
		http.Error(w,"{\"status\":\"Bad Request\"}" , 400)
		return errors.New("Bad Request")
	}
	return nil
}

