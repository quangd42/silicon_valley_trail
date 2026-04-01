package main

import (
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/quangd42/silicon_valley_trail/internal/config"
	"github.com/quangd42/silicon_valley_trail/internal/gamedef"
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
	rng := rand.New(rand.NewPCG(
		uint64(time.Now().UnixNano()),
		uint64(time.Now().UnixNano()),
	))
	def := gamedef.Load()
	program.New(
		renderer,
		saver,
		weather,
		rng,
		def,
	).Run()
}
