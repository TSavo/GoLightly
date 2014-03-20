package main

import (
	"fmt"
	"math/rand"
	"time"
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
		p.Jump(p.Registers.Get(0))
	})
	i.Movement("jmpz", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) == 0 {
			p.Jump(p.Registers.Get(0))
		}
	})
	i.Movement("jmpnz", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) != 0 {
			p.Jump(p.Registers[1])
		}
	})
	i.Movement("call", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Call(p.Registers.Get((*m).Get(0)))
	})
	i.Movement("ret", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Return()
	})
	i.Operator("set", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), (*m).Get(1))
	})
	i.Operator("store", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Heap.Set(p.Registers.Get(1), p.Registers.Get(0))
	})
	i.Operator("load", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set(0, p.Heap.Get(p.Registers.Get(1)))
	})
	i.Operator("swap", func(p *vm.ProcessorCore, m *vm.Memory) {
		x := p.Registers.Get((*m).Get(0))
		p.Registers.Set((*m).Get(0), (*m).Get(1))
		p.Registers.Set((*m).Get(1), x)
	})
	i.Operator("push", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Stack.Push(p.Registers.Get((*m).Get(0)))
	})
	i.Operator("pop", func(p *vm.ProcessorCore, m *vm.Memory) {
		if x, err := p.Stack.Pop(); !err {
			p.Registers.Set((*m).Get(0), x)
		}
	})
	i.Operator("inc", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Increment((*m).Get(0))
	})
	i.Operator("dec", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Decrement((*m).Get(0))
	})
	i.Operator("add", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) + p.Registers.Get((*m).Get(1)))
	})
	i.Operator("sub", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) - p.Registers.Get((*m).Get(1)))
	})
	i.Operator("mul", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) * p.Registers.Get((*m).Get(1)))
	})
	i.Operator("div", func(p *vm.ProcessorCore, m *vm.Memory) {
	    defer func(){
	      recover()
	    }()
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) / p.Registers.Get((*m).Get(1)))
	})
	i.Operator("mod", func(p *vm.ProcessorCore, m *vm.Memory) {
		defer func(){
	      recover()
	    }()
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) % p.Registers.Get((*m).Get(1)))
	})
	
	i.Operator("and", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) & p.Registers.Get((*m).Get(1)))
	})
	i.Operator("or", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) | p.Registers.Get((*m).Get(1)))
	})
	i.Operator("xor", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) ^ p.Registers.Get((*m).Get(1)))
	})
	i.Operator("xnot", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) &^ p.Registers.Get((*m).Get(1)))
	})
	
	i.Operator("lsft", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) << uint(p.Registers.Get((*m).Get(1))))
	})
	i.Operator("rsft", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0)) >> uint(p.Registers.Get((*m).Get(1))))
	})
	
	return
}

func main() {
	fmt.Println("ok")
	p := &vm.ProcessorCore{}
	instructionSet := DefineInstructions()
	p.Init(4, instructionSet)
	fmt.Println(p)
	
	pro := &vm.Program{
		instructionSet.Assemble("set", &vm.Memory{0, 1}),
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for x := 0; x < 100; x++ {
		*pro = append(*pro, instructionSet.Encode(&vm.Memory{r.Int(), r.Int(), r.Int()}))
	}
	p.LoadProgram(pro)
	p.Run()
	fmt.Println(p)
}
