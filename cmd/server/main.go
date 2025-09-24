package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Regras de probabilidade: código -> peso
var statusRules = map[int]int{
	200: 90, // 85%
	500: 5,  // 5%
	404: 2,  // 5%
	429: 3,  // 5%
}

var statusList []int
var rnd *rand.Rand

func init() {
	for code, weight := range statusRules {
		for i := 0; i < weight; i++ {
			statusList = append(statusList, code)
		}
	}
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func handler(w http.ResponseWriter, r *http.Request) {
	code := statusList[rnd.Intn(len(statusList))]
	w.WriteHeader(code)
	_, _ = w.Write([]byte(fmt.Sprintf("%d - %s", code, http.StatusText(code))))

	// fmt.Printf("[%s] %s %s -> %d %s\n",
	// 	time.Now().Format("15:04:05"), // horário da requisição
	// 	r.Method,                      // método HTTP (GET, POST, etc.)
	// 	r.URL.Path,                    // rota chamada
	// 	code,
	// 	http.StatusText(code),
	// )

}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
