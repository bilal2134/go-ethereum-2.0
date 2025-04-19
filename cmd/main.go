package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bilal2134/Blockchain_A3/internal/blockchain"
)

// Entry point for the blockchain node
func main() {
	// Initialize the blockchain
	bc := NewBlockchain()

	// Interactive CLI loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nOptions:")
		fmt.Println("1) Show blockchain")
		fmt.Println("2) Add block")
		fmt.Println("3) Exit")
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
			block := blockchain.NewBlock(index, prev, []string{data}, []byte{}, [][]byte{}, 1.0)
			bc.AddBlock(block)
			fmt.Println("Block added:", block)
		case "3":
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
