package govirtual

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	Evaluate(Processor, ...Pointer)
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

type Assembler interface {
	Assemble(name string, data *Memory) Operation
}

type Argument struct {
	Name, Type string
}

type Closure func(*Processor, ...Pointer) Memory

type Instruction struct {
	id        int
	Name      string
	Movement  int
	Closure   Closure
	Arguments []Argument
	Infix     bool
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

func (i *InstructionSet) Define(name string, movement int, infix bool, closure Closure, args ...Argument) {
	id := i.Len()
	(*i)[id] = &Instruction{id: id, Name: name, Movement: movement, Closure: closure, Infix: infix, Arguments: args}
}
func (i *InstructionSet) Movement(name string, closure Closure, args ...Argument) {
	i.Define(name, 0, false, closure, args...)
}
func (i *InstructionSet) Instruction(name string, closure Closure, args ...Argument) {
	i.Define(name, 1, false, closure, args...)
}

func (i *InstructionSet) Infix(name string, closure Closure, left Argument, right Argument) {
	i.Define(name, 1, true, closure, left, right)
}

func (i *InstructionSet) Assemble(id int, args ...Pointer) *Operation {
	if op, ok := (*i)[id]; ok {
		a := Memory(args)
		return &Operation{Instruction: op, Data: &a}
	}
	panic("No such Instruction")
}

func (i *InstructionSet) Compile(name string, args ...Pointer) *Operation {
	for x, n := range *i {
		if n.Name == name {
			return i.Assemble(x, args...)
		}
	}
	panic("No such instruction " + name)
}

func (i *InstructionSet) CompileLabel(label string) *Operation {
	o := i.Compile("noop")
	o.Label = label
	return o
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

func Coherse(arg string, heap *Memory) Pointer {
	if strings.HasPrefix(arg, ":") {
		return &Literal{arg}
	} else if strings.HasPrefix(arg, "#") {
		index, _ := strconv.Atoi(arg[1:])
		return &MemoryPointer{heap, index, "#"}
	} else if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
		return &Literal{arg[1:len(arg)-1]}
	} else {
		if strings.Contains(arg, ".") {
			n, _ := strconv.ParseFloat(arg, 64)
			return &Literal{n}
		} else {
			n, _ := strconv.Atoi(arg)
			return &Literal{n}
		}
	}
}

func (i *InstructionSet) CompileProgram(s string, heap *Memory) *Program {
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
			args := make([]Pointer, len(c))
			for x, arg := range c {
				args[x] = Coherse(arg, heap)
			}
			p.Append(i.Compile(o[0], args...))
		} else {
			panic("Don't know how to compile" + line)
		}
	}
	return p
}
