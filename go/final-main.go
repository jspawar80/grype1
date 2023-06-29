package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os/exec"
        "strings"
        "time"
        "github.com/google/uuid"

)

type Scan struct {
        DOCKER_USERNAME  string
        DOCKER_TOKEN     string
        SCANNER          string
        IMAGE_TO_SCAN    string
        REQUEST_ID       string
}

type DockerInfo struct {
        CommandOutput string
}

type Response struct {
        Message string
        DockerInfo DockerInfo
	Request_id uuid.UUID
        Results map[string]string
}

func handleScan(w http.ResponseWriter, r *http.Request) {
        decoder := json.NewDecoder(r.Body)
        var t Scan
        err := decoder.Decode(&t)
        if err != nil {
                log.Println("Error decoding request body:", err)
                panic(err)
        }

        cmd := exec.Command("docker", "run", "-v", "/home/ishu/grype/grype3/previous:/output",
                "-v", "/var/run/docker.sock:/var/run/docker.sock",
                "-e", fmt.Sprintf("DOCKER_USERNAME=%s", t.DOCKER_USERNAME),
                "-e", fmt.Sprintf("DOCKER_TOKEN=%s", t.DOCKER_TOKEN),
                "jspawar80/interlynk_scanner_grype",
                t.SCANNER,
                t.IMAGE_TO_SCAN)

        output, err := cmd.CombinedOutput()
        if err != nil {
                log.Printf("Error: %s\n", err)
                http.Error(w, "Error executing Docker command", http.StatusInternalServerError)
                return
        }

        log.Println("Docker command executed successfully")
        outputStr := string(output)

        dockerInfo := DockerInfo{
                CommandOutput: outputStr,
        }

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
                Message: "Scan completed successfully",
                DockerInfo: dockerInfo,
                Request_id: requestID,

                Results: results,
        }

        json.NewEncoder(w).Encode(response)
        log.Println("Response sent successfully")
}

func main() {
        http.HandleFunc("/scan", handleScan)
        log.Println("Starting server on port 3000")
        log.Fatal(http.ListenAndServe(":3000", nil))
}
