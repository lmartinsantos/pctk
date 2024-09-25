package pack

import (
	"path/filepath"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

// ImageData is the data for a room resource.
type ImageData struct {
	Resource *pctk.Image

	workingDir string
}

// NewImageData creates a new image data associated with a working directory.
func NewImageData(workingDir string) *ImageData {
	return &ImageData{workingDir: workingDir}
}

func (d *ImageData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Source string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	d.Resource = pctk.LoadImageFromFile(filepath.Join(d.workingDir, data.Source))

	return nil
}
