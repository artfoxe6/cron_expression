package cron_expression

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

//表达式
type Expression struct {
	Minute         []int  //分 规则
	Hour           []int  //时 规则
	Dom            []int  //日 规则
	DomCopy        []int  //日 规则拷贝
	Month          []int  //月 规则
	Dow            []int  //周 规则
	LocationOffset int    //当前时区与UTC的偏移秒数
	LocationName   string //当前时区名
	CronRule       string //原始的cron规则
	IsParse        bool   //是否解析
}

//规则键位
const (
	MinuteKey = iota
	HourKey
	DomKey
	MonthKey
	DowKey
)

//创建一个表达式
func NewExpression(cronExpr, locationName string, locationOffset int) *Expression {
	return &Expression{
		Minute:         make([]int, 0),
		Hour:           make([]int, 0),
		Dom:            make([]int, 0),
		DomCopy:        make([]int, 0),
		Month:          make([]int, 0),
		Dow:            make([]int, 0),
		LocationOffset: locationOffset,
		LocationName:   locationName,
		CronRule:       cronExpr,
		IsParse:        false,
	}
}

//解析
func (expr *Expression) Parse() error {
	ruleItems := strings.Split(expr.CronRule, " ")
	if len(ruleItems) != 5 {
		return errors.New("cron规则格式错误")
	}
	var err error
	expr.Minute, err = cronRuleParse(ruleItems[MinuteKey], []int{0, 59})
	if err != nil {
		return errors.New("分钟" + err.Error())
	}
	expr.Hour, err = cronRuleParse(ruleItems[HourKey], []int{0, 23})
	if err != nil {
		return errors.New("小时" + err.Error())
	}
	expr.Dom, err = cronRuleParse(ruleItems[DomKey], []int{1, 31})
	if err != nil {
		return errors.New("日(月)" + err.Error())
	}
	expr.DomCopy = make([]int, len(expr.Dom))
	_ = copy(expr.DomCopy, expr.Dom)
	monthExpr := monthAliasToNumber(ruleItems[MonthKey])
	expr.Month, err = cronRuleParse(monthExpr, []int{1, 12})
	if err != nil {
		return errors.New("月" + err.Error())
	}
	dowExpr := DowAliasToNumber(ruleItems[DowKey])
	expr.Dow, err = cronRuleParse(dowExpr, []int{0, 6})
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

func (expr *Expression) nextMonth(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	change.month = 0
	change.day = 0
	change.hour = 0
	change.minute = 0
	nextAt.year = now.Year()
	current := int(now.Month())
	if *jump == "month" {
		change.month = 1
		nextAt.month = getMinValueInArray(expr.Month, current)
		if nextAt.month <= current {
			nextAt.year += 1
		}
	} else {
		nextAt.month = current
		if !existsInArray(expr.Month, current) {
			nextAt.month = getMinValueInArray(expr.Month, current)
			change.month = 1
		}
	}
	expr.Dom = make([]int, len(expr.DomCopy))
	_ = copy(expr.Dom, expr.DomCopy)
	expr.weekToDay(now, change, nextAt)
	return true
}
func (expr *Expression) nextDay(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	current := now.Day()
	if *jump != "day" {
		if change.month == 1 {
			nextAt.day = expr.Dom[0]
			if nextAt.day != current {
				change.day = 1
			}
			return true
		}
		if !existsInArray(expr.Dom, current) {
			nextAt.day = getMinValueInArray(expr.Dom, current)
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
		nextAt.day = getMinValueInArray(expr.Dom, current)
		if nextAt.day <= current {
			*jump = "month"
			return false
		}
		change.day = 1
	}
	return true
}
func (expr *Expression) nextHour(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	current := now.Hour()
	if *jump != "hour" {
		if change.day == 1 || change.month == 1 {
			nextAt.hour = expr.Hour[0]
			if nextAt.hour != current {
				change.hour = 1
			}
			return true
		}
		if !existsInArray(expr.Hour, current) {
			nextAt.hour = getMinValueInArray(expr.Hour, current)
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
		nextAt.hour = getMinValueInArray(expr.Hour, current)
		if nextAt.hour <= current {
			*jump = "day"
			return false
		}
		change.hour = 1
	}
	return true
}
func (expr *Expression) nextMinute(now time.Time, change *change, nextAt *nextAt, jump *string) bool {
	current := now.Minute()
	if change.hour == 1 || change.day == 1 || change.month == 1 {
		nextAt.minute = expr.Minute[0]
	} else {
		nextAt.minute = getMinValueInArray(expr.Minute, current)
		if nextAt.minute <= current {
			*jump = "hour"
			return false
		}
	}
	return true
}

//dow转dom
func (expr *Expression) weekToDay(now time.Time, change *change, nextAt *nextAt) {

	ruleItems := strings.Split(expr.CronRule, " ")
	if ruleItems[DowKey] == "*" {
		return
	}
	//dom和dow任意一个存在 间隔符 / 将形成交集
	days := getDayByWeek(now.Year(), nextAt.month, expr.Dow, expr.LocationName, expr.LocationOffset)
	fmt.Println(days)
	if strings.Contains(ruleItems[DowKey], "*/") || strings.Contains(ruleItems[DomKey], "*/") {
		expr.Dom = arrayIntersect(expr.Dom, days)
	} else {
		if ruleItems[DomKey] != "*" {
			expr.Dom = arrayMerge(expr.Dom, days)
		} else {
			expr.Dom = days
		}

	}
}

//计算下一个执行日期
func (expr *Expression) Next(current time.Time, nextStep int, dst *[]string) error {
	if expr.IsParse == false {
		err := expr.Parse()
		if err != nil {
			return err
		}
	}
	now := current.In(time.FixedZone(expr.LocationName, expr.LocationOffset))
	var (
		change = &change{}
		nextAt = &nextAt{}
	)
	//当前周期的上一个周期是否需要跳跃
	var jump = ""
	expr.Recursive(now, change, nextAt, &jump)
	//fmt.Println("分",expr.Minute)
	//fmt.Println("时",expr.Hour)
	fmt.Println("日",expr.Dom)
	fmt.Println("月",expr.Month)
	//fmt.Println("周",expr.Dow)
	res := time.Date(nextAt.year, time.Month(nextAt.month), nextAt.day, nextAt.hour, nextAt.minute, 0, 0, now.Location())
	*dst = append(*dst, res.Format("2006-01-02 15:04:05"))
	if nextStep == 1 {
		return nil
	}
	return expr.Next(res, nextStep-1, dst)
}

//计算下一个日期
func (expr *Expression) Recursive(current time.Time, change *change, at *nextAt, jump *string) {
	expr.nextMonth(current, change, at, jump)
	if b := expr.nextDay(current, change, at, jump); !b {
		expr.Recursive(current, change, at, jump)
	}
	if b := expr.nextHour(current, change, at, jump); !b {
		expr.Recursive(current, change, at, jump)
	}
	if b := expr.nextMinute(current, change, at, jump); !b {
		expr.Recursive(current, change, at, jump)
	}
}
