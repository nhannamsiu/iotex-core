// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package blockchain

import (
	"bytes"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/blockchain/action"
	"github.com/iotexproject/iotex-core/crypto"
	"github.com/iotexproject/iotex-core/iotxaddress"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/state"
)

// Validator is the interface of validator
type Validator interface {
	// Validate validates the given block's content
	Validate(block *Block, tipHeight uint64, tipHash hash.Hash32B, containCoinbase bool) error
}

type validator struct {
	sf            state.Factory
	validatorAddr string
}

var (
	// ErrInvalidTipHeight is the error returned when the block height is not valid
	ErrInvalidTipHeight = errors.New("invalid tip height")
	// ErrInvalidBlock is the error returned when the block is not valid
	ErrInvalidBlock = errors.New("failed to validate the block")
	// ErrActionNonce is the error when the nonce of the action is wrong
	ErrActionNonce = errors.New("invalid action nonce")
	// ErrGasHigherThanLimit indicates the error of gas value
	ErrGasHigherThanLimit = errors.New("invalid gas for action")
	// ErrInsufficientGas indicates the error of insufficient gas value for data storage
	ErrInsufficientGas = errors.New("insufficient intrinsic gas value")
	// ErrBalance indicates the error of balance
	ErrBalance = errors.New("invalid balance")
	// ErrDKGSecretProposal indicates the error of DKG secret proposal
	ErrDKGSecretProposal = errors.New("invalid DKG secret proposal")
)

// Validate validates the given block's content
func (v *validator) Validate(blk *Block, tipHeight uint64, tipHash hash.Hash32B, containCoinbase bool) error {
	if err := verifyHeightAndHash(blk, tipHeight, tipHash); err != nil {
		return errors.Wrap(err, "failed to verify block's height and hash")
	}
	if blk.IsDummyBlock() {
		return nil
	}
	if err := verifySigAndRoot(blk); err != nil {
		return errors.Wrap(err, "failed to verify block's signature and merkle root")
	}

	if v.sf != nil {
		return v.verifyActions(blk, containCoinbase)
	}

	return nil
}

