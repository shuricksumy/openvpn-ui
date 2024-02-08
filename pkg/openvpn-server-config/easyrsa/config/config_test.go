package config_test

import (
	"io/ioutil"
	"testing"

	easyrsaconfig "github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/easyrsa/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := easyrsaconfig.New()
	assert.Equal(t, c.Auth, "SHA256")
}

func TestTemplateGeneration(t *testing.T) {
	c := easyrsaconfig.New()
	txt, err := ioutil.ReadFile("./templates/easyrsa-vars.tpl")
	assert.Nil(t, err)

	_, err = easyrsaconfig.GetText(string(txt), c)
	assert.Nil(t, err)
}

func TestBrokenTemplate(t *testing.T) {
	c := config.New()

	_, err := config.GetText("{{ ", c)
	assert.NotNil(t, err, "Parser should fail on broken template")
}
