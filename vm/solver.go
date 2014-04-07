package vm

import (
	"fmt"
)

type Solver struct {
	Id, RegisterLength   int
	InstructionSet       *InstructionSet
	Breeder              *Breeder
	Evaluator            *Evaluator
	Selector             *Selector
	TerminationCondition *TerminationCondition
	ControlChan          chan bool
	SolverReportChan     chan *SolverReport
	Heap                 *Memory
}

type Solution struct {
	Reward  int64
	Program string
}

type SolutionList []*Solution

func (sol *SolutionList) GetPrograms() []string {
	x := make([]string, len(*sol))
	for i, solution := range *sol {
		x[i] = solution.Program
	}
	return x
}

type SolverReport struct {
	Id int
	SolutionList
}

func (s SolutionList) Len() int           { return len(s) }
func (s SolutionList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SolutionList) Less(i, j int) bool { return s[i].Reward > s[j].Reward }

func NewSolver(id int, sharedMemory *Memory, rl int, is *InstructionSet, term TerminationCondition, gen Breeder, eval Evaluator, selector Selector) *Solver {
	return &Solver{id, rl, is, &gen, &eval, &selector, &term, make(chan bool, 1), make(chan *SolverReport, 1), sharedMemory}
}

func (s *Solver) SolveOneAtATime() {
	programs := (*s.Breeder).Breed(nil)
	processors := make([]*ProcessorCore, 0)
	for {
		solutions := make(SolutionList, len(programs))
		for len(processors) < len(solutions) {
			c := NewProcessorCore(s.RegisterLength, s.InstructionSet, s.Heap, s.TerminationCondition)
			processors = append(processors, c)
		}
		if len(processors) > len(solutions) {
			processors = processors[:len(solutions)]
		}
		for x, pro := range processors {
			select {
			case <-s.ControlChan:
				return
			default:
			}
			fmt.Printf("#%d: %d\n", s.Id, x)
			pro.Reset()
			pro.CompileAndLoad(programs[x])
			pro.Run()
			solutions[x] = &Solution{(*s.Evaluator).Evaluate(pro), pro.Program.Decompile()}
		}
		select {
		case s.SolverReportChan <- &SolverReport{s.Id, solutions}:
		default:
		}
		programs = (*s.Breeder).Breed((*s.Selector).Select(&solutions).GetPrograms())
	}
}
