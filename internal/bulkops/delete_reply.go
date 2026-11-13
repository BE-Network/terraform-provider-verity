package bulkops

import (
	"regexp"
	"strings"
)

type DeleteBlockCategory string

const (
	DeleteBlockedPredefined DeleteBlockCategory = "predefined"
	DeleteBlockedReferenced DeleteBlockCategory = "referenced"
	DeleteBlockedInUse      DeleteBlockCategory = "in_use"
	DeleteBlockedUnresolved DeleteBlockCategory = "unresolved_reference"
)

type DeleteBlock struct {
	Category DeleteBlockCategory
	Reason   string
}

var (
	deleteReplyHeader     = regexp.MustCompile(`^Object\(s\) could not be deleted as they are (.+):$`)
	deleteReplyEntry      = regexp.MustCompile(`^([A-Za-z0-9_]+)\|([^|]*)\|(.*)$`)
	deleteReplyPredefined = regexp.MustCompile(`^Cannot delete ([A-Za-z0-9_]+) "([^"]*)"\s+it is a predefined object$`)
	deleteReplyUnresolved = regexp.MustCompile(`^Cannot resolve object reference: ([A-Za-z0-9_]+) "([^"]*)" is missing from the database or import payload, but is referenced by ([A-Za-z0-9_]+) "([^"]*)"$`)
)

var deleteReplyCategories = map[string]DeleteBlockCategory{
	"predefined objects":          DeleteBlockedPredefined,
	"referenced by other objects": DeleteBlockedReferenced,
	"being used by users":         DeleteBlockedInUse,
}

var entryReasonPrefix = map[DeleteBlockCategory]string{
	DeleteBlockedReferenced: " could not be deleted as it was referenced by ",
	DeleteBlockedInUse:      " could not be deleted as it is being viewed by users ",
}

func ParseDeleteReply(message string, requested []string) (map[string]DeleteBlock, bool) {
	inRequest := make(map[string]bool, len(requested))
	for _, name := range requested {
		inRequest[name] = true
	}
	blocks := make(map[string]DeleteBlock)
	var category DeleteBlockCategory
	entries := 0
	record := func(name string, block DeleteBlock) bool {
		if !inRequest[name] {
			return false
		}
		if existing, seen := blocks[name]; seen && existing.Category != block.Category {
			return false
		}
		blocks[name] = block
		return true
	}

	for _, raw := range strings.Split(message, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if match := deleteReplyHeader.FindStringSubmatch(line); match != nil {
			known, ok := deleteReplyCategories[match[1]]
			if !ok || (category != "" && entries == 0) {
				return nil, false
			}
			category = known
			entries = 0
			continue
		}
		if match := deleteReplyPredefined.FindStringSubmatch(line); match != nil {
			if !record(match[2], DeleteBlock{Category: DeleteBlockedPredefined, Reason: "it is a predefined object"}) {
				return nil, false
			}
			continue
		}
		if match := deleteReplyUnresolved.FindStringSubmatch(line); match != nil {
			reason := "deleting it would leave " + match[3] + " \"" + match[4] + "\" with a reference to a missing object"
			if !record(match[2], DeleteBlock{Category: DeleteBlockedUnresolved, Reason: reason}) {
				return nil, false
			}
			continue
		}
		match := deleteReplyEntry.FindStringSubmatch(line)
		if match == nil || category == "" {
			return nil, false
		}
		name, rest := match[2], match[3]
		switch category {
		case DeleteBlockedPredefined:
			if strings.TrimSpace(rest) != "" {
				return nil, false
			}
			rest = "it is a predefined object"
		default:
			prefix := entryReasonPrefix[category]
			if !strings.HasPrefix(rest, prefix) {
				return nil, false
			}
			rest = strings.TrimSpace(rest)
		}
		if !record(name, DeleteBlock{Category: category, Reason: rest}) {
			return nil, false
		}
		entries++
	}
	if category != "" && entries == 0 {
		return nil, false
	}
	if len(blocks) == 0 {
		return nil, false
	}
	return blocks, true
}
