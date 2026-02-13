package app

import (
	"context"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"

	"github.com/ggfevans/endorse/internal/config"
	"github.com/ggfevans/endorse/internal/linkedin"
	"github.com/ggfevans/endorse/internal/ui/compose"
	"github.com/ggfevans/endorse/internal/ui/convlist"
	"github.com/ggfevans/endorse/internal/ui/header"
	"github.com/ggfevans/endorse/internal/ui/layout"
	"github.com/ggfevans/endorse/internal/ui/modal"
	"github.com/ggfevans/endorse/internal/ui/statusbar"
	"github.com/ggfevans/endorse/internal/ui/styles"
	"github.com/ggfevans/endorse/internal/ui/thread"
	"github.com/ggfevans/endorse/internal/util"
)

// Model is the root application model.
type Model struct {
	// State
	state    AppState
	focus    FocusedPanel
	ready    bool
	quitting bool

	// Config
	cfg   config.Config
	theme config.Theme

	// Styles
	styles styles.Styles

	// Layout
	dims layout.Dimensions

	// Child components
	header       header.Model
	statusBar    statusbar.Model
	convList     convlist.Model
	thread       thread.Model
	compose      compose.Model
	authModal    modal.AuthModel
	confirmModal modal.ConfirmModel

	// LinkedIn client
	client *linkedin.Client
	ctx    context.Context

	// User info
	username string
	userURN  linkedingo.URN

	// Conversation data (mapped by ID for quick lookup)
	conversations []linkedin.DisplayConversation
	prevCursor    string // for message pagination

	// Pending delete (conversation ID awaiting confirmation)
	pendingDeleteID string

	// Pending credentials (saved between auth submit and validation)
	pendingCreds *config.Credentials
}

// New creates a new application model.
func New() Model {
	cfg, _ := config.Load()
	theme := config.ThemeByName(cfg.ThemeName)
	s := styles.New(theme)
	creds, _ := config.LoadCredentials()

	noop := zerolog.Nop()
	ctx := noop.WithContext(context.Background())

	startState := StateAuth
	if !creds.IsEmpty() {
		startState = StateLoading
	}

	m := Model{
		state:        startState,
		focus:        FocusConvList,
		cfg:          cfg,
		theme:        theme,
		styles:       s,
		ctx:          ctx,
		header:       header.New(s),
		statusBar:    statusbar.New(s),
		convList:     convlist.New(s),
		thread:       thread.New(s),
		compose:      compose.New(s),
		authModal:    modal.NewAuth(s),
		confirmModal: modal.NewConfirm(s),
	}

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	creds, _ := config.LoadCredentials()
	if !creds.IsEmpty() {
		// Try to validate stored credentials
		client, err := linkedin.New(m.ctx, creds.Cookie, creds.PageInstance, creds.XLiTrack)
		if err != nil {
			m.state = StateAuth
			return nil
		}
		m.client = client
		return client.ValidateAuth()
	}
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.dims = layout.Calculate(msg.Width, msg.Height)
		m.ready = true
		m.updateSizes()
		m.authModal.SetSize(msg.Width, msg.Height)
		m.confirmModal.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case ProgramRefMsg:
		if m.client != nil {
			m.client.SetProgram(msg.Program)
		}
		return m, nil

	case ClearErrorMsg:
		m.statusBar.ClearError()
		return m, nil

	// Auth flow messages
	case modal.AuthSubmitMsg:
		return m.handleAuthSubmit(msg)

	case linkedin.AuthValidatedMsg:
		return m.handleAuthValidated(msg)

	case linkedin.AuthFailedMsg:
		return m.handleAuthFailed(msg)

	// Conversation messages
	case linkedin.ConversationsLoadedMsg:
		return m.handleConversationsLoaded(msg)

	case linkedin.ConversationsLoadFailedMsg:
		m.statusBar.SetError(msg.Err.Error())
		return m, clearErrorAfter()

	// Message messages
	case linkedin.MessagesLoadedMsg:
		return m.handleMessagesLoaded(msg)

	case linkedin.MessagesLoadFailedMsg:
		m.statusBar.SetError(msg.Err.Error())
		return m, clearErrorAfter()

	// Send messages
	case linkedin.MessageSentMsg:
		return m.handleMessageSent(msg)

	case linkedin.MessageSendFailedMsg:
		m.statusBar.SetError("Failed to send: " + msg.Err.Error())
		return m, clearErrorAfter()

	// Real-time messages
	case linkedin.RealtimeMessageMsg:
		return m.handleRealtimeMessage(msg)

	case linkedin.RealtimeTypingMsg:
		// Could show typing indicator in thread
		return m, nil

	case linkedin.RealtimeConnectedMsg:
		m.header.SetConnected(true)
		m.statusBar.SetConnected(true)
		return m, nil

	case linkedin.RealtimeDisconnectedMsg:
		m.header.SetConnected(false)
		m.statusBar.SetConnected(false)
		return m, nil

	case linkedin.ConversationDeletedMsg:
		return m.handleConversationDeleted(msg)

	case linkedin.ConversationDeleteFailedMsg:
		m.statusBar.SetError("Failed to delete: " + msg.Err.Error())
		return m, clearErrorAfter()

	case linkedin.SessionExpiredMsg:
		_ = config.ClearCredentials()
		m.state = StateAuth
		m.client = nil
		m.header.SetConnected(false)
		m.statusBar.SetConnected(false)
		m.authModal.SetError("Session expired. Please re-authenticate.")
		return m, nil
	}

	// Forward to focused component
	return m.updateFocused(msg)
}

