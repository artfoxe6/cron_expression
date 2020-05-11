package cron_expression

import (
	"errors"
	"strings"
	"time"
)

//表达式
type expression struct {
	minute         []int  //分 规则
	hour           []int  //时 规则
	dom            []int  //日 规则
	domCopy        []int  //日 规则拷贝
	month          []int  //月 规则
	dow            []int  //周 规则
	locationOffset int    //当前时区与UTC的偏移秒数
	locationName   string //当前时区名
	cronRule       string //原始的cron规则
	isParse        bool   //是否解析
}

//规则键位
const (
	minuteKey = iota
	hourKey
	domKey
	monthKey
	dowKey
)

//创建一个表达式 参数实例: "* 1-10/2,1,2,3 * */2 *", "CST", 8*3600
func NewExpression(cronExpr, locationName string, locationOffset int) *expression {
	return &expression{
		minute:         make([]int, 0),
		hour:           make([]int, 0),
		dom:            make([]int, 0),
		domCopy:        make([]int, 0),
		month:          make([]int, 0),
		dow:            make([]int, 0),
		locationOffset: locationOffset,
		locationName:   locationName,
		cronRule:       cronExpr,
		isParse:        false,
	}
}

//解析
func (expr *expression) parse() error {
	if strings.ContainsRune(expr.cronRule, '@') {
		v, ok := alias[expr.cronRule]
		if !ok {
			return errors.New("不支持的宏")
		}
		expr.cronRule = v
	}
	ruleItems := strings.Split(expr.cronRule, " ")
	if len(ruleItems) != 5 {
		return errors.New("cron规则格式错误")
	}
	var err error
	expr.minute, err = cronRuleParse(ruleItems[minuteKey], []int{0, 59})
	if err != nil {
		return errors.New("分钟" + err.Error())
	}
	expr.hour, err = cronRuleParse(ruleItems[hourKey], []int{0, 23})
	if err != nil {
		return errors.New("小时" + err.Error())
	}
	expr.dom, err = cronRuleParse(ruleItems[domKey], []int{1, 31})
	if err != nil {
		return errors.New("日(月)" + err.Error())
	}
	expr.domCopy = make([]int, len(expr.dom))
	_ = copy(expr.domCopy, expr.dom)
	monthExpr := monthAliasToNumber(ruleItems[monthKey])
	expr.month, err = cronRuleParse(monthExpr, []int{1, 12})
	if err != nil {
		return errors.New("月" + err.Error())
	}
	dowExpr := dowAliasToNumber(ruleItems[dowKey])
	expr.dow, err = cronRuleParse(dowExpr, []int{0, 6})
	if err != nil {
		return errors.New("周(月)" + err.Error())
	}
	return nil
}

//时间点变更记录
type change struct {
	month  int
	day    int
	hour   int
	minute int
}

//下一个执行时间点
type nextAt struct {
	year   int
	month  int
	day    int
	hour   int
	minute int
	week   int
}

