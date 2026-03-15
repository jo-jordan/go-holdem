package screens

import (
	"fmt"
	"strconv"

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
}

type RoomSetupOps struct {
	Screen screen
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
	roomSetup.screen = ops.Screen
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
			Title:     "Initial Account: ",
			IsNum:     true,
			TextWidth: ACCOUNT_WIDTH,
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
			Title:     "Small Blind: ",
			IsNum:     true,
			TextWidth: ACCOUNT_WIDTH,
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
						initAccount, _ := strconv.Atoi(roomSetup.initAccountInput.Value())
						smallBlind, _ := strconv.Atoi(roomSetup.smallBlindInput.Value())
						if initAccount < smallBlind {
							return nil, nil
						}
						return NewRoom(RoomOps{
							Name:        roomSetup.nameInput.Value(),
							SmallBlind:  uint(smallBlind),
							InitAccount: uint(initAccount),
							Player:      roomSetup.player,
							Screen:      roomSetup.screen,
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
	var cmd tea.Cmd
	var model tea.Model
	switch msg := msg.(type) {
	case tea.KeyPressMsg, tea.WindowSizeMsg:
		cmd := roomSetup.screen.Update(msg)
		if cmd != nil {
			return roomSetup, cmd
		}
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
	game := &JoinGame{
		name:   opt.Name,
		player: opt.Player,
	}
	game.style = opt.Style
	return game.initTarget().
		initJoinButton().
		initCancelButton().
		initUI()
}

func (g *JoinGame) initTarget() *JoinGame {
	g.target = ui.NewInputText(ui.InputTextOption{
		Title:     "Target: ",
		Focus:     true,
		TextWidth: ACCOUNT_WIDTH,
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.EnterToNext,
			ui.ShiftTabToPrev,
		},
	},
	)
	return g
}

func (g *JoinGame) initJoinButton() *JoinGame {
	g.joinButton = ui.NewButton(ui.ButtonOption{
		Value: "Join",
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.ShiftTabToPrev,
			{
				Msg: "enter",
				Act: func() (tea.Model, tea.Cmd) {
					target := g.target.Value()
					if target == "" {
						return nil, nil
					}
					return NewRoom(RoomOps{Name: target}), nil
				},
			},
		},
	})
	return g
}

func (g *JoinGame) initCancelButton() *JoinGame {
	g.cancelButton = ui.NewButton(ui.ButtonOption{
		Value: "Cancel",
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.ShiftTabToPrev,
			{
				Msg: "enter",
				Act: func() (tea.Model, tea.Cmd) {
					return NewStartScreen(StartScreenOpt{
						PlayerName: g.player.Name,
					}).WithStyle(&g.style), nil
				},
			},
		},
	})
	return g
}

func (g *JoinGame) initUI() *JoinGame {
	g.CursorMove = ui.NewCursorMove(ui.CursorMoveOption{
		Models: []ui.Elementer{
			g.target,
			g.joinButton,
			g.cancelButton,
		},
	})

	return g
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

const ACCOUNT_WIDTH = 32
