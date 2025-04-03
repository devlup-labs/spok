// Copyright 2024 OpenPubkey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package cosigner

import (
	"github.com/openpubkey/openpubkey/pktoken"
)

type AuthStateStore interface {
	CreateNewAuthSession(pkt *pktoken.PKToken, ruri string, nonce string) (authID string, err error)
	LookupAuthState(authID string) (*AuthState, bool)
	UpdateAuthState(authID string, authState AuthState) error
	CreateAuthcode(authID string) (authcode string, err error)
	RedeemAuthcode(authcode string) (authState AuthState, authID string, err error)
}
