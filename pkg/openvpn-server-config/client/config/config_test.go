package config_test

import (
	"io/ioutil"
	"testing"

	"github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/client/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := config.New()
	assert.Equal(t, c.Auth, "SHA256")
}

func TestTemplateGeneration(t *testing.T) {
	c := config.New()
	txt, err := ioutil.ReadFile("./templates/openvpn-client-config.tpl")
	assert.Nil(t, err)

	_, err = config.GetText(string(txt), c)
	assert.Nil(t, err)
}

func TestBrokenTemplate(t *testing.T) {
	c := config.New()

	_, err := config.GetText("{{ ", c)
	assert.NotNil(t, err, "Parser should fail on broken template")
}
