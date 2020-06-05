package readline

import (
	"errors"
	"fmt"
	"internal/fancy_printer"

	"github.com/loremdipso/go_utils"

	"github.com/nsf/termbox-go"
)

type Reader struct {
	History   []string
	completer func(line string) []string
}

func NewReadline() *Reader {
	termbox.Init()
	return &Reader{History: make([]string, 0)}
}

func (reader *Reader) Close() {
	// TODO: is this necessary every time? And are we even going to have multiple times?
	// termbox.Close()
}

func (reader *Reader) SetCompleter(completer func(string) []string) {
	reader.completer = completer
}

func (reader *Reader) Readline(prompt string) (string, error) {
	/*
		Strategy: listen for keyboardinputs
	*/
	// ioreader := bufio.NewReader(os.Stdin)
	var query string
	var history_position = len(reader.History)
	var cursor_position = 5

	// termbox.Init()
	// defer termbox.Close()

	// var cursor
	term_height := 5
	for {
		// TODO: something smarter than this
		fmt.Printf("\r\r\r\r\r\r\r\r\r%s", query)
		fmt.Println(term_height)
		termbox.SetCursor(term_height, cursor_position)
		// termbox.GetCursor()
		// termbox.SetCursor()

		// cursor_position := len(query)
		event := termbox.PollEvent()
		//termbox.SetCursor(event.Height-1, cursor_position)
		if event.Key > 0 {
			switch event.Key {
			case termbox.KeyCtrlU:
				query = ""

			case termbox.KeyArrowUp:
				history_position -= 1
				history_position = go_utils.Max(0, history_position-1)
				if len(reader.History) > history_position {
					query = reader.History[history_position]
				}

			case termbox.KeyArrowDown:
				history_position += 1
				history_position = go_utils.Min(len(reader.History), history_position)
				if len(reader.History) > history_position {
					query = reader.History[history_position]
				}

			case termbox.KeyEnter:
				return query, nil

			case termbox.KeyBackspace:
				// TODO: implement this and cursor history_position
				if len(query) > 0 {
					query = query[0 : len(query)-1]
				}
			case termbox.KeyCtrlC:
				return "", errors.New("Aborted")
			}
		} else if event.Ch > 0 {
			query += string(event.Ch)
		}

		// c, _, err := ioreader.ReadRune()
		// if err != nil {
		// 	fmt.Println("uh-oh", err)
		// 	return "", err
		// }
		// fmt.Println("Character: %d", c)
		// if key, ok := next.(rune); ok {
		// 	if key == tab {
		// 		direction = tabForward
		// 		continue
		// 	}
		// 	if key == esc {
		// 		return line, pos, rune(esc), nil
		// 	}
		// }
		// pick. err :=
	}
	// return fake()
}

func fake() (string, error) {
	tags := []string{"a", "b", "c"}
	fancy_printer.Print(tags, false, false)
	return "q", nil
}
