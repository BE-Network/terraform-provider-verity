package utils

import (
	"errors"
	"strings"
	"testing"
)

type apiError struct {
	status string
	body   string
}

func (e apiError) Error() string { return e.status }
func (e apiError) Body() []byte  { return []byte(e.body) }

const (
	deleteReply  = `{"status": "failed", "message": "Object(s) could not be deleted as they are predefined objects:</br>packet_queue|(Packet Queue)|"}`
	putReply     = `{"Error":"Failed to decrypt some encrypted fields."}`
	payloadReply = `{"payload":"something went wrong"}`
)

func TestAPIErrorMessageReadsEveryReplyShape(t *testing.T) {
	for name, tc := range map[string]struct{ body, want string }{
		"DELETE message": {deleteReply, "Object(s) could not be deleted as they are predefined objects:\npacket_queue|(Packet Queue)|"},
		"PUT Error":      {putReply, "Failed to decrypt some encrypted fields."},
		"payload":        {payloadReply, "something went wrong"},
		"no reason":      {`{"status":"failed"}`, ""},
		"not JSON":       {`<html>bad gateway</html>`, ""},
	} {
		t.Run(name, func(t *testing.T) {
			if got := APIErrorMessage(apiError{"400 Bad Request", tc.body}); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
	if got := APIErrorMessage(errors.New("plain")); got != "" {
		t.Fatalf("an error without a body gave %q", got)
	}
}

func TestFormatOpenAPIErrorShowsTheReason(t *testing.T) {
	diags := FormatOpenAPIError(apiError{"400 Bad Request", putReply}, "Failed to create switchpoint leaf1")
	if !diags.HasError() {
		t.Fatal("expected an error")
	}
	if detail := diags[0].Detail(); !strings.Contains(detail, "Failed to decrypt some encrypted fields.") || !strings.Contains(detail, "400 Bad Request") {
		t.Fatalf("detail does not carry the status and the API reason: %q", detail)
	}
}

func TestFormatOpenAPIErrorTruncatesLongReasons(t *testing.T) {
	var lines []string
	for i := 0; i < 50; i++ {
		lines = append(lines, "threshold|T| could not be deleted as it was referenced by threshold_group|G|")
	}
	body := `{"message": "Object(s) could not be deleted as they are referenced by other objects:</br>` + strings.Join(lines, "</br>") + `"}`
	detail := FormatOpenAPIError(apiError{"400 Bad Request", body}, "Failed to delete threshold T")[0].Detail()
	if !strings.Contains(detail, "… and 31 more lines") {
		t.Fatalf("long reason was not truncated as expected: %q", detail[len(detail)-80:])
	}
}
