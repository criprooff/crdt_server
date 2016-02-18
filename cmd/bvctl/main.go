// Package main provides a client for the biolviewity server
package main

import (
	pb "github.com/kitschysynq/biolviewity/setserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSetClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func add(i string) {
	present, err := c.Add(i)
	if err != nil {
		fmt.Printf("Error adding %q: %s")
	}
	if c.Add(i) == true {
		fmt.Printf("Added %q\n", i)
	} else {
	}
}
