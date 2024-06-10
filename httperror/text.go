package httperror

import (
	"artk.dev/apperror"
	"artk.dev/assume"
	"io"
	"net/http"
	"strings"
)

// EncodeToText encodes an error into plain text.
// No further writes to the ResponseWriter w should happen after this function.
func EncodeToText(w http.ResponseWriter, err error) {
	assume.NotZero(w)

	kind := apperror.KindOf(err)
	status := EncodeKind(kind)

	switch kind {
	case apperror.OK:
		// The only difference compared to just returning is that
		// the content-type is set.
		http.Error(w, "", status)
	case apperror.UnknownError:
		http.Error(w, "Internal Server Error", status)
	default:
		http.Error(w, err.Error(), status)
	}
}

// DecodeFromText decodes an error from plain text.
func DecodeFromText(response *http.Response) error {
	assume.NotZero(response)

	kind := DecodeKind(response.StatusCode)
	if kind == apperror.OK {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return apperror.Unknownf("cannot parse HTTP error: %w", err)
	}

	// The http.Error used in encoding adds an extra newline. Remove it.
	msg := string(body)
	msg = strings.TrimSuffix(msg, "\n")

	return apperror.New(kind, msg)
}
