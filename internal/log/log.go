// Copyright 2023 The fhub-runtime-go Authors
// This file is part of fhub-runtime-go.
//
// This file is part of fhub-runtime-go.
// fhub-runtime-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// fhub-runtime-go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with fhub-runtime-go. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

var (
	root zerolog.Logger
)

type FmtOutput struct{}

// Support webassembly
func (f *FmtOutput) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: &FmtOutput{}, NoColor: false, TimeFormat: time.RFC3339}
	// output := zerolog.ConsoleWriter{Out: &FmtOutput{}, TimeFormat: time.RFC3339}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	root = zerolog.New(output).With().Timestamp().Logger()
}

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
}

type Log struct {
	logger zerolog.Logger
}

func (log *Log) Trace(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Trace(), msg, args)
}

func (log *Log) Debug(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Debug(), msg, args)
}

func (log *Log) Info(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Info(), msg, args)
}

func (log *Log) Warn(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Warn(), msg, args)
}

func (log *Log) Error(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Error(), msg, args)
}

func (log *Log) Fatal(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Fatal(), msg, args)
}

func (log *Log) Panic(msg string, args ...interface{}) {
	log.iterateLog(log.logger.Panic(), msg, args)
}

func (log *Log) iterateLog(loggerLevel *zerolog.Event, msg string, args []interface{}) {
	for i := 0; i < len(args); i += 2 {
		if len(args) <= i+1 {
			loggerLevel.Interface(args[i].(string), "<empty>")
		} else {
			loggerLevel.Interface(args[i].(string), args[i+1])
		}
	}

	loggerLevel.Msg(msg)
}

func New(context string) Logger {
	return &Log{
		logger: root.With().Str("context", context).Logger(),
	}
}
