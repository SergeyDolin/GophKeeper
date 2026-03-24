package models

import (
	"time"
)

type Secret struct {
	ID        string
	UserLogin string
	Type      string
	Data      []byte
	Meta      string
	UpdatedAt time.Time
}
