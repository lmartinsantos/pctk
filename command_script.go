package pctk

import "log"

// ScriptRun is a command to run a script.
type ScriptRun struct {
	ScriptRef ResourceRef
}

func (c ScriptRun) Execute(app *App, prom *Promise) {
	script, ok := app.scripts[c.ScriptRef]
	if ok {
		prom.CompleteWithValue(script)
		return
	}

	// The script was not loaded yet, so we load and execute it now.
	script = app.res.LoadScript(c.ScriptRef)
	if script == nil {
		log.Panicf("Script not found: %s", c.ScriptRef)
	}
	script.init(app, c.ScriptRef)
	app.scripts[c.ScriptRef] = script
	script.run(app, prom)
}

// ScriptCall is a command to call a script function.
type ScriptCall struct {
	ScriptRef ResourceRef
	Function  FieldAccessor
}

func (c ScriptCall) Execute(app *App, prom *Promise) {
	script, ok := app.scripts[c.ScriptRef]
	if !ok {
		log.Panicf("Script not found: %s", c.ScriptRef)
	}

	prom.Bind(script.Call(c.Function, nil, false))
}
