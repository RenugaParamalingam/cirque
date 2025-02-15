package cirque

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestCirque(t *testing.T) {
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedOutput := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

	var measuredParallelism int64 = 0

	var maxParallelism int64 = 3
	inputChannel, outputChannel := NewCirque(maxParallelism, func(i interface{}) interface{} {
		atomic.AddInt64(&measuredParallelism, 1)
		time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
		atomic.AddInt64(&measuredParallelism, -1)
		return i.(int) * 2
	})

	go func() {
		for _, i := range inputs {
			inputChannel <- i
			if atomic.LoadInt64(&measuredParallelism) > maxParallelism {
				t.Error("SO MUCH CANNOT ABLE TO HANDLE!")
			}
		}
		close(inputChannel)
	}()

	var actualOutput []int
	for i := range outputChannel {
		actualOutput = append(actualOutput, i.(int))
	}
	if !reflect.DeepEqual(expectedOutput, actualOutput) {
		t.Error("WRONG WRONG WRONG")
		t.Log(expectedOutput)
		t.Log(actualOutput)
	}
}

func ExampleNewCirque() {
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	inputChannel, outputChannel := NewCirque(3, func(i interface{}) interface{} {
		time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
		return i.(int) * 2
	})

	go func() {
		for _, i := range inputs {
			inputChannel <- i
		}
		close(inputChannel)
	}()

	var output []int
	for i := range outputChannel {
		output = append(output, i.(int))
	}
	fmt.Println(output)

	// Output: [2 4 6 8 10 12 14 16 18 20]
}
