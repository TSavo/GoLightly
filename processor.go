package govirtual

import (
	"fmt"
)

type InstructionPipeline chan Operation

//A Processor contains an instruction set, some memory, a Program, an Instruction Pointer
//(an index for where in the program it is), and a Termination Condition which
//let's it know when to stop.
type Processor struct {
	*InstructionSet
	*Memory
	*InstructionPipeline
}

func (p *Processor) String() string {
	return fmt.Sprintf("Processor [Memory: %v]", p.Memory)
}

//Create a new Processor with a memory, an instruction set, and a termination condition
func NewProcessor(instructions *InstructionSet, instructionPipeline *InstructionPipeline, memory *Memory, stop *TerminationCondition) *Processor {
	p := new(Processor)
	if instructions == nil {
		p.InstructionSet = new(InstructionSet)
	} else {
		p.InstructionSet = instructions
	}
	p.InstructionPipeline = instructionPipeline
	p.Memory = memory
	return p
}

//Start running the program. This won't return until the TerminationCondition.ShouldTerminate() returns true.
func (p *Processor) Run() {
	for x := range *p.InstructionPipeline {
		x.Instruction.Closure(p, x.Data...)
	}
}
