package interpreter

type Symbol struct {
	Name, Scope string
	Index       int
}

type SymbalTable struct {
	outer *SymbalTable
	store map[string]*Symbol
	size  int
}

func NewEnclosedSymbalTable(outer *SymbalTable) *SymbalTable {
	st := NewSymbolTable()
	st.outer = outer
	return st
}

func NewSymbolTable() *SymbalTable {
	return &SymbalTable{
		store: make(map[string]*Symbol),
		size:  0,
	}
}

func (s *SymbalTable) Define(name string) *Symbol {
	symbol := Symbol{Name: name, Scope: GlobalScope, Index: s.size}
	if s.outer != nil {
		symbol.Scope = LocalScope
	}
	s.store[name] = &symbol
	s.size++
	return &symbol
}

func (s *SymbalTable) Resolve(name string) (*Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.outer != nil {
		symbol, ok = s.outer.Resolve(name)
	}
	return symbol, ok
}
