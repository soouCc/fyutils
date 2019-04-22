package fyutils

import (
	"fmt"
	"github.com/google/gops/agent"
	"github.com/liu-junyong/go-logger/logger"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
)

func CacheError() {
	if err := recover(); err != nil {
		logger.Error(fmt.Sprintf("CacheError, err:%v ", err))
		logger.Error(string(debug.Stack()))
	}
}

//gops监听
func StartGops(port string) {
	fn := func() {

		defer CacheError()
		addr := fmt.Sprintf(":%s", port)
		logger.Debug("----------启动 gops---------port:", addr)

		if err := agent.Listen(agent.Options{
			ShutdownCleanup: true,
			Addr:            addr,
		}); err != nil {
			logger.Error(err.Error())
		}
	}
	go fn()
}

//通过两个经纬度计算距离   返回值的单位为米
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := float64(6371000) // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

	return dist * radius
}

func String2SliceFloat64(str string) []float64 {
	s:=strings.Split(str,",")
	ssf:=make([]float64,len(s))
	for i:=0;i<len(s);i++{
		ssf[i] = Json2float64(s[i])
	}
	return ssf
}

func Json2float64(hh interface{}) float64 {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()

	if hh == nil {
		return 0
	}

	heifan := 0
	switch hh.(type) {
	case float64:
		heifan = int(hh.(float64))
		return float64(heifan)
	case float32:
		heifan = int(hh.(float64))
		return float64(heifan)
	case int32:
		heifan = int(hh.(int32))
		return float64(heifan)
	case int64:
		heifan = int(hh.(int64))
		return float64(heifan)
	case string:
		str := hh.(string)
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return float
	}
	return float64(hh.(int))
}

