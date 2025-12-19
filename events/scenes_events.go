package events

import (
	"github.com/jo-jordan/go-holdem/cmd"
)

type SceneEvent interface{ isSceneEvent() }

type SceneReady struct{}
type SceneChatMessage struct {
	ChatCmd cmd.ChatCmd
}

type GameStarted struct{}

type PlayerReady struct{}

type PlayerCall struct{}

type PlayerRaise struct{}

type PlayerCheck struct{}

type PlayerFold struct{}

type GameEnded struct{}

// markers
func (GameStarted) isSceneEvent()      {}
func (SceneChatMessage) isSceneEvent() {}
func (SceneReady) isSceneEvent()       {}
