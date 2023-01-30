package time

import (
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"time"
)

// https://goframe.org/os/gtime/func

// CN,TW,JP,KR,国服,台服,日服,韩服
func GetNowTimeToZone(zone string) time.Time {
	zone = countryToFile(zone)
	var cstSh, _ = time.LoadLocation(zone)
	return time.Now().In(cstSh)
}

// CN,TW,JP,KR,国服,台服,日服,韩服
func GetTimeAtUnixToZone(zone string, unixTime int64) time.Time {
	zone = countryToFile(zone)
	var cstSh, _ = time.LoadLocation(zone)
	return time.Unix(unixTime, 0).In(cstSh)
}

// 解析指定地区时间
func TimeParseToZone(zone string, timeStr string, format string) time.Time {
	zone = countryToFile(zone)
	cstSh, _ := time.LoadLocation(zone)
	time1, _ := time.ParseInLocation(format, timeStr, cstSh)
	return time1
}

func StrToTimeFormatTheZone(strTime string, zone string) (time.Time, error) {
	zone = countryToFile(zone)
	var cstSh, _ = time.LoadLocation(zone) //上海
	if t, err := gtime.StrToTime(strTime); err == nil {
		return t.In(cstSh), nil
	} else {
		glog.Error(zone, err)
		return time.Time{}, err
	}
	return time.Time{}, nil
}

// 从时间戳获取PCR日数
func GetPcrUnixToDay(gameServer string, unixTime int64) int {
	return GetTimeAtUnixToZone(gameServer, unixTime).Day()
}

// 获取PCR今日开始时间
func GetPcrDayStartTimeToUnix(gameServer string) int64 {
	t := TimeParseToZone(gameServer, GetNowTimeToZone(gameServer).Format("2006-01-02"), "2006-01-02")
	// 每天5点刷新
	unix := t.Unix() + 18000
	if unix > time.Now().Unix() {
		unix -= 24 * 60 * 60
	}
	return unix
}

// 获取PCR今日结束时间
func GetPcrDayEndTimeToUnix(gameServer string) int64 {
	return GetPcrDayStartTimeToUnix(gameServer) + 86400
}

// 获取从文本日期获取PCR开始时间
func GetPcrStartTimeToUnixAtStr(gameServer string, timeStr string) int64 {
	t := TimeParseToZone(gameServer, timeStr, "2006-01-02")
	// 每天5点刷新
	unix := t.Unix() + 18000
	return unix
}

// 获取PCR结束时间
func GetPcrEndTimeToUnixAtStr(gameServer string, timeStr string) int64 {
	return GetPcrStartTimeToUnixAtStr(gameServer, timeStr) + 86400
}

// 获取PCR昨日开始时间
func GetPcrYesterdayStartTimeToUnix(gameServer string) int64 {
	t := TimeParseToZone(gameServer, GetNowTimeToZone(gameServer).Format("2006-01-02"), "2006-01-02")
	// 每天5点刷新
	unix := t.Unix() - 24*60*60 + 18000
	if unix > time.Now().Unix() {
		unix -= 24 * 60 * 60
	}
	return unix
}

// 获取PCR昨日结束时间
func GetPcrYesterdayEndTimeToUnix(gameServer string) int64 {
	return GetPcrYesterdayStartTimeToUnix(gameServer) + 86400
}

func countryToFile(name string) string {
	switch name {
	case "CN":
		return "Asia/Shanghai"
	case "TW":
		return "Asia/Taipei"
	case "JP":
		return "Asia/Tokyo"
	case "KR":
		return "Asia/Seoul"
	}
	return name
}
