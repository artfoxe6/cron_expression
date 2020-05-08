package cron_expression

import (
	"log"
	"testing"
	"time"
)

func TestExpression(t *testing.T) {
	expr := NewExpression("* 1-10/2 * */2 *", "CST", 8*3600)
	dst := make([]string, 0)
	err := expr.Next(time.Now(), 5, &dst)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, v := range dst {
		log.Println(v)
	}
}
