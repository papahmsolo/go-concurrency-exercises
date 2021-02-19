//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	const checkInterval = 1
	tick := time.NewTicker(checkInterval * time.Second)
	defer tick.Stop()

	finish := make(chan struct{})

	go func() {
		process()
		finish <- struct{}{}
	}()

	for {
		select {
		case <-tick.C:
			if !u.IsPremium && atomic.AddInt64(&u.TimeUsed, int64(checkInterval)) >= 10 {
				return false
			}
		case <-finish:
			return true
		}
	}
}

func main() {
	RunMockServer()
}
