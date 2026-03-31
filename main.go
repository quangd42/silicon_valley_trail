package main

import (
	"log"
	"os"

	"github.com/quangd42/silicon_valley_trail/internal/config"
	"github.com/quangd42/silicon_valley_trail/internal/content"
	"github.com/quangd42/silicon_valley_trail/internal/program"
	"github.com/quangd42/silicon_valley_trail/internal/save"
	"github.com/quangd42/silicon_valley_trail/internal/ui"
	"github.com/quangd42/silicon_valley_trail/internal/weather"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	renderer := ui.NewTerminal(os.Stdin, os.Stdout)
	saver := save.NewJSONSaver(cfg.SavePath)
	weather := weather.NewWeatherService(cfg.WeatherAPIKey, cfg.WeatherAPIUseMock, cfg.WeatherAPITimeout)
	gameCopy := content.Load()
	program.Run(
		renderer,
		saver,
		weather,
		gameCopy,
	)
}
