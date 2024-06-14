package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	// "sync"
)

// handleConnection takes a connection, containerId, and port. It runs `socat` on the docker container and manages
// forwarding data between the intup conn and the running `socat` process.
func handleConnection(conn net.Conn, containerId string, port int) {
	defer conn.Close()

	// Define the command to execute
	// TODO: Error handling doesn't really work, as it will just sit and hang even on a bad command.
	// I need to do something to detect when a bad exit code happens
	cmd := exec.Command("docker", "exec", "-i", containerId,
		"bash", "-c", fmt.Sprintf("su - root -c 'socat - TCP:localhost:%v'", port))

	// Get the stdin and stdout pipes of the command
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	// Create a channel to signal when done
	done := make(chan struct{})

	go func() {
		io.Copy(cmdStdin, conn)
		cmdStdin.Close()
		done <- struct{}{}
	}()
	go func() {
		io.Copy(conn, cmdStdout)
		done <- struct{}{}
	}()

	// Wait for both to complete
	<-done
	<-done

	// Wait for the command to exit
	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

func start_listener(stop chan struct{}, containerId string, port int) {
	// TODO: Should probably do something to check that the containerId is valid

	log.Printf("Starting listener on port %v", port)

	// Listen on TCP port 3000
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}
	defer listener.Close()

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)

			// Signal all listeners to stop
			stop <- struct{}{}

			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, containerId, port)
	}
}

func main() {

	portsPtr := flag.String("p", "", "A port or list of ports. Can be like: '8080', '8080:1234', or a comma-seperated list")
	containerIdPtr := flag.String("c", "", "The ID of the docker container")
	flag.Parse()

	if *portsPtr == "" {
		panic("You must specify a port, or list of ports")
	}
	if *containerIdPtr == "" {
		panic("You must specify a container id")
	}

	ports := []int{}
	parts := strings.Split(*portsPtr, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])

		num, err := strconv.Atoi(parts[i])
		if err != nil {
			log.Fatalf("Failed to parse port %v: %v", parts[i], err)
		}

		ports = append(ports, num)
	}

	stop := make(chan struct{})
	for _, port := range ports {
		go start_listener(stop, *containerIdPtr, port)
	}

	// Wait for any listeners to give the stop signal
	<-stop
}
