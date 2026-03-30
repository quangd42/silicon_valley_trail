package save

import (
	"errors"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

var ErrSaveCorrupted = errors.New("save file is corrupted")

type Saver interface {
	Save(*model.State) error
	Load(*model.State) error
}
