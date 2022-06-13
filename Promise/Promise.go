package Promise

import (
	"context"
	"sync"
	"time"
)

const (
	Running int = iota
	Pending
	Stopped
	Completed
)

type Task func() (interface{}, error)

type Feature struct {
	Task   Task
	Result interface{}
	Error  error
	Status int
	TaskId int
}

type FeatureList struct {
	Futures  []*Feature
	Result   []interface{}
	Error    error
	Status   int
	Schedule sync.WaitGroup
}

type FeatureMap struct {
	Futures  []*Feature
	Result   []interface{}
	Error    error
	Status   int
	Schedule sync.WaitGroup
	KeyMap   map[int]string
}

func All(fs []Task) *FeatureList {
	featureList := &FeatureList{
		Status:  Running,
		Futures: make([]*Feature, len(fs)),
	}
	featureList.Schedule.Add(len(fs))
	for i, f := range fs {
		featureList.Futures[i] = &Feature{
			Task:   f,
			Status: Pending,
			TaskId: i,
		}
		go func(f Task, i int) {
			defer featureList.Schedule.Done()
			result, err := f()
			status := Completed
			if err != nil {
				status = Stopped
			}
			featureList.Futures[i].Result = result
			featureList.Futures[i].Error = err
			featureList.Futures[i].Status = status

		}(f, i)
	}
	return featureList
}

func (p *FeatureList) Await() []interface{} {
	p.Result = make([]interface{}, len(p.Futures))
	p.Schedule.Wait()
	for _, v := range p.Futures {
		if v.Status == Completed {
			p.Result[v.TaskId] = v.Result
		}
		if v.Error != nil {
			p.Error = v.Error
		}
	}
	return p.Result
}

func (p *FeatureList) WaitATime(millisecond int) []interface{} {
	p.Result = make([]interface{}, len(p.Futures))
	ctx, cancle := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(millisecond))
	go func(ctx context.Context, cancle context.CancelFunc) {
		defer cancle()
		p.Schedule.Wait()
	}(ctx, cancle)
	<-ctx.Done()
	for _, v := range p.Futures {
		if v.Status == Completed {
			p.Result[v.TaskId] = v.Result
		}
		if v.Error != nil {
			p.Error = v.Error
		}
	}
	return p.Result
}

func Props(fs map[string]Task) *FeatureMap {
	featureMap := &FeatureMap{
		Status:  Running,
		Futures: make([]*Feature, len(fs)),
		KeyMap:  make(map[int]string, len(fs)),
	}
	featureMap.Schedule.Add(len(fs))
	i := 0
	for k, f := range fs {
		featureMap.KeyMap[i] = k
		featureMap.Futures[i] = &Feature{
			Task:   f,
			Status: Pending,
			TaskId: i,
		}
		go func(f Task, i int) {
			defer featureMap.Schedule.Done()
			result, err := f()
			status := Completed
			if err != nil {
				status = Stopped
			}
			featureMap.Futures[i].Result = result
			featureMap.Futures[i].Error = err
			featureMap.Futures[i].Status = status
		}(f, i)
		i++
	}
	return featureMap
}
func (p *FeatureMap) Await() map[string]interface{} {
	p.Result = make([]interface{}, len(p.Futures))
	p.Schedule.Wait()
	result := make(map[string]interface{}, len(p.Futures))
	for _, v := range p.Futures {
		if v.Status == Completed {
			p.Result[v.TaskId] = v.Result
			result[p.KeyMap[v.TaskId]] = v.Result
		}
		if v.Error != nil {
			p.Error = v.Error
		}
	}
	return result
}

func (p *FeatureMap) WaitATime(millisecond int) map[string]interface{} {
	p.Result = make([]interface{}, len(p.Futures))
	result := make(map[string]interface{}, len(p.Futures))
	ctx, cancle := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(millisecond))
	go func(ctx context.Context, cancle context.CancelFunc) {
		defer cancle()
		p.Schedule.Wait()
	}(ctx, cancle)
	<-ctx.Done()
	for _, v := range p.Futures {
		if v.Status == Completed {
			p.Result[v.TaskId] = v.Result
			result[p.KeyMap[v.TaskId]] = v.Result
		}
		if v.Error != nil {
			p.Error = v.Error
		}
	}
	return result
}
