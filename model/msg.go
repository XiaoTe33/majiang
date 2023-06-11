package model

import "time"

type Msg struct {
	Type    int       `json:"type"`
	Time    time.Time `json:"-"`
	Content string    `json:"content"`
}
