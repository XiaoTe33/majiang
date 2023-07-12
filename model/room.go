package model

import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct {
	Id             string
	IsGaming       bool
	Owner          *Player
	Players        []*Player
	CurrentCard    Card
	RemainingCards Cards
	PastCards      Cards
	Messages       []string
	CommandNum     int

	PlayerTurns int
	// 奇数中间操作，偶数玩家操作
	Turns int
}

func NewCards() Cards {
	myLog.Info("making new Cards ......")
	var res Cards
	var pre Cards
	idx := 0
	for i := 1; i < 10; i++ {
		for j := 1; j < 4; j++ {
			for k := 0; k < 4; k++ {
				pre = append(pre, Card{
					Type: j,
					Num:  i,
				})
				idx++
			}
		}
	}
	for i := 0; i < 108; i++ {
		rand.Seed(time.Now().UnixNano())
		randIdx := rand.Intn(len(pre))
		res = append(res, pre[randIdx])
		pre = pre.Remove(randIdx)
	}
	return res
}

var rooms = [1000]Room{}

func Rooms(idx int) Room {
	return rooms[idx]
}

// Info 信息打印
func (r *Room) Info(p *Player) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("turns:", r.Turns, " playerTurns:", r.PlayerTurns+1, " remainingCards:", len(r.RemainingCards))
	myDeskNum := 0
	for i, player := range r.Players {
		if player == p {
			myDeskNum = i + 1
			continue
		}
		fmt.Println("【桌", i+1, "】(用户名:"+player.Username, ")剩余", len(player.Cards), "张")
		fmt.Println("desktop:", player.Down)
	}
	fmt.Println("【桌", myDeskNum, "】(用户名:"+p.Username, "【我】)剩余", len(p.Cards), "张")
	fmt.Println("desktop:", p.Down)

	fmt.Print("myCards:", p.Cards)
	if p.IsYouTurn() {
		fmt.Print(" [", p.CurrentCard, "]")
	}
	fmt.Println()
	if r.IsOperateTurn() {
		fmt.Println("currentCards [", r.CurrentCard, "]")
	}
	fmt.Println("pastCards", r.PastCards)
	fmt.Println("\n【消息记录】")
	for i := len(r.Messages) - 6; i < len(r.Messages); i++ {
		if len(r.Messages) == 0 {
			break
		}
		if i < 0 {
			i = 0
		}
		fmt.Print(r.Messages[i])
	}
	//for _, player := range r.Players {
	//	myLog.Info(player.Username, player.IsOwn, player.IsReady)
	//}
}

// Say 系统消息
func (r *Room) Say(content ...any) {
	r.Messages = append(r.Messages, fmt.Sprint("\x1b[31m"+fmt.Sprint("[system] ")+fmt.Sprintln(content...)+"\x1b[0m"))
}

// StoreCmd 用于处理间隔ignore操作
func (r *Room) StoreCmd(cmd string) {
	if cmd != CmdIgnore {
		r.CommandNum = 0
		return
	}
	r.CommandNum++
	if r.CommandNum == 3 {
		r.NextTurn()
		r.NextPlayerTurn()
		r.CommandNum = 0
	}
}

func (r *Room) Start() {

	newCards := NewCards()
	r.Turns = 0
	r.PlayerTurns = 0
	r.IsGaming = true
	r.Players[0].Cards = newCards[0:13]
	r.Players[1].Cards = newCards[13:26]
	r.Players[2].Cards = newCards[26:39]
	r.Players[3].Cards = newCards[39:52]
	SortCards(&r.Players[0].Cards)
	SortCards(&r.Players[1].Cards)
	SortCards(&r.Players[2].Cards)
	SortCards(&r.Players[3].Cards)
	r.RemainingCards = newCards[52:]

	r.Say("新的一局开始了")

}
func (r *Room) Over(p *Player) {
	r.Say("[player][" + p.Username + "]胡了")
	r.IsGaming = false
}

// NextTurn 游戏换下一个轮次
func (r *Room) NextTurn() {
	if r.IsOperateTurn() {
		r.PastCards = append(r.PastCards, r.CurrentCard)
	}
	r.Turns++
	if len(r.RemainingCards) == 0 {
		r.Over(r.Players[r.PlayerTurns%4])
	}
}

// NextPlayerTurn 玩家轮回
func (r *Room) NextPlayerTurn() {
	if r.PlayerTurns >= 3 {
		r.PlayerTurns = 0
	} else {
		r.PlayerTurns++
	}
	r.Players[r.PlayerTurns].CardIn()
}

func (r *Room) NextCard() {
	r.RemainingCards = r.RemainingCards.Remove(0)
}

// PlayerPeng 有玩家碰
func (r *Room) PlayerPeng(p *Player) {
	r.NextTurn()
	r.PlayerTurns = r.Index(p)
	r.NextPlayerTurn()
}

// PlayerGang 有玩家杠
func (r *Room) PlayerGang(p *Player) {
	r.Turns += 2
	r.PlayerTurns = r.Index(p)
}

// Index 获取玩家桌子序号
func (r *Room) Index(p *Player) int {
	for i, player := range r.Players {
		if p == player {
			return i
		}
	}
	return -1
}

func (r *Room) IsOperateTurn() bool {
	return r.Turns%2 == 1
}
