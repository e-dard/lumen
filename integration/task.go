package integration

import (
	"encoding/json"
	"io/ioutil"
)

type Runnable interface {
	Run() error
}

type SublimeTask struct {
	prefsLocation string
	preferences   []byte
}

func NewSublimeTask(loc string, prefs []byte) *SublimeTask {
	return &SublimeTask{
		prefsLocation: loc,
		preferences:   prefs,
	}
}

func (t *SublimeTask) Run() error {
	// read existing file contents
	current, err := ioutil.ReadFile(t.prefsLocation)
	if err != nil {
		return err
	}

	// patch the JSON
	next, err := PatchJSON(current, t.preferences)
	if err != nil {
		return err
	}

	// save the file
	return ioutil.WriteFile(t.prefsLocation, next, 0644)
}

func PatchJSON(br []byte, pr []byte) ([]byte, error) {
	var (
		base  map[string]interface{}
		patch map[string]interface{}
	)

	if err := json.Unmarshal(br, &base); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(pr, &patch); err != nil {
		return nil, err
	}

	for k, v := range patch {
		base[k] = v
	}

	return json.MarshalIndent(base, "", "  ")
}

// type Preferences struct {
// 	ColorScheme string   `json:"color_scheme"`
// 	FontFace    string   `json:"font_face"`
// 	FontSize    float64  `json:"font_size"`
// 	FontOptions []string `json:"font_options"`
// }
