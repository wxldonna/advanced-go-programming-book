package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	_IPA_UPLOAD_ADDRESS    = "http://localhost:8080/upload" //"https://ipa.wdf.sap.corp:8243/upload2"
	_IPA_TEST_TYPE         = "Single User"
	_IPA_BASE_COUNTER_NAME = "End to End Response Time [s]##"
	_IPA_USER              = "i324483"
	_IPA_PROJECT           = "gobenchTest"
	_RELEASE               = "2211.1.0"
	_IPA_CHARACTERISTIC    = "BeRT"
	_COMPARATOR            = "median"
	_IPA_CHARSET           = "UTF-8"
)

//password Jed8
func main() {
	uploadUrl := _IPA_UPLOAD_ADDRESS
	//file, err := os.Open("/home/vagrant/go/src/github.wdf.sap.corp/velocity/vflow-sub-abap/8.csv")

	paras := map[string]string{
		"_charset_": _IPA_CHARSET,
		"fileCSVUploader-data": _IPA_PROJECT + "###" + "BenchmarkConverterCSV" + "###" +
			_IPA_CHARACTERISTIC + "###" + "2211.2.0" + "###" +
			"8.csv" + "###" + "2022/04/26" + "###" + "10:52:15" + "###" + "*" + "###" + "*" + "###" + "*",
		"username": "i324483",
		"password": "Jed8",
	}
	/*
		respBody, err := UploadFile(uploadUrl, paras, "fileCSVUploader", "8.csv", file)
		if err != nil {
			panic(err)
		}

	*/
	filePath := "/home/vagrant/Postman/postInstall/postman-9.15.2-linux-x64.tar.gz"
	contentType, r, err := createReqBody(paras, "fileCSVUploader", filePath)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", uploadUrl, r)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(content))

}

//注意client 本身是连接池，不要每次请求时创建client
var (
	HttpClient = &http.Client{
		Timeout: 3 * time.Hour,
	}
)

// 上传文件
// url                请求地址
// params        post form里数据
// nameField  请求地址上传文件对应field
// fileName     文件名
// file               文件
func UploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func createReqBody(params map[string]string, nameField string, filePath string) (string, io.Reader, error) {
	var err error
	pr, pw := io.Pipe()
	bw := multipart.NewWriter(pw) // body writer
	f, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}

	go func() {
		defer f.Close()
		// text part1
		for key, val := range params {
			_ = bw.WriteField(key, val)
		}

		// file part1
		_, fileName := filepath.Split(filePath)
		formFile, err := bw.CreateFormFile(nameField, fileName)
		if err != nil {
			pw.CloseWithError(err)
		}

		var buf = make([]byte, 1024)
		cnt, _ := io.CopyBuffer(formFile, f, buf)
		log.Printf("copy %d bytes from file %s in total\n", cnt, fileName)

		bw.Close() //write the tail boundry
		pw.Close()
	}()
	return bw.FormDataContentType(), pr, nil
}
