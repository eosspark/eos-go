package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/common/arithmetic_types"
	. "github.com/eosspark/eos-go/exception"
)

var IsActiveRc bool

var rcInstance *ResourceLimitsManager

type ResourceLimitsManager struct {
	db database.DataBase `json:"db"`
}

func GetResourceLimitsManager() *ResourceLimitsManager {
	if !IsActiveRc {
		rcInstance = newResourceLimitsManager()
	}
	return rcInstance
}

func newResourceLimitsManager() *ResourceLimitsManager {
	IsActiveRc = true
	//control := GetControllerInstance()
	//db := control.DataBase()
	db, _ := database.NewDataBase(common.DefaultConfig.DefaultStateDirName)
	return &ResourceLimitsManager{db: db}
}

func (r *ResourceLimitsManager) InitializeDatabase() {
	config := entity.NewResourceLimitsConfigObject()
	r.db.Insert(&config)

	state := entity.DefaultResourceLimitsStateObject
	state.VirtualCpuLimit = config.CpuLimitParameters.Max
	state.VirtualNetLimit = config.NetLimitParameters.Max
	r.db.Insert(&state)
}

func (r *ResourceLimitsManager) InitializeAccount(account common.AccountName) {
	bl := entity.ResourceLimitsObject{}
	bl.Owner = account
	r.db.Insert(&bl)

	bu := entity.ResourceUsageObject{}
	bu.Owner = account
	r.db.Insert(&bu)
}

func (r *ResourceLimitsManager) SetBlockParameters(cpuLimitParameters types.ElasticLimitParameters, netLimitParameters types.ElasticLimitParameters) {
	cpuLimitParameters.Validate()
	netLimitParameters.Validate()
	config := entity.DefaultResourceLimitsConfigObject
	r.db.Find("id", config, &config)
	r.db.Modify(&config, func(c entity.ResourceLimitsConfigObject) {
		c.CpuLimitParameters = cpuLimitParameters
		c.NetLimitParameters = netLimitParameters
	})
}

func (r *ResourceLimitsManager) UpdateAccountUsage(account []common.AccountName, timeSlot uint32) { //待定
	config := entity.DefaultResourceLimitsConfigObject
	r.db.Find("id", config, &config)
	usage := entity.ResourceUsageObject{}
	for _, a := range account {
		usage.Owner = a
		r.db.Find("ByOwner", usage, &usage)
		r.db.Modify(&usage, func(bu entity.ResourceUsageObject) {
			bu.NetUsage.Add(0, timeSlot, config.AccountNetUsageAverageWindow)
			bu.CpuUsage.Add(0, timeSlot, config.AccountCpuUsageAverageWindow)
		})
	}
}

func (r *ResourceLimitsManager) AddTransactionUsage(account []common.AccountName, cpuUsage uint64, netUsage uint64, timeSlot uint32) {
	state := entity.DefaultResourceLimitsStateObject
	r.db.Find("id", state, &state)
	config := entity.DefaultResourceLimitsConfigObject
	r.db.Find("id", config, &config)
	for _, a := range account {
		usage := entity.ResourceUsageObject{}
		usage.Owner = a
		r.db.Find("byOwner", usage, &usage)
		var unUsed, netWeight, cpuWeight int64
		r.GetAccountLimits(a, &unUsed, &netWeight, &cpuWeight)
		r.db.Modify(&usage, func(bu entity.ResourceUsageObject) {
			bu.CpuUsage.Add(netUsage, timeSlot, config.AccountNetUsageAverageWindow)
			bu.NetUsage.Add(cpuUsage, timeSlot, config.AccountCpuUsageAverageWindow)
		})

		if cpuWeight >= 0 && state.TotalCpuWeight > 0 {
			windowSize := uint64(config.AccountCpuUsageAverageWindow)
			virtualNetworkCapacityInWindow := arithmeticTypes.MulUint64(state.VirtualCpuLimit, windowSize)
			cpuUsedInWindow := arithmeticTypes.MulUint64(usage.CpuUsage.ValueEx, windowSize)
			cpuUsedInWindow, _ = cpuUsedInWindow.Div(arithmeticTypes.Uint128{0, uint64(common.DefaultConfig.RateLimitingPrecision)})
			userWeight := arithmeticTypes.Uint128{0, uint64(cpuWeight)}
			allUserWeight :=  arithmeticTypes.Uint128{0, state.TotalCpuWeight}

			maxUserUseInWindow := virtualNetworkCapacityInWindow.Mul(userWeight)
			maxUserUseInWindow, _ = maxUserUseInWindow.Div(allUserWeight)
			EosAssert(cpuUsedInWindow.Compare(maxUserUseInWindow) < 1, &TxCpuUsageExceed{},
			"authorizing account %s has insufficient cpu resources for this transaction,\n cpu_used_in_window: %s,\n max_user_use_in_window: %s",
			a, cpuUsedInWindow, maxUserUseInWindow)
		}

		if netWeight >= 0 && state.TotalNetWeight > 0 {
			windowSize := uint64(config.AccountNetUsageAverageWindow)
			virtualNetworkCapacityInWindow := arithmeticTypes.MulUint64(state.VirtualNetLimit, windowSize)
			netUsedInWindow := arithmeticTypes.MulUint64(usage.NetUsage.ValueEx, windowSize)
			netUsedInWindow, _ = netUsedInWindow.Div(arithmeticTypes.Uint128{0, uint64(common.DefaultConfig.RateLimitingPrecision)})
			userWeight := arithmeticTypes.Uint128{0, uint64(cpuWeight)}
			allUserWeight :=  arithmeticTypes.Uint128{0, state.TotalCpuWeight}

			maxUserUseInWindow := virtualNetworkCapacityInWindow.Mul(userWeight)
			maxUserUseInWindow, _ = maxUserUseInWindow.Div(allUserWeight)
			EosAssert(netUsedInWindow.Compare(maxUserUseInWindow) < 1, &TxCpuUsageExceed{},
				"authorizing account %s has insufficient cpu resources for this transaction,\n net_used_in_window: %s,\n max_user_use_in_window: %s",
				a, netUsedInWindow, maxUserUseInWindow)
		}
	}

	r.db.Modify(&state, func(rls entity.ResourceLimitsStateObject) {
		rls.PendingCpuUsage += cpuUsage
		rls.PendingNetUsage += netUsage
	})

}

