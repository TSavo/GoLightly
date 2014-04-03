package vm

import (
	"fmt"
	"reflect"
)

type Operation struct {
	Instruction *Instruction
	Data        *Memory
}

func (o Operation) String() string {
	return fmt.Sprintf("%v %v, %v\n", o.Instruction, o.Data.Get(0), o.Data.Get(1))
}

func (o Operation) Similar(p Operation) bool {
	return o.Instruction == p.Instruction
}

func (o Operation) Identical(p Operation) bool {
	return reflect.DeepEqual(o, p)
}

type Assembler interface {
	Assemble(name string, data *Memory) Operation
}

type Instruction struct {
	id       int
	Name     string
	Movement int
	Closure  func(*ProcessorCore, *Memory)
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s", i.Name)
}

type InstructionSet map[int]*Instruction

func NewInstructionSet() (i *InstructionSet) {
	x := make(InstructionSet)
	return &x
}

func (i *InstructionSet) Len() int {
	return len(*i)
}

func (i *InstructionSet) Exists(name int) bool {
	_, ok := (*i)[name]
	return ok
}
func (i *InstructionSet) Define(name string, movement int, closure func(*ProcessorCore, *Memory)) {
	id := i.Len()
	(*i)[id] = &Instruction{id: id, Name: name, Movement: movement, Closure: closure}
}
func (i *InstructionSet) Movement(name string, closure func(*ProcessorCore, *Memory)) {
	i.Define(name, 0, closure)
}
func (i *InstructionSet) Operator(name string, closure func(*ProcessorCore, *Memory)) {
	i.Define(name, 1, closure)
}
func (i *InstructionSet) Assemble(id int, data *Memory) Operation {
	if op, ok := (*i)[id]; ok {
		return Operation{Instruction: op, Data: data}
	}
	panic("No such Instruction")
}
func (i *InstructionSet) Encode(m *Memory) *Operation {
	return &Operation{(*i)[m.Get(0)%i.Len()], &Memory{m.Get(1), m.Get(2)}}
}
func (i *InstructionSet) Decode(o *Operation) *Memory {
	return &Memory{o.Instruction.id, o.Data.Get(0), o.Data.Get(1)}
}
func (i *InstructionSet) CompileMemory(name string, mem *Memory) *Operation {
	for x, n := range *i {
		if n.Name == name {
			return i.Encode(&Memory{x, mem.Get(0), mem.Get(1)})
		}
	}
	panic("No such instruction")
}

func (i *InstructionSet) Compile(name string, args ...int) *Operation {
	switch len(args) {
	case 0:
		return i.CompileMemory(name, &Memory{0,0})
	case 1:
		return i.CompileMemory(name, &Memory{args[0],0})
	case 2:
		return i.CompileMemory(name, &Memory{args[0], args[1]})
	default:
		panic("Arguments > 2 is not supported")
	}
}

func (i *InstructionSet) Decompile(op *Operation) string {
	s := op.Instruction.Name
	if(op.Data.Len() > 0){
		s += " " + string(op.Data.Get(0))
	}
	if(op.Data.Len() > 1){
		s += ", " + string(op.Data.Get(1))
	}
	return s
}

func (i *InstructionSet) DecompileProgram(p Program) (prog string) {
	prog = ""
	for _, v := range p {
		prog += i.Decompile(v) + "\n"
	}
	return

}
