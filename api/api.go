package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ResponseMeta struct {
	Header     http.Header
	Status     string
	StatusCode int
}

type Response struct {
	ResponseMeta ResponseMeta
}

func (resp *Response) SetMeta(meta ResponseMeta) {
	resp.ResponseMeta = meta
}

type MetaSetter interface {
	SetMeta(ResponseMeta)
}

// base64DataRegex is the regular expression for detecting base64 encoded strings.
var base64DataRegex = regexp.MustCompile("^data:([\\w-]+/[\\w\\-+.]+)?(;[\\w-]+=[\\w-]+)*;base64,([a-zA-Z0-9/+\\n=]+)$")

// StructToParams serializes struct to url.Values, which can be further sent to the http client.
func StructToParams(inputStruct interface{}) (url.Values, error) {
	var paramsMap map[string]interface{}
	paramsJSONObj, _ := json.Marshal(inputStruct)
	err := json.Unmarshal(paramsJSONObj, &paramsMap)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for paramName, value := range paramsMap {
		kind := reflect.ValueOf(value).Kind()

		if kind == reflect.Slice || kind == reflect.Array {
			rVal := reflect.ValueOf(value)
			for i := 0; i < rVal.Len(); i++ {
				item := rVal.Index(i)
				val, err := encodeParamValue(item.Interface())
				if err != nil {
					return nil, err
				}

				arrParamName := fmt.Sprintf("%s[%d]", paramName, i)
				params.Add(arrParamName, val)
			}

			continue
		}

		val, err := encodeParamValue(value)
		if err != nil {
			return nil, err
		}

		params.Add(paramName, val)
	}

	return params, nil
}

func encodeParamValue(value interface{}) (string, error) {
	resBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	res := string(resBytes)
	if strings.HasPrefix(res, "\"") { // FIXME: Fix this dirty hack that prevents double quoting of strings
		res, _ = strconv.Unquote(res)
	}

	return res, nil
}

// BuildPath builds (joins) the URL path from the provided parts.
func BuildPath(parts ...interface{}) string {
	var partsSlice []string

	for _, part := range parts {
		partRes := ""
		switch partVal := part.(type) {
		case string:
			partRes = partVal
		case fmt.Stringer:
			partRes = partVal.String()
		default:
			partRes = fmt.Sprintf("%v", partVal)
		}
		if len(partRes) > 0 {
			partsSlice = append(partsSlice, partRes)
		}
	}

	return strings.Join(partsSlice, "/")
}

// DeferredClose is a wrapper around io.Closer.Close method.
func DeferredClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

// IsLocalFilePath determines whether the provided path can be a local file.
func IsLocalFilePath(path interface{}) bool {
	switch pathV := path.(type) {
	case string:
		return !(IsValidURL(pathV) || IsBase64Data(pathV))
	default:
		return false
	}
}

// IsValidURL checks whether urlCandidate string is a valid URL.
func IsValidURL(urlCandidate string) bool {
	urlStruct, err := url.Parse(urlCandidate)
	if err != nil || urlStruct.Scheme == "" {
		return false
	}

	return true
}

// IsBase64Data checks whether base64Candidate represents a valid base64 encoded string.
func IsBase64Data(base64Candidate string) bool {
	return base64DataRegex.MatchString(base64Candidate)
}

func SetResponseMeta(httpResp *http.Response, respStruct MetaSetter) {
	if httpResp == nil {
		return
	}

	meta := ResponseMeta{
		Header:     httpResp.Header,
		Status:     httpResp.Status,
		StatusCode: httpResp.StatusCode,
	}
	respStruct.SetMeta(meta)
}