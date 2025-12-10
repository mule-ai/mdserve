# Configuration

MD Serve can be configured using the following environment variables:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| PORT | Port to listen on | 8080 |
| DOC_DIR | Directory to store markdown documents | ./docs |

Example:
```bash
PORT=3000 DOC_DIR=/var/documents ./md-serve
```