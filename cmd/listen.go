/*
Copyright © 2025 Priyanshu Sharma inbox.priyanshu@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/ccparser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listens for problems from Competitive Companion",
	Run: func(cmd *cobra.Command, args []string) {
		mux := http.NewServeMux()

		server := &http.Server{
			Addr:    ":10045", // fixed to match printed message
			Handler: mux,
		}

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Only POST supported", http.StatusMethodNotAllowed)
				return
			}

			var ccproblem ccparser.CCProblem
			if err := json.NewDecoder(r.Body).Decode(&ccproblem); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			// Debug: print loaded config values
			fmt.Println("------------ Loaded Configuration ------------")
			fmt.Println("root:", viper.GetString("root"))
			fmt.Println("language:", viper.GetString("language"))
			fmt.Println("buildCommand:", viper.GetString("buildCommand"))
			fmt.Println("executeCommand:", viper.GetString("executeCommand"))
			fmt.Println("testCaseInputPrefix:", viper.GetString("testCaseInputPrefix"))
			fmt.Println("testCaseOutputPrefix:", viper.GetString("testCaseOutputPrefix"))
			fmt.Println("port:", viper.GetInt("port"))
			fmt.Println("editorCommand:", viper.GetString("editorCommand"))
			fmt.Println("templatePath:", viper.GetString("templatePath"))
			fmt.Println("------------------------------------------------")

			// Display metadata
			fmt.Printf("\n✅ Problem: %s\n🔗 %s\n🧠 %dMB | ⏱️ %dms\n", ccproblem.Name, ccproblem.URL, ccproblem.MemoryLimit, ccproblem.TimeLimit)
			fmt.Printf("Group: %s | Interactive: %v | Test Type: %s\n", ccproblem.Group, ccproblem.Interactive, ccproblem.TestType)

			for i, test := range ccproblem.Tests {
				fmt.Printf("\n--- Test #%d ---\n📝 Input:\n%s\n✅ Expected Output:\n%s\n", i+1, test.Input, test.Output)
			}

			if java, ok := ccproblem.Languages["java"]; ok {
				fmt.Printf("\n☕ Java TaskClass: %s | MainClass: %s\n", java.TaskClass, java.MainClass)
			}

			logger := log.New(os.Stdout, "", log.LstdFlags)
			parser := ccparser.NewParser(logger)

			problem, err := parser.Parse(&ccproblem)
			if err != nil {
				logger.Printf("Failed to parse problem: %s", err)
				http.Error(w, "Failed to parse problem", http.StatusInternalServerError)
				return
			}

			fmt.Printf("\n📦 Parsed Problem Struct:\n%+v\n", problem)

			go func() {
				time.Sleep(2 * time.Second)
				if err := server.Close(); err != nil {
					logger.Printf("Failed to close server: %v", err)
				} else {
					logger.Println("Server shut down gracefully.")
				}
			}()

			w.Write([]byte("✅ Problem received! You may close Competitive Companion."))
		})

		fmt.Println("🟢 Listening on http://localhost:10045 ...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("❌ ListenAndServe(): %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}
