package main

import (
	"strings"
	"testing"

	"github.com/devlup-labs/sos/internal/pkg/sshcert"
	"github.com/devlup-labs/sos/openpubkey/client/providers"
	"github.com/devlup-labs/sos/openpubkey/pktoken"
	"github.com/devlup-labs/sos/openpubkey/pktoken/mocks"
	"github.com/devlup-labs/sos/openpubkey/util"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"golang.org/x/crypto/ssh"
)

func AllowAllPolicyEnforcer(userDesired string, pkt *pktoken.PKToken) error {
	return nil
}

func TestAuthorizedKeysCommand(t *testing.T) {
	op, err := providers.NewMockOpenIdProvider()
	if err != nil {
		t.Fatal(err)
	}

	principals := []string{"guest", "dev"}
	alg := jwa.ES256

	signer, err := util.GenKeyPair(alg)
	if err != nil {
		t.Fatal(err)
	}

	email := "arthur.aardvark@example.com"

	pkt, err := mocks.GenerateMockPKTokenWithEmail(signer, alg, email)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := sshcert.New(pkt, principals)
	if err != nil {
		t.Error(err)
	}

	sshSigner, err := ssh.NewSignerFromSigner(signer)
	if err != nil {
		t.Error(err)
	}

	signerMas, err := ssh.NewSignerWithAlgorithms(
		sshSigner.(ssh.AlgorithmSigner),
		[]string{ssh.KeyAlgoECDSA256},
	)
	if err != nil {
		t.Error(err)
	}

	sshCert, err := cert.SignCert(signerMas)
	if err != nil {
		t.Error(err)
	}

	certTypeAndCertB64 := ssh.MarshalAuthorizedKey(sshCert)
	typeArg := strings.Split(string(certTypeAndCertB64), " ")[0]
	certB64Arg := strings.Split(string(certTypeAndCertB64), " ")[1]

	userArg := "user"
	pubkeyList, err := authorizedKeysCommand(
		userArg, typeArg, certB64Arg, AllowAllPolicyEnforcer, op,
	)
	if err != nil {
		t.Error(err)
	}

	expectedPubkeyList := "cert-authority ecdsa-sha2-nistp256"
	if !strings.Contains(pubkeyList, expectedPubkeyList) {
		t.Error(err)
	}
}
