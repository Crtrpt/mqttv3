package packets

const (
	PayloadFormatIndicator          = 1
	MessageExpiryInterval           = 2
	ContentType                     = 3
	ResponseTopic                   = 8
	CorrelationData                 = 9
	SubscriptionIdentifier          = 11
	SessionExpiryInterval           = 17
	AssignedClientIdentifier        = 18
	ServerKeepAlive                 = 19
	AuthenticationMethod            = 21
	AuthenticationData              = 22
	RequestProblemInformation       = 23
	WillDelayInterval               = 24
	RequestResponseInformation      = 26
	ResponseInformation             = 28
	ServerReference                 = 30
	ReasonString                    = 31
	ReceiveMaximum                  = 33
	TopicAliasMaximum               = 34
	TopicAlias                      = 35
	MaximumQoS                      = 36
	RetainAvailable                 = 37
	UserProperty                    = 38
	MaximumPacketSize               = 39
	WildcardSubscriptionAvailable   = 40
	SubscriptionIdentifierAvailable = 41
	SharedSubscriptionAvailable     = 42
)
