package govirtual

import (
	"fmt"
)

type InstructionPointer int

//A Processor contains an instruction set, some memory, a Program, an Instruction Pointer
//(an index for where in the program it is), and a Termination Condition which
//let's it know when to stop.
type Processor struct {
	*InstructionSet
	*Memory
	*InstructionPointer
	*Program
	TerminationCondition *TerminationCondition
}

func (p *Processor) String() string {
	return fmt.Sprintf("Processor [Heap: %v, Instruction Pointer: %d]",
		p.Memory,
		p.InstructionPointer)
}

func (t *Processor) SetInstructionPointer(jump InstructionPointer) {
	*t.InstructionPointer = jump
	if *t.InstructionPointer < 0 {
		*t.InstructionPointer = 0
	}
	*t.InstructionPointer = InstructionPointer(int(*t.InstructionPointer) % t.Program.Len())
}

//Change the instruction pointer to the specified index
func (t *Processor) Jump(jump interface{}) {
	switch j := jump.(type) {
	case InstructionPointer:
		t.SetInstructionPointer(j)
	case string:
		t.JumpLabel(j)
	case Value:
		t.Jump(j.Get())
	default:
		panic(fmt.Sprintf("Don't know how to jump to %v", j))
	}
}

func (t *Processor) JumpLabel(label string) {
	defer func() {
		if recover() != nil {
			*t.InstructionPointer++
		}
	}()
	temp := InstructionPointer(t.Program.Labels()[label][0])
	t.InstructionPointer = &temp
}

//Create a new Processor with a memory of length 'registers', an instruction set, a heap, and a termination condition
func NewProcessor(id int, registers int, instructions *InstructionSet, memory *Memory, stop *TerminationCondition) *Processor {
	p := new(Processor)
	p.TerminationCondition = stop
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	p.Memory = memory
	return p
}

//Load a program into the processor. This has the side effect of setting the instruction pointer to 0.
func (p *Processor) LoadProgram(program *Program) {
	p.Program = program.Clone()
	*p.InstructionPointer = InstructionPointer(0)
}

//Compile a program from a string and load it into the processor.
func (p *Processor) CompileAndLoad(prog string) {
	p.LoadProgram(p.InstructionSet.CompileProgram(prog))
}

//Decompile the program on the processor into a string.
func (p *Processor) GetProgramString() string {
	return p.Program.Decompile()
}

//Execute the instruction currently pointed at by the Instruction Pointer
func (this *Processor) Execute() {
	if this.Program.Len() == 0 {
		return
	}
	operation := this.FetchOperation()
	operation.Instruction.Closure(this, operation.Data...)
}

func (this *Processor) FetchOperation() *Operation {
	x := int(*this.InstructionPointer)
	opcode := this.Program.Get(x)
	*this.InstructionPointer = InstructionPointer(x + 1)
	return opcode
}

//Start running the program. This won't return until the TerminationCondition.ShouldTerminate() returns true.
func (p *Processor) Run() {
	for {
		if (*p.TerminationCondition).ShouldTerminate(p) {
			return
		}
		if p.Program.Len() == 0 {
			return
		}
		p.Execute()
	}
}
