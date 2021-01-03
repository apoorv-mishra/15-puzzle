package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"

	"github.com/mattn/go-runewidth"
)

type Cell struct {
	value rune
}

type Square struct {
	cells []*Cell
}

type Board struct {
	width  int
	height int
	rows   [][]*Square
}

func NewBoard(w, h int) *Board {
	board := &Board{
		width:  w,
		height: h,
	}
	letters := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', ' '}
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	for x, letterIdx := 0, 0; x < w; x++ {
		var squares []*Square
		for y := 0; y < h; y++ {
			sq := NewSquare(letters[letterIdx])
			squares = append(squares, sq)
			letterIdx++
		}
		board.rows = append(board.rows, squares)
	}
	return board
}

func NewSquare(letter rune) *Square {
	square := &Square{}
	x, y := 0, 0
	for sqy := y; sqy <= y+2; sqy++ {
		for sqx := x; sqx <= x+4; sqx++ {
			cell := &Cell{}
			if letter == ' ' {
				cell.value = letter
			} else {
				if sqx == x && sqy == y {
					cell.value = tcell.RuneULCorner
				} else if (sqx == x+1 || sqx == x+2 || sqx == x+3) && (sqy == y || sqy == y+2) {
					cell.value = tcell.RuneHLine
				} else if sqx == x+4 && sqy == y {
					cell.value = tcell.RuneURCorner
				} else if (sqx == x || sqx == x+4) && sqy == y+1 {
					cell.value = tcell.RuneVLine
				} else if sqx == x+2 && sqy == y+1 {
					cell.value = letter
				} else if sqx == x && sqy == y+2 {
					cell.value = tcell.RuneLLCorner
				} else if sqx == x+4 && sqy == y+2 {
					cell.value = tcell.RuneLRCorner
				}
			}
			square.cells = append(square.cells, cell)
		}
	}
	return square
}

func (c *Cell) draw(s tcell.Screen, x int, y int) {
	s.SetContent(x, y, c.value, []rune{}, tcell.StyleDefault)
}

func (sq *Square) draw(s tcell.Screen, x int, y int) {
	index := 0
	for cy := y; cy <= y+2; cy += 1 {
		for cx := x; cx <= x+4; cx += 1 {
			sq.cells[index].draw(s, cx, cy)
			index++
		}
	}
}

func (b *Board) draw(s tcell.Screen, x int, y int) {
	s.Clear()
	for sy := 0; sy < b.height; sy += 1 {
		for sx := 0; sx < b.width; sx += 1 {
			b.rows[sy][sx].draw(s, x+sx*5, y+sy*3)
		}
	}
	s.Show()
}

func swap(board *Board, x1 int, y1 int, x2 int, y2 int) {
	board.rows[x1][y1], board.rows[x2][y2] = board.rows[x2][y2], board.rows[x1][y1]
}

func getEmptySquare(b *Board) (int, int) {
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			if b.rows[i][j].cells[7].value == ' ' {
				return i, j
			}
		}
	}
	return -1, -1
}

func checkIfWon(b *Board) bool {
	winningOrder := [][]rune{
		[]rune{'A', 'B', 'C', 'D'},
		[]rune{'E', 'F', 'G', 'H'},
		[]rune{'I', 'J', 'K', 'L'},
		[]rune{'M', 'N', 'O', ' '},
	}
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			if b.rows[i][j].cells[7].value != winningOrder[i][j] {
				return false
			}
		}
	}
	return true
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func displayWinningMessage(s tcell.Screen, board *Board) {
	w, h := s.Size()
	emitStr(s, w/2-4, (h+board.height*3)/2+5, tcell.StyleDefault, "You Won!")
	emitStr(s, w/2-9, (h+board.height*3)/2+6, tcell.StyleDefault, "Press ESC to exit.")
	s.Show()
}

func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	board := NewBoard(4, 4)

	screenW, screenH := s.Size()
	board.draw(s, (screenW-board.width*5)>>1, (screenH-board.height*3)>>1)

	ei, ej := getEmptySquare(board)

	for {
		if checkIfWon(board) {
			displayWinningMessage(s, board)
		}
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			screenW, screenH = s.Size()
			board.draw(s, (screenW-board.width*5)>>1, (screenH-board.height*3)>>1)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			} else if !checkIfWon(board) {
				if ev.Key() == tcell.KeyUp {
					if ei+1 < board.height {
						swap(board, ei, ej, ei+1, ej)
						board.draw(s, (screenW-board.width*5)>>1, (screenH-board.height*3)>>1)
						ei = ei + 1
					}
				} else if ev.Key() == tcell.KeyDown {
					if ei-1 >= 0 {
						swap(board, ei, ej, ei-1, ej)
						board.draw(s, (screenW-board.width*5)>>1, (screenH-board.height*3)>>1)
						ei = ei - 1
					}
				} else if ev.Key() == tcell.KeyRight {
					if ej-1 >= 0 {
						swap(board, ei, ej, ei, ej-1)
						board.draw(s, (screenW-board.width*5)>>1, (screenH-board.height*3)>>1)
						ej = ej - 1
					}
				} else if ev.Key() == tcell.KeyLeft {
					if ej+1 < board.width {
						swap(board, ei, ej, ei, ej+1)
						board.draw(s, (screenW-board.width*5)>>1, (screenH-board.height*3)>>1)
						ej = ej + 1
					}
				}
			}
		}
	}
}
