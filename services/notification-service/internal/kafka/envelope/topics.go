package envelope

type Topic string

func (t Topic) String() string {
	return string(t)
}

const (
	TopicEmailSend   Topic = "notification.email.send"
	TopicEmailSent   Topic = "notification.email.sent"
	TopicEmailFailed Topic = "notification.email.failed"

	TopicSMSSend   Topic = "notification.sms.send"
	TopicSMSSent   Topic = "notification.sms.sent"
	TopicSMSFailed Topic = "notification.sms.failed"

	TopicPushSend   Topic = "notification.push.send"
	TopicPushSent   Topic = "notification.push.sent"
	TopicPushFailed Topic = "notification.push.failed"
)

const (
	TopicUserRegistered    Topic = "user.registered"
	TopicUserPasswordReset Topic = "user.password_reset_requested"
	TopicOrderCreated      Topic = "order.created"
	TopicOrderShipped      Topic = "order.shipped"
	TopicPaymentCompleted  Topic = "payment.completed"
)
