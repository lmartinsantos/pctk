package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

// ResourceType is the type of resource that the manifest describes.
type ResourceType string

const (
	// ManifestTypeCostume is a costume resource.
	ManifestTypeCostume ResourceType = "costume"

	// ManifestTypeMusic is a music resource.
	ManifestTypeMusic ResourceType = "music"

	// ManifestTypeRoom is a room resource.
	ManifestTypeRoom ResourceType = "room"

	// ManifestTypeScript is a script resource.
	ManifestTypeScript ResourceType = "script"

	// ManifestTypeSound is a sound resource.
	ManifestTypeSound ResourceType = "sound"

	// ManifestTypeObject is an object resource.
	ManifestTypeObject ResourceType = "object"
)

// Manifest is the description of a resource.
type Manifest struct {
	Type        ResourceType
	Compression pctk.ResourceCompression
	Data        any

	workingDir string
}

// LoadManifestFromFile loads a manifest from a file.
func LoadManifestFromFile(path string) (*Manifest, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var man Manifest
	man.workingDir = filepath.Dir(path)
	if err := yaml.NewDecoder(data).Decode(&man); err != nil {
		return nil, err
	}
	return &man, nil
}

func (m *Manifest) UnmarshalYAML(n *yaml.Node) error {
	var header struct {
		Type        ResourceType
		Compression string
		Data        yaml.Node
	}
	if err := n.Decode(&header); err != nil {
		return err
	}
	m.Type = header.Type
	switch m.Type {
	case ManifestTypeCostume:
		m.Data = NewCostumeData(m.workingDir)
	case ManifestTypeMusic:
		m.Data = NewMusicData(m.workingDir)
	case ManifestTypeRoom:
		m.Data = NewRoomData(m.workingDir)
	case ManifestTypeScript:
		m.Data = new(ScriptData)
	case ManifestTypeSound:
		m.Data = NewSoundData(m.workingDir)
	case ManifestTypeObject:
		m.Data = NewObjectData(m.workingDir)
	default:
		return fmt.Errorf("unknown manifest type: %s", m.Type)
	}
	if err := header.Data.Decode(m.Data); err != nil {
		return err
	}

	switch strings.ToLower(header.Compression) {
	case "", "none":
		m.Compression = pctk.CompressionNone
	case "gzip":
		m.Compression = pctk.CompressionGzip
	default:
		return fmt.Errorf("unknown compression type: %s", header.Compression)
	}

	return nil
}
