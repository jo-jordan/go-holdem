package screens

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	Cmd "github.com/jo-jordan/go-holdem/cmd"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/ui"
)

type FocusType uint8

const (
	TABLE FocusType = iota
	MESSAGE
)

type Room struct {
	screen

	name string
	game entities.Game

	player        *entities.Player
	seats         [entities.MAX_SEAT]*ui.Seat
	seatsStyle    lipgloss.Style
	actKeys       []string
	acts          []*ui.Button
	msgBox        *ui.Message
	focusOn       FocusType
	fnActive      bool
	selectFnIndex int
}

type RoomOps struct {
	Screen      screen
	Name        string
	SmallBlind  uint
	InitAccount uint
	Player      entities.Player
}

func NewRoom(opt RoomOps) *Room {
	room := &Room{
		name:          opt.Name,
		screen:        opt.Screen,
		selectFnIndex: -1,
	}
	return room.
		initUI().
		initGame(&opt).
		initActs().
		initSeats()
}

func (room *Room) initActs() *Room {
	acts := []struct {
		text    string
		quick   string
		actions []*ui.ActionMap
	}{
		{
			text:  "Check",
			quick: "c",
		},
		{
			text:  "Call",
			quick: "a",
		},
		{
			text:  "Raise",
			quick: "r",
		},
		{
			text:  "Fold",
			quick: "f",
		},
		{
			text:  "Show Cards",
			quick: "s",
		},
		{
			text:  "Chat",
			quick: "i",
			actions: []*ui.ActionMap{
				{
					Msg: "enter",
					Act: func() (tea.Model, tea.Cmd) {
						return nil, func() tea.Msg {
							return Cmd.MsgFocusMsg{}
						}
					},
				},
			},
		},
		{
			text:  "Quit",
			quick: "q",
			actions: []*ui.ActionMap{
				{
					Msg: "enter",
					Act: func() (tea.Model, tea.Cmd) {
						return nil, tea.Quit
					},
				},
			},
		},
	}
	maxLength := 0
	for _, v := range acts {
		if len(v.text) > maxLength {
			maxLength = len(v.text)
		}
	}
	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(maxLength + 2)

	room.actKeys = make([]string, len(acts))
	room.acts = make([]*ui.Button, len(acts))
	for i, v := range acts {
		room.actKeys[i] = v.quick
		room.acts[i] = ui.NewButton(ui.ButtonOption{
			Value:   v.text,
			Style:   &style,
			Actions: v.actions,
		})
	}
	return room
}

func (room *Room) initSeats() *Room {
	for i, player := range room.game.Players {
		s := seatStyle
		if player != nil && *room.player == *player {
			s = focusSeatStyle
		}
		room.seats[i] = ui.NewSeat(ui.SeatOpt{
			Num:    i + 1,
			Player: player,
			Style:  s,
		})
	}

	room.seatsStyle = lipgloss.NewStyle().
		Width(room.screen.style.GetWidth()).
		Border(lipgloss.NormalBorder())
	room.focusOnGame()
	return room
}

func (room *Room) initUI() *Room {
	room.msgBox = ui.NewMessage(room.screen.style)

	width := room.screen.style.GetWidth()
	// calculate the width of each seat
	seatWidth := (width-3)/4 - 2
	seatStyle = seatStyle.Width(seatWidth)
	focusSeatStyle = focusSeatStyle.Width(seatWidth)
	return room
}

func (room *Room) initGame(opt *RoomOps) *Room {
	room.player = &opt.Player
	opt.Player.Account = int(opt.InitAccount)
	room.game = *entities.NewGame()
	room.game.AddPlayer(&opt.Player)
	room.game.SmallBlind = opt.SmallBlind
	room.game.InitAccount = opt.InitAccount
	return room
}

func (room *Room) Init() tea.Cmd {
	// init connection here
	return textinput.Blink
}

func (room *Room) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model := room

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd = room.handleKey(msg)
	case tea.WindowSizeMsg:
		cmd = room.screen.Update(msg)
	case Cmd.TableFocusMsg:
		room.focusOnGame()
	case Cmd.MsgFocusMsg:
		cmd = room.focusOnMsg(msg)
	}

	return model, cmd
}

func (room *Room) View() tea.View {
	gameTable := room.gameTable()
	actTable := room.actTable()
	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Top,
			room.seatsStyle.Render(
				lipgloss.JoinVertical(lipgloss.Top,
					gameTable.Render(),
					actTable.Render(),
				),
			),
			room.msgBox.View(),
		),
	)
	v.AltScreen = true
	return v
}

func (room *Room) gameTable() table.Table {
	seats := make([]string, len(room.seats))
	for i := range seats {
		seats[i] = room.seats[i].View()
	}

	t := table.
		New().
		Border(lipgloss.HiddenBorder()).
		Row(seats[0:4]...).
		Row(seats[4], "", "", seats[5]).
		Row(seats[6:10]...)
	return *t
}

// Maybe optimize this method for performance
func (room *Room) actTable() table.Table {
	acts := make([]string, len(room.acts))
	for i := range acts {
		acts[i] = room.acts[i].View()
	}
	actTable := table.
		New().
		Border(lipgloss.HiddenBorder())

	if room.fnActive {
		actTable.Row(room.actKeys...)
	} else {
		quickKeys := make([]string, len(room.actKeys))
		actTable.Row(quickKeys...)
	}

	actTable.Row(acts...)
	return *actTable
}

func (room *Room) selectFocusAct(k string) tea.Cmd {
	var cmd tea.Cmd
	for index, key := range room.actKeys {
		if key == k {
			// cancel current focus
			if room.selectFnIndex != -1 {
				_, cmd = room.acts[room.selectFnIndex].Update(Cmd.BlurMsg{})
			}
			room.selectFnIndex = index
			_, cmd = room.acts[index].Update(Cmd.FocusMsg{})
			return cmd
		}
	}
	room.selectFnIndex = -1
	return cmd
}

func (room *Room) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	if room.focusOn == TABLE {
		switch key := msg.String(); key {
		case "space":
			room.fnActive = !room.fnActive
		case "c", "a", "r", "f", "s", "i", "q":
			cmd = room.selectFocusAct(key)
		case "enter":
			if 0 <= room.selectFnIndex && room.selectFnIndex < len(room.acts) {
				_, cmd = room.acts[room.selectFnIndex].Update(msg)
			}
		}
	} else {
		_, cmd = room.msgBox.Update(msg)
	}
	return cmd
}

func (room *Room) focusOnGame() {
	room.focusOn = TABLE
	room.seatsStyle = room.seatsStyle.BorderForeground(lipgloss.Color(ui.FOCUS_COLOR))
}

func (room *Room) focusOnMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)
	room.focusOn = MESSAGE
	room.fnActive = false
	room.seatsStyle = room.seatsStyle.BorderForeground(lipgloss.Color(ui.NORMAL_COLOR))
	_, cmd = room.acts[room.selectFnIndex].Update(Cmd.BlurMsg{})
	cmds = append(cmds, cmd)
	room.selectFnIndex = -1

	_, cmd = room.msgBox.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

var (
	seatStyle = lipgloss.
			NewStyle().
			Border(
			lipgloss.NormalBorder(),
		).
		Height(3).
		BorderForeground(lipgloss.Color(ui.NORMAL_COLOR))
	focusSeatStyle = seatStyle.BorderForeground(lipgloss.Color(ui.FOCUS_COLOR))
)
