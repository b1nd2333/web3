package common

import (
	"fmt"
	"time"
)

// ChineseTimeNow 获取当前中国时间
func ChineseTimeNow() time.Time {
	// 加载中国时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("加载时区出错: %v\n", err)
		return time.Time{}
	}

	// 获取当前中国时间
	now := time.Now().In(loc)
	return now
}

// Int64Time2String 时间戳转字符串
func Int64Time2String(nextTime int64) time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("加载时区出错: %v\n", err)
		return time.Time{}
	}

	var timestampSec, timestampNano int64

	// 判断时间戳单位
	switch {
	// 纳秒时间戳
	case nextTime > 1e15:
		timestampSec = nextTime / 1e9
		timestampNano = nextTime % 1e9
	// 微秒时间戳
	case nextTime > 1e12:
		timestampSec = nextTime / 1e6
		timestampNano = (nextTime % 1e6) * 1e3
	// 毫秒时间戳
	case nextTime > 1e9:
		timestampSec = nextTime / 1000
		timestampNano = (nextTime % 1000) * 1e6
	// 秒时间戳
	default:
		timestampSec = nextTime
		timestampNano = 0
	}

	targetTime := time.Unix(timestampSec, timestampNano).In(loc)
	return targetTime
}
