package stringset

type void struct{}

type mapSet struct {
	set map[string]void
}

func NewMapSet() StringSet {
	return &mapSet{set: make(map[string]void)}
}

func (m *mapSet) Add(str string) bool {
	_, ok := m.set[str]
	if !ok {
		m.set[str] = void{}
	}
	return !ok
}

func (m *mapSet) In(str string) bool {
	_, ok := m.set[str]
	return ok
}

func (m *mapSet) Delete(str string) bool {
	_, ok := m.set[str]
	if ok {
		delete(m.set, str)
	}
	return ok
}
