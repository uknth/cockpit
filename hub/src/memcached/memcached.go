package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"log"
)


var mc *memcache.Client

func Init(){
	mc = memcache.New("localhost:11211") // TODO: Need to add it to conf
}

func Get(key string) (string ,error){
	val, err := mc.Get(key)
	if err != nil{
		log.Panic("Error getting value from memcached")
	}
	return string(val.Value),nil
}



func Set(key string , val string) error{
	err := mc.Set(&memcache.Item{Key:key , Value:[]byte(val)})
	if err != nil{
		log.Panic("Error writing data to memcached")
		return err
	}
	return nil
}


