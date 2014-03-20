package vm

type Memory []int

func (m *Memory) Push(p int) {
	l := len(*m) + 1
	n := make(Memory, l, l)
	copy(n, (*m))
	n[l-1] = p
	*m = n
}

func (s *Memory) Pop() (r int, ok bool) {
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
		copy(a[i:n - 1], a[i + 1:n])
		*s = a[:n - 1]
	}
}

func (s *Memory) Resize(size int){
	n := make(Memory, size, size)
	copy(n, (*s))
	*s = n
}

func (m Memory) Len() (l int) {
	l = len(m)
	return
}

func (m Memory) Get(i int) int {
	return m[i%m.Len()]
}

func (m Memory) Set(i int, x int) {
	m[i%m.Len()] = x
}

func (m Memory) Increment(i int){
	m.Set(i, m.Get(i) + 1)
}

func (m Memory) Decrement(i int){
	m.Set(i, m.Get(i) - 1)
}

func (m *Memory) Zero(){
	for i := 0 ; i < m.Len() ; i++ {
		m.Set(i, 0);
	} 
}

func (m *Memory) Reallocate(size int){
	(*m) = make(Memory, size)
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

func (s *Memory) Prepend(v interface{}) {
	switch v := v.(type) {
	case int:
		l := s.Len() + 1
		n := make(Memory, l, l)
		n[0] = v
		copy(n[1:], *s)
		*s = n

	case Memory:
		l := s.Len() + len(v)
		n := make(Memory, l, l)
		copy(n, v)
		copy(n[len(v):], *s)
		*s = n

	case []int:
		s.Prepend(Memory(v))
	default:
		panic(v)
	}
}
