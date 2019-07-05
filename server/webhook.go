package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/patrickmn/go-cache"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	m     msgService
	ipStr string
)

func serviceInit() {
	m = msgService{cache.New(1*time.Hour, 2*time.Hour)}
	ipStr = getInternalIP()
	m.send("caiyingrong@bytedance.com", InitMsg())
	//m.send("caiyingrong@bytedance.com",ExitMsg())

	logs.Info("Init")
}

func Parse(w http.ResponseWriter, req *http.Request) {
	//把  body 内容读入字符串 s
	s, _ := ioutil.ReadAll(req.Body)

	//logs.Info(string(s))
	gmsg := parseJsonGrafana(string(s))
	logs.Info(gmsg.Msg())

	m.send("caiyingrong@bytedance.com", gmsg.Msg())
	//fmt在返回页面中显示内容。
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))

	//处理结束
	log.Printf("Parse Over\n")
}
func service() {
	logs.Info("Run in %s mode", "debug")
	http.HandleFunc("/", Parse)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		logs.Error("ListenAndServe: ", err)
	}
}
func waitSignal() {
	signals := make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1,
		syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGPIPE)

WaitSignal:
	for {
		select {
		case sig := <-signals:
			logs.Error("get a signal [%s]\n", sig)
			m.send("caiyingrong@bytedance.com", ExitMsg())
			log.Printf("Exiting")
			break WaitSignal
		}
	}
}

func main() {
	//初始化服务
	serviceInit()

	//启动服务
	go service()

	//优雅退出
	waitSignal()
}
