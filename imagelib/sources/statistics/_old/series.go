package statistics

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/serialization"
)

type Series struct {
	Key string
	FPS float64

	Cnt          int
	MsPerFrame   float64
	MsPerMegaPix float64
	TriesPerSucc float64
	SuccessRatio float64
	MissTot      float64
	MissTotP     float64
	MissMin      float64
	MissMinP     float64
	MissMax      float64
	MissMaxP     float64
	Details      []string
	MsTotal      float64
	TriesInfo    [][]string `json:",omitempty"`

	Actor    joiner.InterfaceKey
	Scenario string
	DPM      float64
	RGBA     bool `json:",omitempty"`
}

type Accumulator struct {
	Key             string
	StartedAt       time.Time
	Cnt             int
	MegaPixelsCnt   float64
	TriesSuccCnt    int
	TriesCnt        int
	SuccessCnt      int
	MissTot         float64
	MissTotPercents float64
	MissMin         float64
	MissMinPercents float64
	MissMax         float64
	MissMaxPercents float64
	Details         []string
	TimeSuccTotal   time.Duration
	TimeTotal       time.Duration
	TriesInfo       [][]string
	Actor           joiner.InterfaceKey
	Scenario        string
	DPM             float64
	FPS             float64
	RGBA            bool
}

func (acc *Accumulator) Add(megaPixels float64, tries int, success bool, miss, missPercents float64) {
	if acc == nil {
		return
	}
	acc.Cnt++

	timeTotalNew := time.Now().Sub(acc.StartedAt)
	timeNew := timeTotalNew - acc.TimeTotal
	acc.TimeTotal = timeTotalNew
	ms := fmt.Sprintf("%f", float64(timeNew)/float64(time.Millisecond))

	acc.MegaPixelsCnt += megaPixels
	acc.TriesCnt += tries
	var msSuccess string
	if success {
		acc.TimeSuccTotal += timeNew
		acc.TriesSuccCnt += tries
		acc.SuccessCnt++
		msSuccess = ms
	}

	acc.TriesInfo = append(acc.TriesInfo, []string{strconv.Itoa(tries), ms, msSuccess})

	if acc.MissMin == 0 {
		acc.MissMin, acc.MissMinPercents = miss, missPercents
		acc.MissMax, acc.MissMaxPercents = miss, missPercents
	} else {
		if miss < acc.MissMin {
			acc.MissMin = miss
		} else if miss > acc.MissMax {
			acc.MissMax = miss
		}
		if missPercents < acc.MissMinPercents {
			acc.MissMinPercents = missPercents
		} else if missPercents > acc.MissMaxPercents {
			acc.MissMaxPercents = missPercents
		}
	}
	acc.MissTot, acc.MissTotPercents = miss, missPercents

}

func (acc Accumulator) Series(details []string) (*Series, error) {
	if acc.Cnt <= 0 {
		return nil, fmt.Errorf("empty Accumulator: %+v / on Accumulator.Series()", acc)
	}

	msTotal := float64(acc.TimeTotal) / float64(time.Millisecond)
	triesInfo := []string{fmt.Sprintf("%f", float64(acc.TriesCnt)/float64(acc.Cnt)), fmt.Sprintf("%f", msTotal/float64(acc.Cnt))}

	var triesPerSuccess float64
	if acc.SuccessCnt > 0 {
		triesPerSuccess = float64(acc.TriesSuccCnt) / float64(acc.SuccessCnt)
		triesInfo = append(triesInfo, fmt.Sprintf("%f", float64(acc.TimeSuccTotal)/(float64(acc.SuccessCnt)*float64(time.Millisecond))))
	}

	acc.TriesInfo = append(acc.TriesInfo, triesInfo)

	return &Series{
		Actor:        acc.Actor,
		FPS:          acc.FPS,
		Cnt:          acc.Cnt,
		MsTotal:      msTotal,
		MsPerFrame:   float64(acc.TimeTotal) / float64(time.Millisecond*time.Duration(acc.Cnt)),
		MsPerMegaPix: float64(acc.TimeTotal) / (float64(time.Millisecond) * acc.MegaPixelsCnt),
		TriesInfo:    acc.TriesInfo,
		TriesPerSucc: triesPerSuccess,
		SuccessRatio: float64(acc.SuccessCnt) / float64(acc.Cnt),
		MissMin:      acc.MissMin,
		MissMinP:     acc.MissMinPercents,
		MissMax:      acc.MissMax,
		MissMaxP:     acc.MissMaxPercents,
		MissTot:      acc.MissTot,
		MissTotP:     acc.MissTotPercents,
		Details:      details,
		Key:          acc.Key,
		Scenario:     acc.Scenario,
		DPM:          acc.DPM,
		RGBA:         acc.RGBA,
	}, nil
}

const onSave = "on statistics.Series.Save()"

func (series Series) Save(path string) error {
	var timesStr string
	for _, tryInfo := range series.TriesInfo {
		timesStr += strings.Join(tryInfo, "\t") + "\n"
	}
	if err := os.WriteFile(filepath.Join(path, "tries_info_"+time.Now().Format(time.RFC3339)[:19]+".xls"), []byte(timesStr), 0644); err != nil {
		return errors.Wrap(err, onSave)
	}

	series.TriesInfo = nil

	if err := serialization.SavePart(series, serialization.MarshalerJSON, filepath.Join(path, StatisticsFilename)); err != nil {
		return errors.Wrap(err, onSave)
	}

	return nil

}
