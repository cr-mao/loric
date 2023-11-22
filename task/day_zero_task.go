/**
User: cr-mao
Date: 2023/11/17 11:29
Email: crmao@qq.com
Desc: task.go
*/
package task

import (
	"time"

	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/component"
	"github.com/cr-mao/loric/log"
)

// 暂时只为 日任务、周任务用的
type DayZeroTask struct {
	component.Base
	Node      *node.Node
	DailyTask func(*node.Node)
	WeekTask  func(*node.Node)
}

func NewDayZeroTask(node *node.Node, dayTask, weekTask func(*node.Node)) *DayZeroTask {
	return &DayZeroTask{
		Node:      node,
		DailyTask: dayTask,
		WeekTask:  weekTask,
	}
}

func (tInstance *DayZeroTask) Start() {
	go func() {
		for {
			now := time.Now()
			////计算下一个零点,1秒开始算
			next := now.Add(time.Hour*24 + time.Second)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			go func(do func(*node.Node)) {
				defer func() {
					if err := recover(); err != nil {
						log.Errorf("day task error: %v", err)
					}
				}()
				do(tInstance.Node)
			}(tInstance.DailyTask)
			// 第一天则进行周任务
			if now.Weekday() == time.Monday {
				go func(do func(*node.Node)) {
					defer func() {
						if err := recover(); err != nil {
							log.Errorf("week task error: %v", err)
						}
					}()
					do(tInstance.Node)
				}(tInstance.WeekTask)
			}
			time.Sleep(time.Second * 10)
		}
	}()
}
