package main

import (
	"context"
	"fmt"
	"log"

	"github.com/romeq/logfront/internal/consumers/ntfy_sh"
	"github.com/romeq/logfront/internal/domain"
	"github.com/romeq/logfront/internal/pipeline"
	"github.com/romeq/logfront/internal/sources"
	"github.com/romeq/logfront/internal/sources/ftp"
	"github.com/romeq/logfront/internal/sources/lastenareena"
	"github.com/romeq/logfront/internal/sources/ssh"
)

func main() {
	args := parseArguments()

	config, err := parseConfig(args.configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// register known sources
	sources.Register(ssh.ConfigName, ssh.NewSource)
	sources.Register(ftp.ConfigName, ftp.NewSource)

	// check enabled sources and consumers and create the actual interfaces from register
	configSources, err := loadSourcesFromConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	configConsumers, err := loadConsumersFromConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	events := map[string]chan domain.FailedLoginEvent{}
	for _, configConsumer := range configConsumers {
		events[configConsumer.Name()] = make(chan domain.FailedLoginEvent)
	}

	ctx := context.TODO()
	dispatcher := pipeline.NewDispatcher(configConsumers)
	startSourceListener := func(s domain.Source) {
		if err := s.Start(ctx, events); err != nil {
			log.Println(err)
		}
	}

	for _, s := range configSources {
		go startSourceListener(s)
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
	for key, value := range config.Services {
		if key == "ntfy_sh" {
			firstMatch, ok := value["urls"].([]interface{})[0].(string)
			if !ok {
				return nil, fmt.Errorf("invalid ntfysh urls format")
			}
			configConsumers = append(configConsumers, ntfy_sh.NtfyshConsumer{Webhook: firstMatch})
			log.Println("Loaded ntfysh consumer (url", firstMatch+")")
		}
	}
	// at least 1 source must be configured
	if len(configConsumers) == 0 {
		return nil, fmt.Errorf("no sources configured")
	}
	return configConsumers, nil
}
