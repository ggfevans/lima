package linkedin

import (
	"time"

	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"

	"github.com/ggfevans/lima/internal/util"
)

// DisplayConversation is a display-friendly version of a LinkedIn conversation.
type DisplayConversation struct {
	ID             string
	Title          string
	LastMessage    string
	LastActivityAt time.Time
	Unread         bool
	Participants   []DisplayParticipant
	URN            linkedingo.URN
}

// DisplayParticipant is a display-friendly participant.
type DisplayParticipant struct {
	Name      string
	Headline  string
	URN       linkedingo.URN
	IsOwnUser bool
}

// DisplayMessage is a display-friendly message.
type DisplayMessage struct {
	ID         string
	Sender     string
	SenderURN  linkedingo.URN
	Body       string
	Timestamp  time.Time
	IsOwn      bool
	Format     linkedingo.MessageBodyRenderFormat
	MessageURN linkedingo.URN
}

// ConvertConversation converts a linkedingo Conversation to a display type.
func ConvertConversation(conv linkedingo.Conversation, ownURN linkedingo.URN) DisplayConversation {
	dc := DisplayConversation{
		ID:             conv.EntityURN.String(),
		Title:          conv.Title,
		LastActivityAt: conv.LastActivityAt.Time,
		Unread:         !conv.Read,
		URN:            conv.EntityURN,
	}

	for _, p := range conv.ConversationParticipants {
		dp := ConvertParticipant(p, ownURN)
		dc.Participants = append(dc.Participants, dp)
	}

	// Derive title from participant names if empty
	if dc.Title == "" {
		for _, p := range dc.Participants {
			if !p.IsOwnUser && p.Name != "" {
				dc.Title = p.Name
				break
			}
		}
	}
	if dc.Title == "" {
		dc.Title = "Conversation"
	}

	// Extract last message preview
	if len(conv.Messages.Elements) > 0 {
		last := conv.Messages.Elements[len(conv.Messages.Elements)-1]
		dc.LastMessage = last.Body.Text
	}

	return dc
}

// ConvertParticipant converts a linkedingo MessagingParticipant to display type.
func ConvertParticipant(p linkedingo.MessagingParticipant, ownURN linkedingo.URN) DisplayParticipant {
	dp := DisplayParticipant{
		URN:       p.EntityURN,
		IsOwnUser: p.EntityURN.ID() == ownURN.ID(),
	}

	if p.ParticipantType.Member != nil {
		dp.Name = p.ParticipantType.Member.FirstName.Text + " " + p.ParticipantType.Member.LastName.Text
		dp.Headline = p.ParticipantType.Member.Headline.Text
	} else if p.ParticipantType.Organization != nil {
		dp.Name = p.ParticipantType.Organization.Name.Text
	}

	return dp
}

// ConvertMessage converts a linkedingo Message to a display type.
func ConvertMessage(msg linkedingo.Message, ownURN linkedingo.URN) DisplayMessage {
	dm := DisplayMessage{
		ID:         msg.EntityURN.String(),
		Body:       msg.Body.Text,
		Timestamp:  msg.DeliveredAt.Time,
		Format:     msg.MessageBodyRenderFormat,
		SenderURN:  msg.Sender.EntityURN,
		MessageURN: msg.EntityURN,
		IsOwn:      msg.Sender.EntityURN.ID() == ownURN.ID(),
	}

	if msg.Sender.ParticipantType.Member != nil {
		dm.Sender = msg.Sender.ParticipantType.Member.FirstName.Text + " " + msg.Sender.ParticipantType.Member.LastName.Text
	} else if msg.Sender.ParticipantType.Organization != nil {
		dm.Sender = msg.Sender.ParticipantType.Organization.Name.Text
	} else {
		dm.Sender = "Unknown"
	}

	return dm
}

// ToConvListItem converts a DisplayConversation to a convlist.Conversation.
func (dc DisplayConversation) ToConvListItem() (id, name, lastMsg, timestamp string, unread bool) {
	return dc.ID, dc.Title, dc.LastMessage, util.RelativeTime(dc.LastActivityAt), dc.Unread
}

// ToThreadMessage converts a DisplayMessage to a thread.Message.
func (dm DisplayMessage) ToThreadMessage() (id, sender, body, timestamp string, isOwn bool) {
	return dm.ID, dm.Sender, dm.Body, util.RelativeTime(dm.Timestamp), dm.IsOwn
}
