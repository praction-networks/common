package jobtransferevent

// JobTransferOffer is one peer-offer inside a JobTransferRequest. Multiple
// offers may exist; only one transitions to ACCEPTED.
type JobTransferOffer struct {
	ID            string `json:"id"                      bson:"_id"`
	RequestID     string `json:"requestId"               bson:"requestId"`
	ToUserID      string `json:"toUserId"                bson:"toUserId"`
	Status        string `json:"status"                  bson:"status"` // OFFERED | ACCEPTED | DECLINED | EXPIRED
	OfferedAtMs   int64  `json:"offeredAtMs"             bson:"offeredAtMs"`
	ExpiresAtMs   int64  `json:"expiresAtMs"             bson:"expiresAtMs"` // per-offer timeout (tenant policy: peerOfferMinutes)
	RespondedAtMs int64  `json:"respondedAtMs,omitempty" bson:"respondedAtMs,omitempty"`
}
