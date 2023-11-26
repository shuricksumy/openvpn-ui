package models

type OvpnServerBaseSetting struct {
	OvpnEndpoint                       string // external ip
	ApproveIP                          bool   //double check TO DELETE
	DisableDefRouteForClientsByDefault bool
	ClientToClientConfigIsUsed         bool
	OvpnMgmtAddress                    string //use def interna ip
	OvpnIPRange                        string
	OvpnPort                           string
	TunNumber                          string
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

func GetovpnProtocolList(selected string) []string {
	var actualOvpnProtocolList []string
	for _, p := range ovpnProtocolList {
		if selected == p {
			continue
		}
		actualOvpnProtocolList = append(actualOvpnProtocolList, p)
	}
	return actualOvpnProtocolList
}

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

func GetDNSProvidersList(selected string) []DNSProvider {
	var actualdnsProviders []DNSProvider
	for _, d := range dnsProviders {
		if selected == d.Name {
			continue
		}
		actualdnsProviders = append(actualdnsProviders, d)
	}
	return actualdnsProviders
}

func GetDNSProvidersNameByIP(ip string) DNSProvider {
	for _, d := range dnsProviders {
		if ip == d.DNS1 {
			return d
		}
	}
	var customProvider = DNSProvider{"Custom DNS", "", ""}
	return customProvider
}

var compressionList = []string{"Disabled", "lz4-v2", "lz4", "lzo"}

func GetOvpnCompressionList(selected string) []string {
	var actualCompressionList []string
	for _, c := range compressionList {
		if selected == c {
			continue
		}
		actualCompressionList = append(actualCompressionList, c)
	}
	return actualCompressionList
}

type CipherList struct {
	Name      string
	Algorithm []string
}

var HMACAlgorithmList = []string{"SHA256", "SHA384", "SHA512"}

var cipherList = []CipherList{
	{"AES-128-GCM", HMACAlgorithmList},
	{"AES-192-GCM", HMACAlgorithmList},
	{"AES-256-GCM", HMACAlgorithmList},
	{"AES-128-CBC", HMACAlgorithmList},
	{"AES-192-CBC", HMACAlgorithmList},
	{"AES-256-CBC", HMACAlgorithmList},
}

func GetCipherChoiceList(selected string) []CipherList {
	var actualCipherList []CipherList
	for _, c := range cipherList {
		if selected == c.Name {
			continue
		}
		actualCipherList = append(actualCipherList, c)
	}
	return actualCipherList
}

func GetHMACAlgorithmList(selectedCipher string, selectedAlgorithm string) []string {
	var actualAlgorithmList []string
	for _, c := range cipherList {
		if selectedCipher == c.Name {
			for _, a := range c.Algorithm {
				if selectedAlgorithm == a {
					continue
				}
				actualAlgorithmList = append(actualAlgorithmList, a)
			}
		}
	}
	return actualAlgorithmList
}

func GetHMACAlgorithmAllList() []string {
	return HMACAlgorithmList
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

func GetCertTypeList(selected string) []CertType {
	var actualCertTypeList []CertType
	for _, c := range certType {
		if selected == c.Type {
			continue
		}
		actualCertTypeList = append(actualCertTypeList, c)
	}
	return actualCertTypeList
}

var DHCurveList = []string{"prime256v1", "secp384r1", "secp521r1"}

func GetCertParamList(selectedCertType string, selectedParam string) []string {
	var actualCertParamList []string
	for _, c := range certType {
		if selectedCertType == c.Type {
			for _, p := range c.Params {
				if selectedParam == p {
					continue
				}
				actualCertParamList = append(actualCertParamList, p)
			}
		}
	}
	return actualCertParamList
}

func GetCCCipherChoiceList(selectedCertType string, selectedCipher string) []string {
	var actualCertCipherList []string
	for _, c := range certType {
		if selectedCertType == c.Type {
			for _, p := range c.CCCipher {
				if selectedCipher == p {
					continue
				}
				actualCertCipherList = append(actualCertCipherList, p)
			}
		}
	}
	return actualCertCipherList

}

var DHKeySizeList = []string{"2048", "3072", "4096"}

type DHType struct {
	Name   string
	Params []string
}

var dhType = []DHType{
	{"ECDH", DHCurveList},
	{"DH", DHKeySizeList},
}

func GetDHTypeList(selected string) []DHType {
	var actualDHTypeList []DHType
	for _, d := range dhType {
		if selected == d.Name {
			continue
		}
		actualDHTypeList = append(actualDHTypeList, d)
	}
	return actualDHTypeList
}

func GetDHParamList(selectedDHType string, selectedParam string) []string {
	var actualDHParams []string
	for _, d := range dhType {
		if selectedDHType == d.Name {
			for _, p := range d.Params {
				if selectedParam == p {
					continue
				}
				actualDHParams = append(actualDHParams, p)
			}
		}
	}
	return actualDHParams

}

var TLSSigList = []string{"tls-crypt-v2", "tls-crypt", "tls-auth", "no tls"}

func GetTLSsigList(selected string) []string {
	var actualTLSsigList []string
	for _, c := range TLSSigList {
		if selected == c {
			continue
		}
		actualTLSsigList = append(actualTLSsigList, c)
	}
	return actualTLSsigList
}

var defaultOVPNConfig = OvpnServerBaseSetting{
	OvpnEndpoint:                       "TODO",                     //ENDPOINT
	ApproveIP:                          false,                      //APPROVE_IP
	DisableDefRouteForClientsByDefault: true,                       //DISABLE_DEF_ROUTE_FOR_CLIENTS
	ClientToClientConfigIsUsed:         true,                       //CLIENT_TO_CLIENT
	OvpnMgmtAddress:                    "openvpn:2080",             //SET_MGMT
	OvpnIPRange:                        "10.8.0.0",                 //IP_RANGE
	OvpnPort:                           "1194",                     //PORT
	TunNumber:                          "0",                        //TUN_NUMBER
	OvpnProtocol:                       ovpnProtocolList[0],        //PROTOCOL
	OvpnDNS1:                           dnsProviders[0].DNS1,       //DNS1
	OvpnDNS2:                           dnsProviders[0].DNS2,       //DNS2
	OvpnCompression:                    compressionList[0],         //COMPRESSION_ALG
	CipherChoice:                       cipherList[0].Name,         //CIPHER_CHOICE
	HMACAlgorithm:                      cipherList[0].Algorithm[0], //HMAC_ALG_CHOICE
	CertType:                           certType[0].Type,           //CERT_TYPE
	CertCurve:                          certType[0].Params[0],      //CERT_CURVE
	RSAKeySize:                         certType[0].Params[0],      //RSA_KEY_SIZE
	CCCipherChoice:                     certType[0].CCCipher[0],    //CC_CIPHER_CHOICE
	DHType:                             dhType[0].Name,             //DH_TYPE
	DHCurve:                            dhType[0].Params[0],        //DH_CURVE
	DHKeySize:                          dhType[0].Params[0],        //DH_KEY_SIZE
	TLSsig:                             TLSSigList[0],              //TLS_SIG
}

func GetDefWizardSettings() OvpnServerBaseSetting {
	return defaultOVPNConfig
}

var indexMapper = map[string]string{
	"udp":              "1", //PROTOCOL_CHOICE
	"tcp":              "2",
	"Google":           "9", //DNS
	"Cloudflare":       "3",
	"Quad9":            "4",
	"Quad9 uncensored": "5",
	"FDN":              "6",
	"DNS.WATCH":        "7",
	"OpenDNS":          "8",
	"AdGuard DNS":      "11",
	"NextDNS":          "12",
	"Custom DNS":       "13",
	"ECDSA":            "1", //Cert type
	"RSA":              "2",
	"ECDH":             "1", //DH type
	"DH":               "2",
	"tls-crypt":        "1", //TLS sig
	"tls-auth":         "2",
	"tls-crypt-v2":     "3",
	"no tls":           "4",
	"AES-128-GCM":      "1", //CIPHER_CHOICE
	"AES-192-GCM":      "2",
	"AES-256-GCM":      "3",
	"AES-128-CBC":      "4",
	"AES-192-CBC":      "5",
	"AES-256-CBC":      "6",
	"prime256v1":       "1", //CERT_CURVE_CHOICE
	"secp384r1":        "2",
	"secp521r1":        "3",
	"2048":             "1", //RSA_KEY_SIZE
	"3072":             "2",
	"4096":             "3",
	"TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256": "1", //CC_CIPHER_CHOICE
	"TLS-ECDHE-ECDSA-WITH-AES-256-GCM-SHA384": "2",
	"TLS-ECDHE-RSA-WITH-AES-128-GCM-SHA256":   "1",
	"TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384":   "2",
	"SHA256":                                  "1", //HMAC_ALG_CHOICE
	"SHA384":                                  "2",
	"SHA512":                                  "3",
	"lz4-v2":                                  "1", //COMPRESSION_CHOICE
	"lz4":                                     "2",
	"lzo":                                     "3",
}

func GetIndex(key string) string {
	return indexMapper[key]
}

func GetConstEnv() string {

	const constantEnvVars = `#!/bin/bash
export DOCKER_COMMAND="2"
export IP_CHOICE="2"
export IPV6_SUPPORT="n"
export PORT_CHOICE="2"
export DNS="13"
export CUSTOMIZE_ENC="y"
export SET_MGMT="management 127.0.0.1 2080"
`
	return constantEnvVars
}

func GetScriptEnv() string {

	const scriptInstall = `
if [[ ! -e /etc/openvpn/server.conf ]]; then
        echo "Waiting for the good weather... ;)"
        bash /opt/scripts/openvpn-install-v2.sh
        echo "" > /etc/openvpn/.provisioned
        echo "    --= SETUP IS DONE ==-"
fi
`
	return scriptInstall
}

// # Apply FireWall rules
// if [[ -f /etc/openvpn/set_fw.sh ]]; then
//      bash /etc/openvpn/set_fw.sh
// fi

// /usr/sbin/openvpn  --cd /etc/openvpn --config /etc/openvpn/server.conf
