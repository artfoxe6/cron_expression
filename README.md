# cron_expression

#### 实现了cron的标准定义,以及部分非标准宏

标准定义

| Field  |  Required |  Allowed values | Allowed special characters   |
| ------------ | ------------ | ------------ | ------------ |
|  Minutes |  Yes | 0-59  |  \* , - / |
|  Hours |  Yes |   0-23 |  \* , - / |
| Day of month  | Yes  |  1-31 |  \* , - /  |
|  Month | Yes  |  1-12 or JAN-DEC |  \* , - /  |
| Day of week  |  Yes |  0-6 or SUN-SAT |  \* , - / |

支持的非标准宏

| Entry | Description | Equivalent to |
| ------ | ------ | ------ |
| @yearly (or @annually) | Run once a year at midnight of 1 January | 0 0 1 1 \* |
| @monthly | Run once a month at midnight of the first day of the month	 | 0 0 1 \* \* |
| @weekly | Run once a week at midnight on Sunday morning | 0 0 \* \* 0 |
| @daily (or @midnight) | Run once a day at midnight | 0 0 \* \* \* |
| @hourly | Run once an hour at the beginning of the hour | 0 \* \* \* \* |

#### 使用方式

```
go get github.com/classfoxe6/cron_expression@v1.0.5
```

```
//支持标准格式 分 时 日 月 周,以及短语 @monthly @daily 等
expr := cron_expression.NewExpression("* 1-10/2 * */2 *", "CST", 8*3600)
dst := make([]string, 0)
//tip: 当前时间可以向前或向后任意指定,实现时间穿梭,达到计算当前时间之前的执行点
err := expr.Next(time.Now(), 5, &dst)
if err != nil {
    log.Fatalln(err.Error())
}
for _, v := range dst {
    fmt.Println(v)
}
```
```
2020-05-09 09:38:00
2020-05-09 09:39:00
2020-05-09 09:40:00
2020-05-09 09:41:00
2020-05-09 09:42:00
```

#### 不支持

|  实例 | 请替换  |
| ------------ | ------------ |
| \* \* \*,27 \* \*   | \* \* \* \* \*  |
| \* 1-10/2,1,2,3 \* \* \*  | \* 1,2,3,5,7,9 \* \* \*  |

#### 存在争议的地方

当(日)dom和(周)dow均不为*时,存在一些争议, 目前cron实现标准其实不统一,根据wikipedia说明, <br/>
两者的关系应该是逻辑或,也就是满足dom或dow就执行,而本项目参照了 crontab.guru 的实现标准, <br/>
在原有基础上增加了一层判断: 如果任意一方以\*/开头,两者的关系则为逻辑与,相反为逻辑或 <br/>

#### 结果校验
https://crontab.guru/

#### 参考文献
https://en.wikipedia.org/wiki/Cron <br/>
https://pubs.opengroup.org/onlinepubs/007904975/utilities/crontab.html <br/>
https://crontab.guru/cron-bug.html <br/>
https://crontab.guru/tips.html <br/>
