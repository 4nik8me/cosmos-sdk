package params

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/proposal"
)

func NewHandler(k ProposalKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, k ProposalKeeper, msg MsgSubmitProposal) sdk.Result {
	content := NewProposalChange(msg.Title, msg.Description, msg.Changes)
	return proposal.HandleSubmit(ctx, k.proposal, content, msg.Proposer, msg.InitialDeposit)
}

func NewProposalHandler(k ProposalKeeper) proposal.Handler {
	return func(ctx sdk.Context, p proposal.Content) sdk.Error {
		switch p := p.(type) {
		case ProposalChange:
			return handleProposalChange(ctx, k, p)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov proposal type: %T", p)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

func handleProposalChange(ctx sdk.Context, k ProposalKeeper, p ProposalChange) sdk.Error {
	for _, c := range p.Changes {
		s, ok := k.GetSubspace(c.Space)
		if !ok {
			return ErrUnknownSubspace(k.codespace, c.Space)
		}
		var err error
		if len(c.Subkey) == 0 {
			err = s.SetRaw(ctx, c.Key, c.Value)
		} else {
			err = s.SetRawWithSubkey(ctx, c.Key, c.Subkey, c.Value)
		}

		if err != nil {
			return ErrSettingParameter(k.codespace, c.Key, c.Subkey, c.Value, err.Error())
		}
	}

	return nil
}
