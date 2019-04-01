/*
 * MinIO Cloud Storage, (C) 2019 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/minio/minio/pkg/env"
	config "github.com/minio/minio/pkg/server-config"
	ldap "gopkg.in/ldap.v3"
)

const (
	defaultLDAPExpiry = time.Hour * 1
)

// Config contains AD/LDAP server connectivity information.
type Config struct {
	IsEnabled bool `json:"enabled"`

	// E.g. "ldap.minio.io:636"
	ServerAddr string `json:"serverAddr"`

	// STS credentials expiry duration
	STSExpiryDuration string        `json:"stsExpiryDuration"`
	stsExpiryDuration time.Duration // contains converted value

	RootCAs *x509.CertPool `json:"-"`

	// Format string for usernames
	UsernameFormat string `json:"usernameFormat"`

	GroupSearchBaseDN  string `json:"groupSearchBaseDN"`
	GroupSearchFilter  string `json:"groupSearchFilter"`
	GroupNameAttribute string `json:"groupNameAttribute"`
}

// LDAP keys and envs.
const (
	ServerAddr         = "server_addr"
	STSExpiry          = "sts_expiry"
	UsernameFormat     = "username_format"
	GroupSearchFilter  = "group_search_filter"
	GroupNameAttribute = "group_name_attribute"
	GroupSearchBaseDN  = "group_search_base_dn"

	EnvServerAddr         = "MINIO_IDENTITY_LDAP_SERVER_ADDR"
	EnvSTSExpiry          = "MINIO_IDENTITY_LDAP_STS_EXPIRY"
	EnvUsernameFormat     = "MINIO_IDENTITY_LDAP_USERNAME_FORMAT"
	EnvGroupSearchFilter  = "MINIO_IDENTITY_LDAP_GROUP_SEARCH_FILTER"
	EnvGroupNameAttribute = "MINIO_IDENTITY_LDAP_GROUP_NAME_ATTRIBUTE"
	EnvGroupSearchBaseDN  = "MINIO_IDENTITY_LDAP_GROUP_SEARCH_BASE_DN"
)

// Connect connect to ldap server.
func (l *Config) Connect() (ldapConn *ldap.Conn, err error) {
	if l == nil {
		// Happens when LDAP is not configured.
		return
	}
	return ldap.DialTLS("tcp", l.ServerAddr, &tls.Config{RootCAs: l.RootCAs})
}

// GetExpiryDuration - return parsed expiry duration.
func (l Config) GetExpiryDuration() time.Duration {
	return l.stsExpiryDuration
}

// NewConfig - initializes LDAP config, overrides config, if any ENV values are set.
func NewConfig(kvs config.KVS, rootCAs *x509.CertPool) (l Config, err error) {
	l.IsEnabled = kvs.Get(config.State) == config.StateEnabled
	if !l.IsEnabled {
		return l, nil
	}
	l.RootCAs = rootCAs
	ldapServer := env.Get(EnvServerAddr, kvs.Get(ServerAddr))
	if ldapServer == "" {
		return l, errors.New("ldap server cannot initialize with empty LDAP server")
	}

	l.ServerAddr = ldapServer
	l.stsExpiryDuration = defaultLDAPExpiry
	if v := env.Get(EnvSTSExpiry, kvs.Get(STSExpiry)); v != "" {
		expDur, err := time.ParseDuration(v)
		if err != nil {
			return l, errors.New("LDAP expiry time err:" + err.Error())
		}
		if expDur <= 0 {
			return l, errors.New("LDAP expiry time has to be positive")
		}
		l.STSExpiryDuration = v
		l.stsExpiryDuration = expDur
	}

	if v := env.Get(EnvUsernameFormat, kvs.Get(UsernameFormat)); v != "" {
		subs := NewSubstituter("username", "test")
		if _, err := subs.Substitute(v); err != nil {
			return l, errors.New("Only username may be substituted in the username format")
		}
		l.UsernameFormat = v
	}

	grpSearchFilter := env.Get(EnvGroupSearchFilter, kvs.Get(GroupSearchFilter))
	grpSearchNameAttr := env.Get(EnvGroupNameAttribute, kvs.Get(GroupNameAttribute))
	grpSearchBaseDN := env.Get(EnvGroupSearchBaseDN, kvs.Get(GroupSearchBaseDN))

	// Either all group params must be set or none must be set.
	allNotSet := grpSearchFilter == "" && grpSearchNameAttr == "" && grpSearchBaseDN == ""
	allSet := grpSearchFilter != "" && grpSearchNameAttr != "" && grpSearchBaseDN != ""
	if !allNotSet && !allSet {
		return l, errors.New("All group related parameters must be set")
	}

	if allSet {
		subs := NewSubstituter("username", "test", "usernamedn", "test2")
		if _, err := subs.Substitute(grpSearchFilter); err != nil {
			return l, errors.New("Only username and usernamedn may be substituted in the group search filter string")
		}
		l.GroupSearchFilter = grpSearchFilter
		l.GroupNameAttribute = grpSearchNameAttr
		subs = NewSubstituter("username", "test", "usernamedn", "test2")
		if _, err := subs.Substitute(grpSearchBaseDN); err != nil {
			return l, errors.New("Only username and usernamedn may be substituted in the base DN string")
		}
		l.GroupSearchBaseDN = grpSearchBaseDN
	}
	return
}

// Substituter - This type is to allow restricted runtime
// substitutions of variables in LDAP configuration items during
// runtime.
type Substituter struct {
	vals map[string]string
}

// NewSubstituter - sets up the substituter for usage, for e.g.:
//
//  subber := NewSubstituter("username", "john")
func NewSubstituter(v ...string) Substituter {
	if len(v)%2 != 0 {
		log.Fatal("Need an even number of arguments")
	}
	vals := make(map[string]string)
	for i := 0; i < len(v); i += 2 {
		vals[v[i]] = v[i+1]
	}
	return Substituter{vals: vals}
}

// substitute - performs substitution on the given string `t`. Returns
// an error if there are any variables in the input that do not have
// values in the substituter. E.g.:
//
// subber.substitute("uid=${username},cn=users,dc=example,dc=com")
//
// returns "uid=john,cn=users,dc=example,dc=com"
//
// whereas:
//
// subber.substitute("uid=${usernamedn}")
//
// returns an error.
func (s *Substituter) Substitute(t string) (string, error) {
	for k, v := range s.vals {
		re := regexp.MustCompile(fmt.Sprintf(`\$\{%s\}`, k))
		t = re.ReplaceAllLiteralString(t, v)
	}
	// Check if all requested substitutions have been made.
	re := regexp.MustCompile(`\$\{.*\}`)
	if re.MatchString(t) {
		return "", errors.New("unsupported substitution requested")
	}
	return t, nil
}
