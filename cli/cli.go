package main

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// "github.com/devlup-labs/sos/internal/pkg/policy"
	"github.com/devlup-labs/sos/internal/pkg/sshcert"
	"github.com/devlup-labs/sos/openpubkey/client"
	"github.com/devlup-labs/sos/openpubkey/client/providers"
	"github.com/devlup-labs/sos/openpubkey/util"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"golang.org/x/crypto/ssh"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) < 2 {
		fmt.Printf(
			"Secure Openpubkey Shell Client: Command choices are: configure, login",
		)

		return
	}

	command := os.Args[1]

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	scopes := []string{"openid profile email"}
	redirURIPort := os.Getenv("REDIRECT_URI_PORT")
	callbackPath := os.Getenv("CALLBACK_PATH")
	redirectURI := fmt.Sprintf(
		"http://localhost:%v%v", redirURIPort, callbackPath,
	)

	op := providers.GoogleOp{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		RedirURIPort: redirURIPort,
		CallbackPath: callbackPath,
		RedirectURI:  redirectURI,
	}

	switch command {
	case "configure":
		{
			emailArgs := os.Args[2]
			userArgs := os.Args[3]

			principal := strings.Split(userArgs, "@")[0]

			cmd := exec.Command(
				"scp",
				"scripts/configure-opk-server.sh",
				userArgs+":/root/configure-opk-server.sh",
			)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			cmd2 := exec.Command(
				"scp",
				"verifier/verifier",
				userArgs+":/root/verifier",
			)
			err = cmd2.Run()
			if err != nil {
				log.Fatal(err)
			}

			cmd3 := exec.Command(
				"ssh",
				userArgs,
				"chmod",
				"+x",
				"/root/configure-opk-server.sh",
			)
			err = cmd3.Run()
			if err != nil {
				log.Fatal(err)
			}

			cmd4 := exec.Command(
				"ssh",
				userArgs,
				"/root/configure-opk-server.sh",
				emailArgs,
				principal,
			)
			err = cmd4.Run()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Configured SOS server for:", emailArgs)
		}
	case "login":
		{
			if len(os.Args) != 2 {
				fmt.Println("Invalid number of arguments for login, should be `sos login`")

				os.Exit(1)
			}

			// If principals is empty the server does not enforce any principal.
			// The OPK verifier should use policy to make this decision.
			principals := []string{}

			alg := jwa.ES256
			signer, err := util.GenKeyPair(alg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			opkClient, err := client.New(
				&op,
				client.WithSigner(signer, alg),
				client.WithSignGQ(false),
			)
			if err != nil {
				fmt.Println(err)
			}

			certBytes, seckeySshPem, err := createSSHCert(
				context.Background(), opkClient, principals,
			)
			if err != nil {
				fmt.Println(err)

				os.Exit(1)
			}

			err = writeKeysToSSHDir(seckeySshPem, certBytes)
			if err != nil {
				fmt.Println(err)

				os.Exit(1)
			}
			os.Exit(0)
		}
	default:
		fmt.Println("Error! Unrecognized command:", command)
	}
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
