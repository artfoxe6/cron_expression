package cron_expression

import (
	"fmt"
	"testing"
	"time"
)

func TestExpression(t *testing.T) {
	// 参数: cron表达式 当地时区 当地时区和UTC的时差
	expr := NewExpression("0 0 2-10/3 */2 */2", "CST", 8*3600)
	dst := make([]string, 0)
	// 参数: 计算开始时间 计算下几个执行点 接收结果 tip:指定不同的开始时间可以实现时间穿越
	err := expr.Next(time.Now(), 1, &dst)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range dst {
		fmt.Println(v)
	}
}
