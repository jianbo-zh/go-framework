package gencidtask

import (
	"fmt"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jianbo-zh/go-config"
)

// ---------------------
// 1. large ipfs nodes
// ---------------------
func startLargeFileWorker() {

	ipfsNodes1 := config.GetString("ipfs.large_nodes", "")
	if ipfsNodes1 == "" {
		pinLog.Panic("ipfs.large_nodes not config")
	}

	ipfsNodeArr1 := []string{}
	for _, ipfsNode := range strings.Split(ipfsNodes1, ",") {
		ipfsNode = strings.TrimSpace(ipfsNode)
		if ipfsNode != "" {
			ipfsNodeArr1 = append(ipfsNodeArr1, ipfsNode)
		}
	}

	if len(ipfsNodeArr1) == 0 {
		pinLog.Panic("no valid ipfs large nodes")
	}

	workerNum1 := config.GetInt("worker.large_num", 1)
	alShells = make([]*shell.Shell, workerNum1)
	for i := 0; i < workerNum1; i++ {
		alShells[i] = shell.NewShell(ipfsNodeArr1[i%len(ipfsNodeArr1)])
	}

	for i := 0; i < len(alShells); i++ {
		go func(idx int) {
			tag := fmt.Sprintf("al-%d", idx)

			defer func() {
				pinLog.Errorf("[%s] handle file task break", tag)
			}()

			for mufi := range largeFileChannel {
				HandleQueueTask(tag, alShells[idx], mufi)
			}
		}(i)
	}
}

// ---------------------
// 2. normal ipfs nodes
// ---------------------
func startNormalFileWorker() {
	ipfsNodes2 := config.GetString("ipfs.normal_nodes", "")
	if ipfsNodes2 == "" {
		pinLog.Panic("ipfs.normal_nodes not config")
	}
	ipfsNodeArr2 := []string{}
	for _, ipfsNode := range strings.Split(ipfsNodes2, ",") {
		ipfsNode = strings.TrimSpace(ipfsNode)
		if ipfsNode != "" {
			ipfsNodeArr2 = append(ipfsNodeArr2, ipfsNode)
		}
	}
	if len(ipfsNodeArr2) == 0 {
		pinLog.Panic("no valid ipfs normal nodes")
	}

	workerNum2 := config.GetInt("worker.normal_num", 1)
	anShells = make([]*shell.Shell, workerNum2)
	for i := 0; i < workerNum2; i++ {
		anShells[i] = shell.NewShell(ipfsNodeArr2[i%len(ipfsNodeArr2)])
	}
	for i := 0; i < len(anShells); i++ {
		go func(idx int) {
			tag := fmt.Sprintf("am-%d", idx)
			defer func() {
				pinLog.Errorf("[%s] handle file task break", tag)
			}()

			for mufi := range normalFileChannel {
				HandleQueueTask(tag, anShells[idx], mufi)
			}
		}(i)
	}
}

// ---------------------
// 3. small ipfs nodes
// ---------------------
func startSmallFileWorker() {

	ipfsNodes3 := config.GetString("ipfs.small_nodes", "")
	if ipfsNodes3 == "" {
		pinLog.Panic("ipfs.small_nodes not config")
	}
	ipfsNodeArr3 := []string{}
	for _, ipfsNode := range strings.Split(ipfsNodes3, ",") {
		ipfsNode = strings.TrimSpace(ipfsNode)
		if ipfsNode != "" {
			ipfsNodeArr3 = append(ipfsNodeArr3, ipfsNode)
		}
	}

	if len(ipfsNodeArr3) == 0 {
		pinLog.Panic("no valid ipfs small nodes")
	}

	workerNum3 := config.GetInt("worker.small_num", 1)
	asShells = make([]*shell.Shell, workerNum3)
	for i := 0; i < workerNum3; i++ {
		asShells[i] = shell.NewShell(ipfsNodeArr3[i%len(ipfsNodeArr3)])
	}
	for i := 0; i < len(asShells); i++ {
		go func(idx int) {
			tag := fmt.Sprintf("as-%d", idx)
			defer func() {
				pinLog.Errorf("[%s] handle file task break", tag)
			}()

			for mufi := range smallFileChannel {
				HandleQueueTask(tag, asShells[idx], mufi)
			}
		}(i)
	}
}
