package ledger

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromProtoTimestamp(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
