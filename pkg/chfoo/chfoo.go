package chfoo

import (
	"context"
	e "nutcracker/pkg/common"
)

const (
	Title = "Глава Foo"
)

type LocPart struct {
	e.Part
}

func New() *LocPart {
	part := &LocPart{}
	part.Name = Title
	part.Character = make(map[string]*e.Character)
	return part
}

func (p *LocPart) Tell() {
	p.Title()

	ctx := context.Background()

	foo := p.Instance(e.Character{
		Name: "Foo",
		Type: "fake",
	})

	foo.Say(ctx, func() {
		bar := p.Instance(e.Character{
			Name: "Bar",
			Type: "fake",
		})
		foo.Like(bar)
	})

}
