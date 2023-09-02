package web

import (
	"fmt"
	"time"
)

type Filter HandlerFunc
type FilterBuilder func(next Filter) Filter

func MetricsFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		start := time.Now().UnixNano()
		next(c)
		end := time.Now().UnixNano()
		fmt.Printf("run time: %d ns\n", end-start)
	}
}
