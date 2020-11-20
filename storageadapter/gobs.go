/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package storageadapter

import (
	"bytes"
	"encoding/gob"
	"github.com/newity/crawler/parser"
)

func Encode(data *parser.Data) ([]byte, error) {
	var bytebuffer bytes.Buffer
	e := gob.NewEncoder(&bytebuffer)
	if err := e.Encode(data); err != nil {
		return nil, err
	}
	return bytebuffer.Bytes(), nil
}

func Decode(data []byte) (*parser.Data, error) {
	decoded := &parser.Data{}
	bytebuffer := bytes.NewBuffer(data)
	d := gob.NewDecoder(bytebuffer)
	if err := d.Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}
