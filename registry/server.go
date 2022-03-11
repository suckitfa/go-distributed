package registry

import (
	"encoding/json"
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
