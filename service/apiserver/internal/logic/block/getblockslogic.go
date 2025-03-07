package block

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/bnb-chain/zkbas/service/apiserver/internal/logic/utils"
	"github.com/bnb-chain/zkbas/service/apiserver/internal/svc"
	"github.com/bnb-chain/zkbas/service/apiserver/internal/types"
	types2 "github.com/bnb-chain/zkbas/types"
)

type GetBlocksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlocksLogic {
	return &GetBlocksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBlocksLogic) GetBlocks(req *types.ReqGetRange) (*types.Blocks, error) {
	total, err := l.svcCtx.MemCache.GetBlockTotalCountWithFallback(func() (interface{}, error) {
		return l.svcCtx.BlockModel.GetCurrentHeight()
	})
	if err != nil {
		return nil, types2.AppErrInternal
	}

	resp := &types.Blocks{
		Blocks: make([]*types.Block, 0),
		Total:  uint32(total),
	}
	if total == 0 || total <= int64(req.Offset) {
		return resp, nil
	}

	blocks, err := l.svcCtx.BlockModel.GetBlocksList(int64(req.Limit), int64(req.Offset))
	if err != nil {
		return nil, types2.AppErrInternal
	}
	for _, b := range blocks {
		block := &types.Block{
			Commitment:                      b.BlockCommitment,
			Height:                          b.BlockHeight,
			StateRoot:                       b.StateRoot,
			PriorityOperations:              b.PriorityOperations,
			PendingOnChainOperationsHash:    b.PendingOnChainOperationsHash,
			PendingOnChainOperationsPubData: b.PendingOnChainOperationsPubData,
			CommittedTxHash:                 b.CommittedTxHash,
			CommittedAt:                     b.CommittedAt,
			VerifiedTxHash:                  b.VerifiedTxHash,
			VerifiedAt:                      b.VerifiedAt,
			Status:                          b.BlockStatus,
		}
		for _, t := range b.Txs {
			tx := utils.DbtxTx(t)
			tx.AccountName, _ = l.svcCtx.MemCache.GetAccountNameByIndex(tx.AccountIndex)
			block.Txs = append(block.Txs, tx)
		}
		resp.Blocks = append(resp.Blocks, block)
	}
	return resp, nil
}
