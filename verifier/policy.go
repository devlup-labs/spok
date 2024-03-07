package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/devlup-labs/sos/internal/pkg/policy"
	"github.com/devlup-labs/sos/openpubkey/pktoken"
	"golang.org/x/exp/slices"
)

type simpleFilePolicyEnforcer struct {
	PolicyFilePath string
}

func (p *simpleFilePolicyEnforcer) readPolicyFile() (*policy.Policy, error) {
	info, err := os.Stat(p.PolicyFilePath)
	if err != nil {
		return nil, err
	}

	mode := info.Mode()

	if mode.Perm() != fs.FileMode(0600) {
		return nil, fmt.Errorf(
			"policy file has insecure permissions, expected (0600), got (%o)",
			mode.Perm(),
		)
	}

	allowedPolicy := new(policy.Policy)
	err = allowedPolicy.Unmarshal("/etc/sos/policy.yml")
	if err != nil {
		return nil, err
	}

	return allowedPolicy, err
}

func (p *simpleFilePolicyEnforcer) checkPolicy(
	principalDesired string, pkt *pktoken.PKToken,
) error {
	allowedPolicy, err := p.readPolicyFile()
	if err != nil {
		return err
	}

	var claims struct {
		Email string `json:"email"`
	}

	if err := json.Unmarshal(pkt.Payload, &claims); err != nil {
		return err
	}

	for _, u := range allowedPolicy.User {
		if u.Email == claims.Email {
			if slices.Contains(u.Principals, principalDesired) {
				return nil
			} else {
				return fmt.Errorf(
					"no policy to allow %s to assume %s, check policy config",
					claims.Email,
					principalDesired,
				)
			}
		}
	}

	return fmt.Errorf("no email or policy found")
}

type policyCheck func(userDesired string, pkt *pktoken.PKToken) error
