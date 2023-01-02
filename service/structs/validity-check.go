package structs

import (
	"errors"
	_"fmt"
	"net/http"
	"regexp"
	"strings"
)

func TokenCheck( r *http.Request) (string, error) {
	prefix := "Bearer "

	// If the header doesn't contain Authorization, it returns an empty string ""
	authHeader := r.Header.Get("Authorization")
	reqToken := strings.TrimPrefix(authHeader, "\"")
	reqToken = strings.TrimPrefix(reqToken, prefix)
	reqToken = strings.TrimSuffix(reqToken, "\"")


	// If the authHeader does not contain "Bearer ", then reqToken will be equal to authHeader ("Bearer " won't be trimmed off)
	// if authHeader == "" || reqToken == authHeader {
	// 	return nil, BadReqErr
	// }

	err := UuidCheck(reqToken)
	
	if err == nil{
		return reqToken, nil

	} else if errors.Is(err, BadReqErr) {
		return "", BadReqErr

	}

	return "", err

}

func UuidCheck(uid string) error {


	uid = strings.TrimPrefix(uid, "\"")
	uid = strings.TrimSuffix(uid, "\"")


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

	str = strings.TrimPrefix(str, "\"")
	str = strings.TrimSuffix(str, "\"")

	if str == "Female" || str == "Male" {
		return nil
	}

	return errors.New("invalid gender")
}

func PatternCheck(pattern, name string, min, max int) error {

	name = strings.TrimPrefix(name, "\"")
	name = strings.TrimSuffix(name, "\"")

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

	date = strings.TrimPrefix(date, "\"")
	date = strings.TrimSuffix(date, "\"")

	match, err := regexp.MatchString(pattern, date)
	if err != nil {
		return err

	} else if !match {
		return BadReqErr

	}

	return nil
}
