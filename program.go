//	TODO:	Rewrite Transform to run in parallel

package govirtual

type Program struct {
	Operations []*Operation
}

func NewProgram(size int) *Program {
	return &Program{make([]*Operation, size)}
}

func (p *Program) Len() int {
	return len(p.Operations)
}

func (p *Program) Equals(other *Program) bool {
	if len(p.Operations) != len(other.Operations) {
		return false
	}
	for i, x := range p.Operations {
		if *x != *other.Operations[i] {
			return false
		}
	}
	return true
}

func (p *Program) Append(op *Operation) *Program {
	p.Operations = append(p.Operations, op)
	return p
}

func (p *Program) Get(index int) *Operation {
	if index < 0 {
		index *= -1
	}
	return p.Operations[index%p.Len()]
}

func (p *Program) Clone() *Program {
	pr := NewProgram(0)
	for x := 0; x < p.Len(); x++ {
		pr.Append(p.Get(x))
	}
	return pr
}

func (p *Program) Decompile() string {
	pro := ""
	for _, x := range p.Operations {
		pro += x.String() + "\n"
	}
	return pro
}

func (p *Program) Labels() (labels map[string][]int) {
	labels = make(map[string][]int)
	for x, op := range p.Operations {
		if len(op.Label) > 0 {
			labels[op.Label] = append(labels[op.Label], x)
		}
	}
	return
}

func (p *Program) LabelNames() []string {
	labels := p.Labels()
	names := make([]string, len(labels))
	x := 0
	for k, _ := range labels {
		names[x] = k
		x++
	}
	return names
}
