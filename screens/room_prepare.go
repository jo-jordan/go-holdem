package screens

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/ui"
)

type RoomSetup struct {
	screen
	ui.CursorMove

	nameInput        *ui.InputText
	initAccountInput *ui.InputText
	smallBlindInput  *ui.InputText
	startButton      *ui.Button
	backButton       *ui.Button

	player entities.Player
	style  lipgloss.Style
}

type RoomSetupOps struct {
	Style  lipgloss.Style
	Player entities.Player
}

func NewRootSetup(ops RoomSetupOps) *RoomSetup {
	roomSetup := &RoomSetup{}
	roomSetup.initName().
		initAccount().
		initSmallBlind().
		initStartButton().
		initCancelButton()

	roomSetup.player = ops.Player
	roomSetup.style = ops.Style
	roomSetup.CursorMove = ui.NewCursorMove(ui.CursorMoveOption{
		Models: []ui.Elementer{
			roomSetup.nameInput,
			roomSetup.initAccountInput,
			roomSetup.smallBlindInput,
			roomSetup.startButton,
			roomSetup.backButton,
		},
	})
	return roomSetup
}

func (roomSetup *RoomSetup) initName() *RoomSetup {
	roomSetup.nameInput = ui.NewInputText(
		ui.InputTextOption{
			Title: "Root Name: ",
			Focus: true,
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.EnterToNext,
				ui.ShiftTabToPrev,
			},
		},
	)
	return roomSetup
}

func (roomSetup *RoomSetup) initAccount() *RoomSetup {
	roomSetup.initAccountInput = ui.NewInputText(
		ui.InputTextOption{
			Title: "Initial Account: ",
			IsNum: true,
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.EnterToNext,
				ui.ShiftTabToPrev,
			},
		},
	)
	return roomSetup
}

func (roomSetup *RoomSetup) initSmallBlind() *RoomSetup {
	roomSetup.smallBlindInput = ui.NewInputText(
		ui.InputTextOption{
			Title: "Small Blind: ",
			IsNum: true,
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.EnterToNext,
				ui.ShiftTabToPrev,
			},
		},
	)
	return roomSetup
}

func (roomSetup *RoomSetup) initStartButton() *RoomSetup {
	roomSetup.startButton = ui.NewButton(
		ui.ButtonOption{
			Value: "Start",
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.ShiftTabToPrev,
				{
					Msg: "enter",
					Act: func() (tea.Model, tea.Cmd) {
						return NewRoom(RoomOps{
							Name:  roomSetup.nameInput.Value(),
							Style: roomSetup.style,
						}), nil
					},
				},
			},
		},
	)
	return roomSetup
}

func (roomSetup *RoomSetup) initCancelButton() *RoomSetup {
	roomSetup.backButton = ui.NewButton(
		ui.ButtonOption{
			Value: "Back",
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.ShiftTabToPrev,
				{
					Msg: "enter",
					Act: func() (tea.Model, tea.Cmd) {
						start := NewStartScreen(StartScreenOpt{
							PlayerName: roomSetup.player.Name,
						}).WithStyle(&roomSetup.style)
						return start, nil
					},
				},
			},
		},
	)
	return roomSetup
}

func (roomSetup *RoomSetup) Init() tea.Cmd {
	return textinput.Blink
}

func (roomSetup *RoomSetup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	cmd := roomSetup.screen.Update(msg)
	if cmd != nil {
		return roomSetup, cmd
	}

	model, cmd = roomSetup.CursorMove.Update(msg)
	if model == nil {
		model = roomSetup
	}
	return model, cmd
}

func (roomSetup *RoomSetup) View() tea.View {
	v := tea.NewView(roomSetup.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			roomSetup.nameInput.View(),
			roomSetup.initAccountInput.View(),
			roomSetup.smallBlindInput.View(),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				roomSetup.startButton.View(),
				roomSetup.backButton.View(),
			),
		),
	))
	v.AltScreen = true
	return v
}

type JoinGame struct {
	screen
	ui.CursorMove
	name         string
	target       *ui.InputText
	joinButton   *ui.Button
	cancelButton *ui.Button
	player       entities.Player
}

type JoinGameOption struct {
	Name   string
	Style  lipgloss.Style
	Player entities.Player
}

func NewJoinGame(opt JoinGameOption) *JoinGame {
	game := JoinGame{
		name: opt.Name,
		target: ui.NewInputText(ui.InputTextOption{
			Title: "Target: ",
			Focus: true,
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.EnterToNext,
				ui.ShiftTabToPrev,
			},
		}),
	}

	game.player = opt.Player
	game.joinButton = ui.NewButton(ui.ButtonOption{
		Value: "Join",
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.ShiftTabToPrev,
			{
				Msg: "enter",
				Act: func() (tea.Model, tea.Cmd) {
					target := game.target.Value()
					if target == "" {
						return nil, nil
					}
					return NewRoom(RoomOps{Name: target, Style: game.style}), nil
				},
			},
		},
	})
	game.cancelButton = ui.NewButton(ui.ButtonOption{
		Value: "Cancel",
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.ShiftTabToPrev,
			{
				Msg: "enter",
				Act: func() (tea.Model, tea.Cmd) {
					return NewStartScreen(StartScreenOpt{
						PlayerName: game.player.Name,
					}).WithStyle(&game.style), nil
				},
			},
		},
	})
	game.CursorMove = ui.NewCursorMove(ui.CursorMoveOption{
		Models: []ui.Elementer{
			game.target,
			game.joinButton,
			game.cancelButton,
		},
	})
	game.style = opt.Style
	return &game
}

func (g *JoinGame) Init() tea.Cmd {
	return textinput.Blink
}

func (g *JoinGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	cmd := g.screen.Update(msg)
	if cmd != nil {
		return g, cmd
	}

	model, cmd = g.CursorMove.Update(msg)
	if model == nil {
		model = g
	}
	return model, cmd
}

func (g *JoinGame) View() tea.View {
	v := tea.NewView(g.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf("User: %s", g.name),
			g.target.View(),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				g.joinButton.View(),
				g.cancelButton.View(),
			),
		),
	))
	v.AltScreen = true
	return v
}
