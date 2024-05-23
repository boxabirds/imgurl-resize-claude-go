package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	claude "github.com/potproject/claude-sdk-go"
)

const DEFAULT_CLAUDE_MODEL = "claude-3-haiku-20240307"

// Sample URLs as constants
var sampleURLs = []string{
	"https://assets.newatlas.com/dims4/default/40c2d71/2147483647/strip/true/crop/864x576+80+0/resize/1920x1280!/quality/90/?url=http%3A%2F%2Fnewatlas-brightspot.s3.amazonaws.com%2F22%2Fe5%2F3e95ec3d4968b61caaa35f6e8868%2Famogy-03-edited-1024x576.jpg",
	"https://static.wixstatic.com/media/95ac29_46a94cb8d61c4ad98470784946371f9a~mv2.jpg/v1/fill/w_1108,h_582,al_c,q_85,usm_0.66_1.00_0.01,enc_auto/95ac29_46a94cb8d61c4ad98470784946371f9a~mv2.jpg",
	"https://heatmap.news/media-library/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpbWFnZSI6Imh0dHBzOi8vYXNzZXRzLnJibC5tcy81MjIxNjQwNi9vcmlnaW4uanBnIiwiZXhwaXJlc19hdCI6MTc3NjQ5NTE1Mn0.BzNvkei_5BccvZSez6hxLWwh-2dRwlJX6GSqxyP5iTk/image.jpg?width=1200&height=800&quality=90&coordinates=0%2C0%2C0%2C0",
	"https://images.fastcompany.com/image/upload/f_auto,c_fit,w_2048,q_auto/wp-cms-2/2024/05/p-1-2024-WCI-Energy_Fourth-Power.jpg",
	"https://i.guim.co.uk/img/media/4a92d2eb1982161d58c1311aedda5f096814c1ce/0_190_5705_3423/master/5705.jpg?width=620&dpr=2&s=none",
	"https://i.abcnewsfe.com/a/ab9d610a-b509-443f-a97d-c91b1168bbae/cork3-abc-ml-240424_1713974902805_hpMain_16x9.jpg?w=992",
}

func main() {
	apiKey, exists := os.LookupEnv("ANTHROPIC_API_KEY")
	if !exists {
		log.Fatal("ANTHROPIC_API_KEY environment variable is not set")
	}

	inputFile := flag.String("input-file", "", "Path to the input text file")
	flag.Parse()

	client := claude.NewClient(apiKey)
	ctx := context.Background()

	urls := sampleURLs
	if *inputFile != "" {
		fileUrls, err := readUrlsFromFile(*inputFile)
		if err != nil {
			log.Fatalf("Error reading input file: %v\n", err)
		}
		urls = fileUrls
	}

	for _, url := range urls {
		upscaledUrl, err := upscaleImgUrl(ctx, client, url)
		if err != nil {
			log.Printf("Error processing URL %s: %v\n", url, err)
			continue
		}
		// print the original url and the upscaled url
		fmt.Println("===")
		fmt.Printf("Original URL: %s\n", url)
		fmt.Printf("Upscaled URL: %s\n", upscaledUrl)

	}
}

func readUrlsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls, scanner.Err()
}

func upscaleImgUrl(ctx context.Context, client *claude.Client, url string) (string, error) {
	const SYSTEM_PROMPT = `"The following image url contains image dimensions, 
	possibly in multiple parts of the URL. 
	Rewrite the url for a 4k screen by changing the last set of dimensions:"`

	req := claude.RequestBodyMessages{
		Model:     DEFAULT_CLAUDE_MODEL,
		MaxTokens: 1000,
		System:    SYSTEM_PROMPT,
		Messages: []claude.RequestBodyMessagesMessages{
			{
				Role:    claude.MessagesRoleUser,
				Content: url,
			},
		},
	}

	resp, err := client.CreateMessages(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Content) > 0 {
		return resp.Content[0].Text, nil
	}
	return "", fmt.Errorf("no response from Claude")
}
