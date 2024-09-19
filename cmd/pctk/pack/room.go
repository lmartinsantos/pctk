package pack

import (
	"path/filepath"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

// RoomData is the data for a room resource.
type RoomData struct {
	Resource *pctk.Room

	workingDir string
}

// NewRoomData creates a new room data associated with a working directory.
func NewRoomData(workingDir string) *RoomData {
	return &RoomData{workingDir: workingDir}
}

func (d *RoomData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Background string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	d.Resource = pctk.NewRoom(pctk.LoadImageFromFile(filepath.Join(d.workingDir, data.Background)))

	return nil
}
