package govirtual

import (
	"math"
)

func EmulationInstructions(i *InstructionSet) {
	i.Instruction("noop", func(p *Processor, m *Memory) {
	})
	i.Movement("jump", func(p *Processor, m *Memory) {
		p.Jump(m.Get(0))
	})
	i.Movement("jumpIfZero", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(0)) == 0 {
			p.Jump(m.Get(1))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfNotZero", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(0)) != 0 {
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
		if p.Registers.Get((*m).Get(0)) > p.Registers.Get((*m).Get(1)) {
			p.Jump(m.Get(2))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfLessThan", func(p *Processor, m *Memory) {
		if p.Registers.Get((*m).Get(0)) < p.Registers.Get((*m).Get(1)) {
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
	})
	i.Instruction("swap", func(p *Processor, m *Memory) {
		x := p.Registers.Get((*m).Get(0))
		p.Registers.Set((*m).Get(0), (*m).Get(1))
		p.Registers.Set((*m).Get(1), x)
	})
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
	})
	i.Instruction("add", func(p *Processor, m *Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))+p.Registers.Get((*m).Get(1)))
	})
	i.Instruction("subtract", func(p *Processor, m *Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))-p.Registers.Get((*m).Get(1)))
	})
}

func FloatMathInstructions(is *InstructionSet) {
	is.Instruction("floatAdd", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), p.FloatHeap.Get(m.Get(0))+p.FloatHeap.Get(m.Get(1)))
	})
	is.Instruction("floatSubtract", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), p.FloatHeap.Get(m.Get(0))-p.FloatHeap.Get(m.Get(1)))
	})
	is.Instruction("floatMultiply", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), p.FloatHeap.Get(m.Get(0))*p.FloatHeap.Get(m.Get(1)))
	})
	is.Instruction("floatDivide", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), p.FloatHeap.Get(m.Get(0))/p.FloatHeap.Get(m.Get(1)))
	})
	is.Instruction("floatSet", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(0), float64(m.Get(1))*0.0000001)
	})
	is.Instruction("floatCopy", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), p.FloatHeap.Get(m.Get(0)))
	})
	is.Instruction("floatAbs", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Abs(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatAcos", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Acos(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatAcosh", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Acosh(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatAsin", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Asin(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatAsinh", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Asinh(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatCbrt", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Cbrt(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatCeil", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Ceil(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatCos", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Cos(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatDim", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Dim(p.FloatHeap.Get(m.Get(0)), math.Abs(p.FloatHeap.Get(m.Get(1)))))
	})
	is.Instruction("floatErf", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Erf(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatExp", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Exp(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatExp2", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Exp2(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatExpm1", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Abs(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatFloor", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Floor(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatGamma", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Gamma(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatHypot", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Hypot(p.FloatHeap.Get(m.Get(0)), math.Abs(p.FloatHeap.Get(m.Get(1)))))
	})
	is.Instruction("floatJ0", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.J0(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatJ1", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.J1(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatLog", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Log(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatLog10", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Log10(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatLog1p", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Log1p(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatLog2", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Log2(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatLogb", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Logb(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatMax", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Max(p.FloatHeap.Get(m.Get(0)), p.FloatHeap.Get(m.Get(1))))
	})
	is.Instruction("floatMin", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Min(p.FloatHeap.Get(m.Get(0)), p.FloatHeap.Get(m.Get(1))))
	})
	is.Instruction("floatMod", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Mod(p.FloatHeap.Get(m.Get(0)), p.FloatHeap.Get(m.Get(1))))
	})
	is.Instruction("floatNextAfter", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Nextafter(p.FloatHeap.Get(m.Get(0)), p.FloatHeap.Get(m.Get(1))))
	})
	is.Instruction("floatPow", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Pow(p.FloatHeap.Get(m.Get(0)), p.FloatHeap.Get(m.Get(1))))
	})
	is.Instruction("floatRemainder", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(2), math.Remainder(p.FloatHeap.Get(m.Get(0)), p.FloatHeap.Get(m.Get(1))))
	})
	is.Instruction("floatSin", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Sin(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatSinh", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Sinh(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatSqrt", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Sqrt(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatTan", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Tan(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatTanh", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Tanh(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatTrunc", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Trunc(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatY0", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Y0(p.FloatHeap.Get(m.Get(0))))
	})
	is.Instruction("floatY1", func(p *Processor, m *Memory) {
		p.FloatHeap.Set(m.Get(1), math.Y1(p.FloatHeap.Get(m.Get(0))))
	})
}
