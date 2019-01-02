package unittests

import (
	"encoding/json"
	"github.com/docker/docker/pkg/testutil/assert"
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"testing"
)

/*
var eosioToken = common.AccountName(common.N("eosio.token"))
var DEFAULT_EXPIRATION_DELTA uint32 = 6
var DEFAULT_BILLED_CPU_TIME_US uint32 = 2000*/

type CurrencyTester struct {
	abiSer           abi_serializer.AbiDef
	eosioToken       string
	validatingTester *ValidatingTester
}

/*func initBaseTester() *BaseTester {
	bt := newBaseTester(true, SPECULATIVE)
	return bt
}*/

func NewCurrencyTester() *CurrencyTester {
	ct := &CurrencyTester{}
	ct.eosioToken = eosioToken.String()
	bt := newValidatingTester(true, SPECULATIVE)
	//ct.abiSer = abi_serializer.NewABI()
	bt.CreateDefaultAccount(common.N(ct.eosioToken))
	ct.validatingTester = bt
	//bt.SetCode2(common.N("eosio"),eosioTokenWast)
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	bt.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	bt.SetAbi(common.AccountName(eosioToken), abi, nil)
	accountName := common.AccountName(common.N(ct.eosioToken))

	createData := common.Variants{
		"issuer":         eosioToken,
		"maximum_supply": "1000000000.0000 EOS",
		"can_freeze":     0,
		"can_recall":     0,
		"can_whitelist":  0,
	}
	cn := common.ActionName(common.N("create"))
	result := ct.PushAction(&accountName, &cn, &createData)
	log.Info("NewCurrencyTester push action issue:%v", result.BlockNum)
	data := common.Variants{
		"to":       eosioToken,
		"quantity": "1000000.0000 EOS",
		"memo":     "test",
	}
	in := common.ActionName(common.N("issue"))
	result = ct.PushAction(&accountName, &in, &data)

	log.Info("NewCurrencyTester push action issue::%v", result.BlockNum)
	signedBlock := bt.DefaultProduceBlock()
	log.Info("NewCurrencyTester produceBlock result:%v", signedBlock.Producer.String())
	ct.validatingTester = bt
	return ct
}

func (c *CurrencyTester) PushAction(signer *common.AccountName, name *common.ActionName, data *common.Variants) *types.TransactionTrace {
	action := types.Action{eosioToken, *name, []types.PermissionLevel{{*signer, common.DefaultConfig.ActiveName}}, nil}
	acnt := c.validatingTester.Control.GetAccount(eosioToken)
	a := acnt.GetAbi()
	buf, _ := json.Marshal(data)
	action.Data, _ = a.EncodeAction(*name, buf)
	trx := types.NewSignedTransactionNil()
	trx.Actions = append(trx.Actions, &action)

	c.validatingTester.SetTransactionHeaders(&trx.Transaction, c.validatingTester.DefaultBilledCpuTimeUs, 0)
	key := c.validatingTester.getPrivateKey(*signer, "active")
	chainID := c.validatingTester.Control.GetChainId()
	trx.Sign(&key, &chainID)
	return c.validatingTester.PushTransaction(trx, common.MaxTimePoint(), c.validatingTester.DefaultBilledCpuTimeUs)
}

func (ct *CurrencyTester) GetBalance(account *common.AccountName) *common.Asset {
	symbol := common.Symbol{Precision: 4, Symbol: "EOS"}
	actionName := common.N(ct.eosioToken)
	asset := ct.validatingTester.GetCurrencyBalance(&actionName, &symbol, account)
	return &asset
}

func (c *CurrencyTester) Transfer(from *common.AccountName, to *common.AccountName, quantity string, memo string) *types.TransactionTrace {
	q := common.N("transfer")
	data := common.Variants{
		"from":     from,
		"to":       to,
		"quantity": quantity,
		"memo":     memo}
	trace := c.PushAction(from, &q, &data)
	c.validatingTester.DefaultProduceBlock()
	return trace
}

