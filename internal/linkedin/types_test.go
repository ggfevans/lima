package linkedin

import (
	"testing"
	"time"

	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"
	"go.mau.fi/util/jsontime"
)

// Helper to build a member participant with first and last name.
func makeMemberParticipant(urn string, first, last string) linkedingo.MessagingParticipant {
	return linkedingo.MessagingParticipant{
		EntityURN: linkedingo.NewURN(urn),
		ParticipantType: linkedingo.ParticipantType{
			Member: &linkedingo.MemberParticipantInfo{
				FirstName: linkedingo.AttributedText{Text: first},
				LastName:  linkedingo.AttributedText{Text: last},
			},
		},
	}
}

// Helper to build an organization participant.
func makeOrgParticipant(urn string, name string) linkedingo.MessagingParticipant {
	return linkedingo.MessagingParticipant{
		EntityURN: linkedingo.NewURN(urn),
		ParticipantType: linkedingo.ParticipantType{
			Organization: &linkedingo.OrganizationParticipantInfo{
				Name: linkedingo.AttributedText{Text: name},
			},
		},
	}
}

func TestConvertParticipant_SamePrefixSameID_IsOwnUser(t *testing.T) {
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")
	p := makeMemberParticipant("urn:li:fsd_profile:123", "Alice", "Smith")

	dp := ConvertParticipant(p, ownURN)

	if !dp.IsOwnUser {
		t.Error("expected IsOwnUser to be true when URNs have same prefix and same ID")
	}
}

func TestConvertParticipant_DifferentPrefixSameID_IsOwnUser(t *testing.T) {
	// This tests the bug fix: URNs with different prefixes but the same ID
	// should match. LinkedIn uses different prefixes in different contexts
	// (e.g. "urn:li:fsd_profile:123" vs "urn:li:member:123").
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")
	p := makeMemberParticipant("urn:li:member:123", "Alice", "Smith")

	dp := ConvertParticipant(p, ownURN)

	if !dp.IsOwnUser {
		t.Error("expected IsOwnUser to be true when URNs have different prefixes but same ID")
	}
}

func TestConvertParticipant_DifferentID_NotOwnUser(t *testing.T) {
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")
	p := makeMemberParticipant("urn:li:fsd_profile:456", "Bob", "Jones")

	dp := ConvertParticipant(p, ownURN)

	if dp.IsOwnUser {
		t.Error("expected IsOwnUser to be false when URN IDs differ")
	}
}

func TestConvertParticipant_MemberNameExtraction(t *testing.T) {
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:999")
	p := makeMemberParticipant("urn:li:fsd_profile:456", "Jane", "Doe")

	dp := ConvertParticipant(p, ownURN)

	expected := "Jane Doe"
	if dp.Name != expected {
		t.Errorf("expected Name %q, got %q", expected, dp.Name)
	}
}

func TestConvertMessage_OwnMessage(t *testing.T) {
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")
	msg := linkedingo.Message{
		EntityURN:   linkedingo.NewURN("urn:li:msg_message:(123,100)"),
		Body:        linkedingo.AttributedText{Text: "Hello"},
		DeliveredAt: jsontime.UnixMilli{Time: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)},
		Sender:      makeMemberParticipant("urn:li:fsd_profile:123", "Alice", "Smith"),
	}

	dm := ConvertMessage(msg, ownURN)

	if !dm.IsOwn {
		t.Error("expected IsOwn to be true when sender URN ID matches own URN ID")
	}
	if dm.Sender != "Alice Smith" {
		t.Errorf("expected Sender %q, got %q", "Alice Smith", dm.Sender)
	}
	if dm.Body != "Hello" {
		t.Errorf("expected Body %q, got %q", "Hello", dm.Body)
	}
}

func TestConvertMessage_OtherMessage(t *testing.T) {
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")
	msg := linkedingo.Message{
		EntityURN:   linkedingo.NewURN("urn:li:msg_message:(456,200)"),
		Body:        linkedingo.AttributedText{Text: "Hi there"},
		DeliveredAt: jsontime.UnixMilli{Time: time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)},
		Sender:      makeMemberParticipant("urn:li:fsd_profile:456", "Bob", "Jones"),
	}

	dm := ConvertMessage(msg, ownURN)

	if dm.IsOwn {
		t.Error("expected IsOwn to be false when sender URN ID differs from own URN ID")
	}
	if dm.Sender != "Bob Jones" {
		t.Errorf("expected Sender %q, got %q", "Bob Jones", dm.Sender)
	}
}

func TestConvertMessage_DifferentPrefixSameID_IsOwn(t *testing.T) {
	// Same bug-fix scenario as the participant test: different URN prefix,
	// same underlying ID should still be recognised as own message.
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")
	msg := linkedingo.Message{
		EntityURN:   linkedingo.NewURN("urn:li:msg_message:(123,300)"),
		Body:        linkedingo.AttributedText{Text: "Cross-prefix"},
		DeliveredAt: jsontime.UnixMilli{Time: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)},
		Sender:      makeMemberParticipant("urn:li:member:123", "Alice", "Smith"),
	}

	dm := ConvertMessage(msg, ownURN)

	if !dm.IsOwn {
		t.Error("expected IsOwn to be true when sender has different prefix but same ID as own URN")
	}
}

func TestConvertConversation_TitleFromNonOwnParticipant(t *testing.T) {
	ownURN := linkedingo.NewURN("urn:li:fsd_profile:123")

	conv := linkedingo.Conversation{
		EntityURN:      linkedingo.NewURN("urn:li:msg_conversation:100"),
		LastActivityAt: jsontime.UnixMilli{Time: time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)},
		Read:           true,
		ConversationParticipants: []linkedingo.MessagingParticipant{
			makeMemberParticipant("urn:li:fsd_profile:123", "Alice", "Smith"),
			makeMemberParticipant("urn:li:fsd_profile:456", "Bob", "Jones"),
		},
		// Title is intentionally left empty so it gets derived.
	}

	dc := ConvertConversation(conv, ownURN)

	expected := "Bob Jones"
	if dc.Title != expected {
		t.Errorf("expected Title %q (non-own participant), got %q", expected, dc.Title)
	}
}
