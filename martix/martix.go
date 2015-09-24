package martix

import (
	//"fmt"
	"math/rand"
	//"runtime"
	"time"
)

type row []int
type Martix []row

func Init_martix(l int) (Martix, error) {
	martix := make([]row, l)
	for index, _ := range martix {
		martix[index] = make(row, l)
	}
	return martix, nil
}

func (this Martix) Mirror() {
	for index, _ := range this {
		for i := 0; i < len(this[index])/2; i++ {
			swap(&this[index][i], &this[index][len(this[index])-i-1])
		}
	}
}

func (this Martix) Left90() {
	for x, row := range this {
		for y, _ := range row {
			if x > y {
				swap(&this[x][y], &this[y][x])
			}
		}
	}
}

func (this Martix) Right90() {
	for x, row := range this {
		for y, _ := range row {
			if x+y < len(this)-1 {
				swap(&this[x][y], &this[len(this)-y-1][len(this)-x-1])
			}
		}
	}
}

func (this Martix) Combin() bool {
	ch := make(chan bool, len(this))
	defer close(ch)
	change := false
	for _, row := range this {
		go row.combin_row(ch)
	}
	for i := 0; i < len(this); i++ {
		if <-ch == true {
			change = true
		}
	}
	return change
}

func (this row) combin_row(ch chan bool) {
	change := false
	table := make([]int, 0)
	for i, node := range this {
		if node != 0 {
			table = append(table, i)
		}
	}
	for {
		flag := len(table)
		for num, _ := range table {
			if len(table) == 1 || num == len(table)-1 {
				break
			}
			if this[table[num]] == this[table[num+1]] {
				this[table[num]] *= 2
				this[table[num+1]] = 0
				table = append(table[:num+1], table[num+2:]...)
				change = true
				break
			}
		}
		if len(table) == flag {
			break
		}
	}
	i := 0
	for num, _ := range table {
		this[i] = this[table[num]]
		if i != table[num] {
			change = true
			this[table[num]] = 0
		}
		i++
	}
	ch <- change
}

//func main() {
//runtime.GOMAXPROCS(runtime.NumCPU())
//a, _ := Init_martix(4)
//a[0][0] = 2
//a[0][1] = 2
//a[0][2] = 4
//a[0][3] = 8
//a[1][0] = 2
//a[1][1] = 6
//a[1][2] = 7
//a[1][3] = 8
//a[2][0] = 4
//a[2][1] = 4
//a[2][2] = 11
//a[2][3] = 12
//a[3][0] = 8
//a[3][1] = 14
//a[3][2] = 15
//a[3][3] = 16
//fmt.Println(a)
//a.Left90()
//fmt.Println(a)
//a.Right90()
//fmt.Println(a)
//a.Combin()
//fmt.Println(a)
//a.AddNum()
//fmt.Println(a)
//}

func swap(a, b *int) {
	tmp := *b
	*b = *a
	*a = tmp
}

func (this Martix) AddNum(num int) {
	table := make([]int, 0)
	for x, row := range this {
		for y, _ := range row {
			if this[x][y] == 0 {
				table = append(table, x*len(this)+y)
			}
		}
	}

	ch := make(chan bool, 2)
	defer close(ch)
	for i := 0; i < num; i++ {
		go func() {
			ran := rand.Intn(len(table))
			this[table[ran]/len(this)][table[ran]%len(this)] = 2 << (rand.Uint32() % 2)
			ch <- true
		}()
	}
	for i := 0; i < num; i++ {
		<-ch
	}
}

func Init() {
	rand.Seed(time.Now().UnixNano())
}
