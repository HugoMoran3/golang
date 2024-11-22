package main

import (
	"context"
	"fmt"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,"Hi there, I love to transcribe audio!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func transcribe() {
	c := openai.NewClient("your token")
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "recording.mp3",
	}
	resp, err := c.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp.Text)
}
