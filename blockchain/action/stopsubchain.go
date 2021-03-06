// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"

	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/pkg/version"
	"github.com/iotexproject/iotex-core/proto"
)

// StopSubChain defines the action to stop sub chain
type StopSubChain struct {
	action
	chainID    uint32
	stopHeight uint64
}

// NewStopSubChain returns a StopSubChain instance
func NewStopSubChain(senderAddress string, nonce uint64, chainID uint32, chainAddress string, stopHeight uint64, gasLimit uint64, gasPrice *big.Int) (*StopSubChain, error) {
	return &StopSubChain{
		action: action{
			version:  version.ProtocolVersion,
			nonce:    nonce,
			srcAddr:  senderAddress,
			dstAddr:  chainAddress,
			gasLimit: gasLimit,
			gasPrice: gasPrice,
		},
		chainID:    chainID,
		stopHeight: stopHeight,
	}, nil
}

// ChainAddress returns the address of the sub chain
func (ssc *StopSubChain) ChainAddress() string {
	return ssc.dstAddr
}

// ChainID returns the id of the sub chain
func (ssc *StopSubChain) ChainID() uint32 {
	return ssc.chainID
}

// StopHeight returns the height to stop the sub chain
func (ssc *StopSubChain) StopHeight() uint64 {
	return ssc.stopHeight
}

// TotalSize returns the total size of this instance
func (ssc *StopSubChain) TotalSize() uint32 {
	size := NonceSizeInBytes
	size += VersionSizeInBytes
	size += len(ssc.srcPubkey)
	size += len(ssc.srcAddr)
	size += len(ssc.dstAddr)
	size += GasSizeInBytes
	if ssc.gasPrice != nil && len(ssc.gasPrice.Bytes()) > 0 {
		size += len(ssc.gasPrice.Bytes())
	}
	size += len(ssc.signature)
	return uint32(size) + 4 + 8 // chain id size + stop height size
}

// ByteStream returns a raw byte stream of this instance
func (ssc *StopSubChain) ByteStream() []byte {
	stream := byteutil.Uint32ToBytes(ssc.version)
	stream = append(stream, byteutil.Uint64ToBytes(ssc.nonce)...)
	stream = append(stream, byteutil.Uint64ToBytes(ssc.gasLimit)...)
	stream = append(stream, ssc.srcPubkey[:]...)
	stream = append(stream, ssc.srcAddr...)
	stream = append(stream, ssc.dstAddr...)
	if ssc.gasPrice != nil && len(ssc.gasPrice.Bytes()) > 0 {
		stream = append(stream, ssc.gasPrice.Bytes()...)
	}
	stream = append(stream, byteutil.Uint32ToBytes(ssc.chainID)...)

	return append(stream, byteutil.Uint64ToBytes(ssc.stopHeight)...)
}

// ConvertToActionPb converts StopSubChain to protobuf's ActionPb
func (ssc *StopSubChain) ConvertToActionPb() *iproto.ActionPb {
	pbSSC := &iproto.ActionPb{
		Action: &iproto.ActionPb_StopSubChain{
			StopSubChain: &iproto.StopSubChainPb{
				ChainID:         ssc.chainID,
				StopHeight:      ssc.stopHeight,
				Owner:           ssc.srcAddr,
				OwnerPublicKey:  ssc.srcPubkey[:],
				SubChainAddress: ssc.dstAddr,
			},
		},
		Version:   ssc.version,
		Nonce:     ssc.nonce,
		GasLimit:  ssc.gasLimit,
		Signature: ssc.signature,
	}
	if ssc.gasPrice != nil {
		pbSSC.GasPrice = ssc.gasPrice.Bytes()
	}
	return pbSSC
}

// Serialize returns a serialized byte stream for the StopSubChain
func (ssc *StopSubChain) Serialize() ([]byte, error) {
	return proto.Marshal(ssc.ConvertToActionPb())
}

// ConvertFromActionPb converts a protobuf's ActionPb to StopSubChain
func (ssc *StopSubChain) ConvertFromActionPb(pbAct *iproto.ActionPb) {
	ssc.version = pbAct.Version
	ssc.nonce = pbAct.Nonce
	ssc.gasLimit = pbAct.GasLimit
	if ssc.gasPrice == nil {
		ssc.gasPrice = big.NewInt(0)
	}
	if len(pbAct.GasPrice) > 0 {
		ssc.gasPrice.SetBytes(pbAct.GasPrice)
	}
	ssc.signature = pbAct.Signature
	pbSSC := pbAct.GetStopSubChain()
	if pbSSC != nil {
		ssc.chainID = pbSSC.ChainID
		ssc.stopHeight = pbSSC.StopHeight
		ssc.srcAddr = pbSSC.Owner
		copy(ssc.srcPubkey[:], pbSSC.OwnerPublicKey)
		ssc.dstAddr = pbSSC.SubChainAddress
	}
}

// Deserialize parse the byte stream into StopSubChain
func (ssc *StopSubChain) Deserialize(buf []byte) error {
	pbSSC := &iproto.ActionPb{}
	if err := proto.Unmarshal(buf, pbSSC); err != nil {
		return err
	}
	ssc.ConvertFromActionPb(pbSSC)
	return nil
}

// Hash returns the hash of the StopSubChain
func (ssc *StopSubChain) Hash() hash.Hash32B {
	return blake2b.Sum256(ssc.ByteStream())
}

// IntrinsicGas returns the intrinsic gas of a StopSubChain
func (ssc *StopSubChain) IntrinsicGas() (uint64, error) {
	return StopSubChainIntrinsicGas, nil
}

// Cost returns the total cost of a StopSubChain
func (ssc *StopSubChain) Cost() (*big.Int, error) {
	intrinsicGas, err := ssc.IntrinsicGas()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get intrinsic gas for the StopSubChain action")
	}
	fee := big.NewInt(0).Mul(ssc.GasPrice(), big.NewInt(0).SetUint64(intrinsicGas))
	return fee, nil
}
