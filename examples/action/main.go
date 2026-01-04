package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mirkobrombin/go-module-router/v2/pkg/core"
	"github.com/mirkobrombin/go-module-router/v2/pkg/logger"
	"github.com/mirkobrombin/go-module-router/v2/pkg/transport/action"
)

// --- Domain ---

type Document struct {
	Name    string
	Content string
	Dirty   bool
}

func (d *Document) Save() error {
	d.Dirty = false
	return nil
}

// --- Actions ---

type SaveAction struct {
	Meta     core.Pattern `action:"file.save" keys:"ctrl+s"`
	Document *Document
}

func (a *SaveAction) Handle(ctx context.Context) (any, error) {
	if err := a.Document.Save(); err != nil {
		return nil, err
	}
	return map[string]string{"status": "saved", "file": a.Document.Name}, nil
}

type NewFileAction struct {
	Meta     core.Pattern `action:"file.new" keys:"ctrl+n"`
	Document *Document
}

func (a *NewFileAction) Handle(ctx context.Context) (any, error) {
	a.Document.Name = "Untitled"
	a.Document.Content = ""
	a.Document.Dirty = false
	return map[string]string{"status": "new file created"}, nil
}

// --- Main ---

func main() {
	// Create document
	doc := &Document{Name: "test.txt", Content: "Hello World", Dirty: true}

	// Create action transport
	t := action.New()
	t.Logger = logger.NewSlog(slog.Default())

	// Register dependencies
	t.Provide("Document", doc)

	// Register actions
	t.Register(&SaveAction{})
	t.Register(&NewFileAction{})

	// Show registered actions
	fmt.Println("Registered actions:", t.Actions())
	fmt.Println("Key bindings:", t.KeyBindings())

	// Simulate user pressing Ctrl+S
	fmt.Println("\n--- Dispatching 'ctrl+s' ---")
	result, err := t.DispatchKey(context.Background(), "ctrl+s")
	if err != nil {
		slog.Error("dispatch failed", "error", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	// Simulate dispatching by action name
	fmt.Println("\n--- Dispatching 'file.new' ---")
	result, err = t.Dispatch(context.Background(), "file.new")
	if err != nil {
		slog.Error("dispatch failed", "error", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}
}
