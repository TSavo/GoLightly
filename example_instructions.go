package govirtual

import (
	"math"
//	"fmt"
)

func EmulationInstructions(i *InstructionSet) {
	i.Instruction("noop", func(p *Processor, m *Memory) {
	})
	i.Movement("jump", func(p *Processor, m *Memory) {
		p.Jump(m.Get(0))
	})
	i.Movement("jumpIfZero", func(p *Processor, m *Memory) {
		if p.Registers.GetCardinal((*m).GetCardinal(0)) == 0 {
			p.Jump(m.Get(1))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfNotZero", func(p *Processor, m *Memory) {
		if p.Registers.GetCardinal((*m).Get(0)) != 0 {
			p.Jump(m.Get(1))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfEquals", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(0)) == p.Registers.Get((*m).Get(1)) {
			p.Jump(m.Get(2))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfNotEquals", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(0)) != p.Registers.Get((*m).Get(1)) {
			p.Jump(m.Get(2))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfGreaterThan", func(p *Processor, m *Memory) {
		if p.Registers.GetCardinal((*m).Get(0)) > p.Registers.GetCardinal((*m).Get(1)) {
			p.Jump(m.Get(2))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfLessThan", func(p *Processor, m *Memory) {
		if p.Registers.GetCardinal((*m).Get(0)) < p.Registers.GetCardinal((*m).Get(1)) {
			p.Jump(m.Get(2))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("call", func(p *Processor, m *Memory) {
		p.Call((*m).Get(0))
	})
	i.Movement("return", func(p *Processor, m *Memory) {
		p.Return()
	})
	i.Instruction("set", func(p *Processor, m *Memory) {
		p.Registers.Set((*m).Get(0), (*m).Get(1))
	})
	i.Instruction("store", func(p *Processor, m *Memory) {
		p.Heap.Set(m.Get(0), p.Registers.Get(m.Get(1)))
	})
	i.Instruction("load", func(p *Processor, m *Memory) {
		p.Registers.Set(m.Get(0), p.Heap.Get(m.Get(1)))
		//fmt.Println(m.Get(0), "=", p.Heap.Get(m.Get(1)))
	})
	i.Instruction("swap", func(p *Processor, m *Memory) {
		x := p.Registers.Get((*m).Get(0))
		p.Registers.Set((*m).Get(0), (*m).Get(1))
		p.Registers.Set((*m).Get(1), x)
	})
	/*
		i.Instruction("push", func(p *Processor, m *Memory) {
			p.Stack.Push(p.Registers.Get((*m).Get(0)))
		})
		i.Instruction("pop", func(p *Processor, m *Memory) {
			if x, err := p.Stack.Pop(); !err {
				p.Registers.Set((*m).Get(0), x)
			}
		})
		i.Instruction("increment", func(p *Processor, m *Memory) {
			p.Registers.Increment((*m).Get(0))
		})
		i.Instruction("decrement", func(p *Processor, m *Memory) {
			p.Registers.Decrement((*m).Get(0))
		})*/
	i.Instruction("add", func(p *Processor, m *Memory) {
		switch l := p.Registers.Get((*m).Get(0)).(type) {
		case int:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l+r)
			case float64:
				p.Registers.Set((*m).Get(2), float64(l)+r)
			}
		case float64:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l+float64(r))
			case float64:
				p.Registers.Set((*m).Get(2), l+r)
			}
		}
	})
	i.Instruction("subtract", func(p *Processor, m *Memory) {
		switch l := p.Registers.Get((*m).Get(0)).(type) {
		case int:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l-r)
			case float64:
				p.Registers.Set((*m).Get(2), float64(l)-r)
			}
		case float64:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l-float64(r))
			case float64:
				p.Registers.Set((*m).Get(2), l-r)
			}
		}
	})
	i.Instruction("multiply", func(p *Processor, m *Memory) {
		switch l := p.Registers.Get((*m).Get(0)).(type) {
		case int:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l*r)
			case float64:
				p.Registers.Set((*m).Get(2), float64(l)*r)
			}
		case float64:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l*float64(r))
			case float64:
				p.Registers.Set((*m).Get(2), l*r)
			}
		}
	})
	i.Instruction("divide", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(1)) == 0 {
			return
		}
		switch l := p.Registers.Get((*m).Get(0)).(type) {
		case int:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l/r)
			case float64:
				p.Registers.Set((*m).Get(2), float64(l)/r)
			}
		case float64:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l/float64(r))
			case float64:
				p.Registers.Set((*m).Get(2), l/r)
			}
		}

	})
	i.Instruction("modulos", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(1)) == 0 {
			return
		}
		switch l := p.Registers.Get((*m).Get(0)).(type) {
		case int:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), l%r)
			case float64:
				p.Registers.Set((*m).Get(2), math.Mod(float64(l), r))
			}
		case float64:
			switch r := p.Registers.Get((*m).Get(1)).(type) {
			case int:
				p.Registers.Set((*m).Get(2), math.Mod(l, float64(r)))
			case float64:
				p.Registers.Set((*m).Get(2), math.Mod(l, r))
			}
		}

	})
	i.Instruction("abs", func(p *Processor, m *Memory) {
		p.Registers.Set((*m).Get(0), Abs(p.Registers.Get((*m).Get(0))))
	})
}

func Abs(abs interface{}) interface{} {
	switch a := abs.(type) {
	case int:
		if a < 0 {
			return a * -1
		}
		return abs
	case float64:
		return math.Abs(a)
	}
	return abs
}
