package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	
)

type EnvItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// envValue returns the environment variable for the given key, or "Not set" when empty.
func envValue(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "Not set"
	}
	return value
}

// loadDotEnv reads key=value pairs from the given file and sets them if not already present.
func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

func main() {
	loadDotEnv(".env")

	router := gin.Default()

	publicKeys := []string{"APP_NAME", "API_URL", "ENVIRONMENT", "VERSION"}

	router.SetHTMLTemplate(template.Must(template.New("index").Parse(indexTemplate)))

	router.GET("/", func(c *gin.Context) {
		items := make([]EnvItem, 0, len(publicKeys))
		for _, key := range publicKeys {
			items = append(items, EnvItem{Key: key, Value: envValue(key)})
		}

		c.HTML(http.StatusOK, "index", gin.H{
			"title": "Server Compass Demo Environment Variables",
			"envs":  items,
		})
	})

	router.GET("/api/env", func(c *gin.Context) {
		items := make([]EnvItem, 0, len(publicKeys))
		for _, key := range publicKeys {
			items = append(items, EnvItem{Key: key, Value: envValue(key)})
		}
		c.JSON(http.StatusOK, gin.H{"envs": items})
	})

	// The private variables are intentionally not exposed to the browser; they would be used server-side only.

	router.Run()
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ .title }}</title>
    <style>
        body { font-family: Arial, sans-serif; background: #0b1622; color: #e8eef5; margin: 0; padding: 0; }
        header { padding: 24px 32px; background: linear-gradient(135deg, #0f1f3a, #0b1622); border-bottom: 1px solid #1f2d3d; }
        h1 { margin: 0; font-size: 24px; letter-spacing: 0.2px; }
        main { padding: 32px; max-width: 720px; margin: 0 auto; }
        section { background: #101c2b; border: 1px solid #1f2d3d; border-radius: 12px; padding: 24px; box-shadow: 0 12px 30px rgba(0, 0, 0, 0.25); }
        .env { display: flex; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #1b2634; }
        .env:last-child { border-bottom: none; }
        .key { font-weight: 600; color: #7fb4ff; }
        .value { color: #f4f7fb; }
        .not-set { color: #f7b733; }
        p { margin-top: 8px; color: #9fb3c8; line-height: 1.5; }
        .api-hint { margin-top: 16px; font-size: 14px; color: #7fb4ff; }
        @media (max-width: 640px) {
            main { padding: 24px; }
            section { padding: 20px; }
            .env { flex-direction: column; gap: 6px; }
        }
    </style>
</head>
<body>
    <header>
        <h1>{{ .title }}</h1>
        <p>Only public variables are shown here - Test domain after auto commit Private server values stay on the backend. - test new branch</p>
    </header>
    <main>
        <section>
            {{ range .envs }}
            <div class="env">
                <div class="key">{{ .Key }}</div>
                {{ if eq .Value "Not set" }}
                <div class="value not-set">{{ .Value }}</div>
                {{ else }}
                <div class="value">{{ .Value }}</div>
                {{ end }}
            </div>
            {{ end }}
            <div class="api-hint">Try the JSON view at <code>/api/env</code>.</div>
        </section>
    </main>
</body>
</html>`
