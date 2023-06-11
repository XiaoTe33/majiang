package model

type Cards []Card
type Group struct {
	AAA    []Cards
	BCD    []Cards
	EE     []Cards
	Remain Cards
}

func (cs Cards) HasNext(card Card) (idx int) {
	for i, c := range cs {
		if c.Int() == card.Int()+1 {
			return i
		}
	}
	return -1
}

func (cs Cards) HasPre(card Card) (idx int) {
	for i, c := range cs {
		if c.Int() == card.Int()-1 {
			return i
		}
	}
	return -1
}

func (cs Cards) Remove(idx ...int) Cards {
	var dest Cards
	for _, i := range idx {
		cs[i].Type = -1
	}
	for _, c := range cs {
		if c.Type != -1 {
			dest = append(dest, c)
		}
	}
	return dest
}

func (cs Cards) Clone() Cards {
	clone := make(Cards, len(cs))
	copy(clone, cs)
	return clone
}

type ChooseFunc func(group *Group, current Cards) (remain Cards)

func ChooseBCD(group *Group, current Cards) (remain Cards) {
	SortCards(&current)
	idcB := 0
	for {
		if idcB >= len(current) {
			break
		}
		idxC := current.HasNext(current[idcB])
		if idxC == -1 {
			idcB++
			continue
		}
		idxD := current.HasNext(current[idxC])
		if idxD == -1 {
			idcB++
			continue
		}
		group.BCD = append(group.BCD, Cards{current[idcB], current[idxC], current[idxD]})
		current = current.Remove(idcB, idxC, idxD)
	}
	group.Remain = current
	return current
}

func ChooseBCDDesc(group *Group, current Cards) (remain Cards) {
	SortCardsDesc(&current)
	idxB := 0
	for {
		if idxB >= len(current) {
			break
		}
		idxC := current.HasPre(current[idxB])
		if idxC == -1 {
			idxB++
			continue
		}
		idxD := current.HasPre(current[idxC])
		if idxD == -1 {
			idxB++
			continue
		}
		group.BCD = append(group.BCD, Cards{current[idxB], current[idxC], current[idxD]})
		current = current.Remove(idxB, idxC, idxD)
	}
	group.Remain = current
	return current
}

func ChooseAAA(group *Group, current Cards) (remain Cards) {
	SortCards(&current)
	i := 0
	for {
		if i >= len(current)-2 {
			break
		}
		if current[i].Int() == current[i+1].Int() && current[i+1].Int() == current[i+2].Int() {
			group.AAA = append(group.AAA, Cards{current[i], current[i+1], current[i+2]})
			current = current.Remove(i, i+1, i+2)
			continue
		}
		i++
	}
	group.Remain = current
	return current
}

type Checker [][]ChooseFunc

func NewChecker() *Checker {
	return &Checker{
		{ChooseAAA, ChooseBCD},
		{ChooseBCD, ChooseAAA},
		{ChooseBCDDesc, ChooseAAA},
		{ChooseAAA, ChooseBCDDesc},
	}
}

func (c Checker) Check(cards Cards) (result []*Group) {
	for _, chooseFuncs := range c {
		cardsClone := cards.Clone()
		group := &Group{}
		for i := 0; i < len(chooseFuncs); i++ {
			cardsClone = chooseFuncs[i](group, cardsClone)
		}

		result = append(result, group)
	}
	return
}
