package fiddleware

import (
	"github.com/bernsblack/fiddleware/util"
)

type Recording struct {
	RequestId    string
	ResponseDump []byte
	RequestDump  []byte
}

type InMemoryRecorder struct {
	logger     util.Logger
	innerStore map[string]*Recording // might want to replace this db
}

func NewInMemoryRecorder(logger util.Logger) *InMemoryRecorder {
	return &InMemoryRecorder{
		logger:     logger,
		innerStore: make(map[string]*Recording, 0), // change to sync map or db later on
	}
}

func (r *InMemoryRecorder) Record(key string, value interface{}) {
	if recording, ok := r.innerStore[key]; ok {
		if recording.RequestDump == nil {
			recording.RequestDump = value.([]byte)
		} else {
			recording.ResponseDump = value.([]byte)
		}

		r.innerStore[key] = recording
	} else {
		r.innerStore[key] = &Recording{
			RequestId:   key,
			RequestDump: value.([]byte),
		}
	}
}

type Recorder interface {
	Record(key string, value interface{})
}
