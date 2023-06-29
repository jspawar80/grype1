package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Scan struct {
	DOCKER_USERNAME  string
	DOCKER_TOKEN     string
	SCANNER          string
	IMAGE_OF_SCANNER string
	IMAGE_TO_SCAN    string
}

type Response struct {
	Message    string
	Request_id uuid.UUID
	Results    map[string]string
}

func handleScan(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()

	decoder := json.NewDecoder(r.Body)
	var t Scan
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Error decoding request body:", err)
		panic(err)
	}
	log.Println("Request body decoded successfully")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "-v", "/home/ishu/grype/grype3/previous:/output",
		"-e", fmt.Sprintf("DOCKER_USERNAME=%s", t.DOCKER_USERNAME),
		"-e", fmt.Sprintf("DOCKER_TOKEN=%s", t.DOCKER_TOKEN),
		t.IMAGE_OF_SCANNER,
		t.SCANNER,
		t.IMAGE_TO_SCAN)

	_, err = cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, "Error executing docker command", http.StatusInternalServerError)
		return
	}
	log.Println("Docker command executed successfully")

	// Upload the scan results to a server
	uploadFile("/home/ishu/grype/grype3/previous/result.txt", "http://example.com/upload")

	date := time.Now().UTC()
	dateTime := date.Format("2006-01-02")

	var format string
	if t.SCANNER == "grype" {
		format = "txt"
	} else if t.SCANNER == "trivy" {
		format = "json"
	} else {
		http.Error(w, "Invalid scanner", http.StatusBadRequest)
		return
	}
	log.Printf("Scanner is %s, so file format is %s", t.SCANNER, format)

	outputFile := fmt.Sprintf("%s:%s:%s.%s", t.IMAGE_TO_SCAN, dateTime, t.SCANNER, format)
	outputFile = strings.ReplaceAll(outputFile, "/", ":")
	log.Printf("Output file name: %s", outputFile)

	data, err := ioutil.ReadFile(fmt.Sprintf("/home/ishu/grype/grype3/previous/%s", outputFile))

	if err != nil {
		log.Println("File reading error", err)
		http.Error(w, "Error reading scan result files", http.StatusInternalServerError)
		return
	}
	log.Println("File read successfully")

	results := make(map[string]string)
	results[format] = string(data)

	response := Response{
		Message:    "Scan completed successfully",
		Request_id: requestID,
		Results:    results,
	}

	json.NewEncoder(w).Encode(response)
	log.Println("Response sent successfully")
}

func uploadFile(filepath string, url string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	http.HandleFunc("/scan", handleScan)
	log.Println("Starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
