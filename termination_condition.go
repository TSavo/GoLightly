package govirtual

import (
	"time"
)

type TerminationCondition interface {
	ShouldTerminate(*Processor) bool
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

func (term AndTerminationCondition) ShouldTerminate(p *Processor) bool {
	for _, x := range term {
		if !(x).ShouldTerminate(p) {
			return false
		}
	}
	return true
}

func (term OrTerminationCondition) ShouldTerminate(p *Processor) bool {
	for _, x := range term {
		if (x).ShouldTerminate(p) {
			return true
		}
	}
	return false
}

func (term *NotTerminationCondition) ShouldTerminate(p *Processor) bool {
	return !(*(*term).NotCondition).ShouldTerminate(p)
}

type CostTerminationCondition struct {
	MaxCost int64
}

func NewCostTerminationCondition(maxCost int64) *CostTerminationCondition {
	return &CostTerminationCondition{maxCost}
}

func (term CostTerminationCondition) ShouldTerminate(p *Processor) bool {
	return term.MaxCost < p.Cost()
}

type TimeTerminationCondition struct {
	MaxTime   time.Duration
	StartTime int64
}

func NewTimeTerminationCondition(maxTime time.Duration) *TimeTerminationCondition {
	return &TimeTerminationCondition{maxTime, time.Now().UnixNano()}
}

func (term *TimeTerminationCondition) Reset() {
	term.StartTime = time.Now().UnixNano()
}

func (term *TimeTerminationCondition) ShouldTerminate(p *Processor) bool {
	if int64(term.MaxTime)+term.StartTime < time.Now().UnixNano() {
		term.StartTime = time.Now().UnixNano()
		return true
	} else {
		return false
	}
}

type ChannelTerminationCondition chan bool

func NewChannelTerminationCondition() *ChannelTerminationCondition {
	x := make(ChannelTerminationCondition, 1)
	return &x
}

func (term *ChannelTerminationCondition) ShouldTerminate(p *Processor) bool {
	select {
	case x := <-(*term):
		return x
	default:
		return false
	}
}
