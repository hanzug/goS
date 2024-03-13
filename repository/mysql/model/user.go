package model

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/hanzug/goS/consts"
)

type User struct {
	UserID         int64  `gorm:"primarykey"`
	UserName       string `gorm:"unique"`
	NickName       string
	PasswordDigest string
}

// SetPassword 加密密码
func (user *User) SetPassword(password string) error {
	// bcrypt会自动生成一个随机盐，并将这个盐与用户的密码进行组合，然后进行多轮散列处理。
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), consts.PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 检验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
