//	TODO:	Rewrite Transform to run in parallel

package vm

type Program []*OpCode

func (p Program) Equals(other *Program) bool {
	if(len(p) != len(*other)){
		return false
	}
	for i, x := range p {
		if *x != *(*other)[i] {
			return false
		}
	}
	return true
}
