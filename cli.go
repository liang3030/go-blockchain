package main

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct{}

func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	// Create a print chain command
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	// Create a wallet command
	creatWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	// Get balance command
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	// Create a blockchain command
	CreateBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	// list addresses command
	ListAddressCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	// Send command
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	// reindexUTXO command
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	// start node command
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	// add -address flag to get balance command
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")

	// add -address flag to get balance command
	createblockchainAddress := CreateBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	// add -from, -to, -amount flags to send command
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")

	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")

	switch os.Args[1] {

	case "createwallet":
		err := creatWalletCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing createwallet command:", err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing printchain command:", err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing printchain command:", err)
		}
	case "createblockchain":
		err := CreateBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing createblockchain command:", err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing send command:", err)
		}
	case "listaddresses":
		err := ListAddressCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing listaddresses command:", err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing reindexutxo command:", err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error in parsing startnode command:", err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if creatWalletCmd.Parsed() {
		cli.CreateWallet(nodeID)
	}

	if ListAddressCmd.Parsed() {
		cli.ListAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}

	if CreateBlockchainCmd.Parsed() {
		if *createblockchainAddress == "" {
			CreateBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.CreateBlockchainCLI(*createblockchainAddress, nodeID)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.GetBalance(*getBalanceAddress, nodeID)
	}

	if reindexUTXOCmd.Parsed() {
		cli.ReindexUTXO(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.Send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
	}

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.StartNode(nodeID, *startNodeMiner)
	}
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid number of arguments")
		cli.PrintUsage()
		os.Exit(1)
	}
}
