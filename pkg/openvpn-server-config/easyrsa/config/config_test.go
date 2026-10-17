package config_test

import (
	"os"
	"testing"

	easyrsaconfig "github.com/shuricksumy/openvpn-ui/pkg/openvpn-server-config/easyrsa/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := easyrsaconfig.New()
	assert.Equal(t, 2048, c.EasyRSAKeySize)
}

func TestTemplateGeneration(t *testing.T) {
	c := easyrsaconfig.New()
	txt, err := os.ReadFile("./templates/easyrsa-vars.tpl")
	assert.Nil(t, err)

	_, err = easyrsaconfig.GetText(string(txt), c)
	assert.Nil(t, err)
}

func TestBrokenTemplate(t *testing.T) {
	c := easyrsaconfig.New()

	_, err := easyrsaconfig.GetText("{{ ", c)
	assert.NotNil(t, err, "Parser should fail on broken template")
}
