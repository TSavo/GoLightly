//	TODO:	bytecode optimisation
//	TODO:	JIT compilation
//	TODO:	AOT compilation

package vm

import (
	"fmt"
	"reflect"
)

type OpCode struct {
	Instruction *Instruction
	Data *Memory
}


type Executable interface {
	Execute(code *ProcessorCore, data *Memory)
}

func (o *OpCode) Execute(core *ProcessorCore){
	o.Instruction.Closure(core, o.Data)
}

func (o OpCode) String() string {
	return fmt.Sprintf("Instruction: (%v) Data: (%v)", o.Instruction, o.Data)
}

func (o OpCode) Similar(p OpCode) bool {
	return o.Instruction == p.Instruction && reflect.TypeOf(o.Data) == reflect.TypeOf(p.Data)
}
func (o OpCode) Identical(p OpCode) bool {
	return reflect.DeepEqual(o, p)
}

type Assembler interface {
	Assemble(name string, data *Memory) OpCode
}

type Instruction struct {
	Code      string
	Movement int
	Closure func(*ProcessorCore, *Memory)
}

func (i Instruction) String() string {
	return fmt.Sprintf("Code: %s, Movement: %d", i.Code, i.Movement)
}

type InstructionSet struct {
	Instructions map[string]*Instruction
}

func NewInstructionSet() (i *InstructionSet){
	i = new(InstructionSet)
	i.Init()
	return
}

func (i *InstructionSet) Init() {
	i.Instructions = make(map[string]*Instruction)
}
func (i *InstructionSet) Len() int {
	return len(i.Instructions)
}
func (i *InstructionSet) Exists(name string) bool {
	_, ok := i.Instructions[name]
	return ok
}
func (i *InstructionSet) Define(name string, movement int, closure func(*ProcessorCore, *Memory)) (successful bool) {
	if _, ok := i.Instructions[name]; !ok {
		i.Instructions[name] = &Instruction{Code: name, Movement:movement, Closure:closure}
		successful = true
	}
	return
}
func (i *InstructionSet) Movement(name string, closure func(*ProcessorCore, *Memory)) bool {
	return i.Define(name, 0, closure)
}
func (i *InstructionSet) Operator(name string, closure func(*ProcessorCore, *Memory)) bool {
	return i.Define(name, 1, closure)
}
func (i *InstructionSet) Instruction(name string) *Instruction {
	if op, ok := i.Instructions[name]; ok {
		return op
	}
	return nil
}
func (i *InstructionSet) Assemble(name string, data *Memory) OpCode {
	if op := i.Instruction(name); op != nil {
		return OpCode{Instruction:op, Data: data}
	}
	panic("No such Instruction: " + name)
}