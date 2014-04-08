package vm

import (
	"time"
	"fmt"
)

type TerminationCondition interface {
	ShouldTerminate(*ProcessorCore) bool
}

type AndTerminationCondition []TerminationCondition

type OrTerminationCondition []TerminationCondition

type NotTerminationCondition struct {
	NotCondition *TerminationCondition
}

func AndTerminate(term ...TerminationCondition) *AndTerminationCondition {
	out := AndTerminationCondition(term)
	return &out
}

func OrTerminate(term ...TerminationCondition) *OrTerminationCondition {
	out := OrTerminationCondition(term)
	return &out
}

func NotTerminate(term *TerminationCondition) *NotTerminationCondition {
	return &NotTerminationCondition{term}
}

func (term AndTerminationCondition) ShouldTerminate(p *ProcessorCore) bool {
	for _, x := range term {
		if !x.ShouldTerminate(p) {
			return false
		}
	}
	return true
}

func (term OrTerminationCondition) ShouldTerminate(p *ProcessorCore) bool {
	for _, x := range term {
		if x.ShouldTerminate(p) {
			return true
		}
	}
	return false
}

func (term *NotTerminationCondition) ShouldTerminate(p *ProcessorCore) bool {
	return !(*(*term).NotCondition).ShouldTerminate(p)
}

type CostTerminationCondition struct {
	MaxCost int64
}

func NewCostTerminationCondition(maxCost int64) *CostTerminationCondition {
	return &CostTerminationCondition{maxCost}
}

func (term CostTerminationCondition) ShouldTerminate(p *ProcessorCore) bool {
	return term.MaxCost < p.Cost()
}

type TimeTerminationCondition struct {
	MaxTime int64
}

func NewTimeTerminationCondition(maxTime int64) *TimeTerminationCondition {
	fmt.Println(maxTime)
	return &TimeTerminationCondition{maxTime}
}

func (term TimeTerminationCondition) ShouldTerminate(p *ProcessorCore) bool {
	return term.MaxTime > time.Now().UnixNano()-p.StartTime
}

type ChannelTerminationCondition chan bool

func NewChannelTerminationCondition() *ChannelTerminationCondition {
	x := make(ChannelTerminationCondition, 1)
	return &x
}

func (term *ChannelTerminationCondition) ShouldTerminate(p *ProcessorCore) bool {
	select {
	case x := <-(*term):
		return x
	default:
		return false
	}
}
