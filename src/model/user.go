package model

import (
	"fmt"
	"framework/database"
	"info"
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
	kUserTableName       = "user"
	kUserId              = "id"
	kUserOpenId          = "open_id"
	kUserName            = "name"
	kUserSex             = "sex"
	kUserType            = "type"
	kUserBigPicutreURL   = "big_pic"
	kUserSmallPicutreURL = "small_pic"
	kUserLastLoginTime   = "login_time"
	kUserRegisterTime    = "reg_time"
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

func (u *userModel) CreateTable() error {
	if database.DatabaseInstance().DoesTableExist(kUserTableName) {
		return nil
	}
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
		%s int(64) unsigned NOT NULL AUTO_INCREMENT,
		%s varchar(128) NOT NULL,
		%s varchar(1024) NOT NULL,
		%s varchar(32) NOT NULL,
		%s int(32) NOT NULL,
		%s varchar(1024) NOT NULL,
		%s varchar(1024) NOT NULL,
		%s int(64) NOT NULL,
		%s int(64) NOT NULL,
		PRIMARY KEY (%s)
	) CHARSET=utf8;`, kUserTableName, kUserId, kUserOpenId,
		kUserName, kUserSex, kUserType, kUserBigPicutreURL, kUserSmallPicutreURL,
		kUserLastLoginTime, kUserRegisterTime, kUserId)
	_, err := database.DatabaseInstance().DB.Exec(sql)
	return err
}

func (u *userModel) Login(accountType int, userInfo *info.UserInfo) error {
	isLogin, err := u.accountHasLogin(accountType, userInfo.UserOpenID)
	if err != nil {
		return err
	}
	if isLogin {
		return u.updateAccountInfo(accountType, userInfo)
	}
	return u.insertAccountInfo(accountType, userInfo)
}

func (u *userModel) accountHasLogin(accountType int, openId string) (bool, error) {
	sql := fmt.Sprintf("select %s from %s where %s = ? and %s = ?",
		kUserId, kUserTableName, kUserType, kUserOpenId)
	stat, err := database.DatabaseInstance().DB.Prepare(sql)
	if err == nil {
		defer stat.Close()
		rows, err := stat.Query(accountType, openId)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				return true, nil
			}
		}
	}
	return false, err
}

func (u *userModel) updateAccountInfo(accountType int, userInfo *info.UserInfo) error {
	sql := fmt.Sprintf("update %s set %s = ?, %s = ?, %s = ?, %s = ?, %s = ? where %s = ? and %s = ?",
		kUserTableName, kUserName, kUserSex, kUserBigPicutreURL,
		kUserSmallPicutreURL, kUserLastLoginTime, kUserType, kUserOpenId)
	currentTime := time.Now().Unix()
	_, err := database.DatabaseInstance().DB.Exec(sql, userInfo.UserName, userInfo.Sex,
		userInfo.BigFigureurl, userInfo.SmallFigureurl, currentTime, accountType, userInfo.UserOpenID)
	if err == nil {
		sql = fmt.Sprintf("select %s from %s where %s = ? and %s = ?", kUserId,
			kUserTableName, kUserType, kUserOpenId)
		rows, err := database.DatabaseInstance().DB.Query(sql, accountType, userInfo.UserOpenID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				rows.Scan(&userInfo.UserID)
				return nil
			}
		}
	}
	return err
}

func (u *userModel) insertAccountInfo(accountType int, userInfo *info.UserInfo) error {
	sql := fmt.Sprintf(
		"insert into %s(%s, %s, %s, %s, %s, %s, %s, %s) values(?, ?, ?, ?, ?, ?, ?, ?)",
		kUserTableName, kUserType, kUserOpenId, kUserName, kUserSex, kUserBigPicutreURL,
		kUserSmallPicutreURL, kUserLastLoginTime, kUserRegisterTime)
	currentTime := time.Now().Unix()
	stat, err := database.DatabaseInstance().DB.Prepare(sql)
	if err == nil {
		defer stat.Close()
		result, err := stat.Exec(accountType, userInfo.UserOpenID, userInfo.UserName, userInfo.Sex,
			userInfo.BigFigureurl, userInfo.SmallFigureurl, currentTime, currentTime)
		if err != nil {
			fmt.Println("result error: ", err)
			return err
		}
		userInfo.UserID, err = result.LastInsertId()
	}
	fmt.Println("error: ", err)
	return err
}

func (u *userModel) GetUserInfoById(userId int64) (*info.UserInfo, error) {
	sql := fmt.Sprintf("select %s, %s, %s from %s where %s = ?", kUserName, kUserSex,
		kUserSmallPicutreURL, kUserTableName, kUserId)
	rows, err := database.DatabaseInstance().DB.Query(sql, userId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var userInfo info.UserInfo
			rows.Scan(&userInfo.UserName, &userInfo.Sex, &userInfo.SmallFigureurl)
			return &userInfo, nil
		}
	}
	return nil, err
}
