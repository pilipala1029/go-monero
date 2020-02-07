package walletrpc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubAddress(t *testing.T) {

	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	addr, err := rpccl.GetSubAddress("label1")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("addr: %v", addr)

}

func TestGetHeight(t *testing.T) {

	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	height, err := rpccl.GetHeight()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("addr: %v", height)

}

func TestGetBalance(t *testing.T) {
	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	balance, unlocked, err := rpccl.GetBalance()
	assert.NoError(t, err)

	t.Logf("balance: %v unlocked: %v \n", balance, unlocked)
	// // 1 XMR
	// assert.Equal(t, uint64(1000000000000), balance)
	// // 10 XMR
	// assert.Equal(t, uint64(10000000000000), unlocked)

}

func TestIncomingTransfers(t *testing.T) {

	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	transfers, err := rpccl.IncomingTransfers(TransferAll)

	if err != nil {
		t.Fatal(err)
	}

	for k, v := range transfers {
		t.Logf("k: %v v: %v \n", k, v)
	}

}

func TestGetTransfers(t *testing.T) {

	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	req := GetTransfersRequest{
		In: true,
		// Out            bool   `json:"out"`
		// Pending        bool   `json:"pending"`
		// Failed         bool   `json:"failed"`
		// Pool           bool   `json:"pool"`
		// FilterByHeight bool   `json:"filter_by_height"`
		// MinHeight      uint64 `json:"min_height"`
		// MaxHeight      uint64 `json:"max_height"`
	}

	transfers, err := rpccl.GetTransfers(req)

	if err != nil {
		t.Fatal(err)
	}

	for k, v := range transfers.In {
		t.Logf("k: %v v: %v \n", k, v)
	}
}

func TestTransfer(t *testing.T) {
	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	dest := Destination{
		Amount:  1000000000,
		Address: "4328wJW8epC3Kaq9MsUmQSDoquUAWBW3MCf2ZqjEiLZTHUKrkCnnv1CQG3rzxEsNRn4JJ6CGrof2AATrCgmc3GKrM3xhGjD",
	}

	req := TransferRequest{
		Destinations: []Destination{dest},
		GetTxHex:     true,
		GetTxKey:     true,
	}

	t.Log("before transfer----------------------")

	res, err := rpccl.Transfer(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("after transfer------------------------")

	t.Logf("Fee:  %v   txhash: %v  txKey: %v  txBlob: %v \n", res.Fee, res.TxHash, res.TxKey, res.TxBlob)

}

//601b8465cb4da8cf23f99c2b1086c9dd1f8323ccb716068d67703269b7656f9e
func TestGetTransferByTxID(t *testing.T) {

	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	transfer, err := rpccl.GetTransferByTxID("10364eba232b655e82f4e36ce50d99450808eb17efcec7da4e935944a1e5b291")
	//2e4c2edc783287fcb871fa51beb9d854362d7a97ad7b9c5fd794267a75c3cd1c
	if err != nil {
		t.Fatal(err)
	}

	t.Log("after transfer------------------------")

	tranAsJson, err := json.MarshalIndent(transfer, "", "		")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("trnsfer:  %s \n", tranAsJson)
}

func TestCreateNewAccount(t *testing.T) {
	rpccl := New(Config{
		Address: "http://47.245.28.129:18083/json_rpc",
	})

	addr, err := rpccl.CreateNewAccount("account label1")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("addr: %v", addr)

}

func TestClient(t *testing.T) {

	testClientGetAddress(t)
	testClientGetBalance(t)
}

func testClientGetAddress(t *testing.T) {
	//
	// server setup
	sv0 := basicTestServer([]testfn{
		func(method string, params *json.RawMessage, w http.ResponseWriter, r *http.Request) bool {
			if method == "getaddress" {
				r0 := struct {
					Address string `json:"address"`
				}{
					"45eoXYNHC4LcL2Hh42T9FMPTmZHyDEwDbgfBEuNj3RZUek8A4og4KiCfVL6ZmvHBfCALnggWtHH7QHF8426yRayLQq7MLf5",
				}
				writerpcResponseOK(&r0, w)
				return true
			}
			return false
		},
	})
	defer sv0.Close()
	//
	// test starts here
	rpccl := New(Config{
		Address: sv0.URL + "/json_rpc",
	})
	addr, err := rpccl.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "45eoXYNHC4LcL2Hh42T9FMPTmZHyDEwDbgfBEuNj3RZUek8A4og4KiCfVL6ZmvHBfCALnggWtHH7QHF8426yRayLQq7MLf5", addr)
}

func testClientGetBalance(t *testing.T) {
	//
	// server setup
	sv0 := basicTestServer([]testfn{
		func(method string, params *json.RawMessage, w http.ResponseWriter, r *http.Request) bool {
			if method == "getbalance" {
				r0 := struct {
					Balance  uint64 `json:"balance"`
					Unlocked uint64 `json:"unlocked_balance"`
				}{
					1e12,
					1e13,
				}
				writerpcResponseOK(&r0, w)
				return true
			}
			return false
		},
	})
	defer sv0.Close()
	//
	// test starts here
	rpccl := New(Config{
		Address: sv0.URL + "/json_rpc",
	})
	balance, unlocked, err := rpccl.GetBalance()
	assert.NoError(t, err)
	// 1 XMR
	assert.Equal(t, uint64(1000000000000), balance)
	// 10 XMR
	assert.Equal(t, uint64(10000000000000), unlocked)
}

//TODO: write more server stubs
//
//

type testfn = func(method string, params *json.RawMessage, w http.ResponseWriter, r *http.Request) bool

func basicTestServer(tests []testfn) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/json_rpc" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var c clientRequest
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		for _, v := range tests {
			if v(c.Method, c.Params, w, r) {
				return
			}
		}
		// return method not found
		writerpcResponseError(ErrUnknown, "test this in curl with the real rpc", w)
	}))
}

func writerpcResponseOK(result interface{}, w http.ResponseWriter) {
	r := &clientResponse{
		Version: "2.0",
		Result:  result,
	}
	v, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

func writerpcResponseError(code ErrorCode, message string, w http.ResponseWriter) {
	r := &clientResponse{
		Version: "2.0",
		Result:  nil,
		Error: &WalletError{
			Code:    code,
			Message: message,
		},
	}
	v, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// clientRequest represents a JSON-RPC request received by the server.
type clientRequest struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params *json.RawMessage `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}
