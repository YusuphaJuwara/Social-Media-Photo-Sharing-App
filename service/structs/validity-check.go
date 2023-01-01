package structs

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func TokenCheck( r *http.Request) (string, error) {
	prefix := "Bearer "

	// If the header doesn't contain Authorization, it returns an empty string ""
	authHeader := r.Header.Get("Authorization")
	reqToken := strings.TrimPrefix(authHeader, prefix)

	// If the authHeader does not contain "Bearer ", then reqToken will be equal to authHeader ("Bearer " won't be trimmed off)
	// if authHeader == "" || reqToken == authHeader {
	// 	return nil, BadReqErr
	// }

	fmt.Printf("Token is correct: %v\n", reqToken)
	err := UuidCheck(reqToken)
	if err == nil{
		return reqToken, nil

	} else if errors.Is(err, BadReqErr) {
		return "", BadReqErr

	}

	return "", err

}

func UuidCheck(uid string) error {

	pattern := "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

	match, err := regexp.MatchString(pattern, uid)
	if err != nil {
		return err
	}

	if len(uid) != 36 || !match {
		return BadReqErr
	}

	return nil
}

func GenderCheck(str string) error {

	if str == "Female" || str == "Male" {
		return nil
	}

	return errors.New("invalid gender")
}

func PatternCheck(pattern, name string, min, max int) error {

	match, err := regexp.MatchString(pattern, name)
	if err != nil {
		return err
	}

	if len(name) < min || len(name) > max || !match {
		return BadReqErr
	}

	return nil
}

func DateTimeCheck(pattern, date string) error {

	match, err := regexp.MatchString(pattern, date)
	if err != nil {
		return err

	} else if !match {
		return BadReqErr

	}

	return nil
}