func (rlm *ResourceLimitsManager) AddPendingRamUsage(account common.AccountName, ramDelta int64) {
	//if ramDelta == 0 {
	//	return
	//}
	//
	//ruo := entity.ResourceUsageObject{}
	//ruo.Owner = account
	//rlm.db.Find("byOwner", &ruo)
	//
	//if ramDelta > 0 && math.MaxUint64-ruo.RamUsage < uint64(ramDelta) {
	//	fmt.Println("error")
	//}
	//if ramDelta < 0 && ruo.RamUsage < uint64(-ramDelta) {
	//	fmt.Println("error")
	//}
	//
	//rlm.db.Modify(&ruo, func(data interface{}) error {
	//	ruo.RamUsage += uint64(ramDelta)
	//	return nil
	//})
}

func (rlm *ResourceLimitsManager) VerifyAccountRamUsage(account common.AccountName) {
	//var ramBytes, netWeight, cpuWeight int64
	//rlm.GetAccountLimits(account, &ramBytes, &netWeight, &cpuWeight)
	//ruo := entity.ResourceUsageObject{}
	//
	//rlm.db.Find("byOwner", &ruo)
	//
	//if ramBytes >= 0 {
	//	if int64(ruo.RamUsage) > ramBytes {
	//		fmt.Println("error")
	//	}
	//}
}

func (rlm *ResourceLimitsManager) GetAccountRamUsage(account common.AccountName) int64 {
	//ruo := entity.ResourceUsageObject{}
	//ruo.Owner = account
	//rlm.db.Find("byOwner", &ruo)
	//return int64(ruo.RamUsage)
	return 0
}

func (rlm *ResourceLimitsManager) SetAccountLimits(account common.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) bool { //for test
	//pendingRlo := entity.ResourceLimitsObject{}
	//pendingRlo.Owner = account
	//pendingRlo.Pending = true
	//_, err := rlm.db.Find("byOwner", &pendingRlo)
	//if err != nil {
	//	rlo := entity.ResourceLimitsObject{}
	//	rlo.Owner = account
	//	rlo.Pending = false
	//	rlm.db.Find("byOwner", &rlo)
	//	pendingRlo.ID = rlo.ID
	//	pendingRlo.Owner = rlo.Owner
	//	pendingRlo.Pending = true
	//	pendingRlo.CpuWeight = rlo.CpuWeight
	//	pendingRlo.NetWeight = rlo.NetWeight
	//	pendingRlo.RamBytes = rlo.RamBytes
	//	rlm.db.Insert(&pendingRlo)
	//}
	//decreasedLimit := false
	//if ramBytes >= 0 {
	//	decreasedLimit = pendingRlo.RamBytes < 0 || ramBytes < pendingRlo.RamBytes
	//}
	//
	//rlm.db.Modify(&pendingRlo, func(data interface{}) error {
	//	ref := reflect.ValueOf(data).Elem()
	//	if ref.CanSet() {
	//		ref.FieldByName("RamBytes").SetInt(ramBytes)
	//		ref.FieldByName("NetWeight").SetInt(netWeight)
	//		ref.FieldByName("CpuWeight").SetInt(cpuWeight)
	//	}
	//	return nil
	//})
	//return decreasedLimit
	return false
}

