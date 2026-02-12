package linkedin

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"
)

// MessagingClient defines the interface for LinkedIn messaging operations.
// Both the real Client and DemoClient implement this interface.
type MessagingClient interface {
	ValidateAuth() tea.Cmd
	FetchConversations() tea.Cmd
	FetchConversationsBefore(before time.Time) tea.Cmd
	FetchMessages(conversationURN linkedingo.URN, before time.Time, count int) tea.Cmd
	FetchMessagesWithCursor(conversationURN linkedingo.URN, prevCursor string, count int) tea.Cmd
	SendMessage(conversationURN linkedingo.URN, text string) tea.Cmd
	MarkRead(conversationURN linkedingo.URN) tea.Cmd
	MarkUnread(conversationURN linkedingo.URN) tea.Cmd
	StartTyping(conversationURN linkedingo.URN) tea.Cmd
	DeleteConversation(conversationURN linkedingo.URN) tea.Cmd
	ConnectRealtime() tea.Cmd
	DisconnectRealtime()
	SetProgram(p *tea.Program)
	OwnURN() linkedingo.URN
}
