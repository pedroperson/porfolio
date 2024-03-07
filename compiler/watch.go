package main

import (
	"fmt"
	"html/template"
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

					const layoutPath = "templates/layout.html"
					const templatesDir = "templates"
					const pagesDir = "templates/pages"
					const publicDir = "public"

					CompileTemplates(layoutPath, templatesDir, pagesDir, publicDir)
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

// compileAllTemplates compiles the layout and all other component templates.
func compileAllTemplates(baseLayout string, dirs ...string) (*template.Template, error) {
	tmpl := template.New("layout")
	var err error

	// Parse the base layout first
	tmpl, err = tmpl.ParseFiles(baseLayout)
	if err != nil {
		return nil, err
	}

	// Parse all other templates from specified directories
	for _, dir := range dirs {
		pattern := filepath.Join(dir, "*.html")
		tmpl, err = tmpl.ParseGlob(pattern)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// CompileTemplates compiles HTML templates from the pages directory using a common layout and additional components.
func CompileTemplates(layoutPath, templatesDir, pagesDir, publicDir string) {
	// Compile layout and all component templates
	tmpl, err := compileAllTemplates(layoutPath, filepath.Join(templatesDir, "components"), filepath.Join(templatesDir, "icons"))
	if err != nil {
		panic(err) // Replace with proper error handling
	}

	// Iterate over each HTML file in the pages directory
	files, err := os.ReadDir(pagesDir)
	if err != nil {
		panic(err) // Replace with proper error handling
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".html") {
			continue
		}

		pageName := file.Name()
		pagePath := filepath.Join(pagesDir, pageName)

		// Parse the page-specific content as a new template
		pageContent, err := os.ReadFile(pagePath)
		if err != nil {
			panic(err)
		}

		pageTmpl, err := tmpl.Clone()
		if err != nil {
			panic(err)
		}

		_, err = pageTmpl.New("content").Parse(string(pageContent))
		if err != nil {
			panic(err)
		}

		// Create the output file in the public directory
		outputPath := filepath.Join(publicDir, pageName)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		// Execute the combined template with the layout
		err = pageTmpl.ExecuteTemplate(outputFile, "layout", pageData())
		if err != nil {
			panic(err)
		}

		println("Compiled:", pageName)
	}
}

type Project struct {
	Name      string
	URL       string
	Role      string
	Business  string
	Tasks     string
	Image1    string
	Image1Alt string
	Image2    string
	Image2Alt string
	Color     string
}

func pageData() interface{} {
	return struct {
		Projects []Project
		Title    string
		Content  template.HTML
	}{
		Projects: []Project{
			{
				Name:      "bitbu",
				URL:       "https://bitbu.io/",
				Role:      "Co-founder",
				Business:  "Online tools for musician collaboration",
				Tasks:     "everything tech",
				Image1:    "/bitbu-cover.jpg",
				Image1Alt: "screenshot of ludlow kingsley's website",
				Image2:    "/bitbu-cover.jpg",
				Image2Alt: "screenshot of ludlow kingsley's website",
				Color:     "[13,13,13]",
			},
			{
				Name:      "Ludlow Kingsley",
				URL:       "https://ludlowkingsley.com/",
				Role:      "Frontend Developer",
				Business:  "Corporate design agency",
				Tasks:     "built websites for brands designed in-house",
				Image1:    "/ludlow_home.jpg",
				Image1Alt: "screenshot of ludlow kingsley's website",
				Image2:    "/ludlow_project.jpg",
				Image2Alt: "screenshot of ludlow kingsley's website",
			},
			{
				Name:      "Jerde",
				URL:       "https://jerde.com/",
				Role:      "Frontend Developer",
				Business:  "Architecture design firm",
				Tasks:     "wrote the frontend for the visually striking Ludlow Kingsley design",
				Image1:    "/bitbu-cover.jpg",
				Image1Alt: "screenshot of ludlow kingsley's website",
				Image2:    "/bitbu-cover.jpg",
				Image2Alt: "screenshot of ludlow kingsley's website",
			},
			{
				Name:      "Heloisa Prieto",
				URL:       "https://heloisaprieto.com/?lang=english",
				Role:      "Fullstack developer",
				Business:  "Prolific Brazilian writer",
				Tasks:     "wrote front and backends for eugênia hanitzsch's design",
				Image1:    "/bitbu-cover.jpg",
				Image1Alt: "screenshot of ludlow kingsley's website",
				Image2:    "/bitbu-cover.jpg",
				Image2Alt: "screenshot of ludlow kingsley's website",
			},
			{
				Name:      "JM Agency",
				URL:       "https://jm.agency/",
				Role:      "Fullstack developer",
				Business:  "Prolific Brazilian writer",
				Tasks:     "wrote custom wordpress theme for eugênia hanitzsch's design",
				Image1:    "/bitbu-cover.jpg",
				Image1Alt: "screenshot of ludlow kingsley's website",
				Image2:    "/bitbu-cover.jpg",
				Image2Alt: "screenshot of ludlow kingsley's website",
			},
		},
	}
}
