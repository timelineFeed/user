package biz

import (
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

func Test_pbAny2goAny(t *testing.T) {

	Extra := make(map[string]*anypb.Any)
	Extra["11"] = &anypb.Any{
		TypeUrl: "1111",
		Value:   []byte("hhh"),
	}
	Extra["22"] = &anypb.Any{
		TypeUrl: "2222",
		Value:   []byte("kkk"),
	}
	extra := pbAny2goAny(Extra)
	t.Logf("nice extra=%+v", extra)
}
