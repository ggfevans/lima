package linkedin

import (
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"
)

type scheduledEvent struct {
	delay          time.Duration
	conversationID string
	message        DisplayMessage
}

// DemoClient implements MessagingClient with canned data and scripted events.
type DemoClient struct {
	mu            sync.Mutex
	program       *tea.Program
	ownURN        linkedingo.URN
	conversations []DisplayConversation
	messages      map[string][]DisplayMessage
	autoReplies   map[string][]string
	replyIndex    map[string]int
	timeline      []scheduledEvent
	timers        []*time.Timer
	msgCounter    int
}

// NewDemoClient creates a demo client with ITYSL-themed data.
func NewDemoClient() *DemoClient {
	return &DemoClient{
		ownURN:        demoOwnURN,
		conversations: buildDemoConversations(),
		messages:      buildDemoMessages(),
		autoReplies:   buildDemoAutoReplies(),
		replyIndex:    make(map[string]int),
		timeline:      buildDemoTimeline(),
	}
}

func (c *DemoClient) OwnURN() linkedingo.URN {
	return c.ownURN
}

func (c *DemoClient) SetProgram(p *tea.Program) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.program = p
}

func (c *DemoClient) ValidateAuth() tea.Cmd {
	return func() tea.Msg {
		return AuthValidatedMsg{
			Username: "Demo User",
			UserURN:  demoOwnURN,
		}
	}
}

func (c *DemoClient) FetchConversations() tea.Cmd {
	return func() tea.Msg {
		return ConversationsLoadedMsg{
			Conversations: c.conversations,
		}
	}
}

func (c *DemoClient) FetchConversationsBefore(_ time.Time) tea.Cmd {
	return func() tea.Msg {
		return ConversationsLoadedMsg{
			Conversations: c.conversations,
		}
	}
}

func (c *DemoClient) FetchMessages(conversationURN linkedingo.URN, _ time.Time, _ int) tea.Cmd {
	convID := conversationURN.String()
	msgs := c.messages[convID]
	return func() tea.Msg {
		return MessagesLoadedMsg{
			ConversationID: convID,
			Messages:       msgs,
		}
	}
}

func (c *DemoClient) FetchMessagesWithCursor(conversationURN linkedingo.URN, _ string, _ int) tea.Cmd {
	convID := conversationURN.String()
	msgs := c.messages[convID]
	return func() tea.Msg {
		return MessagesLoadedMsg{
			ConversationID: convID,
			Messages:       msgs,
		}
	}
}

func (c *DemoClient) SendMessage(conversationURN linkedingo.URN, text string) tea.Cmd {
	c.mu.Lock()
	c.msgCounter++
	msgID := fmt.Sprintf("msg-demo-sent-%d", c.msgCounter)
	convID := conversationURN.String()
	c.mu.Unlock()

	sentMsg := DisplayMessage{
		ID:        msgID,
		Sender:    "You",
		SenderURN: c.ownURN,
		Body:      text,
		Timestamp: time.Now(),
		IsOwn:     true,
	}

	// Schedule auto-reply
	c.scheduleAutoReply(convID)

	return func() tea.Msg {
		return MessageSentMsg{
			ConversationID: convID,
			Message:        sentMsg,
		}
	}
}

func (c *DemoClient) scheduleAutoReply(convID string) {
	c.mu.Lock()
	replies, ok := c.autoReplies[convID]
	if !ok || len(replies) == 0 {
		c.mu.Unlock()
		return
	}
	idx := c.replyIndex[convID]
	if idx >= len(replies) {
		idx = 0 // loop back
	}
	replyText := replies[idx]
	c.replyIndex[convID] = idx + 1
	c.msgCounter++
	msgID := fmt.Sprintf("msg-demo-reply-%d", c.msgCounter)

	// Find sender info for this conversation
	var senderName string
	var senderURN linkedingo.URN
	for _, conv := range c.conversations {
		if conv.ID == convID {
			for _, p := range conv.Participants {
				if !p.IsOwnUser {
					senderName = p.Name
					senderURN = p.URN
					break
				}
			}
			break
		}
	}
	c.mu.Unlock()

	// Send typing indicator 800ms after user sends
	typingTimer := time.AfterFunc(800*time.Millisecond, func() {
		c.mu.Lock()
		p := c.program
		c.mu.Unlock()
		if p == nil {
			return
		}
		p.Send(RealtimeTypingMsg{
			ConversationID: convID,
			SenderName:     senderName,
		})
	})

	// Send reply at 3s
	timer := time.AfterFunc(3*time.Second, func() {
		c.mu.Lock()
		p := c.program
		c.mu.Unlock()
		if p == nil {
			return
		}

		p.Send(RealtimeMessageMsg{
			ConversationID: convID,
			Message: DisplayMessage{
				ID:        msgID,
				Sender:    senderName,
				SenderURN: senderURN,
				Body:      replyText,
				Timestamp: time.Now(),
			},
		})
	})

	c.mu.Lock()
	c.timers = append(c.timers, typingTimer, timer)
	c.mu.Unlock()
}

func (c *DemoClient) MarkRead(urn linkedingo.URN) tea.Cmd {
	for i := range c.conversations {
		if c.conversations[i].URN == urn {
			c.conversations[i].Unread = false
			break
		}
	}
	return func() tea.Msg { return nil }
}

func (c *DemoClient) MarkUnread(_ linkedingo.URN) tea.Cmd {
	return func() tea.Msg { return nil }
}

func (c *DemoClient) StartTyping(_ linkedingo.URN) tea.Cmd {
	return func() tea.Msg { return nil }
}

func (c *DemoClient) DeleteConversation(conversationURN linkedingo.URN) tea.Cmd {
	convID := conversationURN.String()
	return func() tea.Msg {
		return ConversationDeletedMsg{ConversationID: convID}
	}
}

func (c *DemoClient) ConnectRealtime() tea.Cmd {
	return func() tea.Msg {
		c.mu.Lock()
		p := c.program
		c.mu.Unlock()

		// Start scripted timeline events
		for _, event := range c.timeline {
			evt := event // capture

			// Send typing indicator 2s before the message
			typingDelay := evt.delay - 2*time.Second
			if typingDelay < 0 {
				typingDelay = 0
			}
			typingTimer := time.AfterFunc(typingDelay, func() {
				c.mu.Lock()
				p := c.program
				c.mu.Unlock()
				if p == nil {
					return
				}
				p.Send(RealtimeTypingMsg{
					ConversationID: evt.conversationID,
					SenderName:     evt.message.Sender,
				})
			})
			c.mu.Lock()
			c.timers = append(c.timers, typingTimer)
			c.mu.Unlock()

			// Send message at scheduled time
			timer := time.AfterFunc(evt.delay, func() {
				c.mu.Lock()
				p := c.program
				c.mu.Unlock()
				if p == nil {
					return
				}
				evt.message.Timestamp = time.Now()
				p.Send(RealtimeMessageMsg{
					ConversationID: evt.conversationID,
					Message:        evt.message,
				})
			})
			c.mu.Lock()
			c.timers = append(c.timers, timer)
			c.mu.Unlock()
		}

		// Signal connected after a brief delay
		if p != nil {
			p.Send(RealtimeConnectedMsg{})
		}
		return nil
	}
}

func (c *DemoClient) DisconnectRealtime() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range c.timers {
		t.Stop()
	}
	c.timers = nil
}
