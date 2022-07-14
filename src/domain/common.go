package domain

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/cjlapao/http-loadtester-go/common"
)

// Interval Entity
type Interval struct {
	value int
}

// Value Gets an interval value
func (s Interval) Value() int {
	return s.value
}

// NewInterval Creates a new interval value
func NewInterval(value int) Interval {
	return Interval{value: value}
}

// ResponseDetails Entity
type ResponseDetails struct {
	IP            string
	TLSCipher     string
	TLSVersion    string
	TLSServerName string
	Body          string
}

func GetRandomBlockInterval(maxInterval Interval, minInterval Interval) int {
	max := maxInterval.Value()
	min := minInterval.Value()

	return common.GetRandomNumber(min, max)
}

func GenerateFormUrlEncoded(values map[string]string) *strings.Reader {
	form := url.Values{}
	for key, val := range values {
		form.Add(key, val)
	}

	return strings.NewReader(form.Encode())
}

func GenerateFormData(values map[string]string) (string, *bytes.Reader) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range values {
		field, err := writer.CreateFormField(key)
		if err != nil {
			return "", nil
		}

		_, err = io.Copy(field, strings.NewReader(val))
		if err != nil {
			return "", nil
		}
	}

	contentType := writer.FormDataContentType()
	writer.Close()

	return contentType, bytes.NewReader(body.Bytes())
}
