package pack

import (
	"fmt"
	"path/filepath"
	"strings"
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
	var data struct {
		Name     string
		Position struct {
			X uint
			Y uint
			Z uint // TODO for z-index ordering
		}
		Sprites struct {
			Sheet  string
			Width  uint
			Height uint
		}
		States []struct {
			Animation *struct {
				Action string
				Dir    string
				Flip   bool
				Frames []struct {
					Row      uint
					Columns  []uint
					Duration uint
				}
			}
			Scripts []struct {
				Verb     string
				Language string
				Code     string
			}
		}
		Classes uint
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	sprites := pctk.LoadSpriteSheetFromFile(
		filepath.Join(d.workingDir, data.Sprites.Sheet),
		pctk.Size{W: int(data.Sprites.Width), H: int(data.Sprites.Height)},
	)
	// TODO Z
	position := pctk.NewPos(int(data.Position.X), int(data.Position.Y))
	d.Resource = pctk.NewObject(data.Name, sprites, position, data.Classes)

	for _, state := range data.States {
		s := pctk.NewState()
		if state.Animation != nil {
			a := pctk.NewAnimation().Flip(state.Animation.Flip)
			for _, frame := range state.Animation.Frames {
				a.WithFrames(
					frame.Row,
					time.Duration(frame.Duration)*time.Millisecond,
					frame.Columns...,
				)
			}

			s.WithAnimation(a)

		}

		for _, script := range state.Scripts {
			verb := func() pctk.VerbType {
				switch strings.ToLower(script.Verb) {
				case "open":
					return pctk.Open
				case "close":
					return pctk.Close
				case "push":
					return pctk.Push
				case "pull":
					return pctk.Pull
				case "walkto":
					return pctk.WalkTo
				case "pickup":
					return pctk.PickUp
				case "talkto":
					return pctk.TalkTo
				case "give":
					return pctk.Give
				case "use":
					return pctk.Use
				case "lookat":
					return pctk.LookAt
				case "turnon":
					return pctk.TurnOn
				case "turnoff":
					return pctk.TurnOff
				case "default":
					return pctk.Default
				default:
					panic(fmt.Sprintf("invalid verb %q", script.Verb))
				}
			}
			sc := pctk.NewScript(pctk.ScriptLua, []byte(script.Code))
			s.WithScript(verb(), sc)
		}
		d.Resource.WithState(s)
	}

	return nil
}
