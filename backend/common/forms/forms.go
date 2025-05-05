package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var MNPhoneRX = regexp.MustCompile("^[7-9][0-9]{7}$")

// Create a custom Form struct, which anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	url.Values
	Errors errors
}

// Define a New function to initialize a custom Form struct. Notice that
// this takes the form data as the parameter?
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Implement a Required method to check that specific fields in the form
// data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Энэ талбарыг бөглөнө үү.")
		}
	}
}

// Implement a MaxLength method to check that a specific field in the form
// contains a maximum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MaxLength(field string, d int) {
	value := strings.TrimSpace(f.Get(field))
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Хэтэрхий урт байна. Хамгийн уртдаа %d үсэгтэй байж болно.", d))
	}
}

// Implement a PermittedValues method to check that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := strings.TrimSpace(f.Get(field))
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "Энд буруу мэдээлэл оруулсан байна.")
}

// Implement a MinLength method to check that a specific field in the form
// contains a minimum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MinLength(field string, d int) {
	value := strings.TrimSpace(f.Get(field))
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("Хэтэрхий богино байна. Хамгийн богинодоо %d үсэгтэй байж болно.", d))
	}
}

// Implement a MinValue method to check that a specific number field
// in the form contains below maximum value. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MaxValue(field string, d int) {
	value, err := strconv.Atoi(f.Get(field))
	if err != nil {
		return
	}
	if value > d {
		f.Errors.Add(field, fmt.Sprintf("Хэтэрхий их байна. Хамгийн ихдээ %d байж болно.", d))
	}
}

// Implement a MinValue method to check that a specific number field
// in the form contains below minimum value. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MinValue(field string, d int) {
	value, err := strconv.Atoi(f.Get(field))
	if err != nil {
		return
	}
	if value < d {
		f.Errors.Add(field, fmt.Sprintf("Хэтэрхий бага байна. Хамгийн багадаа %d байж болно.", d))
	}
}

func (f *Form) Number(fields ...string) {
	for _, field := range fields {
		value := strings.TrimSpace(f.Get(field))
		if value == "" {
			return
		}

		_, err := strconv.Atoi(value)
		if err != nil {
			f.Errors.Add(field, "Энд тоо оруулна уу")
		}
	}
}

func (f *Form) Date(fields ...string) {
	var format = "2006-01-02"
	for _, field := range fields {
		value := strings.TrimSpace(f.Get(field))
		if value == "" {
			return
		}

		_, err := time.Parse(format, value)
		if err != nil {
			f.Errors.Add(field, "Он сар өдрөө зөв оруулна уу.")
		}
	}
}

// Implement a MatchesPattern method to check that a specific field in the form
// matches a regular expression. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := strings.TrimSpace(f.Get(field))
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Энд буруу мэдээлэл оруулсан байна.")
	}
}

// Implement a Valid method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
