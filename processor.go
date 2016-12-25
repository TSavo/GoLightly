package govirtual

import (
	"fmt"
	"runtime"
	"time"
	"golang.org/x/tools/go/gcimporter15/testdata"
)

//A Processor contains an instruction set, some memory, a Program, an Instruction Pointer
//(an index for where in the program it is), and a Termination Condition which
//let's it know when to stop.
type Processor struct {
	Id int
	*InstructionSet
	Registers            Memory
	CallStack            Memory
	Heap                 *Memory
	Stack                Memory
	InstructionPointer   int
	Program              *Program
	StartTime            int
	TerminationCondition *TerminationCondition
	cost                 int
}

//A Processor's cost is the number of operations it's executed + the program length + the stack length + the call stack length
func (p *Processor) Cost() int {
	progLen := p.Len()
	cost := p.cost
	return cost + progLen + p.Stack.Len() + p.CallStack.Len()

}

func (p *Processor) String() string {
	return fmt.Sprintf("Processor [Registers: %v, Heap: %v, Instruction Pointer: %d, Cost: %d]",
		p.Registers,
		//p.CallStack,
		p.Heap,
		//p.Stack,
		p.InstructionPointer,
		p.Cost())
}

func (t *Processor) SetInstructionPointer(jump int) {
	t.InstructionPointer = jump
	if t.InstructionPointer < 0 {
		t.InstructionPointer = 0
	}
	t.InstructionPointer = t.InstructionPointer % t.Program.Len()
}

//Change the instruction pointer to the specified index
func (t *Processor) Jump(jump interface{}) {
	switch j := jump.(type) {
	case int:
		t.SetInstructionPointer(j)
	case float64:
		t.SetInstructionPointer(int(j))
	case string:
		t.JumpLabel(j)
	case Pointer:
		t.Jump(j.Get())
	default:
		panic(fmt.Sprintf("Don't know how to jump to %v", j))
	}
}

func (t *Processor) JumpLabel(label string) {
	defer func() {
		if recover() != nil {
			t.InstructionPointer++
		}
	}()
	t.InstructionPointer = t.Program.Labels()[label][0]
}

//Push the current instruction pointer onto the call stack, and jump to the location in the program specified.
func (t *Processor) Call(location interface{}) {
	t.CallStack.Push(&Literal{t.InstructionPointer})
	t.Jump(location)
}

//Pop the top of the call stack, and jump to that value.
func (t *Processor) Return() {
	if t.CallStack.Len() > 0 {
		x, _ := t.CallStack.Pop()
		t.InstructionPointer = Cardinalize(x)
	}
	t.InstructionPointer++
}

//Create a new Processor with a memory of length 'registers', an instruction set, a heap, and a termination condition
func NewProcessor(id int, registers int, instructions *InstructionSet, heap *Memory, stop *TerminationCondition) *Processor {
	p := new(Processor)
	p.Id = id
	p.TerminationCondition = stop
	p.Registers = make(Memory, registers)
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	p.Stack = make(Memory, 0)
	p.CallStack = make(Memory, 0)
	p.Heap = heap
	return p
}

//Load a program into the processor. This has the side effect of setting the instruction pointer to 0.
func (p *Processor) LoadProgram(program *Program) {
	p.Program = program.Clone()
	p.InstructionPointer = 0
}

//Compile a program from a string and load it into the processor.
func (p *Processor) CompileAndLoad(prog string) {
	p.LoadProgram(p.InstructionSet.CompileProgram(prog, p.Heap))
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
	x := this.InstructionPointer
	x = x % this.Program.Len()
	opcode := this.Program.Get(x)
	this.InstructionPointer += 1
	return opcode
}

//Start running the program. This won't return until the TerminationCondition.ShouldTerminate() returns true.
func (p *Processor) Run() {
	p.StartTime = int(time.Now().UnixNano())
	for {
		runtime.Gosched()
		if (*p.TerminationCondition).ShouldTerminate(p) {
			return
		}
		if p.Program.Len() == 0 {
			return
		}
		p.Execute()
	}
}
