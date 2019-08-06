package aggregate

import (
	"fmt"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	go func() {
		Lock("1", "Test")
		fmt.Println("========locked")
		time.Sleep(2 * time.Second)
		fmt.Println("func1 end...", time.Now())
		UnLock("1", "Test")
	}()
	go func() {
		Lock("1", "Test")
		fmt.Println("========locked")
		time.Sleep(2 * time.Second)
		fmt.Println("func2 end...", time.Now())
		UnLock("1", "Test")
	}()

	select {}
}
