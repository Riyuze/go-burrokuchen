package utils

import (
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

func CatchErr(err error) error {
	pc, filename, line, _ := runtime.Caller(1)
	return errors.Wrap(err, runtime.FuncForPC(pc).Name()+";"+filename+";"+strconv.Itoa(line))
}
