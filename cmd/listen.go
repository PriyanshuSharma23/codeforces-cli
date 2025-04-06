// cmd/listener.go
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/ccparser"
	"github.com/PriyanshuSharma23/codeforces-cli/internal/directorymanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for Competitive Companion problems",
	Run: func(cmd *cobra.Command, args []string) {
		mux := http.NewServeMux()
		port := viper.GetString("port")

		server := &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
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

			logger := log.New(os.Stdout, "", log.LstdFlags)
			parser := ccparser.NewParser(logger)

			parsedProblem, err := parser.Parse(&ccproblem)
			if err != nil {
				http.Error(w, "Failed to parse problem", http.StatusInternalServerError)
				return
			}

			dm := directorymanager.NewDirectoryManager(viper.GetString("root"), logger)
			problemKey := directorymanager.Problem{
				ContestCode: parsedProblem.ContestCode,
				ProblemCode: parsedProblem.ProblemCode,
			}

			if _, err := dm.EnsureDir(problemKey); err != nil {
				http.Error(w, "Could not prepare problem directory", http.StatusInternalServerError)
				return
			}

			if err := dm.WriteTestCases(
				problemKey,
				parsedProblem.TestCases,
				viper.GetString("testCaseInputPrefix"),
				viper.GetString("testCaseOutputPrefix"),
			); err != nil {
				http.Error(w, "Error writing test cases", http.StatusInternalServerError)
				return
			}

			templatePath := viper.GetString("templatePath")
			var template string
			if templatePath != "" {
				t, err := dm.LoadTemplate(templatePath)
				cobra.CheckErr(err)
				template = t
			}

			progFile := fmt.Sprintf("%s.%s", viper.GetString("programFile"), viper.GetString("language"))
			if err := dm.WriteProgramFile(problemKey, progFile, template); err != nil {
				http.Error(w, "Error writing program file", http.StatusInternalServerError)
				return
			}

			if err := dm.WriteMetadata(problemKey, ccproblem); err != nil {
				logger.Printf("Warning: could not write metadata: %v", err)
			}

			go func() {
				time.Sleep(2 * time.Second)
				server.Close()
			}()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":      "success",
				"problemPath": problemKey.RelativeDir(),
				"programFile": progFile,
			})
		})

		fmt.Printf("🟢 Listening on http://localhost:%s...\n", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("❌ Server error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}

