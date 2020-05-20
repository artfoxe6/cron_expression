# cron_expression

Implemented the standard definition of cron expression

#### Supported standard definitions

| Field  |  Required |  Allowed values | Allowed special characters   |
| ------------ | ------------ | ------------ | ------------ |
|  Minutes |  Yes | 0-59  |  \* , - / |
|  Hours |  Yes |   0-23 |  \* , - / |
| Day of month  | Yes  |  1-31 |  \* , - /  |
|  Month | Yes  |  1-12 or JAN-DEC |  \* , - /  |
| Day of week  |  Yes |  0-6 or SUN-SAT |  \* , - / |

#### Supported non-standard definitions

| Entry | Description | Equivalent to |
| ------ | ------ | ------ |
| @yearly (or @annually) | Run once a year at midnight of 1 January | 0 0 1 1 \* |
| @monthly | Run once a month at midnight of the first day of the month	 | 0 0 1 \* \* |
| @weekly | Run once a week at midnight on Sunday morning | 0 0 \* \* 0 |
| @daily (or @midnight) | Run once a day at midnight | 0 0 \* \* \* |
| @hourly | Run once an hour at the beginning of the hour | 0 \* \* \* \* |

#### Install

go get github.com/artfoxe6/cron_expression@v1.0.5

#### Usage

```go
package main

import (
	"fmt"
	"github.com/artfoxe6/cron_expression"
	"log"
	"time"
)
func main() {
	expr := cron_expression.NewExpression("* 1-10/2 * */2 *", "CST", 8*3600)
	dst := make([]string, 0)
	err := expr.Next(time.Now(), 5, &dst)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, v := range dst {
		fmt.Println(v)
	}
}

```
```
2020-05-09 09:38:00
2020-05-09 09:39:00
2020-05-09 09:40:00
2020-05-09 09:41:00
2020-05-09 09:42:00
```

#### Not Support Expression

|  example | replace to  |
| ------------ | ------------ |
| \* \* \*,27 \* \*   | \* \* \* \* \*  |
| \* 1-10/2,1,2,3 \* \* \*  | \* 1,2,3,5,7,9 \* \* \*  |

#### Controversial

On the basis of standards: <br />
If both dom and dow is not \*, anyone starts with \*/ <br />
The relationship becomes logical AND, <br />
Else logical OR <br />

#### Results verification
https://crontab.guru/

#### references
https://en.wikipedia.org/wiki/Cron <br/>
https://pubs.opengroup.org/onlinepubs/007904975/utilities/crontab.html <br/>
https://crontab.guru/cron-bug.html <br/>
https://crontab.guru/tips.html <br/>
