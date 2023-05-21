package user

import (
	"chatgpt-web/pkg/model"
	"chatgpt-web/pkg/password"
)

type User struct {
	model.BaseModel
	Name            string `gorm:"column:name;type:varchar(255);not null;unique" valid:"name"`
	Password        string `gorm:"column:password;type:varchar(255);not null" valid:"password"`
	Realname        string `gorm:"type:varchar(255);comment:Realname" json:"-"`
	ExpireTimestamp int64  `gorm:"type:int(11); default:0; comment:有限期 0:永久 >0:到期结束" json:"expire_timestamp"`
	Stat            int    `gorm:"default:0;comment:0:normal 1:ban 2:waitactive" json:"-"`
}

// ComparePassword 检查密码是否匹配
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}

func (user *User) Save() error {
	return model.DB.Save(user).Error
}
