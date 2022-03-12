package log

import (
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
)

var log *stdlog.Logger

// 将日志写入文件系统
type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
	fmt.Println("func Write in Log Server")
	// 打开文件，写入模式，0666快平台权限
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return 0, err
	}
	// 建议先处理错误后再关闭文件
	defer f.Close()
	return f.Write(data)
}

func Run(destination string) {
	// 创建日志文件 ???
	fmt.Println("func Run in Log Server")
	log = stdlog.New(fileLog(destination), "[go]: ", stdlog.LstdFlags)
}

// 注册处理函数
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// 读取body内容
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// 输出日志
func write(message string) {
	fmt.Println("private func write  in Log Server")
	log.Printf("%v\n", message)
}
