package main

import (
	"github.com/tsavo/golightly/intutil"
	"github.com/tsavo/golightly/vm"
	"runtime"
	"fmt"
)

const (
	POPULATION_SIZE = 20000
	BEST_OF_BREED   = 200
	PROGRAM_LENGTH  = 50
	UNIVERSE_SIZE = 10
	ROUND_LENGTH = 4
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
	i.Operator("multiply", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))*p.Registers.Get((*m).Get(1)))
	})
	i.Operator("divide", func(p *vm.ProcessorCore, m *vm.Memory) {
		defer func() {
			recover()
		}()
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))/p.Registers.Get((*m).Get(1)))
	})
	i.Operator("modulos", func(p *vm.ProcessorCore, m *vm.Memory) {
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
	i.Operator("xand", func(p *vm.ProcessorCore, m *vm.Memory) {
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

func evaluate(p *vm.Memory, cost int) int64 {
	if p.Get(0) > 10000 || p.Get(0) < 0 ||
		p.Get(1) > 10000 || p.Get(1) < 0 ||
		p.Get(2) > 10000 || p.Get(2) < 0 ||
		p.Get(3) > 10000 || p.Get(3) < 0 {
		return 10000
	}

	var fit int64 = intutil.Abs64(255 - int64(p.Get(0)))
	fit += intutil.Abs64(500 - int64(p.Get(1)))
	fit += intutil.Abs64(255 - int64(p.Get(2)))
	fit += intutil.Abs64(255 - int64(p.Get(3)))
	return intutil.Abs64(int64(cost) + fit)
}

func CombinePrograms(s1 []*vm.Program, s2 []*vm.Program) []*vm.Program {
	prog := []*vm.Program{}
	prog = append(prog, s1[0:len(s1)/2]...)
	prog = append(prog, s2[0:len(s2)/2]...)
	return prog
}

type Member struct {
	Solver       *vm.Solver
	Fitness      int64
	SolutionChan chan *vm.Solution
	ControlChan  chan bool
}

type Population []*Member

func main() {
	runtime.LockOSThread()
	instructionSet := DefineInstructions()
	//finish := make(chan *vm.ProcessorCore)
	//results := make(chan int)
	population := make(Population, 0)
	for x := 0; x < UNIVERSE_SIZE; x++ {
		solver := vm.NewSolver(POPULATION_SIZE, BEST_OF_BREED, PROGRAM_LENGTH, 4, 4, 0.1, evaluate, instructionSet)
		solutionChan := make(chan *vm.Solution)
		control := make(chan bool)
		population = append(population, &Member{solver, 10000, solutionChan, control})
		go solver.Solve(solutionChan, control, nil)
	}
	go func() {
		count := 0
		bestFit := int64(100000)
		for {
			count++
			champs := make([]*vm.Program, 0)
			for _, member := range population {
				solution := <-member.SolutionChan
				if(bestFit > solution.Fitness){
					fmt.Println(solution.Fitness)
					bestFit = solution.Fitness
					fmt.Println(solution.Heaps[0])
				}
				if count < ROUND_LENGTH {
					member.ControlChan <- true
				} else {
					member.ControlChan <- false
					champs = append(champs, solution.Champions[0:int(float32(len(solution.Champions))*0.1)]...)
				}
			}
			if count < ROUND_LENGTH {
				continue
			}
			fmt.Printf("recombo: %d\n", len(champs))
			count = 0
			population = make(Population, 0)
			for x := 0; x < UNIVERSE_SIZE-1; x++ {
				solver := vm.NewSolver(POPULATION_SIZE, BEST_OF_BREED, PROGRAM_LENGTH, 4, 4, 0.1, evaluate, instructionSet)
				solutionChan := make(chan *vm.Solution)
				control := make(chan bool)
				population = append(population, &Member{solver, 10000, solutionChan, control})
				go solver.Solve(solutionChan, control, nil)
			}
			solver := vm.NewSolver(POPULATION_SIZE, BEST_OF_BREED, PROGRAM_LENGTH, 4, 4, 0.1, evaluate, instructionSet)
			solutionChan := make(chan *vm.Solution)
			control := make(chan bool)
			population = append(population, &Member{solver, 10000, solutionChan, control})
			go solver.Solve(solutionChan, control, champs)
		}
	}()
	<-make(chan int)
}
