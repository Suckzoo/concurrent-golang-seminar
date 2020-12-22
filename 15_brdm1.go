package main

import (
	"fmt"
	"math/rand"
	"time"
)

func runStep(stepName string) {
	duration := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(duration)
	fmt.Println("stepName: ", stepName, " elapsed: ", duration, "ms")
}
func main() {
	now := time.Now()
	runStep("deployBroserv")
	runStep("deployBroapp")
	runStep("deployBromon")
	runStep("deployPAS")
	elapsed := time.Now().Sub(now)
	fmt.Println(elapsed)
}
