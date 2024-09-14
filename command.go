package pctk

// Command is a command that can be executed by the application. Every action requires to the
// application that has a side effect should be encapsulated in a command in order to ensure
// thread safety.
type Command interface {
	Execute(*App, Promise)
}

// Do will put the given command in the queue to be executed by the application during the next
// frame.
func (a *App) Do(c Command) Future {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.commands.push(c)

}

type commandQueue struct {
	commands []Command
	promises []Promise
}

func (q *commandQueue) push(c Command) Future {
	done := NewPromise()
	q.commands = append(q.commands, c)
	q.promises = append(q.promises, done)
	return done
}

func (q *commandQueue) execute(a *App) {
	commands := q.commands
	promises := q.promises
	q.commands = nil
	q.promises = nil
	for i, c := range commands {
		c.Execute(a, promises[i])
	}
}