//计算下一个month
func (expr *expression) nextMonth(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	change.month = 0
	change.day = 0
	change.hour = 0
	change.minute = 0
	nextAt.year = now.Year()
	current := int(now.Month())
	if *jump == "month" {
		change.month = 1
		nextAt.month = getMinValueInArray(expr.month, current)
		if nextAt.month <= current {
			nextAt.year += 1
		}
	} else {
		nextAt.month = current
		if !existsInArray(expr.month, current) {
			nextAt.month = getMinValueInArray(expr.month, current)
			change.month = 1
		}
	}
	expr.dom = make([]int, len(expr.domCopy))
	_ = copy(expr.dom, expr.domCopy)
	expr.weekToDay(now, change, nextAt)
	return true
}
//计算下一个day
func (expr *expression) nextDay(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	current := now.Day()
	if *jump != "day" {
		if change.month == 1 {
			nextAt.day = expr.dom[0]
			if nextAt.day != current {
				change.day = 1
			}
			return true
		}
		if !existsInArray(expr.dom, current) {
			nextAt.day = getMinValueInArray(expr.dom, current)
			if nextAt.day < current {
				*jump = "month"
				return false
			}
			if nextAt.day != current {
				change.day = 1
			}
		} else {
			nextAt.day = current
		}
	} else {
		nextAt.day = getMinValueInArray(expr.dom, current)
		if nextAt.day <= current {
			*jump = "month"
			return false
		}
		change.day = 1
	}
	return true
}
//计算下一个hour
func (expr *expression) nextHour(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	current := now.Hour()
	if *jump != "hour" {
		if change.day == 1 || change.month == 1 {
			nextAt.hour = expr.hour[0]
			if nextAt.hour != current {
				change.hour = 1
			}
			return true
		}
		if !existsInArray(expr.hour, current) {
			nextAt.hour = getMinValueInArray(expr.hour, current)
			if nextAt.hour < current {
				*jump = "day"
				return false
			}
			if nextAt.hour != current {
				change.hour = 1
			}
		} else {
			nextAt.hour = current
		}
	} else {
		nextAt.hour = getMinValueInArray(expr.hour, current)
		if nextAt.hour <= current {
			*jump = "day"
			return false
		}
		change.hour = 1
	}
	return true
}
//计算下一个minute
func (expr *expression) nextMinute(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	current := now.Minute()
	if change.hour == 1 || change.day == 1 || change.month == 1 {
		nextAt.minute = expr.minute[0]
	} else {
		nextAt.minute = getMinValueInArray(expr.minute, current)
		if nextAt.minute <= current {
			*jump = "hour"
			return false
		}
	}
	return true
}

//dow转dom
func (expr *expression) weekToDay(now time.Time, change *change, nextAt *nextAt) {

	ruleItems := strings.Split(expr.cronRule, " ")
	if ruleItems[dowKey] == "*" {
		return
	}
	//dom和dow任意一个以*/开头将形成交集
	days := getDayByWeek(now.Year(), nextAt.month, expr.dow, expr.locationName, expr.locationOffset)
	if strings.Contains(ruleItems[dowKey], "*/") || strings.Contains(ruleItems[domKey], "*/") {
		expr.dom = arrayIntersect(expr.dom, days)
	} else {
		if ruleItems[domKey] != "*" {
			expr.dom = arrayMerge(expr.dom, days)
		} else {
			expr.dom = days
		}

	}
}

//计算下一个执行日期
//开始时间,下几个执行点 接收结果
//tip:指定不同的开始时间可以实现时间穿梭
func (expr *expression) Next(current time.Time, nextStep int, dst *[]string) error {
	if expr.isParse == false {
		err := expr.parse()
		if err != nil {
			return err
		}
	}
	now := current.In(time.FixedZone(expr.locationName, expr.locationOffset))
	var (
		change = &change{}
		nextAt = &nextAt{}
	)
	//当前周期的上一个周期是否需要跳跃
	var jump = ""
	expr.recursive(now, change, nextAt, &jump)
	res := time.Date(nextAt.year, time.Month(nextAt.month), nextAt.day, nextAt.hour, nextAt.minute, 0, 0, now.Location())
	*dst = append(*dst, res.Format("2006-01-02 15:04:05"))
	if nextStep == 1 {
		return nil
	}
	return expr.Next(res, nextStep-1, dst)
}

//计算下一个日期
func (expr *expression) recursive(current time.Time, change *change, at *nextAt, jump *string) {
	expr.nextMonth(current, change, at, jump)
	if b := expr.nextDay(current, change, at, jump); !b {
		expr.recursive(current, change, at, jump)
	}
	if b := expr.nextHour(current, change, at, jump); !b {
		expr.recursive(current, change, at, jump)
	}
	if b := expr.nextMinute(current, change, at, jump); !b {
		expr.recursive(current, change, at, jump)
	}
}
