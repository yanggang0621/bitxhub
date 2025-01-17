package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/looplab/fsm"
	"github.com/meshplus/bitxhub-core/boltvm"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/bitxhub/internal/repo"
)

type GovStrategy struct {
	boltvm.Stub
}

// Proposal strategy ===============================================================
type ProposalStrategyType string

const (
	SuperMajorityApprove ProposalStrategyType = repo.SuperMajorityApprove
	SuperMajorityAgainst ProposalStrategyType = repo.SuperMajorityAgainst
	SimpleMajority       ProposalStrategyType = repo.SimpleMajority
	ZeroPermission       ProposalStrategyType = repo.ZeroPermission
)

type ProposalStrategy struct {
	Module string               `json:"module"`
	Typ    ProposalStrategyType `json:"typ"`
	// The minimum participation threshold.
	// Only when the number of voting participants reaches this proportion,
	// the proposal will take effect. That is, the proposal can be judged
	// according to the voting situation.
	ParticipateThreshold float64                     `json:"participate_threshold"`
	Extra                []byte                      `json:"extra"`
	Status               governance.GovernanceStatus `json:"status"`
	FSM                  *fsm.FSM                    `json:"fsm"`
}

type UpdateStrategyInfo struct {
	Typ                  UpdateInfo `json:"typ"`
	ParticipateThreshold UpdateInfo `json:"participate_threshold"`
}

func (d *ProposalStrategy) setFSM(lastStatus governance.GovernanceStatus) {
	d.FSM = fsm.NewFSM(
		string(d.Status),
		fsm.Events{
			// register 3
			{Name: string(governance.EventRegister), Src: []string{string(governance.GovernanceUnavailable)}, Dst: string(governance.GovernanceRegisting)},
			{Name: string(governance.EventApprove), Src: []string{string(governance.GovernanceRegisting)}, Dst: string(governance.GovernanceAvailable)},
			{Name: string(governance.EventReject), Src: []string{string(governance.GovernanceRegisting)}, Dst: string(lastStatus)},

			// update 1
			{Name: string(governance.EventUpdate), Src: []string{string(governance.GovernanceAvailable)}, Dst: string(governance.GovernanceUpdating)},
			{Name: string(governance.EventApprove), Src: []string{string(governance.GovernanceUpdating)}, Dst: string(governance.GovernanceAvailable)},
			{Name: string(governance.EventReject), Src: []string{string(governance.GovernanceUpdating)}, Dst: string(governance.GovernanceAvailable)},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.Status = governance.GovernanceStatus(d.FSM.Current()) },
		},
	)
}

func (g *GovStrategy) changeStatus(pt string, trigger, lastStatus string) (bool, []byte) {
	ps := &ProposalStrategy{}
	if ok := g.GetObject(ProposalStrategyKey(pt), ps); !ok {
		return false, []byte("this proposal strategy does not exist")
	}

	ps.setFSM(governance.GovernanceStatus(lastStatus))
	err := ps.FSM.Event(trigger)
	if err != nil {
		return false, []byte(fmt.Sprintf("change status error: %v", err))
	}

	g.SetObject(ProposalStrategyKey(pt), *ps)
	return true, nil
}

func (g *GovStrategy) checkPermission(permissions []string, regulatedAddr, regulatorAddr string, specificAddrsData []byte) error {
	for _, permission := range permissions {
		switch permission {
		case string(PermissionSelf):
			if regulatedAddr == regulatorAddr {
				return nil
			}
		case string(PermissionAdmin):
			res := g.CrossInvoke(constant.RoleContractAddr.Address().String(), "IsAnyAvailableAdmin",
				pb.String(regulatorAddr),
				pb.String(string(GovernanceAdmin)))
			if !res.Ok {
				return fmt.Errorf("cross invoke IsAvailableGovernanceAdmin error:%s", string(res.Result))
			}
			if "true" == string(res.Result) {
				return nil
			}
		case string(PermissionSpecific):
			specificAddrs := []string{}
			if err := json.Unmarshal(specificAddrsData, &specificAddrs); err != nil {
				return err
			}
			for _, addr := range specificAddrs {
				if addr == regulatorAddr {
					return nil
				}
			}
		default:
			return fmt.Errorf("unsupport permission: %s", permission)
		}
	}

	return fmt.Errorf("regulatorAddr(%s) does not have the permission", regulatorAddr)
}

var mgrs = []string{repo.AppchainMgr, repo.NodeMgr, repo.DappMgr, repo.RoleMgr, repo.RuleMgr, repo.ServiceMgr, repo.ProposalStrategyMgr}

