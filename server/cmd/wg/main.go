package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	errs := make(chan error, 1)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go doStuff(i, errs, &wg)
	}
	wg.Wait()
	fmt.Println("finished")
	close(errs)
	for err := range errs {
		fmt.Println(err)
	}
}

func doStuff(id int, errs chan error, wg *sync.WaitGroup) {
	time.Sleep(time.Second)
	if id == 3 {
		errs <- fmt.Errorf("error from %d", id)
	}
	fmt.Printf("hey from %d\n", id)
	wg.Done()
}
