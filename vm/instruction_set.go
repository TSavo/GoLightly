package vm

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
	return fmt.Sprintf("%v %v, %v\n", o.Instruction, o.Data.Get(0), o.Data.Get(1))
}

func (o Operation) Similar(p Operation) bool {
	return o.Instruction == p.Instruction
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

func (i *InstructionSet) Compile(name string, args ...int) (o *Operation) {
	switch len(args) {
	case 0:
		o = i.CompileMemory(name, &Memory{0, 0})
	case 1:
		o = i.CompileMemory(name, &Memory{args[0], 0})
	case 2:
		o = i.CompileMemory(name, &Memory{args[0], args[1]})
	default:
		panic("Arguments > 2 is not supported")
	}
	return
}

func (i *InstructionSet) Decompile(op *Operation) string {
	return op.Instruction.Name + " " + strconv.Itoa(op.Data.Get(0)) + ", " + strconv.Itoa(op.Data.Get(1))
}

func (i *InstructionSet) DecompileProgram(p *Program) (prog string) {
	prog = ""
	for _, v := range *p {
		prog += i.Decompile(v) + "\n"
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
			if len(c) == 1 {
				arg0, _ := strconv.Atoi(strings.TrimSpace(c[0]))
				p = append(p, i.Compile(o[0], arg0))
			} else {
				arg0, _ := strconv.Atoi(strings.TrimSpace(c[0]))
				arg1, _ := strconv.Atoi(strings.TrimSpace(c[1]))
				p = append(p, i.Compile(o[0], arg0, arg1))
			}
		} else {
			c := strings.Split(o[1], ",")
			arg0, _ := strconv.Atoi(strings.TrimSpace(c[0]))
			arg1, _ := strconv.Atoi(strings.TrimSpace(o[2]))
			p = append(p, i.Compile(o[0], arg0, arg1))
		}
	}
	return &p
}