// =========== Manage does some subsequent operations when the proposal is over
func (g *GovStrategy) Manage(eventTyp, proposalResult, lastStatus, objId string, extra []byte) *boltvm.Response {
	// 1. check permission: PermissionSpecific(GovernanceContractAddr)
	specificAddrs := []string{constant.GovernanceContractAddr.Address().String()}
	addrsData, err := json.Marshal(specificAddrs)
	if err != nil {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), fmt.Sprintf("marshal specificAddrs error: %v", err)))
	}
	if err := g.checkPermission([]string{string(PermissionSpecific)}, objId, g.CurrentCaller(), addrsData); err != nil {
		return boltvm.Error(boltvm.ProposalStrategyNoPermissionCode, fmt.Sprintf(string(boltvm.ProposalStrategyNoPermissionMsg), g.CurrentCaller(), err.Error()))
	}

	// 2. change status
	if objId != repo.AllMgr {
		ok, errData := g.changeStatus(objId, proposalResult, lastStatus)
		if !ok {
			return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), fmt.Sprintf("change %s status error: %s", objId, string(errData))))
		}
	} else {
		for _, mgr := range mgrs {
			ok, errData := g.changeStatus(mgr, proposalResult, lastStatus)
			if !ok {
				return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), fmt.Sprintf("change %s status error: %s", objId, string(errData))))
			}
		}
	}

	if proposalResult == string(APPROVED) {
		if eventTyp == string(governance.EventUpdate) {
			updateStrategyInfo := &UpdateStrategyInfo{}
			if err := json.Unmarshal(extra, updateStrategyInfo); err != nil {
				return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), fmt.Sprintf("unmarshal update strategy")))
			}
			if objId == repo.AllMgr {
				for _, mgr := range mgrs {
					ps := &ProposalStrategy{}
					if ok := g.GetObject(ProposalStrategyKey(mgr), ps); !ok {
						return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyIllegalProposalStrategyInfoMsg), fmt.Sprintf("get proposal strategy %s error", objId)))
					}
					if updateStrategyInfo.Typ.IsEdit {
						ps.Typ = ProposalStrategyType(updateStrategyInfo.Typ.NewInfo.(string))
					}
					if updateStrategyInfo.ParticipateThreshold.IsEdit {
						ps.ParticipateThreshold = updateStrategyInfo.ParticipateThreshold.NewInfo.(float64)
					}
					g.SetObject(ProposalStrategyKey(mgr), ps)
				}
			} else {
				ps := &ProposalStrategy{}
				if ok := g.GetObject(ProposalStrategyKey(objId), ps); !ok {
					return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyIllegalProposalStrategyInfoMsg), fmt.Sprintf("get proposal strategy %s error", objId)))
				}
				if updateStrategyInfo.Typ.IsEdit {
					ps.Typ = ProposalStrategyType(updateStrategyInfo.Typ.NewInfo.(string))
				}
				if updateStrategyInfo.ParticipateThreshold.IsEdit {
					ps.ParticipateThreshold = updateStrategyInfo.ParticipateThreshold.NewInfo.(float64)
				}
				g.SetObject(ProposalStrategyKey(ps.Module), ps)
			}
		}
	}
	return boltvm.Success(nil)
}

// update proposal strategy for a proposal type
func (g *GovStrategy) UpdateProposalStrategy(pt string, typ string, participateThreshold float64, reason string) *boltvm.Response {
	if err := g.checkPermission([]string{string(PermissionAdmin)}, pt, g.CurrentCaller(), nil); err != nil {
		return boltvm.Error(boltvm.ProposalStrategyNoPermissionCode, fmt.Sprintf(string(boltvm.ProposalStrategyNoPermissionMsg), g.CurrentCaller(), fmt.Sprintf("check permission error:%v", err)))
	}
	ps := &ProposalStrategy{}
	if ok := g.GetObject(ProposalStrategyKey(pt), ps); !ok {
		return boltvm.Error(boltvm.ProposalStrategyNonexistentProposalStrategyCode, fmt.Sprintf(string(boltvm.ProposalStrategyNonexistentProposalStrategyMsg), pt))
	}
	if ps.Status != governance.GovernanceAvailable {
		return boltvm.Error(boltvm.ProposalStrategyStatusErrorCode, fmt.Sprintf(string(boltvm.ProposalStrategyStatusErrorMsg), typ, ps.Status, governance.EventUpdate))
	}
	info := &UpdateStrategyInfo{}
	info.Typ = UpdateInfo{
		OldInfo: ps.Typ,
		NewInfo: typ,
		IsEdit:  string(ps.Typ) != typ,
	}
	info.ParticipateThreshold = UpdateInfo{
		OldInfo: ps.ParticipateThreshold,
		NewInfo: participateThreshold,
		IsEdit:  ps.ParticipateThreshold != participateThreshold,
	}
	extra, err := json.Marshal(info)
	if err != nil {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf("unmarshal update strategy error: %v", err))
	}
	if err := repo.CheckStrategyInfo(string(ps.Typ), string(ps.Module), ps.ParticipateThreshold); err != nil {
		return boltvm.Error(boltvm.ProposalStrategyIllegalProposalStrategyInfoCode, fmt.Sprintf(string(boltvm.ProposalStrategyIllegalProposalStrategyInfoMsg), err.Error()))
	}

	res := g.CrossInvoke(constant.GovernanceContractAddr.Address().String(), "SubmitProposal",
		pb.String(g.Caller()),
		pb.String(string(governance.EventUpdate)),
		pb.String(string(ProposalStrategyMgr)),
		pb.String(pt),
		pb.String(string(ps.Status)),
		pb.String(reason),
		pb.Bytes(extra),
	)
	if !res.Ok {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), fmt.Sprintf("submit proposal error: %s", string(res.Result))))
	}

	g.changeStatus(pt, string(governance.EventUpdate), string(ps.Status))

	g.CrossInvoke(constant.GovernanceContractAddr.Address().String(), "ZeroPermission", pb.String(string(res.Result)))

	return getGovernanceRet(string(res.Result), []byte(pt))
}

