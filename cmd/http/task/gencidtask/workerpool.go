package gencidtask

import (
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jianbo-zh/go-config"
)

// ---------------------
// 4. direct upload small ipfs nodes
// ---------------------
func createSmallWorkerPool() {

	ipfsNodes4 := config.GetString("ipfs.b_small_nodes", "")
	if ipfsNodes4 == "" {
		pinLog.Panic("ipfs.b_small_nodes not config")
	}
	ipfsNodeArr4 := []string{}
	for _, ipfsNode := range strings.Split(ipfsNodes4, ",") {
		ipfsNode = strings.TrimSpace(ipfsNode)
		if ipfsNode != "" {
			ipfsNodeArr4 = append(ipfsNodeArr4, ipfsNode)
		}
	}

	if len(ipfsNodeArr4) == 0 {
		pinLog.Panic("no valid ipfs small nodes")
	}

	workerNum4 := config.GetInt("worker.b_small_num", 1)
	BsShells = make(chan *shell.Shell, workerNum4)
	for i := 0; i < workerNum4; i++ {
		BsShells <- shell.NewShell(ipfsNodeArr4[i%len(ipfsNodeArr4)])
	}
}
