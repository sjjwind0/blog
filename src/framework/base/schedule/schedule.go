package schedule

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

const (
	Once   = 0x0
	Repeat = 0x1
	time   = 0x1 << 2
	delay  = 0x1 << 3
	second = 0x1 << 4
	minute = 0x1 << 5
	hour   = 0x1 << 6
	day    = 0x1 << 7
	month  = 0x1 << 8
	year   = 0x1 << 9
)

type ScheduleFunc func(param interface{}, stop *bool)
type schedule struct {
	scheduleFunc func(param interface{})
	param        interface{}
	time         int64
	runType      int64
}

func NewScheduleAtDelay(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | time}
}

func NewScheduleAtTime(f ScheduleFunc, param interface{}, t int64) *schedule {
	return &schedule{f, param, t, Once | delay}
}

func NewScheduleAtSecond(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | second}
}

func NewScheduleAtMinute(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | minute}
}

func NewScheduleAtHour(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | hour}
}

func NewScheduleAtDay(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | day}
}

func NewScheduleAtMonth(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | month}
}

func NewScheduleAtYear(f ScheduleFunc, param interface{}, t int64, runType int) *schedule {
	return &schedule{f, param, t, runType | year}
}

type scheduleMgr struct {
	scheduleList *list.List
}

var scheduleOnce sync.Once
var scheduleMgrInstance *scheduleMgr = nil

func GetScheduleMgrInstance() *scheduleMgr {
	scheduleOnce.Do(func() {
		scheduleMgr = &scheduleMgr{}
	})
	return scheduleMgr
}

func (s *scheduleMgr) RegisterSchedule(work ScheduleWork) {
	if s.scheduleList == nil {
		s.scheduleList = list.New()
	}
	s.scheduleList.PushBack(work)
}

func (s *scheduleMgr) runSchedule(schdule ScheduleWork, t int64, runType int) {
	isRepeat := bool(runType & Repeat)
	var nextTime int64 = 0
	switch runType {
	case delay:
		nextTime = t
	case time:
		currentTime := time.Now().Unix()
		if t <= currentTime {
			panic("error")
		}
		nextTime = (t - currentTime)
	case second:
		if t >= time.Second {
			panic("error")
		}
		nextTime = (time.Now().Unix()/time.Second)*time.Second + t
	case minute:
		if t >= time.Minute {
			panic("error")
		}
		nextTime = (time.Now().Unix()/time.Minute)*time.Minute + t
	case hour:
		if t >= time.Hour {
			panic("hour")
		}
		nextTime = (time.Now().Unix()/time.Hour)*time.Hour + t
		// case month:
		// 	if t >= time.Month {
		// 		panic("hour")
		// 	}
		// 	nextTime = (time.Now().Unix()/time.Month)*time.Month + t
		// case year:
		// 	if t >= time.Ye {
		// 		panic("hour")
		// 	}
		// 	nextTime = (time.Now().Unix()/time.Hour)*time.Hour + t
	}
	timer := time.NewTimer(t)
	go func() {
		<-timer.C
		schdule()
	}()
}

func (s *scheduleMgr) Run() {
	for iter := s.scheduleList.Front(); iter != nil; iter = iter.Next() {
		schedule := iter.Value.(ScheduleWork)
		s.runSchedule(schdule)
	}
}
