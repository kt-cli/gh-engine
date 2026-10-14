package ktengine

import (
	"errors"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/extensions"
	"io"
)

type noExtensions struct{}

func (noExtensions) List() []extensions.Extension           { return nil }
func (noExtensions) Install(ghrepo.Interface, string) error { return errExtension }
func (noExtensions) InstallLocal(string) error              { return errExtension }
func (noExtensions) Upgrade(string, bool) error             { return errExtension }
func (noExtensions) Remove(string) error                    { return errExtension }
func (noExtensions) Dispatch([]string, io.Reader, io.Writer, io.Writer) (bool, error) {
	return false, errExtension
}
func (noExtensions) Create(string, extensions.ExtTemplateType) error { return errExtension }
func (noExtensions) EnableDryRunMode()                               {}
func (noExtensions) UpdateDir(string) string                         { return "" }

var errExtension = errors.New("extensions require an independent execution authorization")