// --- Auth handlers ---

func (m Model) handleAuthSubmit(msg modal.AuthSubmitMsg) (tea.Model, tea.Cmd) {
	m.authModal.SetLoading(true)

	client, err := linkedin.New(m.ctx, msg.Cookie, msg.PageInstance, msg.XLiTrack)
	if err != nil {
		m.authModal.SetError(err.Error())
		return m, nil
	}

	m.client = client
	m.pendingCreds = &config.Credentials{
		Cookie:       msg.Cookie,
		PageInstance: msg.PageInstance,
		XLiTrack:     msg.XLiTrack,
	}
	return m, client.ValidateAuth()
}

func (m Model) handleAuthValidated(msg linkedin.AuthValidatedMsg) (tea.Model, tea.Cmd) {
	m.username = msg.Username
	m.userURN = msg.UserURN
	m.header.SetUsername(m.username)
	m.statusBar.SetUsername(m.username)

	// Determine credentials source: pending (from auth modal) or stored (from disk)
	var creds config.Credentials
	if m.pendingCreds != nil {
		creds = *m.pendingCreds
		// Save newly validated credentials
		_ = config.SaveCredentials(creds)
		m.pendingCreds = nil
	} else {
		creds, _ = config.LoadCredentials()
	}

	// Recreate client with proper URN
	newClient, err := linkedin.NewWithURN(m.ctx, creds.Cookie, creds.PageInstance, creds.XLiTrack, msg.UserURN)
	if err == nil {
		// Preserve program reference if the old client had one
		if m.client != nil {
			// The program reference will be set via ProgramRefMsg
		}
		m.client = newClient
	}

	m.state = StateMessaging
	m.header.SetConnected(true)
	m.statusBar.SetConnected(true)
	m.convList.Focus()
	m.setFocus(FocusConvList)

	var cmds []tea.Cmd
	if m.client != nil {
		cmds = append(cmds, m.client.FetchConversations())
		cmds = append(cmds, m.client.ConnectRealtime())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleAuthFailed(msg linkedin.AuthFailedMsg) (tea.Model, tea.Cmd) {
	m.state = StateAuth
	m.authModal.SetError("Authentication failed: " + msg.Err.Error())
	return m, nil
}

// --- Filtering ---

func (m *Model) applyConversationFilter() {
	var items []convlist.Conversation
	for _, dc := range m.conversations {
		// Filter: 0=all, 1=unread only
		if m.convList.FilterTab() == 1 && !dc.Unread {
			continue
		}
		items = append(items, convlist.Conversation{
			ID:          dc.ID,
			Name:        dc.Title,
			LastMessage: dc.LastMessage,
			Timestamp:   util.RelativeTime(dc.LastActivityAt),
			Unread:      dc.Unread,
		})
	}
	m.convList.SetConversations(items)
}

// --- Conversation handlers ---

func (m Model) handleConversationsLoaded(msg linkedin.ConversationsLoadedMsg) (tea.Model, tea.Cmd) {
	m.conversations = msg.Conversations

	// Sort by last activity (most recent first)
	sort.Slice(m.conversations, func(i, j int) bool {
		return m.conversations[i].LastActivityAt.After(m.conversations[j].LastActivityAt)
	})

	// Apply current filter and update convlist
	m.applyConversationFilter()
	m.updateFilterCounts()

	return m, nil
}

// --- Message handlers ---

func (m Model) handleMessagesLoaded(msg linkedin.MessagesLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.ConversationID != m.thread.ConversationID() {
		return m, nil // stale response
	}

	m.prevCursor = msg.PrevCursor

	var msgs []thread.Message
	for _, dm := range msg.Messages {
		msgs = append(msgs, thread.Message{
			ID:        dm.ID,
			Sender:    dm.Sender,
			Body:      dm.Body,
			Timestamp: util.RelativeTime(dm.Timestamp),
			IsOwn:     dm.IsOwn,
		})
	}
	m.thread.SetMessages(msgs)

	return m, nil
}

func (m Model) handleMessageSent(msg linkedin.MessageSentMsg) (tea.Model, tea.Cmd) {
	if msg.ConversationID == m.thread.ConversationID() {
		m.thread.AppendMessage(thread.Message{
			ID:        msg.Message.ID,
			Sender:    msg.Message.Sender,
			Body:      msg.Message.Body,
			Timestamp: util.RelativeTime(msg.Message.Timestamp),
			IsOwn:     msg.Message.IsOwn,
		})
	}
	return m, nil
}

func (m Model) handleRealtimeMessage(msg linkedin.RealtimeMessageMsg) (tea.Model, tea.Cmd) {
	// If the message is for the currently viewed conversation, append it
	if msg.ConversationID == m.thread.ConversationID() {
		m.thread.AppendMessage(thread.Message{
			ID:        msg.Message.ID,
			Sender:    msg.Message.Sender,
			Body:      msg.Message.Body,
			Timestamp: util.RelativeTime(msg.Message.Timestamp),
			IsOwn:     msg.Message.IsOwn,
		})
	}

	// Refresh conversations to update order and unread counts
	if m.client != nil {
		return m, m.client.FetchConversations()
	}
	return m, nil
}

// --- Key handling ---

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Ctrl+C always quits, regardless of state
	if msg.String() == "ctrl+c" {
		m.quitting = true
		if m.client != nil {
			m.client.DisconnectRealtime()
		}
		return m, tea.Quit
	}

	// Auth state: forward keys to auth modal
	if m.state == StateAuth {
		var cmd tea.Cmd
		m.authModal, cmd = m.authModal.Update(msg)
		return m, cmd
	}

	// Confirm modal intercepts all keys when active
	if m.confirmModal.Active() {
		return m.handleConfirmKey(msg)
	}

	// Global keys (messaging state)
	if isQuitKey(msg) && m.focus != FocusCompose {
		m.quitting = true
		if m.client != nil {
			m.client.DisconnectRealtime()
		}
		return m, tea.Quit
	}

	if isTabKey(msg) {
		m.cycleFocusForward()
		return m, nil
	}

	if isShiftTabKey(msg) {
		m.cycleFocusBackward()
		return m, nil
	}

	// Route to focused panel
	switch m.focus {
	case FocusConvList:
		return m.handleConvListKey(msg)
	case FocusThread:
		return m.handleThreadKey(msg)
	case FocusCompose:
		return m.handleComposeKey(msg)
	}

	return m, nil
}

func (m Model) handleConvListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isFilterKey(msg):
		m.convList.ToggleFilter()
		m.applyConversationFilter()
		return m, nil
	case isDownKey(msg):
		m.convList.MoveDown()
	case isUpKey(msg):
		m.convList.MoveUp()
	case isTopKey(msg):
		m.convList.MoveToTop()
	case isBottomKey(msg):
		m.convList.MoveToBottom()
	case isEnterKey(msg):
		return m.openSelectedConversation()
	case isReplyKey(msg):
		return m.openSelectedConversationAndReply()
	case isMarkReadKey(msg):
		return m.toggleSelectedReadState()
	case isDeleteKey(msg):
		return m.promptDeleteSelected()
	}
	return m, nil
}

