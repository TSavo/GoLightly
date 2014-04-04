package vm

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"github.com/tsavo/golightly/intutil"
	"math/rand"
	"sort"
	"time"
)

type ProgramGenerator interface {
	GenerateProgram() *Program
}

type Evaluator interface {
	Evaluate(*ProcessorCore) int64
}

type Solver struct {
	Id, PopulationSize, BestOfBreed, RegisterLength int
	ChanceOfMutation                                float64
	InstructionSet                                  *InstructionSet
	rng                                             *rand.Rand
	Running                                         bool
	ProgramGenerator
	Evaluator
	ControlChan  chan bool
	SolutionChan chan *Solution
}

type Processor struct {
	Reward int64
	Core   *ProcessorCore
}

func (p *Processor) String() string {
	return fmt.Sprintf("%d: %v", p.Reward, p.Core)
}

type ProcessorList []*Processor

func (s ProcessorList) Len() int           { return len(s) }
func (s ProcessorList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ProcessorList) Less(i, j int) bool { return s[i].Reward > s[j].Reward }
func (s *ProcessorList) Dedup() {
	res := make(ProcessorList, 0)
outer:
	for x := 0; x < len(*s); x++ {
		for y := x + 1; y < len(*s); y++ {
			if ((*s)[x].Core.Program).Equals((*s)[y].Core.Program) {
				continue outer
			}
		}
		res = append(res, (*s)[x])
	}
	(*s) = res
}

type SolutionProgram struct {
	Reward  int64
	Program string
}

type Solution struct {
	Id       int
	Programs []SolutionProgram
}

func (s *Solver) NewProcessorCore(prog *Program, heap *Memory) *ProcessorCore {
	p := NewProcessorCore(s.RegisterLength, s.InstructionSet, heap)
	p.LoadProgram(prog)
	return p
}

func (s *Solver) SolveOneAtATime(sharedMemory *Memory, coreChan chan *ProcessorCore, solutionChan chan *Solution, controlChan chan bool, stopChan chan bool, populationInfluxChan chan []string, initialPop []*Program) {
	processors := make(ProcessorList, s.PopulationSize)
	for x := 0; x < len(initialPop); x++ {
		c := s.NewProcessorCore(initialPop[x], sharedMemory)
		c.Heap = sharedMemory
		processors[x] = &Processor{int64(0), c}
	}
	for x := 0; x < s.PopulationSize-len(initialPop); x++ {
		c := s.NewProcessorCore(s.GenerateProgram(), sharedMemory)
		c.Heap = sharedMemory
		processors[x] = &Processor{int64(0), c}
	}
	count := 0
outer:
	for {
		select {
		case <-stopChan:
			return
		default:
		}
		if count < s.PopulationSize {
			processors[count].Core.ControlChan = controlChan
			(*processors[count]).Core.Reset()
			(*processors[count]).Reward = int64(0)
			select {
			case coreChan <- (*processors[count]).Core:
			default:
			}
			fmt.Printf("#%d: %d\n", s.Id, count)
			(*processors[count]).Core.Run()
			(*processors[count]).Reward = s.Evaluate((*processors[count]).Core)
			count++
			select {
			case <-stopChan:
				return
			default:
				continue outer
			}
		}

		sort.Sort(processors)
		best := processors[:s.BestOfBreed]
		fmt.Println(best)
		rest := processors[s.BestOfBreed:]
		//
		bestPrograms := make([]SolutionProgram, len(best))
		for i, x := range best {
			bestPrograms[i] = SolutionProgram{x.Reward, x.Core.GetProgramString()}
		}
		select {
		case solutionChan <- &Solution{s.Id, bestPrograms}:
		default:
		}
		select {
		case <-stopChan:
			return
		default:
		}
		count = 0
		select {
		case inFlux := <-populationInfluxChan:
			fmt.Println("INFLUX!")
			for x := 0; x < len(inFlux) && count < len(rest); x++ {
				rest[count].Core.CompileAndLoad(inFlux[x])
				count++
			}
		default:
		}
		for count < len(rest)/2 {
			for _, b := range best {
				rest[count].Core.LoadProgram(s.MutateProgram(b.Core.Program))
				count++
				rest[count].Core.LoadProgram(s.CombinePrograms(b.Core.Program, best[s.rng.Int()%len(best)].Core.Program))
				count++
			}
		}
		for count < len(rest) {
			rest[count].Core.LoadProgram(s.GenerateProgram())
			count++
		}
		count = 0
	}
}

var seed = rand.New(mt19937.New())

func NewSolver(id int, popSize int, bob int, rl int, chance float64, is *InstructionSet, gen ProgramGenerator, eval Evaluator) *Solver {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano() + seed.Int63())
	return &Solver{id, popSize, bob, rl, chance, is, rng, false, gen, eval, make(chan bool, 1), make(chan *Solution, 100)}
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

func (s *Solver) RandomOperation() *Operation {
	return ((*s.InstructionSet).Encode(&Memory{s.rng.Int(), s.rng.Int(), s.rng.Int()}))
}

func (s *Solver) MutateProgram(prog *Program) *Program {
	outProg := make(Program, 0)
	for _, x := range *prog {
		if s.rng.Float64() < s.ChanceOfMutation {
			if s.rng.Float64() < 0.1 {
				for r := s.rng.Int() % 10; r < 10; r++ {
					outProg = append(outProg, s.RandomOperation())
				}
			}
			if s.rng.Float64() < 0.1 && len(outProg) > 0 {
				continue
			}
			decode := s.InstructionSet.Decode(x)
			if s.rng.Float64() < 0.5 {
				decode.Set(0, s.rng.Int())
			}
			if s.rng.Float64() < 0.5 {
				decode.Set(1, s.rng.Int()%2000)
			}
			if s.rng.Float64() < 0.5 {
				decode.Set(2, s.rng.Int()%2000)
			}
			outProg = append(outProg, s.InstructionSet.Encode(decode))
		} else {
			outProg = append(outProg, x)
		}
	}
	return &outProg
}
