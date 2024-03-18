package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	// Generate data for users
	fileUsers, err := os.Create("test_users_data.sql")
	if err != nil {
		panic(err)
	}

	wrtUsers := bufio.NewWriter(fileUsers)
	_, _ = wrtUsers.Write([]byte("INSERT INTO users (id, name, password) VALUES\n"))
	for i := 0; i <= 24; i++ {
		_, _ = wrtUsers.Write(
			[]byte(fmt.Sprintf("(%d, '%s%d', '%s%d'),\n", i, "user", i, "password", i)),
		)
	}
	_, _ = wrtUsers.Write([]byte(fmt.Sprintf("(%d, '%s', '%s');\n", 25, "admin", "admin")))
	wrtUsers.Flush()
	fileUsers.Close()

	// Generate data for conversation
	fileConversation, err := os.Create("test_conversation_data.sql")
	if err != nil {
		panic(err)
	}

	wrtConversation := bufio.NewWriter(fileConversation)
	_, _ = wrtConversation.Write(
		[]byte("INSERT INTO conversation (id, title, creator_id) VALUES\n"),
	)
	for i := 0; i <= 12; i++ {
		_, _ = wrtConversation.Write(
			[]byte(fmt.Sprintf("(%d, '%s %d', %d),\n", i, "Conversation", i, i)),
		)
	}
	_, _ = wrtConversation.Write([]byte(fmt.Sprintf("(%d, '%s', %d);\n", 13, "General", 13)))
	wrtConversation.Flush()
	fileConversation.Close()

	// Generate data for participants
	fileParticipants, err := os.Create("test_participants_data.sql")
	if err != nil {
		panic(err)
	}

	wrtParticipants := bufio.NewWriter(fileParticipants)
	_, _ = wrtParticipants.Write(
		[]byte("INSERT INTO participants (id, users_id, conversation_id) VALUES\n"),
	)
	for i := 0; i <= 12; i++ {
		_, _ = wrtParticipants.Write(
			[]byte(fmt.Sprintf("(%d, %d, %d),\n", 4*i, 2*i, i)),
		)
		_, _ = wrtParticipants.Write(
			[]byte(fmt.Sprintf("(%d, %d, %d),\n", 4*i+1, 2*i+1, i)),
		)

		_, _ = wrtParticipants.Write(
			[]byte(fmt.Sprintf("(%d, %d, %d),\n", 4*i+2, 2*i, 13)),
		)
		_, _ = wrtParticipants.Write(
			[]byte(fmt.Sprintf("(%d, %d, %d),\n", 4*i+3, 2*i+1, 13)),
		)

		_, _ = wrtParticipants.Write(
			[]byte(fmt.Sprintf("(%d, %d, %d),\n", 52+i, 25, i)),
		)
	}
	_, _ = wrtParticipants.Write(
		[]byte(fmt.Sprintf("(%d, %d, %d);\n", 65, 25, 13)),
	)
	wrtParticipants.Flush()
	fileConversation.Close()

	// Generate data for messages
	fileMessages, err := os.Create("test_messages_data.sql")
	if err != nil {
		panic(err)
	}
	wrtMessages := bufio.NewWriter(fileMessages)

	_, _ = wrtMessages.Write(
		[]byte(
			"INSERT INTO messages (guid, sender_id, conversation_id, message, created_at) VALUES\n",
		),
	)
	for i := 0; i < 1000; i++ {
		_, _ = wrtMessages.Write(
			[]byte(
				fmt.Sprintf(
					"('%s%d', %d, %d, '%s%d', NOW()),\n",
					"message_guid_",
					i,
					rand.Intn(26),
					rand.Intn(14),
					"This is message ",
					i,
				),
			),
		)
	}
	_, _ = wrtMessages.Write(
		[]byte(
			fmt.Sprintf(
				"('%s%d', %d, %d, '%s', NOW());\n",
				"message_guid_",
				1000,
				25,
				13,
				"End message",
			),
		),
	)
	wrtMessages.Flush()
	fileMessages.Close()
}
