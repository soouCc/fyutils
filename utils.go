package fyutils

import (
	"fmt"
	"github.com/google/gops/agent"
	"github.com/liu-junyong/go-logger/logger"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
	"math/rand"
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

func FindIntIn(arr []int32, cd int32) bool {
	for i := 0; i < len(arr); i++ {
		if cd == arr[i] {
			return true
		}
	}

	return false
}

func RemoveCards(cds []int32, value int32, count int32) []int32 {
	res := make([]int32, 0)
	cnt := int32(0)
	for _, v := range cds {
		if v != value {
			res = append(res, v)
		} else {
			if count == -1 {
				continue
			}
			if cnt < count {
				cnt++
			} else {
				res = append(res, v)
			}
		}
	}
	return res
}

//获取一个次数范围外的随机数
func RandExp(min,max int,p int) int{

	n := rand.Intn(max-min+1) + min
	if n==p{
		return RandExp(min,max ,p )
	}
	return n
}

func RemoveSlice(source []int32, target []int32) []int32 {
	for i := 0; i < len(target); i++ {
		for j := len(source) - 1; j >= 0; j-- {
			if target[i] == source[j] {
				source = append(source[:j], source[j+1:]...)
				break
			}
		}
	}
	return source
}
