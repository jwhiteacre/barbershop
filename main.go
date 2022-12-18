package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const WaitRoomSize = 6
const NumberBarbers = 1

// Customer struct for defining the customer and the service
type Customer struct {
	service string
}

// Barber struct empty for now but could allow for an abstraction
type Barber struct {
}

// PerformaService is what the barber does
func (b Barber) PerformService(c Customer) {
	fmt.Printf("-- Performing Service %s\n", c.service)
	time.Sleep(5 * time.Second)
	fmt.Printf("-- Done with Service %s\n", c.service)
}

// BarberShop struct for holding the waitroom
type BarberShop struct {
	waitRoom chan Customer
}

// Puts a customer in the waitroom as long as it is not full
func (b *BarberShop) AddCustomer(c Customer) error {

	// Try to add customer to wait room
	select {
	case b.waitRoom <- c:
	default:
		return fmt.Errorf("WaitRoom is full, hire more barbers")
	}
	return nil
}

// Opens the barbershop
func (b BarberShop) Open() {

	w := Barber{}
	for {
		c, ok := <-b.waitRoom
		if !ok {
			break
		}
		w.PerformService(c)
	}
}

// CLoses the barbershop
func (b BarberShop) Close() {
	close(b.waitRoom)
}

func main() {

	// Open the barbershop
	b := &BarberShop{waitRoom: make(chan Customer, WaitRoomSize)}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Open()
	}()

	// Have 10 customes show up randomly
	var cwg sync.WaitGroup
	for i := 0; i < 10; i++ {
		cwg.Add(1)
		go func(c Customer) {
			defer cwg.Done()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			e := b.AddCustomer(c)
			if e != nil {
				fmt.Println(e)
			}
		}(Customer{service: "haircut"})
	}

	// Wait for last customer to get into wait room then close the shop
	cwg.Wait()
	fmt.Println("Closing the barber shop")
	b.Close()

	fmt.Println("Waiting for last customers")
	wg.Wait()
}
