#Scan Docker Images for Vulnerabilities with Grype and trivy

This repository contains a Docker image, entrypoint script, and Node.js API to enable vulnerability scanning of Docker images using Grype and Trivy.

Here are the steps required to set up and run the vulnerability scanner:


### 1. create the build_grype.sh script 

The build_grype.sh script will build and push Docker images using the Grype scanner.

Use the following script, ensuring to replace DOCKER_USERNAME, DOCKER_PASSWORD, and DOCKER_TAG with your own credentials and desired image tag.

```
#!/bin/bash


DOCKER_USERNAME='username_of_dockerhub'
DOCKER_PASSWORD='password_of_dockerhub'
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


### 2. create the build_trivy.sh script 

The build_trivy.sh script will build and push Docker images using the Trivy scanner.

Use the following script, making sure to replace DOCKER_USERNAME, DOCKER_PASSWORD, and DOCKER_TAG with your own credentials and desired image tag.

```
#!/bin/bash
DOCKER_USERNAME='username_of_dockerhub'
DOCKER_PASSWORD='password_of_dockerhub'
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

### 3. change the path where you want to save the output of scanner in index.js file
in line 17
```
        '/path/of/output:/output',
```
And line 57
```
            fs.readFile(`/path/of/output/${outputFile}`, 'utf8', (err, data) => {
```

### 4. Run the build_grype.sh and build_trivy.sh  

Make build_grype.sh and build_trivy.sh executable and run them:
```
sudo chmod +x build_grype.sh
sudo chmod +x build_trivy.sh  

sudo ./build_grype.sh
sudo ./build_trivy.sh
```

### 5. install the all necessary dependencies for index.js

Install the necessary dependencies for index.js using the following commands:
```
npm init -y
npm install express body-parser child_process fs
node index.js
```
### 6. open new terminal and run this command to Test the API
Open a new terminal and run the following command to check if the API is functioning correctly:

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

Replace your_username, your_password, your_image_to_scan, and your_scanner_image with your own values.
