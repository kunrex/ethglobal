# In-Memory Git Server

A lightweight Git server written in Go that stores repositories in memory using the `go-git` library. This server provides a simple HTTP interface for Git operations and repository management.

## Features

- **In-Memory Storage**: All repositories are stored in memory using `go-git`'s memory storage backend
- **HTTP Git Protocol**: Supports standard Git HTTP operations (clone, fetch, push)
- **Repository Management**: Create, list, and manage repositories via REST API
- **Web Interface**: Simple HTML interface for testing and management
- **Concurrent Safe**: Thread-safe repository access using mutex locks

## Dependencies

- [go-git/v5](https://github.com/go-git/go-git) - Pure Go implementation of Git
- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router and URL matcher

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone or download this repository
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Starting the Server

```bash
go run main.go
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

### Web Interface

Visit `http://localhost:8080` in your browser to see the web interface with available endpoints and usage examples.

### API Endpoints

#### Repository Management

- **GET /repos** - List all repositories
- **POST /repos?name=repo-name** - Create a new repository
- **POST /{repo}/files** - Add a file to a repository (form data: filename, content)

#### Git Protocol Endpoints

- **GET /{repo}/info/refs?service=git-upload-pack** - Get refs for clone/fetch operations
- **POST /{repo}/git-upload-pack** - Handle clone/fetch operations
- **POST /{repo}/git-receive-pack** - Handle push operations

### Example Usage

#### List Repositories
```bash
curl http://localhost:8080/repos
```

#### Create a New Repository
```bash
curl -X POST "http://localhost:8080/repos?name=my-new-repo"
```

#### Add a File to Repository
```bash
curl -X POST "http://localhost:8080/test-repo/files" \
  -d "filename=README.md" \
  -d "content=# My Repository\nThis is a test repository."
```

#### Clone a Repository (Git)
```bash
git clone http://localhost:8080/test-repo
```

**Note:** Use the repository name without `.git` suffix for cloning.

## Architecture

The server consists of:

1. **InMemoryGitServer**: Main server struct that manages repositories in memory
2. **Repository Storage**: Uses `go-git`'s memory storage backend
3. **HTTP Handlers**: Handle both Git protocol and management operations
4. **Concurrent Access**: Thread-safe operations using mutex locks

## Limitations

This is a simplified implementation for demonstration purposes. A production Git server would need:

- Full Git pack protocol implementation
- Authentication and authorization
- Persistent storage options
- Better error handling and logging
- Support for all Git operations (branches, tags, etc.)

## Development

To build the server:

```bash
go build -o git-server main.go
```

To run the demo:

```bash
./demo.sh
```

To run tests:

```bash
./test.sh
```

## License

This project is licensed under the same license as the main repository.
