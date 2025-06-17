package walletservice

import (
	"context"
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/errors"
	"time"
)

// WalletPool 資源池操作命名參考：
//
// 從池中拿出一個地址:      Allocate, Acquire, Get, Pop
// 將地址歸還到池中:       Release, Free, Return, Recycle, PutBack
// 查看某個地址是否在池內:   Contains, Exists, IsAvailable
// 列出目前可用的地址:      ListAvailable, AvailableAddrs
// 初始化池資源:          Init, Load, Seed, Bootstrap
// 清除/重置池:           Reset, Clear, Flush
// 鎖定資源避免重複使用:     Reserve, Lock
// 永久刪除某地址:         Remove, Ban, Blacklist

func (s *walletService) InitWalletAddressPools(ctx context.Context, chain string, count, batchSize int) ([]uint64, error) {
	coinType, ok := wallet.SLIP0044SymbolToType[chain]
	if !ok {
		return nil, errors.ErrWalletUnsupportedChain
	}

	if count <= 0 {
		return nil, errors.ErrWalletInvalidAddressCount
	}

	if batchSize <= 0 {
		batchSize = 100 // 默認批量大小
	}

	// 獲取助記詞
	mnemonic := s.GetMnemonic()
	if mnemonic == "" {
		return nil, errors.ErrWalletMnemonicRequired
	}

	var allIDs []uint64
	// 初始化的startIndex一定是零
	startIndex := 0
	// start tx
	tx := s.repo.Begin()
	defer s.repo.RollBack(tx)
	// 分批處理
	for i := 0; i < count; i += batchSize {
		// 處理含尾段的當前batch大小
		currentBatchSize := min(batchSize, count-i)
		pools := make([]*gormrepo.WalletAddressPool, 0, currentBatchSize)
		// 為當前批次生成地址
		for j := 0; j < currentBatchSize; j++ {
			index := startIndex + i + j
			hdPath := wallet.NewStandardPath(coinType, 0, uint32(index))
			keyPair, err := s.DeriveKeyFromPath(mnemonic, hdPath)
			if err != nil {
				return allIDs, err
			}

			pool := &gormrepo.WalletAddressPool{
				Address:       keyPair.Address,
				Chain:         chain,
				Path:          hdPath.String(),
				Index:         index,
				CurrentStatus: wallet.AddressStatusAvailable,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			pools = append(pools, pool)
		}

		// 保存到資料庫
		ids, err := s.repo.CreateWalletAddressPools(ctx, tx, pools)
		if err != nil {
			return allIDs, err
		}

		allIDs = append(allIDs, ids...)
	}

	_ = s.repo.Commit(tx)

	return allIDs, nil
}

// 獲取並佔用地址資源
func (s *walletService) AcquireAddress(ctx context.Context, opts ...wallet.AcquireOption) (*wallet.AddressReservation, error) {
	// 應用默認選項
	options := wallet.AcquireOptions{
		CoinType:  wallet.CoinTypeETH,
		ExpiresIn: 24 * time.Hour,
	}

	// 應用自定義選項
	for _, opt := range opts {
		opt(&options)
	}

	// tx start
	tx := s.repo.Begin()
	defer s.repo.RollBack(tx)

	// 拿取Available address
	status := wallet.AddressStatusAvailable
	pool, err := s.repo.GetWalletAddressPool(ctx, tx, &gormrepo.QueryWalletAddressPoolsParams{
		CurrentStatus: &status,
	})
	if err != nil {
		return nil, err
	}

	// 改為佔用
	expiryTime := time.Now().Add(options.ExpiresIn)
	newStatus := wallet.AddressStatusReserved
	rowsAffected, err := s.repo.UpdateWalletAddressPools(ctx, tx, &gormrepo.UpdateWalletAddressPoolsParams{
		Where: gormrepo.QueryWalletAddressPoolsParams{
			ID: &pool.ID,
		},
		CurrentStatus: &newStatus,
		ReservedUntil: &expiryTime,
	})
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 {
		return nil, errors.ErrWalletAddressPoolUpdate
	}

	// 增加調用紀錄
	log := gormrepo.NewWalletAddressLog(
		&gormrepo.NewWalletAddressLogParams{
			AddressID:    pool.ID,
			Operation:    wallet.AddressLogOperationPayment,
			StatusAfter:  wallet.AddressStatusReserved,
			StatusBefore: wallet.AddressStatusAvailable,
			OperationAt:  time.Now(),
			ValidUntil:   &expiryTime,
			OrderID:      options.OrderID,
			UserID:       options.UserID,
		},
	)

	// create log
	ids, err := s.repo.CreateWalletAddressLogs(ctx, tx, []*gormrepo.WalletAddressLog{log})
	if err != nil {
		return nil, err
	}

	if len(ids) != 1 {
		return nil, errors.ErrWalletAddressLogCreate
	}

	_ = s.repo.Commit(tx)

	// 產生地址佔用資訊
	reservation := wallet.NewAddressReservation(&wallet.NewAddressReservationParams{
		Address:    pool.Address,
		AddressID:  pool.ID,
		ReservedAt: time.Now(),
		ExpiresAt:  expiryTime,
		CoinType:   options.CoinType,
	})
	return reservation, nil
}

func (s *walletService) ReleaseAddress(ctx context.Context) {

}

func (s *walletService) GetMnemonic() string {
	return s.config.Wallet.Mnemonic
}
