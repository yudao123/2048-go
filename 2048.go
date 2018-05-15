package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/xuzhenglun/2048-Go/martix"
	"math"
	"strconv"
	"time"
)

const MAX_LEN int = 4
const Add_NUM int = 1

var step int
var output_mode = termbox.OutputNormal

var colorTable = [...]termbox.Attribute{
	termbox.ColorMagenta,
	termbox.ColorGreen,
	termbox.ColorBlue,
	termbox.ColorCyan,
	termbox.ColorYellow,
	termbox.ColorBlack,
	termbox.ColorMagenta}

type Go2048 struct {
	martix.Martix
}

func (this Go2048) GoUp() bool {
	this.Left90()
	change := this.Combin()
	this.Left90()
	return change
}

func (this Go2048) GoDown() bool {
	this.Right90()
	change := this.Combin()
	this.Right90()
	return change
}

func (this Go2048) GoLeft() bool {
	change := this.Combin()
	return change
}

func (this Go2048) GoRight() bool {
	this.Mirror()
	change := this.Combin()
	this.Mirror()
	return change
}

func (this Go2048) CheckWinOrLose() bool {
	for x, row := range this.Martix {
		for y, _ := range row {
			if this.Martix[x][y] == 0 {
				return true
				//true = Have not been dead yet
			}
		}
	}
	return false
	//false = Lose
}

func (this Go2048) Init_termbox(x, y int) error {
	fg := termbox.ColorYellow
	bg := termbox.ColorBlack
	err := termbox.Clear(fg, bg)
	if err != nil {
		return err
	}
	str := "Enter: restart game"
	for n, c := range str {
		termbox.SetCell(x+n, y-1, c, fg, bg)
	}

	str = "ESC: quit game" + "  Step: " + strconv.Itoa(step)
	for n, c := range str {
		termbox.SetCell(x+n, y-2, c, fg, bg)
	}

	str = "Play with Arrow Key"
	for n, c := range str {
		termbox.SetCell(x+n, y-3, c, fg, bg)
	}

	fg = termbox.ColorBlack
	bg = termbox.ColorGreen
	for i := 0; i <= len(this.Martix); i++ {
		for t := 0; t < 6*len(this.Martix); t++ {
			if t%6 != 0 {
				termbox.SetCell(x+t, y+i*2, '-', fg, bg)
			}
		}
		for t := 0; t <= 2*len(this.Martix); t++ {
			if t%2 == 0 {
				termbox.SetCell(x+i*6, y+t, '+', fg, bg)
			} else {
				termbox.SetCell(x+i*6, y+t, '|', fg, bg)
			}
		}
	}

	for i, row := range this.Martix {
		for j, _ := range row {
			if this.Martix[i][j] > 0 {
				str := fmt.Sprintf("%-5d", this.Martix[i][j])
				for n, char := range str {
					if output_mode == termbox.Output256 {
						termbox.SetCell(x+j*6+1+n, y+i*2+1, char, 0x10+termbox.Attribute(this.Martix[i][j]%256), 0xe0-termbox.Attribute(this.Martix[i][j]*2%256))
					} else {
						termbox.SetCell(x+j*6+1+n, y+i*2+1, char, termbox.ColorWhite, colorTable[int(math.Log2(float64(this.Martix[i][j])))%7])
					}
				}
			}
		}
	}
	return termbox.Flush()
}

func converPrintStr(x, y int, str string, fg, bg termbox.Attribute) error {
	xx := x
	for n, c := range str {
		if c == '\n' {
			y++
			xx = x - n - 1
		}
		termbox.SetCell(xx+n, y, c, fg, bg)
	}
	return termbox.Flush()
}

func (t *Go2048) ListernKey() chan termbox.Event {
	//ev := termbox.PollEvent()
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent() // 开始监听键盘事件
		}
	}()
	return event_queue
}

func (t *Go2048) ActionAndReturnKey(event_queue chan termbox.Event) termbox.Key {
	for {
		ev := <-event_queue
		changed := false

		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				changed = t.GoUp()
			case termbox.KeyArrowDown:
				changed = t.GoDown()
			case termbox.KeyArrowLeft:
				changed = t.GoLeft()
			case termbox.KeyArrowRight:
				changed = t.GoRight()
			case termbox.KeyEsc, termbox.KeyEnter:
				changed = true
			default:
				changed = false
			}

			// 如果元素的值没有任何更改，则从新开始循环
			if !changed && t.CheckWinOrLose() {
				continue
			}

		case termbox.EventResize:
			x, y := termbox.Size()
			t.Init_termbox(x/2-10, y/2-4)
			continue
		case termbox.EventError:
			panic(ev.Err)
		}
		return ev.Key
	}
}

func main() {
    
    fmt.Println("welcome to 2048 game , forked by yinghua 1707")
    
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	x, y := termbox.Size()

	output_mode = termbox.SetOutputMode(termbox.Output256)

	martix.Init()
	var t Go2048
	t.Martix, _ = martix.Init_martix(MAX_LEN)
	t.AddNum(Add_NUM)
	step = 0
	ch := t.ListernKey()
	defer close(ch)

	for {
		t.Init_termbox(x/2-10, y/2-4)

		key := t.ActionAndReturnKey(ch)

		if t.CheckWinOrLose() == false {
			str := "Lose!"
			strlen := len(str)
			converPrintStr(x/2-strlen/2, y/2, str, termbox.ColorBlack, termbox.ColorRed)
			for {
				key = t.ActionAndReturnKey(ch)
				if key == termbox.KeyEnter || key == termbox.KeyEsc {
					break
				}
			}
		}

		if key == termbox.KeyEnter {
			t.Martix, _ = martix.Init_martix(MAX_LEN)
			step = -1
		}
		if key == termbox.KeyEsc {
			return
		}

		step++

		t.Init_termbox(x/2-10, y/2-4)
		time.Sleep(500 * time.Millisecond)
		t.AddNum(Add_NUM)
	}
}
