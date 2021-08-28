package file

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

type FileWriter struct {
	file io.WriteCloser
}

func (f *FileWriter) Open(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
	}
	f.file = file
}

func (f *FileWriter) Close() {
	f.file.Close()
}

func (f *FileWriter) Copy(reader io.ReadCloser) {
	defer reader.Close()
	r, err := gzip.NewReader(reader)
	if err != nil {
		log.Print(err)
	}
	io.Copy(f.file, r)
}

func (f *FileWriter) CopyAll(readers []io.ReadCloser) {
	for _, r := range readers {
		f.Copy(r)
	}
}

func (f *FileWriter) Run(path string, readers []io.ReadCloser) {
	f.Open(path)
	defer f.Close()
	f.CopyAll(readers)
}
