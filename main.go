package main

import (
	"nutcracker/pkg/ch11"
	"nutcracker/pkg/chfoo"
	"nutcracker/pkg/common"
)

func main() {

	story := common.Story{Name: "Щелкунчик и мышиный король"}
	story.Tell(ch11.New())
	story.Tell(chfoo.New())
}
