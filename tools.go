package http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/Congenital/log/v0.2/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	DisableCompression: true,
}

var client = &http.Client{Transport: tr}

func HttpPost(url string, param string) ([]byte, error) {
	resp, err := client.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(param))

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, errors.New("Err - " + resp.Status)
	}

	return body, nil
}

func HttpPostJson(url string, buff []byte) ([]byte, error) {
	bf := bytes.NewBuffer(*(*[]byte)(unsafe.Pointer(&buff)))

	resp, err := http.Post(url,
		"application/json;charset=utf-8",
		bf)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Err - " + resp.Status)
	}

	return body, nil
}

func HttpGet(url string, param string) ([]byte, error) {
	resp, err := client.Get(url + "?" + param)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Err - " + resp.Status)
	}

	return body, nil
}

func HttpDo(method string, url string, param string) ([]byte, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(param))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Err - " + resp.Status)
	}

	return body, nil
}

func UploadFile(url string, fields []string, fieldsvalue []string, filefield []string, filename []string) ([]byte, error) {
	buf := new(bytes.Buffer)

	w := multipart.NewWriter(buf)

	for i := 0; i < len(fields); i++ {
		w.WriteField(fields[i], fieldsvalue[i])
	}

	for i := 0; i < len(filefield); i++ {
		fw, err := w.CreateFormFile(filefield[i], filename[i])

		if err != nil {
			log.Error(err)
			return nil, err
		}

		fd, err := os.Open(filename[i])
		if err != nil {
			log.Error(err)
			return nil, err
		}

		defer fd.Close()

		_, err = io.Copy(fw, fd)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	w.Close()

	resp, err := client.Post(url, w.FormDataContentType(), buf)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Err - " + resp.Status)
	}

	return buff, nil
}
