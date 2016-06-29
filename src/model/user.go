package model

import (
	"fmt"
	"framework/database"
	"sync"
	"time"
)

// 登陆过期时间为7天
const kLoginRepireTime = time.Hour * 24 * 7

const (
	kAccounQQ     = 0x1
	kAccountWeibo = 0x2
)

const (
	kUserTableName     = "user"
	kUserId            = "id"
	kUserOpenId        = "open_id"
	kUserName          = "name"
	kUserType          = "type"
	kUserPicutreURL    = "pic"
	kUserLastLoginTime = "login_time"
	kUserRegisterTime  = "reg_time"
)

type userModel struct {
}

var userModelInstance *userModel = nil

var userOnce sync.Once

func ShareUserModel() *userModel {
	userOnce.Do(func() {
		userModelInstance = &userModel{}
	})
	return userModelInstance
}

func (c *userModel) CreateTable() error {
	if database.DatabaseInstance().DoesTableExist(kUserTableName) {
		return nil
	}
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
		%s int(32) unsigned NOT NULL AUTO_INCREMENT,
		%s int(32) unsigned NOT NULL,
		%s varchar(256) NOT NULL,
		%s int(32) NOT NULL,
		%s varchar(1024) NOT NULL,
		%s int(64) NOT NULL,
		%s int(64) NOT NULL,
		PRIMARY KEY (%s)
	) ENGINE=MyISAM DEFAULT CHARSET=utf8;`, kUserTableName, kUserId, kUserOpenId,
		kUserName, kUserType, kUserPicutreURL, kUserLastLoginTime, kUserRegisterTime)
	_, err := database.DatabaseInstance().DB.Exec(sql)
	return err
}

func (u *userModel) Login() error {
	return nil
}