func (rlm *ResourceLimitsManager) GetAccountLimits(account common.AccountName, ramBytes *int64, netWeight *int64, cpuWeight *int64) {
	//pendingRlo := entity.ResourceLimitsObject{}
	//pendingRlo.Owner = account
	//pendingRlo.Pending = true
	//_, err := rlm.db.Find("byOwner", &pendingRlo)
	//if err == nil {
	//	*ramBytes = pendingRlo.RamBytes
	//	*netWeight = pendingRlo.NetWeight
	//	*cpuWeight = pendingRlo.CpuWeight
	//} else {
	//	rlo := entity.ResourceLimitsObject{}
	//	rlo.Owner = account
	//	rlo.Pending = false
	//	rlm.db.Find("byOwner", &rlo)
	//	*ramBytes = rlo.RamBytes
	//	*netWeight = rlo.NetWeight
	//	*cpuWeight = rlo.CpuWeight
	//}
}

func (rlm *ResourceLimitsManager) ProcessAccountLimitUpdates() {
	//updateStateAndValue := func(total *uint64, value *int64, pendingValue int64, debugWhich string) {
	//	if *value > 0 {
	//		if *total < uint64(*value) {
	//			fmt.Println("error")
	//		}
	//		*total -= uint64(*value)
	//	}
	//
	//	if pendingValue > 0 {
	//		if math.MaxUint64-*total < uint64(pendingValue) {
	//			fmt.Println("error")
	//		}
	//		*total += uint64(pendingValue)
	//	}
	//
	//	*value = pendingValue
	//}
	//var pendingRlo []entity.ResourceLimitsObject
	//rlm.db.Get("Pending", true, &pendingRlo)
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("ID", ResourceLimitsState, &state)
	//rlm.db.Update(&state, func(data interface{}) error {
	//	for _, itr := range pendingRlo {
	//		rlo := ResourceLimitsObject{}
	//		rlm.db.Find("Rlo", RloIndex{ResourceLimits, itr.Owner, false}, &rlo)
	//		rlm.db.Update(&rlo, func(data interface{}) error {
	//			updateStateAndValue(&state.TotalRamBytes, &rlo.RamBytes, itr.RamBytes, "ram_bytes")
	//			updateStateAndValue(&state.TotalCpuWeight, &rlo.CpuWeight, itr.CpuWeight, "cpu_weight")
	//			updateStateAndValue(&state.TotalNetWeight, &rlo.NetWeight, itr.NetWeight, "net_weight")
	//			return nil
	//		})
	//	}
	//	return nil
	//})
}

func (rlm *ResourceLimitsManager) ProcessBlockUsage(blockNum uint32) {
	//config := entity.ResourceLimitsConfigObject{}
	//rlm.db.Find("byId", &config)
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//rlm.db.Modify(&state, func(data interface{}) error {
	//
	//	state.AverageBlockCpuUsage.Add(state.PendingCpuUsage, blockNum, config.CpuLimitParameters.Periods)
	//	state.UpdateVirtualCpuLimit(config)
	//	state.PendingCpuUsage = 0
	//
	//	state.AverageBlockNetUsage.Add(state.PendingNetUsage, blockNum, config.NetLimitParameters.Periods)
	//	state.UpdateVirtualNetLimit(config)
	//	state.PendingNetUsage = 0
	//
	//	return nil
	//})
}

func (rlm *ResourceLimitsManager) GetVirtualBlockCpuLimit() uint64 {
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//return state.VirtualCpuLimit
	return 0
}

func (rlm *ResourceLimitsManager) GetVirtualBlockNetLimit() uint64 {
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//return state.VirtualNetLimit
	return 0
}

func (rlm *ResourceLimitsManager) GetBlockCpuLimit() uint64 {
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//config := entity.ResourceLimitsConfigObject{}
	//rlm.db.Find("byId", &config)
	//return config.CpuLimitParameters.Max - state.PendingCpuUsage
	return 0
}

func (rlm *ResourceLimitsManager) GetBlockNetLimit() uint64 {
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//config := entity.ResourceLimitsConfigObject{}
	//rlm.db.Find("byId", &config)
	//return config.NetLimitParameters.Max - state.PendingNetUsage
	return 0
}

