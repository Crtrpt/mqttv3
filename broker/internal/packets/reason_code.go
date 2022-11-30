package packets

type ReasonCode byte

const (
	Success                   = 0  //CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK, AUTH
	NormalDisconnection       = 0  //DISCONNECT
	GrantedQoS0               = 0  //SUBACK
	GrantedQoS1               = 1  //SUBACK
	GrantedQoS2               = 2  //SUBACK
	DisconnectWithWillMessage = 4  //DISCONNECT
	NoMatchingSubscribers     = 16 //PUBACK, PUBREC
	NoSubscriptionExisted     = 17 //UNSUBACK
	ContinueAuthentication    = 24 //AUTH
	ReAuthenticate            = 25 //AUTH
)

const (
	UnspecifiedError            = iota + 128 //CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
	MalformedPacket                          //CONNACK, DISCONNECT
	ProtocolError                            //CONNACK, DISCONNECT
	ImplementationSpecificError              //CONNACK, PUBACK, PUBREC, SUBACK, UNSUBACK, DISCONNECT
	UnsupportedProtocolVersion
	ClientIdentifierNotValid
	BadUserNameOrPassword
	NotAuthorized
	ServerUnavailable
	ServerBusy
	Banned
	ServerShuttingDown
	BadAuthenticationMethod
	KeepAliveTimeout
	SessionTakenOver
	TopicFilterInvalid
	TopicNameInvalid
	PacketIdentifierInUse
	PacketIdentifierNotFound
	ReceiveMaximumExceeded
	TopicAliasInvalid
	PacketTooLarge
	MessageRateTooHigh
	QuotaExceeded
	AdministrativeAction
	PayloadFormatInvalid
	RetainNotSupported
	QoSNotSupported
	UseAnotherServer
	ServerMoved
	SharedSubscriptionsNotSupported
	ConnectionRateExceeded
	MaximumConnectTime
	SubscriptionIdentifiersNotSupported
	WildcardSubscriptionsNotSupported
)
