package main

import (
	"errors"
	"ethglobal/pkg/config"
	"ethglobal/pkg/contract"
	"ethglobal/pkg/controllers"
	"ethglobal/pkg/lighthouse"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	configuration := config.LoadConfig()

	actions, rootCtx, err := contract.InitContractActions(&configuration)
	if err != nil {
		return
	}

	lighthouseClient := lighthouse.InitLightHouseClient(configuration)

	controller := controllers.Controller{
		ActionContracts: actions,
		Lighthouse:      lighthouseClient,
	}

	var push = &cobra.Command{
		Use:   "push",
		Short: "Push the latest git history on cold storage",
		Long:  "Push the latest git history on cold storage on Lighthouse",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New(fmt.Sprintf("expected 2 arguments, got %d", len(args)))
			}

			storage, err := controller.PushColdStorage(args[0], args[1])
			if err != nil {
				return err
			}

			(*rootCtx).Done()
			log.Printf(storage)
			return nil
		},
	}

	var pull = &cobra.Command{
		Use:   "pull",
		Short: "Pull the latest git history from cold storage",
		Long:  "Pull the latest git history from cold storage on Lighthouse",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New(fmt.Sprintf("expected 2 arguments, got %d", len(args)))
			}

			if err := controller.RetrieveColdStorage(args[0], args[1]); err != nil {
				return err
			}

			(*rootCtx).Done()
			return nil
		},
	}

	var root = &cobra.Command{
		Use:   "ccg",
		Short: "ccg",
	}

	root.AddCommand(push)
	root.AddCommand(pull)

	_ = root.Execute()
}
