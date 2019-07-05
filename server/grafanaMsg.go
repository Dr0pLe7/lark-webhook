package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

type tags struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
}
type evalMatch struct {
	Value  float64 `json:"value"`
	Metric string  `json:"metric"`
	Tags   tags    `json:"tags"`
}

type grafanaMsg struct {
	EvalMatches []evalMatch `json:"evalMatches"`
	ImageUrl    string      `json:"imageUrl"`
	Message     string      `json:"message"`
	RuleId      int         `json:"ruleId"`
	RuleName    string      `json:"ruleName"`
	RuleUrl     string      `json:"ruleUrl"`
	State       string      `json:"state"`
	Title       string      `json:"title"`
}

func (i *evalMatch) Msg() string {
	return "(" + i.Metric + ":" + fmt.Sprintf("%f", i.Value) + ")"
}
func parseJsonGrafana(str string) grafanaMsg {
	var gmsg grafanaMsg
	err := json.Unmarshal([]byte(str), &gmsg)
	if err != nil {
		logs.Error("Error: %s\n", err)
	}
	return gmsg
}

func (g *grafanaMsg) Msg() string {
	evalMsg := ""
	for _, i := range g.EvalMatches {
		evalMsg += i.Msg()
	}
	return "【" + g.Title + "】\n" + "--" + g.RuleName + "\n" + "「" + g.Message + "」\n" + "{" + evalMsg + "}\n" + g.ImageUrl
}

func InitMsg() string {
	return time.Now().Format("2006-01-02 15:04:05") + "==>\tStart WebHook Service @ " + ipStr + "\n"
}
func ExitMsg() string {
	return time.Now().Format("2006-01-02 15:04:05") + "==>\tStop  WebHook Service @ " + ipStr + "\n"
}

//func main(){
//	str:=`
//	{"evalMatches":[{"value":100,"metric":"High value","tags":null},{"value":200,"metric":"Higher Value","tags":null}],"imageUrl":"https://grafana.com/assets/img/blog/mixed_styles.png","message":"Someone is testing the alert notification within grafana.","ruleId":0,"ruleName":"Test notification","ruleUrl":"http://localhost:3000/","state":"alerting","title":"[Alerting] Test notification"}
//	`
//	g := parseJsonGrafana(str)
//	fmt.Println(g.Msg())
//}
