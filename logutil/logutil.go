package logutil

import (
	"fmt"
	"io"
	"log"
)

var (
	i *log.Logger
	e *log.Logger
)

func Init(infoDestination io.Writer, errorDestination io.Writer, serviceName, serviceVersion string) {
	infoPrefix := fmt.Sprint("info | ", serviceName, " | ", "v", serviceVersion, " | ")
	i = log.New(infoDestination,
		infoPrefix,
		log.Ldate|log.Ltime|log.Lshortfile)

	errorPrefix := fmt.Sprint("error | ", serviceName, " | ", serviceVersion, " | ")
	e = log.New(errorDestination,
		errorPrefix,
		log.Ldate|log.Ltime|log.Lshortfile)
}

func SetOutputInfo(w io.Writer) {
	i.SetOutput(w)
}

func SetOutputErr(w io.Writer) {
	e.SetOutput(w)
}

func Info(format string, v ...interface{}) {
	i.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	e.Output(2, fmt.Sprintf(format, v...))
}
