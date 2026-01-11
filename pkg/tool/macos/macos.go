package macos

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListDirHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read directory: %v", err)), nil
	}

	type FileInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
		Size int64  `json:"size"`
		Ext  string `json:"ext"`
	}

	var fileInfos []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileType := "file"
		if entry.IsDir() {
			fileType = "directory"
		}

		ext := filepath.Ext(entry.Name())
		fullPath := filepath.Join(path, entry.Name())

		fileInfos = append(fileInfos, FileInfo{
			Name: entry.Name(),
			Path: fullPath,
			Type: fileType,
			Size: info.Size(),
			Ext:  ext,
		})
	}

	jsonData, err := json.MarshalIndent(fileInfos, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func Init(mcpServer *server.MCPServer) {
	tool := mcp.NewTool("list_macos_dir",
		mcp.WithDescription("list directory contents with file type, size, extension, and path as JSON"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to list directory contents"),
		),
	)
	mcpServer.AddTool(tool, ListDirHandler)

	// Add resource for directory listing
	mcpServer.AddResource(
		mcp.NewResource(
			"macos://dir/{path}",
			"Directory Contents",
			mcp.WithResourceDescription("Directory contents with file information"),
			mcp.WithMIMEType("application/json"),
		),
		handleDirResource,
	)
}

func handleDirResource(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract path from URI (format: macos://dir/{path})
	uri := req.Params.URI
	path := extractPathFromURI(uri)
	
	if path == "" {
		return nil, fmt.Errorf("invalid URI format: %s", uri)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	type FileInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
		Size int64  `json:"size"`
		Ext  string `json:"ext"`
	}

	var fileInfos []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileType := "file"
		if entry.IsDir() {
			fileType = "directory"
		}

		ext := filepath.Ext(entry.Name())
		fullPath := filepath.Join(path, entry.Name())

		fileInfos = append(fileInfos, FileInfo{
			Name: entry.Name(),
			Path: fullPath,
			Type: fileType,
			Size: info.Size(),
			Ext:  ext,
		})
	}

	jsonData, err := json.MarshalIndent(fileInfos, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func extractPathFromURI(uri string) string {
	// Extract path from URI like "macos://dir/Users/jun/Documents"
	prefix := "macos://dir/"
	if len(uri) > len(prefix) && uri[:len(prefix)] == prefix {
		return uri[len(prefix):]
	}
	return ""
}
