//	TODO:	storing and retrieving pointers to memory buffers
//	TODO:	cloning should create a comms channel by which the parent and child cores can communicate
//	TODO:	should always have stdin, stdout and stderr channels

package vm

import "fmt"

type ProcessorCore struct {
	*InstructionSet
	IOController
	Running            bool
	Registers          Memory
	CallStack          Memory
	MemorySegment      Memory
	InstructionPointer int
	Program
}

func (p ProcessorCore) String() string {
	return fmt.Sprintf("ProcessorCore [Running: %t, Registers: %v, Call Stack: %v, Memory Segment: %v, Instruction Pointer: %d]",
	p.Running,
	p.Registers,
	p.CallStack,
	p.MemorySegment,
	p.InstructionPointer)
}

func (t *ProcessorCore) Call(location int) {
	fmt.Printf("Jumping to %d\n", location)
	t.CallStack = append(t.CallStack, t.InstructionPointer)
	t.InstructionPointer = location
}

func (t *ProcessorCore) Return() {
	fmt.Println("returning")
	if t.CallStack.Len() > 0 {
		t.InstructionPointer, _ = t.CallStack.Pop()
		t.InstructionPointer++
	} else {
		fmt.Println("Can't return with no return stack.")
		panic(t)
	}
}

func (t *ProcessorCore) Jump(jump int) {
	t.InstructionPointer = jump
}

func (p *ProcessorCore) Init(registers int, instructions *InstructionSet) {
	p.Registers = make(Memory, registers)

	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
}

//	Make a copy of the current processor, binding it to the current processor with
//	the supplied io channel
func (p *ProcessorCore) Clone(c chan []int) (q *ProcessorCore, i int) {
	q = new(ProcessorCore)
	q.Init(len(p.Registers), p.InstructionSet)
	q.IOController = append(q.IOController, c)
	p.IOController = append(p.IOController, c)
	i = len(p.IOController) - 1
	return
}

func (p *ProcessorCore) LoadProgram(program *Program) {
	p.Program = make(Program, len(*program))
	copy(p.Program, *program)
	p.InstructionPointer = 0
}
func (p *ProcessorCore) ResetState() {
	p.Registers.Zero()
	p.CallStack.Reallocate(0)
	p.InstructionPointer = 0
}
func (p *ProcessorCore) Execute() {
	o := p.Program[p.InstructionPointer%len(p.Program)]
	fmt.Println(o)
	o.Execute(p)
	p.InstructionPointer += o.Instruction.Movement
}
func (p *ProcessorCore) Run() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Panic in execution detected.")
			fmt.Println(x)
			p.Running = false
		}
	}()
	p.Running = true
	for p.Running {
		fmt.Println("running")
		p.Execute()
	}
	fmt.Println("done")
}
