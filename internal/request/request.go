package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestText := string(rawBytes)
	requestLine, err := parseRequestLine(requestText)
	if err != nil {
		return nil, fmt.Errorf("error parsing request-line: %w", err)
	}

	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(requestText string) (*RequestLine, error) {
	requestParts := strings.Split(requestText, "\r\n")
	requestLineParts := strings.Split(requestParts[0], " ")
	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request line format")
	}

	// Validate if Method is uppercase
	method := requestLineParts[0]
	if method != strings.ToUpper(method) {
		return nil, fmt.Errorf("malformed request-line method: %s", method)
	}

	// Validate if the HttpVersion is HTTP/1.1
	httpVersion := requestLineParts[2]
	httpVersionParts := strings.Split(httpVersion, "/")

	httpVersionProtocol := httpVersionParts[0]
	if httpVersionProtocol != "HTTP" {
		return nil, fmt.Errorf("malformed request-line http version: %s", httpVersionProtocol)
	}

	httpVersionNumber := httpVersionParts[1]
	if httpVersionNumber != "1.1" {
		return nil, fmt.Errorf("invalid request-line http version number: %s", httpVersionNumber)
	}

	requestLine := &RequestLine{
		Method:        method,
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersionNumber,
	}

	return requestLine, nil
}
