package events

import "github.com/jo-jordan/go-holdem/cmd"

type SceneEvent interface{ isSceneEvent() }

type SceneReady struct{}
type SceneChatMessage struct {
	ChatCmd cmd.ChatCmd
}

// markers
func (SceneChatMessage) isSceneEvent() {}
func (SceneReady) isSceneEvent()       {}
