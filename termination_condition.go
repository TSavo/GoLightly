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
	MaxCost int
}

func NewCostTerminationCondition(maxCost int) *CostTerminationCondition {
	return &CostTerminationCondition{maxCost}
}

func (term CostTerminationCondition) ShouldTerminate(p *Processor) bool {
	return term.MaxCost < p.Cost()
}

type TimeTerminationCondition struct {
	MaxTime   time.Duration
	StartTime int
}

func NewTimeTerminationCondition(maxTime time.Duration) *TimeTerminationCondition {
	return &TimeTerminationCondition{maxTime, int(time.Now().UnixNano())}
}

func (term *TimeTerminationCondition) Reset() {
	term.StartTime = int(time.Now().UnixNano())
}

func (term *TimeTerminationCondition) ShouldTerminate(p *Processor) bool {
	if int(term.MaxTime)+term.StartTime < int(time.Now().UnixNano()) {
		term.StartTime = int(time.Now().UnixNano())
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
