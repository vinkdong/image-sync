package log

import (
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	Info("this is info")
}

func TestError(t *testing.T) {
	Error("this is error")
}

func TestSuccess(t *testing.T)  {
	Success("this is success")
}

func TestSuccessf(t *testing.T)  {
	Successf("this is success with args type: %s", "string")
}

func TestLock(t*testing.T) {
	Info("this is first log next lock log")
	Lock()
	go func() {
		for {
			Info("this is sub task")
			time.Sleep(time.Second *1)
		}
	}()
	time.Sleep(time.Second *1)
	Unlock()
	Info("this is last log")
	Lock()
	time.Sleep(time.Second*3)
	Unlock()
	time.Sleep(time.Second*5)
}