package save

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

const DefaultSavePath = "svt_save.json"

var ErrSaveCorrupted = errors.New("save file is corrupted")

// JSONSaver persists game state as JSON on the local filesystem.
//
// To avoid corrupting the save file, it writes to a temp file in the target
// directory and then swap the temp file with the save file once writing succeeded.
// This strategy is intended for Unix-like systems and might not work on Windows.
//
// NOTE: Internal filesystem errors are being exposed to the player, which is fine for the
// current stage.
type JSONSaver struct {
	path string
}

func NewJSONSaver(path string) *JSONSaver {
	if path == "" {
		path = DefaultSavePath
	}
	return &JSONSaver{
		path: path,
	}
}

func (js *JSONSaver) Save(state *model.State) error {
	tmp, err := os.CreateTemp(filepath.Dir(js.path), "svt_save_*.json")
	if err != nil {
		return err
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	err = json.NewEncoder(tmp).Encode(state)
	if err != nil {
		return err
	}

	// Sync commits file to disk
	err = tmp.Sync()
	if err != nil {
		return err
	}
	// Overwrite save file
	return os.Rename(tmp.Name(), js.path)
}

func (js *JSONSaver) Load(state *model.State) error {
	file, err := os.Open(js.path)
	if err != nil {
		return err
	}
	defer file.Close()
	var loaded model.State
	err = json.NewDecoder(file).Decode(&loaded)
	if err != nil {
		return ErrSaveCorrupted
	}
	*state = loaded
	return nil
}
