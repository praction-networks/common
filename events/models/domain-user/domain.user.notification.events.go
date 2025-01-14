package domainUserEventModdel

type DomainUserNoEvent struct {
	NotificationId primitive.ObjectID `bson:"notificationId" json:"notificationId"`
	UserId         primitive.ObjectID `bson:"userId" json:"userId"`
}
