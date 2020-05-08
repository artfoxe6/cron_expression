package cron_expression

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	numberList = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	}
	numberMaps = map[string]int{
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"00": 0, "01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19,
		"20": 20, "21": 21, "22": 22, "23": 23, "24": 24, "25": 25, "26": 26, "27": 27, "28": 28, "29": 29,
		"30": 30, "31": 31, "32": 32, "33": 33, "34": 34, "35": 35, "36": 36, "37": 37, "38": 38, "39": 39,
		"40": 40, "41": 41, "42": 42, "43": 43, "44": 44, "45": 45, "46": 46, "47": 47, "48": 48, "49": 49,
		"50": 50, "51": 51, "52": 52, "53": 53, "54": 54, "55": 55, "56": 56, "57": 57, "58": 58, "59": 59,
		"1970": 1970, "1971": 1971, "1972": 1972, "1973": 1973, "1974": 1974, "1975": 1975, "1976": 1976, "1977": 1977, "1978": 1978, "1979": 1979,
		"1980": 1980, "1981": 1981, "1982": 1982, "1983": 1983, "1984": 1984, "1985": 1985, "1986": 1986, "1987": 1987, "1988": 1988, "1989": 1989,
		"1990": 1990, "1991": 1991, "1992": 1992, "1993": 1993, "1994": 1994, "1995": 1995, "1996": 1996, "1997": 1997, "1998": 1998, "1999": 1999,
		"2000": 2000, "2001": 2001, "2002": 2002, "2003": 2003, "2004": 2004, "2005": 2005, "2006": 2006, "2007": 2007, "2008": 2008, "2009": 2009,
		"2010": 2010, "2011": 2011, "2012": 2012, "2013": 2013, "2014": 2014, "2015": 2015, "2016": 2016, "2017": 2017, "2018": 2018, "2019": 2019,
		"2020": 2020, "2021": 2021, "2022": 2022, "2023": 2023, "2024": 2024, "2025": 2025, "2026": 2026, "2027": 2027, "2028": 2028, "2029": 2029,
		"2030": 2030, "2031": 2031, "2032": 2032, "2033": 2033, "2034": 2034, "2035": 2035, "2036": 2036, "2037": 2037, "2038": 2038, "2039": 2039,
		"2040": 2040, "2041": 2041, "2042": 2042, "2043": 2043, "2044": 2044, "2045": 2045, "2046": 2046, "2047": 2047, "2048": 2048, "2049": 2049,
		"2050": 2050, "2051": 2051, "2052": 2052, "2053": 2053, "2054": 2054, "2055": 2055, "2056": 2056, "2057": 2057, "2058": 2058, "2059": 2059,
		"2060": 2060, "2061": 2061, "2062": 2062, "2063": 2063, "2064": 2064, "2065": 2065, "2066": 2066, "2067": 2067, "2068": 2068, "2069": 2069,
		"2070": 2070, "2071": 2071, "2072": 2072, "2073": 2073, "2074": 2074, "2075": 2075, "2076": 2076, "2077": 2077, "2078": 2078, "2079": 2079,
		"2080": 2080, "2081": 2081, "2082": 2082, "2083": 2083, "2084": 2084, "2085": 2085, "2086": 2086, "2087": 2087, "2088": 2088, "2089": 2089,
		"2090": 2090, "2091": 2091, "2092": 2092, "2093": 2093, "2094": 2094, "2095": 2095, "2096": 2096, "2097": 2097, "2098": 2098, "2099": 2099,
	}
	monthMaps = map[string]string{
		`jan`: `1`, `january`: `1`,
		`feb`: `2`, `february`: `2`,
		`mar`: `3`, `march`: `3`,
		`apr`: `4`, `april`: `4`,
		`may`: `5`,
		`jun`: `6`, `june`: `6`,
		`jul`: `7`, `july`: `7`,
		`aug`: `8`, `august`: `8`,
		`sep`: `9`, `september`: `9`,
		`oct`: `10`, `october`: `10`,
		`nov`: `11`, `november`: `11`,
		`dec`: `12`, `december`: `12`,
	}
	dowMaps = map[string]string{
		`sun`: `0`, `sunday`: `0`,
		`mon`: `1`, `monday`: `1`,
		`tue`: `2`, `tuesday`: `2`,
		`wed`: `3`, `wednesday`: `3`,
		`thu`: `4`, `thursday`: `4`,
		`fri`: `5`, `friday`: `5`,
		`sat`: `6`, `saturday`: `6`,
	}
	alias = map[string]string{
		"@yearly":   "0 0 1 1 *",
		"@annually": "0 0 1 1 *",
		"@monthly":  "0 0 1 * *",
		"@weekly":   "0 0 * * 0",
		"@midnight": "0 0 * * *",
		"@daily":    "0 0 * * *",
		"@hourly":   "0 * * * *",
	}
)

