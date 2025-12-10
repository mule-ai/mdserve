package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Server handles HTTP requests for the markdown server
type Server struct {
	docDir    string
	publicURL string
}

// New creates a new Server instance
func New() *Server {
	docDir := os.Getenv("DOC_DIR")
	if docDir == "" {
		docDir = "./docs"
	}

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8080"
	}

	// Ensure docs directory exists
	if err := os.MkdirAll(docDir, 0755); err != nil {
		log.Fatalf("Failed to create docs directory: %v", err)
	}

	return &Server{
		docDir:    docDir,
		publicURL: publicURL,
	}
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle API routes
	if strings.HasPrefix(r.URL.Path, "/api/") {
		s.handleAPI(w, r)
		return
	}
	
	// Handle document viewing
	s.handleDocumentView(w, r)
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/documents":
		switch r.Method {
		case http.MethodGet:
			s.listDocuments(w, r)
		case http.MethodPost:
			s.createDocument(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/api/document":
		if r.Method == http.MethodGet {
			s.getDocument(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleDocumentView(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		// Show index page
		s.showIndex(w, r)
		return
	}
	
	docPath := filepath.Join(s.docDir, path+".md")
	content, err := ioutil.ReadFile(docPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to read document", http.StatusInternalServerError)
		return
	}
	
	// Render markdown to HTML
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)
	
	doc := parser.Parse(content)
	htmlContent := string(markdown.Render(doc, renderer))
	
	// Serve HTML page with rendered content
	htmlPage := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        :root {
            --background-color: #ffffff;
            --text-color: #212529;
            --link-color: #0d6efd;
            --border-color: #dee2e6;
        }

        .dark-mode {
            --background-color: #1a1a1a;
            --text-color: #e0e0e0;
            --link-color: #6ea8fe;
            --border-color: #495057;
        }

        body {
            background-color: var(--background-color);
            color: var(--text-color);
        }

        a {
            color: var(--link-color);
        }

        .theme-toggle {
            position: fixed;
            top: 10px;
            right: 10px;
            z-index: 1000;
        }
    </style>
</head>
<body>
    <div class="theme-toggle">
        <button id="themeToggle" class="btn btn-sm btn-outline-secondary" title="Toggle theme">
            <span id="themeIcon">🌙</span>
        </button>
    </div>

    <div class="container mt-4">
        <div class="mt-4">%s</div>
    </div>

    <script>
        // Function to set theme
        function setTheme(theme) {
            if (theme === 'dark') {
                document.body.classList.add('dark-mode');
                localStorage.setItem('theme', 'dark');
                document.getElementById('themeIcon').textContent = '☀️';
            } else {
                document.body.classList.remove('dark-mode');
                localStorage.setItem('theme', 'light');
                document.getElementById('themeIcon').textContent = '🌙';
            }
        }

        // Function to detect system theme
        function detectSystemTheme() {
            if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                return 'dark';
            } else {
                return 'light';
            }
        }

        // Initialize theme
        function initTheme() {
            const savedTheme = localStorage.getItem('theme');
            if (savedTheme) {
                setTheme(savedTheme);
            } else {
                // Follow system theme by default
                setTheme(detectSystemTheme());
            }
        }

        // Toggle theme
        function toggleTheme() {
            if (document.body.classList.contains('dark-mode')) {
                setTheme('light');
            } else {
                setTheme('dark');
            }
        }

        // Set up event listeners
        document.addEventListener('DOMContentLoaded', function() {
            initTheme();
            document.getElementById('themeToggle').addEventListener('click', toggleTheme);

            // Listen for system theme changes
            if (window.matchMedia) {
                window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
                    // Only change theme automatically if user hasn't manually set a theme
                    if (!localStorage.getItem('theme')) {
                        setTheme(e.matches ? 'dark' : 'light');
                    }
                });
            }
        });
    </script>
