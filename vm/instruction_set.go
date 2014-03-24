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

func (o OpCode) Execute(core *ProcessorCore){
	o.Instruction.Closure(core, o.Data)
}

func (o OpCode) String() string {
	return fmt.Sprintf("%v %v, %v\n", o.Instruction, o.Data.Get(0), o.Data.Get(1))
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
	id int
	Name string
	Movement int
	Closure func(*ProcessorCore, *Memory)
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s", i.Name)
}

type InstructionSet struct {
	Instructions map[int]*Instruction
	opcode int
}

func NewInstructionSet() (i *InstructionSet){
	i = new(InstructionSet)
	i.Init()
	return
}

func (i *InstructionSet) Init() {
	i.Instructions = make(map[int]*Instruction)
}
func (i *InstructionSet) Len() int {
	return len(i.Instructions)
}
func (i *InstructionSet) Exists(name int) bool {
	_, ok := i.Instructions[name]
	return ok
}
func (i *InstructionSet) Define(name string, movement int, closure func(*ProcessorCore, *Memory)) {
	i.Instructions[i.opcode] = &Instruction{id: i.opcode, Name:name, Movement:movement, Closure:closure}
	i.opcode++
}
func (i *InstructionSet) Movement(name string, closure func(*ProcessorCore, *Memory)) {
	i.Define(name, 0, closure)
}
func (i *InstructionSet) Operator(name string, closure func(*ProcessorCore, *Memory)) {
	i.Define(name, 1, closure)
}
func (i *InstructionSet) Instruction(name int) *Instruction {
	if op, ok := i.Instructions[name]; ok {
		return op
	}
	return nil
}
func (i *InstructionSet) Assemble(name int, data *Memory) OpCode {
	if op := i.Instruction(name); op != nil {
		return OpCode{Instruction:op, Data: data}
	}
	panic("No such Instruction")
}
func (i *InstructionSet) Encode(m *Memory) *OpCode {
	return &OpCode{i.Instructions[m.Get(0) % len(i.Instructions)], &Memory{m.Get(1), m.Get(2)}}
}
func (i *InstructionSet) Decode(o *OpCode) *Memory {
	return &Memory{o.Instruction.id, o.Data.Get(0), o.Data.Get(1)}
}
