package govirtual

import . "math/big"

type Memory []int

type FloatMemory []*Rat

func NewFloatMemory(size int) *FloatMemory {
	f := make(FloatMemory, size)
	for x, _ := range f {
		rat := Rat{}
		rat.SetFloat64(0.0)
		f[x] = &rat
	}
	return &f
}

func (m *Memory) Push(p int) {
	*m = append(*m, p)
}

func (m *FloatMemory) Push(p *Rat) {
	r := &Rat{}
	r.Set(p)
	*m = append(*m, r)
}

func (s *Memory) Pop() (r int, ok bool) {
	if end := s.Len() - 1; end > -1 {
		r = (*s)[end]
		*s = (*s)[:end]
		ok = true
	}
	return
}

func (s *FloatMemory) Pop() (r *Rat, ok bool) {
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

func (s *FloatMemory) Delete(i int) {
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

func (s *FloatMemory) Resize(size int) {
	n := make(FloatMemory, size, size)
	copy(n, (*s))
	*s = n
}

func (m *Memory) Len() (l int) {
	l = len(*m)
	return
}

func (m *FloatMemory) Len() (l int) {
	l = len(*m)
	return
}

func (m *Memory) Get(i int) int {
	l := m.Len()
	if l < 1 {
		panic("Memory is of size < 1")
	}
	if i < 0 {
		i *= -1
	}
	i = i % l
	defer recover()
	return (*m)[i]
}

func (m *FloatMemory) Get(i int) *Rat {
	l := m.Len()
	if l < 1 {
		panic("Memory is of size < 1")
	}
	if i < 0 {
		i *= -1
	}
	i = i % l
	defer recover()
	return (*m)[i]
}

func (m *Memory) Set(i int, x int) {
	l := m.Len()
	if l < 1 {
		panic("Memory is of size < 1")
	}
	if i < 0 {
		i *= -1
	}
	i = i % l
	defer recover()
	(*m)[i] = x
}

func (m *FloatMemory) Set(i int, x *Rat) {
	l := m.Len()
	if l < 1 {
		panic("Memory is of size < 1")
	}
	if i < 0 {
		i *= -1
	}
	i = i % l
	defer recover()
	(*m)[i].Set(x)
}

func (m *FloatMemory) SetFloat(i int, x float64) {
	m.Get(i).SetFloat64(x)
}

func (m *Memory) Increment(i int) {
	m.Set(i, m.Get(i)+1)
}

func (m *Memory) Decrement(i int) {
	m.Set(i, m.Get(i)-1)
}

func (m *Memory) Zero() {
	for i := 0; i < m.Len(); i++ {
		m.Set(i, 0)
	}
}

func (m *FloatMemory) Zero() {
	for i := 0; i < m.Len(); i++ {
		m.Get(i).SetFloat64(0.0)
	}
}

func (m *Memory) Reallocate(size int) {
	(*m) = make(Memory, size)
}

func (m *FloatMemory) Reallocate(size int) {
	*m = *NewFloatMemory(size)
}

func (s *Memory) Append(v interface{}) {
	switch v := v.(type) {
	case int:
		*s = append(*s, v)
	case Memory:
		*s = append(*s, v...)
	case []int:
		s.Append(Memory(v))
	default:
		panic(v)
	}
}

func (m *FloatMemory) Append(v interface{}) {
	switch v := v.(type) {
	case *Rat:
		*m = append(*m, v)
	case FloatMemory:
		*m = append(*m, v...)
	case []*Rat:
		m.Append(FloatMemory(v))
	default:
		panic(v)
	}
}

func (s *Memory) Prepend(v interface{}) {
	switch v := v.(type) {
	case int:
		*s = append([]int{v}, *s...)
	case Memory:
		*s = append(v, *s...)
	case []int:
		s.Prepend(Memory(v))
	default:
		panic(v)
	}
}
