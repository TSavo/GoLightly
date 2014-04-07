//	TODO:	Rewrite Transform to run in parallel

package vm

type Program []*Operation

func (p *Program) Equals(other *Program) bool {
	if(len(*p) != len(*other)){
		return false
	}
	for i, x := range *p {
		if *x != *(*other)[i] {
			return false
		}
	}
	return true
}

func (p *Program) Hashcode() int {
	h := 0
	for i, x := range *p {
		h += x.Instruction.id * (i+1)
	}
	return h
}

func (p *Program) Decompile() string {
	pro := ""
	for _, x := range *p {
		pro += x.String() + "\n"
	}
	return pro
}
