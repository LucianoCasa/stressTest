package cli

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executa o teste de carga",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Erro: parâmetro --url é obrigatório")
			return
		}

		if requests <= 0 || concurrency <= 0 {
			fmt.Println("Erro: --requests e --concurrency devem ser maiores que zero")
			return
		}

		start := time.Now()
		results := make(chan Result, requests)
		var wg sync.WaitGroup

		reqsPerWorker := requests / concurrency
		extra := requests % concurrency

		for i := 0; i < concurrency; i++ {
			n := reqsPerWorker
			if i < extra {
				n++
			}
			wg.Add(1)
			go worker(url, n, results, &wg)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		statusCount := make(map[int]int)
		var success200 int
		var total int
		var totalDuration time.Duration

		for r := range results {
			total++
			if r.StatusCode == 200 {
				success200++
			}
			statusCount[r.StatusCode]++
			totalDuration += r.Duration
		}

		elapsed := time.Since(start)
		success200Perc := float64(success200) / float64(total) * 100

		// relatório
		fmt.Println("===== Relatório de Teste de Carga =====")
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("Requests realizados: %d\n", total)
		fmt.Printf("Total de Workers: %d\n\n", concurrency)

		fmt.Printf("Tempo total gasto: %v\n", elapsed)
		fmt.Printf("Tempo médio por request: %v\n\n", totalDuration/time.Duration(total))

		fmt.Printf("Sucessos: %d (%.2f%%) | Erros: %d (%.2f%%)\n", success200, success200Perc, total-success200, float64(total-success200)/float64(total)*100)
		fmt.Println("Distribuição de status codes:")
		for code, count := range statusCount {
			if code == 0 {
				fmt.Printf("  Falhas de conexão: %d\n", count)
			} else {
				fmt.Printf("  %d: %d (%.2f%%)\n", code, count, float64(count)/float64(total)*100)
			}
		}
	},
}

func worker(url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < requests; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{StatusCode: 0, Duration: duration, Error: err}
			continue
		}

		// fmt.Printf("Request to %s time: %v, status: %d\n", url, duration, resp.StatusCode)
		results <- Result{StatusCode: resp.StatusCode, Duration: duration, Error: nil}
		_ = resp.Body.Close()
	}
}

func init() {
	runCmd.Flags().StringVarP(&url, "url", "u", "http://localhost:8080", "URL do serviço a ser testado")
	runCmd.Flags().IntVarP(&requests, "requests", "r", 1000, "Número total de requests")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Número de chamadas simultâneas")
}
