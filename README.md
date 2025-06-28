# Xuan
This repository contains a simple web shell written in **Go (Golang)**. This tool is designed to demonstrate how a basic web shell works, enabling remote command execution and file uploads via a web interface with an enhanced user experience.

Disclaimer: This tool is for educational and ethical hacking purposes only. Using this tool on any system without explicit, prior authorization is illegal and unethical. The developer is not responsible for any misuse or damage caused by this software.

## Features
Interactive Web Interface: A modern, terminal-like web interface for command execution and file uploads.

Command Execution: Execute arbitrary shell commands on the remote server by typing them directly into a dedicated input field.

File Upload: Upload files from your local machine to the remote server via a file selection field, with an option to specify the destination path.

Real-time Output: Displays command execution results directly on the web page.

Visual Feedback: Clear success and error messages are displayed for each operation.

Single Binary: Compiled into a single executable file, making it easy to deploy on various target systems without external dependencies.

Cross-Platform: Can be compiled for different operating systems (Linux, Windows, macOS) and architectures.

Lightweight HTTP Server: Includes a built-in HTTP server to serve the web interface and handle requests.

## How it Works
The Xuan operates by:

Serving an Enhanced Web Interface: Upon execution, it starts an HTTP server (defaulting to port 8080) that presents a user-friendly HTML form with a dark, terminal-inspired theme for a more immersive experience.

Handling Commands via POST: When a command is submitted through the "Execute Command" form (using an HTTP POST request), the shell captures the input from the designated text field. It then executes the command on the host system using sh -c (for Unix-like systems like Linux/macOS) and returns the combined standard output and error back to the web interface, displayed in a dedicated, monospace output area.

Handling File Uploads via POST: When a file is uploaded using the "File Upload" form (also via an HTTP POST request with multipart/form-data), the shell receives the file data. It then saves the file to the specified destination path on the server. If no specific destination is provided, it defaults to the original filename in the current working directory where the shell is running.

Dynamic Page Rendering: After each command execution or file upload, the entire web page is re-rendered to display the results, status messages (success or error), and keep the forms ready for the next action.
