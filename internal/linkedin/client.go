package linkedin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"
)

// Client wraps linkedingo.Client and produces tea.Cmd factories.
type Client struct {
	raw     *linkedingo.Client
	ctx     context.Context
	ownURN  linkedingo.URN
	program *tea.Program
}

// --- tea.Msg types produced by this client ---

type AuthValidatedMsg struct {
	Username string
	UserURN  linkedingo.URN
}

type AuthFailedMsg struct {
	Err error
}

type ConversationsLoadedMsg struct {
	Conversations []DisplayConversation
}

type ConversationsLoadFailedMsg struct {
	Err error
}

type MessagesLoadedMsg struct {
	ConversationID string
	Messages       []DisplayMessage
	PrevCursor     string
}

type MessagesLoadFailedMsg struct {
	Err error
}

type MessageSentMsg struct {
	ConversationID string
	Message        DisplayMessage
}

type MessageSendFailedMsg struct {
	Err error
}

type RealtimeMessageMsg struct {
	ConversationID string
	Message        DisplayMessage
}

type RealtimeTypingMsg struct {
	ConversationID string
	SenderName     string
}

type RealtimeSeenMsg struct {
	ConversationID string
}

type RealtimeConnectedMsg struct{}

type RealtimeDisconnectedMsg struct {
	Err error
}

type SessionExpiredMsg struct{}

type ConversationDeletedMsg struct {
	ConversationID string
}

type ConversationDeleteFailedMsg struct {
	Err error
}

type MarkReadFailedMsg struct {
	Err error
}

type MarkUnreadFailedMsg struct {
	Err error
}

// normalizeCookieInput takes flexible user input and returns a valid cookie header string.
// Accepts: full Cookie header, just li_at=value, just the li_at token value, etc.
func normalizeCookieInput(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("cookie input is empty")
	}

	// Try parsing as-is first (full cookie header)
	cookies, err := http.ParseCookie(raw)
	if err == nil && len(cookies) > 0 {
		// Check if we have li_at
		for _, c := range cookies {
			if c.Name == "li_at" {
				return raw, nil // valid full header
			}
		}
		// Parsed but no li_at found - might be a bare token value
	}

	// Check if it looks like "li_at=..." (single cookie assignment)
	if strings.HasPrefix(raw, "li_at=") {
		// Try parsing as a cookie header
		cookies, err := http.ParseCookie(raw)
		if err == nil && len(cookies) > 0 {
			return raw, nil
		}
		// If parsing fails, extract the value and construct a clean cookie
		val := strings.TrimPrefix(raw, "li_at=")
		return "li_at=" + val, nil
	}

	// Maybe it's a bare li_at token value (no name= prefix)
	// li_at tokens typically start with "AQED" or "AQEF"
	if strings.HasPrefix(raw, "AQE") || strings.HasPrefix(raw, "AQD") {
		return "li_at=" + raw, nil
	}

	// Last resort: try to extract li_at and JSESSIONID from a potentially malformed header
	liAt := extractCookieValue(raw, "li_at")
	jsessionID := extractCookieValue(raw, "JSESSIONID")

	if liAt != "" {
		result := "li_at=" + liAt
		if jsessionID != "" {
			result += "; JSESSIONID=" + jsessionID
		}
		return result, nil
	}

	// Nothing worked - return original and let linkedingo try
	return raw, nil
}

// extractCookieValue extracts a named cookie value from a raw header string
// using string splitting (more tolerant than http.ParseCookie).
func extractCookieValue(header, name string) string {
	// Split by ; and look for name=value
	parts := strings.Split(header, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, name+"=") {
			return strings.TrimPrefix(part, name+"=")
		}
	}
	return ""
}

// New creates a new Client from credentials.
// The userEntityURN can be empty for initial auth validation.
func New(ctx context.Context, cookie, pageInstance, xLiTrack string) (*Client, error) {
	normalized, err := normalizeCookieInput(cookie)
	if err != nil {
		return nil, err
	}

	jar, err := linkedingo.NewJarFromCookieHeader(normalized)
	if err != nil {
		return nil, fmt.Errorf("invalid cookie header: %w", err)
	}

	c := &Client{ctx: ctx}

	handlers := linkedingo.Handlers{
		DecoratedEvent:      c.onDecoratedEvent,
		BadCredentials:      c.onBadCredentials,
		TransientDisconnect: c.onTransientDisconnect,
		ClientConnection:    c.onClientConnection,
	}

	// Create client with empty URN initially
	c.raw = linkedingo.NewClient(ctx, linkedingo.URN{}, jar, pageInstance, xLiTrack, handlers)

	return c, nil
}

// NewWithURN creates a client with a known user URN.
func NewWithURN(ctx context.Context, cookie, pageInstance, xLiTrack string, urn linkedingo.URN) (*Client, error) {
	normalized, err := normalizeCookieInput(cookie)
	if err != nil {
		return nil, err
	}

	jar, err := linkedingo.NewJarFromCookieHeader(normalized)
	if err != nil {
		return nil, fmt.Errorf("invalid cookie header: %w", err)
	}

	c := &Client{ctx: ctx, ownURN: urn}

	handlers := linkedingo.Handlers{
		DecoratedEvent:      c.onDecoratedEvent,
		BadCredentials:      c.onBadCredentials,
		TransientDisconnect: c.onTransientDisconnect,
		ClientConnection:    c.onClientConnection,
	}

	c.raw = linkedingo.NewClient(ctx, urn, jar, pageInstance, xLiTrack, handlers)

	return c, nil
}

// SetProgram sets the tea.Program for sending real-time messages.
func (c *Client) SetProgram(p *tea.Program) {
	c.program = p
}

