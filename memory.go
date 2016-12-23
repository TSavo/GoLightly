package govirtual

import (
	"fmt"
	"hash/crc32"
	"strings"
)

func CompletelyDereference(value interface{}) interface{} {
	switch x := value.(type) {
	case Pointer:
		return CompletelyDereference(x.Get())
	default:
		return x
	}
}

func Dereference(value interface{}) interface{} {
	switch x := value.(type) {
	case Pointer:
		return x.Get()
	default:
		return x
	}
}

type Pointer interface {
	Get() interface{}
	Set(interface{})
}

type Variable struct {
	Pointer
	Name string
}

func (variable *Variable) Get() interface{} {
	return variable.Pointer.Get()
}

func (variable *Variable) Set(value interface{}) {
	variable.Pointer.Set(value)
}

func (variable *Variable) String() string {
	return variable.Name
}

type Literal struct {
	Value interface{}
}

func (literal *Literal) Get() interface{} {
	return literal.Value
}

func (literal *Literal) Set(value interface{}) {
	literal.Value = value
}

func (literal *Literal) String() string {
	switch x := literal.Value.(type) {
	case string:
		if strings.HasPrefix(x, ":") {
			return x
		} else {
			return fmt.Sprintf("\"%v\"", x)
		}
	default:
		return fmt.Sprintf("%v", x)
	}

}

type Reference struct {
	Pointer
}

func (v *Reference) Get() interface{} {
	return v.Pointer.Get()
}

func (v *Reference) Set(value interface{}) {
	v.Pointer.Set(value)
}

func (v *Reference) String() string {
	return fmt.Sprintf("&%v", v.Pointer)
}

type Memory []Pointer

type MemoryPointer struct {
	*Memory
	Index interface{}
	Name  string
}

func (memory *MemoryPointer) Get() interface{} {
	return memory.Memory.Get(memory.Index).Get()
}

func (memory *MemoryPointer) Set(value interface{}) {
	memory.Memory.Get(memory.Index).Set(value)
}

func (memory *MemoryPointer) String() string {
	return fmt.Sprintf("%s[%d]", memory.Name, memory.Index)
}

func Booleanize(in interface{}) bool {
	return Cardinalize(in)%2 == 0
}

func Cardinalize(in interface{}) int {
	switch x := in.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		return int(crc32.ChecksumIEEE([]byte(x)))
	case Pointer:
		return Cardinalize(x.Get())
	default:
		return 0
	}
}

func (m *Memory) Push(p Pointer) {
	*m = append(*m, p)
}

func (s *Memory) Pop() (r interface{}, ok bool) {
	if end := s.Len() - 1; end > -1 {
		r = (*s)[end]
		*s = (*s)[:end]
		ok = true
	}
	return
}

func (s *Memory) Delete(i int) {
	a := *s
	n := len(a)
	if i > -1 && i < n {
		copy(a[i:n-1], a[i+1:n])
		*s = a[:n-1]
	}
}

func (s *Memory) Resize(size int) {
	n := make(Memory, size, size)
	copy(n, (*s))
	*s = n
}

func (m *Memory) Len() (l int) {
	l = len(*m)
	return
}

func (m *Memory) Get(i interface{}) Pointer {
	x := Cardinalize(i)
	l := m.Len()
	if l < 1 {
		panic("Memory is of size < 1")
	}
	if x < 0 {
		x *= -1
	}
	x = x % l
	defer recover()
	return (*m)[x]
}

func (m *Memory) GetCardinal(i interface{}) int {
	return Cardinalize(m.Get(i))
}

func (m *Memory) Set(i interface{}, x Pointer) {
	xx := Cardinalize(i)
	l := m.Len()
	if l < 1 {
		panic("Memory is of size < 1")
	}
	if xx < 0 {
		xx *= -1
	}
	xx = xx % l
	defer recover()
	(*m)[xx] = x
}

func (m *Memory) Reallocate(size int) {
	(*m) = make(Memory, size)
}

func (s *Memory) Append(v Pointer) {
	*s = append(*s, v)
}

func (s *Memory) Prepend(v Pointer) {
	*s = append([]Pointer{v}, *s...)
}

func (s *Memory) Zero() {
	for x := 0; x < s.Len(); x++ {
		(*s)[x] = &Literal{0}
	}
}
