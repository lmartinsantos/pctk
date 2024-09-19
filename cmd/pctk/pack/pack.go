package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apoloval/pctk"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func do(src string, output string) error {
	rl.SetTraceLogLevel(rl.LogNone)

	idxFile, datFile, err := createOutputFiles(output)
	if err != nil {
		return err
	}
	defer idxFile.Close()
	defer datFile.Close()

	enc, err := pctk.NewResourceEncoder(idxFile, datFile)
	if err != nil {
		return err
	}

	manifests, err := listManifests(src)
	if err != nil {
		return err
	}
	for _, manifest := range manifests {
		ref := pctk.ResourceRef(strings.TrimSuffix(manifest, filepath.Ext(manifest)))
		fmt.Printf("Packing %s...", ref)
		man, err := LoadManifestFromFile(filepath.Join(src, manifest))
		if err != nil {
			return err
		}

		switch data := man.Data.(type) {
		case *CostumeData:
			enc.EncodeCostume(ref, data.Resource, man.Compression)
		case *MusicData:
			enc.EncodeMusic(ref, data.Resource, man.Compression)
		case *RoomData:
			enc.EncodeRoom(ref, data.AsResource(), man.Compression)
		case *ScriptData:
			enc.EncodeScript(ref, data.AsResource(), man.Compression)
		case *SoundData:
			enc.EncodeSound(ref, data.Resource, man.Compression)
		}
		fmt.Printf(" Done\n")
	}

	fmt.Printf("%d data bytes written\n", enc.DataBytesWritten())
	return nil
}

func createOutputFiles(output string) (*os.File, *os.File, error) {
	idx, err := os.Create(output + ".idx")
	if err != nil {
		return nil, nil, err
	}
	dat, err := os.Create(output + ".dat")
	if err != nil {
		return nil, nil, err
	}
	return idx, dat, nil
}

func listManifests(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var manifests []string
	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			others, err := listManifests(filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			for i := range others {
				others[i] = filepath.Join(name, others[i])
			}
			manifests = append(manifests, others...)
			continue
		}

		if len(name) > 5 && name[len(name)-4:] == ".yml" {
			manifests = append(manifests, name)
		}
		if len(name) > 6 && name[len(name)-5:] == ".yaml" {
			manifests = append(manifests, name)
		}
	}
	return manifests, nil
}
