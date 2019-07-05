package main

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Content struct {
	Text string `json:"text"`
}

type messageV3 struct {
	OpenId  string  `json:"open_id"`
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type messageV3Resp struct {
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	OpenMessageId string `json:"open_message_id"`
}
type msgService struct {
	c *cache.Cache
}

func (c *msgService) send(email string, text string) string {
	return sendMsgV3(getOpenId(email, c.c), getToken(c.c), text)
}

//func send(email string ,text string,c*cache.Cache) string {
//	return sendMsgV3(getOpenId(email,c),getToken(c),text)
//}
func sendMsgV3(openId string, token string, text string) string {
	const url = "https://open.feishu.cn/open-apis/message/v3/send/"
	log.Printf(url)
	log.Printf(openId)
	log.Printf(token)
	log.Printf(text)

	s := messageV3{OpenId: openId, MsgType: "text", Content: Content{Text: text}}
	b, _ := json.Marshal(s)
	//fmt.log.Printf(string(b))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Printf(string(body))
	if err != nil {
		// handle error
		fmt.Printf("Error: %s", err)
	}
	var msgV3resp messageV3Resp

	err = json.Unmarshal([]byte(string(body)), &msgV3resp)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	//log.Println(msgV3resp)
	return msgV3resp.Msg
}
