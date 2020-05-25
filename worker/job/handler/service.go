package handler

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/kenza-ai/kenza/worker/job"
	"github.com/kenza-ai/kenza/worker/job/handler/sagify"
	yaml "gopkg.in/yaml.v2"
)

// Service determines the job's service provider e.g. "sagify".
type Service struct {
	next job.Handler
}

type buildConfig struct {
	Sagify sagify.BuildConfiguration
}

// Errors encountered by the Service handler.
var (
	// ErrUnknownServiceProvider is returned when the service provider info cannot be extracted from the build file
	ErrUnknownServiceProvider = errors.New(`Cannot determine service provider e.g. "sagify", ensure the build file is valid`)
)

// NewService — Service constructor
func NewService() *Service {
	return &Service{}
}

// Handle - Handler implementation
//
// Parses build file to service info.
func (h *Service) Handle(r *job.Request) {
	data, err := ioutil.ReadFile(filepath.Join(r.WorkDir, ".kenza.yml"))
	if err != nil {
		r.Fail(err)
		return
	}

	config := &buildConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		r.Fail(err)
		return
	}

	if sage := config.Sagify; (sage != sagify.BuildConfiguration{}) {
		r.Service = "sagify"

		if err := r.Notify(); err != nil {
			e(err.Error())
		}

		sagify.NewSagify(sage).Handle(r)

		if h.next != nil {
			h.Handle(r)
		}
		return
	}

	r.Fail(ErrUnknownServiceProvider) // currently only supporting Sagify
}

// SetNext sets the next Handler
func (h *Service) SetNext(next job.Handler) {
	h.next = next
}
