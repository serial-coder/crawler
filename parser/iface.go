/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package parser

import "github.com/hyperledger/fabric-protos-go/common"

// Parser serves as a contract for the implementation of the logic responsible for processing data from the blockchain.
// Simply put, this is about how exactly and into what constituent parts we will disassemble the blocks.
type Parser interface {
	// Parse is responsible for parsing passed block
	Parse(block *common.Block) (*Data, error)
}
