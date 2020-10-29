/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package blocklib

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestFromFabricBlock(t *testing.T) {
	file, err := ioutil.ReadFile("./mock/genesis.pb")
	assert.NoError(t, err)

	fabBlock := &common.Block{}
	assert.NoError(t, proto.Unmarshal(file, fabBlock))

	block, err := FromFabricBlock(fabBlock)
	assert.NoError(t, err)
	assert.NotNil(t, block.Data)
	assert.NotNil(t, block.Metadata)
	assert.NotNil(t, block.Signatures)
}

func TestBlockSignatures(t *testing.T) {
	file, err := ioutil.ReadFile("./mock/genesis.pb")
	assert.NoError(t, err)

	fabBlock := &common.Block{}
	assert.NoError(t, proto.Unmarshal(file, fabBlock))

	block, err := FromFabricBlock(fabBlock)
	assert.NoError(t, err)

	sigs, err := block.BlockSignatures()
	assert.NoError(t, err)

	for _, sig := range sigs {
		assert.Equal(t, uint64(12652116863344733010), binary.BigEndian.Uint64(sig.Nonce))
		assert.Equal(t, "OrdererMSP", sig.MSPID)
		certhash := sha256.New()
		_, err = certhash.Write(sig.Cert)
		assert.NoError(t, err)
		sighash := sha256.New()
		_, err = sighash.Write(sig.Signature)
		assert.NoError(t, err)
		assert.Equal(t, "9fa97f0795f3ade55fbf89419e7c8afaf2135b58bdf8ad39728a173ec498bfaf", hex.EncodeToString(certhash.Sum(nil)))
		assert.Equal(t, "3d82fb80a17f2c189d088330d2459efe66e7a409ebd4ccbf2f32f5a768059345", hex.EncodeToString(sighash.Sum(nil)))
	}
}

func TestTxs(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		file, err := ioutil.ReadFile("./mock/sampleblock.pb")
		assert.NoError(t, err)

		fabBlock := &common.Block{}
		assert.NoError(t, proto.Unmarshal(file, fabBlock))

		block, err := FromFabricBlock(fabBlock)
		assert.NoError(t, err)

		txs, err := block.Txs()
		assert.NoError(t, err)

		for _, tx := range txs {
			assert.NotNil(t, tx.Data)
			assert.Equal(t, int32(0), tx.ValidationCode())
			assert.Equal(t, "VALID", tx.ValidationStatus())
		}
	})

	t.Run("with MVCC_READ_CONFLICT", func(t *testing.T) {
		file, err := ioutil.ReadFile("./mock/mvcc_read_conflict.pb")
		assert.NoError(t, err)

		fabBlock := &common.Block{}
		assert.NoError(t, proto.Unmarshal(file, fabBlock))

		block, err := FromFabricBlock(fabBlock)
		assert.NoError(t, err)

		txs, err := block.Txs()
		assert.NoError(t, err)

		for _, tx := range txs {
			assert.NotNil(t, tx.Data)
			assert.Equal(t, int32(11), tx.ValidationCode())
			assert.Equal(t, "MVCC_READ_CONFLICT", tx.ValidationStatus())
		}
	})
}
