package main

import (
	"Go-EventService/config"
	"Go-EventService/services"
	"Go-EventService/utils"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	startBlockBigInt := new(big.Int).SetUint64(uint64(cfg.STARTBLOCK))

	eventService, err := services.NewEventService(cfg.RPCURL, cfg.ContractAddress, "StakingContract", big.NewInt(42), utils.STAKINGABIJSON, startBlockBigInt)
	if err != nil {
		log.Fatal("Error in creating staking events service", eventService)
	}
	tokenService, err := services.NewEventService(cfg.RPCURL, cfg.TokenContract, "TokenContract", big.NewInt(43), utils.TOKENABIJSON, big.NewInt(20417192))
	if err != nil {
		log.Fatal("Error in creating token events service", tokenService)

	}
	c := cron.New()
	_, err = c.AddFunc(cfg.CRONDURATION, func() {
		fmt.Println("Cron Running")
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := eventService.FetchAndProcessEvents(); err != nil {
				log.Printf("Error fetching staking service events: %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			if err := tokenService.FetchAndProcessEvents(); err != nil {
				log.Printf("Error fetching token service events: %v", err)
			}
		}()
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}
	c.Start()
	select {}
}
