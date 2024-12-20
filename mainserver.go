package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new TCP server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// Initialize the AOF file for persistence
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// Replay commands from the AOF log to restore state
	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command in AOF: ", command)
			return
		}

		handler(args)
	})

	// Listen for incoming client connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, aof) // Handle each connection in a separate goroutine
	}
}

func handleConnection(conn net.Conn, aof *Aof) {
	defer conn.Close()

	resp := NewResp(conn)
	writer := NewWriter(conn)

	for {
		// Read the client's request
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}
			fmt.Println("Error reading request: ", err)
			writer.Write(Value{typ: "error", str: "ERR invalid request"})
			continue
		}

		// Check if the request is a pipeline (array of commands)
		if value.typ == "array" && len(value.array) > 0 && value.array[0].typ == "array" {
			handlePipeline(value, writer, aof)
			continue
		}

		// Handle single command
		handleCommand(value, writer, aof)
	}
}

func handlePipeline(value Value, writer *Writer, aof *Aof) {
	responses := make([]Value, len(value.array))

	for i, commandValue := range value.array {
		if commandValue.typ != "array" || len(commandValue.array) == 0 {
			responses[i] = Value{typ: "error", str: "ERR invalid pipeline command format"}
			continue
		}

		command := strings.ToUpper(commandValue.array[0].bulk)
		args := commandValue.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			responses[i] = Value{typ: "error", str: "ERR unknown command"}
			continue
		}

		// Persist state-changing commands
		if command == "SET" || command == "HSET" {
			err := aof.Write(commandValue)
			if err != nil {
				fmt.Println("Error writing to AOF: ", err)
				responses[i] = Value{typ: "error", str: "ERR internal server error"}
				continue
			}
		}

		// Execute the command
		responses[i] = handler(args)
	}

	// Write all responses as a single array
	writer.Write(Value{typ: "array", array: responses})
}

var pubsub = NewPubSub() // Initialize PubSub

func handleCommand(value Value, writer *Writer, aof *Aof) {
	// Ensure the request is a valid RESP array
	if value.typ != "array" || len(value.array) == 0 {
		writer.Write(Value{typ: "error", str: "ERR invalid request format"})
		return
	}

	command := strings.ToUpper(value.array[0].bulk)
	args := value.array[1:]

	switch command {
	case "SUBSCRIBE":
		handleSubscribe(args, writer)
	case "UNSUBSCRIBE":
		handleUnsubscribe(args, writer)
	case "PUBLISH":
		handlePublish(args, writer)
	default:
		// Delegate to other command handlers
		handler, ok := Handlers[command]
		if !ok {
			writer.Write(Value{typ: "error", str: "ERR unknown command"})
			return
		}
		if command == "SET" || command == "HSET" {
			aof.Write(value) // Persist state-changing commands
		}
		result := handler(args)
		writer.Write(result)
	}
}

func handleSubscribe(args []Value, writer *Writer) {
	if len(args) < 1 {
		writer.Write(Value{typ: "error", str: "ERR wrong number of arguments for 'SUBSCRIBE' command"})
		return
	}

	for _, arg := range args {
		channel := arg.bulk
		sub := pubsub.Subscribe(channel)
		go func(channel string, sub <-chan Value) {
			for msg := range sub {
				writer.Write(Value{typ: "array", array: []Value{
					{typ: "bulk", bulk: "message"},
					{typ: "bulk", bulk: channel},
					msg,
				}})
			}
		}(channel, sub)
		writer.Write(Value{typ: "array", array: []Value{
			{typ: "bulk", bulk: "subscribe"},
			{typ: "bulk", bulk: channel},
		}})
	}
}

func handleUnsubscribe(args []Value, writer *Writer) {
	if len(args) < 1 {
		writer.Write(Value{typ: "error", str: "ERR wrong number of arguments for 'UNSUBSCRIBE' command"})
		return
	}

	for _, arg := range args {
		channel := arg.bulk
		pubsub.Unsubscribe(channel, nil) // Remove the subscriber
		writer.Write(Value{typ: "array", array: []Value{
			{typ: "bulk", bulk: "unsubscribe"},
			{typ: "bulk", bulk: channel},
		}})
	}
}

func handlePublish(args []Value, writer *Writer) {
	if len(args) < 2 {
		writer.Write(Value{typ: "error", str: "ERR wrong number of arguments for 'PUBLISH' command"})
		return
	}

	channel := args[0].bulk
	message := args[1]
	pubsub.Publish(channel, message)

	writer.Write(Value{typ: "integer", num: len(pubsub.subscribers[channel])})
}

