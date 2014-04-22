package govirtual

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	Evaluate(Processor, ... Pointer) 
}

type Operation struct {
	Instruction *Instruction
	Data        *Memory
	Label       string
	Infix       bool
}

func (o Operation) String() string {
	if len(o.Label) > 0 {
		return o.Label
	}
	if o.Infix {
		return fmt.Sprintf("%v %v %v", o.Data.Get(0), o.Instruction.Name, o.Data.Get(1))
	}
	return fmt.Sprintf("%v(%v)", o.Instruction.Name, o.Data.Get(0), o.Data.Get(1), o.Data.Get(2))
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

type Argument struct {
	Name, Type string
}

type Instruction struct {
	id        int
	Name      string
	Movement  int
	Closure   func(*Processor, *Memory)
	Arguments []Argument
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

func (i *InstructionSet) Define(name string, movement int, closure func(*Processor, *Memory), format func(*Memory) string) {
	id := i.Len()
	(*i)[id] = &Instruction{id: id, Name: name, Movement: movement, Closure: closure, Format: format}
}
func (i *InstructionSet) Movement(name string, closure func(*Processor, *Memory), format func(*Memory) string) {
	i.Define(name, 0, closure, format)
}
func (i *InstructionSet) Instruction(name string, closure func(*Processor, *Memory), format func(*Memory) string) {
	i.Define(name, 1, closure, format)
}

func (i *InstructionSet) Assemble(id int, data *Memory) *Operation {
	if op, ok := (*i)[id]; ok {
		return &Operation{Instruction: op, Data: data}
	}
	panic("No such Instruction")
}

func (i *InstructionSet) Encode(m *Memory) *Operation {
	return i.Assemble(m.GetCardinal(0)%i.Len(), &Memory{m.Get(1), m.Get(2), m.Get(3)})
}

func (i *InstructionSet) CompileMemory(name string, mem *Memory) *Operation {
	for x, n := range *i {
		if n.Name == name {
			return i.Encode(&Memory{x, mem.Get(0), mem.Get(1), mem.Get(2)})
		}
	}
	panic("No such instruction " + name)
}

func (i *InstructionSet) CompileLabel(label string) *Operation {
	o := i.CompileMemory("noop", &Memory{0, 0, 0})
	o.Label = label
	return o
}

func (i *InstructionSet) Compile(name string, args ...interface{}) (o *Operation) {
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

func UnlabelProgram(program string) (string, map[string]int) {
	prog, labels := UnlabelProgramRecurse(strings.Split(program, "\n"), make(map[string]int))
	return strings.Join(prog, "\n"), labels
}

func UnlabelProgramRecurse(program []string, labels map[string]int) ([]string, map[string]int) {
	for x, line := range program {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ":") {
			labels[line] = x
			return UnlabelProgramRecurse(append(program[:x], program[x+1:]...), labels)
		}
	}
	return program, labels
}

func Coherse(arg string) interface{} {
	if strings.HasPrefix(arg, ":") {
		return arg
	} else {
		if strings.Contains(arg, ".") {
			n, _ := strconv.ParseFloat(arg, 64)
			return n
		} else {
			n, _ := strconv.Atoi(arg)
			return n
		}
	}
}

func (i *InstructionSet) CompileProgram(s string) *Program {
	p := NewProgram(0)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, ":") {
			p.Append(i.CompileLabel(line))
			continue
		}
		o := strings.Split(line, " ")

		if len(o) == 1 {
			p.Append(i.Compile(o[0]))
		} else if len(o) == 2 {
			c := strings.Split(o[1], ",")
			args := make([]interface{}, len(c))
			for x, arg := range c {
				args[x] = Coherse(arg)
			}
			p.Append(i.Compile(o[0], args...))
		} else {
			panic("Don't know how to compile" + line)
		}
	}
	return p
}
