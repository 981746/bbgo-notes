## 要撰寫一個BBGO的內建策略要從何開始呢？

### 假設我要撰寫一個m15與H2的SMA交叉的BBGO內建策略，策略目前進度—拿取m5的K棒資料

BBGO的策略有兩種類型  
一種是built-in strategy(內建策略) - 內建策略會被包在預先編譯的binary檔案裡，這份文件要教學的是如何撰寫內建策略  
另一種是external strategy(外部策略) - 外部策略可以讓你放在自己建立的repository裡，因此與BBGO專案是分離的。如果不想公開的話就可以使用外部策略這個類型  
bbgo專案說明頁([https://github.com/c9s/bbgo](https://github.com/c9s/bbgo))的「Write your own private strategy」，就是在說外部策略怎麼弄，之後我也會再寫詳細的文件，說明如何撰寫私人的外部策略

*   內建策略撰寫參考資料：  
    [https://github.com/c9s/bbgo/blob/v1.39.0/doc/topics/developing-strategy.md](https://github.com/c9s/bbgo/blob/v1.39.0/doc/topics/developing-strategy.md)
    
*   BBGO專案資料夾位置： ~/bbgo
    
*   Mysql資料庫：使用者名稱：root；密碼：Your Password
    
*   策略名稱： smacross
    
*   import策略位置： ~/bbgo/pkg/cmd/strategy/builtin.go
    
*   策略本身檔案位置： ~/bbgo/pkg/strategy/smacross/strategy.go
    
*   策略設定檔位置： ~/bbgo/config/smacross.yaml
    
*   BBGO「.env.local」設定檔位置： ~/bbgo/.env.local
    
*   開發環境：Windows10 Wsl2 Ubuntu 20.04
    
*   幣安API權限： 只有開Enable Reading權限
    
*   bbgo run指令(一啟動會去sync，要等一個多小時)：  
    cd ~/bbgo  
    go run ./cmd/bbgo run --config ~/bbgo/config/smacross.yaml
    
*   bbgo run指令，跳過一啟動就sync：  
    cd ~/bbgo  
    go run ./cmd/bbgo run --config ~/bbgo/config/smacross.yaml --no-sync
    
*   bbgo backtest指令：  
    cd ~/bbgo  
    go run ./cmd/bbgo backtest --config ~/bbgo/config/smacross.yaml
    
*   bbgo backtest指令加上`--sync`，「先sync、後回測」：  
    cd ~/bbgo  
    go run ./cmd/bbgo backtest --sync --config ~/bbgo/config/smacross.yaml
    
*   phpmyadmin網址：
    
        localhost/phpmyadmin
        
    
*   查目前有哪些幣種已新增過K線資料–幣安交易所
    
        SELECT DISTINCT(symbol) FROM binance_klines;
        
    
*   查某個交易對的K線資料是否新增成功–幣安交易所
    
        SELECT * FROM binance_klines WHERE symbol = 'XRPUSDT'
        
    

### 請預先裝好：

Wsl2 Ubuntu 20.04、GO 開發環境(GO SDK + VS CODE)、Mysql安裝在wsl2上、Redis安裝在wsl2上(查看redis狀態，發現被拒絕連線;已解決，請看「解決redis服務與wsl2的連線問題」文件)、phpmyadmin(方便查看資料庫狀況)

### 步驟

我使用VS CODE remote到wsl2，然後開terminal來做相關的操作

1.  從BBGO專案git clone開始
    
        cd ~
        git clone https://github.com/c9s/bbgo.git
        
    
    git clone出現奇怪的錯誤(以前一切順利，沒有遇到這個錯誤過)
    
        jenny@LAPTOP-HF07DJBU:~$ git clone https://github.com/c9s/bbgo.git
        Cloning into 'bbgo'...
        remote: Enumerating objects: 53429, done.
        remote: Counting objects: 100% (53429/53429), done.
        remote: Compressing objects: 100% (13020/13020), done.
        error: RPC failed; curl 92 HTTP/2 stream 0 was not closed cleanly: CANCEL (err 8)
        fatal: the remote end hung up unexpectedly
        fatal: early EOF
        fatal: index-pack failed
        
    
    解法：  
    使用SSH與Github連線，再用SSH git clone
    
    SSH git clone  
    (在98%處跑很久，我有先去做別的事)
    
        git clone git@github.com:c9s/bbgo.git
        
    
2.  定義策略名稱  
    我目前要的是是m15與H2的SMA交叉的策略  
    所以我把策略取名為：smacross  
    那我必須去pkg/strategy這個資料夾底下 建一個名為smacross的資料夾
    
        cd ~/bbgo/pkg/strategy 
        mkdir smacross
        
    
3.  在~/bbgo/pkg/strategy/smacross資料夾底下  
    建一個strategy.go檔案，我們的策略就會寫在這個檔案
    
        cd smacross
        vim strategy.go
        :wq
        
    
4.  現在到~/bbgo/pkg/cmd/strategy/builtin.go這個檔案，import我們的smacross策略
    
        vim ~/bbgo/pkg/cmd/strategy/builtin.go
        
    
    檔案內容  
    主要就是有一個import function  
    我們在import()的參數填入我們剛剛建立的smacross資料夾所相依的專案路徑  
    格式為：`_ "github.com/c9s/bbgo/pkg/strategy/策略名稱"`
    
        package strategy
        
        // import built-in strategies
        import (
                _ "github.com/c9s/bbgo/pkg/strategy/atrpin"
                _ "github.com/c9s/bbgo/pkg/strategy/audacitymaker"
                _ "github.com/c9s/bbgo/pkg/strategy/autoborrow"
                ......................................
                _ "github.com/c9s/bbgo/pkg/strategy/smacross"
        )
        
    
5.  現在回到~/bbgo/pkg/strategy/smacross/strategy.go檔案撰寫我們的策略 我先寫一個簡單的拿BTC  
    m5 收盤價資料這樣就好
    
        VS Code→open folder→切到/home/jenny/bbgo/
        接著左側 EXPLORER 面板→找出/pkg/strategy/smacross/strategy.go檔案
        
    
    貼上我們策略的程式碼：
    
        package smacross
        
        import (
        	"context"
        	"fmt"
        
        	"github.com/c9s/bbgo/pkg/bbgo"
        	"github.com/c9s/bbgo/pkg/types"
        )
        
        const ID = "smacross"
        
        func init() {
        	// Register our struct type to BBGO
        	// Note that you don't need to field the fields.
        	// BBGO uses reflect to parse your type information.
        	bbgo.RegisterStrategy(ID, &Strategy{})
        }
        
        type Strategy struct {
        	Symbol   string         `json:"symbol"`   //交易對名稱會從yaml設定檔傳過來
        	Interval types.Interval `json:"interval"` //這個接收過來的參數暫時沒用到，我後面直接寫死
        }
        
        func (s *Strategy) ID() string {
        	return ID
        }
        
        // func (s *Strategy) InstanceID() string {
        // 	return fmt.Sprintf("%s:%s:%s", ID, s.Symbol, s.Interval)
        // }
        
        func (s *Strategy) Subscribe(session *bbgo.ExchangeSession) {
        	// session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.Interval})
        	session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: "5m"})
        }
        
        func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {
        	session.MarketDataStream.OnKLineClosed(func(k types.KLine) {
        		fmt.Println(k, "End of m5 data")
        	})
        	return nil
        }
        
    
    ※注意事項  
    第一、如果程式碼沒有ID()function，會出現panic
    
        func (s *Strategy) ID() string {
        return ID
        }
        
    
    > panic: \*smacross.Strategy does not implement SingleExchangeStrategy  
    > or CrossExchangeStrategy
    
    第二、 InstanceID()目前是註解掉的，因為我不確定這是幹嘛的，也不確定是不是必須的東西
    
        // func (s *Strategy) InstanceID() string {
        
        // return fmt.Sprintf("%s:%s:%s", ID, s.Symbol, s.Interval)
        
        // }
        
    
6.  在~/bbgo/config資料夾底下 我建一個smacross.yaml作為這個策略的設定檔
    
        cd ~/bbgo/config
        vim smacross.yaml
        
    
    貼上內容，存檔離開：
    
        exchangeStrategies:
        - on: binance
          smacross:
            symbol: BTCUSDT
            interval: 5m
        
    
7.  「.env.local」設定檔 修改~/bbgo/.env.local內容，填入幣安的API EKY、API SECRET與Mysql連線密碼
    
        vim ~/bbgo/.env.local
        
    
    貼上以下內容：
    
        BINANCE_API_KEY=Your API KEY
        BINANCE_API_SECRET=Your API Secret
        
        DB_DRIVER=mysql
        DB_DSN="root:Your Mysql Password@tcp(127.0.0.1:3306)/bbgo"
        
    
8.  查看mysql狀態，如果沒開請打開  
    wsl2一開始是用init.d處理系統服務，但是這樣msyql並不會自動啟動，必須手動啟動
    
        sudo /etc/init.d/mysql start
        
    
    而在wsl2設定systemctl之後，mysql在wsl2 terminal一打開就會自動啟動，不用手動特別去打開  
    請參考「解決redis服務與wsl2的連線問題」這份文件打開wsl2的systemctl
    
    以systemctl查看mysql狀態
    
        sudo systemctl status mysql
        
    
    以systemctl來啟動mysql
    
        sudo systemctl start mysql
        
    
9.  在執行bbgo run指令之前，我們先在mysql建一個bbgo資料庫
    
    直接下wsl2下指令(會要求我們輸入mysql root的密碼，而後自動完成新增)：
    
        mysql -uroot -p -e "CREATE DATABASE bbgo CHARSET utf8"
        
    
    接著可以去phpmyadmin確認bbgo資料庫是否成功新增  
    發現成功建了一個名為bbgo的資料庫，裡面還沒有任何的table  
    這個資料庫目前是空的
    
10.  執行此smacross策略
    
	     cd ~/bbgo
	     go run ./cmd/bbgo run --config ~/bbgo/config/smacross.yaml
        
    
	    比較奇怪的是，直接執行bbgo run指令，會卡在sync，策略被執行應該要可以看到拿到K線資料，但沒有
	    
	    喔～我以為sync是卡住了  
	    等了一陣子發現其實會動欸XDDDDD  
	    再等等看  
	    我好像是從中午12:40開始sync，現在時間是下午13:16，還在sync  
	    看log紀錄，好像是13:47:34才sync完成  
	    大概要花一個多小時
	    
	    去sync到底是去同步了甚麼資料啊？  
	    bbgo run指令所執行的sync，感覺是去抓連接的交易所的我自己帳戶曾經交易過的幣種資料，會把資料建在bbgo資料庫底下的「orders」與「trades」資料表  
	    bbgo run先去拿了帳戶底下的相關交易紀錄，而後直接跟交易所去連線要即時的K線資料，這個即時的K線資料並不會新增進資料庫  
	    策略會根據即時的K線資料來運行
	    
	    bbgo backtest所執行的sync，應該是抓K線歷史資料，我連的交易所是幣安，幣安現貨API可以支援的K線資料的時框有1m/5m/15m/30m/1h/2h/4h/6h/12h/1d/3d/1w  
	    那回測所使用到的幣種的相關的K線資料(有支援的時框都會新增)就會建立在「\[binance\_klines\]資料表裡
	    
	    所以  
	    bbgo run執行策略拿到的K線資料是即時K線資料  
	    bbgo backtest執行策略使用的是歷史K線資料來做回測
	    
	    bbgo run 的 sync 訊息大概是長這樣：
    
	     [0005]  INFO querying closed orders PEOPLEUSDT from 2023-11-02 12:40:07.523423859 +0800 CST <=> 2023-12-07 12:40:10.588001727 +0800 CST m=+3.091242680 ... exchange=binance
	     [0005]  INFO querying closed orders PEOPLEUSDT from 2023-12-02 12:40:07.523423859 +0800 CST <=> 2023-12-07 12:40:10.588001727 +0800 CST m=+3.091242680 ... exchange=binance
	     [0005]  INFO syncing binance PEOPLEBUSD trades from 2022-12-07 12:40:07.523423859 +0800 CST...
	     [0369]  INFO syncing binance PEOPLEBUSD orders from 2022-12-07 12:40:07.523423859 +0800 CST...
        
    
	    sync完成後的訊息，最後也出現了我要的m5的K棒資料：  
	    看起來是每五分鐘會吐一次K棒資料
    
         [4049]  INFO attaching strategy *smacross.Strategy on binance...
         [4049]  INFO loading strategies states...
         [4049]  INFO found symbol BTCUSDT based strategy from smacross.Strategy
         [4049]  INFO querying kline BTCUSDT 5m {1000 <nil> 2023-12-07 12:40:07.523423859 +0800 CST m=+0.026664912} exchange=binance
         [4049]  INFO querying kline BTCUSDT 1m {1000 <nil> 2023-12-07 12:40:07.523423859 +0800 CST m=+0.026664912} exchange=binance
         [4049]  INFO BTCUSDT last price: 44036.79
         [4049]  WARN StandardIndicatorSet() is deprecated in v1.49.0 and which will be removed in the next version, please use Indicators() instead
         [4049]  WARN messenger is not set, skip initializing
         [4049]  INFO subscribing BTCUSDT kline 5m session=binance
         [4049]  INFO connecting binance market data stream... session=binance
         [4050]  INFO websocket: connected, public = true, read timeout = 2m0s
         [4050]  INFO subscribing channels: [btcusdt@kline_5m] exchange=binance
          [4050]  INFO connecting binance user data stream... session=binance
          [4050]  INFO websocket: connected, public = false, read timeout = 2m0s
          [4050]  INFO session binance user data stream connected
         {0 binance BTCUSDT 2023-12-07 13:45:00 +0800 CST 2023-12-07 13:49:59.999 +0800 CST 5m 43990.01 43981.19 43993.8 43978.62 41.92447 1844072.8996063 20.328 894103.7920169 3310580153 2748 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 13:50:00 +0800 CST 2023-12-07 13:54:59.999 +0800 CST 5m 43981.2 43986.42 43986.42 43972.39 76.99942 3386283.1330992 42.79076 1881845.575337 3310582961 2808 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 13:55:00 +0800 CST 2023-12-07 13:59:59.999 +0800 CST 5m 43986.42 44007.25 44007.25 43971.11 98.32445 4325227.5461898 63.32751 2785712.405214 3310585967 3006 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 14:00:00 +0800 CST 2023-12-07 14:04:59.999 +0800 CST 5m 44007.24 43990.6 44030.01 43984.31 85.6543 3769715.5620298 43.08634 1896202.569543 3310589588 3621 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 14:05:00 +0800 CST 2023-12-07 14:09:59.999 +0800 CST 5m 43990.61 43972.8 43990.61 43972.79 43.83153 1927721.9720138 16.16159 710753.1156457 3310591848 2260 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 14:10:00 +0800 CST 2023-12-07 14:14:59.999 +0800 CST 5m 43972.8 43956 43972.8 43946.9 73.05359 3211144.2532015 32.62827 1434121.8323922 3310595105 3257 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 14:15:00 +0800 CST 2023-12-07 14:19:59.999 +0800 CST 5m 43956 43952.34 43983.2 43948.32 90.99357 4000647.049075 50.43495 2217482.3644649 3310598519 3414 true} End of m5 data
         {0 binance BTCUSDT 2023-12-07 14:20:00 +0800 CST 2023-12-07 14:24:59.999 +0800 CST 5m 43952.34 43950.08 43952.35 43932 104.6113 4597025.8217293 61.67755 2710256.3121995 3310601857 3338 true} End of m5 data
        
    
11.  接著按ctrl-c，結束bbgo run
    
	     ^C[6898]  WARN interrupt
	     [6898]  INFO bbgo shutting down...
	     [6899]  INFO session binance user data stream disconnected
        
    
	    下指令加上-h，可以知道某個指令搭配的參數
    
	     go run ./cmd/bbgo run -h
        
    
	    terminal顯示：
    
         jenny@LAPTOP-HF07DJBU:~/bbgo$ go run ./cmd/bbgo run -h
         run strategies from config file
           
         Usage:
           bbgo run [flags]
           
         Flags:
               --enable-grpc                enable grpc server
               --enable-web-server          legacy option, this is renamed to --enable-webserver
               --enable-webserver           enable webserver
               --grpc-bind string           grpc server binding (default ":50051")
           -h, --help                       help for run
               --lightweight                lightweight mode
               --no-compile                 do not compile wrapper binary
               --no-sync                    do not sync on startup
               --setup                      use setup mode
               --totp-account-name string   
               --totp-issuer string         
               --totp-key-url string        time-based one-time password key URL, if defined, it will be used for restoring the otp key
               --webserver-bind string      webserver binding (default ":8080")
           
         Global Flags:
               --binance-api-key string           binance api key
               --binance-api-secret string        binance api secret
               --config string                    config file (default "bbgo.yaml")
               --cpu-profile string               cpu profile
               --debug                            debug mode
               --dotenv string                    the dotenv file you want to load (default ".env.local")
               --log-formatter string             configure log formatter
               --max-api-key string               max api key
               --max-api-secret string            max api secret
               --metrics                          enable prometheus metrics
               --metrics-port string              prometheus http server port (default "9090")
               --no-dotenv                        disable built-in dotenv
               --rollbar-token string             rollbar token
               --slack-channel string             slack trading channel (default "dev-bbgo")
               --slack-error-channel string       slack error channel (default "bbgo-error")
               --slack-token string               slack token
               --telegram-bot-auth-token string   telegram auth token
               --telegram-bot-token string        telegram bot token from bot father
        
    
	    所以其實bbgo run指令有一個`--no-sync`可以用  
	    sync過一次之後，應該不用每次都sync  
	    所以可以這樣下指令：
    
	     cd ~/bbgo 
	     go run ./cmd/bbgo run --config ~/bbgo/config/smacross.yaml --no-sync
        
    
	    每五分鐘，應該就可以看到會吐出一筆K棒資料：
    
	     {0 binance BTCUSDT 2023-12-07 14:45:00 +0800 CST 2023-12-07 14:49:59.999 +0800 CST 5m 43930.35 43862 43930.35 43862 108.48748 4761001.2754831 45.88704 2013538.5591972 3310620281 4776 true} End of m5 data
        
    
12.  那如果用bbgo的回測來跑smacross策略呢？  
    那我必須先在策略設定檔(~/bbgo/config/smacross.yaml)去增加backtest區塊
    
	     backtest:
	      startTime: "2022-01-01"
	      endTime: "2022-03-01"
	      symbols:
	      - BTCUSDT
	      sessions: [binance]
	      # syncSecKLines: true
	      accounts:
	        binance:
	          makerFeeRate: 0.0%
	          takerFeeRate: 0.075%
	          balances:
	            BTC: 0.0
	            USDT: 10_000.0
        
    
	    然後執行bbgo backtest指令(如果執行後沒K線資料，請往下繼續看如何sync K線資料)
    
	     cd ~/bbgo
	     go run ./cmd/bbgo backtest --config ~/bbgo/config/smacross.yaml
        
    
	    接著terminal會出現startTime與endTime之間，所有m5的K棒資料。比較特別的是，bbgo的回測系統還會多吐一個m1的K棒資料，因為bbgo的回測是以m1的K線來進行
    
         ....................................................................
         {209883 binance BTCUSDT 2022-02-28 23:50:00 +0000 UTC 2022-02-28 23:50:59.999 +0000 UTC 1m 43216.18 43220.25 43240 43207.08 27.67119 1195985.4666678 16.92496 731513.3896009 0 0 true} End of m5 data
         {209884 binance BTCUSDT 2022-02-28 23:51:00 +0000 UTC 2022-02-28 23:51:59.999 +0000 UTC 1m 43220.26 43247.04 43300 43193.52 112.55121 4867321.6908315 90.18516 3900348.4362712 0 0 true} End of m5 data
         {209885 binance BTCUSDT 2022-02-28 23:52:00 +0000 UTC 2022-02-28 23:52:59.999 +0000 UTC 1m 43247.03 43195.49 43247.04 43195.49 39.90645 1724993.9675383 13.99468 604808.4974494 0 0 true} End of m5 data
         {209886 binance BTCUSDT 2022-02-28 23:53:00 +0000 UTC 2022-02-28 23:53:59.999 +0000 UTC 1m 43195.49 43191.73 43228 43177 27.9691 1208069.9812698 12.68336 547840.6760792 0 0 true} End of m5 data
         {209887 binance BTCUSDT 2022-02-28 23:54:00 +0000 UTC 2022-02-28 23:54:59.999 +0000 UTC 1m 43191.72 43193.52 43201.92 43171.63 35.24152 1521902.6165683 10.20492 440734.3986899 0 0 true} End of m5 data
         {235811 binance BTCUSDT 2022-02-28 23:50:00 +0000 UTC 2022-02-28 23:54:59.999 +0000 UTC 5m 43216.18 43193.52 43300 43171.63 243.33947 10518273.7228757 143.99308 6225245.3980906 0 0 true} End of m5 data
         {209888 binance BTCUSDT 2022-02-28 23:55:00 +0000 UTC 2022-02-28 23:55:59.999 +0000 UTC 1m 43193.52 43193.57 43213.23 43186.93 35.77888 1545684.6202191 10.20424 440802.9358222 0 0 true} End of m5 data
         {209889 binance BTCUSDT 2022-02-28 23:56:00 +0000 UTC 2022-02-28 23:56:59.999 +0000 UTC 1m 43193.56 43149.02 43193.57 43114.75 52.81375 2278620.1326194 15.95052 688070.8600808 0 0 true} End of m5 data
         {209890 binance BTCUSDT 2022-02-28 23:57:00 +0000 UTC 2022-02-28 23:57:59.999 +0000 UTC 1m 43149.02 43162.85 43162.86 43129.88 37.29367 1609075.9717301 12.26873 529373.2192626 0 0 true} End of m5 data
         {209891 binance BTCUSDT 2022-02-28 23:58:00 +0000 UTC 2022-02-28 23:58:59.999 +0000 UTC 1m 43162.86 43205.75 43209.48 43151.41 41.5637 1795043.8670604 23.99486 1036377.0488003 0 0 true} End of m5 data
         [0019]  WARN session has no BTCUSDT trades
         BACK-TEST REPORT
         ===============================================
         START TIME: Sat, 01 Jan 2022 08:00:00 CST
         END TIME: Tue, 01 Mar 2022 08:00:00 CST
         INITIAL TOTAL BALANCE: BalanceMap[BTC: 0, USDT: 10000]
         FINAL TOTAL BALANCE: BalanceMap[USDT: 10000, BTC: 0]
        
    
	    回測會以111115分鐘，這樣的順序吐資料，不斷重複：
    
        1m(K棒資料)
        1m(K棒資料)
        1m(K棒資料)
        1m(K棒資料)
        1m(K棒資料)
        5m(K棒資料)
        
    
	    最後再出一個回測報告：  
	    那因為策略目前尚未撰寫下單的邏輯，所以餘額沒有任何變化，也沒有其他回測相關的數據可以顯示
    
	     BACK-TEST REPORT 
	     =============================================== 
	     START TIME: Sat, 01 Jan 2022 08:00:00 CST 
	     END TIME: Tue, 01 Mar 2022 08:00:00 CST 
	     INITIAL TOTAL BALANCE: BalanceMap[BTC: 0, USDT: 10000] 
	     FINAL TOTAL BALANCE: BalanceMap[USDT: 10000, BTC: 0]
        
    
13.  問題：smacross.yaml裡，交易對改成BNBUSDT，再執行回測指令，發現BNB不會吐m1m5的K棒資料，去資料庫檢查發現根本沒有K線資料，因此不可能成功回測策略
    
	    解法：  
	    bbgo backtest指令加上–sync，就可以「先sync、後回測」
    
         go run ./cmd/bbgo backtest --sync --config ~/bbgo/config/smacross.yaml
        
    
	    檢測：如何查詢某個交易對是否新增K線資料到資料庫？
    
	    方法一：查目前有哪些幣種已新增過K線資料–幣安交易所
    
         SELECT DISTINCT(symbol) FROM binance_klines;
        
    
	    方法二：查某個交易對的K線資料是否新增成功–幣安交易所
    
         SELECT * FROM binance_klines WHERE symbol = 'XRPUSDT'
        
    
	    smacross.yaml裡的幣改成XRPUSDT，可以改成任何其他想回測與同步K線資料的幣，有兩處要改
    
	     exchangeStrategies:
	     - on: binance
	       smacross:
	         symbol: XRPUSDT
	         interval: 5m
	             
	     backtest:
	       startTime: "2022-01-01"
	       endTime: "2022-03-01"
	       symbols:
	       - XRPUSDT
	       sessions: [binance]
	       # syncSecKLines: true
	       accounts:
	         binance:
	           makerFeeRate: 0.0%
	           takerFeeRate: 0.075%
	           balances:
	             BTC: 0.0
	             USDT: 10_000.0
        
    

> Written with [StackEdit](https://stackedit.io/).