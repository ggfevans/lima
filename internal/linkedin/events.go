package linkedin

import (
	"context"

	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"
)

// onDecoratedEvent handles real-time decorated events from LinkedIn SSE.
func (c *Client) onDecoratedEvent(ctx context.Context, evt *linkedingo.DecoratedEvent) {
	if c.program == nil {
		return
	}

	data := evt.Payload.Data

	switch {
	case data.DecoratedMessage != nil:
		msg := data.DecoratedMessage.Result
		dm := ConvertMessage(msg, c.ownURN)
		convID := msg.BackendConversationURN.String()
		if convID == "" {
			convID = msg.Conversation.EntityURN.String()
		}
		c.program.Send(RealtimeMessageMsg{
			ConversationID: convID,
			Message:        dm,
		})

	case data.DecoratedTypingIndicator != nil:
		ti := data.DecoratedTypingIndicator.Result
		name := ""
		if ti.TypingParticipant.ParticipantType.Member != nil {
			name = ti.TypingParticipant.ParticipantType.Member.FirstName.Text
		}
		convID := ti.Conversation.EntityURN.String()
		c.program.Send(RealtimeTypingMsg{
			ConversationID: convID,
			SenderName:     name,
		})

	case data.DecoratedSeenReceipt != nil:
		receipt := data.DecoratedSeenReceipt.Result
		convID := receipt.Message.Conversation.EntityURN.String()
		if convID == "" {
			convID = receipt.Message.BackendConversationURN.String()
		}
		c.program.Send(RealtimeSeenMsg{
			ConversationID: convID,
		})
	}
}

// onBadCredentials handles credential expiry.
func (c *Client) onBadCredentials(_ context.Context, _ error) {
	if c.program == nil {
		return
	}
	c.program.Send(SessionExpiredMsg{})
}

// onTransientDisconnect handles temporary disconnects.
func (c *Client) onTransientDisconnect(_ context.Context, err error) {
	if c.program == nil {
		return
	}
	c.program.Send(RealtimeDisconnectedMsg{Err: err})
}

// onClientConnection handles successful connection events.
func (c *Client) onClientConnection(_ context.Context, _ *linkedingo.ClientConnection) {
	if c.program == nil {
		return
	}
	c.program.Send(RealtimeConnectedMsg{})
}
