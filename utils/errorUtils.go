package utils

import (
	"log"
	"runtime"
)

// https://stackoverflow.com/questions/24809287/how-do-you-get-a-golang-program-to-print-the-line-number-of-the-error-it-just-ca
func Handle(err error) error {
	if err != nil {
		pc, filename, line, _ := runtime.Caller(1)
		log.Printf("[errors] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)

		return err
	}
	return nil
}