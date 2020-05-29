package cron_expression

import (
	"fmt"
	"testing"
	"time"
)

//测试接下来5个时间点
func TestNextAny(t *testing.T) {
	location := time.FixedZone("CST", 3600*8)
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-05-09 10:47:00", location)

	list := map[string][]string{
		"* */2 * */2 *": {
			"2020-05-09 10:48:00",
			"2020-05-09 10:49:00",
			"2020-05-09 10:50:00",
			"2020-05-09 10:51:00",
			"2020-05-09 10:52:00",
		},
		"* 1-10/2 * */2 *": {
			"2020-05-10 01:00:00",
			"2020-05-10 01:01:00",
			"2020-05-10 01:02:00",
			"2020-05-10 01:03:00",
			"2020-05-10 01:04:00",
		},
		"5 4 * * sun": {
			"2020-05-10 04:05:00",
			"2020-05-17 04:05:00",
			"2020-05-24 04:05:00",
			"2020-05-31 04:05:00",
			"2020-06-07 04:05:00",
		},
	}
	for k, v := range list {
		expr := NewExpression(k, "CST", 8*3600)
		dst := make([]string, 0)
		err := expr.NextAny(currentTime, 5, &dst)
		if err != nil {
			t.Fail()
		} else {
			for i := 0; i < 5; i++ {
				if dst[i] != v[i] {
					fmt.Println(dst[i], v[i])
					t.Fail()
				}
			}
		}
	}
}

//测试下一个时间点
func TestNext(t *testing.T) {
	location := time.FixedZone("CST", 3600*8)
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-05-09 10:47:00", location)

	list := map[string]string{
		"* */2 * */2 *":    "2020-05-09 10:48:00",
		"* 1-10/2 * */2 *": "2020-05-10 01:00:00",
		"5 4 * * sun":      "2020-05-10 04:05:00",
	}
	for k, v := range list {
		expr := NewExpression(k, "CST", 8*3600)
		dst, err := expr.Next(currentTime)
		if err != nil {
			t.Fail()
		} else {
			if dst.Format("2006-01-02 15:04:05") != v {
				t.Fail()
			}
		}
	}
}

//压测
func BenchmarkExpression(b *testing.B) {

	for i := 0; i < b.N; i++ {
		expr := NewExpression("* 1-10/2 * */2 *", "CST", 8*3600)
		dst := make([]string, 0)
		_ = expr.NextAny(time.Now(), 1, &dst)
	}
}
