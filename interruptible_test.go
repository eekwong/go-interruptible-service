package interruptible

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

type testService struct {
	Service
	started          bool
	cleanup          bool
	goRoutineStopped bool
	errInRun         error
}

func (ts *testService) Run() (Stop, error) {
	ts.started = true
	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		<-done
		ts.goRoutineStopped = true
		close(stopped)
	}()
	return func() error {
		close(done)
		<-stopped
		ts.cleanup = true
		return nil
	}, ts.errInRun
}

func TestNormal(t *testing.T) {
	ts := &testService{}
	go func() {
		time.Sleep(time.Second * 1)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()
	Run(ts)

	if !ts.started {
		t.Fatal("ts.started should be true")
	}
	if !ts.goRoutineStopped {
		t.Fatal("ts.goRoutineStopped should be true")
	}
	if !ts.cleanup {
		t.Fatal("ts.cleanup should be true")
	}
}

func TestRunError(t *testing.T) {
	ts := &testService{
		errInRun: errors.New(""),
	}
	Run(ts)

	if !ts.started {
		t.Fatal("ts.started should be true")
	}
	if !ts.goRoutineStopped {
		t.Fatal("ts.goRoutineStopped should be true")
	}
	if !ts.cleanup {
		t.Fatal("ts.cleanup should be true")
	}
}