func TestBootstrap(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		s := "1000000.0000 EOS"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		actionName := common.N(ct.eosioToken)
		accountName := common.N(ct.eosioToken)

		actual := ct.validatingTester.GetCurrencyBalance(&actionName, &expected.Symbol, &accountName)
		assert.Equal(t, expected, actual)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestCurrencyTransfer(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		alice := common.N("alice")
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		accountName := common.AccountName(eosioToken)
		actionName := common.N("transfer")
		data := common.Variants{
			"from":     eosioToken,
			"to":       "alice",
			"quantity": "100.0000 EOS",
			"memo":     "fund Alice"}
		trace := ct.PushAction(&accountName, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock()

		s := "100.0000 EOS"
		expected := common.Asset{}.FromString(&s)
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestDuplicateTransfer(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		accountName := common.AccountName(common.N(ct.eosioToken))
		actionName := common.ActionName(common.N("transfer"))
		asset := common.Asset{}
		s := "100.0000 EOS"
		expected := asset.FromString(&s)
		alice := common.N("alice")
		data := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": "100.0000 EOS",
			"memo":     "fund Alice"}
		trace := ct.PushAction(&accountName, &actionName, &data)

		try.Try(func() {
			ct.PushAction(&accountName, &actionName, &data)
		}).Catch(func(e error) {
			assert.Error(t, e, "Duplicate transaction")
		})

		ct.validatingTester.DefaultProduceBlock()
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestAddTransfer(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		alice := common.N("alice")
		actionName := common.ActionName(common.N("transfer"))
		asset := common.Asset{}
		s := "100.0000 EOS"
		expected := asset.FromString(&s)
		data := common.Variants{
			"from":     eosioToken,
			"to":       "alice",
			"quantity": s,
			"memo":     "fund Alice"}
		trace := ct.PushAction(&eosioToken, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock() //
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)

		asset2 := common.Asset{}
		st := "110.0000 EOS"

		exp := asset2.FromString(&st)
		transferData := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": "10.0000 EOS",
			"memo":     "fund Alice"}

		try.Try(func() {
			ct.PushAction(&eosioToken, &actionName, &transferData)
		}).Catch(func(e exception.TxDuplicate) {
			log.Error(e.String())
		})

		ct.validatingTester.DefaultProduceBlock()

		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), exp)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestOverspend(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
		accountName := common.AccountName(common.N(ct.eosioToken))
		actionName := common.ActionName(common.N("transfer"))
		alice := common.AccountName(common.N("alice"))
		asset := common.Asset{}
		s := "100.0000 EOS"
		expected := asset.FromString(&s)
		data := common.Variants{
			"from":     ct.eosioToken,
			"to":       "alice",
			"quantity": s,
			"memo":     "fund Alice"}
		trace := ct.PushAction(&accountName, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock() //
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)

		s2 := "101.0000 EOS"
		//expected2 := asset2.FromString(&s2)
		data2 := common.Variants{
			"from":     "alice",
			"to":       "bob",
			"quantity": s2,
			"memo":     "fund Alice"}
		returning := false
		try.Try(func() {
			ct.PushAction(&alice, &actionName, &data2)
			bob := common.AccountName(common.N("bob"))
			ct.validatingTester.DefaultProduceBlock() //
			tt := "0.0000 EOS"
			assert.Equal(t, *ct.GetBalance(&alice), expected)
			assert.Equal(t, *ct.GetBalance(&bob), asset.FromString(&tt))
		}).Catch(func(e exception.EosioAssertMessageException) {
			if inString(e.DetailMessage(), "overdrawn balance") {
				returning = true
			}
		}).End()
		assert.Equal(t, true, returning)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestFullspend(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
		actionName := common.ActionName(common.N("transfer"))
		alice := common.AccountName(common.N("alice"))
		zero := "0.0000 EOS"
		val := "100.0000 EOS"
		data := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": val,
			"memo":     "all in! Alice"}
		trace := ct.PushAction(&eosioToken, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock()
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		z := common.Asset{}.FromString(&zero)

		assert.Equal(t, *ct.GetBalance(&alice), common.Asset{}.FromString(&val))
		bob := common.AccountName(common.N("bob"))
		data2 := common.Variants{
			"from":     alice,
			"to":       bob,
			"quantity": "100.0000 EOS",
			"memo":     "all in! Alice"}
		trace2 := ct.PushAction(&alice, &actionName, &data2)
		log.Info("trace2 id:%d", trace2.ID)
		ct.validatingTester.DefaultProduceBlock()
		s := "100.0000 EOS"
		expected := common.Asset{}.FromString(&s)
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace2.ID))
		assert.Equal(t, *ct.GetBalance(&bob), expected)
		assert.Equal(t, *ct.GetBalance(&alice), z)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestSymbol(t *testing.T) {
	{
		dollar := common.Symbol{Precision: 2, Symbol: "DLLR"}
		sy := "2,DLLR"
		dollar2 := common.Symbol{}.FromString(&sy)
		assert.Equal(t, dollar2, dollar)
		assert.Equal(t, dollar.Decimals(), uint8(2))
		assert.Equal(t, dollar.Name(), "DLLR")
		assert.Equal(t, dollar.Valid(), true)
	}
	{
		def := CORE_SYMBOL
		assert.Equal(t, def.Decimals(), uint8(4))
		assert.Equal(t, def.Name(), CORE_SYMBOL_NAME)
	}
	{
		returning := false
		try.Try(func() {
			sy := ""
			common.Symbol{}.FromString(&sy)
		}).Catch(func(e exception.SymbolTypeException) {
			returning = true
		}).End()
		if returning {
			assert.Equal(t, true, returning)
		}
	}
	{
		returning := false
		try.Try(func() {
			sy := "RND"
			common.Symbol{}.FromString(&sy)
		}).Catch(func(e exception.SymbolTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		returning := false
		try.Try(func() {
			sy := "6,EoS"
			common.Symbol{}.FromString(&sy)
		}).Catch(func(e exception.SymbolTypeException) {

			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		str := "10 CUR"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(10))
		assert.Equal(t, asset.Decimals(), uint8(0))
		assert.Equal(t, asset.Symbol.Symbol, "CUR")
	}
	{
		returning := false
		try.Try(func() {
			str := "10CUR"
			common.Asset{}.FromString(&str)
		}).Catch(func(e exception.AssetTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		returning := false
		try.Try(func() {
			str := "10. CUR"
			common.Asset{}.FromString(&str)
		}).Catch(func(e exception.AssetTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		returning := false
		try.Try(func() {
			str := "10"
			common.Asset{}.FromString(&str)
		}).Catch(func(e exception.AssetTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		str := "-001000000.00010 CUR"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(-100000000010))
		assert.Equal(t, asset.Decimals(), uint8(5))
		assert.Equal(t, asset.Symbol.Symbol, "CUR")
		assert.Equal(t, asset.String(), "-1000000.00010 CUR")
	}
	{
		str := "-000000000.00100 CUR"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(-100))
		assert.Equal(t, asset.Decimals(), uint8(5))
		assert.Equal(t, asset.Symbol.Symbol, "CUR")
		assert.Equal(t, asset.String(), "-0.00100 CUR")
	}

	{
		str := "-0.0001 PPP"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(-1))
		assert.Equal(t, asset.Decimals(), uint8(4))
		assert.Equal(t, asset.Symbol.Symbol, "PPP")
		assert.Equal(t, asset.String(), "-0.0001 PPP")
	}
}

func TestProxy(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.ProduceBlocks(2, false)
		alice := common.N("alice")
		proxy := common.N("proxy")
		ct.validatingTester.CreateAccounts([]common.AccountName{alice, proxy}, false, true)
		ct.validatingTester.DefaultProduceBlock()
		wasmName := "test_contracts/proxy.wasm"
		code, _ := ioutil.ReadFile(wasmName)
		ct.validatingTester.SetCode(proxy, code, nil)
		abiName := "test_contracts/proxy.abi"
		abi, _ := ioutil.ReadFile(abiName)

		ct.validatingTester.SetAbi(proxy, abi, nil)

		{
			act := types.Action{}
			act.Account = proxy
			act.Name = common.N("setowner")
			act.Authorization = []types.PermissionLevel{{common.N("alice"), common.DefaultConfig.ActiveName}}
			data := common.Variants{
				"owner": alice,
				"delay": 10,
			}

			//trace := ct.validatingTester.PushAction(&alice, &act.Name, &data)
			trace := ct.validatingTester.PushAction2(&proxy, &act.Name, alice, &data, ct.validatingTester.DefaultExpirationDelta, 0)
			ct.validatingTester.ProduceBlocks(1, false)
			assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
			ct.validatingTester.ProduceBlocks(1, false)
		}

		{
			act1 := types.Action{}
			act1.Account = eosioToken
			act1.Name = common.N("transfer")
			act1.Authorization = []types.PermissionLevel{{eosioToken, common.DefaultConfig.ActiveName}}
			data1 := common.Variants{
				"from":     eosioToken,
				"to":       proxy,
				"quantity": "5.0000 EOS",
				"memo":     "fund Proxy",
			}

			//trace1 := ct.PushAction(&eosioToken, &act1.Name, &data1)
			trace1 := ct.validatingTester.PushAction2(&eosioToken, &act1.Name, eosioToken, &data1, ct.validatingTester.DefaultExpirationDelta, 0)
			tt := ct.validatingTester.Control.HeadBlockTime().TimeSinceEpoch().Count()
			expectedDelivery := tt + common.Seconds(10).Count()
			s := "5.0000 EOS"
			expected := common.Asset{}.FromString(&s)
			s1 := "0.0000 EOS"
			expected1 := common.Asset{}.FromString(&s1)
			for ct.validatingTester.Control.HeadBlockTime().TimeSinceEpoch().Count() < expectedDelivery {
				ct.validatingTester.ProduceBlocks(1, false)
				assert.Equal(t, *ct.GetBalance(&proxy), expected)
				assert.Equal(t, *ct.GetBalance(&alice), expected1)
			}

			ct.validatingTester.ProduceBlocks(1, false)
			assert.Equal(t, *ct.GetBalance(&proxy), expected1)
			assert.Equal(t, *ct.GetBalance(&alice), expected)
			assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace1.ID))
			ct.validatingTester.close()
		}

	}).FcLogAndRethrow().End()
}