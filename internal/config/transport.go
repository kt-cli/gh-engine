package config

import (
	"errors"
	"github.com/cli/cli/v2/internal/gh"
)

// NewTransportConfig constructs an in-memory configuration without reading the
// environment, config files, or the operating system credential store.
func NewTransportConfig() gh.Config {
	return &transportConfig{Config: &cfg{fallbackConfig()}}
}

type transportConfig struct{ gh.Config }

func (c *transportConfig) Authentication() gh.AuthConfig { return transportAuth{} }
func (c *transportConfig) Write() error {
	return errors.New("transport configuration cannot be persisted")
}

type transportAuth struct{}

func (transportAuth) HasActiveToken(string) bool                               { return false }
func (transportAuth) ActiveToken(string) (string, string)                      { return "", "transport" }
func (transportAuth) HasEnvToken() bool                                        { return false }
func (transportAuth) TokenFromKeyring(string) (string, error)                  { return "", errHeld }
func (transportAuth) TokenFromKeyringForUser(string, string) (string, error)   { return "", errHeld }
func (transportAuth) ActiveUser(string) (string, error)                        { return "", errHeld }
func (transportAuth) Hosts() []string                                          { return []string{"github.com"} }
func (transportAuth) DefaultHost() (string, string)                            { return "github.com", "transport" }
func (transportAuth) Login(string, string, string, string, bool) (bool, error) { return false, errHeld }
func (transportAuth) SwitchUser(string, string) error                          { return errHeld }
func (transportAuth) Logout(string, string) error                              { return errHeld }
func (transportAuth) UsersForHost(string) []string                             { return nil }
func (transportAuth) TokenForUser(string, string) (string, string, error) {
	return "", "transport", errHeld
}
func (transportAuth) SetActiveToken(string, string) {}
func (transportAuth) SetHosts([]string)             {}
func (transportAuth) SetDefaultHost(string, string) {}

var errHeld = errors.New("credentials are held by the configured transport")

func (transportAuth) APIHostForHost(string) (string, bool) { return "", false }
func (transportAuth) HostForAPIHost(string) (string, bool) { return "", false }
func (transportAuth) ActiveTokenType(string) gh.TokenType  { return gh.TokenType("") }
