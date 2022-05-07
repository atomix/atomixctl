// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package version

import sdkversion "github.com/atomix/sdk/pkg/version"

var (
	version   string
	commit    string
	isRelease bool
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func SDKVersion() string {
	return sdkversion.Version()
}

func IsRelease() bool {
	return isRelease
}
