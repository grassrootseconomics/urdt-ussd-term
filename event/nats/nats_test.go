package nats

import (
	"context"
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
	"git.grassecon.net/term/event"
)

func init() {
}

const (
	aliceChecksum = "0xeae046BF396e91f5A8D74f863dC57c107c8a4a70"
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
			ContractAddress: "0xeE0A29AE1BB7a033c8277C04780c4aBcf4388E93",
			TokenSymbol: "FOO",
			TokenDecimals: "6",
			Balance: "362436",
		},
	}, nil
}

func(m mockApi) FetchTransactions(ctx context.Context, publicKey string) ([]dataserviceapi.Last10TxResponse, error) {
	return nil, nil
}

func TestHandleMsg(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	event.Api = mockApi{}

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

	aliceSession := "5553425"
	userDb.SetSession(alice)
	userDb.SetPrefix(db.DATATYPE_USERDATA)
	err = userDb.Put(ctx, common.PackKey(common.DATA_PUBLIC_KEY_REVERSE, []byte{}), []byte(aliceSession))
	if err != nil {
		t.Fatal(err)
	}

	sub := NewNatsSubscription(userDb)

	data := `{
	"block": 42,
	"contractAddress": "0x765DE816845861e75A25fCA122bb6898B8B1282a",
	"success": true,
	"timestamp": 1730592500,
	"transactionHash": "0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
	"transactionType": "TOKEN_TRANSFER",
	"payload": {
		"from": "0xeae046BF396e91f5A8D74f863dC57c107c8a4a70",
		"to": "B3117202371853e24B725d4169D87616A7dDb127",
		"value": "1337"
	}
}`
	msg := &testMsg{
		data: []byte(data),
	}
	sub.handleEvent(msg)
}
