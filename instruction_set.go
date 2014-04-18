package govirtual

import (
	"fmt"
	"strconv"
	"strings"
)

type Operation struct {
	Instruction *Instruction
	Data        *Memory
}

func (o Operation) String() string {
	return fmt.Sprintf("%v %v,%v,%v", o.Instruction.Name, o.Data.Get(0), o.Data.Get(1), o.Data.Get(2))
}

func (o Operation) Similar(p Operation) bool {
	return o.Instruction == p.Instruction
}

func (o *Operation) Decode() *Memory {
	return &Memory{o.Instruction.id, o.Data.Get(0), o.Data.Get(1), o.Data.Get(2)}
}

type Assembler interface {
	Assemble(name string, data *Memory) Operation
}

type Instruction struct {
	id       int
	Name     string
	Movement int
	Closure  func(*Processor, *Memory)
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
func (i *InstructionSet) Define(name string, movement int, closure func(*Processor, *Memory)) {
	id := i.Len()
	(*i)[id] = &Instruction{id: id, Name: name, Movement: movement, Closure: closure}
}
func (i *InstructionSet) Movement(name string, closure func(*Processor, *Memory)) {
	i.Define(name, 0, closure)
}
func (i *InstructionSet) Instruction(name string, closure func(*Processor, *Memory)) {
	i.Define(name, 1, closure)
}

func (i *InstructionSet) Assemble(id int, data *Memory) *Operation {
	if op, ok := (*i)[id]; ok {
		return &Operation{Instruction: op, Data: data}
	}
	panic("No such Instruction")
}

func (i *InstructionSet) Encode(m *Memory) *Operation {
	return i.Assemble(m.Get(0)%i.Len(), &Memory{m.Get(1), m.Get(2), m.Get(3)})
}

func (i *InstructionSet) CompileMemory(name string, mem *Memory) *Operation {
	for x, n := range *i {
		if n.Name == name {
			return i.Encode(&Memory{x, mem.Get(0), mem.Get(1), mem.Get(2)})
		}
	}
	panic("No such instruction")
}

func (i *InstructionSet) Compile(name string, args ...int) (o *Operation) {
	switch len(args) {
	case 0:
		o = i.CompileMemory(name, &Memory{0, 0, 0})
	case 1:
		o = i.CompileMemory(name, &Memory{args[0], 0, 0})
	case 2:
		o = i.CompileMemory(name, &Memory{args[0], args[1], 0})
	case 3:
		o = i.CompileMemory(name, &Memory{args[0], args[1], args[2]})
	default:
		panic("Arguments > 3 is not supported")
	}
	return
}

func (i *InstructionSet) CompileProgram(s string) *Program {
	p := make(Program, 0)
	for _, x := range strings.Split(s, "\n") {
		o := strings.Split(x, " ")
		if len(strings.TrimSpace(o[0])) == 0 {
			continue
		}
		if len(o) == 1 {
			p = append(p, i.Compile(o[0]))
		} else if len(o) == 2 {
			c := strings.Split(o[1], ",")
			args := make([]int, len(c))
			for x, _ := range c {
				args[x], _ = strconv.Atoi(strings.TrimSpace(c[x]))
			}
			p = append(p, i.Compile(o[0], args...))
		} else {
			panic("Don't know how to compile" + x)
		}
	}
	return &p
}
