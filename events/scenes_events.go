package events

type SceneEvent interface{ isSceneEvent() }

type SceneReadyToAnnounce struct{}
type SceneChatMessage struct{ Text string }
