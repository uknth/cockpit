package HTTPHandler

import(
	"net/http"
	"log"
	"memcached"
	"encoding/json"
	"github.com/gorilla/mux"
	"errors"
)


func Auth(w http.ResponseWriter,r *http.Request){
	user := r.FormValue("user")
	token := r.FormValue("token")
	// now put it in memcached
	if token != "" && user != "" {
		// TODO: We can add multiple server, but i guess we don't need them
		err := memcached.Set(user,token)
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



func validateRequest(w http.ResponseWriter,r *http.Request) error{
	// Make a request to memcached and get the token based on key
	user := r.FormValue("user")
	token := r.FormValue("token")
	if user != "" && token != "" {
		// Check it the token exists in memcached
		val, err := memcached.Get(user)
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
		http.Error(w,"{\"status\":\"Bad Request\"}" , 400)
		return errors.New("Bad Request")
	}
	return nil
}

