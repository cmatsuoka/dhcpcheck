package dhcp

const (
	// RFC 1497 Vendor Extensions
	PadOption              = 0
	EndOption              = 255
	SubnetMask             = 1
	TimeOffset             = 2
	Router                 = 3
	TimeServer             = 4
	NameServer             = 5
	DomainNameServer       = 6
	LogServer              = 7
	CookieServer           = 8
	LPRServer              = 9
	ImpressServer          = 10
	ResourceLocationServer = 11
	HostName               = 12
	BootFileSize           = 13
	MeritDumpFile          = 14
	DomainName             = 15
	SwapServer             = 16
	RootPath               = 17
	ExtensionsPath         = 18

	// IP Layer Parameters per Host
	IPForwarding                  = 19
	NonLocalSourceRouting         = 20
	PolicyFilter                  = 21
	MaximumDatagramReassemblySize = 22
	DefaultIPTimeToLive           = 23
	PathMTUAgingTimeout           = 24
	PathMTUPlateauTable           = 25

	// IP Layer Parameters per Interface
	InterfaceMTU              = 26
	AllSubnetsAreLocal        = 27
	BroadcastAddress          = 28
	PerformMaskDiscovery      = 29
	MaskSupplier              = 30
	PerformRouterDiscovery    = 31
	RouterSolicitationAddress = 32
	StaticRoute               = 33

	// Link Layer Parameters per Interface
	TrailerEncapsulation  = 34
	ARPCacheTimeout       = 35
	EthernetEncapsulation = 36

	// TCP Parameters
	TCPDefaultTTL        = 37
	TCPKeepaliveInterval = 38
	TCPKeepaliveGarbage  = 39

	// Application and Service Parameters
	NISDomain                 = 40
	NISServers                = 41
	NTPServers                = 42
	VendorSpecificInformation = 43
	NetBIOSNameServer         = 44
	NetBIOSDatagramServer     = 45
	NetBIOSNodeType           = 46
	NetBIOSScope              = 47
	XFontServer               = 48
	XDisplayManager           = 49
	NISPlusDomain             = 64
	NISPlusServers            = 65
	MobileIPHomeAgent         = 68
	SMTPServer                = 69
	POP3Server                = 70
	NNTPServer                = 71
	DefaultWWWServer          = 72
	DefaultFingerServer       = 73
	DefaultIRCServer          = 74
	StreetTalkServer          = 75
	STDAServer                = 76

	// DHCP Extensions
	RequestedIPAddress     = 50
	IPAddressLeaseTime     = 51
	OptionOverload         = 52
	TFTPServerName         = 66
	BootfileName           = 67
	DHCPMessageType        = 53
	ServerIdentifier       = 54
	ParameterRequestList   = 55
	Message                = 56
	MaximumDHCPMessageSize = 57
	RenewalTimeValue       = 58
	RebindingTimeValue     = 59
	VendorClassIdentifier  = 60
	ClientIdentifier       = 61

	// RFC3397
	DomainSearch = 119

	// Web Proxy Auto-Discovery Protocol (ietf-wrec-wpad-01)
	WebProxyServer = 252
)
