// SPDX-License-Identifier: Apache-2.0

package scan

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/complytime/complytime/cmd/openscap-plugin/config"
	"github.com/complytime/complytime/cmd/openscap-plugin/oscap"
	"github.com/complytime/complytime/cmd/openscap-plugin/xccdf"
)

func isXMLFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				return true, nil
			}
			return false, fmt.Errorf("invalid XML file %s: %w", filePath, err)
		}
	}
}

func validateOpenSCAPFiles(cfg *config.Config) (map[string]string, error) {
	if _, err := isXMLFile(cfg.Files.Datastream); err != nil {
		return nil, err
	}

	if _, err := os.Stat(cfg.Files.Policy); err != nil {
		return nil, err
	}

	if _, err := isXMLFile(cfg.Files.Policy); err != nil {
		return nil, err
	}

	return map[string]string{
		"datastream": cfg.Files.Datastream,
		"policy":     cfg.Files.Policy,
		"results":    cfg.Files.Results,
		"arf":        cfg.Files.ARF,
	}, nil
}

func ScanSystem(cfg *config.Config, profile string) ([]byte, error) {
	openscapFiles, err := validateOpenSCAPFiles(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid openscap files: %w", err)
	}

	tailoringProfile := fmt.Sprintf("%s_%s", profile, xccdf.XCCDFTailoringSuffix)
	// In the future, we can add an integrity check to confirm if the expected tailoring profile
	// id exists in the tailoring file. It is not a common case but a guardrail to prevent manual
	// manipulation of the tailoring file would be good.

	output, err := oscap.OscapScan(openscapFiles, tailoringProfile)
	if err != nil {
		return output, fmt.Errorf("failed during scan: %w", err)
	}

	return output, nil
}