func (rlm *ResourceLimitsManager) GetAccountCpuLimit(name common.AccountName, elastic bool) int64 {
	arl := rlm.GetAccountCpuLimitEx(name, elastic)
	return arl.Available
}

func (rlm *ResourceLimitsManager) GetAccountCpuLimitEx(name common.AccountName, elastic bool) AccountResourceLimit {
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//config := entity.ResourceLimitsConfigObject{}
	//rlm.db.Find("byId", &config)
	//ruo := entity.ResourceUsageObject{}
	//rlm.db.Find("byOwner", &ruo)
	//
	//var cpuWeight, x, y int64
	//rlm.GetAccountLimits(name, &x, &y, &cpuWeight)
	//
	//if cpuWeight < 0 || state.TotalCpuWeight == 0 {
	//	return AccountResourceLimit{-1, -1, -1}
	//}

	arl := AccountResourceLimit{}
	//windowSize := new(big.Int).SetUint64(uint64(config.AccountCpuUsageAverageWindow))
	//virtualCpuCapacityInWindow := new(big.Int)
	//if elastic {
	//	virtualCpuCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(state.VirtualCpuLimit), windowSize)
	//} else {
	//	virtualCpuCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(config.CpuLimitParameters.Max), windowSize)
	//}
	//userWeight := new(big.Int).SetUint64(uint64(cpuWeight))
	//allUserWeight := new(big.Int).SetUint64(state.TotalCpuWeight)
	//
	//maxUserUseInWindow := new(big.Int).Div(new(big.Int).Mul(virtualCpuCapacityInWindow, userWeight), allUserWeight)
	//cpuUsedInWindow := IntegerDivideCeil(
	//	new(big.Int).Mul(new(big.Int).SetUint64(ruo.CpuUsage.ValueEx), windowSize),
	//	new(big.Int).SetUint64(uint64(common.DefaultConfig.RateLimitingPrecision)))
	//
	//if maxUserUseInWindow.Cmp(cpuUsedInWindow) != 1 {
	//	arl.Available = 0
	//} else {
	//	arl.Available = DowngradeCast(new(big.Int).Sub(maxUserUseInWindow, cpuUsedInWindow))
	//}
	//
	//arl.Used = DowngradeCast(cpuUsedInWindow)
	//arl.Max = DowngradeCast(maxUserUseInWindow)
	return arl
}

func (rlm *ResourceLimitsManager) GetAccountNetLimit(name common.AccountName, elastic bool) int64 {
	arl := rlm.GetAccountNetLimitEx(name, elastic)
	return arl.Available
}

func (rlm *ResourceLimitsManager) GetAccountNetLimitEx(name common.AccountName, elastic bool) AccountResourceLimit {
	//state := entity.ResourceLimitsStateObject{}
	//rlm.db.Find("byId", &state)
	//config := entity.ResourceLimitsConfigObject{}
	//rlm.db.Find("byId", &config)
	//ruo := entity.ResourceUsageObject{}
	//rlm.db.Find("byOwner", &ruo)
	//
	//var netWeight, x, y int64
	//rlm.GetAccountLimits(name, &x, &y, &netWeight)
	//
	//if netWeight < 0 || state.TotalNetWeight == 0 {
	//	return AccountResourceLimit{-1, -1, -1}
	//}

	arl := AccountResourceLimit{}
	//windowSize := new(big.Int).SetUint64(uint64(config.AccountNetUsageAverageWindow))
	//virtualNetCapacityInWindow := new(big.Int)
	//if elastic {
	//	virtualNetCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(state.VirtualNetLimit), windowSize)
	//} else {
	//	virtualNetCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(config.NetLimitParameters.Max), windowSize)
	//}
	//userWeight := new(big.Int).SetUint64(uint64(netWeight))
	//allUserWeight := new(big.Int).SetUint64(state.TotalNetWeight)
	//
	//maxUserUseInWindow := new(big.Int).Div(new(big.Int).Mul(virtualNetCapacityInWindow, userWeight), allUserWeight)
	//netUsedInWindow := IntegerDivideCeil(
	//	new(big.Int).Mul(new(big.Int).SetUint64(ruo.NetUsage.ValueEx), windowSize),
	//	new(big.Int).SetUint64(uint64(common.DefaultConfig.RateLimitingPrecision)))
	//if maxUserUseInWindow.Cmp(netUsedInWindow) != 1 {
	//	arl.Available = 0
	//} else {
	//	arl.Available = DowngradeCast(new(big.Int).Sub(maxUserUseInWindow, netUsedInWindow))
	//}
	//
	//arl.Used = DowngradeCast(netUsedInWindow)
	//arl.Max = DowngradeCast(maxUserUseInWindow)
	return arl
}
