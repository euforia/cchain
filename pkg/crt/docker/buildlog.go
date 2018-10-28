package docker

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// BuildStep ...
type BuildStep struct {
	Index     int
	Cmd       string
	UsedCache bool
	RanIn     string
	Digest    string
	Logs      []string
}

// BuildLogLine is a single build log line
type BuildLogLine struct {
	Stream string
}

// StructedBuildLogReader ...
type StructedBuildLogReader struct {
	r     io.ReadCloser
	w     io.Writer
	Steps []*BuildStep
}

// NewStructedBuildLogReader ... It takes the build output reader from docker and
// a writer to write out the log line if needed.
func NewStructedBuildLogReader(r io.ReadCloser, output io.Writer) *StructedBuildLogReader {
	return &StructedBuildLogReader{
		r:     r,
		w:     output,
		Steps: make([]*BuildStep, 0),
	}
}

func (r *StructedBuildLogReader) Read() error {
	defer r.r.Close()

	var (
		dec  = json.NewDecoder(r.r)
		step *BuildStep
		err  error
	)

	for {

		var bll BuildLogLine
		if err = dec.Decode(&bll); err != nil {
			if err == io.EOF {
				r.Steps = append(r.Steps, step)
				err = nil
			}
			break
		}

		// Write to provided writer
		r.w.Write([]byte(bll.Stream))

		// Start parsing
		line := strings.TrimSuffix(bll.Stream, "\n")
		if line == "" {
			continue
		}

		if len(line) < 5 {
			step.Logs = append(step.Logs, line)
			continue
		}

		switch line[:5] {
		case "Step ":
			if step != nil {
				r.Steps = append(r.Steps, step)
			}
			step = parseStepLine(line[5:])

		case " --->":
			str := line[6:]

			switch {
			case strings.HasPrefix(str, "Using cache"):
				step.UsedCache = true

			case strings.HasPrefix(str, "Running in"):
				step.RanIn = strings.TrimPrefix(str[:len(str)-1], "Running in ")

			default:
				if _, err := hex.DecodeString(str); err == nil {
					step.Digest = str
				}
			}

		default:
			step.Logs = append(step.Logs, line)

		}

	}

	return err
}

func parseStepLine(str string) *BuildStep {
	i := strings.IndexRune(str, '/')
	step, _ := strconv.ParseInt(str[:i], 10, 64)

	i = strings.IndexRune(str, ':')
	cmd := str[i+2:]

	return &BuildStep{
		Index: int(step),
		Cmd:   cmd,
	}
}
