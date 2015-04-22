// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zk

// DO NOT EDIT.
// FILE GENERATED BY BSONGEN.

import (
	"bytes"

	"github.com/youtube/vitess/go/bson"
	"github.com/youtube/vitess/go/bytes2"
)

// MarshalBson bson-encodes ZkNodeV.
func (zkNodeV *ZkNodeV) MarshalBson(buf *bytes2.ChunkedWriter, key string) {
	bson.EncodeOptionalPrefix(buf, bson.Object, key)
	lenWriter := bson.NewLenWriter(buf)

	// []*ZkNode
	{
		bson.EncodePrefix(buf, bson.Array, "Nodes")
		lenWriter := bson.NewLenWriter(buf)
		for _i, _v1 := range zkNodeV.Nodes {
			// *ZkNode
			if _v1 == nil {
				bson.EncodePrefix(buf, bson.Null, bson.Itoa(_i))
			} else {
				(*_v1).MarshalBson(buf, bson.Itoa(_i))
			}
		}
		lenWriter.Close()
	}

	lenWriter.Close()
}

// UnmarshalBson bson-decodes into ZkNodeV.
func (zkNodeV *ZkNodeV) UnmarshalBson(buf *bytes.Buffer, kind byte) {
	switch kind {
	case bson.EOO, bson.Object:
		// valid
	case bson.Null:
		return
	default:
		panic(bson.NewBsonError("unexpected kind %v for ZkNodeV", kind))
	}
	bson.Next(buf, 4)

	for kind := bson.NextByte(buf); kind != bson.EOO; kind = bson.NextByte(buf) {
		switch bson.ReadCString(buf) {
		case "Nodes":
			// []*ZkNode
			if kind != bson.Null {
				if kind != bson.Array {
					panic(bson.NewBsonError("unexpected kind %v for zkNodeV.Nodes", kind))
				}
				bson.Next(buf, 4)
				zkNodeV.Nodes = make([]*ZkNode, 0, 8)
				for kind := bson.NextByte(buf); kind != bson.EOO; kind = bson.NextByte(buf) {
					bson.SkipIndex(buf)
					var _v1 *ZkNode
					// *ZkNode
					if kind != bson.Null {
						_v1 = new(ZkNode)
						(*_v1).UnmarshalBson(buf, kind)
					}
					zkNodeV.Nodes = append(zkNodeV.Nodes, _v1)
				}
			}
		default:
			bson.Skip(buf, kind)
		}
	}
}
