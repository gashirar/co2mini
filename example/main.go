package main

import (
	"github.com/gashirar/co2mini"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
)

func main() {
	var co2mini co2mini.Co2mini
	var co2 int
	var temp float64

	zerolog.TimeFieldFormat = ""

	if err := co2mini.Connect(); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	go func() {
		if err := co2mini.Start(); err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}()

	for {
		select {
		case co2 = <-co2mini.Co2Ch:
		case temp = <-co2mini.TempCh:
			log.Log().
				Int("co2", co2).
				Float64("temp", math.Round(temp*10)/10).
				Msg("")
		}
	}
}
