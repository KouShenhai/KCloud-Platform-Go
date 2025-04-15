package core

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	now := time.Now()
	t1 := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	t2 := now.AddDate(0, 0, -7)
	fmt.Println(t1.Unix(), t2.Unix())
}
