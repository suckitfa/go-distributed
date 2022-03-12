package registry

import (
	"bytes"
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
	mutex         *sync.RWMutex
	// 保证线程安全,动态变化，加上互斥锁
}

func (r *registry) add(reg Registration) error {
	// 上锁
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// 发送请求:通知所有的依赖服务
	err := r.sendRequiredServices(reg)
	// 发送请求:通知服务更新
	r.notify(patch{
		Added: []pathEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return err
}

func (r registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		// 并发发送通知
		go func(reg Registration) {
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []pathEntry{}, Removed: []pathEntry{}}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
					}
				}
			}

		}(reg)
	}

}
func (r registry) sendRequiredServices(reg Registration) error {
	// 上读的锁
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		if serviceReg.ServiceName == reg.ServiceName {
			p.Added = append(p.Added, pathEntry{
				Name: serviceReg.ServiceName,
				URL:  serviceReg.ServiceURL,
			})
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	return err
}

func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// 发送请求
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	return err
}

func (r *registry) remove(url string) error {
	for i := range r.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			// 删除index为i的元素
			reg.registrations = append(reg.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
		}
	}
	return fmt.Errorf("no service found with URL: %s", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
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
