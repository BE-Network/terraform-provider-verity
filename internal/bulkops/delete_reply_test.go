package bulkops

import (
	"errors"
	"reflect"
	"testing"

	"terraform-provider-verity/internal/utils"
)

const (
	replyPredefined       = `{"status": "failed", "message": "Object(s) could not be deleted as they are predefined objects:</br>eth_device_profile|(Device Settings)|"}`
	replyReferenced       = `{"status": "failed", "message": "Object(s) could not be deleted as they are referenced by other objects:</br>threshold|Threshold A| could not be deleted as it was referenced by threshold_group|Group A|</br>threshold|Threshold B| could not be deleted as it was referenced by threshold_group|Group B|</br>"}`
	replyInUse            = `{"status": "failed", "message": "Object(s) could not be deleted as they are being used by users:</br>archipelago|Fabric A| could not be deleted as it is being viewed by users user_state|User State 8| and user_state|User State 18|"}`
	replyUnresolved       = `{"status": "failed", "message": "Cannot resolve object reference: badges \"Badge A\" is missing from the database or import payload, but is referenced by system_sections \"System Sections BlackBox\""}`
	replyPredefinedSingle = `{"status": "failed", "message": "Cannot delete packet_queue \"(Packet Queue)\"  it is a predefined object"}`
	replyMixed            = `{"status": "failed", "message": "Object(s) could not be deleted as they are predefined objects:</br>packet_queue|(Packet Queue)|</br>Object(s) could not be deleted as they are referenced by other objects:</br>packet_queue|Queue A| could not be deleted as it was referenced by eth_device_profile|Settings A|"}`
)

type replyError struct{ body string }

func (e replyError) Error() string { return "400 Bad Request" }
func (e replyError) Body() []byte  { return []byte(e.body) }

func reason(body string) string { return utils.APIErrorMessage(replyError{body}) }

func TestParseDeleteReplyAcceptsKnownShapes(t *testing.T) {
	for name, tc := range map[string]struct {
		body      string
		requested []string
		want      map[string]DeleteBlock
	}{
		"predefined": {
			body:      replyPredefined,
			requested: []string{"Save30minutes", "(Device Settings)"},
			want:      map[string]DeleteBlock{"(Device Settings)": {Category: DeleteBlockedPredefined, Reason: "it is a predefined object"}},
		},
		"referenced, referencing objects are not blocked names": {
			body:      replyReferenced,
			requested: []string{"Threshold A", "Threshold B", "Threshold C"},
			want: map[string]DeleteBlock{
				"Threshold A": {Category: DeleteBlockedReferenced, Reason: "could not be deleted as it was referenced by threshold_group|Group A|"},
				"Threshold B": {Category: DeleteBlockedReferenced, Reason: "could not be deleted as it was referenced by threshold_group|Group B|"},
			},
		},
		"in use by GUI sessions": {
			body:      replyInUse,
			requested: []string{"Fabric A"},
			want: map[string]DeleteBlock{"Fabric A": {Category: DeleteBlockedInUse,
				Reason: "could not be deleted as it is being viewed by users user_state|User State 8| and user_state|User State 18|"}},
		},
		"unresolved reference from an unmanaged object": {
			body:      replyUnresolved,
			requested: []string{"Badge A", "Badge B"},
			want: map[string]DeleteBlock{"Badge A": {Category: DeleteBlockedUnresolved,
				Reason: `deleting it would leave system_sections "System Sections BlackBox" with a reference to a missing object`}},
		},
		"predefined, single-line format": {
			body:      replyPredefinedSingle,
			requested: []string{"(Packet Queue)"},
			want:      map[string]DeleteBlock{"(Packet Queue)": {Category: DeleteBlockedPredefined, Reason: "it is a predefined object"}},
		},
		"several categories in one reply": {
			body:      replyMixed,
			requested: []string{"(Packet Queue)", "Queue A", "Queue B"},
			want: map[string]DeleteBlock{
				"(Packet Queue)": {Category: DeleteBlockedPredefined, Reason: "it is a predefined object"},
				"Queue A":        {Category: DeleteBlockedReferenced, Reason: "could not be deleted as it was referenced by eth_device_profile|Settings A|"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := ParseDeleteReply(reason(tc.body), tc.requested)
			if !ok {
				t.Fatalf("reply was rejected: %q", reason(tc.body))
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %#v\nwant %#v", got, tc.want)
			}
		})
	}
}

func TestParseDeleteReplyRejectsAnythingNotFullyUnderstood(t *testing.T) {
	for name, tc := range map[string]struct {
		message   string
		requested []string
	}{
		"unknown category": {
			message:   "Object(s) could not be deleted as they are locked by an upgrade:\nthreshold|Threshold A|",
			requested: []string{"Threshold A"},
		},
		"blocked name not in the request": {
			message:   reason(replyReferenced),
			requested: []string{"Threshold A"},
		},
		"referencing object named only, not the blocked one": {
			message:   "Object(s) could not be deleted as they are referenced by other objects:\nthreshold_group|Group A|",
			requested: []string{"Group A"},
		},
		"truncated entry": {
			message:   "Object(s) could not be deleted as they are referenced by other objects:\nthreshold|Threshold A| could not be",
			requested: []string{"Threshold A"},
		},
		"entry without a header": {
			message:   "threshold|Threshold A| could not be deleted as it was referenced by threshold_group|Group A|",
			requested: []string{"Threshold A"},
		},
		"predefined entry with trailing text": {
			message:   "Object(s) could not be deleted as they are predefined objects:\npacket_queue|(Packet Queue)| because reasons",
			requested: []string{"(Packet Queue)"},
		},
		"unrelated error": {
			message:   "Failed to decrypt some encrypted fields.",
			requested: []string{"Threshold A"},
		},
		"empty": {message: "", requested: []string{"Threshold A"}},
		"single-line predefined, name not in the request": {
			message:   reason(replyPredefinedSingle),
			requested: []string{"Queue A"},
		},
		"single-line predefined with trailing text": {
			message:   `Cannot delete packet_queue "(Packet Queue)"  it is a predefined object, sorry`,
			requested: []string{"(Packet Queue)"},
		},
		"single-line with a different verb": {
			message:   `Cannot update packet_queue "(Packet Queue)"  it is a predefined object`,
			requested: []string{"(Packet Queue)"},
		},
		"category header with no entries after a valid one (truncated)": {
			message:   "Object(s) could not be deleted as they are predefined objects:\npacket_queue|(Packet Queue)|\nObject(s) could not be deleted as they are referenced by other objects:",
			requested: []string{"(Packet Queue)", "Queue A"},
		},
		"category header with no entries before another header": {
			message:   "Object(s) could not be deleted as they are referenced by other objects:\nObject(s) could not be deleted as they are predefined objects:\npacket_queue|(Packet Queue)|",
			requested: []string{"(Packet Queue)"},
		},
		"same name in two categories": {
			message:   "Object(s) could not be deleted as they are predefined objects:\nqueue|Q|\nObject(s) could not be deleted as they are referenced by other objects:\nqueue|Q| could not be deleted as it was referenced by x|Y|",
			requested: []string{"Q"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got, ok := ParseDeleteReply(tc.message, tc.requested); ok {
				t.Fatalf("reply should be rejected, got %#v", got)
			}
		})
	}
}

func TestParseDeleteReplyRejectsErrorsWithoutBody(t *testing.T) {
	if _, ok := ParseDeleteReply(utils.APIErrorMessage(errors.New("connection reset")), []string{"x"}); ok {
		t.Fatal("an error without a reply body must not be parsed")
	}
}
