package handlers

import (
	"fmt"
	"github.com/pablogolobaro/pdfcomposer"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

const topic = "notifications"

type Handler struct {
	HttpClient       *http.Client
	PdfComposeCLient pdfcomposer.PdfComposeClient
}
type fileWithName struct {
	file io.ReadCloser
	name string
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(256)
	mForm := r.MultipartForm
	files := []fileWithName{}
	for k, _ := range mForm.File {
		// k is the key of file part
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			fmt.Println("inovke FormFile error:", err)
			return
		}
		fmt.Printf("the uploaded file: name[%s], size[%d], header[%#v]\n",
			fileHeader.Filename, fileHeader.Size, fileHeader.Header)
		files = append(files, fileWithName{name: fileHeader.Filename, file: file})
		defer file.Close()
	}
	pdfFile, err := uploadImage(h.PdfComposeCLient, files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := ioutil.ReadAll(pdfFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteToTopic(topic, "Got PDF file")
	w.Write(bytes)

	return
}

func (h *Handler) Web(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/form.html")
	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, nil)
}
