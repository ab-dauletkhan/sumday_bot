package models

import "time"

type Message struct {
	Timestamp time.Time
	Text      string
}
