package WebAssets

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var hotReloadMu sync.Mutex
var (
	Html     *[]byte
	Css      *[]byte
	JsBundle *[]byte
	Assets   map[string]*[]byte
	ReloadWS *websocket.Conn
)

func HandleSendReact(r *gin.Engine) {
	Html = &[]byte{}
	Css = &[]byte{}
	JsBundle = &[]byte{}
	Assets = make(map[string]*[]byte)

	//initial build
	HtmlGotten, err := os.ReadFile("./public/index.html")
	*Html = HtmlGotten

	Assets = make(map[string]*[]byte)

	if err != nil {
		panic(err)
	}
	Rebuild()

	r.GET("/:assetName", func(c *gin.Context) {
		assetName := c.Param("assetName")
		fmt.Println("assetName", assetName)
		if assetData, exists := Assets[assetName]; exists {
			fmt.Println("exists in assets")

			c.Data(http.StatusOK, "image/svg+xml", *assetData)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/Html; charset=utf-8", *Html)
	})

	r.GET("/bundle.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/javascript", *JsBundle)
	})

	r.Static("/public", "./public")

	r.GET("/styles.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css", *Css)
	})

	SERVER_ENV := os.Getenv("SERVER_ENV")

	if SERVER_ENV == "development" {
		r.GET("/ws", func(c *gin.Context) {

			var upgrader = websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}
			reloadWS, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				panic(err)
			}
			go handleHotReload(r, c, reloadWS)

		})
	}

}

func handleHotReload(r *gin.Engine, c *gin.Context, reloadWS *websocket.Conn) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	err = watchDir("./src", watcher)
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
				ext := filepath.Ext(event.Name)
				if !isValidExtension(ext) {
					continue
				}
				UpdateAssets(reloadWS)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}
func watchDir(path string, watcher *fsnotify.Watcher) error {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
func UpdateAssets(reloadWS *websocket.Conn) {
	hotReloadMu.Lock()
	defer hotReloadMu.Unlock()
	fmt.Println("updating assets")

	updatedFiles, updatedAssets, buildErrors := Rebuild()

	if buildErrors != nil {
		fmt.Println("buildErrors", buildErrors)

		message := map[string]interface{}{
			"type":   "error",
			"errors": buildErrors,
		}
		err := reloadWS.WriteJSON(message)
		if err != nil {
			panic(err)
		}
	}

	for _, fileType := range updatedFiles {
		switch fileType {
		case "js":
			sendFile(reloadWS, fileType, *JsBundle)
		case "css":
			sendFile(reloadWS, fileType, *Css)
		case "assets":
			message := map[string]interface{}{
				"type":       "update",
				"assetPaths": updatedAssets,
			}
			err := reloadWS.WriteJSON(message)
			if err != nil {
				panic(err)
			}

		}
	}

}

func isValidExtension(ext string) bool {
	return ext == ".js" || ext == ".jsx" || ext == ".tsx" || ext == ".ts" || ext == ".css"
}
func sendFile(reloadWS *websocket.Conn, assetType string, assetData []byte) {
	header := []byte(assetType + ":") // Header to specify the type of asset
	finalData := append(header, assetData...)

	err := reloadWS.WriteMessage(websocket.BinaryMessage, finalData)
	if err != nil {
		panic(err)
	}
}
