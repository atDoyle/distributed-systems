package main

import (
    "fmt"
    "log"
    "net"

    "github.com/atDoyle/distributed-systems/key-value-store" // Import your generated protobuf code
    "google.golang.org/protobuf/proto"
)

// In-memory key-value store
var store = make(map[string]string)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Client connected:", conn.RemoteAddr())

	func handleConnection(conn net.Conn) {
		defer conn.Close()
		fmt.Println("Client connected:", conn.RemoteAddr())
	
		for {
			// 1. Read the length of the incoming message
			var length int32
			err := binary.Read(conn, binary.BigEndian, &length)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("Client disconnected.")
					return
				}
				log.Printf("Error reading message length: %v", err)
				return
			}
	
			// 2. Read the message data
			buffer := make([]byte, length)
			_, err = io.ReadFull(conn, buffer)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
	
			// 3. Deserialize the Command message
			var command kvstore.Command
			err = proto.Unmarshal(buffer, &command)
			if err != nil {
				log.Printf("Error unmarshaling command: %v", err)
				return
			}
	
			// 4. Process the command
			var response *kvstore.Command
			switch req := command.Request.(type) {
			case *kvstore.Command_Set:
				store[req.Set.Key] = req.Set.Value
				response = &kvstore.Command{
					Response: &kvstore.Command_SetResponse{
						SetResponse: &kvstore.SetResponse{Success: true},
					},
				}
			case *kvstore.Command_Get:
				value, found := store[req.Get.Key]
				response = &kvstore.Command{
					Response: &kvstore.Command_GetResponse{
						GetResponse: &kvstore.GetResponse{Value: &value, Found: found},
					},
				}
			default:
				log.Println("Received unknown command.")
				response = &kvstore.Command{
					Response: &kvstore.Command_SetResponse{
						SetResponse: &kvstore.SetResponse{Success: false}, // Indicate failure for unknown command
					},
				}
			}
	
			// 5. Serialize the response
			responseBytes, err := proto.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				return
			}
	
			// 6. Write the length of the response
			responseLength := int32(len(responseBytes))
			err = binary.Write(conn, binary.BigEndian, &responseLength)
			if err != nil {
				log.Printf("Error writing response length: %v", err)
				return
			}
	
			// 7. Write the response data
			_, err = conn.Write(responseBytes)
			if err != nil {
				log.Printf("Error writing response: %v", err)
				return
			}
		}
	}
}

func main() {
    listener, err := net.Listen("tcp", ":8080") // Listen on port 8080
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    defer listener.Close()

    fmt.Println("Server listening on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        go handleConnection(conn) // Handle each connection in a new goroutine
    }
}