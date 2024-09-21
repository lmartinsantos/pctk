package main

import (
	"github.com/apoloval/pctk"
)

func main() {
	loader := pctk.NewResourceFileLoader("./")

	app := pctk.New(loader)
	app.Do(pctk.ScriptRun{ScriptResource: pctk.NewResourceRef("resources", "scripts/LostMyKeys")})
	app.Run()
	/*

		objects_sprites := pctk.LoadSpriteSheetFromFile("items.png", pctk.Size{W: 30, H: 18})
			bundle.PutSpriteSheet("/objects/sprites", objects_sprites)

			bucket := pctk.NewObject("bucket").WithAnimation(pctk.NewAnimation("/objects/sprites").
				WithFrame(5, 6, 0*time.Second)).WithScript(pctk.LookAt, &pctk.Script{
				Language: pctk.ScriptLua,
				Code: []byte(`
					ActorSpeak("guybrush", "Mmmm, nice bucket!", {delay=1000}).Wait()
				`)})
			bundle.PutObject("/objects/bucket", bucket)
	*/
}
