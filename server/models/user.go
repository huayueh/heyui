package models

import (
	"errors"
	"fmt"
	"heyui/server/auth"
	"html"
	"strings"
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

func (u *User) BeforeSave() error {
	hashedPassword, err := auth.Hash(u.Pwd)
	if err != nil {
		return err
	}
	u.Pwd = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Acct = html.EscapeString(strings.TrimSpace(u.Acct))
	u.Fullname = html.EscapeString(strings.TrimSpace(u.Fullname))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if u.Acct == "" {
			return errors.New("account required")
		}
		if u.Fullname == "" {
			return errors.New("fullname required")
		}
		if u.Pwd == "" {
			return errors.New("password required")
		}
		return nil
	case "update":
		if u.Fullname == "" || u.Pwd == "" {
			return errors.New("nothing to update")
		}
		return nil
	case "fullname":
		if u.Fullname == "" {
			return errors.New("fullname required to update")
		}
		return nil
	case "login":
		if u.Acct == "" || u.Pwd == "" {
			return errors.New("missing account or password for login")
		}
		return nil

	default:
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid string) (int64, error) {

	db = db.Debug().Model(&User{}).Where("acct = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) UpdateUser(db *gorm.DB, uid string) (*User, error) {
	// hash the password
	err := u.BeforeSave()
	if err != nil {
		fmt.Print(err)
	}
	db = db.Debug().Model(&User{}).Where("acct = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"pwd":        u.Pwd,
			"fullname":   u.Fullname,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// return updated user
	err = db.Debug().Model(&User{}).Where("acct = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
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

func (u *User) FindUserByFullName(db *gorm.DB, fullname string) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Where("fullname = ?", fullname).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid string) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("acct = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("user not found")
	}
	return u, err
}
