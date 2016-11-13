package model

import (
	"container/list"
	"fmt"
	"framework/database"
	"info"
	"sync"
	"time"
)

type pluginModel struct {
}

const (
	kPluginTableName    = "plugin"
	kPluginId           = "id"
	kPluginUUID         = "uuid"
	kPluginType         = "type"
	kPluginName         = "name"
	kPluginVersion      = "version"
	kPluginTime         = "time"
	kPluginVisitCount   = "visit"
	kPluginPraiseCount  = "praise"
	kPluginDissentCount = "dissent"
)

var pluginModelInstance *pluginModel = nil

var pluginOnce sync.Once

func SharePluginModel() *pluginModel {
	pluginOnce.Do(func() {
		pluginModelInstance = &pluginModel{}
	})
	return pluginModelInstance
}

func (c *pluginModel) CreateTable() error {
	if database.DatabaseInstance().DoesTableExist(kPluginTableName) {
		return nil
	}
	fmt.Println("Hello World")
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
		%s int(32) unsigned NOT NULL AUTO_INCREMENT,
		%s varchar(128) NOT NULL,
		%s varchar(256) NOT NULL,
		%s varchar(256) NOT NULL,
		%s varchar(128) DEFAULT '1.0.0',
		%s int(64) NOT NULL,
		%s int(32) DEFAULT '0',
		%s int(32) DEFAULT '0',
		%s int(32) DEFAULT '0',
		PRIMARY KEY (%s)
	) CHARSET=utf8;`, kPluginTableName, kPluginId,
		kPluginUUID, kPluginName, kPluginType, kPluginVersion, kPluginTime, kPluginVisitCount,
		kPluginPraiseCount, kPluginDissentCount, kPluginId)
	_, err := database.DatabaseInstance().DB.Exec(sql)
	return err
}

func (b *pluginModel) InsertPlugin(uuid string, title string, pluginType int,
	pluginVersion string) (int, error) {
	currentTime := time.Now().Unix()
	sql := fmt.Sprintf("insert into %s(%s, %s, %s, %s, %s) values(?, ?, ?, ?, ?)",
		kPluginTableName, kPluginUUID, kPluginName, kPluginType, kPluginVersion, kPluginTime)
	fmt.Println(sql)
	stat, err := database.DatabaseInstance().DB.Prepare(sql)
	if err == nil {
		defer stat.Close()
		result, err := stat.Exec(uuid, title, pluginType, pluginVersion, currentTime)
		insertId, _ := result.LastInsertId()
		return int(insertId), err
	}
	return -1, err
}

func (b *pluginModel) UpdatePlugin(uuid string, title string, pluginType int,
	pluginVersion string) (int, error) {
	currentTime := time.Now().Unix()
	sql := fmt.Sprintf("update %s set %s = ?, %s = ?, %s = ?, %s = ? where %s = ?",
		kPluginTableName, kPluginName, kPluginType, kPluginVersion, kPluginTime, kPluginUUID)
	result, err := database.DatabaseInstance().DB.Exec(sql, title, pluginType, pluginVersion, currentTime, uuid)
	updateId, _ := result.RowsAffected()
	fmt.Println("updateId: ", updateId)
	return int(updateId), err
}

func (b *pluginModel) PluginIsExistByUUID(uuid string) (bool, error) {
	sql := fmt.Sprintf("select * from %s where %s = ?", kPluginTableName, kPluginUUID)
	fmt.Println("sql: ", sql)
	rows, err := database.DatabaseInstance().DB.Query(sql, uuid)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			return true, nil
		}
	}
	return false, err
}

func (b *pluginModel) PluginIsExistByPluginID(pluginId int) (bool, error) {
	sql := fmt.Sprintf("select * from %s where %s = ?", kPluginTableName, kPluginId)
	rows, err := database.DatabaseInstance().DB.Query(sql, pluginId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			return true, nil
		}
	}
	return false, err
}

func (b *pluginModel) FetchAllPlugin() (*list.List, error) {
	sql := fmt.Sprintf("select * from %s order by %s desc", kPluginTableName, kPluginId)
	rows, err := database.DatabaseInstance().DB.Query(sql)
	if err == nil {
		defer rows.Close()
		var pluginList *list.List = list.New()
		for rows.Next() {
			var plugin info.PluginInfo
			err = rows.Scan(&plugin.PluginID, &plugin.PluginUUID, &plugin.PluginName,
				&plugin.PluginType, &plugin.PluginVersion, &plugin.PluginTime, &plugin.PluginVisitCount,
				&plugin.PluginPraiseCount, &plugin.PluginDissentCount)
			if err == nil {
				pluginList.PushBack(plugin)
			}
		}
		return pluginList, err
	}
	fmt.Println(err)
	return nil, err
}

func (b *pluginModel) FetchPluginByPluginID(pluginID int) (*info.PluginInfo, error) {
	sql := fmt.Sprintf("select * from %s where %s = ?", kPluginTableName, kPluginId)
	rows, err := database.DatabaseInstance().DB.Query(sql, pluginID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var plugin info.PluginInfo
			err = rows.Scan(&plugin.PluginID, &plugin.PluginUUID, &plugin.PluginName,
				&plugin.PluginType, &plugin.PluginVersion, &plugin.PluginTime, &plugin.PluginVisitCount,
				&plugin.PluginPraiseCount, &plugin.PluginDissentCount)
			if err == nil {
				return &plugin, nil
			}
			break
		}
	}
	return nil, err
}

func (b *pluginModel) GetPluginUUIDByPluginID(pluginID int) (string, error) {
	sql := fmt.Sprintf("select %s from %s where %s = ?", kPluginUUID, kPluginTableName, kPluginId)
	rows, err := database.DatabaseInstance().DB.Query(sql, pluginID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var uuid string
			err = rows.Scan(&uuid)
			if err == nil {
				return uuid, nil
			}
			break
		}
	}
	return "", err
}

func (b *pluginModel) FetchPluginByUUID(uuid string) (*info.PluginInfo, error) {
	sql := fmt.Sprintf("select* from %s where %s = ?", kPluginTableName, kPluginUUID)
	rows, err := database.DatabaseInstance().DB.Query(sql, uuid)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var plugin info.PluginInfo
			err = rows.Scan(&plugin.PluginID, &plugin.PluginUUID, &plugin.PluginName,
				&plugin.PluginType, &plugin.PluginVersion, &plugin.PluginTime, &plugin.PluginVisitCount,
				&plugin.PluginPraiseCount, &plugin.PluginDissentCount)
			if err == nil {
				return &plugin, nil
			}
			break
		}
	}
	return nil, err
}

func (b *pluginModel) FetchAllType() ([]int, error) {
	sql := fmt.Sprintf("select %s from %s distinct", kPluginType, kPluginTableName)
	rows, err := database.DatabaseInstance().DB.Query(sql)
	if err == nil {
		defer rows.Close()
		var typeList []int
		for rows.Next() {
			var pluginType int
			err = rows.Scan(&pluginType)
			if err == nil {
				typeList = append(typeList, pluginType)
			}
		}
		return typeList, nil
	}
	return nil, err
}

func (b *pluginModel) FetchAllPluginBySortType(pluginType int) (*list.List, error) {
	sql := fmt.Sprintf("select * from %s where %s = ? order by %s desc",
		kPluginTableName, kPluginType, kPluginId)
	rows, err := database.DatabaseInstance().DB.Query(sql, pluginType)
	if err == nil {
		defer rows.Close()
		var pluginList *list.List = list.New()
		for rows.Next() {
			var plugin info.PluginInfo
			err = rows.Scan(&plugin.PluginID, &plugin.PluginUUID, &plugin.PluginName,
				&plugin.PluginType, &plugin.PluginVersion, &plugin.PluginTime, &plugin.PluginVisitCount,
				&plugin.PluginPraiseCount, &plugin.PluginDissentCount)
			if err == nil {
				pluginList.PushBack(plugin)
			}
		}
		return pluginList, err
	}
	return nil, err
}

func (b *pluginModel) FetchAllPluginByTime(beginTime int64, endTime int64) (*list.List, error) {
	sql := fmt.Sprintf("select * from %s where %s >= ? and %s <= ? order by %s desc", kPluginTableName, kPluginTime, kPluginTime, kPluginId)
	rows, err := database.DatabaseInstance().DB.Query(sql, beginTime, endTime)
	if err == nil {
		defer rows.Close()
		var pluginList *list.List = list.New()
		for rows.Next() {
			var plugin info.PluginInfo
			err = rows.Scan(&plugin.PluginID, &plugin.PluginUUID, &plugin.PluginName,
				&plugin.PluginType, &plugin.PluginVersion, &plugin.PluginTime, &plugin.PluginVisitCount,
				&plugin.PluginPraiseCount, &plugin.PluginDissentCount)
			if err == nil {
				pluginList.PushBack(plugin)
			}
		}
		return pluginList, err
	}
	return nil, err
}

func (b *pluginModel) AddVisitCount(pluginId int) error {
	sql := fmt.Sprintf("update %s set visit = visit + 1 where %s = ?", kPluginTableName, kPluginId)
	_, err := database.DatabaseInstance().DB.Exec(sql, pluginId)
	return err
}

func (b *pluginModel) DeletePlugin(pluginId int) error {
	sql := fmt.Sprintf("delete from %s where %s = ?", kPluginTableName, kPluginId)
	_, err := database.DatabaseInstance().DB.Exec(sql, pluginId)
	return err
}
