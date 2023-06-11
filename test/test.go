package main

import (
	"fmt"
	"majiang/model"
	"time"
)

func main() {
	a := &model.Player{Username: "1"}
	a.Say("1111")
	r := model.Room{}
	r.Say("222")
}
func main02() {
	cards := model.Cards{
		{1, 2},
		{1, 3},
		{1, 2},
		{1, 3},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
		{2, 7},
		{2, 8},
		{2, 9},
	}
	groups := model.NewChecker().Check(cards)
	for _, group := range groups {
		fmt.Println(group)
	}
	time.Sleep(6 * time.Second)
	fmt.Print("\033[H\033[2J")
}

func main01() {
	pre := model.Cards{
		{1, 2},
		{1, 3},
		{1, 2},
		{1, 3},
		{1, 2},
		{1, 3},
	}
	model.SortCards(&pre)
	fmt.Println(pre)

	model.SortCardsDesc(&pre)
	fmt.Println(pre)
}