// update proposal strategy for a proposal type
func (g *GovStrategy) UpdateAllProposalStrategy(typ string, participateThreshold float64, reason string) *boltvm.Response {
	pt := repo.AllMgr
	if err := g.checkPermission([]string{string(PermissionAdmin)}, pt, g.CurrentCaller(), nil); err != nil {
		return boltvm.Error(boltvm.ProposalStrategyNoPermissionCode, fmt.Sprintf(string(boltvm.ProposalStrategyNoPermissionMsg), g.CurrentCaller(), fmt.Sprintf("check permission error:%v", err)))
	}
	info := &UpdateStrategyInfo{}
	info.Typ = UpdateInfo{
		NewInfo: typ,
		IsEdit:  true,
	}
	info.ParticipateThreshold = UpdateInfo{
		NewInfo: participateThreshold,
		IsEdit:  true,
	}
	extra, err := json.Marshal(info)
	if err != nil {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf("unmarshal update strategy error: %v", err))
	}
	if err := repo.CheckStrategyInfo(typ, pt, participateThreshold); err != nil {
		return boltvm.Error(boltvm.ProposalStrategyIllegalProposalStrategyInfoCode, fmt.Sprintf(string(boltvm.ProposalStrategyIllegalProposalStrategyInfoMsg), err.Error()))
	}

	res := g.CrossInvoke(constant.GovernanceContractAddr.Address().String(), "SubmitProposal",
		pb.String(g.Caller()),
		pb.String(string(governance.EventUpdate)),
		pb.String(string(ProposalStrategyMgr)),
		pb.String(pt),
		pb.String(string(governance.GovernanceAvailable)),
		pb.String(reason),
		pb.Bytes(extra),
	)
	if !res.Ok {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), fmt.Sprintf("submit proposal error: %s", string(res.Result))))
	}

	for _, mgr := range mgrs {
		g.changeStatus(mgr, string(governance.EventUpdate), string(governance.GovernanceAvailable))
	}

	g.CrossInvoke(constant.GovernanceContractAddr.Address().String(), "ZeroPermission", pb.String(string(res.Result)))

	return getGovernanceRet(string(res.Result), []byte(pt))
}

func (g *GovStrategy) GetAllProposalStrategy() *boltvm.Response {
	ret := make([]*ProposalStrategy, 0)
	ok, value := g.Query(PROPOSALSTRATEGY_PREFIX)
	if ok {
		for _, data := range value {
			ps := &ProposalStrategy{}
			if err := json.Unmarshal(data, ps); err != nil {
				return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf("unmarshal proposal strategy error: %v", err))
			}
			ret = append(ret, ps)
		}
	}
	data, err := json.Marshal(ret)
	if err != nil {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf("marshal proposal strategy error: %v", err))
	}
	return boltvm.Success(data)
}

func (g *GovStrategy) GetProposalStrategy(pt string) *boltvm.Response {
	if err := repo.CheckManageModule(pt); err != nil {
		return boltvm.Error(boltvm.ProposalStrategyIllegalProposalTypeCode, fmt.Sprintf(string(boltvm.ProposalStrategyIllegalProposalTypeMsg), pt))
	}

	ps := &ProposalStrategy{}
	if !g.GetObject(ProposalStrategyKey(pt), ps) {
		return boltvm.Error(boltvm.ProposalStrategyNonexistentProposalStrategyCode, fmt.Sprintf(string(boltvm.ProposalStrategyNonexistentProposalStrategyMsg), pt))
	}

	pData, err := json.Marshal(ps)
	if err != nil {
		return boltvm.Error(boltvm.ProposalStrategyInternalErrCode, fmt.Sprintf(string(boltvm.ProposalStrategyInternalErrMsg), err.Error()))
	}
	return boltvm.Success(pData)
}

func ProposalStrategyKey(ps string) string {
	return fmt.Sprintf("%s-%s", PROPOSALSTRATEGY_PREFIX, ps)
}
