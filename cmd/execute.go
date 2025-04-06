/*
Copyright © 2025 Priyanshu Sharma inbox.priyanshu@gmail.com
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/PriyanshuSharma23/codeforces-cli/internal/execution"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run test cases for a problem",
	Long: `Executes all test cases defined for a problem within its directory.

To use this command, navigate your terminal to the problem's directory, typically located at /<contest>/<problem_code>. 

The command utilizes the 'buildCommand' specified in the configuration to compile the program and the 'executeCommand' to run the compiled executable.

It processes all test files within the directory, comparing the program's output to the expected results. A summary report is then displayed, indicating the number of passed and failed test cases. For failed cases, both the expected output and the program's output are shown for debugging.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootPath := viper.GetString("root")
		buildCommand := viper.GetString("buildCommand")
		executionCommand := viper.GetString("executeCommand")
		inputPrefix := viper.GetString("testCaseInputPrefix")
		outputPrefix := viper.GetString("testCaseOutputPrefix")
		programFile := viper.GetString("programFile")
		language := viper.GetString("language")

		testCasesDir, err := os.Getwd()
		cobra.CheckErr(err)

		pathStr := filepath.Join(testCasesDir, fmt.Sprintf("%s.%s", programFile, language))

		variables := map[string]string{
			"Path": pathStr,
			"Dir":  testCasesDir,
		}

		buildTemplate, err := template.New("build").Parse(buildCommand)
		if err != nil {
			logger.Printf("Invalid buildCommand template: %v", err)
		}

		var buildCmdBuf bytes.Buffer
		err = buildTemplate.Execute(&buildCmdBuf, variables)
		cobra.CheckErr(err)

		execTemplate, err := template.New("exec").Parse(executionCommand)
		if err != nil {
			logger.Printf("Invalid execCommand template: %v", err)
		}

		var execCmdBuf bytes.Buffer
		err = execTemplate.Execute(&execCmdBuf, variables)
		cobra.CheckErr(err)

		fmt.Println(execCmdBuf.String())

		em := execution.NewEngine(
			rootPath,
			testCasesDir,
			buildCmdBuf.String(),
			execCmdBuf.String(),
			inputPrefix,
			outputPrefix,
			logger,
		)

		res, err := em.Execute()
		cobra.CheckErr(err)

		printResults(res)
	},
}

func printResults(results []execution.Result) {
	passedCount := 0
	failedCount := 0

	for _, result := range results {
		if result.Ok {
			color.Green("Test Case %d: Passed", result.TestCase)
			passedCount++
		} else {
			color.Red("Test Case %d: Failed", result.TestCase)
			fmt.Println(color.YellowString("Expected Output:"))
			fmt.Println(result.ExpectedOutput)
			fmt.Println(color.YellowString("Program Output:"))
			fmt.Println(result.ProgramOutput)
			failedCount++
		}
		fmt.Println("---") // Separator
	}

	fmt.Println(color.CyanString("Summary:"))
	fmt.Printf("Passed: %d, Failed: %d, Total: %d\n", passedCount, failedCount, len(results))
}

func init() {
	rootCmd.AddCommand(executeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
