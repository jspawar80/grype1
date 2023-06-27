#Scan Docker Images for Vulnerabilities with Grype and trivy
This repository provides a Docker image, entrypoint script and node js API to scan Docker images for vulnerabilities using Grype.

### 1. create the build_grype.sh script 
```
#!/bin/bash


DOCKER_USERNAME='jspawar80'
DOCKER_PASSWORD='dckr_pat_ecD6s3cZVofGBBE7jSAR_ykbxL4'
DOCKER_TAG='interlynk_scanner_grype'


if [ -z "${DOCKER_USERNAME}" ] || [ -z "${DOCKER_PASSWORD}" ]; then
  echo "Docker credentials are not set."
  exit 1
fi

if [ -z "${DOCKER_TAG}" ]; then
  echo "Docker tag is not provided."
  exit 1
fi

LATEST_GRYPE=$(curl --silent "https://api.github.com/repos/anchore/grype/releases/latest" | jq -r '.tag_name')

echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

docker build -t ${DOCKER_USERNAME}/${DOCKER_TAG}:latest .
docker tag ${DOCKER_USERNAME}/${DOCKER_TAG}:latest ${DOCKER_USERNAME}/${DOCKER_TAG}:${LATEST_GRYPE}
docker push ${DOCKER_USERNAME}/${DOCKER_TAG}:latest
docker push ${DOCKER_USERNAME}/${DOCKER_TAG}:${LATEST_GRYPE}
```

DOCKER_USERNAME='jspawar80'
DOCKER_PASSWORD='dckr_pat_ecD6s3cZVofGBBE7jSAR_ykbxL4'
DOCKER_TAG='interlynk_scanner_grype'
change the value of DOCKER_USERNAME, DOCKER_PASSWORD and DOCKER_TAG
DOCKER_TAG is the tag you want to give to your image

### 2. create the build_trivy.sh script 
```
#!/bin/bash
DOCKER_USERNAME='jspawar80'
DOCKER_PASSWORD='dckr_pat_ecD6s3cZVofGBBE7jSAR_ykbxL4'

DOCKER_TAG='interlynk_scanner_trivy'
if [ -z "${DOCKER_USERNAME}" ] || [ -z "${DOCKER_PASSWORD}" ]; then
  echo "Docker credentials are not set."
  exit 1
fi

if [ -z "${DOCKER_TAG}" ]; then
  echo "Docker tag is not provided."
  exit 1
fi

LATEST_TRIVY=$(curl --silent "https://api.github.com/repos/aquasecurity/trivy/releases/latest" | jq -r '.tag_name')
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

docker build -t ${DOCKER_USERNAME}/${DOCKER_TAG}:latest -f Dockerfiletrivy .
docker tag ${DOCKER_USERNAME}/${DOCKER_TAG}:latest ${DOCKER_USERNAME}/${DOCKER_TAG}:${LATEST_TRIVY}
docker push ${DOCKER_USERNAME}/${DOCKER_TAG}:latest
docker push ${DOCKER_USERNAME}/${DOCKER_TAG}:${LATEST_TRIVY}
````
DOCKER_USERNAME='jspawar80'
DOCKER_PASSWORD='dckr_pat_ecD6s3cZVofGBBE7jSAR_ykbxL4'

DOCKER_TAG='interlynk_scanner_trivy'
change the value of DOCKER_USERNAME, DOCKER_PASSWORD and DOCKER_TAG
DOCKER_TAG is the tag you want to give to your image


### 3. changes the path where you want to save the output of scanner in index.js file
in line 17
```
        '/home/ishu/grype/grype3/previous:/output',
```
And line 57
```
            fs.readFile(`/home/ishu/grype/grype3/previous/${outputFile}`, 'utf8', (err, data) => {
```

### 4. Run the build_grype.sh and build_trivy.sh  
```
sudo chmod +x build_grype.sh
sudo chmod +x build_trivy.sh  

sudo build_grype.sh
sudo build_trivy.sh
```

### 5. install the all necessary dependencies for index.js
```
npm init -y
npm install express body-parser child_process fs
```
### 6. open new terminal and run this command to check the api
```
curl -X POST -H "Content-Type: application/json" -d '{
   "DOCKER_USERNAME": "riteshnoronha2022",
   "DOCKER_TOKEN": "dckr_pat_wHRzkibVa0gsYiIV-ue7IxBesO4",
   "SCANNER": "grype",
   "IMAGE_TO_SCAN": "riteshnoronha2022/sbomqs:v0.0.17",
   "IMAGE_OF_SCANNER": "jspawar80/interlynk_scanner_grype"
}' http://localhost:3000/scan
```
