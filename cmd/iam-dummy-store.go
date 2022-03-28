// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"sync"

	"github.com/minio/minio/internal/auth"
)

type iamDummyStore struct {
	sync.RWMutex
	*IAMCache
	usersSysType UsersSysType
}

func newIAMDummyStore(usersSysType UsersSysType) *iamDummyStore {
	return &iamDummyStore{
		IAMCache:     newIamCache(),
		usersSysType: usersSysType,
	}
}

func (ids *iamDummyStore) rlock() *IAMCache {
	ids.RLock()
	return ids.IAMCache
}

func (ids *iamDummyStore) runlock() {
	ids.RUnlock()
}

func (ids *iamDummyStore) lock() *IAMCache {
	ids.Lock()
	return ids.IAMCache
}

func (ids *iamDummyStore) unlock() {
	ids.Unlock()
}

func (ids *iamDummyStore) getUsersSysType() UsersSysType {
	return ids.usersSysType
}

func (ids *iamDummyStore) migrateBackendFormat(context.Context) error {
	return nil
}

func (ids *iamDummyStore) loadPolicyDoc(ctx context.Context, policy string, m map[string]PolicyDoc) error {
	v, ok := ids.PolicyDocsMap[policy]
	if !ok {
		return errNoSuchPolicy
	}
	m[policy] = v
	return nil
}

func (ids *iamDummyStore) loadPolicyDocs(ctx context.Context, m map[string]PolicyDoc) error {
	for k, v := range ids.PolicyDocsMap {
		m[k] = v
	}
	return nil
}

func (ids *iamDummyStore) loadUser(ctx context.Context, user string, userType IAMUserType, m map[string]auth.Credentials) error {
	u, ok := ids.UsersMap[user]
	if !ok {
		return errNoSuchUser
	}
	ids.UsersMap[user] = u
	return nil
}

func (ids *iamDummyStore) loadUsers(ctx context.Context, userType IAMUserType, m map[string]auth.Credentials) error {
	for k, v := range ids.UsersMap {
		m[k] = v
	}
	return nil
}

func (ids *iamDummyStore) loadGroup(ctx context.Context, group string, m map[string]GroupInfo) error {
	g, ok := ids.GroupsMap[group]
	if !ok {
		return errNoSuchGroup
	}
	m[group] = g
	return nil
}

func (ids *iamDummyStore) loadGroups(ctx context.Context, m map[string]GroupInfo) error {
	for k, v := range ids.GroupsMap {
		m[k] = v
	}
	return nil
}

func (ids *iamDummyStore) loadMappedPolicy(ctx context.Context, name string, userType IAMUserType, isGroup bool, m map[string]MappedPolicy) error {
	if isGroup {
		g, ok := ids.GroupPolicyMap[name]
		if !ok {
			return errNoSuchPolicy
		}
		m[name] = g
	} else {
		u, ok := ids.UserPolicyMap[name]
		if !ok {
			return errNoSuchPolicy
		}
		m[name] = u
	}
	return nil
}

func (ids *iamDummyStore) loadMappedPolicies(ctx context.Context, userType IAMUserType, isGroup bool, m map[string]MappedPolicy) error {
	if !isGroup {
		for k, v := range ids.UserPolicyMap {
			m[k] = v
		}
	} else {
		for k, v := range ids.GroupPolicyMap {
			m[k] = v
		}
	}
	return nil
}

func (ids *iamDummyStore) saveIAMConfig(ctx context.Context, item interface{}, path string, opts ...options) error {
	return nil
}

func (ids *iamDummyStore) loadIAMConfig(ctx context.Context, item interface{}, path string) error {
	return nil
}

func (ids *iamDummyStore) deleteIAMConfig(ctx context.Context, path string) error {
	return nil
}

func (ids *iamDummyStore) savePolicyDoc(ctx context.Context, policyName string, p PolicyDoc) error {
	return nil
}

func (ids *iamDummyStore) saveMappedPolicy(ctx context.Context, name string, userType IAMUserType, isGroup bool, mp MappedPolicy, opts ...options) error {
	return nil
}

func (ids *iamDummyStore) saveUserIdentity(ctx context.Context, name string, userType IAMUserType, u UserIdentity, opts ...options) error {
	return nil
}

func (ids *iamDummyStore) saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error {
	return nil
}

func (ids *iamDummyStore) deletePolicyDoc(ctx context.Context, policyName string) error {
	return nil
}

func (ids *iamDummyStore) deleteMappedPolicy(ctx context.Context, name string, userType IAMUserType, isGroup bool) error {
	return nil
}

func (ids *iamDummyStore) deleteUserIdentity(ctx context.Context, name string, userType IAMUserType) error {
	return nil
}

func (ids *iamDummyStore) deleteGroupInfo(ctx context.Context, name string) error {
	return nil
}
