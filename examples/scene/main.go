package main

import (
	"github.com/apoloval/pctk"
)

func main() {
	loader := pctk.NewResourceFileLoader("./")

	app := pctk.New(loader)
	app.Do(pctk.ScriptRun{ScriptResource: pctk.NewResourceRef("resources", "scripts/LostMyKeys")})
	app.Run()
}