func (v *validator) verifyActions(blk *Block, containCoinbase bool) error {
	// Verify transfers, votes, executions, witness, and secrets (balance is checked in RunActions)
	confirmedNonceMap := make(map[string]uint64)
	accountNonceMap := make(map[string][]uint64)
	var wg sync.WaitGroup
	wg.Add(len(blk.Transfers) + len(blk.Votes) + len(blk.Executions))
	var correctAction uint64
	var coinbaseCount uint64
	for _, tsf := range blk.Transfers {
		// Verify Address
		// Verify Gas
		// Verify Nonce
		// Verify Signature
		// Verify Coinbase transfer

		if !tsf.IsCoinbase() {
			if _, err := iotxaddress.GetPubkeyHash(tsf.Sender()); err != nil {
				return errors.Wrapf(err, "failed to validate transfer sender's address %s", tsf.Sender())
			}
			if _, err := iotxaddress.GetPubkeyHash(tsf.Recipient()); err != nil {
				return errors.Wrapf(err, "failed to validate transfer recipient's address %s", tsf.Recipient())
			}
		}

		if blk.Header.height > 0 && !tsf.IsCoinbase() {
			// Reject over-gassed transfer
			if tsf.GasLimit() > action.GasLimit {
				return errors.Wrapf(ErrGasHigherThanLimit, "gas is higher than gas limit")
			}
			intrinsicGas, err := tsf.IntrinsicGas()
			if intrinsicGas > tsf.GasLimit() || err != nil {
				return errors.Wrapf(ErrInsufficientGas, "insufficient gas for transfer")
			}
			// Store the nonce of the sender and verify later
			if _, ok := confirmedNonceMap[tsf.Sender()]; !ok {
				accountNonce, err := v.sf.Nonce(tsf.Sender())
				if err != nil {
					return errors.Wrap(err, "failed to get the nonce of transfer sender")
				}
				confirmedNonceMap[tsf.Sender()] = accountNonce
				accountNonceMap[tsf.Sender()] = make([]uint64, 0)
			}
			accountNonceMap[tsf.Sender()] = append(accountNonceMap[tsf.Sender()], tsf.Nonce())
		}

		go func(tsf *action.Transfer, correctTsf *uint64, correctCoinbase *uint64) {
			defer wg.Done()
			// Verify coinbase transfer
			if tsf.IsCoinbase() {
				pkHash := keypair.HashPubKey(blk.Header.Pubkey)
				addr := address.New(blk.Header.chainID, pkHash[:])
				if addr.IotxAddress() != tsf.Recipient() {
					return
				}
				atomic.AddUint64(correctCoinbase, uint64(1))
				return
			}
			if err := action.Verify(tsf); err != nil {
				return
			}
			atomic.AddUint64(correctTsf, uint64(1))
		}(tsf, &correctAction, &coinbaseCount)
	}
	for _, vote := range blk.Votes {
		// Verify Address
		// Verify Gas
		// Verify Nonce
		// Verify Signature

		if _, err := iotxaddress.GetPubkeyHash(vote.Voter()); err != nil {
			return errors.Wrapf(err, "failed to validate voter's address %s", vote.Voter())
		}
		if vote.Votee() != action.EmptyAddress {
			if _, err := iotxaddress.GetPubkeyHash(vote.Votee()); err != nil {
				return errors.Wrapf(err, "failed to validate votee's address %s", vote.Votee())
			}
		}

		if blk.Header.height > 0 {
			// Reject over-gassed vote
			if vote.GasLimit() > action.GasLimit {
				return errors.Wrapf(ErrGasHigherThanLimit, "gas is higher than gas limit")
			}
			intrinsicGas, err := vote.IntrinsicGas()
			if intrinsicGas > vote.GasLimit() || err != nil {
				return errors.Wrapf(ErrInsufficientGas, "insufficient gas for vote")
			}
			// Store the nonce of the voter and verify later
			voterAddress := vote.Voter()
			if _, ok := confirmedNonceMap[voterAddress]; !ok {
				accountNonce, err := v.sf.Nonce(voterAddress)
				if err != nil {
					return errors.Wrap(err, "failed to get the nonce of the voter")
				}
				confirmedNonceMap[voterAddress] = accountNonce
				accountNonceMap[voterAddress] = make([]uint64, 0)
			}
			accountNonceMap[voterAddress] = append(accountNonceMap[voterAddress], vote.Nonce())
		}

		// Verify signature
		go func(vote *action.Vote, correctVote *uint64) {
			defer wg.Done()
			if err := action.Verify(vote); err != nil {
				return
			}
			atomic.AddUint64(correctVote, uint64(1))
		}(vote, &correctAction)
	}
	for _, execution := range blk.Executions {
		// Verify Address
		// Verify Nonce
		// Verify Signature
		// Verify Gas
		// Verify Amount

		if _, err := iotxaddress.GetPubkeyHash(execution.Executor()); err != nil {
			return errors.Wrapf(err, "failed to validate executor's address %s", execution.Executor())
		}
		if execution.Contract() != action.EmptyAddress {
			if _, err := iotxaddress.GetPubkeyHash(execution.Contract()); err != nil {
				return errors.Wrapf(err, "failed to validate contract's address %s", execution.Contract())
			}
		}

		if blk.Header.height > 0 {
			// Store the nonce of the executor and verify later
			executor := execution.Executor()
			if _, ok := confirmedNonceMap[executor]; !ok {
				accountNonce, err := v.sf.Nonce(executor)
				if err != nil {
					return errors.Wrap(err, "failed to get the nonce of the executor")
				}
				confirmedNonceMap[executor] = accountNonce
				accountNonceMap[executor] = make([]uint64, 0)
			}
			accountNonceMap[executor] = append(accountNonceMap[executor], execution.Nonce())
		}

		// Verify signature
		go func(execution *action.Execution, correctVote *uint64) {
			defer wg.Done()
			if err := action.Verify(execution); err != nil {
				return
			}
			atomic.AddUint64(correctVote, uint64(1))
		}(execution, &correctAction)

		// Reject over-gassed execution
		if execution.GasLimit() > action.GasLimit {
			return errors.Wrapf(ErrGasHigherThanLimit, "gas is higher than gas limit")
		}
		intrinsicGas, err := execution.IntrinsicGas()
		if intrinsicGas > execution.GasLimit() || err != nil {
			return errors.Wrapf(ErrInsufficientGas, "insufficient gas for execution")
		}

		// Reject execution of negative amount
		if execution.Amount().Sign() < 0 {
			return errors.Wrapf(ErrBalance, "negative value")
		}
	}
	wg.Wait()
	// Verify coinbase transfer count
	if (containCoinbase && coinbaseCount != 1) || (!containCoinbase && coinbaseCount != 0) {
		return errors.Wrapf(
			ErrInvalidBlock,
			"wrong number of coinbase transfers")
	}
	if correctAction+coinbaseCount != uint64(len(blk.Transfers)+len(blk.Votes)+len(blk.Executions)) {
		return errors.Wrapf(
			ErrInvalidBlock,
			"failed to verify actions signature")
	}

	// Verify Witness
	if blk.SecretWitness != nil {
		// Verify witness sender address
		if _, err := iotxaddress.GetPubkeyHash(blk.SecretWitness.SrcAddr()); err != nil {
			return errors.Wrapf(err, "failed to validate witness sender's address %s", blk.SecretWitness.SrcAddr())
		}
		// Store the nonce of the witness sender and verify later
		if _, ok := confirmedNonceMap[blk.SecretWitness.SrcAddr()]; !ok {
			accountNonce, err := v.sf.Nonce(blk.SecretWitness.SrcAddr())
			if err != nil {
				return errors.Wrap(err, "failed to get the nonce of secret sender")
			}
			confirmedNonceMap[blk.SecretWitness.SrcAddr()] = accountNonce
			accountNonceMap[blk.SecretWitness.SrcAddr()] = make([]uint64, 0)
		}
		accountNonceMap[blk.SecretWitness.SrcAddr()] = append(accountNonceMap[blk.SecretWitness.SrcAddr()], blk.SecretWitness.Nonce())
	}

	// Verify Secrets
	for _, sp := range blk.SecretProposals {
		// Verify address
		if _, err := iotxaddress.GetPubkeyHash(sp.SrcAddr()); err != nil {
			return errors.Wrapf(err, "failed to validate secret sender's address %s", sp.SrcAddr())
		}
		if _, err := iotxaddress.GetPubkeyHash(sp.DstAddr()); err != nil {
			return errors.Wrapf(err, "failed to validate secret recipient's address %s", sp.DstAddr())
		}

		// Store the nonce of the sender and verify later
		if _, ok := confirmedNonceMap[sp.SrcAddr()]; !ok {
			accountNonce, err := v.sf.Nonce(sp.SrcAddr())
			if err != nil {
				return errors.Wrap(err, "failed to get the nonce of secret sender")
			}
			confirmedNonceMap[sp.SrcAddr()] = accountNonce
			accountNonceMap[sp.SrcAddr()] = make([]uint64, 0)
		}
		accountNonceMap[sp.SrcAddr()] = append(accountNonceMap[sp.SrcAddr()], sp.Nonce())

		// verify secret if the validator is recipient
		if v.validatorAddr == sp.DstAddr() {
			validatorID := iotxaddress.CreateID(v.validatorAddr)
			result, err := crypto.DKG.ShareVerify(validatorID, sp.Secret(), blk.SecretWitness.Witness())
			if err == nil {
				err = ErrDKGSecretProposal
			}
			if !result {
				return errors.Wrap(err, "failed to verify the DKG secret share")
			}
		}
	}

	if blk.Header.height > 0 {
		//Verify each account's Nonce
		for address := range confirmedNonceMap {
			// The nonce of each action should be increasing, unique and consecutive
			confirmedNonce := confirmedNonceMap[address]
			receivedNonce := accountNonceMap[address]
			sort.Slice(receivedNonce, func(i, j int) bool { return receivedNonce[i] < receivedNonce[j] })
			for i, nonce := range receivedNonce {
				if nonce != confirmedNonce+uint64(i+1) {
					return errors.Wrap(ErrActionNonce, "the nonce of the action is invalid")
				}
			}
		}
	}
	return nil
}