// OwnURN returns the authenticated user's URN.
func (c *Client) OwnURN() linkedingo.URN {
	return c.ownURN
}

// ValidateAuth validates credentials by fetching the current user profile.
func (c *Client) ValidateAuth() tea.Cmd {
	return func() tea.Msg {
		profile, err := c.raw.GetCurrentUserProfile(c.ctx)
		if err != nil {
			return AuthFailedMsg{Err: err}
		}

		c.ownURN = profile.MiniProfile.EntityURN
		username := profile.MiniProfile.PublicIdentifier
		if username == "" {
			username = profile.MiniProfile.FirstName + " " + profile.MiniProfile.LastName
		}

		return AuthValidatedMsg{
			Username: username,
			UserURN:  profile.MiniProfile.EntityURN,
		}
	}
}

// FetchConversations fetches the conversation list.
func (c *Client) FetchConversations() tea.Cmd {
	return func() tea.Msg {
		resp, err := c.raw.GetConversations(c.ctx)
		if err != nil {
			return ConversationsLoadFailedMsg{Err: err}
		}

		var convs []DisplayConversation
		for _, conv := range resp.Elements {
			convs = append(convs, ConvertConversation(conv, c.ownURN))
		}

		return ConversationsLoadedMsg{Conversations: convs}
	}
}

// FetchConversationsBefore fetches older conversations.
func (c *Client) FetchConversationsBefore(before time.Time) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.raw.GetConversationsUpdatedBefore(c.ctx, before)
		if err != nil {
			return ConversationsLoadFailedMsg{Err: err}
		}

		var convs []DisplayConversation
		for _, conv := range resp.Elements {
			convs = append(convs, ConvertConversation(conv, c.ownURN))
		}

		return ConversationsLoadedMsg{Conversations: convs}
	}
}

// FetchMessages fetches messages for a conversation.
func (c *Client) FetchMessages(conversationURN linkedingo.URN, before time.Time, count int) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.raw.GetMessagesBefore(c.ctx, conversationURN, before, count)
		if err != nil {
			return MessagesLoadFailedMsg{Err: err}
		}

		var msgs []DisplayMessage
		for _, msg := range resp.Elements {
			msgs = append(msgs, ConvertMessage(msg, c.ownURN))
		}

		return MessagesLoadedMsg{
			ConversationID: conversationURN.String(),
			Messages:       msgs,
			PrevCursor:     resp.Metadata.PrevCursor,
		}
	}
}

// FetchMessagesWithCursor fetches older messages using cursor pagination.
func (c *Client) FetchMessagesWithCursor(conversationURN linkedingo.URN, prevCursor string, count int) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.raw.GetMessagesWithPrevCursor(c.ctx, conversationURN, prevCursor, count)
		if err != nil {
			return MessagesLoadFailedMsg{Err: err}
		}

		var msgs []DisplayMessage
		for _, msg := range resp.Elements {
			msgs = append(msgs, ConvertMessage(msg, c.ownURN))
		}

		cursor := ""
		if resp != nil {
			cursor = resp.Metadata.PrevCursor
		}

		return MessagesLoadedMsg{
			ConversationID: conversationURN.String(),
			Messages:       msgs,
			PrevCursor:     cursor,
		}
	}
}

// SendMessage sends a text message to a conversation.
func (c *Client) SendMessage(conversationURN linkedingo.URN, text string) tea.Cmd {
	return func() tea.Msg {
		body := linkedingo.SendMessageBody{Text: text}
		resp, err := c.raw.SendMessage(c.ctx, conversationURN, body, nil, "")
		if err != nil {
			return MessageSendFailedMsg{Err: err}
		}

		dm := ConvertMessage(resp.Data, c.ownURN)
		return MessageSentMsg{
			ConversationID: conversationURN.String(),
			Message:        dm,
		}
	}
}

// MarkRead marks a conversation as read.
func (c *Client) MarkRead(conversationURN linkedingo.URN) tea.Cmd {
	return func() tea.Msg {
		_, err := c.raw.MarkConversationRead(c.ctx, conversationURN)
		if err != nil {
			return MarkReadFailedMsg{Err: err}
		}
		return nil
	}
}

// MarkUnread marks a conversation as unread.
func (c *Client) MarkUnread(conversationURN linkedingo.URN) tea.Cmd {
	return func() tea.Msg {
		_, err := c.raw.MarkConversationUnread(c.ctx, conversationURN)
		if err != nil {
			return MarkUnreadFailedMsg{Err: err}
		}
		return nil
	}
}

// StartTyping sends a typing indicator.
func (c *Client) StartTyping(conversationURN linkedingo.URN) tea.Cmd {
	return func() tea.Msg {
		_ = c.raw.StartTyping(c.ctx, conversationURN)
		return nil
	}
}

// ConnectRealtime starts the real-time SSE connection.
func (c *Client) ConnectRealtime() tea.Cmd {
	return func() tea.Msg {
		err := c.raw.RealtimeConnect(c.ctx)
		if err != nil {
			return RealtimeDisconnectedMsg{Err: err}
		}
		return nil
	}
}

// DisconnectRealtime stops the real-time connection.
func (c *Client) DisconnectRealtime() {
	c.raw.RealtimeDisconnect()
}

// DeleteConversation deletes a conversation.
func (c *Client) DeleteConversation(conversationURN linkedingo.URN) tea.Cmd {
	return func() tea.Msg {
		err := c.raw.DeleteConversation(c.ctx, conversationURN)
		if err != nil {
			return ConversationDeleteFailedMsg{Err: err}
		}
		return ConversationDeletedMsg{ConversationID: conversationURN.String()}
	}
}

