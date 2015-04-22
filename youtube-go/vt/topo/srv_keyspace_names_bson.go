// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

// DO NOT EDIT.
// FILE GENERATED BY BSONGEN.

import (
	"bytes"

	"github.com/youtube/vitess/go/bson"
	"github.com/youtube/vitess/go/bytes2"
)

// MarshalBson bson-encodes SrvKeyspaceNames.
func (srvKeyspaceNames *SrvKeyspaceNames) MarshalBson(buf *bytes2.ChunkedWriter, key string) {
	bson.EncodeOptionalPrefix(buf, bson.Object, key)
	lenWriter := bson.NewLenWriter(buf)

	// []string
	{
		bson.EncodePrefix(buf, bson.Array, "Entries")
		lenWriter := bson.NewLenWriter(buf)
		for _i, _v1 := range srvKeyspaceNames.Entries {
			bson.EncodeString(buf, bson.Itoa(_i), _v1)
		}
		lenWriter.Close()
	}

	lenWriter.Close()
}

// UnmarshalBson bson-decodes into SrvKeyspaceNames.
func (srvKeyspaceNames *SrvKeyspaceNames) UnmarshalBson(buf *bytes.Buffer, kind byte) {
	switch kind {
	case bson.EOO, bson.Object:
		// valid
	case bson.Null:
		return
	default:
		panic(bson.NewBsonError("unexpected kind %v for SrvKeyspaceNames", kind))
	}
	bson.Next(buf, 4)

	for kind := bson.NextByte(buf); kind != bson.EOO; kind = bson.NextByte(buf) {
		switch bson.ReadCString(buf) {
		case "Entries":
			// []string
			if kind != bson.Null {
				if kind != bson.Array {
					panic(bson.NewBsonError("unexpected kind %v for srvKeyspaceNames.Entries", kind))
				}
				bson.Next(buf, 4)
				srvKeyspaceNames.Entries = make([]string, 0, 8)
				for kind := bson.NextByte(buf); kind != bson.EOO; kind = bson.NextByte(buf) {
					bson.SkipIndex(buf)
					var _v1 string
					_v1 = bson.DecodeString(buf, kind)
					srvKeyspaceNames.Entries = append(srvKeyspaceNames.Entries, _v1)
				}
			}
		default:
			bson.Skip(buf, kind)
		}
	}
}
