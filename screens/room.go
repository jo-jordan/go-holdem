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

type focusType uint8

const (
	TABLE focusType = iota
	MESSAGE
)

type Room struct {
	screen

	name string
	game entities.Game

	player  *entities.Player
	seat    Seat
	actions Actions
	msgBox  *ui.Message
	status  Status
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
		name:   opt.Name,
		screen: opt.Screen,
		status: initStatus(),
	}
	return room.
		initGame(&opt).
		initUI().
		initActs().
		initSeats()
}

func (room *Room) initActs() *Room {
	acts := []actConf{
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

	room.actions = initActions(&style, acts)
	return room
}

func (room *Room) initSeats() *Room {
	room.seat = initSeat(room.style.GetWidth())
	for i, player := range room.game.Players {
		s := seatStyle
		if player != nil && *room.player == *player {
			s = focusSeatStyle
		}
		room.seat.setSeat(i, ui.NewSeat(ui.SeatOpt{
			Num:    i + 1,
			Player: player,
			Style:  s,
		}))
	}

	room.focusOnGame()
	return room
}

func (room *Room) initUI() *Room {
	room.msgBox = ui.NewMessage(room.player, room.screen.style)

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
			room.seat.style.Render(
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
	seats := room.seat.seatView()
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
	acts := room.actions.actionView()
	actTable := table.
		New().
		Border(lipgloss.HiddenBorder())

	actTable.Row(room.actions.keyView(room.status.isActActive)...)
	actTable.Row(acts...)
	return *actTable
}

func (room *Room) selectFocusAct(k string) tea.Cmd {
	for index, key := range room.actions.actKeys {
		if key != k {
			continue
		}
		if room.actions.current == index {
			return room.actions.blur()
		}
		return tea.Batch(room.actions.blur(), room.actions.focusAt(index))
	}
	return nil
}

func (room *Room) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	if room.status.isFocusOnTable() {
		switch key := msg.String(); key {
		case "space":
			room.status.toggle()
			if !room.status.isActActive {
				room.actions.blur()
			}
		case "c", "a", "r", "f", "s", "i", "q":
			cmd = room.selectFocusAct(key)
		case "enter":
			_, cmd = room.actions.update(msg)
		}
	} else {
		_, cmd = room.msgBox.Update(msg)
	}
	return cmd
}

func (room *Room) focusOnGame() {
	room.status.focusOn(TABLE)
	room.seat.focus()
}

func (room *Room) focusOnMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)
	room.status.focusOn(MESSAGE)
	room.status.isActActive = false
	room.seat.blur()
	cmd = room.actions.blur()
	cmds = append(cmds, cmd)

	_, cmd = room.msgBox.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

type Seat struct {
	seats [entities.MAX_SEAT]*ui.Seat
	style lipgloss.Style
}

func initSeat(width int) Seat {
	seat := Seat{
		style: lipgloss.NewStyle().
			Width(width).
			Border(lipgloss.NormalBorder()),
	}

	return seat
}

func (seat *Seat) setSeat(index int, s *ui.Seat) {
	seat.seats[index] = s
}

func (seat *Seat) seatView() []string {
	seats := make([]string, len(seat.seats))
	for i := range seats {
		seats[i] = seat.seats[i].View()
	}
	return seats
}

func (seat *Seat) focus() {
	seat.style = seat.style.BorderForeground(lipgloss.Color(ui.FOCUS_COLOR))
}

func (seat *Seat) blur() {
	seat.style = seat.style.BorderForeground(lipgloss.Color(ui.NORMAL_COLOR))
}

type actConf struct {
	text    string
	quick   string
	actions []*ui.ActionMap
}

type Actions struct {
	actKeys []string
	acts    []*ui.Button
	current int
}

func initActions(style *lipgloss.Style, acts []actConf) Actions {
	actions := Actions{
		current: -1,
	}
	actions.actKeys = make([]string, len(acts))
	actions.acts = make([]*ui.Button, len(acts))
	for i, v := range acts {
		actions.actKeys[i] = v.quick
		actions.acts[i] = ui.NewButton(ui.ButtonOption{
			Value:   v.text,
			Style:   style,
			Actions: v.actions,
		})
	}

	return actions
}

func (actions Actions) actionView() []string {
	acts := make([]string, len(actions.acts))
	for i := range acts {
		acts[i] = actions.acts[i].View()
	}
	return acts
}

func (actions Actions) keyView(active bool) []string {
	if active {
		return actions.actKeys
	} else {
		return make([]string, len(actions.actKeys))
	}
}

func (actions *Actions) focusAt(index int) tea.Cmd {
	actions.current = index
	_, cmd := actions.acts[actions.current].Update(Cmd.FocusMsg{})
	return cmd
}

func (actions *Actions) blur() tea.Cmd {
	if actions.current == -1 {
		return nil
	}
	_, cmd := actions.acts[actions.current].Update(Cmd.BlurMsg{})
	actions.current = -1
	return cmd
}

func (actions *Actions) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if actions.current != -1 {
		return actions.acts[actions.current].Update(msg)
	}
	return nil, nil
}

type Status struct {
	focusOnType focusType
	isActActive bool
}

func initStatus() Status {
	return Status{
		focusOnType: TABLE,
		isActActive: false,
	}
}

func (s Status) isFocusOnTable() bool {
	return s.focusOnType == TABLE
}

func (s *Status) toggle() {
	s.isActActive = !s.isActActive
}

func (s *Status) focusOn(t focusType) {
	s.focusOnType = t
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
