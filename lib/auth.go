package lib

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/shuricksumy/openvpn-ui/models"
	"github.com/go-ldap/ldap/v3"
	"gopkg.in/hlandau/passlib.v1"
)

var authError error

func init() {
	authError = errors.New("invalid login or password")
}

func Authenticate(login string, password string, authType string) (*models.User, error) {
	logs.Info("auth type: ", authType)
	if authType == "ldap" {
		return authenticateLdap(login, password)
	} else {
		return authenticateSimple(login, password)
	}
}

func authenticateSimple(login string, password string) (*models.User, error) {
	user := &models.User{Login: login}
	err := user.Read("Login")
	if err != nil {
		logs.Error(err)
		return nil, authError
	}
	if user.Id < 1 {
		logs.Error(err)
		return nil, authError
	}
	if _, err := passlib.Verify(password, user.Password); err != nil {
		logs.Error(err)
		return nil, authError
	}
	return user, nil
}

func authenticateLdap(login string, password string) (*models.User, error) {
	address, _ := web.AppConfig.String("LdapAddress")
	var connection *ldap.Conn
	var err error
	ldapTransport, _ := web.AppConfig.String("LdapTransport")
	skipVerify, err := web.AppConfig.Bool("LdapInsecureSkipVerify")
	if err != nil {
		logs.Error("LDAP Dial:", err)
		return nil, authError
	}

	if ldapTransport == "tls" {
		connection, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: skipVerify})
	} else {
		connection, err = ldap.Dial("tcp", address)
	}

	if err != nil {
		logs.Error("LDAP Dial:", err)
		return nil, authError
	}

	if ldapTransport == "starttls" {
		err = connection.StartTLS(&tls.Config{InsecureSkipVerify: skipVerify})
		if err != nil {
			logs.Error("LDAP Start TLS:", err)
			return nil, authError
		}
	}

	defer connection.Close()

	bindDn, _ := web.AppConfig.String("LdapBindDn")

	err = connection.Bind(fmt.Sprintf(bindDn, login), password)
	if err != nil {
		logs.Error("LDAP Bind:", err)
		return nil, authError
	}

	user := &models.User{Login: login}
	err = user.Read("Login")
	if err == orm.ErrNoRows {
		err = user.Insert()
	}
	if err != nil {
		logs.Error(err)
		return nil, authError
	}

	return user, nil
}
