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
	"regexp"
	"strings"
)

const API_VERSION string = "1.0"
const URL string = "https://ml.charrel.fr/api.php"
const MAX_SIZE int = 5000000

func submitHash(filepath string) (string, string, string) {
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

	return result["status"], result["result"], result["vt_positives"]
}

func submitFile(filepath string, vt bool) string {
	client := &http.Client{
        Timeout: time.Second * 60,
    }

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormField("vt")
    if err != nil {
    }

    var vtValue string

    if vt {
    	vtValue = "1"
    } else {
    	vtValue = "0"
    }
    _, err = io.Copy(fw, strings.NewReader(vtValue))
    if err != nil {
        return "ERROR"
    }

    fw, err = writer.CreateFormFile("file", filepath)
    if err != nil {
    }
    file, err := os.Open(filepath)
    if err != nil {
        panic(err)
    }
    _, err = io.Copy(fw, file)
    if err != nil {
        return "ERROR"
    }

    // Close multipart writer.
    writer.Close()
    req, err := http.NewRequest("POST", URL, bytes.NewReader(body.Bytes()))
    if err != nil {
    	log.Println(err)
        return "ERROR"
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rsp, err := client.Do(req)
    if err == nil {
	   	if rsp.StatusCode != http.StatusOK {
	        log.Printf("Request failed with response code: %d", rsp.StatusCode)
	        return "ERROR 2"
	    } else {
		    b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				log.Println(err)
				return "ERROR 1"
			}
			//Convert the body to type string

			// <p id="result">Your file seems clean.</p>
			r, _ := regexp.Compile("<p id=\"result\">(.+)</p>")

			match := r.FindStringSubmatch(string(b))

			if match != nil && len(match) > 1 {
				return match[1]
			} else {
				r, _ := regexp.Compile("<div class=\"alert alert-danger\" id=\"errors\">[\t\n\f\r ]+<ul>(.+)</ul>[\t\n\f\r ]+</div>")

				match := r.FindStringSubmatch(string(b))
				if match != nil && len(match) > 1 {
					return match[1]
				} else {
					return string(b)
				}
			}
	    }
    } else {
    	log.Println(err)
    	return "ERROR 3"
    }
}