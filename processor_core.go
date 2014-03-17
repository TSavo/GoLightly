//	TODO:	storing and retrieving pointers to memory buffers
//	TODO:	cloning should create a comms channel by which the parent and child cores can communicate
//	TODO:	should always have stdin, stdout and stderr channels

package golightly

import "time"
import "fmt"

type IsExecutable interface {
	Map(f interface{}) interface{}
	Reduce(f interface{}) interface{}
}

type Thread struct {
	Running bool
	R       []int
	M       []int
	CS      []int
	DS      []int
	PC      int
	Program
}

func (t *Thread) I() OpCode     { return t.Program[t.PC] }
func (t *Thread) ValidPC() bool { return t.PC > -1 && t.PC < len(t.Program) }
func (t *Thread) Call(location int) {
	fmt.Printf("Jumping to %d\n", location)
	t.CS = append(t.CS, t.PC)
	t.PC = location
}
func (t *Thread) Return() {
    fmt.Println("returning")
	if len(t.CS) > 0 {
		t.PC = t.CS[len(t.CS)-1] + 1
		t.CS = t.CS[:len(t.CS)-1]
	} else {
		fmt.Println("Can't return with no return stack.")
		panic(t)
	}
}

type ProcessorCore struct {
	*InstructionSet
	IOController
	Thread
}

func (p *ProcessorCore) Init(registers int, instructions *InstructionSet) {
	data := make([]int, registers)
	p.Thread = Thread{R: data}
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
		p.InstructionSet.Init()
		p.DefineInstructions()
	} else {
		p.InstructionSet = instructions
	}
}

//	Make a copy of the current processor, binding it to the current processor with
//	the supplied io channel
func (p *ProcessorCore) Clone(c chan []int) (q *ProcessorCore, i int) {
	q = new(ProcessorCore)
	q.Init(len(p.R), p.InstructionSet)
	q.IOController = append(q.IOController, c)
	p.IOController = append(p.IOController, c)
	i = len(p.IOController) - 1
	return
}
func (p *ProcessorCore) DefineInstructions() {
	p.InstructionSet.Operator("noop", func() {})                                      //	NOOP
	p.InstructionSet.Operator("nsleep", func(n time.Duration) { time.Sleep(n) })      //	NSLEEP	n
	p.InstructionSet.Operator("sleep", func(n time.Duration) { time.Sleep(n << 32) }) //	SLEEP	n
	p.InstructionSet.Movement("halt", func() { p.Running = false })                   //	HALT
	p.InstructionSet.Movement("jmp", func(n int) { p.PC = n })                        //	JMP		n
	p.InstructionSet.Movement("jmpz", func(o []int) {
		if p.R[o[0]] == 0 {
			p.PC = o[1]
		}
	}) //	JMPZ	r, n
	p.InstructionSet.Movement("jmpnz", func(o []int) {
		if p.R[o[0]] != 0 {
			p.PC += o[1]
		}
	}) //	JMPNZ	r, n
	p.InstructionSet.Movement("call", func(n int) { p.Call(n) })                   //	CALL	n
	p.InstructionSet.Movement("ret", func() { p.Return() })                        //	RET
	p.InstructionSet.Operator("push", func(r int) { p.DS = append(p.DS, p.R[r]) }) //	PUSH	r
	p.InstructionSet.Operator("pop", func(r int) {
		p.R[r] = p.DS[len(p.DS)-1]
		p.DS = p.DS[:len(p.DS)-1]
	}) //	POP		r
	p.InstructionSet.Operator("cld", func(o []int) { p.R[o[0]] = o[1] })               //	CLD		r, v
	p.InstructionSet.Operator("send", func(c int) { p.IOController.Send(c, p.M) })     //	SEND	c
	p.InstructionSet.Operator("recv", func(c int) { p.M = p.IOController.Receive(c) }) //	RECV	c
	p.InstructionSet.Operator("inc", func(r int) { p.R[r] = p.R[r] + 1 })                 //	INC		r
	p.InstructionSet.Operator("dec", func(r int) { p.R[r] = p.R[r] - 1 })                 //	DEC		r
}
func (p *ProcessorCore) LoadProgram(program Program) {
	p.Program = make(Program, len(program))
	copy(p.Program, program)
	//slices.ClearAll(p.R)
	p.M = nil
	p.PC = 0
}
func (p *ProcessorCore) ResetState() {
	//slices.ClearAll(p.R)
	p.M = nil
	p.PC = 0
}
func (p *ProcessorCore) Execute() {
	o := p.Program[p.PC]
	fmt.Println(o)
	switch data := o.data.(type) {
	case int:
		p.ops[o.code].(func(int))(data)
	case []int:
		p.ops[o.code].(func([]int))(data)
	case nil:
		p.ops[o.code].(func())()
    default:
        panic("No case for op type")
	}
	p.PC += o.movement
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
