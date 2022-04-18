package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	_ = godotenv.Load()
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	eg := new(errgroup.Group)

	eg.Go(func() error {
		port, ok := os.LookupEnv("PORT")
		if !ok {
			port = "8080"
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s! Connecting via 0.0.0.0:%s.", r.RemoteAddr, port)
		})
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), mux)
	})
	log.Println("started standard server")

	eg.Go(func() error {
		port := "11106"

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s! Connecting via 127.0.0.1:%s.", r.RemoteAddr, port)
		})

		// todo: explicitly listen on tailscale address
		// for some reason, interfaces.Tailscale() doesn't find interface right now
		// (when running on Railway).
		return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), mux)
	})
	log.Println("started local only server")

	return eg.Wait()
}
