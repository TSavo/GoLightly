package vm

import (
	"fmt"
	"runtime"
	"time"
)

type ProcessorCore struct {
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

func NewProcessorCore(registers int, instructions *InstructionSet, heap *Memory, stop *TerminationCondition) *ProcessorCore {
	p := new(ProcessorCore)
	p.TerminationCondition = stop
	p.Registers = make(Memory, registers)
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	return p
}

func (p *ProcessorCore) LoadProgram(program *Program) {
	pr := make(Program, len(*program))
	p.Program = &pr
	copy(*p.Program, *program)
}

func (p *ProcessorCore) CompileAndLoad(prog string) {
	p.LoadProgram(p.InstructionSet.CompileProgram(prog))
}

func (p *ProcessorCore) GetProgramString() string {
	return p.Program.Decompile()
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
	p.StartTime = time.Now().UnixNano()
	for {
		runtime.Gosched()
		if (*p.TerminationCondition).ShouldTerminate(p) {
			return
		}
		p.Execute()
	}
}
