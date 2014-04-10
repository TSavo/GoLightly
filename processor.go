package govirtual

import (
	"fmt"
	"runtime"
	"time"
)

//A Processor contains an instruction set, some memory, a Program, an Instruction Pointer 
//(an index for where in the program it is), and a Termination Condition which 
//let's it know when to stop.
type Processor struct {
	*InstructionSet
	Registers            Memory
	CallStack            Memory
	Heap                 *Memory
	Stack                Memory
	InstructionPointer   int
	Program              *Program
	StartTime            int64
	TerminationCondition *TerminationCondition
	cost                 int64
}


//A Processor's cost is the number of operations it's executed + the program length + the stack length + the call stack length
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

//Change the instruction pointer to the specified index
func (t *Processor) Jump(jump int) {
	t.InstructionPointer = jump
	if t.InstructionPointer < 0 {
		t.InstructionPointer = 0
	}
	t.InstructionPointer = t.InstructionPointer % len(*t.Program)
}

//Push the current instruction pointer onto the call stack, and jump to the location in the program specified.
func (t *Processor) Call(location int) {
	t.CallStack.Push(t.InstructionPointer)
	t.Jump(location)
}

//Pop the top of the call stack, and jump to that value.
func (t *Processor) Return() {
	if t.CallStack.Len() > 0 {
		t.InstructionPointer, _ = t.CallStack.Pop()
	}
	t.InstructionPointer++
}

//Create a new Processor with a memory of length 'registers', an instruction set, a heap, and a termination condition
func NewProcessor(registers int, instructions *InstructionSet, heap *Memory, stop *TerminationCondition) *Processor {
	p := new(Processor)
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
	pr := make(Program, len(*program))
	p.Program = &pr
	copy(*p.Program, *program)
	p.InstructionPointer = 0
}

//Compile a program from a string and load it into the processor.
func (p *Processor) CompileAndLoad(prog string) {
	p.LoadProgram(p.InstructionSet.CompileProgram(prog))
}

//Decompile the program on the processor into a string.
func (p *Processor) GetProgramString() string {
	return p.Program.Decompile()
}

//Zero the registers, reset the stack and the call stack, and set the cost and instruction pointer to 0
func (p *Processor) Reset() {
	p.Registers.Zero()
	p.CallStack.Reallocate(0)
	p.Stack.Reallocate(0)
	p.InstructionPointer = 0
	p.cost = 0
}

//Execute the instruction currently pointed at by the Instruction Pointer
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

//Start running the program. This won't return until the TerminationCondition.ShouldTerminate() returns true.
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