//将Month的英文格式转为数字
func monthAliasToNumber(str string) string {
	for k, v := range monthMaps {
		str = strings.ReplaceAll(str, k, v)
	}
	return str
}

//将Dow的英文格式转为数字
func DowAliasToNumber(str string) string {
	for k, v := range dowMaps {
		str = strings.ReplaceAll(str, k, v)
	}
	return str
}

//从规则中解析出范围内符合的数字
func cronRuleParse(str string, limit []int) ([]int, error) {

	limitList := numberList[limit[0] : limit[1]+1]
	// "*"
	if b, _ := regexp.MatchString(`^\*$`, str); b {
		return limitList, nil
	}
	// "10"
	if b, _ := regexp.MatchString(`^\d{1,2}$`, str); b {
		return []int{numberMaps[str]}, nil
	}
	// "1-12"
	if b, _ := regexp.MatchString(`^\d{1,2}-\d{1,2}$`, str); b {
		arr := strings.Split(str, "-")
		return numberList[numberMaps[arr[0]] : numberMaps[arr[1]]+1], nil
	}
	// "1,2,3,30"
	if b, _ := regexp.MatchString(`^(\d{1,2},)+\d{1,2}$`, str); b {
		arr := strings.Split(str, ",")
		temp := make([]int, 0)
		for i := 0; i < len(arr); i++ {
			temp = append(temp, numberMaps[arr[i]])
		}
		sort.Ints(temp)
		return temp, nil
	}
	// "*/3"
	if b, _ := regexp.MatchString(`^\*/\d{1,2}$`, str); b {
		arr := strings.Split(str, "/")
		temp := make([]int, 0)
		for i := 0; i < len(limitList); i += numberMaps[arr[1]] {
			temp = append(temp, limitList[i])
		}
		return temp, nil
	}
	// "3-10/3"
	if b, _ := regexp.MatchString(`^\d{1,2}-\d{1,2}/\d{1,2}$`, str); b {
		arr := strings.Split(str, "/")
		temp := make([]int, 0)
		subArr := strings.Split(arr[0], "-")
		for i := numberMaps[subArr[0]]; i < numberMaps[subArr[1]]; i += numberMaps[arr[1]] {
			temp = append(temp, numberMaps[strconv.Itoa(i)])
		}
		return temp, nil
	}
	return nil, errors.New("格式设置错误")
}

//判断一个元素是否在数组中
func existsInArray(arr []int, item int) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

//在数组中寻找比给定目标大的最小的元素
func getMinValueInArray(arr []int, item int) int {
	sort.Ints(arr)
	for _, v := range arr {
		if v > item {
			return v
		}
	}
	//如果目标是当前范围的最大值,就返回最小
	return arr[0]
}

//获取某月下面有多少天
func getDayCountInMonth(year, month int) int {
	monthDay := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	monthDayLeapYear := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month > 12 {
		month = month - 12
	}
	if isLeapYear(year) {
		return monthDayLeapYear[month-1]
	}
	return monthDay[month-1]
}

//判断某年是不是闰年
func isLeapYear(year int) bool {
	if year%100 != 0 && year%4 == 0 {
		return true
	}
	if year%100 == 0 && year%400 == 0 {
		return true
	}
	return false
}

//通过当月的weekdays,计算出对等的days
func getDayByWeek(year int, month int, weekdays []int, locationName string, locationOffset int) []int {
	var days = make([]int, 0)
	monthHasDay := getDayCountInMonth(year, month)
	location := time.FixedZone(locationName, locationOffset)
	for i := 1; i <= monthHasDay; i++ {
		t := time.Date(year, time.Month(month), i, 0, 0, 0, 0, location)
		if existsInArray(weekdays, int(t.Weekday())) {
			days = append(days, i)
		}
	}
	return days
}

//两个切片的交集
func arrayIntersect(a, b []int) []int {
	res := make([]int, 0)
	for _, v := range a {
		if existsInArray(b, v) {
			res = append(res, v)
		}
	}
	return res
}

//两个切片的并集
func arrayMerge(a, b []int) []int {
	res := make([]int, 0)
	for _, v := range a {
		if !existsInArray(res, v) {
			res = append(res, v)
		}
	}
	for _, v := range b {
		if !existsInArray(res, v) {
			res = append(res, v)
		}
	}
	return res
}
