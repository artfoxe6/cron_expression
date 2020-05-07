# cron_expression

实现了cron的标准定义,不支持任何非标准宏
>维基百科: https://en.wikipedia.org/wiki/Cron

## 使用方式

```go
    // 参数: cron表达式 当地时区 当地时区和UTC的时差
	expr := NewExpression("5 4 * * sun", "CST", 8*3600)
	dst := make([]string, 0)
	// 参数: 计算开始时间 计算下几个执行点 接收结果 tip:指定不同的开始时间可以实现时间穿梭
	err := expr.Next(time.Now(), 5, &dst)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range dst {
		fmt.Println(v)
	}
```

## 容易误解的点

如果dom日(月) 和 dow周(月)都不为*, 那么任何一方出现 */num 语法,两者关系变为且,否则为或

## 结果校验
>crontab guru: https://crontab.guru/
>
>在线工具: https://tool.lu/crontab/