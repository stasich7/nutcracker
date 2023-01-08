package common

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
)

type Character struct {
	Name     string
	Alias    string
	Type     string
	Kinds    []string
	Likes    []*Character
	Dislikes []*Character
	Relative []*Character
}

type Part struct {
	mx        sync.Mutex
	Character map[string]*Character
}

type Story struct {
	Part []*Part
}

func (p *Part) FanIn(ctx context.Context, wg *sync.WaitGroup, ch chan<- *Character, things ...*Character) {
	go func(things []*Character, wg *sync.WaitGroup) {
		fmt.Printf("--- Мысленно запущено ---\n")
		defer wg.Done()
		var iThings uint
		wg.Add(1)
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				iThings++
				thing := things[rand.Intn(len(things))]
				item := p.Instance(Character{Name: fmt.Sprintf("%s_%d", thing.Name, iThings), Type: thing.Type})
				ch <- item
			}
		}
	}(things, wg)
}

func (p *Part) FanOut(ctx context.Context, wg *sync.WaitGroup, ch <-chan *Character, consumer *Character) {
	go func(wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
			fmt.Printf("--- Мысленно остановлено ---\n")
		}()
		wg.Add(1)
		for item := range ch {
			p.Destroy(item, consumer)
		}
	}(wg)
}

func (c *Character) Say(ctx context.Context, what func()) {
	fmt.Printf("\n---------\n➜ %s говорит:\n", c.Name)
	what()
}

func (c *Character) Like(cWhom *Character) {
	isExist := false
	for _, v := range c.Likes {
		if v.Name == cWhom.Name {
			isExist = true
		}
	}
	if !isExist {
		c.Likes = append(c.Likes, cWhom)
	}
	fmt.Printf("%s нравится %s, [%s]\n", c.Name, cWhom.Name, printFields(*c))
}

func (c *Character) Dislike(cWhom *Character) {
	c.Dislikes = append(c.Dislikes, cWhom)
	fmt.Printf("%s не нравится %s, [%s]\n", c.Name, cWhom.Name, printFields(*c))
}

func (c *Character) AddKind(k string, originator ...*Character) {
	var originatorName string
	if len(originator) > 0 && originator[0] != nil {
		originatorName = fmt.Sprintf(" by %s", originator[0].Name)
	}
	c.Kinds = append(c.Kinds, k)
	fmt.Printf("%s получает новое свойство \"%s\"%s, [%s]\n", c.Name, k, originatorName, printFields(*c))
}

func (p *Part) Destroy(c *Character, originator ...*Character) {
	var originatorName string
	if len(originator) > 0 && originator[0] != nil {
		originatorName = fmt.Sprintf(" кем %s", originator[0].Name)
	}
	defer p.mx.Unlock()
	p.mx.Lock()
	fmt.Printf("%s уничтожен%s\n", c.Name, originatorName)
	delete(p.Character, c.Name)
}

func (p *Part) Instance(newC Character) *Character {
	if p == nil || newC.Name == "" {
		panic(fmt.Errorf("Cannot create instance of %#v in part %#v", newC, p))
	}

	defer p.mx.Unlock()
	p.mx.Lock()
	if _, ok := p.Character[newC.Name]; ok {
		return p.Character[newC.Name]
	}
	fmt.Printf("новый персонаж %s [%s]\n", newC.Name, printFields(newC))
	p.Character[newC.Name] = &newC
	return &newC
}

func (p *Part) IsExist(c string) bool {
	if _, ok := p.Character[c]; ok {
		return true
	}
	return false
}

func printFields(v interface{}) string {
	var result string
	s := reflect.ValueOf(v)
	for i := 0; i < s.NumField(); i++ {
		switch s.Field(i).Kind() {
		case reflect.String:
			val := s.Field(i).String()
			if val != "" {
				result += fmt.Sprintf("%s=%v ", s.Type().Field(i).Name, s.Field(i))
			}
		case reflect.Slice:
			if s.Field(i).Len() > 0 {
				var arrVal []string
				f := s.Field(i)
				for j := 0; j < f.Len(); j++ {
					val, ok := f.Index(j).Interface().(*Character)
					if !ok {
						valS := f.Index(j).String()
						arrVal = append(arrVal, valS)
						continue
					}
					arrVal = append(arrVal, val.Name)
				}
				result += fmt.Sprintf("%s=%s ",
					s.Type().Field(i).Name, arrVal)
			}
		}
	}
	return result
}
