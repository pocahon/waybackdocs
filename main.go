package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// downloadTask holds the information for each download.
type downloadTask struct {
	timestamp   string
	originalURL string
}

// worker processes download tasks from the channel.
func worker(id int, tasks <-chan downloadTask, outputDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	for task := range tasks {
		// Wait 10 seconds before starting this download.
		time.Sleep(10 * time.Second)

		downloadURL := fmt.Sprintf("https://web.archive.org/web/%s/%s", task.timestamp, task.originalURL)
		fmt.Printf("[Worker %d] Downloading: %s\n", id, downloadURL)

		req, err := http.NewRequest("GET", downloadURL, nil)
		if err != nil {
			log.Printf("[Worker %d] Error creating request for %s: %v", id, downloadURL, err)
			continue
		}
		// Set headers to simulate curl behavior.
		req.Header.Set("User-Agent", "curl/7.64.1")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Connection", "keep-alive")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Worker %d] Error downloading %s: %v", id, downloadURL, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("[Worker %d] Download failed for %s: status %s", id, downloadURL, resp.Status)
			resp.Body.Close()
			continue
		}

		// Determine the filename; by default, use the basename from the original URL.
		fileName := path.Base(task.originalURL)
		// Try to get the filename from the Content-Disposition header, if available.
		if cd := resp.Header.Get("Content-Disposition"); cd != "" {
			if _, params, err := mime.ParseMediaType(cd); err == nil {
				if fname, ok := params["filename"]; ok {
					fileName = fname
				}
			}
		}

		outFilePath := filepath.Join(outputDir, fileName)
		outFile, err := os.Create(outFilePath)
		if err != nil {
			log.Printf("[Worker %d] Could not create file %s: %v", id, outFilePath, err)
			resp.Body.Close()
			continue
		}

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			log.Printf("[Worker %d] Error writing to file %s: %v", id, outFilePath, err)
		}
		outFile.Close()
		resp.Body.Close()

		fmt.Printf("[Worker %d] Completed: %s\n", id, outFilePath)
	}
}

func main() {
	// Command-line parameters.
	domain := flag.String("d", "", "Target domain (e.g. example.com)")
	maxDownloads := flag.Int("n", 0, "Maximum downloads (0 for unlimited)")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Usage: waybackdocs -d <domain> [-n <max downloads>]")
		os.Exit(1)
	}

	// Create the output directory if it doesn't exist.
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Build the CDX API URL.
	cdxURL := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=*.%s&collapse=urlkey", *domain)
	fmt.Println("Fetching list from:", cdxURL)

	// Get the list from the Wayback Machine.
	resp, err := http.Get(cdxURL)
	if err != nil {
		log.Fatalf("Error fetching CDX API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch CDX API: status %s", resp.Status)
	}

	// Read download tasks from the CDX response.
	var tasksSlice []downloadTask
	scanner := bufio.NewScanner(resp.Body)
	downloadCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		timestamp := fields[1]
		originalURL := fields[2]
		lowerURL := strings.ToLower(originalURL)
		// Only include .doc, .docx, and .pdf files (exclude .txt and .eml).
		if !(strings.HasSuffix(lowerURL, ".doc") ||
			strings.HasSuffix(lowerURL, ".docx") ||
			strings.HasSuffix(lowerURL, ".pdf")) {
			continue
		}
		tasksSlice = append(tasksSlice, downloadTask{timestamp: timestamp, originalURL: originalURL})
		downloadCount++
		if *maxDownloads > 0 && downloadCount >= *maxDownloads {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading CDX response: %v", err)
	}
	fmt.Printf("Total tasks: %d\n", len(tasksSlice))

	// Set up a worker pool (5 concurrent downloads).
	numWorkers := 5
	tasksChan := make(chan downloadTask)
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasksChan, outputDir, &wg)
	}

	// Send tasks to the workers.
	for _, task := range tasksSlice {
		tasksChan <- task
	}
	close(tasksChan)
	wg.Wait()

	fmt.Printf("Done! Total downloaded files: %d\n", downloadCount)
}
