package cron_expression

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestExpression(t *testing.T) {
	expr := NewExpression("0 0 * * *", "CST", 8*3600)
	dst := make([]string, 0)
	err := expr.Next(time.Now(), 5, &dst)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, v := range dst {
		fmt.Println(v)
	}
}
