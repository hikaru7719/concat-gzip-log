package file

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

type FileWriter struct {
	name string
	file io.WriteCloser
}

func NewFileWriter(name string) *FileWriter {
	return &FileWriter{
		name: name,
	}
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

func (f *FileWriter) Copy(reader io.Reader) {
	r, err := gzip.NewReader(reader)
	if err != nil {
		log.Print(err)
	}
	io.Copy(f.file, r)
	f.file.Write([]byte("\n"))
}

func (f *FileWriter) CopyAll(readers []io.Reader) {
	for _, r := range readers {
		f.Copy(r)
	}
}

func (f *FileWriter) Run(readers []io.Reader) {
	f.Open(f.name)
	defer f.Close()
	f.CopyAll(readers)
}
