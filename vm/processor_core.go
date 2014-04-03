package vm

import (
	"fmt"
	"runtime"
)

const (
	MAX_COST = 300
	MIN_COST = 25
)

type ProcessorCore struct {
	*InstructionSet
	Registers          Memory
	CallStack          Memory
	Heap               *Memory
	Stack              Memory
	InstructionPointer int
	cost               int64
	Program            *Program
	ControlChan        chan bool
}

func (p *ProcessorCore) Cost() int64 {
	progLen := int64(len(*p.Program))
	cost := p.cost
	return cost + progLen + int64(p.Stack.Len()+p.CallStack.Len())

}

func (p *ProcessorCore) String() string {
	return fmt.Sprintf("ProcessorCore [Registers: %v, Heap: %v, Instruction Pointer: %d Cost: %d]",
		p.Registers,
		//p.CallStack,
		p.Heap,
		//p.Stack,
		p.InstructionPointer,
		p.Cost())
}

func (t *ProcessorCore) Call(location int) {
	t.CallStack.Push(t.InstructionPointer)
	t.Jump(location)
}

func (t *ProcessorCore) Return() {
	if t.CallStack.Len() > 0 {
		t.InstructionPointer, _ = t.CallStack.Pop()
	}
	t.InstructionPointer++

}

func (t *ProcessorCore) Jump(jump int) {
	t.InstructionPointer = jump
	if t.InstructionPointer < 0 {
		t.InstructionPointer = 0
	}
	t.InstructionPointer = t.InstructionPointer % len(*t.Program)
}

func NewProcessorCore(registers int, instructions *InstructionSet, heap *Memory) *ProcessorCore {
	p := new(ProcessorCore)
	p.Registers = make(Memory, registers)
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	p.ControlChan = make(chan bool, 1)
	return p
}

func (p *ProcessorCore) LoadProgram(program *Program) {
	pr := make(Program, len(*program))
	p.Program = &pr
	copy(*p.Program, *program)
	for i, x := range *program {
		(*program)[i] = p.InstructionSet.Encode(p.InstructionSet.Decode(x))
	}
}

func (p *ProcessorCore) Reset() {
	p.Registers.Zero()
	p.CallStack.Reallocate(0)
	p.InstructionPointer = 0
	p.cost = 0
}

func (p *ProcessorCore) Execute() {
	if len(*p.Program) == 0 {
		return
	}
	x := p.InstructionPointer
	if x >= len(*p.Program) {
		x = x % len(*p.Program)
	}
	opcode := (*p.Program)[x]
	opcode.Instruction.Closure(p, opcode.Data)
	p.cost++
	p.InstructionPointer += opcode.Instruction.Movement
}

func (p *ProcessorCore) Run() {
outer:
	for {
		runtime.Gosched()
		select {
		case <-p.ControlChan:
			break outer
		default:
			p.Execute()
		}
	}
}
