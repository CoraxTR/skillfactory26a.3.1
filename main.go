package main

import (
	"bufio"
	"fmt"
	"os"
	"skillfactory/20.2.1/buffer"
	"strconv"
	"sync"
	"time"
)

func filterNegativeInt(wg *sync.WaitGroup, done <-chan struct{}, input <-chan int) <-chan int {
	wg.Add(1)
	output := make(chan int)
	go func() {
		defer wg.Done()
		defer close(output)
		for {
			select {
			case <-done:
				return
			case v, isChannelOpen := <-input:
				if !isChannelOpen {
					return
				}
				if v >= 0 {
					select {
					case output <- v:
					case <-done:
						return

					}
				}
			}
		}
	}()
	return output
}

func filterZeroesAndNotDivideableByThree(wg *sync.WaitGroup, done <-chan struct{}, input <-chan int) <-chan int {
	output := make(chan int)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(output)
		for {
			select {
			case <-done:
				return
			case v, isChannelOpen := <-input:
				if !isChannelOpen {
					return
				}
				if v != 0 && v%3 == 0 {
					select {
					case output <- v:
					case <-done:
						return

					}
				}
			}
		}
	}()
	return output
}

func main() {
	b := buffer.CreateBuffer()
	input := make(chan int)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	bufferticker := time.NewTicker(b.GetBufferLifetime())
	defer bufferticker.Stop()

	buffermaterial := filterZeroesAndNotDivideableByThree(&wg, done, filterNegativeInt(&wg, done, input))

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-bufferticker.C:
				v, err := b.ReadAllFromBuffer()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("Полученные значения: ")
				fmt.Println(v)
				b.ClearBuffer()
			case v := <-buffermaterial:
				b.WriteToBuffer(v)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			close(done)
			break
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Неверный тип данных")
			continue
		}
		select {
		case input <- v:
		case <-done:
			break

		}
	}

	wg.Wait()
	os.Exit(0)

}
