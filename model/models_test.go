package model_test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"majiang/model"
	"majiang/router"
	"testing"
	"time"
)

func command(cmd string) map[string]any {
	return map[string]any{
		"type":    1,
		"content": cmd,
	}
}

var tokens = []string{
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2ODkyMTI5MjcsImlhdCI6MTY4OTEyNjUyN30.pYmTsolgdbM_MWNbiW1f7QmIQH6qp6MeGGOrnL0keJs",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI3Mjc1OTMxNjQ4LCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2ODkyMTI5NTAsImlhdCI6MTY4OTEyNjU1MH0.AwNQ6sDTL5lQrw8sd-eksPTFXqd3-MjHj-kq3EJIGmk",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI4MzAxNjAyODE2LCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2ODkyMTI5NzcsImlhdCI6MTY4OTEyNjU3N30.wd9iq5L5r0kg31JbABwykyt-OMdAJjfKHmWPX6bYRoU",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI4MzAxNjM1NTg0LCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2ODkyMTI5OTIsImlhdCI6MTY4OTEyNjU5Mn0.gr_j7MTBfbivDnbomSoDJ9DOvuHjx64J9kY0TeSQeJQ",
}
var (
	isGoing = false
	players []*websocket.Conn
)

func ini(t *testing.T) {
	if !isGoing {
		go router.InitRouters()
		isGoing = true
	}
	if len(players) != 4 {
		for i := 0; i < len(tokens); i++ {
			dialer := &websocket.Dialer{}
			conn, response, err := dialer.Dial("ws://127.0.0.1:8080/join/1?accessToken="+tokens[i], nil)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			fmt.Println("status:", response.Status)
			err = conn.WriteJSON(command("ready"))
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			players = append(players, conn)
		}
	}

}

func TestPlayer_JoinRoom(t *testing.T) {
	ini(t)
	time.Sleep(2 * time.Second)
}

func TestPlayer_StartGame(t *testing.T) {
	ini(t)
	time.Sleep(time.Second)
	err := players[0].WriteJSON(command("start"))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	time.Sleep(2 * time.Second)
}

func TestPlayer_CanGang(t *testing.T) {
	time.Sleep(2 * time.Second)
	ini(t)
	room := model.Rooms(1)
	p := room.Players[room.PlayerTurns]
	p.Cards = model.Cards{l1, l1, l1}
	p.CurrentCard = l1
	if !p.CanGang() {
		t.Errorf("cannot gang when l1,l1,l1,l1")
		t.Fail()
	}
}
func TestPlayer_CanPeng(t *testing.T) {
}
func TestPlayer_CanHu(t *testing.T) {
}
func TestPlayer_Connect(t *testing.T) {
}