func verifyHeightAndHash(blk *Block, tipHeight uint64, tipHash hash.Hash32B) error {
	if blk == nil {
		return ErrInvalidBlock
	}
	// verify new block has height incremented by 1
	if blk.Header.height != 0 && blk.Header.height != tipHeight+1 {
		return errors.Wrapf(
			ErrInvalidTipHeight,
			"wrong block height %d, expecting %d",
			blk.Header.height,
			tipHeight+1)
	}
	// verify new block has correctly linked to current tip
	if blk.Header.prevBlockHash != tipHash {
		return errors.Wrapf(
			ErrInvalidBlock,
			"wrong prev hash %x, expecting %x",
			blk.Header.prevBlockHash,
			tipHash)
	}
	return nil
}

func verifySigAndRoot(blk *Block) error {
	if blk.Header.height > 0 {
		// verify new block's signature is correct
		blkHash := blk.HashBlock()
		if !crypto.EC283.Verify(blk.Header.Pubkey, blkHash[:], blk.Header.blockSig) {
			return errors.Wrapf(
				ErrInvalidBlock,
				"failed to verify block's signature with public key: %x",
				blk.Header.Pubkey)
		}
	}

	hashExpect := blk.Header.txRoot
	hashActual := blk.TxRoot()
	if !bytes.Equal(hashExpect[:], hashActual[:]) {
		return errors.Wrapf(
			ErrInvalidBlock,
			"wrong tx hash %x, expecting %x",
			hashActual,
			hashActual)
	}
	return nil
}
