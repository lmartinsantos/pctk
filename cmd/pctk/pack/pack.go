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
		id := pctk.ResourceID(strings.TrimSuffix(manifest, filepath.Ext(manifest)))
		fmt.Printf("Packing %s...", id)
		man, err := LoadManifestFromFile(filepath.Join(src, manifest))
		if err != nil {
			return err
		}

		switch data := man.Data.(type) {
		case *CostumeData:
			err = enc.EncodeCostume(id, data.Resource, man.Compression)
		case *ImageData:
			err = enc.EncodeImage(id, data.Resource, man.Compression)
		case *MusicData:
			err = enc.EncodeMusic(id, data.Resource, man.Compression)
		case *ScriptData:
			err = enc.EncodeScript(id, data.Resource, man.Compression)
		case *SoundData:
			err = enc.EncodeSound(id, data.Resource, man.Compression)
		}
		if err != nil {
			fmt.Printf(" Failed!\n")
			return err
		}
		fmt.Printf(" Done\n")
	}

	luaScripts, err := listLuaScripts(src)
	if err != nil {
		return err
	}
	for _, script := range luaScripts {
		id := pctk.ResourceID(strings.TrimSuffix(script, filepath.Ext(script)))
		fmt.Printf("Packing %s...", id)
		code, err := os.ReadFile(filepath.Join(src, script))
		if err != nil {
			return err
		}

		script := pctk.NewScript(pctk.ScriptLua, code)
		err = enc.EncodeScript(id, script, pctk.CompressionNone)
		if err != nil {
			fmt.Printf(" Failed!\n")
			return err
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
	return listFiles(dir, ".yml", ".yaml")
}

func listLuaScripts(dir string) ([]string, error) {
	return listFiles(dir, ".lua")
}

func listFiles(dir string, extensions ...string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			others, err := listFiles(filepath.Join(dir, name), extensions...)
			if err != nil {
				return nil, err
			}
			for i := range others {
				others[i] = filepath.Join(name, others[i])
			}
			files = append(files, others...)
			continue
		}

		for _, ext := range extensions {
			if len(name) > len(ext) && name[len(name)-len(ext):] == ext {
				files = append(files, name)
				break
			}
		}
	}
	return files, nil
}
