// Code generated by "stringer -type=STSErrorCode -trimprefix=Err sts-errors.go"; DO NOT EDIT.

package cmd

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrSTSNone-0]
	_ = x[ErrSTSAccessDenied-1]
	_ = x[ErrSTSMissingParameter-2]
	_ = x[ErrSTSInvalidParameterValue-3]
	_ = x[ErrSTSWebIdentityExpiredToken-4]
	_ = x[ErrSTSClientGrantsExpiredToken-5]
	_ = x[ErrSTSInvalidClientGrantsToken-6]
	_ = x[ErrSTSMalformedPolicyDocument-7]
	_ = x[ErrSTSInsecureConnection-8]
	_ = x[ErrSTSInvalidClientCertificate-9]
	_ = x[ErrSTSNotInitialized-10]
	_ = x[ErrSTSIAMNotInitialized-11]
	_ = x[ErrSTSUpstreamError-12]
	_ = x[ErrSTSInternalError-13]
}

const _STSErrorCode_name = "STSNoneSTSAccessDeniedSTSMissingParameterSTSInvalidParameterValueSTSWebIdentityExpiredTokenSTSClientGrantsExpiredTokenSTSInvalidClientGrantsTokenSTSMalformedPolicyDocumentSTSInsecureConnectionSTSInvalidClientCertificateSTSNotInitializedSTSIAMNotInitializedSTSUpstreamErrorSTSInternalError"

var _STSErrorCode_index = [...]uint16{0, 7, 22, 41, 65, 91, 118, 145, 171, 192, 219, 236, 256, 272, 288}

func (i STSErrorCode) String() string {
	if i < 0 || i >= STSErrorCode(len(_STSErrorCode_index)-1) {
		return "STSErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _STSErrorCode_name[_STSErrorCode_index[i]:_STSErrorCode_index[i+1]]
}
