package goFunc

import "github.com/sirupsen/logrus"

func Go(f func() error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("PANIC: in goroutine %v", r)

			}
		}()
		err := f()
		if err != nil {
			logrus.Errorf("Error in goroutine: %v", err)
			return
		}
	}()
}
