package command

import (
	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"sync"
)

// commandContext is an Atomix command context
type commandContext struct {
	isShell  bool
	shellCmd *cobra.Command
	shell    *ishell.Shell
	shellCtx *ishell.Context
}

var manager *contextManager

func init() {
	manager = &contextManager{
		context: &commandContext{},
	}
}

func setContext(context commandContext) {
	manager.setContext(context)
}

func setContextFunc(f func(context *commandContext)) {
	manager.setContextFunc(f)
}

func getContext() commandContext {
	return manager.getContext()
}

type contextManager struct {
	context *commandContext
	mu      sync.RWMutex
}

func (m *contextManager) setContext(context commandContext) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.context = &context
}

func (m *contextManager) setContextFunc(f func(context *commandContext)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	context := *m.context
	f(&context)
	m.context = &context
}

func (m *contextManager) getContext() commandContext {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return *m.context
}
