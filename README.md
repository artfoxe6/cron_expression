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
go get github.com/codeab/cron_expression@v1.0.1
```

```
expr := cron_expression.NewExpression("* 1-10/2 * */2 *", "CST", 8*3600)
dst := make([]string, 0)
err := expr.Next(time.Now(), 5, &dst)
if err != nil {
    log.Fatalln(err.Error())
}
for _, v := range dst {
    fmt.Println(v)
}
```

#### 不支持

|  实例 | 请替换  |
| ------------ | ------------ |
| \* \* \*,27 \* \*   | \* \* \* \* \*  |
| \* 1-10/2,1,2,3 \* \* \*  | \* 1,2,3,5,7,9 \* \* \*  |

#### 存在争议的地方

当(日)dom和(周)dow均不为*时,存在一些争议, <br/>
首先说一下cron的实现标准,目前有多个版本标准(也可以说没有标准),根据wikipedia说明, <br/>
两者的关系应该是逻辑或,也就是满足dom或dow就执行 <br/>
而本项目实现了 crontab.guru 的标准, <br/>
在原有基础上增加了一条: <br/>
如果任意一方出现\*/num语法,两者的关系为逻辑与,相反为逻辑或 <br/>

#### 结果校验
由于实现标准的区别,复杂的表达式可能存在不相同的结果<br/>
https://crontab.guru/ <br/>
https://tool.lu/crontab/

#### 参考文献
https://en.wikipedia.org/wiki/Cron <br/>
https://pubs.opengroup.org/onlinepubs/007904975/utilities/crontab.html <br/>
https://crontab.guru/cron-bug.html <br/>
https://crontab.guru/tips.html <br/>