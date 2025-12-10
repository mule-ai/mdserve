# MD Serve

A simple markdown server built with Go and HTMX.

## Features

- REST API for managing markdown documents
- Web interface with HTMX for dynamic updates
- Automatic markdown rendering
- Filesystem-based storage
- Statically compiled binary

## API Endpoints

- `GET /api/documents` - List all documents
- `GET /api/document?filename={name}` - Get raw markdown content
- `POST /api/documents` - Create a new document
- `GET /{document}` - View rendered document in browser

## Building

To build the statically compiled binary:

```bash
make build-static
```

To build for all platforms:

```bash
make build-all
```

## Running

To run the server:

```bash
./md-serve
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable:

```bash
PORT=3000 ./md-serve
```

## Usage

### Web Interface

Visit `http://localhost:8080` in your browser to access the web interface. From there you can:
- Create new documents
- View existing documents
- See a list of all documents

### API Usage

List all documents:
```bash
curl http://localhost:8080/api/documents
```

Get a document's raw markdown:
```bash
curl "http://localhost:8080/api/document?filename=README"
```

Create a new document:
```bash
curl -X POST http://localhost:8080/api/documents -d "filename=mydoc&content=# My Document%0A%0AThis is my document."
```

Response (JSON):
```json
{
  "message": "Document 'mydoc' created successfully",
  "filename": "mydoc",
  "url": "http://localhost:8080/mydoc",
  "api_url": "http://localhost:8080/api/document?filename=mydoc"
}
```

View a rendered document:
```bash
curl http://localhost:8080/mydoc
```