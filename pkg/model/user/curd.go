package user

import (
	"chatgpt-web/pkg/model"
)

// GetByName 根据名称获取用户
func GetByName(name string) (user *User, err error) {
	user = &User{}
	err = model.DB.Where("name = ?", name).First(user).Error
	return
}

// CreateUser 创建用户
func CreateUser(name, password, realname string, expire int64) (user *User, err error) {
	user = &User{}
	user.Name = name
	user.Password = password
	user.Realname = realname
	user.ExpireTimestamp = expire
	result := model.DB.Create(user)
	err = result.Error
	return
}

func InitUser(name, password, realname string, expire int64) (user *User, err error) {
	user = &User{}
	user.Name = name
	user.Password = password
	user.Realname = realname
	user.ExpireTimestamp = expire
	user.Stat = 2 //set waitactive
	result := model.DB.Create(user)
	err = result.Error
	return
}

func DeleteUser(name string) {
	var user *User
	err := model.DB.Where("name = ?", name).First(&user).Error
	if err != nil {
		return
	}
	if user.ID > 0 {
		model.DB.Delete(&user)
	}
	return
}
