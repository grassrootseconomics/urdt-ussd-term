package nats

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	dataserviceapi "github.com/grassrootseconomics/ussd-data-service/pkg/api"
	"git.defalsify.org/vise.git/db"
	memdb "git.defalsify.org/vise.git/db/mem"
	"git.grassecon.net/urdt/ussd/common"
	"git.grassecon.net/urdt/ussd/config"
	"git.grassecon.net/urdt/ussd/models"
	"git.grassecon.net/term/lookup"
	"git.grassecon.net/term/event"
)

const (
	txBlock = 42
	tokenAddress = "0x765DE816845861e75A25fCA122bb6898B8B1282a"
	tokenSymbol = "FOO"
	tokenName = "Foo Token"
	tokenDecimals = 6
	txValue = 1337
	tokenBalance = 362436
	txTimestamp = 1730592500
	txHash = "0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	sinkAddress = "0xb42C5920014eE152F2225285219407938469BBfA"
	aliceChecksum = "0xeae046BF396e91f5A8D74f863dC57c107c8a4a70"
	bobChecksum = "0xB3117202371853e24B725d4169D87616A7dDb127"
	aliceSession = "5553425"
)

// TODO: jetstream, would have been nice of you to provide an easier way to make a mock msg
type testMsg struct {
	data []byte
}

func(m *testMsg) Ack() error {
	return nil
}

func(m *testMsg) Nak() error {
	return nil
}

func(m *testMsg) NakWithDelay(time.Duration) error {
	return nil
}

func(m *testMsg) Data() []byte {
	return m.data
}

func(m *testMsg) Reply() string {
	return ""
}

func(m *testMsg) Subject() string {
	return ""
}

func(m *testMsg) Term() error {
	return nil
}

func(m *testMsg) TermWithReason(string) error {
	return nil
}

func(m *testMsg) DoubleAck(ctx context.Context) error {
	return nil
}

func(m *testMsg) Headers() nats.Header {
	return nats.Header{}
}

func(m *testMsg) InProgress() error {
	return nil
}

func(m *testMsg) Metadata() (*jetstream.MsgMetadata, error) {
	return nil, nil
}

type mockApi struct {
}

func(m mockApi) CheckBalance(ctx context.Context, publicKey string) (*models.BalanceResult, error) {
	return nil, nil
}

func(m mockApi) CreateAccount(ctx context.Context) (*models.AccountResult, error) {
	return nil, nil
}

func(m mockApi) TrackAccountStatus(ctx context.Context, publicKey string) (*models.TrackStatusResult, error) {
	return nil, nil
}

func(m mockApi) FetchVouchers(ctx context.Context, publicKey string) ([]dataserviceapi.TokenHoldings, error) {
	logg.DebugCtxf(ctx, "mockapi fetchvouchers", "key", publicKey)
	return []dataserviceapi.TokenHoldings{
		dataserviceapi.TokenHoldings{
			ContractAddress: tokenAddress,
			TokenSymbol: tokenSymbol,
			TokenDecimals: strconv.Itoa(tokenDecimals),
			Balance: strconv.Itoa(tokenBalance),
		},
	}, nil
}

func(m mockApi) FetchTransactions(ctx context.Context, publicKey string) ([]dataserviceapi.Last10TxResponse, error) {
	logg.DebugCtxf(ctx, "mockapi fetchtransactions", "key", publicKey)
	return []dataserviceapi.Last10TxResponse{
		dataserviceapi.Last10TxResponse{
			Sender: aliceChecksum,
			Recipient: bobChecksum,
			TransferValue: strconv.Itoa(txValue),
			ContractAddress: tokenAddress,
			TxHash: txHash,
			DateBlock: time.Unix(txTimestamp, 0),
			TokenSymbol: tokenSymbol,
			TokenDecimals: strconv.Itoa(tokenDecimals),
		},
	}, nil
}

func(m mockApi) VoucherData(ctx context.Context, address string) (*models.VoucherDataResult, error) {
	return &models.VoucherDataResult{
		TokenSymbol: tokenSymbol,
		TokenName: tokenName,
		TokenDecimals: strconv.Itoa(tokenDecimals),
		SinkAddress: sinkAddress,
	}, nil
}

func TestHandleMsg(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	lookup.Api = mockApi{}

	ctx := context.Background()
	userDb := memdb.NewMemDb()
	err = userDb.Connect(ctx, "")
	if err != nil {
		panic(err)
	}

	alice, err := common.NormalizeHex(aliceChecksum)
	if err != nil {
		t.Fatal(err)
	}

	userDb.SetSession(alice)
	userDb.SetPrefix(db.DATATYPE_USERDATA)
	err = userDb.Put(ctx, common.PackKey(common.DATA_PUBLIC_KEY_REVERSE, []byte{}), []byte(aliceSession))
	if err != nil {
		t.Fatal(err)
	}

	sub := NewNatsSubscription(userDb)

	data := fmt.Sprintf(`{
	"block": %d,
	"contractAddress": "%s",
	"success": true,
	"timestamp": %d,
	"transactionHash": "%s",
	"transactionType": "TOKEN_TRANSFER",
	"payload": {
		"from": "%s",
		"to": "%s",
		"value": "%d"
	}
}`, txBlock, tokenAddress, txTimestamp, txHash, aliceChecksum, bobChecksum, txValue)
	msg := &testMsg{
		data: []byte(data),
	}
	sub.handleEvent(msg)

	store := common.UserDataStore{
		Db: userDb,
	}
	v, err := store.ReadEntry(ctx, aliceSession, common.DATA_ACTIVE_SYM)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(v, []byte(tokenSymbol)) {
		t.Fatalf("expected '%s', got %s", tokenSymbol, v)
	}

	v, err = store.ReadEntry(ctx, aliceSession, common.DATA_ACTIVE_BAL)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(v, []byte(strconv.Itoa(tokenBalance))) {
		t.Fatalf("expected '%d', got %s", tokenBalance, v)
	}

	v, err = store.ReadEntry(ctx, aliceSession, common.DATA_TRANSACTIONS)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(v, []byte("abcdef")) {
		t.Fatal("no transaction data")
	}

	userDb.SetPrefix(event.DATATYPE_USERSUB)
	userDb.SetSession(aliceSession)
	k := append([]byte("vouchers"), []byte("sym")...)
	v, err = userDb.Get(ctx, k)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(v, []byte(fmt.Sprintf("1:%s", tokenSymbol))) {
		t.Fatalf("expected '1:%s', got %s", tokenSymbol, v)
	}
}
