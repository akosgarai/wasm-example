package client

import "github.com/akosgarai/wasm-example/pkg/page"

// Client represents the client side of the application
type Client struct {
	Page page.Page
}

// New returns a new Client
func New(p page.Page) *Client {
	return &Client{
		Page: p,
	}
}

// Run runs the client side of the application
func (c *Client) Run() {
	c.Page.LoadPage()
	c.Page.Run()
}
