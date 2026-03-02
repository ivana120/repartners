package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ivana120/repartners/internal/api"
	"github.com/ivana120/repartners/internal/calculator"
)

//go:embed static
var staticFiles embed.FS

func main() {
	port := flag.Int("port", 8080, "Server port")
	packSizesStr := flag.String("pack-sizes", "250,500,1000,2000,5000", "Comma-separated pack sizes")
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			*port = p
		}
	}

	if envPackSizes := os.Getenv("PACK_SIZES"); envPackSizes != "" {
		*packSizesStr = envPackSizes
	}

	packSizes := parsePackSizes(*packSizesStr)
	if len(packSizes) == 0 {
		log.Fatal("No valid pack sizes provided")
	}

	log.Printf("Starting server with pack sizes: %v", packSizes)

	calc := calculator.New(packSizes)
	handler := api.NewHandler(calc)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/calculate", handler.Calculate)
	mux.HandleFunc("/api/pack-sizes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetPackSizes(w, r)
		case http.MethodPut:
			handler.UpdatePackSizes(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/health", handler.HealthCheck)

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	corsHandler := api.EnableCORS(mux)

	addr := ":" + strconv.Itoa(*port)
	log.Printf("Server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatal(err)
	}
}

func parsePackSizes(s string) []int {
	var sizes []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if size, err := strconv.Atoi(part); err == nil && size > 0 {
			sizes = append(sizes, size)
		}
	}
	return sizes
}
