const express = require('express');
const bodyParser = require('body-parser');
const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');

const app = express();

app.use(bodyParser.json());

app.post('/scan', (req, res) => {
    const { DOCKER_USERNAME, DOCKER_TOKEN, SCANNER, IMAGE_OF_SCANNER, IMAGE_TO_SCAN } = req.body;

    const child = spawn('docker', [
        'run',
        '-v',
        '/home/ishu/grype/grype3/previous:/output',
        '-e',
        `DOCKER_USERNAME=${DOCKER_USERNAME}`,
        '-e',
        `DOCKER_TOKEN=${DOCKER_TOKEN}`,
        IMAGE_OF_SCANNER,
        SCANNER,
        IMAGE_TO_SCAN
    ]);

    let output = '';

    child.stdout.on('data', (data) => {
        output += data.toString();
    });

    child.stderr.on('data', (data) => {
        console.error(`Error: ${data}`);
    });

    child.on('exit', (code) => {
        if (code !== 0) {
            return res.status(500).send({ message: 'Error executing docker command' });
        }

        const date = new Date();
        const year = date.getUTCFullYear();
        const month = ("0" + (date.getUTCMonth() + 1)).slice(-2); // Months are zero-based, so +1
        const day = ("0" + date.getUTCDate()).slice(-2);

        const dateTime = `${year}-${month}-${day}`;

        const formats = ['txt', 'json'];
        const results = {};

        let readCount = 0;
        formats.forEach((format) => {
            let outputFile = `${IMAGE_TO_SCAN}:${dateTime}:${SCANNER}.${format}`;
            outputFile = outputFile.replace(/\//g, ':');

            fs.readFile(`/home/ishu/grype/grype3/previous/${outputFile}`, 'utf8', (err, data) => {
                readCount++;
                if (!err) {
                    results[format] = data;
                }
                
                if(readCount === formats.length) {
                    if (Object.keys(results).length === 0) {
                        return res.status(500).send({ message: 'Error reading scan result files' });
                    }
                    res.send({ message: 'Scan completed successfully', results: results });
                }
            });
        });
    });
});

app.listen(3000, () => {
    console.log('Server started on port 3000');
});
