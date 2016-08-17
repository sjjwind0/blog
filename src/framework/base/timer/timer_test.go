package timer

import (
	"fmt"
	"testing"
	"time"
)

func Test_OneShotTimer(t *testing.T) {
	var testChan chan bool = make(chan bool)
	t1 := NewOneShotTimer()
	fmt.Println("startTimer: ", time.Now().Unix())
	t1.Start(5*time.Second, func() {
		fmt.Println("onTimer: ", time.Now().Unix())
		<-testChan
	})
	testChan <- true
}

func Test_RepeatTimer(t *testing.T) {
	var testChan chan bool = make(chan bool)
	t1 := NewRepeatingTimer()
	fmt.Println("startTimer: ", time.Now().Unix())
	var index int = 0
	t1.Start(time.Second, func() {
		fmt.Println("onTimer: ", time.Now().Unix())
		if index == 5 {
			<-testChan
		}
		index++
	})
	testChan <- true
}
