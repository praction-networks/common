package events

import (
	"context"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

// FailedNATSEvent uses a string _id like "AuthStream|tenantuser.created:abcd1234..."

// ReplayFallbackOnce processes up to `limit` documents once.
// On success -> delete doc. On failure -> Attempts++, LastError, Timestamp=now.
func ReplayFallbackOnce(ctx context.Context, jsm *JsStreamManager, coll *mongo.Collection, limit int) (processed, ok, failed int, err error) {
	if coll == nil {
		return 0, 0, 0, nil
	}
	if limit <= 0 {
		limit = 100
	}

	cur, err := coll.Find(ctx, bson.M{}, mopt.Find().
		SetSort(bson.D{{Key: "timestamp", Value: 1}, {Key: "attempts", Value: 1}}).
		SetLimit(int64(limit)))
	if err != nil {
		return 0, 0, 0, err
	}
	defer cur.Close(ctx)

	now := time.Now()

	for cur.Next(ctx) {
		processed++

		var d FailedNATSEvent
		if derr := cur.Decode(&d); derr != nil {
			logger.Error("fallback decode failed", derr)
			failed++
			continue
		}

		// MsgID is the part after the first '|'
		msgID := d.ID
		if i := strings.IndexByte(d.ID, '|'); i >= 0 && i+1 < len(d.ID) {
			msgID = d.ID[i+1:]
		}

		jsOpts := []jetstream.PublishOpt{
			jetstream.WithMsgID(msgID),               // idempotent store on retries
			jetstream.WithExpectStream(d.StreamName), // safety: ensure correct stream
		}

		if ack, perr := jsm.JsClient.Publish(ctx, d.Subject, d.Payload, jsOpts...); perr == nil {
			// success (Duplicate=true is also success)
			if _, derr := coll.DeleteOne(ctx, bson.M{"_id": d.ID}); derr != nil {
				logger.Error("fallback delete failed after success", derr, "id", d.ID)
				// not counting as failure of publish; keep ok++
				ok++
				continue
			}
			logger.Info("republished from fallback",
				"stream", ack.Stream, "seq", ack.Sequence, "duplicate", ack.Duplicate,
				"subject", d.Subject, "id", d.ID)
			ok++
			continue
		} else {
			// one-shot failed: bump attempts and lastError
			_, uerr := coll.UpdateByID(ctx, d.ID, bson.M{
				"$set": bson.M{
					"lastError": perr.Error(),
					"timestamp": now,
				},
				"$inc": bson.M{"attempts": 1},
			})
			if uerr != nil {
				logger.Error("fallback update failed", uerr, "id", d.ID)
			}
			failed++
		}
	}
	return
}
