package main

import (
	"fmt"
	"github.com/tsavo/golightly/vm"
)

func DefineInstructions() (i *vm.InstructionSet) {
	i = vm.NewInstructionSet()
	i.Operator("noop", func(p *vm.ProcessorCore, m *vm.Memory) {
	})
	i.Movement("halt", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Running = false
	})
	i.Movement("jmp", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Jump(p.Registers.Get((*m).Get(0)))
	})
	i.Movement("jmpz", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) == 0 {
			p.Jump(p.Registers.Get((*m).Get(1)))
		}
	})
	i.Movement("jmpnz", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers[(*m).Get(0)] != 0 {
			p.Jump(p.Registers[(*m).Get(1)])
		}
	})
	i.Movement("call", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Call(p.Registers.Get((*m).Get(0)))
	})
	i.Movement("ret", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Return()
	})
	i.Operator("push", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.MemorySegment.Push(p.Registers.Get((*m).Get(0)))
	})
	i.Operator("pop", func(p *vm.ProcessorCore, m *vm.Memory) {
		x, _ := p.MemorySegment.Pop()
		p.Registers.Set((*m).Get(0), x)
	})
	i.Operator("mov", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), (*m).Get(1))
	})
	i.Operator("inc", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Increment((*m).Get(0))
	})
	i.Operator("dec", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Decrement((*m).Get(0))
	})
	return
}

func main() {
	fmt.Println("ok")
	x := &vm.ProcessorCore{}
	instructionSet := DefineInstructions()
	x.Init(4, instructionSet)
	fmt.Println(x)
	p := &vm.Program{
		instructionSet.Assemble("mov", &vm.Memory{0, 1}),
		instructionSet.Assemble("halt", nil),
	}
	x.LoadProgram(p)
	x.Run()
	fmt.Println(x)
}
