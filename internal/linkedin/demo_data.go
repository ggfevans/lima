package linkedin

import (
	"time"

	"go.mau.fi/mautrix-linkedin/pkg/linkedingo"
)

// Demo URNs — using LinkedIn's URN format with fake IDs
var (
	demoOwnURN   = linkedingo.NewURN("urn:li:member:demo-self")
	demoKarlURN  = linkedingo.NewURN("urn:li:member:demo-karl")
	demoHowieURN = linkedingo.NewURN("urn:li:member:demo-howie")
	demoTammyURN = linkedingo.NewURN("urn:li:member:demo-tammy")
	demoBrianURN = linkedingo.NewURN("urn:li:member:demo-brian")

	demoConvKarlURN  = linkedingo.NewURN("urn:li:conversation:conv-karl")
	demoConvHowieURN = linkedingo.NewURN("urn:li:conversation:conv-howie")
	demoConvTammyURN = linkedingo.NewURN("urn:li:conversation:conv-tammy")
	demoConvBrianURN = linkedingo.NewURN("urn:li:conversation:conv-brian")
)

func buildDemoConversations() []DisplayConversation {
	now := time.Now()
	return []DisplayConversation{
		{
			ID:             demoConvKarlURN.String(),
			Title:          "Karl Havoc",
			LastMessage:    "Do you think they know I'm just a guy?",
			LastActivityAt: now.Add(-2 * time.Minute),
			Unread:         true,
			URN:            demoConvKarlURN,
			Participants: []DisplayParticipant{
				{Name: "Karl Havoc", URN: demoKarlURN},
				{Name: "You", URN: demoOwnURN, IsOwnUser: true},
			},
		},
		{
			ID:             demoConvTammyURN.String(),
			Title:          "Tammy Craps",
			LastMessage:    "I'm done. This is the maddest I've ever been",
			LastActivityAt: now.Add(-8 * time.Minute),
			Unread:         true,
			URN:            demoConvTammyURN,
			Participants: []DisplayParticipant{
				{Name: "Tammy Craps", URN: demoTammyURN},
				{Name: "You", URN: demoOwnURN, IsOwnUser: true},
			},
		},
		{
			ID:             demoConvHowieURN.String(),
			Title:          "Howie",
			LastMessage:    "Tell the kid",
			LastActivityAt: now.Add(-45 * time.Minute),
			Unread:         false,
			URN:            demoConvHowieURN,
			Participants: []DisplayParticipant{
				{Name: "Howie", URN: demoHowieURN},
				{Name: "You", URN: demoOwnURN, IsOwnUser: true},
			},
		},
		{
			ID:             demoConvBrianURN.String(),
			Title:          "Brian",
			LastMessage:    "If you don't give, the whole thing falls apart",
			LastActivityAt: now.Add(-3 * time.Hour),
			Unread:         false,
			URN:            demoConvBrianURN,
			Participants: []DisplayParticipant{
				{Name: "Brian", URN: demoBrianURN},
				{Name: "You", URN: demoOwnURN, IsOwnUser: true},
			},
		},
	}
}

