package wallet

import (
	"ec-wallet/internal/errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// 硬化衍生常量 (2^31)
	HardenedBit = uint32(0x80000000)
	// 非硬化掩碼
	NonHardenedMask = uint32(0x7FFFFFFF)

	// BIP44 路徑組件索引
	PurposeIndex      = 0
	CoinTypeIndex     = 1
	AccountIndex      = 2
	ChangeIndex       = 3
	AddressIndexIndex = 4

	// 常用幣種類型
	CoinTypeBTC = uint32(0)
	CoinTypeETH = uint32(60)
	CoinTypeBSC = uint32(60)
)

var SLIP0044SymbolToType = map[string]uint32{
	"BTC": CoinTypeBTC,
	"ETH": CoinTypeETH,
	"BSC": CoinTypeBSC,
}

type HDPath struct {
	// 存儲原始路徑字符串，例如 "m/44'/60'/0'/0/0"
	Original string

	// 存儲解析後的路徑組件數組
	Components []uint32

	// BIP44 路徑具體部分
	Purpose      uint32 // 目的 (通常是 44')
	CoinType     uint32 // 幣種類型 (例如 60' 表示 Ethereum)
	Account      uint32 // 帳戶編號
	Change       uint32 // 變更類型 (0 外部鏈/收款地址，1 內部鏈/找零地址)
	AddressIndex uint32 // 地址索引
}

// NewStandardPath 創建指定幣種的標準路徑
func NewStandardPath(coinType uint32, accountIndex uint32, addressIndex uint32) *HDPath {
	return NewHDPathFromComponents(44, coinType, accountIndex, 0, addressIndex)
}

// ParsePath 從路徑字符串解析 HDPath
func NewHDPathFromString(path string) (*HDPath, error) {
	if !strings.HasPrefix(path, "m/") {
		return nil, errors.ErrHDPathInvalidFormat
	}

	// 移除 'm/' 前綴
	path = path[2:]

	// 分割路徑組件
	components := strings.Split(path, "/")
	result := make([]uint32, len(components))

	for i, comp := range components {
		hardened := false
		if strings.HasSuffix(comp, "'") || strings.HasSuffix(comp, "h") {
			hardened = true
			comp = comp[:len(comp)-1]
		}

		val, err := strconv.ParseUint(comp, 10, 32)
		if err != nil {
			return nil, errors.ErrHDPathInvalidComponent
		}

		// 前三個組件（purpose, coinType, account）必須是硬化衍生的
		if i <= AccountIndex && !hardened {
			return nil, errors.ErrHDPathInvalidComponent
		}

		if hardened {
			result[i] = uint32(val) + HardenedBit
		} else {
			result[i] = uint32(val)
		}
	}

	hdPath := &HDPath{
		Original:   "m/" + path,
		Components: result,
	}

	// 設置 BIP44 特定欄位
	if len(result) > PurposeIndex {
		hdPath.Purpose = result[PurposeIndex]
	}
	if len(result) > CoinTypeIndex {
		hdPath.CoinType = result[CoinTypeIndex]
	}
	if len(result) > AccountIndex {
		hdPath.Account = result[AccountIndex]
	}
	if len(result) > ChangeIndex {
		hdPath.Change = result[ChangeIndex]
	}
	if len(result) > AddressIndexIndex {
		hdPath.AddressIndex = result[AddressIndexIndex]
	}

	return hdPath, nil
}

// NewHDPath 從路徑組件直接創建 HDPath
func NewHDPathFromComponents(purpose, coinType, account, change, addressIndex uint32) *HDPath {
	components := make([]uint32, 5)

	// 始終應用硬化衍生到前三個組件
	purpose += HardenedBit
	coinType += HardenedBit
	account += HardenedBit

	components[PurposeIndex] = purpose
	components[CoinTypeIndex] = coinType
	components[AccountIndex] = account
	components[ChangeIndex] = change
	components[AddressIndexIndex] = addressIndex

	// 構建原始路徑字符串
	path := fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		purpose&NonHardenedMask,
		coinType&NonHardenedMask,
		account&NonHardenedMask,
		change,
		addressIndex)

	return &HDPath{
		Original:     path,
		Components:   components,
		Purpose:      purpose,
		CoinType:     coinType,
		Account:      account,
		Change:       change,
		AddressIndex: addressIndex,
	}
}

// String 返回路徑字符串表示
func (p *HDPath) String() string {
	return p.Original
}

// GetUnhardenedValue 獲取組件值（去除硬化位）
func (p *HDPath) GetUnhardenedValue(index int) uint32 {
	if index < 0 || index >= len(p.Components) {
		return 0
	}
	return p.Components[index] & NonHardenedMask
}

// WithAddressIndex 創建新的衍生路徑（適用於迭代地址）
func (p *HDPath) WithAddressIndex(index uint32) *HDPath {
	newPath := &HDPath{
		Components:   make([]uint32, len(p.Components)),
		Purpose:      p.Purpose,
		CoinType:     p.CoinType,
		Account:      p.Account,
		Change:       p.Change,
		AddressIndex: index,
	}

	// 複製組件
	copy(newPath.Components, p.Components)

	// 更新地址索引
	if len(newPath.Components) > AddressIndexIndex {
		newPath.Components[AddressIndexIndex] = index
	}

	// 更新原始路徑字符串
	newPath.Original = fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		p.GetUnhardenedValue(PurposeIndex),
		p.GetUnhardenedValue(CoinTypeIndex),
		p.GetUnhardenedValue(AccountIndex),
		p.Change,
		index)

	return newPath
}

// IsValidBIP44Path 檢查路徑是否符合 BIP44 標準
func (p *HDPath) IsValidBIP44Path() bool {
	// BIP44 路徑必須至少有 5 個組件
	if len(p.Components) < 5 {
		return false
	}

	// 檢查第一個組件是否為 44'（硬化衍生的）
	if p.GetUnhardenedValue(PurposeIndex) != 44 || (p.Components[PurposeIndex]&HardenedBit) == 0 {
		return false
	}

	// 檢查前三個組件是否為硬化衍生
	return (p.Components[CoinTypeIndex]&HardenedBit) != 0 && (p.Components[AccountIndex]&HardenedBit) != 0
}