</body>
</html>`, path, htmlContent)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlPage)
}

func (s *Server) showIndex(w http.ResponseWriter, r *http.Request) {
	documents, err := s.getAllDocuments()
	if err != nil {
		http.Error(w, "Failed to list documents", http.StatusInternalServerError)
		return
	}

	// Create HTML list of documents
	docList := ""
	if len(documents) == 0 {
		docList = "<li class='list-group-item'>No documents found</li>"
	} else {
		for _, doc := range documents {
			docList += fmt.Sprintf("<li class='list-group-item'><a href='/%s'>%s</a></li>", doc, doc)
		}
	}

	htmlPage := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>MD Serve</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <h1>MD Serve</h1>

        <div class="row mt-4">
            <div class="col-md-6">
                <h2>Create Document</h2>
                <form id="createForm">
                    <div class="mb-3">
                        <label for="filename" class="form-label">Filename (without .md extension)</label>
                        <input type="text" class="form-control" id="filename" name="filename" required>
                    </div>
                    <div class="mb-3">
                        <label for="content" class="form-label">Content</label>
                        <textarea class="form-control" id="content" name="content" rows="10" required></textarea>
                    </div>
                    <button type="submit" class="btn btn-primary"
                            hx-post="/api/documents"
                            hx-target="#result"
                            hx-vals="js:{filename: document.getElementById('filename').value, content: document.getElementById('content').value}"
                            hx-headers='{"Content-Type": "application/json"}'>
                        Create Document
                    </button>
                </form>
                <div id="result" class="mt-3"></div>
            </div>

            <div class="col-md-6">
                <h2>Documents</h2>
                <ul class="list-group">
                    %s
                </ul>
            </div>
        </div>
    </div>

    <script>
        // Handle form submission with JSON payload
        document.getElementById('createForm').addEventListener('submit', function(e) {
            e.preventDefault(); // Prevent default form submission

            const filename = document.getElementById('filename').value;
            const content = document.getElementById('content').value;

            // Send JSON request using htmx
            htmx.ajax('POST', '/api/documents', {
                target: '#result',
                swap: 'innerHTML',
                headers: {
                    'Content-Type': 'application/json'
                },
                values: {
                    filename: filename,
                    content: content
                }
            }).then(function() {
                // Clear form after successful submission
                document.getElementById('createForm').reset();
            });
        });
    </script>
</body>
</html>`, docList)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlPage)
}

func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := s.getAllDocuments()
	if err != nil {
		http.Error(w, "Failed to list documents", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (s *Server) getDocument(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename parameter is required", http.StatusBadRequest)
		return
	}
	
	docPath := filepath.Join(s.docDir, filename+".md")
	content, err := ioutil.ReadFile(docPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to read document", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(content)
}

// DocumentRequest represents the JSON structure for creating a document
type DocumentRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func (s *Server) createDocument(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req DocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	filename := req.Filename
	content := req.Content

	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	// Sanitize filename
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	docPath := filepath.Join(s.docDir, filename+".md")

	if err := ioutil.WriteFile(docPath, []byte(content), 0644); err != nil {
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	// Return JSON response with document info
	docURL := fmt.Sprintf("%s/%s", strings.TrimRight(s.publicURL, "/"), filename)
	response := map[string]string{
		"message":   fmt.Sprintf("Document '%s' created successfully", filename),
		"filename":  filename,
		"url":       docURL,
		"api_url":   fmt.Sprintf("%s/api/document?filename=%s", strings.TrimRight(s.publicURL, "/"), filename),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) getAllDocuments() ([]string, error) {
	var documents []string
	
	err := filepath.Walk(s.docDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			relPath, err := filepath.Rel(s.docDir, path)
			if err != nil {
				return err
			}
			// Remove .md extension
			docName := strings.TrimSuffix(relPath, ".md")
			documents = append(documents, docName)
		}
		
		return nil
	})
	
	return documents, err
}