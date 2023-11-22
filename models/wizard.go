package models

type FormData struct {
	Step1Data string
	Step2Data string
	Step3Data string
}

type OvpnServerBaseSetting struct {
	OvpnEndpoint                       string // external ip
	OvpnTunNumber                      int
	ApproveIP                          bool //double check
	OvpnIPV6Support                    bool
	DisableDefRouteForClientsByDefault bool
	ClientToClientConfigIsUsed         bool
	OvpnMgmtAddress                    string //use def interna ip
	OvpnIPRange                        string
	OvpnPort                           string
	OvpnProtocol                       string
	OvpnDNS1                           string //use struct
	OvpnDNS2                           string //use struct
	OvpnCompression                    string //use string array
	CipherChoice                       string //use struct
	HMACAlgorithm                      string //use struct
	CertType                           string //use struct
	CertCurve                          string //use struct for ECDSA
	RSAKeySize                         string //use struct for RSA
	CCCipherChoice                     string //use sruct
	DHType                             string //use struct
	DHCurve                            string //use struct for ECDH
	DHKeySize                          string //use struct for DH
	TLSsig                             string //use array
}

var ovpnProtocolList = []string{"udp", "tcp"}

type DNSProvider struct {
	Name string
	DNS1 string
	DNS2 string
}

var dnsProviders = []DNSProvider{
	{"Google", "8.8.8.8", "8.8.4.4"},
	{"Cloudflare", "1.0.0.1", "1.1.1.1"},
	{"Quad9", "9.9.9.9", "149.112.112.112"},
	{"Quad9 uncensored", "9.9.9.10", "149.112.112.10"},
	{"FDN", "80.67.169.40", "80.67.169.12"},
	{"DNS.WATCH", "84.200.69.80", "84.200.70.40"},
	{"OpenDNS", "208.67.222.222", "208.67.220.220"},
	{"AdGuard DNS", "94.140.14.14", "94.140.15.15"},
	{"NextDNS", "45.90.28.167", "45.90.30.167"},
	{"Custom DNS", "", ""},
}

var compressionList = []string{"Disabled", "lz4-v2", "lz4", "lzo"}

type CipherList struct {
	Name      string
	Algorithm []string
}

var HMACAlgorithmList = []string{"SHA256", "SHA384", "SHA512"}

var cipherList = []CipherList{
	{"AES-128-GCM", HMACAlgorithmList},
	{"AES-192-GCM", HMACAlgorithmList},
	{"AES-256-GCM", HMACAlgorithmList},
	{"AES-128-CBC", nil},
	{"AES-192-CBC", nil},
	{"AES-256-CBC", nil},
}

var CertCurveList = []string{"prime256v1", "secp384r1", "secp521r1"}
var CCCipherECDSAList = []string{"TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256", "TLS-ECDHE-ECDSA-WITH-AES-256-GCM-SHA384"}
var RSAKeySizeList = []string{"2048", "3072", "4096"}
var CCCipherRSAList = []string{"TLS-ECDHE-RSA-WITH-AES-128-GCM-SHA256", "TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384"}

type CertType struct {
	Type     string
	Params   []string
	CCCipher []string
}

var certType = []CertType{
	{"ECDSA", CertCurveList, CCCipherECDSAList},
	{"RSA", RSAKeySizeList, CCCipherRSAList},
}

var DHCurveList = []string{"prime256v1", "secp384r1", "secp521r1"}
var DHKeySizeList = []string{"2048", "3072", "4096"}

type DHType struct {
	Name   string
	Params []string
}

var dhType = []DHType{
	{"ECDH", DHCurveList},
	{"DH", DHKeySizeList},
}

var TLSSigList = []string{"tls-crypt-v2", "tls-crypt", "tls-auth", "no tls"}

var defaultOVPNConfig = OvpnServerBaseSetting{
	OvpnEndpoint:                       "TODO",                     //ENDPOINT
	OvpnTunNumber:                      0,                          //TUN_NUMBER
	ApproveIP:                          false,                      //APPROVE_IP
	OvpnIPV6Support:                    false,                      //IPV6_SUPPORT
	DisableDefRouteForClientsByDefault: true,                       //DISABLE_DEF_ROUTE_FOR_CLIENTS
	ClientToClientConfigIsUsed:         true,                       //CLIENT_TO_CLIENT
	OvpnMgmtAddress:                    "openvpn:2080",             //SET_MGMT
	OvpnIPRange:                        "10.8.0.0",                 //IP_RANGE
	OvpnPort:                           "1194",                     //PORT
	OvpnProtocol:                       ovpnProtocolList[0],        //PROTOCOL
	OvpnDNS1:                           dnsProviders[0].DNS1,       //DNS1
	OvpnDNS2:                           dnsProviders[0].DNS2,       //DNS2
	OvpnCompression:                    compressionList[0],         //COMPRESSION_ALG
	CipherChoice:                       cipherList[0].Name,         //CIPHER_CHOICE
	HMACAlgorithm:                      cipherList[0].Algorithm[0], //HMAC_ALG_CHOICE
	CertType:                           certType[0].Type,           //CERT_TYPE
	CertCurve:                          certType[0].Params[0],      //CERT_CURVE
	RSAKeySize:                         certType[0].Params[0],      //RSA_KEY_SIZE
	CCCipherChoice:                     certType[0].CCCipher[0],    //CC_CIPHER
	DHType:                             dhType[0].Name,             //DH_TYPE
	DHCurve:                            dhType[0].Params[0],        //DH_CURVE
	DHKeySize:                          dhType[0].Params[0],        //DH_KEY_SIZE
	TLSsig:                             TLSSigList[0],              //TLS_SIG
}

func GetDefWizardSettings() OvpnServerBaseSetting {
	return defaultOVPNConfig
}
