package managers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type Composer struct {
	Files Files
}

func ComposerDefault() PackageManager {
	c := &Composer{}
	p := PackageManager(c)
	return p
}

func (c *Composer) Filter(files Files) {
	for _, file := range files {
		if file.Name == "composer.lock" {
			c.Match(file)
		}
	}
}

func (pm *Composer) Match(f File) {
	pm.Files = append(pm.Files, f)
}

func (c *Composer) Report() {
	for _, file := range c.Files {
		c.CheckFile(file)
	}
}

func (c *Composer) GetFiles() Files {
	return c.Files
}

func (c *Composer) SetFiles(files Files) {
	c.Files = files
}

func (c *Composer) CheckFile(file File) {
	url := "https://security.sensiolabs.org/check_lock"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("lock", file.Name)
	if err != nil {
		log.Fatal(err)
	}
	f := bytes.NewBufferString(file.Data)
	_, err = io.Copy(part, f)

	//	_ = writer.WriteField("lock", file.Data)
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Length", string(req.ContentLength))

	//	fmt.Println("response Status:", req.Body)
	//	fmt.Println("response Headers:", req.Header)
	//	b, _ := ioutil.ReadAll(req.Body)
	//	fmt.Println("response Body:", string(b))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(b))
}
