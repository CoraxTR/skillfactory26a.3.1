package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"skillfactory/20.2.1/buffer"
	"strconv"
	"sync"
	"time"
)

func filterNegativeInt(wg *sync.WaitGroup, done <-chan struct{}, input <-chan int) <-chan int {
	log.Printf("[INFO] Starting filtering negative ints")
	wg.Add(1)
	output := make(chan int)
	go func() {
		defer wg.Done()
		defer close(output)
		defer log.Printf("[INFO] Stopping filtering negative ints")
		for {
			select {
			case <-done:
				return
			case v, isChannelOpen := <-input:
				if !isChannelOpen {
					return
				}
				log.Printf("[INFO] Filtering %d as a negative int", v)
				if v >= 0 {
					select {
					case output <- v:
						log.Printf("[INFO] %d went through filtering negative ints", v)
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
	log.Printf("[INFO] Starting filtering zeroes and not-dividable by 3 ints")
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(output)
		defer log.Printf("[INFO] Stopping filtering zeroes and not-dividable by 3 ints")
		for {
			select {
			case <-done:
				return
			case v, isChannelOpen := <-input:
				if !isChannelOpen {
					return
				}
				log.Printf("[INFO] Filtering %d as a zero or not-dividable by 3 int", v)
				if v != 0 && v%3 == 0 {
					select {
					case output <- v:
						log.Printf("[INFO] %d went through filtering zeroes and not-dividable by 3 ints", v)
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
	log.Printf("[INFO] Starting pipeline")
	b := buffer.CreateBuffer()
	input := make(chan int)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	log.Printf("[INFO] Starting buffer clean-up ticker")
	bufferticker := time.NewTicker(b.GetBufferLifetime())
	defer bufferticker.Stop()

	buffermaterial := filterZeroesAndNotDivideableByThree(&wg, done, filterNegativeInt(&wg, done, input))

	wg.Add(1)
	go func() {
		log.Printf("[INFO] Starting bufferfilter")
		defer wg.Done()
		defer log.Printf("[INFO] Stopping bufferfilter")
		for {
			select {
			case <-done:
				return
			case <-bufferticker.C:
				log.Printf("[INFO] Buffer clean-up timer went off")
				v, err := b.ReadAllFromBuffer()
				if err != nil {
					log.Printf("[ERROR] Error reading from buffer")
					fmt.Println(err)
				}
				fmt.Print("Полученные значения: ")
				fmt.Println(v)
				b.ClearBuffer()
			case v := <-buffermaterial:
				log.Printf("[INFO] %d went through bufferfilter", v)
				b.WriteToBuffer(v)
			}
		}
	}()
	log.Printf("[INFO] Waiting for user input")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			log.Printf("[INFO] Pipeline shut-down dur to user command")
			close(done)
			break
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("[ERROR] Error parsing %v to int", v)
			fmt.Println("Неверный тип данных")
			continue
		}
		select {
		case input <- v:
			log.Printf("[INFO] %d goes to pipeline", v)
		case <-done:
			break

		}
	}

	wg.Wait()
	log.Printf("[INFO] Pipeline shut down")
	os.Exit(0)

}
