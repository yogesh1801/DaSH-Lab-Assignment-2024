package main

import (
	"log"
	"strings"

	"github.com/jpoz/groq"
)

func groqClient(groq_client *groq.Client, message string) (string, error) {
	client := groq_client

	response, err := client.CreateChatCompletion(groq.CompletionCreateParams{
		Model: "llama3-8b-8192",
		Messages: []groq.Message{
			{
				Role:    "user",
				Content: message,
			},
		},
	})

	if err != nil {
		log.Printf("Error encoutered while talking to the api: %v", err)
		return "", err
	}

	result := response.Choices[0].Message.Content
	lines := strings.Split(result, "\n")
	result = lines[0]

	return result, nil
}
