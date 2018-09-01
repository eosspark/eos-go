package chain

import (
	"github.com/eosspark/eos-go/chain/database"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
	"fmt"
)

type DBReadMode int8

const (
	SPECULATIVE = DBReadMode(iota)
	HEADER      //HEAD
	READONLY
	IRREVERSIBLE
)

type HandlerKey struct {
	handKey map[common.AccountName]common.AccountName
}

type applyCon struct {
	handlerKey map[common.AccountName]common.AccountName
	applyContext types.ApplyContext
}

//apply_context
type ApplyHandler struct {
	applyHandler map[common.AccountName]applyCon
	scopeName common.AccountName
}

type Controller struct {
	db 					  eosiodb.Database
	dbsession             *eosiodb.Session
	reversibledb		  eosiodb.Database
	reversibleBlocks      *eosiodb.Session
	blog                  string //TODO
	pending               *types.PendingState
	head                  types.BlockState
	forkDB                database.ForkDatabase
	wasmif                string //TODO
	resourceLimist        types.ResourceLimitsManager
	authorization         string //TODO AuthorizationManager
	config                string //TODO	Config
	chainID               common.ChainIDType
	rePlaying             bool
	replayHeadTime        common.Tstamp //optional<common.Tstamp>
	readMode              DBReadMode
	inTrxRequiringChecks  bool	//if true, checks that are normally skipped on replay (e.g. auth checks) cannot be skipped
	subjectiveCupLeeway   common.Tstamp //optional<common.Tstamp>
	handlerKey            HandlerKey
	applyHandlers         ApplyHandler
	unappliedTransactions map[[4]uint64]types.TransactionMetadata
}

func NewController() *Controller {

	db, err := eosiodb.NewDatabase("./", "shared_memory.bin", true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:", err)
		return  nil
	}
	defer db.Close()

	session, err := db.Start_Session()

	if err != nil {
		log.Debug("db start session is error detail:", err.Error(),session)
		return  nil
	}
	defer session.Undo()

	session.Commit()
	return &Controller{inTrxRequiringChecks : false}
}

func (self *Controller) PopBlock() {

	prev, err := self.forkDB.GetBlock(self.head.Header.Previous)
	if err != nil {
		log.Error("PopBlock GetBlockByID is error,detail:", err)
	}
	var r types.ReversibleBlockObject
	errs := self.reversibleBlocks.Find("NUM", self.head.BlockNum, r)
	if errs != nil {
		log.Error("PopBlock ReversibleBlocks Find is error,detail:", errs)
	}
	if &r != nil {
		self.reversibleBlocks.Remover(&r)
	}

	if self.readMode == SPECULATIVE {
		var trx []types.TransactionMetadata = self.head.Trxs
		step := 0
		for ; step < len(trx); step++ {
			self.unappliedTransactions[trx[step].SignedID] = trx[step]
		}
	}
	self.head = prev
	self.dbsession.Undo() //TODO
}

func newApplyCon(ac types.ApplyContext) *applyCon{
	a :=applyCon{}
	a.applyContext = ac
	return &a
}
func (self *Controller) SetApplayHandler(receiver common.AccountName, contract common.AccountName, action common.AccountName, handler types.ApplyContext) {
	h:=make(map[common.AccountName]common.AccountName)
	h[receiver] = contract
	apply := newApplyCon(handler)
	apply.handlerKey = h
	t := make(map[common.AccountName]applyCon)
	t[receiver]= *apply
	self.applyHandlers = ApplyHandler{t,receiver}
	fmt.Println(self.applyHandlers)
}

func (self *Controller) AbortBlock() {
	if self.pending != nil {
		if self.readMode == SPECULATIVE {
			trx := append(self.pending.PendingBlockState.Trxs)
			step := 0
			for ; step < len(trx); step++ {
				self.unappliedTransactions[trx[step].SignedID] = trx[step]
			}
		}
	}
}

func (self *Controller) StartBlock(when common.BlockTimeStamp,confirmBlockCount uint16,s types.BlockStatus){
	if self.pending !=nil{
		fmt.Println("pending block already exists")
		return
	}
	// defer self.peding.reset()


}

func  Close(db eosiodb.Database,session eosiodb.Session){
	//session.close() 	//db close前关闭session
	db.Close()
}

/*func main(){
	c := new(Controller)

	fmt.Println("asdf",c)
}*/
/*"github.com/eos-go/chain/types".TransactionMetadata)
"github.com/eosspark/eos-go/chain/types".TransactionMetadata*/