package geth_test

import (
    "context"
    "math/big"
    "testing"
    "time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/stretchr/testify/require"
)

const (
    gethEndpoint = "http://localhost:8545" // Update this with your Geth endpoint
    testTimeout  = 5 * time.Minute
)

func TestGethLoadTest(t *testing.T) {
    client, err := ethclient.Dial(gethEndpoint)
    require.NoError(t, err)
    defer client.Close()

    t.Run("Block Generation Test", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
        defer cancel()

        // Get initial block number
        startBlock, err := client.BlockNumber(ctx)
        require.NoError(t, err)

        // Monitor block generation for 1 minute
        time.Sleep(1 * time.Minute)

        endBlock, err := client.BlockNumber(ctx)
        require.NoError(t, err)

        blocksGenerated := endBlock - startBlock
        t.Logf("Blocks generated in 1 minute: %d", blocksGenerated)
        require.Greater(t, blocksGenerated, uint64(0), "Expected at least one block to be generated")
    })

    t.Run("Transaction Response Time", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
        defer cancel()

        // Test transaction query response time
        address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e") // Example address
        for i := 0; i < 10; i++ {
            start := time.Now()
            
            _, err := client.BalanceAt(ctx, address, nil)
            require.NoError(t, err)
            
            duration := time.Since(start)
            t.Logf("Query %d took: %v", i+1, duration)
            require.Less(t, duration, 1*time.Second, "Query took too long")
        }
    })

    t.Run("Pending Transactions Test", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
        defer cancel()

        // Monitor pending transactions
        start := time.Now()
        count := 0
        for time.Since(start) < 30*time.Second {
            pending, err := client.PendingTransactionCount(ctx)
            require.NoError(t, err)
            
            if pending > 0 {
                count++
                t.Logf("Pending transactions: %d", pending)
            }
            
            time.Sleep(1 * time.Second)
        }
        t.Logf("Total measurements with pending transactions: %d", count)
    })

    t.Run("Gas Price Monitoring", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
        defer cancel()

        gasPrice, err := client.SuggestGasPrice(ctx)
        require.NoError(t, err)
        require.NotNil(t, gasPrice)
        
        t.Logf("Current gas price: %s wei", gasPrice.String())
        require.Greater(t, gasPrice.Cmp(big.NewInt(0)), 0, "Gas price should be greater than 0")
    })
}