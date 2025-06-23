package chainlistener

// type TxConfirmationTracker struct {
// 	txHash           common.Hash
// 	firstSeenBlock   uint64
// 	confirmations    uint64
// 	maxConfirmations uint64 // 達到此確認數後認為交易最終確定
// 	lastCheckedBlock uint64
// 	confirmed        bool
// }

// func (t *TxConfirmationTracker) CheckBlock(ctx context.Context, block *types.Block) (bool, error) {
// 	// 如果已經確認，直接返回
// 	if t.confirmed {
// 		return true, nil
// 	}

// 	// 獲取交易收據
// 	receipt, err := client.TransactionReceipt(ctx, t.txHash)
// 	if err != nil {
// 		if err.Error() == "not found" {
// 			// 交易可能被重組掉了
// 			return false, fmt.Errorf("交易 %s 在區塊 %d 中未找到，可能已被重組",
// 				t.txHash.Hex(), block.NumberU64())
// 		}
// 		return false, err
// 	}

// 	// 檢查交易是否在當前區塊的區塊鏈上
// 	if receipt.BlockHash != block.Hash() {
// 		// 交易存在，但不在當前區塊鏈上，可能發生了重組
// 		return false, fmt.Errorf("交易 %s 不在主鏈上，可能發生了重組", t.txHash.Hex())
// 	}

// 	// 更新確認數
// 	currentConfirmations := block.NumberU64() - receipt.BlockNumber.Uint64() + 1
// 	if currentConfirmations > t.confirmations {
// 		t.confirmations = currentConfirmations
// 		t.lastCheckedBlock = block.NumberU64()
// 	}

// 	// 檢查是否達到最終確認
// 	if t.confirmations >= t.maxConfirmations {
// 		t.confirmed = true
// 		return true, nil
// 	}

// 	return false, nil
// }

// // 添加交易追蹤
// func (l *EVMChainListener) TrackTransaction(txHash common.Hash, firstSeenBlock uint64) {
// 	l.trackerMutex.Lock()
// 	defer l.trackerMutex.Unlock()

// 	if _, exists := l.txTrackers[txHash]; !exists {
// 		l.txTrackers[txHash] = &TxConfirmationTracker{
// 			txHash:           txHash,
// 			firstSeenBlock:   firstSeenBlock,
// 			maxConfirmations: l.maxConfirmations,
// 			lastCheckedBlock: firstSeenBlock,
// 		}
// 	}
// }

// // 在處理新區塊時檢查所有追蹤的交易
// func (l *EVMChainListener) checkTrackedTransactions(ctx context.Context, block *types.Block) {
// 	l.trackerMutex.RLock()
// 	defer l.trackerMutex.RUnlock()

// 	for _, tracker := range l.txTrackers {
// 		if tracker.confirmed {
// 			continue
// 		}

// 		confirmed, err := tracker.CheckBlock(ctx, block)
// 		if err != nil {
// 			log.Printf("檢查交易 %s 時出錯: %v", tracker.txHash.Hex(), err)
// 			// 可以考慮從追蹤中移除或標記為失敗
// 			continue
// 		}

// 		if confirmed {
// 			log.Printf("交易 %s 已達到最終確認 (%d/%d 確認)",
// 				tracker.txHash.Hex(),
// 				tracker.confirmations,
// 				tracker.maxConfirmations)
// 			// 發送確認通知
// 			l.notifyTxConfirmed(tracker.txHash, block)
// 		}
// 	}
// }

// func (l *EVMChainListener) handleReorg(oldBlocks []*types.Block, newBlocks []*types.Block) {
// 	log.Printf("檢測到區塊重組: 移除了 %d 個區塊，添加了 %d 個新區塊",
// 		len(oldBlocks), len(newBlocks))

// 	// 更新追蹤器中的區塊高度
// 	l.trackerMutex.Lock()
// 	defer l.trackerMutex.Unlock()

// 	for _, tracker := range l.txTrackers {
// 		// 重置確認計數，強制重新檢查
// 		tracker.confirmations = 0
// 		tracker.confirmed = false
// 	}
// }
