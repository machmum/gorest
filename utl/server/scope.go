package server

import (
	"errors"
)

var (
	ErrScopeNotFound = errors.New("Failed to get scope")

	ScopeLogin         = map[string]interface{}{"profile": false}
	ScopeProfile       = map[string]interface{}{"profile": false}
	ScopeProfileDetail = map[string]interface{}{"profile": false, "product": false}
)

func NewScope(req []string, def map[string]interface{}) *Scope {
	return &Scope{
		ScopeReq: req,
		ScopeDef: def,
	}
}

type Scope struct {
	// ScopeSet contains maps of scope-name value true
	ScopeSet map[string]interface{}
	// ScopeReq contains slice of scope-request
	ScopeReq []string
	// ScopeDef / scope-default contains default scope
	ScopeDef map[string]interface{}
}

// Check: return value in ScopeSet to true
// if meet conditions /
// key in ScopeSet exist in ScopeReq
func (s *Scope) Check() error {
	n := 0
	s.ScopeSet = make(map[string]interface{}, len(s.ScopeDef))
	for i := range s.ScopeReq {
		for idx := range s.ScopeDef {
			if idx == s.ScopeReq[i] {
				s.ScopeSet[idx] = true
			}
			n++
		}
	}

	if len(s.ScopeSet) < 1 {
		return ErrScopeNotFound
	}

	return nil
}
