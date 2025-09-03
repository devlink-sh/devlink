// controller/main.go
package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	Name  string `json:"name"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

type Hive struct {
	Name     string              `json:"name"`
	Services map[string]Service  `json:"services"`
}

var (
	hives = make(map[string]*Hive)
	mu    sync.Mutex
)

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// POST /hives/create?name=<name>
func createHive(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing hive name", http.StatusBadRequest)
		return
	}
	token := name + "-" + randString(6)

	mu.Lock()
	defer mu.Unlock()
	hives[token] = &Hive{
		Name:     name,
		Services: make(map[string]Service),
	}
	log.Printf("Hive created: %s", token)
	w.Write([]byte(token))
}

// POST /hives/contribute?hive=<token>&service=<name>&port=<port>&token=<shareToken>
func contribute(w http.ResponseWriter, r *http.Request) {
	hiveToken := r.URL.Query().Get("hive")
	service := r.URL.Query().Get("service")
	port := r.URL.Query().Get("port")
	shareToken := r.URL.Query().Get("token")

	if hiveToken == "" || service == "" || port == "" || shareToken == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	h, ok := hives[hiveToken]
	if !ok {
		http.Error(w, "hive not found", http.StatusNotFound)
		return
	}

	h.Services[service] = Service{
		Name:  service,
		Port:  port,
		Token: shareToken,
	}

	log.Printf("Service %s added to hive %s", service, hiveToken)
	w.Write([]byte("ok"))
}

// GET /hives/services?hive=<token>
func getServices(w http.ResponseWriter, r *http.Request) {
	hiveToken := r.URL.Query().Get("hive")
	mu.Lock()
	defer mu.Unlock()

	h, ok := hives[hiveToken]
	if !ok {
		http.Error(w, "hive not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(h.Services)
}

func main() {
	http.HandleFunc("/hives/create", createHive)
	http.HandleFunc("/hives/contribute", contribute)
	http.HandleFunc("/hives/services", getServices)

	log.Println("Hive Controller running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
