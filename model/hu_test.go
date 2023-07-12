package model_test

import (
	"fmt"
	"majiang/model"
	"testing"
)

func TestChecker_Check(t *testing.T) {
	type group struct {
		input model.Cards
		want  []*model.Group
	}
	groups := []group{
		{
			model.Cards{l1},
			[]*model.Group{
				{[]model.Cards{}, []model.Cards{}, []model.Cards{}, model.Cards{l1}},
				{[]model.Cards{}, []model.Cards{}, []model.Cards{}, model.Cards{l1}},
				{[]model.Cards{}, []model.Cards{}, []model.Cards{}, model.Cards{l1}},
				{[]model.Cards{}, []model.Cards{}, []model.Cards{}, model.Cards{l1}},
			},
		},
		{
			model.Cards{l1, l1, l1, l2, l3, l5, l6, l7, l8, l8, l8, l9, l9, l9},
			[]*model.Group{
				{
					[]model.Cards{{l1, l1, l1}, {l8, l8, l8}, {l9, l9, l9}},
					[]model.Cards{{l5, l6, l7}},
					[]model.Cards{},
					model.Cards{l2, l3},
				},
				{
					[]model.Cards{{l8, l8, l8}, {l9, l9, l9}},
					[]model.Cards{{l1, l2, l3}, {l5, l6, l7}},
					[]model.Cards{},
					model.Cards{l1, l1},
				},
				{
					[]model.Cards{},
					[]model.Cards{{l9, l8, l7}, {l3, l2, l1}},
					[]model.Cards{},
					model.Cards{l1, l1, l5, l6, l8, l8, l9, l9},
				},
				{
					[]model.Cards{{l1, l1, l1}, {l8, l8, l8}, {l9, l9, l9}},
					[]model.Cards{{l7, l6, l5}},
					[]model.Cards{},
					model.Cards{l3, l2},
				},
			},
		},
	}

	for i := range groups {
		got := model.NewChecker().Check(groups[i].input)
		for i2 := range got {
			for i3 := range (*got[i2]).AAA {
				for i4 := range (*got[i2]).AAA[i3] {
					if (*got[i2]).AAA[i3][i4].Int() != (*(groups[i].want[i2])).AAA[i3][i4].Int() {
						goto Fail
					}
				}
			}
			for i3 := range (*got[i2]).BCD {
				for i4 := range (*got[i2]).BCD[i3] {
					if (*got[i2]).BCD[i3][i4].Int() != (*(groups[i].want[i2])).BCD[i3][i4].Int() {
						goto Fail
					}
				}
			}
			for i3 := range (*got[i2]).EE {
				for i4 := range (*got[i2]).EE[i3] {
					if (*got[i2]).EE[i3][i4].Int() != (*(groups[i].want[i2])).EE[i3][i4].Int() {
						goto Fail
					}
				}
			}
			for i3 := range (*got[i2]).Remain {
				if got[i2].Remain[i3].Int() != (groups[i].want[i2]).Remain[i3].Int() {
					goto Fail
				}
			}
		}
	}
	return
Fail:
	//打印
	t.Error("all cards are as follow:")
	fmt.Println("got:want")
	for i := range groups {
		got := model.NewChecker().Check(groups[i].input)
		for i2 := range got {
			for i3 := range (*got[i2]).AAA {
				for i4 := range (*got[i2]).AAA[i3] {
					fmt.Println((*got[i2]).AAA[i3][i4], ":", (*(groups[i].want[i2])).AAA[i3][i4])
				}
			}
			for i3 := range (*got[i2]).BCD {
				for i4 := range (*got[i2]).BCD[i3] {
					fmt.Println((*got[i2]).BCD[i3][i4], ":", (*(groups[i].want[i2])).BCD[i3][i4])
				}
			}

			for i3 := range (*got[i2]).Remain {
				fmt.Println(got[i2].Remain[i3], ":", (groups[i].want[i2]).Remain[i3])
			}
			fmt.Println()
		}
		fmt.Println()
	}
	t.Fail()
}
