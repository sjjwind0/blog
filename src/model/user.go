package model

import (
	"fmt"
	"framework/database"
	"sync"
)

const (
	kUserTableName = "user"
	kUserId        = "user_id"
	kUserName      = "name"
	kUserFrom      = "from"
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
		%s varchar(256) NOT NULL,
		%s int(32) NOT NULL,
		PRIMARY KEY (%s)
	) ENGINE=MyISAM DEFAULT CHARSET=utf8;`, kUserTableName, kUserId,
		kUserName, kUserFrom)
	_, err := database.DatabaseInstance().DB.Exec(sql)
	return err
}

func (u *userModel) Login() error {
	return nil
}
