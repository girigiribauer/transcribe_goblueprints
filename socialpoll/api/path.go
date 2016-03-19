package main

import (
	"strings"
)

// PathSeparator ...
const PathSeparator = "/"

// Path 型
type Path struct {
	Path string
	ID   string
}

// NewPath パスの文字列を受け取り、Path 型のインスタンスを返す
func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{
		Path: p,
		ID:   id,
	}
}

// HasID は Path 型が ID を持っているかどうかを返す
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
