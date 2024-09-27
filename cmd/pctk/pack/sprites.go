package pack

import (
	"path/filepath"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

type SpriteSheetData struct {
	Resource *pctk.SpriteSheet

	workingDir string
}

func NewSpriteSheetData(workingDir string) *SpriteSheetData {
	return &SpriteSheetData{workingDir: workingDir}
}

func (d *SpriteSheetData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Frames struct {
			Width  uint
			Height uint
		}
		Source string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	d.Resource = pctk.LoadSpriteSheetFromFile(
		filepath.Join(d.workingDir, data.Source),
		pctk.Size{W: int(data.Frames.Width), H: int(data.Frames.Height)},
	)
	return nil
}
