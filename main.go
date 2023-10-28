package main

import (
	WebAssets "github.com/S-IR/grit-template/lib/compile"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Example function that accepts only CSS files
func ProcessCSSFile(file []byte) error {
	// Process CSS file here
	return nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	// Initialize Gin
	r := gin.Default()

	defer WebAssets.Ctx.Dispose()

	WebAssets.HandleSendReact(r)
	// Serve React App (assuming built React app is in a folder named "build")

	// Start Gin on port 3000
	err = r.Run(":3000")
	if err != nil {
		panic(err)
	}

}
