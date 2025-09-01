package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/brutella/can"
	"golang.org/x/sys/unix"
)

var lastEvent can.Frame

func eventsHandler(w http.ResponseWriter, r *http.Request, frameCh chan can.Frame) {
	// todo: if can server is not running, return error

	log.Println("Client connected")

	// Set necessary headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins for simplicity

	// Ensure the ResponseWriter supports Flusher for sending data immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Channel to signal when the client disconnects
	closeNotify := r.Context().Done()

	// wait for signals from the can listener

	// todo: listen for event until a configured timeout, or listen for a 'stop' signal from UI / API call.
	for {
		select {
			case <-closeNotify:
				// Client disconnected
				log.Println("Client disconnected")
				return
			case frame := <-frameCh:
				fmt.Printf("Frame: %s\n", frame)
				currentTime := time.Now().Format("15:04:05")
				fmt.Fprintf(w, "event: can-msg\ndata: %s (%s)\n\n", frame, currentTime)
				flusher.Flush()
		}
	}

	// Loop to send events periodically
	// for i := 0; i < 5; i++ {
	// 	select {
	// 	case <-closeNotify:
	// 		// Client disconnected
	// 		log.Println("Client disconnected")
	// 		return
	// 	default:
	// 		// Send a new event with current time
	// 		currentTime := time.Now().Format("15:04:05")
	// 		fmt.Fprintf(w, "event: can-msg\ndata: Last event: %s (current time: %s)\n\n", lastEvent, currentTime)

	// 		// Flush the buffer to send the data to the client immediately
	// 		flusher.Flush()

	// 		// Wait before sending the next event
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }

}

func startCanServer(ctx context.Context, wg *sync.WaitGroup, sig chan os.Signal, frameCh chan can.Frame) {
	defer wg.Done()

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <can-interface>")
	}
	iface := os.Args[1]

	// Create a new CAN bus for the specified interface.
	bus, err := can.NewBusForInterfaceWithName(iface)
	if err != nil {
		log.Fatalf("Failed to create CAN bus for interface %s: %v", iface, err)
	}

	// // Create a channel to receive frames from the bus.
	// frameCh := make(chan can.Frame)

	// Subscribe to all CAN frames and publish them to our channel.
	// `can.NewSubscriber` can also be used for advanced filtering.
	bus.SubscribeFunc(func(f can.Frame) {
		frameCh <- f
	})

	// Use a context to handle graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the CAN bus connection in a separate goroutine.
	go func() {
		if err := bus.ConnectAndPublish(); err != nil {
			log.Printf("CAN bus connection error: %v", err)
			cancel()
		}
	}()

	fmt.Printf("CAN reader started on interface '%s'. Press Ctrl+C to exit.\n", iface)

	<-ctx.Done()
	fmt.Println("Context done, shutting down...")

	// for {
	// 	select {
	// 	case frame := <-frameCh:
	// 		// Handle the incoming CAN frame.
	// 		lastEvent = frame
	// 		fmt.Printf("Received: %s\n", frame)

	// 	case <-sig:
	// 		fmt.Println("\nReceived interrupt signal, shutting down...")
	// 		cancel()
	// 		return

	// 	case <-ctx.Done():
	// 		fmt.Println("Context done, shutting down...")
	// 		return
	// 	}
	// }
}

func startWebServer(ctx context.Context, wg *sync.WaitGroup, sig chan os.Signal, frameCh chan can.Frame) {
	defer wg.Done()

	if frameCh == nil {
		fmt.Println("Warning: No CAN listening.")
	}

	eventsHandlerWrapper := func(w http.ResponseWriter, r *http.Request) {
		eventsHandler(w, r, frameCh)
    }

	mux := http.NewServeMux()
	mux.HandleFunc("/events", eventsHandlerWrapper)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     mux,
		IdleTimeout: 5 * time.Minute, // Set the idle timeout for keep-alive connections
	}

	log.Println("Starting API server on port 8080 ...")
	log.Fatal(server.ListenAndServe())

	// wait for a signal to stop

	// todo: this doesn't work?
	<-ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run the web and can servers in two separate goroutines
	var wg sync.WaitGroup

	// Handle OS signals for graceful shutdown (e.g., Ctrl+C).
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, unix.SIGINT, unix.SIGTERM, syscall.SIGINT, syscall.SIGTERM)

	// add the can server if requested
	if len(os.Args) == 2 { // CAN server requested
		// global frame channel to share between can and web server
		frameCh := make(chan can.Frame)
		wg.Add(2)
		go startCanServer(ctx, &wg, sig, frameCh)
		go startWebServer(ctx, &wg, sig, frameCh)
	} else {
		fmt.Println("CAN server not started.")
		wg.Add(1)
		go startWebServer(ctx, &wg, sig, nil)
	}



	// Wait for an interrupt signal to initiate graceful shutdown
	<-sig
	// Handle shutdown signal (Ctrl+C or SIGTERM)
	fmt.Println("Received shutdown signal. Shutting down gracefully...")
	ctx.Done()
	//sig <- unix.SIGTERM
	//<-ctx.Done()
	wg.Done()
	cancel()

	//wg.Wait()

	// blocks indefinitely, or until an error occurs

}
