package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"

const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutext        *sync.Mutex
	// 保证线程安全,动态变化，加上互斥锁
}

func (r *registry) add(reg Registration) error {
	// 上锁
	r.mutext.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutext.Unlock()
	return nil
}

func (r *registry) remove(url string) error {
	for i := range r.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutext.Lock()
			// 删除index为i的元素
			reg.registrations = append(reg.registrations[:i], r.registrations[:i+1]...)
			r.mutext.Unlock()
		}
	}
	return fmt.Errorf("no service found with URL: %s", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutext:        new(sync.Mutex),
}

// 空struct
type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		// 声明一个结构体Registration
		var registration Registration
		err := dec.Decode(&registration)
		if err != nil {
			// 记录日志
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Printf("Adding service: %v with URL: %s.\n", registration.ServiceName, registration.ServiceURL)
		// registry.add(registration)
		err = reg.add(registration)
		if err != nil {
			// 记录日志
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	case http.MethodDelete:
		// 获取url
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
