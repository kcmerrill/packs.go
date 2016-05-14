package packs

import (
	"sort"
)

type trigger struct {
	events map[int][]*event
}

func (t *trigger) add(e *event) {
	t.events[e.Priority] = append(t.events[e.Priority], e)
}

func (t *trigger) priorities() []int {
	p := []int{}
	for priority := range t.events {
		p = append(p, priority)
	}
	sort.Ints(p)
	return p
}

func NewTrigger() *trigger {
	t := &trigger{}
	t.events = make(map[int][]*event)
	return t
}
