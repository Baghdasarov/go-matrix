package main

import (
	"testing"
	"time"
)

func TestLineSpeed_IsOld(t *testing.T) {
	now := time.Now()

	ls := LineSpeed{
		LineId:    0,
		Speed:     0,
		Timestamp: msTimeStamp(),
	}

	if ls.IsOld() {
		t.Fatal("LineSpeed with current timestamp is old")
	}

	ls.Timestamp = uint(now.Add(-1*time.Minute).UnixNano() / int64(time.Millisecond))

	if ls.IsOld() {
		t.Fatal("should not be old")
	}

	ls.Timestamp = uint(now.Add(-1*time.Hour).UnixNano() / int64(time.Millisecond))

	if !ls.IsOld() {
		t.Fatal("LineSpeed with old timestamp is not old")
	}
}

func TestLineSpeedGroup_Items(t *testing.T) {
	lsg := LineSpeedGroup{LineId: 1}

	if lsg.Items != nil {
		t.Fatal("Newly created LineSpeedGroup must not have any items")
	}

	ls := LineSpeed{
		LineId:    1,
		Speed:     1,
		Timestamp: msTimeStamp(),
	}

	_ = lsg.AddLs(ls)

	if len(lsg.Items) == 0 {
		t.Fatal("LineSpeedGroup must not be empty after adding an item")
	}
}

func TestLineSpeedGroup_AddLs(t *testing.T) {
	lsg := LineSpeedGroup{LineId: 1}

	ls := LineSpeed{
		LineId:    1,
		Speed:     1,
		Timestamp: msTimeStamp(),
	}

	e := lsg.AddLs(ls)

	if e != nil {
		t.Fatal("must not return error")
	}

	ls.LineId = 2

	e = lsg.AddLs(ls)

	if e == nil {
		t.Fatal("must return error")
	}
}

func msTimeStamp() uint {
	return uint(time.Now().UnixNano() / int64(time.Millisecond))
}
