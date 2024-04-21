package converter

import (
	"fmt"
	"math"
)

type MapCords struct {
	LongDir                   string
	LongDeg, LongMin, LongSec float64
	LatDir                    string
	LatDeg, LatMin, LatSec    float64
	CordsReqString            string
}

func ConvertCordsToMapCords(longitude, latitude float64) *MapCords {
	longDir := "E"
	if longitude < 0 {
		longDir = "W"
	}
	latDir := "N"
	if latitude < 0 {
		latDir = "S"
	}

	absLong := math.Abs(longitude)
	absLat := math.Abs(latitude)

	longDeg := math.Floor(absLong)
	longMin := math.Floor((absLong - longDeg) * 60)
	longSec := (absLong - longDeg - longMin/60) * 3600

	latDeg := math.Floor(absLat)
	latMin := math.Floor((absLat - latDeg) * 60)
	latSec := (absLat - latDeg - latMin/60) * 3600

	cordsReqString := fmt.Sprintf("lat:%v%s%.3f;lon:%v%s%.3f",
		latDeg,
		latDir,
		(absLat-latDeg)*60,
		longDeg,
		longDir,
		(absLong-longDeg)*60,
	)

	return &MapCords{
		longDir,
		longDeg,
		longMin,
		longSec,
		latDir,
		latDeg,
		latMin,
		latSec,
		cordsReqString,
	}
}
