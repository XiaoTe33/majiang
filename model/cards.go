package model

import (
	"sort"
	"strconv"
)

var typeMap = map[int]string{
	1: "W", //万
	2: "L", //条
	3: "T", //筒
}

type Card struct {
	Type int //花色
	Num  int //数字
}

func SortCards(cards *Cards) {
	var intSlice []int
	for _, card := range *cards {
		intSlice = append(intSlice, card.Type*10+card.Num)
	}
	sort.Ints(intSlice)
	var res Cards
	for _, i := range intSlice {
		res = append(res, Card{Type: i / 10, Num: i % 10})
	}
	*cards = res
}
func SortCardsDesc(cards *Cards) {
	var intSlice []int
	for _, card := range *cards {
		intSlice = append(intSlice, card.Type*10+card.Num)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(intSlice)))
	var res Cards
	for _, i := range intSlice {
		res = append(res, Card{Type: i / 10, Num: i % 10})
	}
	*cards = res
}

func (c Card) Int() int {
	return 10*c.Type + c.Num
}

func (c Card) String() string {
	return typeMap[c.Type] + strconv.Itoa(c.Num)
}
