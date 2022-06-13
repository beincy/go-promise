package main

import (
	"go-promise/Promise"
	"time"
)

func main() {
	getData1 := func() (interface{}, error) {
		return "data1", nil
	}
	getData2 := func() (interface{}, error) {
		time.Sleep(time.Second * 10)
		return "data2", nil
	}
	getData3 := func() (interface{}, error) {
		return 4, nil
	}
	// tasks := []Promise.Task{getData1, getData2, getData3}

	// result := Promise.All(tasks).WaitATime(300)

	tasks := map[string]Promise.Task{"1": getData1, "2": getData2, "3": getData3}
	result := Promise.Props(tasks).WaitATime(300)

	for k, r := range result {
		if r == nil {
			continue
		}
		if k == "3" {
			println(r.(int))
		} else {
			println(r.(string))
		}
	}
}

// 模仿js的promise的包，实验性的，不知道性能的