func buildDemoMessages() map[string][]DisplayMessage {
	now := time.Now()
	return map[string][]DisplayMessage{
		demoConvKarlURN.String(): {
			{
				ID:        "msg-karl-1",
				Sender:    "Karl Havoc",
				SenderURN: demoKarlURN,
				Body:      "They moved my desk again. I don't even want to be around anymore",
				Timestamp: now.Add(-10 * time.Minute),
			},
			{
				ID:        "msg-karl-2",
				Sender:    "Karl Havoc",
				SenderURN: demoKarlURN,
				Body:      "I'm not going to do the voice",
				Timestamp: now.Add(-5 * time.Minute),
			},
			{
				ID:        "msg-karl-3",
				Sender:    "Karl Havoc",
				SenderURN: demoKarlURN,
				Body:      "Do you think they know I'm just a guy?",
				Timestamp: now.Add(-2 * time.Minute),
			},
		},
		demoConvTammyURN.String(): {
			{
				ID:        "msg-tammy-1",
				Sender:    "Tammy Craps",
				SenderURN: demoTammyURN,
				Body:      "You told me the gift receipt meant I could return it",
				Timestamp: now.Add(-15 * time.Minute),
			},
			{
				ID:        "msg-tammy-2",
				Sender:    "Tammy Craps",
				SenderURN: demoTammyURN,
				Body:      "That one egg was 40 eggs???",
				Timestamp: now.Add(-12 * time.Minute),
			},
			{
				ID:        "msg-tammy-3",
				Sender:    "Tammy Craps",
				SenderURN: demoTammyURN,
				Body:      "I'm done. This is the maddest I've ever been",
				Timestamp: now.Add(-8 * time.Minute),
			},
		},
		demoConvHowieURN.String(): {
			{
				ID:        "msg-howie-1",
				Sender:    "You",
				SenderURN: demoOwnURN,
				Body:      "How was dinner?",
				Timestamp: now.Add(-50 * time.Minute),
				IsOwn:     true,
			},
			{
				ID:        "msg-howie-2",
				Sender:    "Howie",
				SenderURN: demoHowieURN,
				Body:      "The tables!! They were so small!",
				Timestamp: now.Add(-49 * time.Minute),
			},
			{
				ID:        "msg-howie-3",
				Sender:    "Howie",
				SenderURN: demoHowieURN,
				Body:      "I couldn't fit my laptop on the table. The waiter tried to move my stuff",
				Timestamp: now.Add(-48 * time.Minute),
			},
			{
				ID:        "msg-howie-4",
				Sender:    "Howie",
				SenderURN: demoHowieURN,
				Body:      "I've got triples of the barracuda. Triples is best",
				Timestamp: now.Add(-47 * time.Minute),
			},
			{
				ID:        "msg-howie-5",
				Sender:    "You",
				SenderURN: demoOwnURN,
				Body:      "I don't think you should bring a laptop to dinner",
				Timestamp: now.Add(-46 * time.Minute),
				IsOwn:     true,
			},
			{
				ID:        "msg-howie-6",
				Sender:    "Howie",
				SenderURN: demoHowieURN,
				Body:      "Tell the kid",
				Timestamp: now.Add(-45 * time.Minute),
			},
		},
		demoConvBrianURN.String(): {
			{
				ID:        "msg-brian-1",
				Sender:    "Brian",
				SenderURN: demoBrianURN,
				Body:      "Have you seen CalicosCutPants.com?",
				Timestamp: now.Add(-4 * time.Hour),
			},
			{
				ID:        "msg-brian-2",
				Sender:    "Brian",
				SenderURN: demoBrianURN,
				Body:      "You HAVE to give",
				Timestamp: now.Add(-3*time.Hour - 50*time.Minute),
			},
			{
				ID:        "msg-brian-3",
				Sender:    "You",
				SenderURN: demoOwnURN,
				Body:      "What exactly is the business model?",
				Timestamp: now.Add(-3*time.Hour - 30*time.Minute),
				IsOwn:     true,
			},
			{
				ID:        "msg-brian-4",
				Sender:    "Brian",
				SenderURN: demoBrianURN,
				Body:      "It's like a tip jar but for the whole internet",
				Timestamp: now.Add(-3*time.Hour - 10*time.Minute),
			},
			{
				ID:        "msg-brian-5",
				Sender:    "Brian",
				SenderURN: demoBrianURN,
				Body:      "If you don't give, the whole thing falls apart",
				Timestamp: now.Add(-3 * time.Hour),
			},
		},
	}
}

func buildDemoAutoReplies() map[string][]string {
	return map[string][]string{
		demoConvKarlURN.String(): {
			"I don't WANT to be around anymore",
			"Oh yeah that would be great actually",
			"I'm worried that the baby thinks people can't change",
		},
		demoConvHowieURN.String(): {
			"They have no good food there",
			"I don't care about it, but it's not good behavior",
		},
		demoConvTammyURN.String(): {
			"I feel like I'm back in the pants",
			"You ruined my funeral!",
		},
		demoConvBrianURN.String(): {
			"You gotta give though",
			"The website works because of you",
		},
	}
}

func buildDemoTimeline() []scheduledEvent {
	return []scheduledEvent{
		{
			delay:          10 * time.Second,
			conversationID: demoConvKarlURN.String(),
			message: DisplayMessage{
				ID:        "msg-karl-rt-1",
				Sender:    "Karl Havoc",
				SenderURN: demoKarlURN,
				Body:      "They said I have to do the voice or I'm fired",
				Timestamp: time.Now(),
			},
		},
		{
			delay:          15 * time.Second,
			conversationID: demoConvTammyURN.String(),
			message: DisplayMessage{
				ID:        "msg-tammy-rt-1",
				Sender:    "Tammy Craps",
				SenderURN: demoTammyURN,
				Body:      "I need you to tell me right now that that's not a lot of eggs",
				Timestamp: time.Now(),
			},
		},
	}
}
