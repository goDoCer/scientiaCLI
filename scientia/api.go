package scientia

import (
	"net/http"
	"time"
)

var token string

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

const (
	baseURL = "https://api-materials.doc.ic.ac.uk/"
)

type File struct {
	Name     string
	path     string
	children []*File
	searched bool
	isFolder bool
}

func GetFiles() []File {
	return []File{}
}

func (file *File) Download() error {
	if !file.searched {
		//Perform get request
	}
	if !file.isFolder {
		return nil
	}
	return nil
}
