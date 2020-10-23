// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"bytes"
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"io"
	"os"
	"strings"
	"sync"
)

var mgr = &shellContextManager{}

type shellContextManager struct {
	ctx *shellContext
	mu  sync.RWMutex
}

func (m *shellContextManager) getContext() *shellContext {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ctx
}

func (m *shellContextManager) setContext(context *shellContext) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ctx = context
}

func getContext() *shellContext {
	return mgr.getContext()
}

func setContext(context *shellContext) {
	mgr.setContext(context)
}

// newContext creates a new shell context
func newContext(names ...string) *shellContext {
	name := strings.Join(names, ":")
	return &shellContext{
		name:   name,
		args:   []string{},
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// shellContext is a shell context
type shellContext struct {
	name    string
	args    []string
	stdin   io.ReadCloser
	stdout  io.Writer
	stderr  io.Writer
	context *ishell.Context
	parent  *shellContext
}

func (c *shellContext) withCommand(names ...string) *shellContext {
	name := strings.Join(append([]string{c.name}, names...), ":")
	args := append(c.args, names...)
	return &shellContext{
		name:   name,
		args:   args,
		stdin:  c.stdin,
		stdout: c.stdout,
		stderr: c.stderr,
		parent: c,
	}
}

func (c *shellContext) withContext(context *ishell.Context) *shellContext {
	return &shellContext{
		name:    c.name,
		args:    c.args,
		stdin:   c.stdin,
		stdout:  c.stdout,
		stderr:  c.stderr,
		context: context,
		parent:  c,
	}
}

func (c *shellContext) run(flags ...string) error {
	args := append(c.args, flags...)
	historyFile, err := getConfigFile("history.log")
	if err != nil {
		return err
	}

	shell := ishell.NewWithConfig(&readline.Config{
		Prompt:      fmt.Sprintf("%s> ", c.name),
		HistoryFile: historyFile,
		Stdin:       c.stdin,
		Stdout:      c.stdout,
		Stderr:      c.stderr,
	})

	shell.NotFound(func(context *ishell.Context) {
		setContext(c.withContext(context))
		args := append(args, context.RawArgs...)
		cmd := GetRootCommand()
		cmd.SetArgs(args)
		err := cmd.Execute()
		if err != nil {
			context.Println(err)
		}
		setContext(c)
	})

	shell.Interrupt(func(context *ishell.Context, count int, input string) {
		if c.parent != nil {
			context.Stop()
		} else if count >= 2 {
			context.Println("Interrupted")
			os.Exit(1)
		} else {
			context.Println("Input Ctrl-c once more to exit")
		}
	})

	shell.EOF(func(context *ishell.Context) {
		context.Stop()
	})

	shell.DeleteCmd("help")

	setContext(c)
	shell.Run()
	setContext(c.parent)
	return nil
}

func (c *shellContext) Print(values ...interface{}) {
	fmt.Fprint(c.stdout, values...)
}

func (c *shellContext) Println(values ...interface{}) {
	fmt.Fprintln(c.stdout, values...)
}

func (c *shellContext) Printlns(lines ...interface{}) {
	if c.context != nil {
		buf := bytes.Buffer{}
		for _, line := range lines {
			buf.WriteString(fmt.Sprint(line))
			buf.WriteByte('\n')
		}
		c.context.ShowPaged(buf.String())
	} else {
		for _, line := range lines {
			c.Println(line)
		}
	}
}

func (c *shellContext) Printf(format string, args ...interface{}) {
	fmt.Fprintf(c.stdout, format, args...)
}
