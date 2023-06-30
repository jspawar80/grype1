install golang and run this command before running the api
```
go get github.com/google/uuid
```

```
docker run -v /home/ishu/grype/grype3/previous:/output -v /var/run/docker.sock:/var/run/docker.sock -e DOCKER_USERNAME=riteshnoronha2022 -e DOCKER_TOKEN=dckr_pat_wHRzkibVa0gsYiIV-ue7IxBesO4 jspawar80/interlynk_scanner_grype grype riteshnoronha2022/sbomqs:v0.0.17
```
GRYPE

```
curl -X POST -H "Content-Type: application/json" -d '{
   "DOCKER_USERNAME": "riteshnoronha2022",
   "DOCKER_TOKEN": "dckr_pat_wHRzkibVa0gsYiIV-ue7IxBesO4",
   "SCANNER": "grype",
   "IMAGE_TO_SCAN": "riteshnoronha2022/sbomqs:v0.0.17",
   "IMAGE_OF_SCANNER": "jspawar80/interlynk_scanner_grype"
}' http://localhost:3000/scan
```

TRIVY
```
curl -X POST -H "Content-Type: application/json" -d '{
   "DOCKER_USERNAME": "riteshnoronha2022",
   "DOCKER_TOKEN": "dckr_pat_wHRzkibVa0gsYiIV-ue7IxBesO4",
   "SCANNER": "trivy",
   "IMAGE_TO_SCAN": "riteshnoronha2022/sbomqs:v0.0.17",
   "IMAGE_OF_SCANNER": "jspawar80/interlynk_scanner_trivy"
}' http://localhost:3000/scan
```
