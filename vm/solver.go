package vm

import (
	"fmt"
)

type Solver struct {
	Id, RegisterLength int
	InstructionSet                                  *InstructionSet
	Breeder
	Evaluator
	Selector
	ControlChan      chan bool
	SolverReportChan chan *SolverReport
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
	Id        int
	SolutionList
}

func (s SolutionList) Len() int           { return len(s) }
func (s SolutionList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SolutionList) Less(i, j int) bool { return s[i].Reward > s[j].Reward }

func NewSolver(id int, rl int, is *InstructionSet, gen Breeder, eval Evaluator, selector Selector) *Solver {
	return &Solver{id, rl, is, gen, eval, selector, make(chan bool, 1), make(chan *SolverReport, 1)}
}

func (s *Solver) SolveOneAtATime(sharedMemory *Memory, solutionChan chan *SolverReport, stopCondition *TerminationCondition, stopChan chan bool, populationInfluxChan chan []*Solution, initialPop []*Program) {
	processors := make([]*ProcessorCore, 0)
	programs := s.Breed(nil)
	for {
		solutions := make(SolutionList, len(programs))
		for len(processors) < len(solutions) {
			c := NewProcessorCore(s.RegisterLength, s.InstructionSet, sharedMemory, stopCondition)
			processors = append(processors, c)
		}
		if len(processors) > len(solutions) {
			processors = processors[:len(solutions)]
		}
		for x, pro := range processors {
			select {
			case <-stopChan:
				return
			default:
			}
			fmt.Printf("#%d: %d\n", s.Id, x)
			pro.CompileAndLoad(programs[x])
			pro.Run()
			solutions[x] = &Solution{s.Evaluate(pro), pro.Program.Decompile()}
		}
		select {
		case solutionChan <- &SolverReport{s.Id, solutions}:
		default:
		}
		programs = s.Breed(s.Select(&solutions).GetPrograms())
	}
}

