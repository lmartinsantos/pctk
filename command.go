package pctk

import "sync"

// Command is a command that can be executed by the application. Every action requires to the
// application that has a side effect should be encapsulated in a command in order to ensure
// thread safety.
type Command interface {
	Execute(*App, *Promise)
}

// AppContext is an interface that defines the method to execute a command.
type AppContext interface {
	Do(Command) Future
}

// Do will put the given command in the queue to be executed by the application during the next
// frame. This function must not be called from a command handler, as it will cause a deadlock.
// If one command have to execute another command, use doNow function instead.
func (a *App) Do(c Command) Future {
	return a.commands.push(c)
}

func (a *App) doNow(c Command) Future {
	done := NewPromise()
	c.Execute(a, done)
	return done
}

type commandQueue struct {
	mutex    sync.Mutex
	commands []Command
	promises []*Promise
}

func (q *commandQueue) push(c Command) Future {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	done := NewPromise()
	q.commands = append(q.commands, c)
	q.promises = append(q.promises, done)
	return done
}

func (q *commandQueue) execute(a *App) {
	q.mutex.Lock()
	commands := q.commands
	promises := q.promises
	q.commands = nil
	q.promises = nil
	q.mutex.Unlock()

	for i, c := range commands {
		c.Execute(a, promises[i])
	}
}
