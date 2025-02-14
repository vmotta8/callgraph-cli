package detectors

import (
	"os"
)

func DetectLanguage() string {
	if hasGoMod() {
		return "go"
	}
	if hasPackageJSON() {
		return "javascript"
	}
	if hasRequirementsTxt() {
		return "python"
	}
	if hasCargoToml() {
		return "rust"
	}
	return "unknown"
}

func hasGoMod() bool {
	_, err := os.Stat("go.mod")
	return !os.IsNotExist(err)
}

func hasPackageJSON() bool {
	_, err := os.Stat("package.json")
	return !os.IsNotExist(err)
}

func hasRequirementsTxt() bool {
	_, err := os.Stat("requirements.txt")
	return !os.IsNotExist(err)
}

func hasCargoToml() bool {
	_, err := os.Stat("Cargo.toml")
	return !os.IsNotExist(err)
}
