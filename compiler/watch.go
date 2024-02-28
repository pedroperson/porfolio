package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					// Replace "other_script.go" with the path to your Go script
					// // cmd := exec.Command("go", "run", "other_script.go")
					// // output, err := cmd.CombinedOutput()
					// if err != nil {
					// 	fmt.Println("Error running script:", err)
					// }
					// fmt.Println("Script output:", string(output))

					Compile()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Replace "/path/to/target/folder" with your target folder
	err = recursiveAdd(watcher, "templates")

	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// recursiveAdd watches all subdirectories of the given directory path
func recursiveAdd(watcher *fsnotify.Watcher, basePath string) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		return watcher.Add(path)
	})
}

type PageData struct {
	Title   string
	Content template.HTML // Allows for HTML content
}

func Compile() {
	// Define the base directories
	baseTemplateDir := "templates/pages"
	baseLayoutDir := "templates" // Adjust if your directory structure is different
	baseOutputDir := "public"

	// Combine all layout related files into a single template object
	// This should match your layout, header, and footer files
	layoutPattern := filepath.Join(baseLayoutDir, "*.html")

	layoutTemplates, err := template.ParseGlob(layoutPattern)
	if err != nil {
		panic(err)
	}

	// Process each page template
	err = filepath.Walk(baseTemplateDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		// Read and prepare the page-specific content
		pageContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Set up the full page data
		fullPageData := PageData{
			Title:   "Your Page Title",          // Change as needed
			Content: template.HTML(pageContent), // Directly use read content
		}

		// Construct the output file path
		relPath, err := filepath.Rel(baseTemplateDir, path)
		if err != nil {
			return err
		}
		outputFilePath := filepath.Join(baseOutputDir, strings.Replace(relPath, filepath.Ext(relPath), ".html", -1))

		// Create and open the output file
		file, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Execute the layout template with the full page content
		err = layoutTemplates.ExecuteTemplate(file, "layout.html", fullPageData)
		if err != nil {
			return err
		}

		fmt.Println("Compiled to html:", path)

		return nil
	})

	if err != nil {
		panic(err)
	}
}
