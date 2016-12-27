package govirtual

import (
	"fmt"
)

type InstructionPipeline chan Operation

type Processor struct {
	*Memory
	*InstructionPipeline
}

func (p *Processor) String() string {
	return fmt.Sprintf("Processor [Memory: %v]", p.Memory)
}

//Start running the program. This won't return until the InstructionPipeline is closed.
func (p *Processor) Run() {
	for x := range *p.InstructionPipeline {
		x.Instruction.Implementation(p.Memory, x.Data...)
	}
}
