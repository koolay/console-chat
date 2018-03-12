// Package rethink
package rethink

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func checkRoomName(room string) error {
	if match, _ := regexp.MatchString(`[a-zA-Z]+[0-9\-]*[\-a-zA-Z]*`, room); !match {
		return errors.New("Only alphanumeric characters and underscores are valid!")
	}
	return nil
}

func checkUsername(username string) error {
	if match, _ := regexp.MatchString(`[a-zA-Z]+[0-9\-]*[\-a-zA-Z]*`, username); !match {
		return errors.New("Only alphanumeric characters and underscores are valid!")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
