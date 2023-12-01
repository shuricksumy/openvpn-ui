package lib

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/shuricksumy/openvpn-ui/state"
)

// Cert
// https://groups.google.com/d/msg/mailing.openssl.users/gMRbePiuwV0/wTASgPhuPzkJ
type Cert struct {
	EntryType   string
	Expiration  string
	ExpirationT time.Time
	Revocation  string
	RevocationT time.Time
	Serial      string
	FileName    string
	Details     *Details
}

type NameSorter []*Cert

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].Details.CN < a[j].Details.CN }

type Details struct {
	Name             string
	CN               string
	Country          string
	State            string
	City             string
	Organisation     string
	OrganisationUnit string
	Email            string
	LocalIP          string
	Description      string
}

func ReadCerts(path string) ([]*Cert, error) {
	certs_v := make([]*Cert, 0)
	certs_r := make([]*Cert, 0)
	certs_err := make([]*Cert, 0)
	certs := make([]*Cert, 0)
	certsSorted := make([]*Cert, 0)

	text, err := os.ReadFile(path)
	if err != nil {
		return certs, err
	}
	lines := strings.Split(trim(string(text)), "\n")

	// Validate file
	for i, line := range lines {
		//ff := strings.Split(trim(line), "\t")
		ff := strings.Fields(trim(line))

		if ff[0] == "V" {
			if len(ff) == 5 {
				continue
			} else {
				return certs,
					fmt.Errorf("V-validator: Incorrect number of lines - file 'index.txt'. Error in %d-line: \n%s\n. Expected %d, found %d",
						i, line, 5, len(ff))
			}
		}

		if ff[0] == "R" {
			if len(ff) == 6 {
				continue
			} else {
				return certs,
					fmt.Errorf("R-validator: Incorrect number of lines - file 'index.txt'. Error in %d-line: \n%s\n. Expected %d, found %d",
						i, line, 6, len(ff))
			}
		}

		if len(ff) == 5 {
			continue
		} else {
			return certs,
				fmt.Errorf("Other: Incorrect number of lines - file 'index.txt'. Error in %d-line: \n%s\n. Expected %d, found %d",
					i, line, 5, len(ff))
		}
	}

	for i, line := range lines {
		// Skip first item - server cert
		if i == 0 {
			continue
		}

		fields := strings.Fields(trim(line))
		//strings.Split(trim(line), "\t")
		// if cert is valid
		if fields[0] == "V" {
			expT, _ := time.Parse("060102150405Z", fields[1])
			revT, _ := time.Parse("060102150405Z", fields[1])
			c := &Cert{
				EntryType:   fields[0],
				Expiration:  fields[1],
				ExpirationT: expT,
				Revocation:  fields[1],
				RevocationT: revT,
				Serial:      fields[2],
				FileName:    fields[3],
				Details:     parseDetails(fields[4]),
			}
			certs = append(certs, c)
		} else if fields[0] == "R" {
			expT, _ := time.Parse("060102150405Z", fields[1])
			revT, _ := time.Parse("060102150405Z", fields[2])
			c := &Cert{
				EntryType:   fields[0],
				Expiration:  fields[1],
				ExpirationT: expT,
				Revocation:  fields[2],
				RevocationT: revT,
				Serial:      fields[3],
				FileName:    fields[4],
				Details:     parseDetails(fields[5]),
			}
			certs = append(certs, c)
		} else {
			expT, _ := time.Parse("060102150405Z", fields[1])
			revT, _ := time.Parse("060102150405Z", fields[1])
			c := &Cert{
				EntryType:   fields[0],
				Expiration:  fields[1],
				ExpirationT: expT,
				Revocation:  fields[1],
				RevocationT: revT,
				Serial:      fields[2],
				FileName:    fields[3],
				Details:     parseDetails(fields[4]),
			}
			certs = append(certs, c)
		}
	}

	for _, cc := range certs {
		if cc.EntryType == "V" {
			certs_v = append(certs_v, cc)
		} else if cc.EntryType == "R" {
			certs_r = append(certs_r, cc)
		} else {
			certs_err = append(certs_err, cc)
		}
	}

	sort.Sort(NameSorter(certs_v))
	sort.Sort(NameSorter(certs_r))
	sort.Sort(NameSorter(certs_err))

	certsSorted = append(certsSorted, certs_v...)
	certsSorted = append(certsSorted, certs_r...)
	certsSorted = append(certsSorted, certs_err...)

	return certsSorted, nil
}

