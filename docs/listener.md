[使用者] 
   │
   ▼
[電商服務]───(1) 用戶下單───▶
   │
   ▼
(2) 分配一個唯一地址 (動態錢包)
   │
   ▼
(3) 建立訂單資料（包含等待付款狀態）
   │
   ▼
(4) 通知監聽器「啟用地址監聽」
   Redis Stream: "listen:evm"
   消息內容: {
       order_id,
       address,
       token_type,
       chain,
       amount,
       expired_block,
       expected_confirmations
   }
   │
   ▼
[監聽器服務]
   │
   ├─▶(5) 將監聽任務加入動態 watcher（存在內存、可加 Redis 緩存）
   │
   └─▶(6) 開始監聽區塊與地址交易/事件（透過 `eth_subscribe`）
   │
   ▼
[鏈上有用戶付款]
   │
   ▼
(7) 監聽器捕捉到匹配的 tx 或 event
   │
   ▼
(8) 初步過濾符合條件交易
   │
   ▼
(9) 將 tx 資訊推送 Redis Stream「tx:pending」
   消息內容: {
       order_id,
       tx_hash,
       from,
       to,
       amount,
       block_number
   }
   │
   ▼
[確認 Worker (可在監聽器或電商中)]
   │
   └─▶(10) 根據 tx_hash 查詢 tx receipt
         - 是否成功
         - 是否進區塊
         - 確認數是否達標
         - 金額、地址是否符合
   │
   ▼
(11) 若通過驗證，送出 Redis Stream：「tx:confirmed」
   消息內容: {
       order_id,
       tx_hash,
       confirmed: true,
       confirmed_at_block,
       confirmed_at_time
   }
   │
   ▼
[電商服務]
   │
   └─▶(12) 訂單狀態更新為已付款
   │
   ▼
(13) 通知前端 / 發貨等後續流程


| 元件               | 責任                                |
| ---------------- | --------------------------------- |
| **電商服務**         | 分配動態地址、送出監聽請求、處理已付款通知、驅動訂單邏輯      |
| **Redis Stream** | 作為監聽事件與 tx 狀態的中繼總線，保證解耦與可靠傳遞      |
| **監聽器服務**        | 接收訂單指令、實時監聽區塊/交易事件、初步過濾與記錄        |
| **確認 worker**    | 根據 tx\_hash 驗證交易真偽、是否進鏈與確認數，保障安全性 |


| 功能                          | 說明                                        |
| --------------------------- | ----------------------------------------- |
| **動態監聽 TTL**                | 到期後自動移除任務，避免監聽器資源佔用太多                     |
| **監聽結果 timeout / fallback** | 若過期未付款，電商可關閉訂單並通知監聽器停止監聽                  |
| **監聽器持久化支持**                | 若你部署為多實例，可考慮將 watcher 狀態也保存在 Redis，實現水平擴展 |


| Channel 名稱           | 用途                       |
| -------------------- | ------------------------ |
| `chain:watch:start`  | 電商送出新地址要監聽的請求            |
| `chain:tx:pending`   | 監聽器發現疑似 tx，推送給驗證 worker  |
| `chain:tx:confirmed` | 驗證成功的 tx，通知電商更新訂單        |
| `chain:watch:stop`   | 電商通知監聽器可停止監聽該地址（過期或確認完成） |
