// Copyright (c) 2019 Australian Rivers Institute.

package json

import (
	"bufio"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"os"
	"path"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/pkg/errors"
)

const fileType = "json"
const fileTypeExtension = "." + fileType

type Encoder struct {
	loggers.ContainedLogger
	marshaler  Marshaler
	outputPath string
}

func (e *Encoder) WithOutputPath(outputPath string) *Encoder {
	e.outputPath = outputPath
	return e
}

func (e *Encoder) WithLogHandler(logHandler logging.Logger) *Encoder {
	e.SetLogHandler(logHandler)
	return e
}

func (e Encoder) Encode(solution *solution.Solution) error {
	e.LogHandler().Info("Saving [" + solution.Id + "] as [JSON]")
	marshaledSolution, marshalError := e.marshaler.Marshal(solution)
	if marshalError != nil {
		return errors.Wrap(marshalError, fileType+"marshaling of solution")
	}

	outputPath := e.deriveOutputPath(solution)
	e.LogHandler().Debug("Encoding [" + solution.Id + "] to [" + outputPath + "]")
	return e.encodeMarshaled(marshaledSolution, outputPath)
}

func (e Encoder) encodeMarshaled(marshaledSolution []byte, outputPath string) error {
	file, openError := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if openError != nil {
		return errors.Wrap(openError, "opening file for "+fileType+" encoding of solution")
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, writeError := bufferedWriter.Write(marshaledSolution); writeError != nil {
		return errors.Wrap(writeError, "writing marshaled "+fileType+" of solution")
	}

	bufferedWriter.Flush()
	return nil
}

func (e Encoder) deriveOutputPath(solution *solution.Solution) (outputPath string) {
	safeIdBasedFileName := solution.FileNameSafeId() + fileTypeExtension
	return path.Join(e.outputPath, safeIdBasedFileName)
}
