package forging

import (
	"encoding/binary"
	"pandora-pay/blockchain/block/difficulty"
	"pandora-pay/config"
	"pandora-pay/config/stake"
	"pandora-pay/crypto"
	"sync/atomic"
	"time"
)

//inside a thread
func forge(threads, threadIndex int) {

	buf := make([]byte, binary.MaxVarintLen64)

	ForgingW.RLock()
	defer ForgingW.RUnlock()
	defer wg.Done()

	height := Forging.BlkComplete.Block.Height
	serialized := Forging.BlkComplete.Block.SerializeBlock(true, false)
	n := binary.PutUvarint(buf, Forging.BlkComplete.Block.Timestamp)

	serialized = serialized[:len(serialized)-n-20]
	timestamp := Forging.BlkComplete.Block.Timestamp + 1

	for atomic.LoadInt32(&forgingWorking) == 1 {

		if timestamp > uint64(time.Now().Unix())+config.NETWORK_TIMESTAMP_DRIFT_MAX {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		//forge with my wallets
		for i, address := range ForgingW.addresses {

			if i%threads == threadIndex && (address.account != nil || height == 0) {

				var stakingAmount uint64
				if address.account != nil {
					stakingAmount = address.account.GetDelegatedStakeAvailable(height)
				}

				if stakingAmount >= stake.GetRequiredStake(height) {

					if atomic.LoadInt32(&forgingWorking) == 0 {
						break
					}

					n = binary.PutUvarint(buf, timestamp)
					serialized = append(serialized, buf[:n]...)
					serialized = append(serialized, address.publicKeyHash[:]...)
					kernelHash := crypto.SHA3Hash(serialized)

					if height > 0 {
						kernelHash = crypto.ComputeKernelHash(kernelHash, stakingAmount)
					}

					if difficulty.CheckKernelHashBig(kernelHash, Forging.target) {

						Forging.foundSolution(address, timestamp)

					} else {
						// for debugging only
						// gui.Log(hex.EncodeToString(kernelHash[:]))
					}

					serialized = serialized[:len(serialized)-n-20]

				}

			}

		}
		timestamp += 1

	}

}
