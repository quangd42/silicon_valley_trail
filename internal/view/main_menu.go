package view

import "github.com/quangd42/silicon_valley_trail/internal/model"

func MainMenu() PromptView {
	return PromptView{
		Controls: []model.Control{
			model.ControlNewGame,
			model.ControlLoad,
			model.ControlQuitGame,
		},
	}
}
