package main

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type email2idBody struct {
	Email string `json:"email"`
}
type email2idResp struct {
	Code       int    `json:"code"`
	EmployeeId string `json:"employee_id"`
	Msg        string `json:"msg"`
	OpenId     string `json:"open_id"`
	//"code": 0,
	//"employee_id": "63fbee7g",
	//"msg": "ok",
	//"open_id": "ou_9b0f2e2262a9b40f125da2d7ae162932"
}

func email2openId(email string, token string) string {
	const url = "https://open.feishu.cn/open-apis/user/v3/email2id"
	s := email2idBody{Email: email}
	b, _ := json.Marshal(s)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	//client.Do(req)
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Printf("Error: %s", err)
	}
	var user email2idResp
	log.Println(url)
	log.Println(string(body))

	err = json.Unmarshal([]byte(string(body)), &user)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	//log.Println(user)
	return user.OpenId
}

func getOpenId(email string, c *cache.Cache) string {
	var openId string
	t, found := c.Get("open_id_" + email)
	if found {
		openId = t.(string)
	} else {
		openId = email2openId(email, getToken(c))
		c.Set("open_id_"+email, openId, time.Duration(c.Items()["token"].Expiration))
	}
	return openId
}
