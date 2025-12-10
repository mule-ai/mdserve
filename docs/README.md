# MD Serve

A simple markdown server built with Go and HTMX.

## Features

- REST API for managing markdown documents
- Web interface with HTMX for dynamic updates
- Automatic markdown rendering
- Filesystem-based storage

## API Endpoints

- `GET /api/documents` - List all documents
- `GET /api/document?filename={name}` - Get raw markdown content
- `POST /api/documents` - Create a new document
- `GET /{document}` - View rendered document in browser