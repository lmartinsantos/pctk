package pack

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

// CostumeData is the data for a costume resource.
type CostumeData struct {
	Resource *pctk.Costume

	workingDir string
}

// NewCostumeData creates a new costume data associated with a working directory.
func NewCostumeData(workingDir string) *CostumeData {
	return &CostumeData{workingDir: workingDir}
}

func (d *CostumeData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Sprites struct {
			Sheet  string
			Width  uint
			Height uint
		}
		Animations []struct {
			Action string
			Dir    string
			Flip   bool
			Frames []struct {
				Row      uint
				Columns  []uint
				Duration uint
			}
		}
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	sprites := pctk.LoadSpriteSheetFromFile(
		filepath.Join(d.workingDir, data.Sprites.Sheet),
		pctk.Size{W: int(data.Sprites.Width), H: int(data.Sprites.Height)},
	)
	d.Resource = pctk.NewCostume(sprites)

	for _, anim := range data.Animations {
		dir := func() pctk.Direction {
			switch strings.ToLower(anim.Dir) {
			case "right":
				return pctk.DirRight
			case "left":
				return pctk.DirLeft
			case "up":
				return pctk.DirUp
			case "down":
				return pctk.DirDown
			default:
				panic(fmt.Sprintf("invalid direction %q", anim.Dir))
			}
		}

		a := pctk.NewAnimation().Flip(anim.Flip)
		for _, frame := range anim.Frames {
			a.WithFrames(
				frame.Row,
				time.Duration(frame.Duration)*time.Millisecond,
				frame.Columns...,
			)
		}

		var act pctk.CostumeAction
		switch strings.ToLower(anim.Action) {
		case "idle":
			act = pctk.CostumeIdle(dir())
		case "speak":
			act = pctk.CostumeSpeak(dir())
		case "walk":
			act = pctk.CostumeWalk(dir())
		default:
			code, err := strconv.Atoi(anim.Action)
			if err != nil {
				err := fmt.Errorf("neither a default action nor a custom action code: %w", err)
				return fmt.Errorf("invalid action %q: %w", anim.Action, err)
			}
			act = pctk.CostumeAction(code)
		}
		d.Resource.WithAnimation(act, a)
	}

	return nil
}
