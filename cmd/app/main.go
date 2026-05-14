package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/Mkz-Prog/yadro-telecom-test/internal/config"
	"github.com/Mkz-Prog/yadro-telecom-test/internal/parser"
	"github.com/Mkz-Prog/yadro-telecom-test/internal/processor"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	eventsPath := flag.String("events", "events.txt", "path to events log")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}

	engine := processor.NewEngine(cfg)

	eventsFile, err := os.Open(*eventsPath)
	if err != nil {
		log.Fatalf("Fatal error opening events file: %v", err)
	}
	defer func() {
		if err := eventsFile.Close(); err != nil {
			log.Printf("Warning: failed to close events file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(eventsFile)
	for scanner.Scan() {
		line := scanner.Text()

		event, err := parser.ParseLine(line)
		if err != nil {
			log.Printf("Warning: failed to parse line: %v", err)
			continue
		}
		if event == nil {
			continue
		}

		engine.ProcessEvent(event)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading events file: %v", err)
	}

	engine.PrintFinalReport()
}
