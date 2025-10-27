package envelope

// Topic represents a Kafka topic name.
type Topic string

// String returns the string representation of the topic.
func (t Topic) String() string {
	return string(t)
}

// Auth Service Topics - Outbound (Auth Service Publishes)
const (
	TopicUserRegistered     Topic = "user.registered"
	TopicUserPasswordReset  Topic = "user.password_reset_requested"
	TopicUserEmailVerified  Topic = "user.email_verified"
	TopicUserLoggedIn       Topic = "user.logged_in"
	TopicUserLoggedOut      Topic = "user.logged_out"
	TopicUserProfileUpdated Topic = "user.profile_updated"
	TopicUserDeleted        Topic = "user.deleted"
)

// External Event Topics - Inbound (Auth Service Consumes)
const (
	TopicEmailSent   Topic = "notification.email.sent"
	TopicEmailFailed Topic = "notification.email.failed"
)
