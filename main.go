package main 

import (
	"fmt"
	
)



func main() {
	server := NewAPIServer(":3000") // runs on 3000
	server.Run()    
	
	// call run
	fmt.Println("Hi Buddy!!, Server is running")
}