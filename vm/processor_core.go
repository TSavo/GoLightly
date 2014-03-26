//	TODO:	storing and retrieving pointers to memory buffers
//	TODO:	cloning should create a comms channel by which the parent and child cores can communicate
//	TODO:	should always have stdin, stdout and stderr channels

package vm

import "fmt"

const (
	MAX_COST = 300
	MIN_COST = 25
)

type ProcessorCore struct {
	*InstructionSet
	IOController
	Running            bool
	Registers          Memory
	CallStack          Memory
	Heap               *Memory
	Stack              Memory
	InstructionPointer int
	cost               int
	ChanceOfMutation   float64
	OutChan            chan int
	Program
	finished chan *ProcessorCore
}

func (p *ProcessorCore) Cost() int {
	progLen := len(p.Program)
	cost := p.cost
	return cost + progLen + p.Stack.Len() + p.CallStack.Len()

}

func (p *ProcessorCore) String() string {
	return fmt.Sprintf("ProcessorCore [Running: %t, Registers: %v, Heap: %v, Instruction Pointer: %d Cost: %d]",
		p.Running,
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
	if t.InstructionPointer >= len(t.Program) {
		t.InstructionPointer = t.InstructionPointer % (len(t.Program) - 1)
	}
}

func (p *ProcessorCore) Init(registers int, instructions *InstructionSet, outChan chan int, finished chan *ProcessorCore) {
	p.Registers = make(Memory, registers)
	heap := make(Memory, 4)
	p.Heap = &heap
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	p.finished = finished
	p.OutChan = outChan
	p.ChanceOfMutation = 0.1
}


func (p *ProcessorCore) LoadProgram(program *Program) {
	p.Program = make(Program, len(*program))
	copy(p.Program, *program)
	p.InstructionPointer = 0
}
func (p *ProcessorCore) ResetState() {
	p.Registers.Zero()
	p.Heap.Zero()
	p.CallStack.Reallocate(0)
	p.InstructionPointer = 0
	p.cost = 0
}
func (p *ProcessorCore) Execute() {
	if len(p.Program) == 0 {
		p.InstructionPointer++
		p.cost++
		return
	}
	x := p.InstructionPointer
	if x >= len(p.Program) {
		x = x % len(p.Program)
	}
	o := p.Program[x]
	o.Execute(p)
	p.cost++
	p.InstructionPointer += o.Instruction.Movement
}
func (p *ProcessorCore) Run() {
	//	defer func() {
	//		if x := recover(); x != nil {
	//			p.Running = false
	//			panic(x)
	//		}
	//	}()
	p.Running = true
	defer func() {
		p.finished <- p
	}()
	for p.Running {
		p.Execute()
		if p.Cost() > MAX_COST {
			p.Running = false
		}
	}
}
