package cmd

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devlup-labs/sos/internal/pkg/constants"
	"github.com/devlup-labs/sos/internal/pkg/sshcert"
	"github.com/devlup-labs/sos/openpubkey/client"
	"github.com/devlup-labs/sos/openpubkey/util"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Opens the login page for email authentication",
	Long:  `Email auth uses the JWT token to then generate the openpubkey for SSH`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--Running Login----")

		principals := []string{}
		alg := jwa.ES256
		signer, err := util.GenKeyPair(alg)
		cobra.CheckErr(err)

		opkClient, err := client.New(
			&constants.Op,
			client.WithSigner(signer, alg),
			client.WithSignGQ(false),
		)
		cobra.CheckErr(err)

		certBytes, seckeySshPem, err := createSSHCert(
			context.Background(), opkClient, principals,
		)
		cobra.CheckErr(err)

		cobra.CheckErr(writeKeysToSSHDir(seckeySshPem, certBytes))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func createSSHCert(
	cxt context.Context,
	client *client.OpkClient,
	principals []string,
) ([]byte, []byte, error) {
	pkt, err := client.Auth(cxt)
	if err != nil {
		return nil, nil, err
	}

	cert, err := sshcert.New(pkt, principals)
	if err != nil {
		return nil, nil, err
	}

	sshSigner, err := ssh.NewSignerFromSigner(client.GetSigner())
	if err != nil {
		return nil, nil, err
	}

	signerMas, err := ssh.NewSignerWithAlgorithms(
		sshSigner.(ssh.AlgorithmSigner), []string{ssh.KeyAlgoECDSA256},
	)
	if err != nil {
		return nil, nil, err
	}

	sshCert, err := cert.SignCert(signerMas)
	if err != nil {
		return nil, nil, err
	}

	certBytes := ssh.MarshalAuthorizedKey(sshCert)

	seckeySsh, err := ssh.MarshalPrivateKey(
		client.GetSigner(), "openpubkey cert",
	)
	if err != nil {
		return nil, nil, err
	}

	seckeySshBytes := pem.EncodeToMemory(seckeySsh)

	return certBytes, seckeySshBytes, nil
}

func writeKeys(
	seckeyPath string,
	pubkeyPath string,
	seckeySshPem []byte,
	certBytes []byte,
) error {
	if err := os.WriteFile(seckeyPath, seckeySshPem, 0600); err != nil {
		return err
	}

	certBytes = append(certBytes, []byte(" "+"openpubkey")...)

	return os.WriteFile(pubkeyPath, certBytes, 0777)
}

func fileExists(fPath string) bool {
	_, err := os.Open(fPath)

	return !errors.Is(err, os.ErrNotExist)
}

func writeKeysToSSHDir(seckeySshPem []byte, certBytes []byte) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshPath := filepath.Join(homePath, ".ssh")

	for _, keyFilename := range []string{"id_ecdsa", "id_dsa"} {
		seckeyPath := filepath.Join(sshPath, keyFilename)
		pubkeyPath := seckeyPath + ".pub"

		if !fileExists(seckeyPath) {
			return writeKeys(seckeyPath, pubkeyPath, seckeySshPem, certBytes)
		} else if !fileExists(pubkeyPath) {
			continue
		} else {
			sshPubkey, err := os.ReadFile(pubkeyPath)
			if err != nil {
				fmt.Println("Failed to read:", pubkeyPath)

				continue
			}

			sshPubkeySplit := strings.Split(string(sshPubkey), " ")
			if len(sshPubkeySplit) != 3 {
				fmt.Println("Failed to parse:", pubkeyPath)

				continue
			}

			if strings.Contains(sshPubkeySplit[2], ("openpubkey")) {
				return writeKeys(
					seckeyPath, pubkeyPath, seckeySshPem, certBytes,
				)
			}
		}
	}

	return fmt.Errorf("no default ssh key file free for openpubkey")
}
