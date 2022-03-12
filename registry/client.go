package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	// 处理服务发现请求
	http.Handle(serviceUpdateURL.Path, &serviceUpdateURLHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return nil
	}

	//
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service "+" responded with code %v", res.StatusCode)
	}
	if err != nil {
		return err
	}
	return nil
}

type serviceUpdateURLHandler struct{}

func (suh serviceUpdateURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to shutdown service. Registry service "+"responded with code %v", res.StatusCode)
	}
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// 新增的部分
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		// 存在就直接加在后面即可
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		if providerURLS, ok := p.services[patchEntry.Name]; ok {
			for i, url := range providerURLS {
				if url == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerURLS[:i], providerURLS[i+1:]...)
					break
				}
			}
		}
	}

}

// 在本系统中只有一个service， grade
func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no provider found for service %s", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
