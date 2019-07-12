package receiver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"strconv"
)

type MessageReceiver interface {
	FillData(b []byte)
	//FillData(s string)
	OnReceive(c *websocket.Conn, done *chan struct{})
}

type GeneralMessageReceiver struct {
	Command string `json:"command"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type HelloBackMessageReceiver struct {
	Question          string `json:"question"`
	InPlay            int    `json:"inplay"`
	NumberOfQuestions int    `json:"numberofquestions"`
	StreamUrl         string `json:"streamurl"`
	LifePrice         int    `json:"lifeprice"`
	HasLife           int    `json:"haslife"`
	UsedLife          int    `json:"usedlife"`
	Prize             int    `json:"prize"`
}

type QuestionMessageReceiver struct {
	Question       string `json:"question"`
	Answer1        string `json:"ans1"`
	Answer2        string `json:"ans2"`
	Answer3        string `json:"ans3"`
	QuestionNumber int    `json:"qnumber"`
	QuestionType   string `json:"questionType"`
}

func (data *GeneralMessageReceiver) FillData(b []byte) {
	error := json.Unmarshal(b, data)
	if error != nil {
		fmt.Println(error)
		return
	}

	return
}
func (data *HelloBackMessageReceiver) FillData(b []byte) {
	error := json.Unmarshal(b, data)
	if error != nil {
		fmt.Println(error)
		return
	}

	return
}

func (data *QuestionMessageReceiver) FillData(b []byte) {
	error := json.Unmarshal(b, data)
	if error != nil {
		fmt.Println(error)
		return
	}

	return
}

func (data GeneralMessageReceiver) OnReceive(c *websocket.Conn, done *chan struct{}) {
	//fmt.Println(data)
}
func (data HelloBackMessageReceiver) OnReceive(c *websocket.Conn, done *chan struct{}) {
	// do something
}
func (data QuestionMessageReceiver) OnReceive(c *websocket.Conn, done *chan struct{}) {
	max := 4
	min := 1
	answer := rand.Intn(max-min) + min

	if answer > 3 {
		answer = 3
	}

	message := "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTAzOTg3LCJpYXQiOjE1NTI4MzQ2MzcsImV4cCI6MTYxNTkwNjYzN30.00hCKmHkkVkZJyCXP0SUiEXGT-EbHqQJTwDH6DBVTRI\",\"action\":\"ANSWER\",\"debugstring\":\"\",\"question\":" + strconv.Itoa(data.QuestionNumber) + ",\"answer\":" + strconv.Itoa(answer) + "}"
	fmt.Println(message)

	err := c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write:", err)
		return
	}
}
