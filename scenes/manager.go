package scenes

import "github.com/jo-jordan/go-holdem/cmd"

type SceneManager struct {
	out chan<- cmd.GameCmd
	in  <-chan cmd.GameCmd
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		out: make(chan cmd.GameCmd),
		in:  make(chan cmd.GameCmd),
	}
}
