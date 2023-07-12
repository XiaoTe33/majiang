package model_test

import (
	"majiang/model"
	"testing"
)

var (
	w1 = model.Card{Type: 1, Num: 1}
	w2 = model.Card{Type: 1, Num: 2}
	w3 = model.Card{Type: 1, Num: 3}
	w4 = model.Card{Type: 1, Num: 4}
	w5 = model.Card{Type: 1, Num: 5}
	w6 = model.Card{Type: 1, Num: 6}
	w7 = model.Card{Type: 1, Num: 7}
	w8 = model.Card{Type: 1, Num: 8}
	w9 = model.Card{Type: 1, Num: 9}
	l1 = model.Card{Type: 2, Num: 1}
	l2 = model.Card{Type: 2, Num: 2}
	l3 = model.Card{Type: 2, Num: 3}
	l4 = model.Card{Type: 2, Num: 4}
	l5 = model.Card{Type: 2, Num: 5}
	l6 = model.Card{Type: 2, Num: 6}
	l7 = model.Card{Type: 2, Num: 7}
	l8 = model.Card{Type: 2, Num: 8}
	l9 = model.Card{Type: 2, Num: 9}
	t1 = model.Card{Type: 3, Num: 1}
	t2 = model.Card{Type: 3, Num: 2}
	t3 = model.Card{Type: 3, Num: 3}
	t4 = model.Card{Type: 3, Num: 4}
	t5 = model.Card{Type: 3, Num: 5}
	t6 = model.Card{Type: 3, Num: 6}
	t7 = model.Card{Type: 3, Num: 7}
	t8 = model.Card{Type: 3, Num: 8}
	t9 = model.Card{Type: 3, Num: 9}
)

func TestSortCards(t *testing.T) {
	type group struct {
		input *model.Cards
		want  *model.Cards
	}

	g := []group{
		{
			input: &model.Cards{l9, l3, l4, l7, w3, w5, w6, w4, t4, t1, t2, t2, t9},
			want:  &model.Cards{w3, w4, w5, w6, l3, l4, l7, l9, t1, t2, t2, t4, t9},
		},
	}
	for i := range g {
		model.SortCards(g[i].input)
		for i2 := range *g[i].input {
			if (*g[i].input)[i2].Int() != (*g[i].want)[i2].Int() {
				t.Fail()
			}
		}
	}
}

func TestSortCardsDesc(t *testing.T) {
	type group struct {
		input *model.Cards
		want  *model.Cards
	}

	g := []group{
		{
			input: &model.Cards{l9, l3, l4, l7, w3, w5, w6, w4, t4, t1, t2, t2, t9},
			want:  &model.Cards{t9, t4, t2, t2, t1, l9, l7, l4, l3, w6, w5, w4, w3},
		},
	}
	for i := range g {
		model.SortCardsDesc(g[i].input)
		for i2 := range *g[i].input {
			if (*g[i].input)[i2].Int() != (*g[i].want)[i2].Int() {
				t.Fail()
			}
		}
	}
}
