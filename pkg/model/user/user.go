package user

import (
	"chatgpt-web/pkg/model"
	"chatgpt-web/pkg/password"
)

type User struct {
	model.BaseModel
	Name     string `gorm:"column:name;type:varchar(255);not null;unique" valid:"name"`
	Password string `gorm:"column:password;type:varchar(255);not null" valid:"password"`
	Realname string `gorm:"type:varchar(255);comment:Realname" json:"-"`
	Stat     int    `gorm:"default:0;comment:0:normal 1:ban" json:"-"`
}

// ComparePassword 检查密码是否匹配
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}
