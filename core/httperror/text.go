package httperror

import (
	"artk.dev/core/apperror"
	"artk.dev/core/mustbe"
	"io"
	"net/http"
	"strings"
)

// EncodeToText encodes an error into plain text.
// No further writes to the ResponseWriter w should happen after this function.
func EncodeToText(w http.ResponseWriter, err error) {
	mustbe.NotNil(w)

	kind := apperror.KindOf(err)
	status := EncodeKind(kind)

	var msg string
	if err != nil {
		msg = err.Error()
	}

	http.Error(w, msg, status)
}

// DecodeFromText decodes an error from plain text.
func DecodeFromText(response *http.Response) error {
	mustbe.NotNil(response)

	kind := DecodeKind(response.StatusCode)
	if kind == apperror.OK {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return apperror.Unknown("cannot parse HTTP error: %w", err)
	}

	// The http.Error used in encoding adds an extra newline. Remove it.
	msg := string(body)
	msg = strings.TrimSuffix(msg, "\n")

	return apperror.New(kind, msg)
}
