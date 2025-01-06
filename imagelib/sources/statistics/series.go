package statistics

import (
	"fmt"
	"github.com/pavlo67/common/common/joiner"
	"time"
)

type Series struct {
	Path string
	FPS  float64

	Cnt          int
	MsPerFrame   float64
	MsPerMegaPix float64
	TriesPerSucc float64
	SuccessRatio float64
	MissMin      float64
	MissMinP     float64
	MissMax      float64
	MissMaxP     float64
	MissFin      float64
	MissFinP     float64
	MsTotal      float64
	Times        [][]float64 `json:",omitempty"`

	Actor    joiner.InterfaceKey
	Scenario string
	DPM      float64
	RGBA     bool `json:",omitempty"`
}

type Accumulator struct {
	Actor joiner.InterfaceKey
	FPS   float64

	StartedAt       time.Time
	Cnt             int
	MegaPixelsCnt   float64
	TriesSuccCnt    int
	TriesCnt        int
	SuccessCnt      int
	MissMin         float64
	MissMinPercents float64
	MissMax         float64
	MissMaxPercents float64
	MissFin         float64
	MissFinPercents float64
	TimeSuccTotal   time.Duration
	TimeTotal       time.Duration
	Times           [][]float64

	Path     string
	Scenario string
	DPM      float64
	RGBA     bool
}

func (acc *Accumulator) Add(megaPixels float64, tries int, success bool, deviation, deviationPercents float64) {
	if acc == nil {
		return
	}
	acc.Cnt++

	timeTotalNew := time.Now().Sub(acc.StartedAt)
	timeNew := timeTotalNew - acc.TimeTotal
	acc.Times = append(acc.Times, []float64{float64(tries), float64(timeNew) / float64(time.Millisecond)})
	acc.TimeTotal = timeTotalNew

	acc.MegaPixelsCnt += megaPixels
	acc.TriesCnt += tries
	if success {
		acc.TimeSuccTotal += timeNew
		acc.TriesSuccCnt += tries
		acc.SuccessCnt++
	}
	if acc.MissMin == 0 {
		acc.MissMin, acc.MissMinPercents = deviation, deviationPercents
		acc.MissMax, acc.MissMaxPercents = deviation, deviationPercents
	} else {
		if deviation < acc.MissMin {
			acc.MissMin = deviation
		} else if deviation > acc.MissMax {
			acc.MissMax = deviation
		}
		if deviationPercents < acc.MissMinPercents {
			acc.MissMinPercents = deviationPercents
		} else if deviationPercents > acc.MissMaxPercents {
			acc.MissMaxPercents = deviationPercents
		}
	}
	acc.MissFin, acc.MissFinPercents = deviation, deviationPercents

}

func (acc Accumulator) Series() (*Series, error) {
	if acc.Cnt <= 0 {
		return nil, fmt.Errorf("empty Accumulator: %+v / on Accumulator.Series()", acc)
	}

	msTotal := float64(acc.TimeTotal) / float64(time.Millisecond)
	tValues := []float64{float64(acc.TriesCnt) / float64(acc.Cnt), msTotal / float64(acc.Cnt)}

	var triesPerSuccess float64
	if acc.SuccessCnt > 0 {
		triesPerSuccess = float64(acc.TriesSuccCnt) / float64(acc.SuccessCnt)
		tValues = append(tValues, float64(acc.TimeSuccTotal)/(float64(acc.SuccessCnt)*float64(time.Millisecond)))
	}

	acc.Times = append(acc.Times, tValues)

	return &Series{
		Actor:        acc.Actor,
		FPS:          acc.FPS,
		Cnt:          acc.Cnt,
		MsTotal:      msTotal,
		MsPerFrame:   float64(acc.TimeTotal) / float64(time.Millisecond*time.Duration(acc.Cnt)),
		MsPerMegaPix: float64(acc.TimeTotal) / (float64(time.Millisecond) * acc.MegaPixelsCnt),
		Times:        acc.Times,
		TriesPerSucc: triesPerSuccess,
		SuccessRatio: float64(acc.SuccessCnt) / float64(acc.Cnt),
		MissMin:      acc.MissMin,
		MissMinP:     acc.MissMinPercents,
		MissMax:      acc.MissMax,
		MissMaxP:     acc.MissMaxPercents,
		MissFin:      acc.MissFin,
		MissFinP:     acc.MissFinPercents,

		Path:     acc.Path,
		Scenario: acc.Scenario,
		DPM:      acc.DPM,
		RGBA:     acc.RGBA,
	}, nil
}
