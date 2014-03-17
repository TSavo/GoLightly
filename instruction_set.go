//	TODO:	bytecode optimisation
//	TODO:	JIT compilation
//	TODO:	AOT compilation

package golightly

import (
	"fmt"
	"github.com/feyeleanor/slices"
	"reflect"
)

type OpCode struct {
	code     int
	movement int
	data     interface{}
}


func (o OpCode) String() string {
	return fmt.Sprintf("Code: %d Movement: %d Data: %v", o.code, o.movement, o.data)
}

func (o OpCode) Similar(p OpCode) bool {
	return o.code == p.code && o.movement == p.movement && reflect.TypeOf(o.data) == reflect.TypeOf(p.data)
}
func (o OpCode) Identical(p OpCode) bool {
	return reflect.DeepEqual(o, p)
}
func (o *OpCode) Replace(p *OpCode) {
	o.code = p.code
	o.movement = p.movement
	o.data = p.data
}

type Assembler interface {
	Assemble(name string, data interface{}) OpCode
}

type Instruction struct {
	op       int
	movement int
}

type InstructionSet struct {
	ops    slices.Slice
	tokens map[string]*Instruction
}

func (i *InstructionSet) Init() {
	i.tokens = make(map[string]*Instruction)
}
func (i *InstructionSet) Len() int {
	return i.ops.Len()
}
func (i *InstructionSet) Exists(name string) bool {
	_, ok := i.tokens[name]
	return ok
}
func (i *InstructionSet) Define(name string, movement int, closure interface{}) (successful bool) {
	if _, ok := i.tokens[name]; !ok {
		i.ops.Append(closure)
		i.tokens[name] = &Instruction{op: i.ops.Len() - 1, movement: movement}
		successful = true
	}
	return
}
func (i *InstructionSet) Movement(name string, data interface{}) bool {
	return i.Define(name, 0, data)
}
func (i *InstructionSet) Operator(name string, data interface{}) bool {
	return i.Define(name, 1, data)
}
func (i *InstructionSet) Instruction(name string) *Instruction {
	if op, ok := i.tokens[name]; ok {
		return op
	}
	return nil
}
func (i *InstructionSet) Assemble(name string, data interface{}) OpCode {
	if op := i.Instruction(name); op != nil {
		return OpCode{code: op.op, movement: op.movement, data: data}
	}
	panic(name)
}
func (i *InstructionSet) Invoke(o *OpCode) {
	switch data := o.data.(type) {
	case []int:
		i.ops.At(o.code).(func(o []int))(data)
	}
}
