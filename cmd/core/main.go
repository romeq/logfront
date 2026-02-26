package main

import (
	"context"
	"fmt"
	"log"

	"github.com/romeq/logfront/internal/consumers"
	"github.com/romeq/logfront/internal/domain"
	"github.com/romeq/logfront/internal/pipeline"
	"github.com/romeq/logfront/internal/sources"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	args := parseArguments()
	config, err := parseConfig(args.configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// known.go - remember to add all source & consumer implementations here!
	initRegisters()

	// check enabled sources and consumers and create the actual interfaces from register
	configSources, err := loadSourcesFromConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	configConsumers, err := loadConsumersFromConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	// create event channels for each consumer
	events := map[string]chan domain.LogEvent{}
	for _, configConsumer := range configConsumers {
		events[configConsumer.Name()] = make(chan domain.LogEvent)
	}

	// wire things up
	ctx := context.Background()
	dispatcher := pipeline.NewDispatcher(configConsumers)
	for _, s := range configSources {
		go func() {
			if err := s.Start(ctx, events); err != nil {
				log.Println(s.Name(), "returned error:", err) //
			}
		}()
	}
	for {
		dispatcher.Run(ctx, events)
	}
}

func loadSourcesFromConfig(config AppConfig) ([]domain.Source, error) {
	var configSources []domain.Source
	for configSource, configSourceConfig := range config.Sources {
		keyAppropriateSource, sourceFound := sources.Create(configSource, configSourceConfig)
		if !sourceFound {
			return nil, fmt.Errorf("failed to find source implementation with name: %s", configSource)
		}

		configSources = append(configSources, keyAppropriateSource)
	}

	// at least 1 source must be configured
	if len(configSources) == 0 {
		return nil, fmt.Errorf("no sources configured")
	}

	return configSources, nil
}

func loadConsumersFromConfig(config AppConfig) ([]domain.Consumer, error) {
	var configConsumers []domain.Consumer
	for configConsumer, configConsumerConfig := range config.Consumers {
		keyAppropriateConsumer, consumerFound := consumers.Create(configConsumer, configConsumerConfig)
		if !consumerFound {
			return nil, fmt.Errorf("failed to find consumer implementation with name: %s", configConsumer)
		}

		configConsumers = append(configConsumers, keyAppropriateConsumer)
	}

	// at least 1 consumer must be configured
	if len(configConsumers) == 0 {
		return nil, fmt.Errorf("no consumers configured")
	}
	return configConsumers, nil
}
