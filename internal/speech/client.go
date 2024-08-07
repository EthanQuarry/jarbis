package speech

import (
	"context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
)

type Client struct {
	client *speech.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) StreamingRecognize(ctx context.Context) (speechpb.Speech_StreamingRecognizeClient, error) {
	return c.client.StreamingRecognize(ctx)
}