package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

type Cell struct {
	value rune
}

type Square struct {
	cells []*Cell
}

type Board struct {
	width   int
	height  int
	startX  int
	startY  int
	squares []*Square
}

func NewBoard(w, h int) *Board {
	board := &Board{
		width:  w,
		height: h,
	}
	letter := rune(65)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var sq *Square
			if x == w-1 && y == h-1 {
				sq = NewSquare(' ')
			} else {
				sq = NewSquare(letter)
			}
			board.squares = append(board.squares, sq)
			letter++
		}
	}
	return board
}

func NewSquare(letter rune) *Square {
	square := &Square{}
	x, y := 0, 0
	for sqy := y; sqy <= y+2; sqy++ {
		for sqx := x; sqx <= x+4; sqx++ {
			cell := &Cell{}
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
	index := 0
	for sy := 0; sy < b.height; sy += 1 {
		for sx := 0; sx < b.width; sx += 1 {
			b.squares[index].draw(s, x+sx*5, y+sy*3)
			index++
		}
	}
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

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyUp {
				fmt.Println("Up")
			} else if ev.Key() == tcell.KeyDown {
				fmt.Println("Down")
			} else if ev.Key() == tcell.KeyRight {
				fmt.Println("Right")
			} else if ev.Key() == tcell.KeyLeft {
				fmt.Println("Left")
			}
		}
	}
}
