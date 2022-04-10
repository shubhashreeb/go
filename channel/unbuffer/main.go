package main

import (
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
)

var data chan SomeStructWithTags

func main() {
	start := time.Now()

	done := make(chan bool)
	sendDone := make(chan bool)

	numberOfSendThread := 3
	numberOfReceiveThread := 20

	numberOfTask := 800000

	data = make(chan SomeStructWithTags, 6000)
	for i := 0; i < numberOfSendThread; i++ {
		go send(numberOfTask/numberOfSendThread, sendDone)
	}
	// Function to close the channel properly when all the send routines are done sending data
	go func() {
		cnt := 0
		for {
			select {
			case _, ok := <-sendDone:
				if ok {
					cnt++
					if cnt == numberOfSendThread {
						goto exit
					}
				}
			}
		}
	exit:
		close(data)
	}()

	for i := 0; i < numberOfReceiveThread; i++ {
		go receive(done)
	}

	for i := 0; i < numberOfReceiveThread; i++ {
		<-done
	}

	fmt.Println("Total lapse time is", time.Until(start).Milliseconds())
}

type SomeStructWithTags struct {
	FirstName       string `faker:"first_name"`
	FirstNameMale   string `faker:"first_name_male"`
	FirstNameFemale string `faker:"first_name_female"`
	LastName        string `faker:"last_name"`
	Name            string `faker:"name"`
	Password        string `faker:"password"`
	PhoneNumber     string `faker:"phone_number"`
	DomainName      string `faker:"domain_name"`
}

func send(size int, done chan bool) {
	for i := 0; i < size; i++ {
		a := SomeStructWithTags{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println("Running into err", err)
		}

		data <- a
	}
	done <- true
}

func receive(done chan bool) {
	for {
		select {
		case _, ok := <-data:
			if !ok {
				// fmt.Println("channel is closed", ok)
				goto exit
			} else {
				// fmt.Println("Received data is ", p)
			}
		}
	}
exit:
	done <- true
}
