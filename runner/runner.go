package runner

import (
	"io"
	"time"

	"github.com/hikaru7719/concat-gzip-log/aws"
	"github.com/hikaru7719/concat-gzip-log/file"
)

type Reader interface {
	Run() []io.Reader
	RunP() <-chan io.Reader
}

type Writer interface {
	Run(readers []io.Reader)
	RunP(c <-chan io.Reader)
}

type Runner struct {
	reader Reader
	writer Writer
}

func NewRunner(bucket string, date time.Time, name string) *Runner {
	return &Runner{
		reader: aws.NewStorageReader(bucket, date),
		writer: file.NewFileWriter(name),
	}
}

func (r *Runner) Run() {
	result := r.reader.Run()
	r.writer.Run(result)
}

func (r *Runner) RunP() {
	c := r.reader.RunP()
	r.writer.RunP(c)
}
