package web

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/SXUOJ/judge/worker"
	"github.com/google/uuid"
)

const (
	filePerm = 0644
	dirPerm  = 0755
)

type Submit struct {
	SourceCode string `json:"source_code"`
	CodeType   string `json:"code_type"`

	// sample
	Samples []Sample `json:"samples"`

	AllowProc bool `json:"allow_proc"`

	// Limit
	TimeLimit     uint64 `json:"time_limit"`
	RealTimeLimit uint64 `json:"real_time_limit"`
	MemoryLimit   uint64 `json:"memory_limit"`
	OutputLimit   uint64 `json:"output_limit"`
	StackLimit    uint64 `json:"stack_limit"`
}

type Sample struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

func (submit *Submit) Load() (*worker.Worker, error) {
	if submit.TimeLimit == 0 {
		submit.TimeLimit = 1
	}
	if submit.RealTimeLimit == 0 {
		submit.RealTimeLimit = 1
	}
	if submit.MemoryLimit == 0 {
		submit.MemoryLimit = 256
	}
	if submit.OutputLimit == 0 {
		submit.OutputLimit = 256
	}
	if submit.StackLimit == 0 {
		submit.StackLimit = 256
	}

	if submit.RealTimeLimit < submit.TimeLimit {
		submit.RealTimeLimit = submit.TimeLimit + 2
	}

	if submit.StackLimit > submit.MemoryLimit {
		submit.StackLimit = submit.MemoryLimit
	}
	submit.CodeType = strings.ToLower(submit.CodeType)

	submitID := uuid.New()
	wk := &worker.Worker{
		WorkDir: filepath.Join(worker.RunDir, submitID.String()),

		SubmitID: submitID.String(),

		FileName:  strings.Join([]string{submitID.String(), submit.CodeType}, "."),
		Type:      submit.CodeType,
		AllowProc: submit.AllowProc,

		TimeLimit:     submit.TimeLimit,
		RealTimeLimit: submit.RealTimeLimit,
		MemoryLimit:   submit.MemoryLimit,
		OutputLimit:   submit.OutputLimit,
		StackLimit:    submit.StackLimit,
	}
	if err := submit.saveCodeAndSample(wk); err != nil {
		return nil, err
	}

	return wk, nil
}

func (submit *Submit) saveCodeAndSample(wk *worker.Worker) error {
	if err := os.MkdirAll(wk.WorkDir, dirPerm); err != nil {
		return err
	}

	if err := writeFile(filepath.Join(wk.WorkDir, wk.FileName), []byte(submit.SourceCode)); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(wk.WorkDir, "sample"), dirPerm); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i, sample := range submit.Samples {
		wg.Add(1)
		go func(id int, oneSample Sample) error {
			defer wg.Done()

			sampleID := strconv.FormatInt(int64(id), 10)
			sampleInPath := strings.Join([]string{sampleID, "in"}, ".")
			sampleOutPath := strings.Join([]string{sampleID, "out"}, ".")
			if err := writeFile(filepath.Join(wk.WorkDir, "sample", sampleInPath), []byte(oneSample.In)); err != nil {
				return err
			}
			if err := writeFile(filepath.Join(wk.WorkDir, "sample", sampleOutPath), []byte(oneSample.Out)); err != nil {
				return err
			}
			return nil
		}(i+1, sample)
	}
	wg.Wait()

	if err := os.MkdirAll(filepath.Join(wk.WorkDir, "output"), dirPerm); err != nil {
		return err
	}
	return nil
}

func writeFile(path string, text []byte) error {
	err := os.WriteFile(path, text, filePerm)
	return err
}

func remove(path string) error {
	if path != "" {
		return os.Remove(path)
	}
	return nil
}
