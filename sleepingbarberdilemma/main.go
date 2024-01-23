package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numBarbers   = 2
	numChairs    = 5
	closingTime  = 10 * time.Second
	customerTime = 2 * time.Second
)

type BarberShop struct {
	barbers        []*Barber
	waitingRoom    []int
	barberMutex    sync.Mutex
	customerSem    chan int
	barberSem      chan int
	wg             sync.WaitGroup
	closeShopCh    chan struct{}
}

type Barber struct {
	id int
}

func (bs *BarberShop) openShop() {
	for i := 0; i < numBarbers; i++ {
		bs.wg.Add(1)
		barber := &Barber{id: i + 1}
		go bs.barberWork(barber)
	}

	go bs.closeShop()

	// Start closing timer
	go func() {
		time.Sleep(closingTime)
		close(bs.closeShopCh)
	}()

	bs.wg.Wait() // Wait for all barbers to finish
}

func (bs *BarberShop) closeShop() {
	fmt.Println("Closing the shop. Waiting for customers to finish.")
	bs.barberMutex.Lock()
	close(bs.customerSem) // Release any sleeping barber
	bs.barberMutex.Unlock()

	bs.wg.Wait() // Wait for all barbers to finish

	fmt.Println("Shop is closed.")
}

func (bs *BarberShop) barberWork(barber *Barber) {
	defer bs.wg.Done()

	for {
		select {
		case <-bs.customerSem: // Sleep if no customers
			bs.barberMutex.Lock()
			customer := -1
			if len(bs.waitingRoom) > 0 {
				customer, bs.waitingRoom = bs.waitingRoom[0], bs.waitingRoom[1:]
			}
			bs.barberMutex.Unlock()

			if customer != -1 {
				bs.cutHair(barber, customer)
			} else {
				fmt.Printf("Barber %d falls asleep.\n", barber.id)
				bs.barberSem <- barber.id
				return
			}
		case <-bs.closeShopCh:
			fmt.Printf("Barber %d goes home.\n", barber.id)
			bs.barberSem <- barber.id
			return
		}
	}
}

func (bs *BarberShop) cutHair(barber *Barber, customer int) {
	fmt.Printf("Barber %d is cutting hair for customer %d.\n", barber.id, customer)
	time.Sleep(customerTime)
	fmt.Printf("Barber %d finished cutting hair for customer %d.\n", barber.id, customer)
}

func customerArrives(bs *BarberShop, customer int) {
	bs.barberMutex.Lock()
	defer bs.barberMutex.Unlock()

	if len(bs.waitingRoom) < numChairs {
		fmt.Printf("Customer %d arrives and sits in the waiting room.\n", customer)
		bs.waitingRoom = append(bs.waitingRoom, customer)
		bs.customerSem <- customer
	} else {
		fmt.Printf("Waiting room is full. Customer %d leaves.\n", customer)
	}
}

func main() {
	bs := &BarberShop{
		barbers:        make([]*Barber, numBarbers),
		waitingRoom:    make([]int, 0, numChairs),
		barberMutex:    sync.Mutex{},
		customerSem:    make(chan int, numChairs),
		barberSem:      make(chan int),
		closeShopCh:    make(chan struct{}),
	}

	for i := range bs.barbers {
		bs.barbers[i] = &Barber{id: i + 1}
	}

	go bs.openShop()

	// Simulate customers arriving at random intervals
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(randInt(1, 3)) * time.Second)
		bs.wg.Add(1)
		go customerArrives(bs, i+1)
	}

	// Wait for the shop to close
	<-bs.closeShopCh
	close(bs.barberSem)
	bs.wg.Wait() // Ensure all goroutines finish

	fmt.Println("All customers finished.")
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
