package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stresstest",
	Short: "Ferramenta CLI para testes de carga em serviços web",
	Long:  "Stresstest é uma ferramenta simples em Go para realizar testes de carga em serviços web, medindo status codes e tempo de resposta.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
