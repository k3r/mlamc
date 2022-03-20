package main

import (
	"net/http"
	"log"
	"crypto/sha256"
	"io"
	"os"
	"io/ioutil"
	"encoding/hex"
	"encoding/json"
	"mime/multipart"
	"bytes"
	"time"
	"strings"
)

const API_VERSION string = "1.0"
const URL string = "https://ml.charrel.fr/api.php"
const MAX_SIZE int = 5000000

func checkAPIVersion(result map[string]string) {
	if (result["version"] != API_VERSION) {
		log.Fatalln(result, "Please update the version of the client. You can download the lastest version on https://github.com/k3r/mlamc .")
	}
}

/*
	Check if the hash of the file is already known by mlam.
*/
func submitHash(filepath string) (string, string) {
	// Process the hash
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	hash := hex.EncodeToString(h.Sum(nil))

	// Send the hash to the service
	resp, err := http.Get(URL + "?action=submit_hash&sha256=" + hash)
	if err != nil {
	   log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatal(err)
	}

	var result map[string]string

	json.Unmarshal(body, &result)

	checkAPIVersion(result)

	return result["status"], result["message"]
}

/*
	Check if the file is malicious.
*/
func submitFile(filepath string, vt bool) (string, string) {
	client := &http.Client{
        Timeout: time.Second * 60,
    }

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// vt value
	fw, err := writer.CreateFormField("vt")
    if err != nil {
    	return "ERROR", "writer.CreateFormField(\"vt\")"
    }

    var vtValue string

    if vt {
    	vtValue = "1"
    } else {
    	vtValue = "0"
    }
    _, err = io.Copy(fw, strings.NewReader(vtValue))
    if err != nil {
        return "ERROR", "io.Copy(fw, strings.NewReader(vtValue))"
    }

    // file value
    fw, err = writer.CreateFormFile("file", filepath)
    if err != nil {
    }
    file, err := os.Open(filepath)
    if err != nil {
        panic(err)
    }
    _, err = io.Copy(fw, file)
    if err != nil {
        return "ERROR", "io.Copy(fw, file)"
    }

    // Close multipart writer.
    writer.Close()
    req, err := http.NewRequest("POST", URL + "?action=submit_file", bytes.NewReader(body.Bytes()))
    if err != nil {
    	log.Println(err)
        return "ERROR", "http.NewRequest"
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    resp, err := client.Do(req)
    if err == nil {
	   	if resp.StatusCode != http.StatusOK {
	        log.Printf("Request failed with response code: %d", resp.StatusCode)
	        return "ERROR", "http.StatusOK"
	    } else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
			   return "ERROR", "ioutil.ReadAll(resp.Body)"
			}

			var result map[string]string

			json.Unmarshal(body, &result)

			checkAPIVersion(result)
			
			return result["status"], result["message"]
		}
    } else {
    	log.Println(err)
    	return "ERROR", "client.Do(req)"
    }
}