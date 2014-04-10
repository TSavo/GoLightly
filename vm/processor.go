package vm

import (
	"fmt"
	"runtime"
	"time"
)

type Processor struct {
	*InstructionSet
	Registers            Memory
	CallStack            Memory
	Heap                 *Memory
	Stack                Memory
	InstructionPointer   int
	cost                 int64
	Program              *Program
	StartTime            int64
	TerminationCondition *TerminationCondition
}

func (p *Processor) Cost() int64 {
	progLen := int64(len(*p.Program))
	cost := p.cost
	return cost + progLen + int64(p.Stack.Len()+p.CallStack.Len())

}

func (p *Processor) String() string {
	return fmt.Sprintf("Processor [Registers: %v, Heap: %v, Instruction Pointer: %d Cost: %d]",
		p.Registers,
		//p.CallStack,
		p.Heap,
		//p.Stack,
		p.InstructionPointer,
		p.Cost())
}

func (t *Processor) Call(location int) {
	t.CallStack.Push(t.InstructionPointer)
	t.Jump(location)
}

func (t *Processor) Return() {
	if t.CallStack.Len() > 0 {
		t.InstructionPointer, _ = t.CallStack.Pop()
	}
	t.InstructionPointer++

}

func (t *Processor) Jump(jump int) {
	t.InstructionPointer = jump
	if t.InstructionPointer < 0 {
		t.InstructionPointer = 0
	}
	t.InstructionPointer = t.InstructionPointer % len(*t.Program)
}

func NewProcessor(registers int, instructions *InstructionSet, heap *Memory, stop *TerminationCondition) *Processor {
	p := new(Processor)
	p.TerminationCondition = stop
	p.Registers = make(Memory, registers)
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	p.Heap = heap
	return p
}

func (p *Processor) LoadProgram(program *Program) {
	pr := make(Program, len(*program))
	p.Program = &pr
	copy(*p.Program, *program)
}

func (p *Processor) CompileAndLoad(prog string) {
	p.LoadProgram(p.InstructionSet.CompileProgram(prog))
}

func (p *Processor) GetProgramString() string {
	return p.Program.Decompile()
}

func (p *Processor) Reset() {
	p.Registers.Zero()
	p.CallStack.Reallocate(0)
	p.InstructionPointer = 0
	p.cost = 0
}

func (p *Processor) Execute() {
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

func (p *Processor) Run() {
	p.StartTime = time.Now().UnixNano()
	for {
		runtime.Gosched()
		if (*p.TerminationCondition).ShouldTerminate(p) {
			return
		}
		p.Execute()
	}
}
