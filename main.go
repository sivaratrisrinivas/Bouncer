package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// The TokenBucket struct and Take() method are the same as before.
// ... (You can copy them here from your previous code)
type TokenBucket struct {
	capacity       int64
	tokens         int64
	refillRate     time.Duration
	lastRefillTime time.Time
}

func (tb *TokenBucket) Take() bool {
	now := time.Now()
	duration := now.Sub(tb.lastRefillTime)
	tokensToAdd := int64(duration / tb.refillRate)
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefillTime = now
	}
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

var tenantBuckets = map[string]*TokenBucket{
	"free-tier": {
		capacity: 10, tokens: 10, refillRate: time.Second * 15, lastRefillTime: time.Now(),
	},
	"pro-tier": {
		capacity: 50, tokens: 50, refillRate: time.Millisecond * 200, lastRefillTime: time.Now(),
	},
}

// This struct will help us organize the data from our 'products' table.
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// The new queryHandler now needs access to the database connection (*pgx.Conn)
func queryHandler(db_conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.URL.Query().Get("tenant")
		bucket, ok := tenantBuckets[tenantID]
		if !ok {
			http.Error(w, "Invalid tenant", http.StatusForbidden)
			return
		}

		allowed := bucket.Take()
		if !allowed {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// --- NEW DATABASE LOGIC ---
		// Query the database to get all products
		rows, err := db_conn.Query(context.Background(), "select id, name, price from products")
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Loop through the results and build a slice of products
		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		// Convert the products slice to JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func main() {
	// Load .env file for local development.
	// This will be ignored in a production environment like Render.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Please provide your Supabase connection string.")
	}

	// Connect to the database
	db_conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// Close the connection when the application exits
	defer db_conn.Close(context.Background())
	log.Println("Successfully connected to Supabase database!")

	// Pass the database connection to our handler
	http.HandleFunc("/query", queryHandler(db_conn))

	// This tells Go to also serve static files like index.html and script.js
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Hosting platforms provide a PORT environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for local development
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
