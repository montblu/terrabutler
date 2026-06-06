package requirements

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRequirementsValid(t *testing.T) {

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	// Pre-populate the memory filesystem for the test
	// In this case we need the settings.yaml
	configPath := "ROOT/configs/settings.yml"
	_ = afero.WriteFile(fs, configPath, []byte(``), 0644)

	t.Setenv("TERRABUTLER_ROOT", "ROOT")
	t.Setenv("TERRABUTLER_ENABLE", "true")

	assert.NoError(t, Check_requirement(fs), "The test was not supposed to fail.")

}

func TestRequirementsInvalidROOT(t *testing.T) {

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	// Pre-populate the memory filesystem for the test
	// In this case we need the settings.yaml
	configPath := "ROOT/configs/settings.yml"
	_ = afero.WriteFile(fs, configPath, []byte(``), 0644)

	t.Setenv("TERRABUTLER_ENABLE", "true")

	assert.Error(t, Check_requirement(fs), "Failed, it was accepted with TERRABUTLER_ROOT empty")

	t.Setenv("TERRABUTLER_ROOT", "NOROOT")

	assert.Error(t, Check_requirement(fs), "Failed, it was accepted with TERRABUTLER_ROOT set with the wrong directory")

}

func TestRequirementsInvalidENABLE(t *testing.T) {

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	// Pre-populate the memory filesystem for the test
	// In this case we need the settings.yaml
	configPath := "ROOT/configs/settings.yml"
	_ = afero.WriteFile(fs, configPath, []byte(``), 0644)
	t.Setenv("TERRABUTLER_ROOT", "ROOT")

	assert.Error(t, Check_requirement(fs), "Failed, it was accepted with TERRABUTLER_ENABLE empty")

	t.Setenv("TERRABUTLER_ENABLE", "false")

	assert.Error(t, Check_requirement(fs), "Failed, it was accepted with TERRABUTLER_ENABLE set false")

	t.Setenv("TERRABUTLER_ENABLE", "random")

	assert.Error(t, Check_requirement(fs), "Failed,  it was accepted with TERRABUTLER_ENABLE set with a invalid input")

}

func TestRequirementsInvalidConfigPath(t *testing.T) {

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	t.Setenv("TERRABUTLER_ROOT", "ROOT")
	t.Setenv("TERRABUTLER_ENABLE", "true")

	// No file exists
	assert.Error(t, Check_requirement(fs), "Failed, it accepted a file which doesn't exists")

	configPath := "NOROOT/configs/settings.yml"
	_ = afero.WriteFile(fs, configPath, []byte(``), 0644)

	// A file ina different directory
	assert.Error(t, Check_requirement(fs), "Failed, it accepted a file with a different in a different directory")

	configPath = "ROOT/configs/settings.yaml"
	_ = afero.WriteFile(fs, configPath, []byte(``), 0644)

	// A file with a different name
	assert.Error(t, Check_requirement(fs), "Failed, it accepted a file with a different name")
}
