package loop3

import (
	"bytes"
	"crypto/sha512"
	"github.com/google/go-cmp/cmp"
	loop3_pb "github.com/openziti/ziti/ziti-fabric-test/subcmd/loop3/pb"
	"github.com/stretchr/testify/require"
	"math/rand"
	"reflect"
	"testing"
)

type testPeer struct {
	bytes.Buffer
}

func (t *testPeer) Close() error {
	return nil
}

func Test_MessageSerDeser(t *testing.T) {
	req := require.New(t)
	data := make([]byte, 4192)
	rand.Read(data)

	hash := sha512.Sum512(data)

	block := &Block{
		Type:     BlockTypePlain,
		Sequence: 10,
		Hash:     hash[:],
		Data:     data,
	}

	testBuf := &testPeer{}

	p := &protocol{
		peer: testBuf,
		test: &loop3_pb.Test{
			Name: "test",
		},
	}

	req.NoError(block.Tx(p))

	readBlock := &Block{}
	req.NoError(readBlock.Rx(p))

	req.True(reflect.DeepEqual(block, readBlock), cmp.Diff(block, readBlock))

	data = make([]byte, 4192)
	rand.Read(data)
	hash = sha512.Sum512(data)

	block = &Block{
		Type:     BlockTypeLatencyRequest,
		Sequence: 10,
		Hash:     hash[:],
		Data:     data,
	}

	req.NoError(block.Tx(p))

	readBlock = &Block{}
	req.NoError(readBlock.Rx(p))

	req.Equal("", cmp.Diff(block, readBlock))
}
