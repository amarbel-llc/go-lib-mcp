// Package main demonstrates a simple MCP server using go-lib-mcp.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/friedenberg/go-lib-mcp/protocol"
	"github.com/friedenberg/go-lib-mcp/server"
	"github.com/friedenberg/go-lib-mcp/transport"
)

func main() {
	// Create stdio transport (MCP uses newline-delimited JSON)
	t := transport.NewStdio(os.Stdin, os.Stdout)

	// Create tool registry and register example tools
	tools := server.NewToolRegistry()
	registerTools(tools)

	// Create resource registry and register example resources
	resources := server.NewResourceRegistry()
	registerResources(resources)

	// Create prompt registry and register example prompts
	prompts := server.NewPromptRegistry()
	registerPrompts(prompts)

	// Create server with all providers
	srv, err := server.New(t, server.Options{
		ServerName:    "example-server",
		ServerVersion: "1.0.0",
		Tools:         tools,
		Resources:     resources,
		Prompts:       prompts,
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Run server
	log.Println("Starting MCP example server...")
	if err := srv.Run(context.Background()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerTools(tools *server.ToolRegistry) {
	// Add tool: current time
	tools.Register(
		"get_current_time",
		"Returns the current time in RFC3339 format",
		json.RawMessage(`{"type": "object", "properties": {}}`),
		func(ctx context.Context, args json.RawMessage) (*protocol.ToolCallResult, error) {
			currentTime := time.Now().Format(time.RFC3339)
			return &protocol.ToolCallResult{
				Content: []protocol.ContentBlock{
					protocol.TextContent(fmt.Sprintf("Current time: %s", currentTime)),
				},
			}, nil
		},
	)

	// Add tool: echo
	tools.Register(
		"echo",
		"Echoes back the provided message",
		json.RawMessage(`{
			"type": "object",
			"properties": {
				"message": {"type": "string", "description": "The message to echo back"}
			},
			"required": ["message"]
		}`),
		func(ctx context.Context, args json.RawMessage) (*protocol.ToolCallResult, error) {
			var params struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(args, &params); err != nil {
				return protocol.ErrorResult(fmt.Sprintf("Invalid arguments: %v", err)), nil
			}

			return &protocol.ToolCallResult{
				Content: []protocol.ContentBlock{
					protocol.TextContent(fmt.Sprintf("Echo: %s", params.Message)),
				},
			}, nil
		},
	)
}

func registerResources(resources *server.ResourceRegistry) {
	// Add a simple text resource
	resources.RegisterResource(
		protocol.Resource{
			URI:         "example://greeting",
			Name:        "Example Greeting",
			Description: "A simple greeting message",
			MimeType:    "text/plain",
		},
		func(ctx context.Context, uri string) (*protocol.ResourceReadResult, error) {
			return &protocol.ResourceReadResult{
				Contents: []protocol.ResourceContent{
					{
						URI:      uri,
						MimeType: "text/plain",
						Text:     "Hello from the example MCP server!",
					},
				},
			}, nil
		},
	)

	// Add a JSON resource
	resources.RegisterResource(
		protocol.Resource{
			URI:         "example://info",
			Name:        "Server Info",
			Description: "Information about this server",
			MimeType:    "application/json",
		},
		func(ctx context.Context, uri string) (*protocol.ResourceReadResult, error) {
			info := map[string]interface{}{
				"server":  "example-server",
				"version": "1.0.0",
				"tools":   []string{"get_current_time", "echo"},
			}
			data, _ := json.MarshalIndent(info, "", "  ")
			return &protocol.ResourceReadResult{
				Contents: []protocol.ResourceContent{
					{
						URI:      uri,
						MimeType: "application/json",
						Text:     string(data),
					},
				},
			}, nil
		},
	)
}

func registerPrompts(prompts *server.PromptRegistry) {
	// Add a simple prompt
	prompts.Register(
		protocol.Prompt{
			Name:        "greeting",
			Description: "A friendly greeting prompt",
			Arguments: []protocol.PromptArgument{
				{
					Name:        "name",
					Description: "The name to greet",
					Required:    false,
				},
			},
		},
		func(ctx context.Context, args map[string]string) (*protocol.PromptGetResult, error) {
			name := args["name"]
			if name == "" {
				name = "friend"
			}

			return &protocol.PromptGetResult{
				Description: "A friendly greeting",
				Messages: []protocol.PromptMessage{
					{
						Role:    "user",
						Content: protocol.TextContent(fmt.Sprintf("Say hello to %s", name)),
					},
				},
			}, nil
		},
	)
}
