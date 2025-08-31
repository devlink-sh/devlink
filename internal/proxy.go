package internal

import (
	"io"
	"log"
	"net"
	"sync"
)

func Pipe(clientConnection, targetConnection net.Conn){
	defer clientConnection.Close()
	defer targetConnection.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := io.Copy(targetConnection, clientConnection); err != nil {
			log.Printf("Error copying from client to target: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(clientConnection, targetConnection); err != nil {
			log.Printf("Error copying from target to client: %v", err)
		}
	}()

	wg.Wait()
}

