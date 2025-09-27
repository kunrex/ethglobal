package main

import (
	"git-server/pkg/abi"
	"git-server/pkg/eth"
	"git-server/pkg/routes"
	"git-server/pkg/server"
	"git-server/pkg/types"
	"git-server/pkg/utils"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethRPC         = "https://testnet.infura.io/v3/YOUR_KEY" // Ethereum endpoint
	chainID        = big.NewInt(11155111)                    // Example chainId (Sepolia...)
	contractGlobal *abi.Abi                                  // Global contract instance
)

var ProjectWallets = map[string]*types.AnonymousWallet{}

func main() {
	// Load project wallets from file
	var err error
	ProjectWallets, err = utils.LoadProjectsFromFile()
	if err != nil {
		log.Fatalf("Failed to load project wallets: %v", err)
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial(ethRPC)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}
	// Bind the global contract instance ONCE for the whole server
	// (You might want to use a “master” wallet here for contract writes if needed)
	masterWallet := ProjectWallets["master"]
	if masterWallet == nil {
		masterWallet, err = eth.CreateWallet() // Used ONLY for setting up the contract if needed
		if err != nil {
			log.Fatalf("Error creating master wallet: %v", err)
		}
		ProjectWallets["master"] = masterWallet
		if err := utils.SaveProjectsToFile(ProjectWallets); err != nil {
			log.Fatalf("Error saving master wallet: %v", err)
		}
	}

	contract, err := eth.CreateContract(masterWallet, client)
	if err != nil {
		log.Fatalf("Error binding contract: %v", err)
	}
	contractGlobal = contract
	gitServer := server.NewInMemoryGitServer()

	router := routes.SetupRoutes(gitServer)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3370"
	}
	log.Printf("Starting Git server on port %s", port)
	log.Printf("Visit http://localhost:%s to see the web interface", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
