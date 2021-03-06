// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"bytes"
	"math/big"

	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/crypto"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/proto"
)

var (
	// ErrAction indicates error for an action
	ErrAction = errors.New("action error")
	// ErrAddress indicates error of address
	ErrAddress = errors.New("address error")
)

// Action is the generic interface of all types of actions, and defines the common methods of them
type Action interface {
	Version() uint32
	Nonce() uint64
	SrcAddr() string
	SrcPubkey() keypair.PublicKey
	SetSrcPubkey(srcPubkey keypair.PublicKey)
	DstAddr() string
	GasLimit() uint64
	GasPrice() *big.Int
	Signature() []byte
	SetSignature(signature []byte)
	ByteStream() []byte
	Hash() hash.Hash32B
	IntrinsicGas() (uint64, error)
	Cost() (*big.Int, error)
	ConvertToActionPb() *iproto.ActionPb
}

type action struct {
	version   uint32
	nonce     uint64
	srcAddr   string
	srcPubkey keypair.PublicKey
	dstAddr   string
	gasLimit  uint64
	gasPrice  *big.Int
	signature []byte
}

// NewActionFromProto converts a proto message into a corresponding action struct
func NewActionFromProto(pbAct *iproto.ActionPb) Action {
	// TODO: implement the logic
	return nil
}

// Version returns the version
func (act *action) Version() uint32 { return act.version }

// Nonce returns the nonce
func (act *action) Nonce() uint64 { return act.nonce }

// SrcAddr returns the source address
func (act *action) SrcAddr() string { return act.srcAddr }

// SrcPubkey returns the source public key
func (act *action) SrcPubkey() keypair.PublicKey { return act.srcPubkey }

// SetSrcPubkey sets the source public key
func (act *action) SetSrcPubkey(srcPubkey keypair.PublicKey) { act.srcPubkey = srcPubkey }

// DstAddr returns the destination address
func (act *action) DstAddr() string { return act.dstAddr }

// GasLimit returns the gas limit
func (act *action) GasLimit() uint64 { return act.gasLimit }

// GasPrice returns the gas price
func (act *action) GasPrice() *big.Int { return act.gasPrice }

// Signature returns signature bytes
func (act *action) Signature() []byte { return act.signature }

// SetSignature sets the signature bytes
func (act *action) SetSignature(signature []byte) { act.signature = signature }

// Sign signs the action using sender's private key
func Sign(act Action, sk keypair.PrivateKey) error {
	// TODO: remove this conversion once we deprecate old address format
	srcAddr, err := address.IotxAddressToAddress(act.SrcAddr())
	if err != nil {
		return errors.Wrapf(err, "error when converting from old address format")
	}
	// TODO: we should avoid generate public key from private key in each signature
	pk, err := crypto.EC283.NewPubKey(sk)
	if err != nil {
		return errors.Wrapf(err, "error when deriving public key from private key")
	}
	pkHash := keypair.HashPubKey(pk)
	// TODO: abstract action shouldn't be aware that the playload is the hash of public key
	if !bytes.Equal(srcAddr.Payload(), pkHash[:]) {
		return errors.Wrapf(
			ErrAction,
			"signer public key hash %x does not match action source address payload %x",
			pkHash,
			srcAddr.Payload(),
		)
	}
	act.SetSrcPubkey(pk)
	hash := act.Hash()
	if act.SetSignature(crypto.EC283.Sign(sk, hash[:])); act.Signature() == nil {
		return errors.Wrapf(ErrAction, "failed to sign action hash = %x", hash)
	}
	return nil
}

// Verify verifies the action using sender's public key
func Verify(act Action) error {
	hash := act.Hash()
	if success := crypto.EC283.Verify(act.SrcPubkey(), hash[:], act.Signature()); success {
		return nil
	}
	return errors.Wrapf(
		ErrAction,
		"failed to verify action hash = %x and signature = %x",
		act.Hash(),
		act.Signature(),
	)
}
