// Package main provides a set server
package main

import (
	"github.com/kitschysynq/biolviewity"
)

type Set struct {
	items map[string]struct{}
}

func (s *Set) Add(i string) (bool, error) {
	if _, ok := s.items[i]; ok {
		return true, nil
	}
	s.items[i] = struct{}{}
	return false, nil
}

func (s *Set) Remove(i string) (bool, error) {
	if _, ok := s.items[i]; ok {
		delete(s.items, i)
		return true, nil
	}
	return false, nil
}

func (s *Set) Contains(i string) (bool, error) {
	_, ok := s.items[i]
	return ok, nil
}

func main() {
	biolviewity.AddHandler(&Set{})
	biolviewity.ListenAndServe()
}
