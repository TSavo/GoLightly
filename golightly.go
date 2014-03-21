package main

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"github.com/tsavo/golightly/vm"
	"math/rand"
	"sort"
	"time"
)

const (
	POPULATION_SIZE = 10000
	BEST_OF_BREED   = 100
)

func DefineInstructions() (i *vm.InstructionSet) {
	i = vm.NewInstructionSet()
	i.Operator("noop", func(p *vm.ProcessorCore, m *vm.Memory) {
	})
	i.Movement("halt", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Running = false
	})
	i.Movement("jump", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Jump(p.Registers.Get(1))
	})
	i.Movement("jumpIfZero", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) == 0 {
			p.Jump(p.Registers.Get(1))
		}
	})
	i.Movement("jumpIfNotZero", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) != 0 {
			p.Jump(p.Registers[1])
		}
	})
	i.Movement("call", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Call(p.Registers.Get((*m).Get(0)))
	})
	i.Movement("return", func(p *vm.ProcessorCore, m *vm.Memory) {
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
	i.Operator("increment", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Increment((*m).Get(0))
	})
	i.Operator("decrement", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Decrement((*m).Get(0))
	})
	i.Operator("add", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))+p.Registers.Get((*m).Get(1)))
	})
	i.Operator("subtract", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))-p.Registers.Get((*m).Get(1)))
	})
	i.Operator("muliply", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))*p.Registers.Get((*m).Get(1)))
	})
	i.Operator("divide", func(p *vm.ProcessorCore, m *vm.Memory) {
		defer func() {
			recover()
		}()
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))/p.Registers.Get((*m).Get(1)))
	})
	i.Operator("modulo", func(p *vm.ProcessorCore, m *vm.Memory) {
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
	i.Operator("leftShift", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))<<uint(p.Registers.Get((*m).Get(1))))
	})
	i.Operator("rightShift", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))>>uint(p.Registers.Get((*m).Get(1))))
	})

	return
}

func Abs(abs int64) int64 {
	if abs < 0 {
		abs = abs * -1
	}
	return abs
}

func evaluate(p *vm.ProcessorCore) int64 {
	if p.Heap.Get(0) > 10000 || p.Heap.Get(0) < 0 ||
		p.Heap.Get(1) > 10000 || p.Heap.Get(1) < 0 ||
		p.Heap.Get(2) > 10000 || p.Heap.Get(2) < 0 ||
		p.Heap.Get(3) > 10000 || p.Heap.Get(3) < 0 {
		return 10000
	}

	var cost int64 = Abs(255 - int64(p.Heap.Get(0)))
	cost += Abs(500 - int64(p.Heap.Get(1)))
	cost += Abs(255 - int64(p.Heap.Get(2)))
	cost += Abs(255 - int64(p.Heap.Get(3)))
	return Abs(int64(p.Cost()) + cost)
}

type Result struct {
	Cost int64
	Core *vm.ProcessorCore
}

type ResultList []Result

func (s ResultList) Len() int           { return len(s) }
func (s ResultList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ResultList) Less(i, j int) bool { return s[i].Cost < s[j].Cost }
func (s *ResultList) Dedup() {
	res := make(ResultList, 0)
outer:
	for x := 0; x < len(*s); x++ {
		for y := x + 1; y < len(*s); y++ {
			if (*s)[x].Core.Program.Equals(&(*s)[y].Core.Program) {
				continue outer
			}
		}
		res = append(res, (*s)[x])
	}
	(*s) = res
}

var pops = 0
var bestSoFar int64 = 1000

func bestManager(finished chan *vm.ProcessorCore, instructionSet *vm.InstructionSet, results chan int) {
	best := make(ResultList, 0)
	for {
		e := <-finished
		c := evaluate(e)
		if e.Heap.Get(0) == 255 && c < int64(bestSoFar) {
			results <- int(c)
			bestSoFar = c
			fmt.Println(e)
		}
		best = append(best, Result{c, e})
		if len(best) >= POPULATION_SIZE {
			best.Dedup()
			sort.Sort(best)
			ok := BEST_OF_BREED
			if len(best) < BEST_OF_BREED {
				ok = len(best)
			}
			best = best[0:ok]
			fmt.Println(best)
			count := 0
			for _, b := range best {
				for i := 0; i < (POPULATION_SIZE/BEST_OF_BREED)/2; i++ {
					go NewSolver(instructionSet, mutateProgram(&b.Core.Program, instructionSet, 0.1), finished).Run()
					count++
				}
				go NewSolver(instructionSet, &b.Core.Program, finished).Run()
				count++
			}
			for count < POPULATION_SIZE {
				count++
				pro := &vm.Program{}
				for x := 0; x < 100; x++ {
					*pro = append(*pro, *instructionSet.Encode(&vm.Memory{rng.Int() % 10000, rng.Int() % 10000, rng.Int() % 10000}))
				}
				r := NewSolver(instructionSet, pro, finished)
				go r.Run()
			}
			pops++
			best = make(ResultList, 0)

		}
	}
}

func mutateProgram(prog *vm.Program, instructions *vm.InstructionSet, chance float64) *vm.Program {
	outProg := make(vm.Program, 0)

	for _, x := range *prog {
		if rng.Float64() < chance {

			//one in 100 times, add
			if rng.Float64() < 0.1 {

				for r := rng.Int() % 10; r < 10; r++ {
					outProg = append(outProg, *instructions.Encode(&vm.Memory{rng.Int() % 10000, rng.Int() % 10000, rng.Int() % 10000}))
				}
			}

			//one in 100 times, delete
			if rng.Float64() < 0.1 && len(*prog) > 20 {
				continue
			}
			decode := instructions.Decode(&x)
			if rng.Float64() < 0.5 {
				decode.Set(0, rng.Int()%10000)
			}
			if rng.Float64() < 0.5 {
				decode.Set(1, rng.Int()%10000)
			}
			if rng.Float64() < 0.5 {
				decode.Set(2, rng.Int()%10000)
			}
			outProg = append(outProg, *instructions.Encode(decode))

		} else {
			outProg = append(outProg, x)
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

var rng = rand.New(mt19937.New())

func main() {
	rng.Seed(time.Now().UnixNano())
	finish := make(chan *vm.ProcessorCore)
	results := make(chan int)
	instructionSet := DefineInstructions()
	go bestManager(finish, instructionSet, results)

	for y := 0; y < POPULATION_SIZE; y++ {

		pro := &vm.Program{}
		for x := 0; x < 100; x++ {
			*pro = append(*pro, *instructionSet.Encode(&vm.Memory{rng.Int() % 10000, rng.Int() % 10000, rng.Int() % 10000}))
		}
		p := NewSolver(instructionSet, pro, finish)
		go p.Run()
	}
	for {
		fmt.Println(<-results)
	}
}