func (m Model) handleThreadKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isUpKey(msg):
		m.thread.ScrollUp(1)
	case isDownKey(msg):
		m.thread.ScrollDown(1)
	case isPageUp(msg):
		m.thread.ScrollUp(m.thread.VisibleHeight() / 2)
	case isPageDown(msg):
		m.thread.ScrollDown(m.thread.VisibleHeight() / 2)
	case isReplyKey(msg):
		if m.thread.HasConversation() {
			m.activateCompose()
		}
	case isEscapeKey(msg):
		m.setFocus(FocusConvList)
	}
	return m, nil
}

func (m Model) handleComposeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isEscapeKey(msg):
		m.compose.Deactivate()
		m.setFocus(FocusThread)
		return m, nil
	case isSendKey(msg):
		return m.sendMessage()
	}

	// Forward to textarea
	var cmd tea.Cmd
	m.compose, cmd = m.compose.Update(msg)
	return m, cmd
}

// --- Actions ---

func (m Model) openSelectedConversation() (tea.Model, tea.Cmd) {
	conv, ok := m.convList.SelectedConversation()
	if !ok {
		return m, nil
	}

	m.thread.SetConversation(conv.ID, conv.Name)
	m.setFocus(FocusThread)

	var cmds []tea.Cmd
	if m.client != nil {
		urn := m.findConversationURN(conv.ID)
		if !urn.IsEmpty() {
			cmds = append(cmds, m.client.FetchMessages(urn, time.Now(), 20))
			cmds = append(cmds, m.client.MarkRead(urn))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) toggleSelectedReadState() (tea.Model, tea.Cmd) {
	conv, ok := m.convList.SelectedConversation()
	if !ok {
		return m, nil
	}

	// Toggle locally
	var nowUnread bool
	for i := range m.conversations {
		if m.conversations[i].ID == conv.ID {
			m.conversations[i].Unread = !m.conversations[i].Unread
			nowUnread = m.conversations[i].Unread
			break
		}
	}
	m.applyConversationFilter()
	m.updateFilterCounts()

	// Sync to server
	if m.client != nil {
		urn := m.findConversationURN(conv.ID)
		if !urn.IsEmpty() {
			if nowUnread {
				return m, m.client.MarkUnread(urn)
			}
			return m, m.client.MarkRead(urn)
		}
	}

	return m, nil
}

func (m Model) promptDeleteSelected() (tea.Model, tea.Cmd) {
	conv, ok := m.convList.SelectedConversation()
	if !ok {
		return m, nil
	}

	m.pendingDeleteID = conv.ID
	m.confirmModal.Show("Delete conversation with " + conv.Name + "?")
	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isEnterKey(msg):
		m.confirmModal.Hide()
		return m.deleteConversation(m.pendingDeleteID)
	case isEscapeKey(msg):
		m.confirmModal.Hide()
		m.pendingDeleteID = ""
		return m, nil
	}
	return m, nil
}

func (m Model) deleteConversation(id string) (tea.Model, tea.Cmd) {
	// Look up URN before removing from local state
	urn := m.findConversationURN(id)

	// Remove from local state
	for i := range m.conversations {
		if m.conversations[i].ID == id {
			m.conversations = append(m.conversations[:i], m.conversations[i+1:]...)
			break
		}
	}
	m.applyConversationFilter()
	m.updateFilterCounts()

	// If we were viewing this conversation, clear the thread
	if m.thread.ConversationID() == id {
		m.thread.SetConversation("", "")
	}

	m.pendingDeleteID = ""

	// Delete on server
	if m.client != nil && !urn.IsEmpty() {
		return m, m.client.DeleteConversation(urn)
	}

	return m, nil
}

func (m Model) handleConversationDeleted(msg linkedin.ConversationDeletedMsg) (tea.Model, tea.Cmd) {
	// Already removed from local state in deleteConversation; nothing more to do
	return m, nil
}

func (m *Model) updateFilterCounts() {
	unreadCount := 0
	for _, dc := range m.conversations {
		if dc.Unread {
			unreadCount++
		}
	}
	m.convList.SetFilterCounts(len(m.conversations), unreadCount)
}

func (m Model) openSelectedConversationAndReply() (tea.Model, tea.Cmd) {
	newM, cmd := m.openSelectedConversation()
	m = newM.(Model)
	m.activateCompose()
	return m, cmd
}

func (m *Model) activateCompose() {
	// Find conversation title for recipient
	convID := m.thread.ConversationID()
	for _, dc := range m.conversations {
		if dc.ID == convID {
			m.compose.SetRecipient(dc.Title)
			break
		}
	}
	m.compose.Focus()
	m.setFocus(FocusCompose)
}

func (m Model) sendMessage() (tea.Model, tea.Cmd) {
	text := m.compose.Value()
	if text == "" {
		return m, nil
	}

	convID := m.thread.ConversationID()
	m.compose.Reset()
	m.compose.Deactivate()
	m.setFocus(FocusThread)

	if m.client != nil {
		urn := m.findConversationURN(convID)
		if !urn.IsEmpty() {
			return m, m.client.SendMessage(urn, text)
		}
	}

	return m, nil
}

func (m Model) findConversationURN(id string) linkedingo.URN {
	for _, dc := range m.conversations {
		if dc.ID == id {
			return dc.URN
		}
	}
	return linkedingo.URN{}
}

func (m Model) updateFocused(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == StateAuth {
		var cmd tea.Cmd
		m.authModal, cmd = m.authModal.Update(msg)
		return m, cmd
	}

	if m.focus == FocusCompose {
		var cmd tea.Cmd
		m.compose, cmd = m.compose.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) cycleFocusForward() {
	panels := m.availablePanels()
	for i, p := range panels {
		if p == m.focus {
			next := panels[(i+1)%len(panels)]
			m.setFocus(next)
			return
		}
	}
	m.setFocus(panels[0])
}

func (m *Model) cycleFocusBackward() {
	panels := m.availablePanels()
	for i, p := range panels {
		if p == m.focus {
			prev := panels[(i-1+len(panels))%len(panels)]
			m.setFocus(prev)
			return
		}
	}
	m.setFocus(panels[0])
}

func (m Model) availablePanels() []FocusedPanel {
	switch m.dims.Mode {
	case layout.TwoPanel:
		if m.compose.Active() {
			return []FocusedPanel{FocusConvList, FocusThread, FocusCompose}
		}
		return []FocusedPanel{FocusConvList, FocusThread}
	default:
		return []FocusedPanel{FocusConvList}
	}
}

func (m *Model) setFocus(panel FocusedPanel) {
	m.convList.Blur()
	m.thread.Blur()
	if panel != FocusCompose {
		m.compose.Blur()
	}

	m.focus = panel
	switch panel {
	case FocusConvList:
		m.convList.Focus()
	case FocusThread:
		m.thread.Focus()
	case FocusCompose:
		m.compose.Focus()
	}
}

func (m *Model) updateSizes() {
	m.header.SetWidth(m.dims.Width)
	m.statusBar.SetWidth(m.dims.Width)

	composeH := m.compose.ComposeHeight()
	contentH := m.dims.ContentHeight - composeH

	m.convList.SetSize(m.dims.ConvListWidth, contentH)

	threadW := m.dims.ThreadWidth
	m.thread.SetSize(threadW, contentH)
	m.compose.SetSize(threadW, 5)
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "Loading..."
	}

	// Auth state: show auth modal fullscreen
	if m.state == StateAuth {
		return m.authModal.View()
	}

	// Loading state
	if m.state == StateLoading {
		return lipgloss.Place(m.dims.Width, m.dims.Height,
			lipgloss.Center, lipgloss.Center,
			m.styles.AccentText.Render("Connecting to LinkedIn..."))
	}

	// Confirm modal overlays everything
	if m.confirmModal.Active() {
		return m.confirmModal.View()
	}

	// Messaging state
	var contentView string
	switch m.dims.Mode {
	case layout.TwoPanel:
		contentView = lipgloss.JoinHorizontal(lipgloss.Top,
			m.convList.View(),
			m.threadWithCompose(),
		)
	case layout.SinglePanel:
		contentView = m.threadWithCompose()
	}

	headerView := m.header.View()
	statusView := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		contentView,
		statusView,
	)
}

func (m Model) threadWithCompose() string {
	if m.compose.Active() {
		return lipgloss.JoinVertical(lipgloss.Left,
			m.thread.View(),
			m.compose.View(),
		)
	}
	return m.thread.View()
}
