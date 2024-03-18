package admin

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name            string
	Date            string
	Place           string
	NumberOfTickets int32
}
