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
        IMAGE_OF_SCANNER string
}

type DockerInfo struct {
        CommandOutput string
}

type DockerVersionInfo struct {
        Client DockerClient `json:"Client"`
        Server DockerServer `json:"Server"`
}

type DockerClient struct {
        Version string `json:"Version"`
}

type DockerServer struct {
        Engine DockerEngine `json:"Engine"`
}

type DockerEngine struct {
        Version string `json:"Version"`
}

type GrypeVersionInfo struct {
        Application string `json:"Application"`
        Version     string `json:"Version"`
}

type ImageLayer struct {
        Image   string `json:"IMAGE"`
        Created string `json:"CREATED"`
        CreatedBy string `json:"CREATED BY"`
        Size string `json:"SIZE"`
        Comment string `json:"COMMENT"`
}

type Response struct {
        Message        string
        RequestId      string
        DockerInfo     DockerInfo
        Results        map[string]interface{}
        DockerVersion  DockerVersionInfo
        GrypeVersion   GrypeVersionInfo
        ImageLayerInfo []ImageLayer
}

func handleScan(w http.ResponseWriter, r *http.Request) {
        decoder := json.NewDecoder(r.Body)
        var t Scan
        err := decoder.Decode(&t)
        if err != nil {
                log.Println("Error decoding request body:", err)
                panic(err)
        }

        requestId := uuid.New().String()

        cmd := exec.Command("docker", "run", "-v", "/home/ishu/grype/grype3/previous:/output",
                "-v", "/var/run/docker.sock:/var/run/docker.sock",
                "-e", fmt.Sprintf("DOCKER_USERNAME=%s", t.DOCKER_USERNAME),
                "-e", fmt.Sprintf("DOCKER_TOKEN=%s", t.DOCKER_TOKEN),
                t.IMAGE_OF_SCANNER,
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

        // Parse CommandOutput here to extract versions and image layer info
        dockerVersion, grypeVersion, imageLayerInfo := parseOutput(outputStr)

        date := time.Now().UTC()
        dateTime := date.Format("2006-01-02")
        var format string
        if t.SCANNER == "grype" {
                format = "json"
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

        results := make(map[string]interface{})
        err = json.Unmarshal(data, &results)
        if err != nil {
                log.Println("Error decoding file contents:", err)
                http.Error(w, "Error decoding file contents", http.StatusInternalServerError)
                return
        }

        resp := Response{
                Message:        "Scanning request processed successfully",
                RequestId:      requestId,
                DockerInfo:     dockerInfo,
                Results:        results,
                DockerVersion:  dockerVersion,
                GrypeVersion:   grypeVersion,
                ImageLayerInfo: imageLayerInfo,
        }

        respBytes, err := json.Marshal(resp)
        if err != nil {
                log.Println("Error encoding response to JSON:", err)
                http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
                return
        }

        fmt.Fprintf(w, string(respBytes))
}

func main() {
        http.HandleFunc("/scan", handleScan)
        log.Println("Starting server on port 8080")
        log.Fatal(http.ListenAndServe(":8080", nil))
}

// Placeholder function for parsing the Docker command output
func parseOutput(outputStr string) (DockerVersionInfo, GrypeVersionInfo, []ImageLayer) {
        // TODO: Parse the outputStr to extract Docker and Grype versions and image layer info
        return DockerVersionInfo{}, GrypeVersionInfo{}, []ImageLayer{}
}
