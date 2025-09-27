package main

import (
	"encoding/json"
	"errors"
	"ethglobal/pkg/config"
	"ethglobal/pkg/contract"
	"ethglobal/pkg/controllers"
	"ethglobal/pkg/lighthouse"
	"ethglobal/pkg/types"
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
		ActionContracts:    actions,
		Lighthouse:         lighthouseClient,
		EncryptionKeyBytes: []byte(configuration.EncryptionKey),
	}

	var push = &cobra.Command{
		Use:   "push",
		Short: "push [repository identifier] [path/to/.git.zip] [latest commit] -> Transaction Id",
		Long:  "Push the latest git history and metadata on cold storage on Lighthouse, prints Transaction ID",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New(fmt.Sprintf("expected 3 arguments , got %d", len(args)))
			}

			transactionId, err := controller.PushColdStorage(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			(*rootCtx).Done()
			log.Printf(transactionId)
			return nil
		},
	}

	var pull = &cobra.Command{
		Use:   "pull",
		Short: "pull [repository identifier] [path/to/output.git.zip] -> Metadata",
		Long:  "Pull the latest git history from cold storage on Lighthouse, prints Metadata",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New(fmt.Sprintf("expected 2 arguments, got %d", len(args)))
			}

			bytes, err := controller.RetrieveColdStorage(args[0], args[1])
			if err != nil {
				return err
			}

			var versions []types.VersionMetaData
			err = json.Unmarshal(bytes, &versions)
			if err != nil {
				return err
			}

			(*rootCtx).Done()
			log.Printf(string(bytes))
			return nil
		},
	}

	var metadata = &cobra.Command{
		Use:   "metadata",
		Short: "metadata [repository identifier] -> Metadata",
		Long:  "Get metadata of a repository from Lighthouse",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New(fmt.Sprintf("expected 1 arguments, got %d", len(args)))
			}

			bytes, err := controller.RetrieveLatestMetaData(args[0])
			if err != nil {
				return err
			}

			var versions []types.VersionMetaData
			err = json.Unmarshal(bytes, &versions)
			if err != nil {
				return err
			}

			(*rootCtx).Done()
			log.Printf(string(bytes))
			return nil
		},
	}

	var root = &cobra.Command{}

	root.AddCommand(push)
	root.AddCommand(pull)
	root.AddCommand(metadata)

	_ = root.Execute()
}
