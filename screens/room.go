package screens

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/ui"
)

type Room struct {
	screen
	ui.CursorMove

	name string
	game entities.Game

	msgBox *ui.Message
}

type RoomOps struct {
	Name   string
	Style  lipgloss.Style
	Player entities.Player
}

func NewRoom(opt RoomOps) *Room {
	room := &Room{
		name: opt.Name,
	}
	return room.initUI(&opt).
		initGame(&opt)
}

func (room *Room) initUI(opt *RoomOps) *Room {
	room.style = opt.Style
	room.msgBox = ui.NewMessage(opt.Style)

	room.CursorMove = ui.NewCursorMove(ui.CursorMoveOption{
		Models: []ui.Elementer{
			room.msgBox,
		},
	})
	return room
}

func (room *Room) initGame(opt *RoomOps) *Room {
	room.game = *entities.NewGame()
	room.game.AddPlayer(&opt.Player)
	return room
}

func (room *Room) Init() tea.Cmd {
	return textinput.Blink
}

func (room *Room) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmd = room.screen.Update(msg)

	if cmd != nil {
		return room, cmd
	}

	_, cmd = room.msgBox.Update(msg)
	return room, cmd
}

func (room *Room) View() tea.View {
	v := tea.NewView(room.msgBox.View())
	v.AltScreen = true
	return v
}