func parseDetails(d string) *Details {
	details := &Details{}
	lines := strings.Split(trim(d), "/")
	for _, line := range lines {
		if strings.Contains(line, "") {
			fields := strings.Split(trim(line), "=")
			switch fields[0] {
			case "name":
				details.Name = fields[1]
			case "CN":
				details.CN = fields[1]
			case "C":
				details.Country = fields[1]
			case "ST":
				details.State = fields[1]
			case "L":
				details.City = fields[1]
			case "O":
				details.Organisation = fields[1]
			case "OU":
				details.OrganisationUnit = fields[1]
			case "emailAddress":
				details.Email = fields[1]
			case "LocalIP":
				details.LocalIP = fields[1]
			default:
				// logs.Debug(fmt.Sprintf("Undefined entry: %s", line))
			}
		}
	}
	return details
}

func CreateCertificate(name string, passphrase string) error {
	path := filepath.Join(state.GlobalCfg.OVConfigPath, "easy-rsa/pki/index.txt")
	pass := false
	existsError := errors.New("Error! There is already a valid or invalid certificate for the name \"" + name + "\"")
	if passphrase != "" {
		pass = true
	}
	certs, err := ReadCerts(path)
	if err != nil {
		//		web.Debug(string(output))
		logs.Error(err)
		//		return err
	}
	// Dump(certs)
	exists := false
	for _, v := range certs {
		if v.Details.CN == name {
			exists = true
		}
	}
	if !pass {
		if !exists {
			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(
					"cd /opt/scripts/ && export DOCKER_COMMAND=3 && "+
						"export CLIENT=%s && export PASS=1 && "+
						"export OVPN_PATH=%s && "+
						"./openvpn-install-v2.sh", name, state.GlobalCfg.OVConfigPath))
			cmd.Dir = state.GlobalCfg.OVConfigPath
			output, err := cmd.CombinedOutput()
			if err != nil {
				logs.Debug(string(output))
				logs.Error(err)
				return err
			}
			return nil
		}
		return existsError
	} else {
		if !exists {
			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(
					"cd /opt/scripts/ && export DOCKER_COMMAND=3 && "+
						"export CLIENT=%s && export PASS=2 && export CL_PASS=%s && "+
						"export OVPN_PATH=%s && "+
						"./openvpn-install-v2.sh", name, passphrase, state.GlobalCfg.OVConfigPath))
			cmd.Dir = state.GlobalCfg.OVConfigPath
			output, err := cmd.CombinedOutput()
			if err != nil {
				logs.Debug(string(output))
				logs.Error(err)
				return err
			}
			return nil
		}
		return existsError
	}
}

func CreateOVPNFile(name string) error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && export DOCKER_COMMAND=4 && "+
				"export CLIENT=%s && "+
				"export OVPN_PATH=%s && "+
				"./openvpn-install-v2.sh", name, state.GlobalCfg.OVConfigPath))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}

func RevokeCertificate(name string) error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && export DOCKER_COMMAND=5 && "+
				"export CLIENT=%s && "+
				"export OVPN_PATH=%s && "+
				"./openvpn-install-v2.sh", name, state.GlobalCfg.OVConfigPath))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}

func UnRevokeCertificate(name string) error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && export DOCKER_COMMAND=6 && "+
				"export CLIENT=%s && "+
				"export OVPN_PATH=%s && "+
				"./openvpn-install-v2.sh", name, state.GlobalCfg.OVConfigPath))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}

func BurnCertificate(CN string, serial string) error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && "+
				"export OVDIR=%s && "+
				"./rmcert.sh %s %s", state.GlobalCfg.OVConfigPath, CN, serial))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}

func RenewCertificate(name string, serial string) error {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf(
			"cd /opt/scripts/ && "+
				"export KEY_NAME=%s &&"+
				"./renew.sh %s %s", name, name, serial))
	cmd.Dir = state.GlobalCfg.OVConfigPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Debug(string(output))
		logs.Error(err)
		return err
	}
	return nil
}
