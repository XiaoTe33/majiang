package model

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
)

const (
	CmdPeng    = "peng"
	CmdHu      = "hu"
	CmdGang    = "gang"
	CmdIgnore  = "ignore"
	CmdCardOut = "out"
	CmdCardIn  = "in"
	CmdReady   = "ready"
	CmdStart   = "start"
	CmdUnready = "unready"
	CmdSay     = "say"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Player struct {
	Id          int64
	Username    string
	Conn        *Client
	Cards       Cards
	CurrentCard Card
	Down        []Cards
	Room        *Room
	WinGroup    Group
	IsOwn       bool
	IsReady     bool
	IsOnline    bool
}

func NewPlayer(id int64, username string, room string, conn *websocket.Conn) *Player {
	p := &Player{
		Id:          id,
		Username:    username,
		Conn:        NewClient(username, room, conn),
		CurrentCard: Card{},
		Room:        nil,
		IsOwn:       false,
		IsReady:     false,
		IsOnline:    false,
		Down:        []Cards{},
	}
	p.JoinRoom(room)
	p.Connect()
	return p
}

func (p *Player) JoinRoom(id string) error {
	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if len((&rooms[i]).Players) >= 4 {
		return errors.New("房间已满")
	}
	if len((&rooms[i]).Players) == 0 {
		p.IsOwn = true
	}
	//myLog.Info(p.Username, "enter the room", (&rooms[i]).Id)
	(&rooms[i]).Players = append((&rooms[i]).Players, p)
	//myLog.Info((&rooms[i]).Players[0].Username)
	p.Room = &rooms[i]

	return nil
}

// Connect 开启协程接收消息
func (p *Player) Connect() {
	go func() {
		for {
			select {
			case cmd := <-p.Conn.CommandChan:
				p.ResolveCmd(cmd)
				p.Room.Info(p)
			default:
			}

		}
	}()
}

func (p *Player) Say(content ...any) {
	p.Room.Messages = append(p.Room.Messages, fmt.Sprint("\x1b[36m"+fmt.Sprint("[player]"+p.Username+":")+fmt.Sprintln(content...)+"\x1b[0m"))
}

// StartGame 房主开始游戏
func (p *Player) StartGame() {
	for _, player := range p.Room.Players {
		if !player.IsReady {
			p.Say("请尽快准备哦")
			return
		}
	}
	if p.IsOwn {
		p.Room.Start()
	}

}

// ResolveCmd 处理消息
func (p *Player) ResolveCmd(cmd string) {
	myLog.Info(p.Username, "cmd "+cmd)
	cmdSlice := strings.Split(cmd, ",")
	myLog.Info("cmd[0]", cmdSlice[0])
	switch cmdSlice[0] {
	case CmdCardOut:
		i, _ := strconv.Atoi(cmdSlice[1])
		p.CardOut(i)
	case CmdCardIn:
		p.CardIn()
	case CmdPeng:
		i, _ := strconv.Atoi(cmdSlice[1])
		p.Peng(i)
	case CmdGang:
		p.Gang()
	case CmdHu:
		p.Hu()
	case CmdIgnore:
		p.Ignore()
	case CmdReady:
		p.Ready()
	case CmdStart:
		p.StartGame()
	case CmdUnready:
		p.Unready()
	case CmdSay:
		content := ""
		for _, s := range cmdSlice[1:] {
			content += " " + s
		}
		p.Say(content)
	}
}

// Hu 胡牌
func (p *Player) Hu() {
	if p.CanHu() {
		p.Room.StoreCmd(CmdHu)
		p.Room.Over(p)
	} else {
		p.Room.Say("还不能胡")
	}

}

// CardOut 出牌
func (p *Player) CardOut(num int) {
	if num == 0 {
		p.Room.CurrentCard = p.CurrentCard
		p.Room.NextTurn()
		return
	}
	p.Room.CurrentCard = p.Cards[num-1]
	p.Cards = append(p.Cards.Remove(num-1), p.CurrentCard)
	SortCards(&p.Cards)
	p.Room.NextTurn()

}

// CardIn 抓牌
func (p *Player) CardIn() {
	p.CurrentCard = p.Room.RemainingCards[0]
	p.Room.NextCard()
}

// Peng 碰
func (p *Player) Peng(num int) {
	if !p.CanPeng() {
		p.Room.Say("不能碰哦")
		return
	}
	var rm []int
	for i, card := range p.Cards {
		if p.Room.CurrentCard.Int() == card.Int() {
			rm = append(rm, i)
			if len(rm) == 2 {
				break
			}
		}
	}
	rm = append(rm, num-1)
	p.Down = append(p.Down, Cards{p.Room.CurrentCard, p.Room.CurrentCard, p.Room.CurrentCard})
	p.Room.CurrentCard = p.Cards[num-1]
	p.Cards = p.Cards.Remove(rm...)
	p.Room.PlayerTurns = p.Room.Index(p)
	SortCards(&p.Cards)
	p.Room.NextCard()
	p.Room.StoreCmd(CmdPeng)
}

// Gang 杠
func (p *Player) Gang() {
	if !p.CanGang() {
		p.Room.Say("不能杠哦")
	}
	currentCard := Card{}
	if !p.Room.IsOperateTurn() {
		currentCard = p.CurrentCard
	} else {
		currentCard = p.Room.CurrentCard
	}
	var rm []int
	for i, card := range p.Cards {
		if currentCard.Int() == card.Int() {
			rm = append(rm, i)
		}
	}
	p.Down = append(p.Down, Cards{currentCard, currentCard, currentCard, currentCard})
	p.Cards = p.Cards.Remove(rm...)
	p.Room.PlayerTurns = p.Room.Index(p)
	p.Room.NextTurn()
	p.Room.NextCard()
	p.Room.StoreCmd(CmdGang)
}

// Ready 准备
func (p *Player) Ready() {
	p.IsReady = true
}

// Unready 取消准备
func (p *Player) Unready() {
	p.IsReady = false
}

// IsYouTurn 当前是你的回合吗
func (p *Player) IsYouTurn() bool {
	return p.Room.Players[p.Room.PlayerTurns] == p && p.Room.Turns%2 == 0
}

// Ignore 过
func (p *Player) Ignore() {
	p.Room.StoreCmd(CmdIgnore)
}

// CanHu 判断牌是否能胡
func (p *Player) CanHu() bool {
	if !p.IsYouTurn() && !p.Room.IsOperateTurn() {
		return false
	}
	result := NewChecker().Check(p.AllCards())
	for _, g := range result {
		if len(g.Remain) != 2 {
			return false
		}
		if g.Remain[0].Int() == g.Remain[1].Int() {
			return true
		}
	}
	return false
}

// CanGang 判断是否能杠
func (p *Player) CanGang() bool {
	currentCard := Card{}
	if !p.Room.IsOperateTurn() {
		currentCard = p.CurrentCard
	} else {
		currentCard = p.Room.CurrentCard
	}
	t := 0
	for _, card := range p.Cards {
		if currentCard.Int() == card.Int() {
			t++
		}
		if t == 3 {
			return true
		}
	}
	return false
}

// CanPeng 判断是否能碰
func (p *Player) CanPeng() bool {
	if !p.Room.IsOperateTurn() {
		return false
	}
	t := 0
	for _, card := range p.Cards {
		if p.Room.CurrentCard.Int() == card.Int() {
			t++
		}
		if t == 2 {
			return true
		}
	}
	return false
}

func (p *Player) AllCards() Cards {
	all := append(p.Cards, p.CurrentCard)
	SortCards(&all)
	return all
}
