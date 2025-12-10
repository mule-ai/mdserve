# MD Serve - Project Overview

MD Serve is a lightweight markdown server implemented in Go that serves markdown documents via a REST API and renders them as HTML in the browser. It uses HTMX for dynamic web interface updates and stores documents in the filesystem.

## Project Structure

```
md-serve/
├── cmd/
│   └── md-serve/
│       └── main.go          # Application entry point
├── internal/
│   └── server/
│       └── server.go        # HTTP server implementation
├── docs/                    # Default document storage directory
│   ├── README.md            # Basic project documentation
│   └── configuration.md     # Configuration documentation
├── Makefile                 # Build instructions
├── Dockerfile               # Containerization setup
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksums
├── README.md                # Project documentation
└── .gitignore               # Git ignore rules
```

## Core Components

### Server Implementation (internal/server/server.go)

The server package contains all the core functionality:
- HTTP request routing and handling
- Document management (list, get, create)
- Filesystem storage operations
- Markdown rendering to HTML
- Web interface with HTMX integration

### Main Entry Point (cmd/md-serve/main.go)

Simple initialization of the server and startup on the configured port.

## Usage

### Running the Service

#### Using the Binary
```bash
# Start with defaults (port 8080, docs in ./docs)
./md-serve

# Custom port
PORT=3000 ./md-serve
# Or using flags
./md-serve -port 3000

# Custom document directory
DOC_DIR=/path/to/documents ./md-serve
# Or using flags
./md-serve -doc-dir /path/to/documents

# Custom public URL
PUBLIC_URL=https://mydocs.example.com ./md-serve
# Or using flags
./md-serve -public-url https://mydocs.example.com

# All customizations with environment variables
PORT=3000 DOC_DIR=/path/to/documents PUBLIC_URL=https://mydocs.example.com ./md-serve

# All customizations with flags
./md-serve -port 3000 -doc-dir /path/to/documents -public-url https://mydocs.example.com
```

#### Using Go Run
```bash
# Development mode
go run cmd/md-serve/main.go

# With environment variables
PORT=3000 DOC_DIR=/path/to/documents go run cmd/md-serve/main.go
```

#### Using Docker
```bash
# Build the image
docker build -t md-serve .

# Run the container
docker run -p 8080:8080 md-serve

# With volume mount for documents
docker run -p 8080:8080 -v /host/docs:/root/docs md-serve

# With environment variables
docker run -p 3000:3000 -e PORT=3000 md-serve
```

### API Endpoints

1. **List Documents**
   - Method: GET
   - Path: `/api/documents`
   - Response: JSON array of document names
   - Example: `curl http://localhost:8080/api/documents`

2. **Get Document Content**
   - Method: GET
   - Path: `/api/document`
   - Query Parameter: `filename` (without .md extension)
   - Response: Raw markdown content
   - Example: `curl "http://localhost:8080/api/document?filename=README"`

3. **Create Document**
   - Method: POST
   - Path: `/api/documents`
   - Form Data: `filename` (without .md extension), `content` (markdown text)
   - Response: Success message
   - Example: `curl -X POST http://localhost:8080/api/documents -d "filename=test&content=# Test%0AContent"`

4. **View Rendered Document**
   - Method: GET
   - Path: `/{document-name}` (without .md extension)
   - Response: HTML-rendered document
   - Example: `curl http://localhost:8080/README`

### Web Interface

Access `http://localhost:8080` in a browser to use the web interface:
- Create new documents using the form
- View list of existing documents
- Click on document names to view rendered versions

## Making Updates

### Development Workflow

1. **Code Changes**
   - Modify files in `internal/server/` for core functionality
   - Update `cmd/md-serve/main.go` for application startup changes
   - Add new dependencies with `go get`

2. **Testing Changes**
   ```bash
   # Run in development mode
   go run cmd/md-serve/main.go
   
   # Test API endpoints
   curl http://localhost:8080/api/documents
   ```

3. **Building**
   ```bash
   # Build for current platform
   make build
   
   # Build statically linked binary
   make build-static
   
   # Build for all platforms
   make build-all
   ```

### Adding New Features

1. **New API Endpoints**
   - Add route handling in `server.ServeHTTP`
   - Implement handler methods in the Server struct
   - Follow existing patterns for error handling and response formatting

2. **Enhanced Document Operations**
   - Extend the Server struct methods
   - Maintain consistency with existing filesystem operations
   - Add appropriate error handling

3. **UI Improvements**
   - Modify HTML templates in handler methods
   - Update styling by modifying Bootstrap classes or adding custom CSS
   - Enhance HTMX interactions as needed

### Dependencies

Current dependencies:
- `github.com/gomarkdown/markdown` for markdown parsing and rendering

To update dependencies:
```bash
go get -u
go mod tidy
```

## Integration Guide

### Embedding in Other Applications

MD Serve can be embedded in other Go applications:

```go
import "github.com/mule-ai/md-serve/internal/server"

// Create and configure server
srv := server.New()

// Use with custom HTTP server
http.ListenAndServe(":8080", srv)
```

### API Integration

Other services can integrate with MD Serve through its REST API:

1. **Service-to-Service Communication**
   - List documents: `GET /api/documents`
   - Retrieve content: `GET /api/document?filename={name}`
   - Create content: `POST /api/documents`

2. **Automation Scripts**
   ```bash
   # Backup all documents
   curl -s http://localhost:8080/api/documents | jq -r '.[]' | while read doc; do
     curl -s "http://localhost:8080/api/document?filename=$doc" > "$doc.md"
   done
   ```

3. **CI/CD Integration**
   - Deploy documentation updates by calling the create endpoint
   - Validate document availability through the list endpoint

### Extending Functionality

1. **Authentication/Authorization**
   - Add middleware to the Server struct
   - Implement token validation or session management
   - Protect specific endpoints as needed

2. **Storage Backends**
   - Replace filesystem operations with database calls
   - Implement interfaces for storage abstraction
   - Maintain the same API contract

3. **Additional Formats**
   - Extend markdown rendering with additional processors
   - Add support for other document formats
   - Implement format detection and appropriate rendering

## Configuration

Environment Variables:
- `PORT`: Server port (default: 8080)
- `DOC_DIR`: Document storage directory (default: ./docs)
- `PUBLIC_URL`: Public URL for the service (default: http://localhost:8080)

## Deployment Options

1. **Binary Deployment**
   - Build static binary with `make build-static`
   - Deploy binary to target system
   - Configure as service or run directly

2. **Container Deployment**
   - Build Docker image
   - Deploy to container orchestration platform
   - Mount volumes for persistent document storage

3. **Cloud Deployment**
   - Package as container or binary
   - Deploy to cloud provider of choice
   - Configure load balancing and scaling as needed