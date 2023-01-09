package ch11

import (
	"context"
	e "nutcracker/pkg/common"
	"sync"
	"time"
)

const (
	Title = "Глава 11. Фрагмент про кота."
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

	mother := p.Instance(e.Character{
		Name:  "Мама",
		Type:  "человек",
		Alias: "Советница",
	})

	fritz := p.Instance(e.Character{
		Name:     "Фриц",
		Type:     "человек",
		Relative: []*e.Character{mother},
	})

	// -- Фриц --
	// Ну как не быть? — весело воскликнул Фриц. — Внизу у булочника есть отличный
	// серый кот; надо его взять к нам наверх, а он уж отлично обделает дело, будь
	// даже эта скверная мышь сама крыса Мышильда или сын ее, мышиный король!
	fritz.Say(ctx, func() {
		greyCat := p.Instance(e.Character{
			Name:     "Кот",
			Type:     "животное",
			Kinds:    []string{"Серый"},
			Relative: []*e.Character{{Name: "Булочник", Type: "человек"}},
		})
		fritz.Like(greyCat)
		rodent := make(chan *e.Character)
		wg := sync.WaitGroup{}
		ctxM, cancel := context.WithTimeout(ctx, time.Microsecond*500)
		defer cancel()
		p.FanIn(ctxM, &wg, rodent,
			&e.Character{Name: "Мышильда", Type: "животное"},
			&e.Character{Name: "Мышиный король", Type: "животное"})
		p.FanOut(ctxM, &wg, rodent, greyCat)
		wg.Wait()

	})

	// -- мама --
	//— Да, — прибавила мама, смеясь, — а заодно начнет прыгать по стульям, столам и перебьет всю посуду и чашки.
	ctxM, cancelM := context.WithCancel(ctx)
	defer cancelM()
	wgM := sync.WaitGroup{}
	mother.Say(ctx, func() {
		mother.Dislike(p.Instance(e.Character{Name: "Кот"}))
		cups := make(chan *e.Character)
		p.FanIn(ctxM, &wgM, cups, &e.Character{Name: "Чашка"}, &e.Character{Name: "Тарелка"})
		p.FanOut(ctxM, &wgM, cups, p.Instance(e.Character{Name: "Кот"}))
	})

	// -- Фриц --
	// — О нет! — возразил Фриц. — Это очень ловкий кот; я бы хотел сам уметь так
	// лазать по крышам, как он.
	fritz.Say(ctx, func() {
		greycat := p.Instance(e.Character{Name: "Кот"})
		fritz.Like(greycat)
		greycat.AddKind("Прыгать по крышам", fritz)
	})

	// -- Луиза --
	// — Нет уж, пожалуйста, нельзя ли обойтись ночью без кошек, — сказала Луиза, которая их очень не любила.
	louisa := p.Instance(e.Character{
		Name:     "Луиза",
		Type:     "человек",
		Relative: []*e.Character{mother},
	})

	louisa.Say(ctx, func() {
		<-time.After(time.Microsecond * 500)
		cancelM()
		wgM.Wait()
		p.Destroy(p.Instance(e.Character{Name: "Кот"}))
	})

	// -- советник --
	// — Я думаю, — сказал советник, — что Фриц прав, а пока можно будет поставить и мышеловку; ведь у нас она есть?
	father := p.Instance(e.Character{
		Name:  "Отец",
		Type:  "человек",
		Alias: "Советник",
	})
	father.Say(ctx, func() {
		father.Like(fritz)
		if p.IsExist("Мышеловка") {
			mousetrap := p.Instance(e.Character{Name: "Мышеловка"})
			mice := make(chan *e.Character)
			wg := sync.WaitGroup{}
			ctxM, cancel := context.WithTimeout(ctx, time.Microsecond*500)
			defer cancel()
			p.FanIn(ctxM, &wg, mice, &e.Character{Name: "Мышь", Type: "животное"})
			p.FanOut(ctxM, &wg, mice, mousetrap)
			wg.Wait()
		}
	})

	// -- Фриц --
	// — Что за беда, если и нет, — закричал Фриц, — крестный тотчас сделает новую! Ведь он же их выдумал!
	fritz.Say(ctx, func() {
		if !p.IsExist("MouseTrap") {
			drosselmeier := p.Instance(e.Character{
				Name:  "Дроссельмейер",
				Type:  "человек",
				Alias: "Крестный",
			})
			mousetrap := p.Instance(e.Character{Name: "Мышеловка"})
			drosselmeier.Like(mousetrap)
			mice := make(chan *e.Character)
			wg := sync.WaitGroup{}
			ctxM, cancel := context.WithTimeout(ctx, time.Microsecond*500)
			defer cancel()
			p.FanIn(ctxM, &wg, mice, &e.Character{Name: "Мышь", Type: "животное"})
			p.FanOut(ctxM, &wg, mice, mousetrap)
			wg.Wait()
		}
	})

	// -- советница --
	// Все засмеялись, когда же советница сказала, что у них в самом деле нет мышеловок
	mother.Say(ctx, func() {
		mother.Like(fritz)
		father.Like(fritz)
		p.Instance(e.Character{Name: "Мышеловка"}).Like(fritz)
		louisa.Like(fritz)
		p.Destroy(p.Instance(e.Character{Name: "Мышеловка"}))

	})

	// -- крестный --
	// крестный объявил, что у него в доме их много, и тотчас велел принести одну, отлично сделанную.
	drosselmeier := p.Instance(e.Character{Name: "Дроссельмейер"})
	drosselmeier.Say(ctx, func() {
		mousetrap := p.Instance(e.Character{
			Name:     "Мышеловка",
			Relative: []*e.Character{{Name: "Дом", Type: "строение"}},
		})
		mousetrapC := make(chan *e.Character)
		wg := sync.WaitGroup{}
		ctxM, cancel := context.WithTimeout(ctx, time.Microsecond*500)
		defer cancel()
		p.FanIn(ctxM, &wg, mousetrapC, mousetrap)
		p.FanOut(ctxM, &wg, mousetrapC, mousetrap)
		wg.Wait()
		drosselmeier.Like(mousetrap)
	})
}
