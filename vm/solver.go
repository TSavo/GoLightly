package vm

import (
	"github.com/seehuhn/mt19937"
	"github.com/tsavo/golightly/intutil"
	"math/rand"
	"sort"
	"time"
)

type Solver struct {
	PopulationSize, BestOfBreed, ProgramLength, HeapLength, RegisterLength int
	MutationChance                                                         float64
	Evaluator                                                              func(*Memory, int) int64
	instructionSet                                                         *InstructionSet
	rng                                                                    *rand.Rand
	Running                                                                bool
}

func (s *Solver) RandomProgram() *Program {
	pro := Program{}
	for x := 0; x < s.ProgramLength; x++ {
		pro = append(pro, s.RandomOpCode())
	}
	return &pro
}

func (s *Solver) RandomOpCode() *OpCode {
	return s.instructionSet.Encode(&Memory{s.rng.Int() % 2000, s.rng.Int() % 2000, s.rng.Int() % 2000})
}

func (s *Solver) NewProcessorCore(prog *Program, c chan *ProcessorCore) *ProcessorCore {
	p := &ProcessorCore{}
	p.Init(s.RegisterLength, s.instructionSet, c)
	p.LoadProgram(prog)
	return p
}

type Result struct {
	Cost int64
	Core *ProcessorCore
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

type Solution struct {
	Fitness int64
	Heaps []*Memory
	Champions []*Program
}

func (s *Solver) Solve(solutionChan chan *Solution, keepGoing chan bool, initialPop []*Program) {
	coreChan := make(chan *ProcessorCore)
	for _, x := range initialPop {
		go s.NewProcessorCore(x, coreChan).Run()
	}
	for y := len(initialPop); y < s.PopulationSize; y++ {
		go s.NewProcessorCore(s.RandomProgram(), coreChan).Run()
	}
	best := make(ResultList, 0)
	s.Running = true
	for s.Running {
		e := <-coreChan
		c := s.Evaluator(e.Heap, e.Cost())
		best = append(best, Result{c, e})
		if len(best) >= s.PopulationSize {
			sort.Sort(best)
			ok := s.BestOfBreed
			if len(best) < s.BestOfBreed {
				ok = len(best)
			}
			best = best[0:ok]
			best.Dedup()

			count := 0
			solution := &Solution{}
			solution.Fitness = int64(1000000)
			for _, b := range best {
				solution.Champions = append(solution.Champions, &b.Core.Program)
				solution.Heaps = append(solution.Heaps, b.Core.Heap)
				solution.Fitness = intutil.Min64(solution.Fitness, b.Cost)
			}
			solutionChan <- solution
			if !<-keepGoing {
				return
			}
			for _, b := range best {
				for i := 0; i < (s.PopulationSize/s.BestOfBreed)/3; i++ {
					co := s.NewProcessorCore(s.MutateProgram(&b.Core.Program), coreChan)
					go co.Run()
					count++
					co = s.NewProcessorCore(s.CombinePrograms(&b.Core.Program, &best[s.rng.Int()%len(best)].Core.Program), coreChan)
					go co.Run()
					count++
				}
				co := s.NewProcessorCore(&b.Core.Program, coreChan)
				go co.Run()
				count++
			}
			for count < s.PopulationSize {
				co := s.NewProcessorCore(s.RandomProgram(), coreChan)
				go co.Run()
				count++
			}
			best = make(ResultList, 0)
		}
	}
}

var seed = rand.New(mt19937.New())

func NewSolver(popSize int, bob int, pl int, rl int, hl int, chance float64, eval func(*Memory, int) int64, is *InstructionSet) *Solver {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano() + int64(seed.Int()))
	return &Solver{popSize, bob, pl, hl, rl, chance, eval, is, rng, false}
}

func (s *Solver) CombinePrograms(prog1 *Program, prog2 *Program) *Program {
	l1 := len(*prog1)
	l2 := len(*prog2)
	prog := make(Program, intutil.Max(l1, l2))
	split := s.rng.Int() % intutil.Min(l1, l2)
	endSplit := (s.rng.Int()%intutil.Min(l1, l2) - split) + split
	for x := 0; x < len(prog); x++ {
		if x > len(*prog1)-1 || (x < endSplit && x >= split && x < len(*prog2)) {
			prog[x] = (*prog2)[x]
		} else {
			prog[x] = (*prog1)[x]
		}
	}
	return &prog
}

func (s *Solver) MutateProgram(prog *Program) *Program {
	outProg := make(Program, 0)
	for _, x := range *prog {
		if s.rng.Float64() < s.MutationChance {
			if s.rng.Float64() < 0.1 {
				for r := s.rng.Int() % 10; r < 10; r++ {
					outProg = append(outProg, s.RandomOpCode())
				}
			}
			if s.rng.Float64() < 0.1 && len(*prog) > 50 {
				continue
			}
			decode := s.instructionSet.Decode(x)
			if s.rng.Float64() < 0.5 {
				decode.Set(0, s.rng.Int()%2000)
			}
			if s.rng.Float64() < 0.5 {
				decode.Set(1, s.rng.Int()%2000)
			}
			if s.rng.Float64() < 0.5 {
				decode.Set(2, s.rng.Int()%2000)
			}
			outProg = append(outProg, s.instructionSet.Encode(decode))

		} else {
			outProg = append(outProg, x)
		}
	}
	return &outProg
}
