package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

/*
type Car struct {
	model string
}

func (c Car) PrintModel() {
	fmt.Println(c.model)
}
func main() {
	c := Car{model: "new model"}
	defer c.PrintModel()
	c.model = "second model"
	c.PrintModel()

}
*/
/* output
new model
*/
/*
type Slice []int

func NewSlice() Slice {
	return make(Slice, 0)
}
func (s *Slice) add(elem int) *Slice {
	*s = append(*s, elem)
	fmt.Print(elem)
	return s
}
func main() {
	s := NewSlice()
	defer s.add(1).add(2)
	s.add(3)
}
*/
/*132*/
/*
func main() {
	strs := []string{"1", "2", "3"}
	for _, s := range strs {
		go func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("%s", s)
		}()
	}
	time.Sleep(3 * time.Second)
}
*/
/*
333*/
/*
func app() func(string) string {
	t := "Hi"
	c := func(b string) string {
		t = t + " " + b
		return t
	}
	return c
}
func main() {
	a := app()
	b := app()
	a("go")
	fmt.Println(b("All"))
}
*/
/*
func main() {
	v := []int{1, 2, 3}
	for i := range v {
		fmt.Printf("%v\n", i)
		v = append(v, i)
	}
}
*/
/* 0,1,2*/
/*
func main() {
	nums := []int{1, 2, 3, 4, 5}
	sum := 0
	for i, n := range nums {
		i = 6
		sum += n
	}
	fmt.Println(sum)

}
*/
/*couldn't complie */
/*
func main() {
	int_chain := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chain <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chain:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}
}
*/
/*
panic
*/
/*
func main() {
	ch1 := make(chan int)
	go fmt.Println(<-ch1)
	ch1 <- 5
	time.Sleep(1 * time.Second)
}
*/
/*dead lock */
/*
func main() {
	str1 := []string{"a", "b", "c"}
	str2 := str1[1:]
	str2[1] = "new"
	fmt.Println(str1)
	str2 = append(str2, "z", "x", "y")
	fmt.Println(str1)
}
*/
/*
[a b new]
[a b new]
*/
/*
declare a pointer which point to an interger array
func main() {
	var p_int *[4]int

	myArray := [4]int{1, 2, 3, 4}
	p_int = &myArray
	fmt.Println(*p_int)

}
*/
/*declar an array of interger pointer
func main() {
	var p_array [4]*int
	elem1 := 1
	p_array[0] = &elem1
	fmt.Println(p_array)
}
*/
/*
func main() {
	var p_int *int
	i := 10
	p_int = &i
	fmt.Println(*p_int)
	var pp_int **int
	pp_int = &p_int
	fmt.Println(**pp_int)
}

*/

func main() {
	n := 1
	sub(n)
	time.Sleep(2 * time.Minute)
}
func sub(n int) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		timer := time.NewTimer(10 * time.Second)
		log.Printf("subfunction start %d", n)
		defer func() {
			timer.Stop()
			if r := recover(); r != nil {
				n = n + 1
				sub(n)
			}
			log.Printf("subfunction finished")
		}()
		for {
			select {
			case <-ticker.C:
				log.Printf("subfunction is running")
			case <-timer.C:
				panic(errors.New("Raise panic"))

			}
		}
	}()
}
