package main

import (
	"log"
	"net/http"
	"os"

	js "github.com/S-IR/grit-template/lib/compile"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type HTMLFile = []byte
type CSSFile = []byte

// Example function that accepts only CSS files
func ProcessCSSFile(file CSSFile) error {
	// Process CSS file here
	return nil
}

var (
	html     HTMLFile
	css      CSSFile
	jsBundle []byte
	reloadWS *websocket.Conn
)

func main() {

	var err error
	html, err = os.ReadFile("./public/index.html")

	if err != nil {
		panic(err)
	}
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Initialize Gin
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		reloadWS, err = upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}
	})
	go handleCompilation(reloadWS)

	r.GET("/bundle.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/javascript", jsBundle)
	})

	r.Static("/public", "./public")

	r.GET("/styles.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css", css)
	})
	// Serve React App (assuming built React app is in a folder named "build")
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
	})

	// Start Gin on port 3000
	err = r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

func handleCompilation(reloadWS *websocket.Conn) {
	defer reloadWS.Close()

	js.Recompile(&css, &jsBundle)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add("./src")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// File has been modified, trigger recompilation
				js.Recompile(&css, &jsBundle) // A function where you use esbuild API to recompile
				err := reloadWS.WriteJSON(map[string]interface{}{
					"type": "update",
					"html": string(html),
					"css":  string(css),
					"js":   string(jsBundle),
				})
				if err != nil {
					panic(err)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}
