package govirtual

import (
	"hash/crc32"
)

func Dereference(value interface{}) interface{} {
	x, ok := value.(Pointer)
	if ok {
		return Dereference(x.Get())
	}
	return value
}

type Pointer interface {
	Get() interface{}
	Set(interface{})
}

type VariablePointer struct {
	Pointer
	Name string
}

func (v *VariablePointer) Get() interface{} {
	return v.Pointer.Get()
}

func (v *VariablePointer) Set(value interface{}) {
	v.Pointer.Set(value)
}

func (v *VariablePointer) String() string {
	return v.Name
}

type Memory []interface{}

type MemoryPointer struct {
	Memory
	Index interface{}
	Name  string
}

func (memory *MemoryPointer) Get() interface{} {
	return memory.Memory.Get(memory.Index)
}

func (memory *MemoryPointer) Set(value interface{}) {
	memory.Memory.Set(memory.Index, value)
}

func (memory *MemoryPointer) String() string {
	return memory.Name
}

func Cardinalize(in interface{}) int {
	in = Dereference(in)
	switch x := in.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		return int(crc32.ChecksumIEEE([]byte(x)))
	default:
		return 0
	}
}

func (m *Memory) Push(p int) {
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

func (m *Memory) Get(i interface{}) interface{} {
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

func (m *Memory) Set(i interface{}, x interface{}) {
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

func (s *Memory) Append(v interface{}) {
	*s = append(*s, v)
}

func (s *Memory) Prepend(v interface{}) {
	*s = append([]interface{}{v}, *s...)
}

func (s *Memory) Zero() {
	for x := 0; x < s.Len(); x++ {
		(*s)[x] = 0
	}
}
