package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/devlup-labs/sos/internal/pkg/constants"
	"github.com/devlup-labs/sos/internal/pkg/policy"
	"github.com/devlup-labs/sos/internal/pkg/sshcert"
	"github.com/devlup-labs/sos/openpubkey/client"
	"github.com/devlup-labs/sos/openpubkey/pktoken"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
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

func authorizedKeysCommand(
	userArg string,
	typArg string,
	certB64Arg string,
	policyEnforcer policyCheck,
	op client.OpenIdProvider,
) (string, error) {
	cert, err := sshcert.NewFromAuthorizedKey(typArg, certB64Arg)
	if err != nil {
		return "", err
	}

	if pkt, err := cert.VerifySshPktCert(op); err != nil {
		return "", err
	} else if err := policyEnforcer(userArg, pkt); err != nil {
		return "", err
	} else {
		pubkeyBytes := ssh.MarshalAuthorizedKey(cert.SshCert.SignatureKey)

		return "cert-authority " + string(pubkeyBytes), nil
	}
}

func Log(line string) {
	f, err := os.OpenFile(
		"/var/log/sos.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0700,
	)
	if err != nil {
		fmt.Println("Couldn't write to file")
	} else {
		defer f.Close()

		if _, err = f.WriteString(line + "\n"); err != nil {
			fmt.Println("Couldn't write to file")
		}
	}
}

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies OPK tokens",
	Long: `
This command is called by the SSH server as the authorizedKeysCommand:

The following lines are added to /etc/ssh/sshd_config:

AuthorizedKeysCommand /etc/opk/opkssh ver %u %k %t
AuthorizedKeysCommandUser root

The parameters specified in the config map the parameters sent to the function below.
We prepend "Arg" to specify which ones are arguments sent by sshd. They are:

%u The username (requested principal) - userArg
%t The public key type - typArg - in this case a certificate being used as a public key
%k The base64-encoded public key for authentication - certB64Arg - the public key is also a certificate
	`,
	Run: func(cmd *cobra.Command, args []string) {
		Log(strings.Join(os.Args, " "))

		policyEnforcer := simpleFilePolicyEnforcer{
			PolicyFilePath: "/etc/sos/policy.yml",
		}

		if len(os.Args) != 5 {
			fmt.Println("Invalid number of arguments for verify, should be `verifier verify <User (TOKEN u)> <Cert (TOKEN k)> <Key type (TOKEN t)>`")

			os.Exit(1)
		}

		userArg := os.Args[2]
		certB64Arg := os.Args[3]
		typArg := os.Args[4]

		authKey, err := authorizedKeysCommand(
			userArg,
			typArg,
			certB64Arg,
			policyEnforcer.checkPolicy,
			&constants.Op,
		)
		cobra.CheckErr(err)

		fmt.Println(authKey)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
