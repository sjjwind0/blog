package database

import (
	"container/list"
	"fmt"
	"sync"
)

type databaseRunner struct {
	modeList *list.List
}

var databaseRunnerInstance *databaseRunner = nil

var databaseRunnerOnce sync.Once

func ShareDatabaseRunner() *databaseRunner {
	databaseRunnerOnce.Do(func() {
		databaseRunnerInstance = &databaseRunner{nil}
	})
	return databaseRunnerInstance
}

func (d *databaseRunner) RegisterModel(model DatabaseInterface) {
	if d.modeList == nil {
		d.modeList = list.New()
	}
	d.modeList.PushBack(model)
}

func (d *databaseRunner) Start() {
	if d.modeList != nil {
		for model := d.modeList.Front(); model != nil; model = model.Next() {
			databaseInterface := model.Value.(DatabaseInterface)
			err := databaseInterface.CreateTable()
			if err != nil {
				fmt.Println("CreateTable error: ", err)
			}
		}
	}
}
