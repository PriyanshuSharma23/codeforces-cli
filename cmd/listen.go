// cmd/listener.go
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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
			var templateStr string
			if templatePath != "" {
				t, err := dm.LoadTemplate(templatePath)
				cobra.CheckErr(err)
				templateStr = t
			}

			progFile := fmt.Sprintf("%s.%s", viper.GetString("programFile"), viper.GetString("language"))
			if err := dm.WriteProgramFile(problemKey, progFile, templateStr); err != nil {
				http.Error(w, "Error writing program file", http.StatusInternalServerError)
				return
			}

			if err := dm.WriteMetadata(problemKey, ccproblem); err != nil {
				logger.Printf("Warning: could not write metadata: %v", err)
			}

			go func() {
				// 🔥 Open the editor using editorCommand
				editorCmdTemplate := viper.GetString("editorCommand")
				editorTemplate, err := template.New("editor").Parse(editorCmdTemplate)
				if err != nil {
					logger.Printf("Invalid editorCommand template: %v", err)
				} else {
					var cmdBuf bytes.Buffer
					err = editorTemplate.Execute(&cmdBuf, map[string]string{
						"Path": filepath.Join(dm.FullProblemPath(problemKey), progFile),
						"Dir":  filepath.Join(dm.FullProblemPath(problemKey)),
					})
					if err != nil {
						logger.Printf("Failed to render editor command: %v", err)
					} else {
						editorArgs := strings.Fields(cmdBuf.String())
						cmd := exec.Command(editorArgs[0], editorArgs[1:]...)
						cmd.Dir = dm.FullProblemPath(problemKey)
						err = cmd.Start()
						if err != nil {
							logger.Printf("Failed to launch editor: %v", err)
						}
					}
				}

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
