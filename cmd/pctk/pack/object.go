package pack

import (
	"path/filepath"
	"time"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

// ObjectData is the data for an object resource.
type ObjectData struct {
	Resource *pctk.Object

	workingDir string
}

// NewObjectData creates a new object data associated with a working directory.
func NewObjectData(workingDir string) *ObjectData {
	return &ObjectData{workingDir: workingDir}
}

func (d *ObjectData) UnmarshalYAML(n *yaml.Node) error {
	// TODO
	var data struct {
		Name    string
		Sprites struct {
			Sheet  string
			Width  uint
			Height uint
		}
		Animation struct {
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
	d.Resource = pctk.NewObject(data.Name, sprites)
	a := pctk.NewAnimation().Flip(data.Animation.Flip)
	for _, frame := range data.Animation.Frames {
		a.WithFrames(
			frame.Row,
			time.Duration(frame.Duration)*time.Millisecond,
			frame.Columns...,
		)
	}
	d.Resource.WithAnimation(a)

	return nil
}
