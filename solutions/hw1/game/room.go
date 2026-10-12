package main

import "strings"

type Room struct {
	ID               string
	EntryDescription string
	EmptyDescription string
	Furniture        []*Furniture
	Targets          map[string]*Interaction
	Describe         func(*Game) string
}

type Furniture struct {
	Notation string
	Items    []*Item
}

func (r *Room) FindItem(name string) (*Item, *Furniture) {
	for _, furniture := range r.Furniture {
		for _, item := range furniture.Items {
			if item.Name == name {
				return item, furniture
			}
		}
	}
	return nil, nil
}

func (f *Furniture) RemoveItem(item *Item) {
	for i, candidate := range f.Items {
		if candidate == item {
			f.Items = append(f.Items[:i], f.Items[i+1:]...)
			return
		}
	}
}

func (r *Room) Contents() string {
	var groups []string
	for _, furniture := range r.Furniture {
		if len(furniture.Items) == 0 {
			continue
		}
		names := make([]string, 0, len(furniture.Items))
		for _, item := range furniture.Items {
			names = append(names, item.Name)
		}
		groups = append(groups, furniture.Notation+": "+strings.Join(names, ", "))
	}
	return strings.Join(groups, ", ")
}

func (r *Room) LookDescription(g *Game) string {
	if r.Describe != nil {
		return r.Describe(g)
	}
	if contents := r.Contents(); contents != "" {
		return contents
	}
	if r.EmptyDescription != "" {
		return r.EmptyDescription
	}
	return r.EntryDescription
}

func (r *Room) WithExits(transitions []*Transition, description string) string {
	if len(transitions) == 0 {
		return description
	}
	names := make([]string, 0, len(transitions))
	for _, exit := range transitions {
		names = append(names, exit.Name)
	}
	return description + ". можно пройти - " + strings.Join(names, ", ")
}
