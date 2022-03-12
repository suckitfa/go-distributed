package log

import (
	"bytes"
	"fmt"
	stdlog "log"
	"net/http"

	"test.com/registry"
)

func SetClient(serviceURL string, clientService registry.ServiceName) {
	stdlog.SetPrefix(fmt.Sprintf("[%s] ", clientService))
	stdlog.SetFlags(0)
	stdlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl *clientLogger) Write(data []byte) (n int, err error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to log. Registry service "+"responded with code %v", res.StatusCode)
	}
	return len(data), nil
}
