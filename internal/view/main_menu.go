package view

import "github.com/quangd42/silicon_valley_trail/internal/model"

func MainMenuPrompt() PromptView {
	return PromptView{
		Title: "SILICON VALLEY TRAIL - Main Menu",
		Sections: []PromptSectionView{
			{
				Items: controlItems([]model.Control{
					model.ControlNewGame,
					model.ControlLoad,
					model.ControlQuitGame,
				}),
			},
		},
	}
}
