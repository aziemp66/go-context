package gocontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f")) //dapat
	fmt.Println(contextF.Value("c")) //dapat milik parent
	fmt.Println(contextF.Value("b")) //tidak dapat, beda parent
	fmt.Println(contextA.Value("b")) //tidak dapat mengambil child
}

func CreateCounter() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()

	return destination
}

func TestWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine: ", runtime.NumGoroutine())

	destination := CreateCounter()

	for n := range destination {
		fmt.Println("Counter: ", n)
		if n == 10 {
			break
		}
	}

	fmt.Println("Final Goroutine: ", runtime.NumGoroutine())
}

func CreateCounterWithContext(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestWithCancelContext(t *testing.T) {
	fmt.Println("Total Goroutine: ", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounterWithContext(ctx)

	fmt.Println("Total Goroutine: ", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter: ", n)
		if n == 10 {
			break
		}
	}
	cancel() //mengirim sinyal ke context

	time.Sleep(1 * time.Second)

	fmt.Println("Final Goroutine: ", runtime.NumGoroutine())
}
