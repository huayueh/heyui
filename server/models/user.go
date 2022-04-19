package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	Acct      string    `gorm:"size:30;not null;primary_key;" json:"acct"`
	Fullname  string    `gorm:"size:50;not null" json:"fullname"`
	Pwd       string    `gorm:"size:100;not null;" json:"pwd"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type UserRep struct {
	Acct      string    `json:"acct"`
	Fullname  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserRep {
	return UserRep{
		Acct:      u.Acct,
		Fullname:  u.Fullname,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) FetchUsers(db *gorm.DB, limit int) (*[]User, error) {
	var users []User
	var size int

	if limit > 0 {
		size = limit
	} else {
		size = 100
	}
	err := db.Debug().Model(&User{}).Limit(size).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}
