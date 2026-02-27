package screens

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/ui"
)

type Room struct {
	screen
	ui.CursorMove

	name  string
	game  entities.Game
	seats [entities.MAX_SEAT]*ui.Seat

	msgBox *ui.Message
}

type RoomOps struct {
	Screen screen
	Name   string
	Player entities.Player
}

func NewRoom(opt RoomOps) *Room {
	room := &Room{
		name:   opt.Name,
		screen: opt.Screen,
	}
	return room.
		initUI().
		initGame(&opt).
		initSeats()
}

func (room *Room) initSeats() *Room {
	for i, player := range room.game.Players {
		room.seats[i] = ui.NewSeat(ui.SeatOpt{
			Num:    i + 1,
			Player: player,
		})
	}

	return room
}

func (room *Room) initUI() *Room {
	room.msgBox = ui.NewMessage(room.screen.style)

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
	// init connection here
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
	seats := make([]string, len(room.seats))
	for i := range seats {
		seats[i] = room.seats[i].View()
	}
	t := table.New().Border(lipgloss.HiddenBorder())
	t.Row(seats[0:4]...)
	t.Row(seats[4], "", "", seats[5])
	t.Row(seats[6:10]...)
	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().
				Height(room.screen.style.GetHeight()-room.msgBox.Height()).
				Width(room.screen.style.GetWidth()).
				Border(lipgloss.NormalBorder()).
				Render(t.Render()),
			room.msgBox.View(),
		),
	)
	v.AltScreen = true
	return v
}
