package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bilal2134/Blockchain_A3/internal/bft"
	"github.com/bilal2134/Blockchain_A3/internal/blockchain"
	"github.com/bilal2134/Blockchain_A3/internal/cap"
	"github.com/bilal2134/Blockchain_A3/internal/consensus"
)

// Entry point for the blockchain node
func main() {
	// Initialize the blockchain
	bc := NewBlockchain()
	// Initialize authentication and reputation
	authMgr := consensus.NewAuthManager()
	repSys := bft.NewReputationSystem()
	// Initialize CAP orchestrator components
	telemetry := cap.NetworkTelemetry{LatencyMs: 0, PacketLoss: 0, Throughput: 0}
	predictor := &cap.SimplePartitionPredictor{Telemetry: telemetry}
	orchestrator := cap.NewOrchestrator(int(cap.EventualConsistency), predictor)
	ac := cap.NewAdaptiveConsistency(cap.EventualConsistency, 2*time.Second, cap.RetryPolicy{MaxRetries: 3, Backoff: 1 * time.Second})

	// Interactive CLI loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nOptions:")
		fmt.Println("1) Show blockchain")
		fmt.Println("2) Add block")
		fmt.Println("3) Register node")
		fmt.Println("4) New auth challenge")
		fmt.Println("5) Validate auth response")
		fmt.Println("6) Show node trust scores")
		fmt.Println("7) Input network telemetry")
		fmt.Println("8) Run CAP orchestration")
		fmt.Println("9) Show consistency level")
		fmt.Println("10) Run hybrid consensus")
		fmt.Println("11) Exit")
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			for _, blk := range bc.Blocks {
				fmt.Printf("%+v\n", blk)
			}
		case "2":
			fmt.Print("Enter transaction data: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimSpace(data)
			index := len(bc.Blocks)
			prev := ""
			if index > 0 {
				prev = bc.Blocks[index-1].Hash
			}
			block := blockchain.NewBlock(index, prev, []string{data}, nil, nil, 0)
			bc.AddBlock(block)
			fmt.Println("Block added:", block)
			// Archive the new block
			if err := blockchain.ArchiveBlock(block); err != nil {
				fmt.Println("ArchiveBlock error:", err)
			} else {
				fmt.Println("Block archived to disk.")
			}
			// Archive current state
			stateMap := make(map[string]interface{})
			for i, blk := range bc.Blocks {
				stateMap[strconv.Itoa(i)] = blk.Hash
			}
			if err := blockchain.ArchiveState(stateMap); err != nil {
				fmt.Println("ArchiveState error:", err)
			} else {
				fmt.Println("State archived to disk.")
			}
		case "3":
			fmt.Print("Node ID: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)
			fmt.Print("Public key: ")
			pub, _ := reader.ReadString('\n')
			pub = strings.TrimSpace(pub)
			authMgr.AddNode(id, pub)
			repSys.UpdateReputation(id, 0)
			fmt.Println("Node registered.")
		case "4":
			fmt.Print("Enter node ID: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)
			chal, err := authMgr.NewChallenge(id)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Challenge:", chal)
			}
		case "5":
			fmt.Print("Enter node ID: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)
			fmt.Print("Enter HMAC response: ")
			resp, _ := reader.ReadString('\n')
			resp = strings.TrimSpace(resp)
			if err := authMgr.ValidateResponse(id, resp); err != nil {
				fmt.Println("Validation failed:", err)
			} else {
				fmt.Println("Validation succeeded.")
			}
		case "6":
			fmt.Println("Trust scores:")
			scores := authMgr.GetTrustScores()
			for id, score := range scores {
				fmt.Printf("%s: %.2f\n", id, score)
			}
		case "7":
			fmt.Print("Latency ms: ")
			fmt.Scanf("%f\n", &telemetry.LatencyMs)
			fmt.Print("Packet loss (0-1): ")
			fmt.Scanf("%f\n", &telemetry.PacketLoss)
			fmt.Print("Throughput: ")
			fmt.Scanf("%f\n", &telemetry.Throughput)
			predictor.Telemetry = telemetry
			fmt.Println("Telemetry updated.")
		case "8":
			ctx := context.Background()
			ac.Orchestrate(ctx, telemetry, orchestrator)
			fmt.Println("CAP orchestration executed.")
		case "9":
			fmt.Println("Current consistency level:", orchestrator.CurrentLevel())
		case "10":
			// Hybrid consensus flow
			fmt.Print("Enter transaction data: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimSpace(data)
			validators := []string{}
			for id := range authMgr.GetTrustScores() {
				validators = append(validators, id)
			}
			if len(validators) == 0 {
				fmt.Println("No validators registered.")
				break
			}
			hc := consensus.NewHybridConsensus(validators)
			hc.StartRound()
			// propose block
			index := len(bc.Blocks)
			prev := ""
			if index > 0 {
				prev = bc.Blocks[index-1].Hash
			}
			block := blockchain.NewBlock(index, prev, []string{data}, nil, nil, 0)
			hc.ProposeBlock(block.Hash)
			// all validators vote yes
			for _, v := range validators {
				hc.Vote(v, "yes")
			}
			// finalize
			hc.FinalizeRound()
			// simple majority check
			if len(validators) > 0 {
				bc.AddBlock(block)
				fmt.Println("Consensus reached, block added:", block)
				// Archive the committed block
				if err := blockchain.ArchiveBlock(block); err != nil {
					fmt.Println("ArchiveBlock error:", err)
				} else {
					fmt.Println("Block archived to disk.")
				}
				// Archive state
				stateMap := make(map[string]interface{})
				for i, blk := range bc.Blocks {
					stateMap[strconv.Itoa(i)] = blk.Hash
				}
				if err := blockchain.ArchiveState(stateMap); err != nil {
					fmt.Println("ArchiveState error:", err)
				} else {
					fmt.Println("State archived to disk.")
				}
			} else {
				fmt.Println("Consensus not reached.")
			}
		case "11":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid option")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func NewBlockchain() *blockchain.Blockchain {
	return &blockchain.Blockchain{
		Blocks: []*blockchain.Block{},
	}
}
