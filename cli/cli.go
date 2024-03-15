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

	clientID := "992028499768-ce9juclb3vvckh23r83fjkmvf1lvjq18.apps.googleusercontent.com"
	// The clientSecret was intentionally checked in. It holds no power and is used for development. Do not report as a security issue
	clientSecret := "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" // Google requires a ClientSecret even if this a public OIDC App
	scopes       := []string{"openid profile email"}
	redirURIPort := "3000"
	callbackPath := "/login-callback"
	redirectURI  := fmt.Sprintf("http://localhost:%v%v", redirURIPort, callbackPath)

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

			privateKeyAuth := false
			privateKeyPath := ""

			if len(os.Args) == 6 {
				if os.Args[4] == "-i" {
					privateKeyAuth = true
					privateKeyPath = os.Args[5]
				}
			}

			var scpCommandScript []string
			var scpCommandVerifier []string
			var sshCommandChmod []string
			var sshCommandConfigure []string

			if privateKeyAuth && privateKeyPath != "" {
				scpCommandScript = []string{
					"scp",
					"-i",
					privateKeyPath,
					"scripts/configure-sos-server.sh",
					userArgs + ":/root/configure-sos-server.sh",
				}
				scpCommandVerifier = []string{
					"scp",
					"-i",
					privateKeyPath,
					"verifier/verifier",
					userArgs + ":/root/verifier",
				}
				sshCommandChmod = []string{
					"ssh",
					"-i",
					privateKeyPath,
					userArgs,
					"chmod",
					"+x",
					"/root/configure-sos-server.sh",
				}
				sshCommandConfigure = []string{
					"ssh",
					"-i",
					privateKeyPath,
					userArgs,
					"/root/configure-sos-server.sh",
					emailArgs,
					principal,
				}
			} else {
				scpCommandScript = []string{
					"scp",
					"scripts/configure-sos-server.sh",
					userArgs + ":/root/configure-sos-server.sh",
				}
				scpCommandVerifier = []string{
					"scp",
					"verifier/verifier",
					userArgs + ":/root/verifier",
				}
				sshCommandChmod = []string{
					"ssh",
					userArgs,
					"chmod",
					"+x",
					"/root/configure-sos-server.sh",
				}
				sshCommandConfigure = []string{
					"ssh",
					userArgs,
					"/root/configure-sos-server.sh",
					emailArgs,
					principal,
				}
			}

			scpCmdScript := exec.Command(
				scpCommandScript[0], scpCommandScript[1:]...,
			)

			fmt.Println("Copying configuration script to server...")

			err := scpCmdScript.Run()
			if err != nil {
				log.Fatal(err)
			}

			scpCmdVerifier := exec.Command(
				scpCommandVerifier[0], scpCommandVerifier[1:]...,
			)

			fmt.Println("Copying verifier to server...")

			err = scpCmdVerifier.Run()
			if err != nil {
				log.Fatal(err)
			}

			sshCmdChmod := exec.Command(
				sshCommandChmod[0], sshCommandChmod[1:]...,
			)

			fmt.Println("Making configuration script executable...")

			err = sshCmdChmod.Run()
			if err != nil {
				log.Fatal(err)
			}

			sshCmdConfigure := exec.Command(
				sshCommandConfigure[0], sshCommandConfigure[1:]...,
			)

			fmt.Println("Configuring SOS server...")

			err = sshCmdConfigure.Run()
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
