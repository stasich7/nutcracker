package ch11

import (
	"context"
	e "nutcracker/pkg/common"
	"sync"
	"time"
)

func New() *e.Part {
	part := &e.Part{}
	part.Character = make(map[string]*e.Character)
	return part
}

func Tell(p *e.Part) {
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
		ctxM, cancel := context.WithCancel(ctx)
		p.FanIn(ctxM, &wg, rodent,
			&e.Character{Name: "Мышильда", Type: "животное"},
			&e.Character{Name: "Мышиный король", Type: "животное"})
		p.FanOut(ctxM, &wg, rodent, greyCat)
		<-time.After(time.Microsecond * 500)
		cancel()
		wg.Wait()

	})

	// -- мама --
	//— Да, — прибавила мама, смеясь, — а заодно начнет прыгать по стульям, столам и перебьет всю посуду и чашки.
	ctxM, cancelM := context.WithCancel(ctx)
	wgM := sync.WaitGroup{}
	mother.Say(ctx, func() {
		mother.Dislike(p.Instance(e.Character{Name: "Кот"}))
		cups := make(chan *e.Character)
		p.FanIn(ctxM, &wgM, cups, &e.Character{Name: "Чашка"}, &e.Character{Name: "Тарелка"})
		p.FanOut(ctxM, &wgM, cups, p.Instance(e.Character{Name: "Кот"}))
		<-time.After(time.Microsecond * 500)
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
		p.Destroy(p.Instance(e.Character{Name: "Мышь"}))
		cancelM()
		wgM.Wait()
	})

	// -- советник --
	// — Я думаю, — сказал советник, — что Фриц прав, а пока можно будет поставить и мышеловку; ведь у нас она есть?
	father := p.Instance(e.Character{
		Name:  "Отец",
		Type:  "человек",
		Alias: "Советник",
	})
	fatherAct4 := func() {
		father.Like(fritz)
		if p.IsExist("Мышеловка") {
			mousetrap := p.Instance(e.Character{Name: "Мышеловка"})
			mice := make(chan *e.Character)
			wg := sync.WaitGroup{}
			ctxM, cancel := context.WithCancel(ctx)
			p.FanIn(ctxM, &wg, mice, &e.Character{Name: "Мышь", Type: "животное"})
			p.FanOut(ctxM, &wg, mice, mousetrap)
			<-time.After(time.Microsecond * 500)
			cancel()
			wg.Wait()
		}
	}
	father.Say(ctx, fatherAct4)

	// -- Фриц --
	// — Что за беда, если и нет, — закричал Фриц, — крестный тотчас сделает новую! Ведь он же их выдумал!
	fritzAct5 := func() {
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
			ctxM, cancel := context.WithCancel(ctx)
			p.FanIn(ctxM, &wg, mice, &e.Character{Name: "Мышь", Type: "животное"})
			p.FanOut(ctxM, &wg, mice, mousetrap)
			<-time.After(time.Microsecond * 500)
			cancel()
			wg.Wait()
		}
	}
	fritz.Say(ctx, fritzAct5)

	// -- советница --
	// Все засмеялись, когда же советница сказала, что у них в самом деле нет мышеловок
	motherAct6 := func() {
		mother.Like(fritz)
		father.Like(fritz)
		p.Instance(e.Character{Name: "Мышеловка"}).Like(fritz)
		louisa.Like(fritz)
		p.Destroy(p.Instance(e.Character{Name: "Мышеловка"}))

	}
	mother.Say(ctx, motherAct6)

	// -- крестный --
	// крестный объявил, что у него в доме их много, и тотчас велел принести одну, отлично сделанную.
	drosselmeier := p.Instance(e.Character{Name: "Дроссельмейер"})
	drosselmeierAct7 := func() {
		mousetrap := p.Instance(e.Character{
			Name:     "Мышеловка",
			Relative: []*e.Character{{Name: "Дом", Type: "строение"}},
		})
		mousetrapC := make(chan *e.Character)
		wg := sync.WaitGroup{}
		ctxM, cancel := context.WithCancel(ctx)
		p.FanIn(ctxM, &wg, mousetrapC, mousetrap)
		p.FanOut(ctxM, &wg, mousetrapC, mousetrap)
		<-time.After(time.Microsecond * 500)
		cancel()
		wg.Wait()
		drosselmeier.Like(mousetrap)
	}
	drosselmeier.Say(ctx, drosselmeierAct7)
}
