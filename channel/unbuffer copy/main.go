package main

import (
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
)

var data chan string

func main() {
	start := time.Now()
	done := make(chan bool)
	data = make(chan string)
	go send(100000)
	go send(100000)
	go send(100000)
	go send(100000)
	go send(100000)

	go receive(done)
	go receive(done)
	go receive(done)
	go receive(done)
	go receive(done)

	for i := 0; i < 5; i++ {
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

func send(size int) {
	for i := 0; i < size; i++ {
		a := SomeStructWithTags{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println("Running into err", err)
		}

		data <- a.Name
	}
	// names := []string{"Aaron", "Alice", "Bob", "Carol", "Dan"}
	// for _, name := range names {
	// 	data <- name
	// }
	// close(data)
}

func receive(done chan bool, total int) {
	cnt := 0
	for {
		select {
		case name, ok := <-data:
			if !ok {
				// fmt.Println("channel is closed", ok)
				goto exit
			} else {
				fmt.Println("Received data is ", name, ok)
				cnt++
				if cnt == total {
					goto exit
				}
			}
		}
	}
exit:
	done <- true
}
