package main

import (
	"errors"
	"time"
)

type LineSpeedAggr map[uint]LineSpeedGroup

func (lsa *LineSpeedAggr) Metrics() (speedMetrics []LineSpeedMetric) {
	for _, lsg := range *lsa {
		speedMetrics = append(speedMetrics, lsg.Metrics())
	}

	return
}

func (lsa *LineSpeedAggr) AddLs(ls LineSpeed) {
	lsg := (*lsa)[ls.LineId]

	lsg.LineId = ls.LineId
	_ = lsg.AddLs(ls)

	(*lsa)[ls.LineId] = lsg
}

type LineSpeedGroup struct {
	LineId uint
	Items  []LineSpeed
}

func (lsg *LineSpeedGroup) AddLs(ls LineSpeed) error {
	if ls.LineId != lsg.LineId {
		return errors.New("LineID mismatch")
	}

	if lsg.Items == nil {
		lsg.Items = []LineSpeed{}
	}

	lsg.Items = append(lsg.Items, ls)

	return nil
}

func (lsg *LineSpeedGroup) Metrics() (m LineSpeedMetric) {
	m.LineId = lsg.LineId
	sum, max, min := 0.0, 0.0, 0.0
	c := 0

	if len(lsg.Items) == 0 {
		return
	}

	if len(lsg.Items) > 0 && !lsg.Items[0].IsOld() {
		min = lsg.Items[0].Speed
	}

	minSet := false

	for _, ls := range lsg.Items {
		if ls.IsOld() {
			continue
		}

		c++
		speed := ls.Speed
		sum += speed

		if !minSet {
			minSet = true
			min = speed
		}

		if speed > max {
			max = speed
		}

		if speed < min {
			min = speed
		}
	}

	avg := 0.0

	if sum > 0 {
		avg = sum / float64(c)
	}

	m.Metrics = Metrics{
		Avg: avg,
		Max: max,
		Min: min,
	}

	return
}

type LineSpeed struct {
	LineId    uint    `json:"line_id"`
	Speed     float64 `json:"speed"`
	Timestamp uint    `json:"timestamp"`
}

func (ls *LineSpeed) Time() time.Time {
	return time.Unix(0, int64(ls.Timestamp)*int64(time.Millisecond))
}

func (ls *LineSpeed) IsOld() bool {
	return time.Now().Sub(ls.Time()).Hours() >= 1
}

type LineSpeedMetric struct {
	LineId  uint    `json:"line_id"`
	Metrics Metrics `json:"metrics"`
}

type Metrics struct {
	Avg float64 `json:"avg"`
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}
