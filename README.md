# codeforces-cli

`codeforces-cli` is a command-line tool designed to enhance competitive programming experiences, particularly for Codeforces users. This tool provides a seamless integration with the Competitive Companion browser extension, helps in managing problems, setting up a working environment with boilerplate code, and allows for easy compilation and execution of solutions with direct feedback from test cases.

## Features

- **Automatic Problem Directory Creation**: Automatically creates a directory structure for problems, organizing them by contest and problem code.
- **Competitive Companion Integration**: Seamlessly imports problem data using the Competitive Companion extension.
- **Customizable Templates**: Utilizes user-provided templates to quickly set up code files for problems.
- **Execution and Testing**: Compiles and runs solutions against sample test cases from the problem statement, providing a summary of passed and failed cases.
- **Flexible Configuration**: Allows customization for various programming languages, code editors, and build & execute commands.

## Installation

### Prerequisites

- Go (Golang) v1.16 or above.
- Competitive Companion browser extension.
- Optional: Your preferred code editor (e.g., Visual Studio Code, Vim, etc.) and relevant compilers/interpreters for your programming language of choice.

### Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/PriyanshuSharma23/codeforces-cli.git
   cd codeforces-cli
   ```

2. **Build the Application**

   ```bash
   go build -o codeforces-cli
   ```

3. **Move the Executable**

   ```bash
   mv codeforces-cli /usr/local/bin/
   ```

4. **Verify Installation**

   ```bash
   codeforces-cli --help
   ```

## Configuration

### Configuration File

The application reads its configuration from a YAML file. The default configuration file is located at `~/.config/codeforces-cli/config.yaml`. You can provide a custom configuration file using the `--config` flag.

### Sample Configuration

```yaml
root: /home/user/codeforces/problems
language: py
programFile: main
buildCommand: ""
executeCommand: "python3 {{.Path}}"
testCaseInputPrefix: "input"
testCaseOutputPrefix: "output"
port: 10045
editorCommand: "nvim {{.Path}}"
templatePath: /home/user/codeforces/templates/main.cpp
```

- **root**: Directory where problems are stored.
- **language**: Programming language file extension.
- **programFile**: Name of the main program file.
- **buildCommand**: Command to compile the program, if necessary.
- **executeCommand**: Command to run the built program or script.
- **testCaseInputPrefix**: Prefix for input test case files.
- **testCaseOutputPrefix**: Prefix for output test case files.
- **port**: Port for Competitive Companion to send data to.
- **editorCommand**: Command template to open the code editor.
- **templatePath**: Path to the code template file that gets copied when a problem is created.

## Usage

### Listening for Problems

Start the tool's server to listen for problem data sent from Competitive Companion:

```bash
codeforces-cli listen
```

### Running Test Cases

Navigate to the problem directory and execute test cases using:

```bash
codeforces-cli execute
```

## Development

### Running Tests

To run tests, use the following command:

```bash
go test ./...
```

### Logging and Debugging

The application logs significant events and errors to the console. Ensure your environment supports console output to monitor logging.

## Contribution

Contributions are welcome. Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
