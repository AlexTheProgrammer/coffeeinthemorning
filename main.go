package main

import (
	"fmt"
	"time"
)

func main() {
	// in the morning I
	// get up.
	getUp()
	// turn on the kettle
	isBoiledChan := turnOnKettle()
	// put the shower on
	isWaterHotChan := turnOnTheShower()
	// think for a second
	time.Sleep(2 * time.Second)
	// grind the coffee
	// if the kettle has boiled, fill the perculator
	filled := fillThePerculatorIfReady(isBoiledChan)
	// wait for the shower to heat up and then get in
	// shower for at least 10mins (we can simulate with 5 seconds)
	waitForShowerAndGetIn(isWaterHotChan)
	// if the perculator is not filled up wait for the kettle to boil and fill up
	if !filled {
		waitForKettleAndFillPerculator(isBoiledChan)
	}
	// go to computer
	goToComputer()
}

func goToComputer() {
	fmt.Println("Going to Computer to start my day")
}

func waitForKettleAndFillPerculator(isBoiledChan chan bool) {
	// wait for kettle to boil
	<-isBoiledChan
	fmt.Println("Filling up the perculator")
}

func waitForShowerAndGetIn(isWaterHotChan chan bool) {
	// wait for water to be hot
	<-isWaterHotChan
	fmt.Println("getting into shower")
	showerTime := 5 * time.Second
	time.Sleep(showerTime)
}

// getUp prints get up action
func getUp() {
	fmt.Println("I got up")
}

// turnOnTheShower opens chan that will recieve when shower is hot
func turnOnTheShower() chan bool {
	fmt.Println("turning on shower")
	isWaterHotChan := make(chan bool)
	heatUpTime := 3
	go logThenWriteChan(isWaterHotChan, heatUpTime, "shower getting hotter", "SHOWER IS HOT")
	return isWaterHotChan
}

// turnOnKettle opens chan that will recieve when kettle is boiled
// kettle is boiled
func turnOnKettle() chan bool {
	fmt.Println("turning on kettle")
	isBoiledChan := make(chan bool)
	boilSeconds := 5
	go logThenWriteChan(isBoiledChan, boilSeconds, "Kettle whistles", "KETTLE BOILED")
	return isBoiledChan
}

// logThenWriteChan takes a chan and will send to it after maxTime has passed
// it also starts a ticker that will log to console whiles waiting to maxTime
func logThenWriteChan(aChan chan bool, maxTime int, processingMessage, readyMessage string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	canExit := make(chan bool)
	go func(aChan, canExit chan bool) {
		for {
			select {
			case res := <-aChan:
				// reload the chan for later use
				aChan <- res
				return
			case <-ticker.C:
				fmt.Println(processingMessage)
			}
		}
	}(aChan, canExit)
	// wait for kettle to boil
	time.Sleep(time.Duration(maxTime) * time.Second)
	ticker.Stop()
	fmt.Println(readyMessage)
	aChan <- true
}

// fillThePerculatorIfReady will return true if the isBoiledChan has a value ready
func fillThePerculatorIfReady(isBoiledChan chan bool) bool {
	// this will return if the chan does not have a value to send
	select {
	case <-isBoiledChan:
		fmt.Println("Filling the perculator")
		return true
	default:
		fmt.Println("Kettle still boiling, will come back")
		return false
	}
}
