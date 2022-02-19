package utils

import "log"

func Handle(err error) error {
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}