package main

import (
	"fmt"
	"github.com/tsavo/golightly/vm"
	"math/rand"
	"sort"
	"time"
)

const (
	POPULATION_SIZE = 1000
	BEST_OF_BREED = 100
)

func DefineInstructions() (i *vm.InstructionSet) {
	i = vm.NewInstructionSet()
	i.Operator("noop", func(p *vm.ProcessorCore, m *vm.Memory) {
	})
	i.Movement("halt", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Running = false
	})
	i.Movement("jmp", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Jump(p.Registers.Get(1))
	})
	i.Movement("jmpz", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) == 0 {
			p.Jump(p.Registers.Get(1))
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
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))+p.Registers.Get((*m).Get(1)))
	})
	i.Operator("sub", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))-p.Registers.Get((*m).Get(1)))
	})
	i.Operator("mul", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))*p.Registers.Get((*m).Get(1)))
	})
	i.Operator("div", func(p *vm.ProcessorCore, m *vm.Memory) {
		defer func() {
			recover()
		}()
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))/p.Registers.Get((*m).Get(1)))
	})
	i.Operator("mod", func(p *vm.ProcessorCore, m *vm.Memory) {
		defer func() {
			recover()
		}()
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))%p.Registers.Get((*m).Get(1)))
	})

	i.Operator("and", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))&p.Registers.Get((*m).Get(1)))
	})
	i.Operator("or", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))|p.Registers.Get((*m).Get(1)))
	})
	i.Operator("xor", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))^p.Registers.Get((*m).Get(1)))
	})
	i.Operator("xnot", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))&^p.Registers.Get((*m).Get(1)))
	})
	i.Operator("lsft", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))<<uint(p.Registers.Get((*m).Get(1))))
	})
	i.Operator("rsft", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))>>uint(p.Registers.Get((*m).Get(1))))
	})

	return
}

func evaluate(p *vm.ProcessorCore) int {
	cost := 255 - p.Heap.Get(0)
	if cost < 0 {
		cost *= -1
	}
	return p.Cost() + cost
}

type Result struct {
	Cost int
	Core *vm.ProcessorCore
}

type ResultList []Result

func (s ResultList) Len() int           { return len(s) }
func (s ResultList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ResultList) Less(i, j int) bool { return s[i].Cost < s[j].Cost }

func bestManager(finished chan *vm.ProcessorCore, instructionSet *vm.InstructionSet, results chan int) {
	best := make(ResultList, 0)
	for {
		e := <-finished
		cost := e.Heap.Get(0)
		if(cost == 255){
			results <- 255
			//fmt.Println(255)
		}
		c := evaluate(e)
		best = append(best, Result{c, e})
		if len(best) >= POPULATION_SIZE {
			sort.Sort(best)
			best = best[0:BEST_OF_BREED]
			for _,b := range best {
				for i := 0; i < (POPULATION_SIZE / BEST_OF_BREED) - 1; i++ {
					prog := mutateProgram(&b.Core.Program, instructionSet, 0.1)
					p := NewSolver(instructionSet, prog, finished)
					go p.Run()
				}
				go NewSolver(instructionSet, &b.Core.Program, finished).Run()
				
			}
			best = make(ResultList, 0)
		}
	}
}

func mutateProgram(prog *vm.Program, instructions *vm.InstructionSet, chance float64) *vm.Program {
	outProg := make(vm.Program, len(*prog))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, x := range *prog {
		if rand.Float64() <= chance {
			decode := instructions.Decode(&x)
			if rand.Float64() <= 0.5 {
				decode.Set(0, r.Int())
			}
			if rand.Float64() <= 0.5 {
				decode.Set(1, r.Int())
			}
			if rand.Float64() <= 0.5 {
				decode.Set(2, r.Int())
			}
			outProg[i] = *instructions.Encode(decode)
		} else {
			outProg[i] = x
		}
	}
	return &outProg
}

func NewSolver(instructionSet *vm.InstructionSet, prog *vm.Program, c chan *vm.ProcessorCore) *vm.ProcessorCore {
	p := &vm.ProcessorCore{}
	p.Init(4, instructionSet, c)
	p.LoadProgram(prog)
	return p
}

func main() {
	finish := make(chan *vm.ProcessorCore)
	results := make(chan int)
	instructionSet := DefineInstructions()
	go bestManager(finish, instructionSet, results)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for y := 0; y < POPULATION_SIZE; y++ {

		pro := &vm.Program{
			instructionSet.Assemble("set", &vm.Memory{0, 1}),
		}
		for x := 0; x < 100; x++ {
			*pro = append(*pro, *instructionSet.Encode(&vm.Memory{r.Int(), r.Int(), r.Int()}))
		}
		p := NewSolver(instructionSet, pro, finish)
		go p.Run()
	}
	fmt.Println(<- results)
}
