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

// Code generated by "stringer -type=APIErrorCode -trimprefix=Err api-errors.go"; DO NOT EDIT.

package cmd

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrNone-0]
	_ = x[ErrAccessDenied-1]
	_ = x[ErrBadDigest-2]
	_ = x[ErrEntityTooSmall-3]
	_ = x[ErrEntityTooLarge-4]
	_ = x[ErrPolicyTooLarge-5]
	_ = x[ErrIncompleteBody-6]
	_ = x[ErrInternalError-7]
	_ = x[ErrInvalidAccessKeyID-8]
	_ = x[ErrInvalidBucketName-9]
	_ = x[ErrInvalidDigest-10]
	_ = x[ErrInvalidRange-11]
	_ = x[ErrInvalidRangePartNumber-12]
	_ = x[ErrInvalidCopyPartRange-13]
	_ = x[ErrInvalidCopyPartRangeSource-14]
	_ = x[ErrInvalidMaxKeys-15]
	_ = x[ErrInvalidEncodingMethod-16]
	_ = x[ErrInvalidMaxUploads-17]
	_ = x[ErrInvalidMaxParts-18]
	_ = x[ErrInvalidPartNumberMarker-19]
	_ = x[ErrInvalidPartNumber-20]
	_ = x[ErrInvalidRequestBody-21]
	_ = x[ErrInvalidCopySource-22]
	_ = x[ErrInvalidMetadataDirective-23]
	_ = x[ErrInvalidCopyDest-24]
	_ = x[ErrInvalidPolicyDocument-25]
	_ = x[ErrInvalidObjectState-26]
	_ = x[ErrMalformedXML-27]
	_ = x[ErrMissingContentLength-28]
	_ = x[ErrMissingContentMD5-29]
	_ = x[ErrMissingRequestBodyError-30]
	_ = x[ErrMissingSecurityHeader-31]
	_ = x[ErrNoSuchBucket-32]
	_ = x[ErrNoSuchBucketPolicy-33]
	_ = x[ErrNoSuchBucketLifecycle-34]
	_ = x[ErrNoSuchLifecycleConfiguration-35]
	_ = x[ErrNoSuchBucketSSEConfig-36]
	_ = x[ErrNoSuchCORSConfiguration-37]
	_ = x[ErrNoSuchWebsiteConfiguration-38]
	_ = x[ErrReplicationConfigurationNotFoundError-39]
	_ = x[ErrRemoteDestinationNotFoundError-40]
	_ = x[ErrReplicationDestinationMissingLock-41]
	_ = x[ErrRemoteTargetNotFoundError-42]
	_ = x[ErrReplicationRemoteConnectionError-43]
	_ = x[ErrBucketRemoteIdenticalToSource-44]
	_ = x[ErrBucketRemoteAlreadyExists-45]
	_ = x[ErrBucketRemoteLabelInUse-46]
	_ = x[ErrBucketRemoteArnTypeInvalid-47]
	_ = x[ErrBucketRemoteArnInvalid-48]
	_ = x[ErrBucketRemoteRemoveDisallowed-49]
	_ = x[ErrRemoteTargetNotVersionedError-50]
	_ = x[ErrReplicationSourceNotVersionedError-51]
	_ = x[ErrReplicationNeedsVersioningError-52]
	_ = x[ErrReplicationBucketNeedsVersioningError-53]
	_ = x[ErrObjectRestoreAlreadyInProgress-54]
	_ = x[ErrNoSuchKey-55]
	_ = x[ErrNoSuchUpload-56]
	_ = x[ErrInvalidVersionID-57]
	_ = x[ErrNoSuchVersion-58]
	_ = x[ErrNotImplemented-59]
	_ = x[ErrPreconditionFailed-60]
	_ = x[ErrRequestTimeTooSkewed-61]
	_ = x[ErrSignatureDoesNotMatch-62]
	_ = x[ErrMethodNotAllowed-63]
	_ = x[ErrInvalidPart-64]
	_ = x[ErrInvalidPartOrder-65]
	_ = x[ErrAuthorizationHeaderMalformed-66]
	_ = x[ErrMalformedPOSTRequest-67]
	_ = x[ErrPOSTFileRequired-68]
	_ = x[ErrSignatureVersionNotSupported-69]
	_ = x[ErrBucketNotEmpty-70]
	_ = x[ErrAllAccessDisabled-71]
	_ = x[ErrMalformedPolicy-72]
	_ = x[ErrMissingFields-73]
	_ = x[ErrMissingCredTag-74]
	_ = x[ErrCredMalformed-75]
	_ = x[ErrInvalidRegion-76]
	_ = x[ErrInvalidServiceS3-77]
	_ = x[ErrInvalidServiceSTS-78]
	_ = x[ErrInvalidRequestVersion-79]
	_ = x[ErrMissingSignTag-80]
	_ = x[ErrMissingSignHeadersTag-81]
	_ = x[ErrMalformedDate-82]
	_ = x[ErrMalformedPresignedDate-83]
	_ = x[ErrMalformedCredentialDate-84]
	_ = x[ErrMalformedCredentialRegion-85]
	_ = x[ErrMalformedExpires-86]
	_ = x[ErrNegativeExpires-87]
	_ = x[ErrAuthHeaderEmpty-88]
	_ = x[ErrExpiredPresignRequest-89]
	_ = x[ErrRequestNotReadyYet-90]
	_ = x[ErrUnsignedHeaders-91]
	_ = x[ErrMissingDateHeader-92]
	_ = x[ErrInvalidQuerySignatureAlgo-93]
	_ = x[ErrInvalidQueryParams-94]
	_ = x[ErrBucketAlreadyOwnedByYou-95]
	_ = x[ErrInvalidDuration-96]
	_ = x[ErrBucketAlreadyExists-97]
	_ = x[ErrMetadataTooLarge-98]
	_ = x[ErrUnsupportedMetadata-99]
	_ = x[ErrMaximumExpires-100]
	_ = x[ErrSlowDown-101]
	_ = x[ErrInvalidPrefixMarker-102]
	_ = x[ErrBadRequest-103]
	_ = x[ErrKeyTooLongError-104]
	_ = x[ErrInvalidBucketObjectLockConfiguration-105]
	_ = x[ErrObjectLockConfigurationNotFound-106]
	_ = x[ErrObjectLockConfigurationNotAllowed-107]
	_ = x[ErrNoSuchObjectLockConfiguration-108]
	_ = x[ErrObjectLocked-109]
	_ = x[ErrInvalidRetentionDate-110]
	_ = x[ErrPastObjectLockRetainDate-111]
	_ = x[ErrUnknownWORMModeDirective-112]
	_ = x[ErrBucketTaggingNotFound-113]
	_ = x[ErrObjectLockInvalidHeaders-114]
	_ = x[ErrInvalidTagDirective-115]
	_ = x[ErrInvalidEncryptionMethod-116]
	_ = x[ErrInsecureSSECustomerRequest-117]
	_ = x[ErrSSEMultipartEncrypted-118]
	_ = x[ErrSSEEncryptedObject-119]
	_ = x[ErrInvalidEncryptionParameters-120]
	_ = x[ErrInvalidSSECustomerAlgorithm-121]
	_ = x[ErrInvalidSSECustomerKey-122]
	_ = x[ErrMissingSSECustomerKey-123]
	_ = x[ErrMissingSSECustomerKeyMD5-124]
	_ = x[ErrSSECustomerKeyMD5Mismatch-125]
	_ = x[ErrInvalidSSECustomerParameters-126]
	_ = x[ErrIncompatibleEncryptionMethod-127]
	_ = x[ErrKMSNotConfigured-128]
	_ = x[ErrKMSAuthFailure-129]
	_ = x[ErrNoAccessKey-130]
	_ = x[ErrInvalidToken-131]
	_ = x[ErrEventNotification-132]
	_ = x[ErrARNNotification-133]
	_ = x[ErrRegionNotification-134]
	_ = x[ErrOverlappingFilterNotification-135]
	_ = x[ErrFilterNameInvalid-136]
	_ = x[ErrFilterNamePrefix-137]
	_ = x[ErrFilterNameSuffix-138]
	_ = x[ErrFilterValueInvalid-139]
	_ = x[ErrOverlappingConfigs-140]
	_ = x[ErrUnsupportedNotification-141]
	_ = x[ErrContentSHA256Mismatch-142]
	_ = x[ErrReadQuorum-143]
	_ = x[ErrWriteQuorum-144]
	_ = x[ErrParentIsObject-145]
	_ = x[ErrStorageFull-146]
	_ = x[ErrRequestBodyParse-147]
	_ = x[ErrObjectExistsAsDirectory-148]
	_ = x[ErrInvalidObjectName-149]
	_ = x[ErrInvalidObjectNamePrefixSlash-150]
	_ = x[ErrInvalidResourceName-151]
	_ = x[ErrServerNotInitialized-152]
	_ = x[ErrOperationTimedOut-153]
	_ = x[ErrClientDisconnected-154]
	_ = x[ErrOperationMaxedOut-155]
	_ = x[ErrInvalidRequest-156]
	_ = x[ErrInvalidStorageClass-157]
	_ = x[ErrBackendDown-158]
	_ = x[ErrMalformedJSON-159]
	_ = x[ErrAdminNoSuchUser-160]
	_ = x[ErrAdminNoSuchGroup-161]
	_ = x[ErrAdminGroupNotEmpty-162]
	_ = x[ErrAdminNoSuchPolicy-163]
	_ = x[ErrAdminInvalidArgument-164]
	_ = x[ErrAdminInvalidAccessKey-165]
	_ = x[ErrAdminInvalidSecretKey-166]
	_ = x[ErrAdminConfigNoQuorum-167]
	_ = x[ErrAdminConfigTooLarge-168]
	_ = x[ErrAdminConfigBadJSON-169]
	_ = x[ErrAdminConfigDuplicateKeys-170]
	_ = x[ErrAdminCredentialsMismatch-171]
	_ = x[ErrInsecureClientRequest-172]
	_ = x[ErrObjectTampered-173]
	_ = x[ErrAdminBucketQuotaExceeded-174]
	_ = x[ErrAdminNoSuchQuotaConfiguration-175]
	_ = x[ErrHealNotImplemented-176]
	_ = x[ErrHealNoSuchProcess-177]
	_ = x[ErrHealInvalidClientToken-178]
	_ = x[ErrHealMissingBucket-179]
	_ = x[ErrHealAlreadyRunning-180]
	_ = x[ErrHealOverlappingPaths-181]
	_ = x[ErrIncorrectContinuationToken-182]
	_ = x[ErrEmptyRequestBody-183]
	_ = x[ErrUnsupportedFunction-184]
	_ = x[ErrInvalidExpressionType-185]
	_ = x[ErrBusy-186]
	_ = x[ErrUnauthorizedAccess-187]
	_ = x[ErrExpressionTooLong-188]
	_ = x[ErrIllegalSQLFunctionArgument-189]
	_ = x[ErrInvalidKeyPath-190]
	_ = x[ErrInvalidCompressionFormat-191]
	_ = x[ErrInvalidFileHeaderInfo-192]
	_ = x[ErrInvalidJSONType-193]
	_ = x[ErrInvalidQuoteFields-194]
	_ = x[ErrInvalidRequestParameter-195]
	_ = x[ErrInvalidDataType-196]
	_ = x[ErrInvalidTextEncoding-197]
	_ = x[ErrInvalidDataSource-198]
	_ = x[ErrInvalidTableAlias-199]
	_ = x[ErrMissingRequiredParameter-200]
	_ = x[ErrObjectSerializationConflict-201]
	_ = x[ErrUnsupportedSQLOperation-202]
	_ = x[ErrUnsupportedSQLStructure-203]
	_ = x[ErrUnsupportedSyntax-204]
	_ = x[ErrUnsupportedRangeHeader-205]
	_ = x[ErrLexerInvalidChar-206]
	_ = x[ErrLexerInvalidOperator-207]
	_ = x[ErrLexerInvalidLiteral-208]
	_ = x[ErrLexerInvalidIONLiteral-209]
	_ = x[ErrParseExpectedDatePart-210]
	_ = x[ErrParseExpectedKeyword-211]
	_ = x[ErrParseExpectedTokenType-212]
	_ = x[ErrParseExpected2TokenTypes-213]
	_ = x[ErrParseExpectedNumber-214]
	_ = x[ErrParseExpectedRightParenBuiltinFunctionCall-215]
	_ = x[ErrParseExpectedTypeName-216]
	_ = x[ErrParseExpectedWhenClause-217]
	_ = x[ErrParseUnsupportedToken-218]
	_ = x[ErrParseUnsupportedLiteralsGroupBy-219]
	_ = x[ErrParseExpectedMember-220]
	_ = x[ErrParseUnsupportedSelect-221]
	_ = x[ErrParseUnsupportedCase-222]
	_ = x[ErrParseUnsupportedCaseClause-223]
	_ = x[ErrParseUnsupportedAlias-224]
	_ = x[ErrParseUnsupportedSyntax-225]
	_ = x[ErrParseUnknownOperator-226]
	_ = x[ErrParseMissingIdentAfterAt-227]
	_ = x[ErrParseUnexpectedOperator-228]
	_ = x[ErrParseUnexpectedTerm-229]
	_ = x[ErrParseUnexpectedToken-230]
	_ = x[ErrParseUnexpectedKeyword-231]
	_ = x[ErrParseExpectedExpression-232]
	_ = x[ErrParseExpectedLeftParenAfterCast-233]
	_ = x[ErrParseExpectedLeftParenValueConstructor-234]
	_ = x[ErrParseExpectedLeftParenBuiltinFunctionCall-235]
	_ = x[ErrParseExpectedArgumentDelimiter-236]
	_ = x[ErrParseCastArity-237]
	_ = x[ErrParseInvalidTypeParam-238]
	_ = x[ErrParseEmptySelect-239]
	_ = x[ErrParseSelectMissingFrom-240]
	_ = x[ErrParseExpectedIdentForGroupName-241]
	_ = x[ErrParseExpectedIdentForAlias-242]
	_ = x[ErrParseUnsupportedCallWithStar-243]
	_ = x[ErrParseNonUnaryAgregateFunctionCall-244]
	_ = x[ErrParseMalformedJoin-245]
	_ = x[ErrParseExpectedIdentForAt-246]
	_ = x[ErrParseAsteriskIsNotAloneInSelectList-247]
	_ = x[ErrParseCannotMixSqbAndWildcardInSelectList-248]
	_ = x[ErrParseInvalidContextForWildcardInSelectList-249]
	_ = x[ErrIncorrectSQLFunctionArgumentType-250]
	_ = x[ErrValueParseFailure-251]
	_ = x[ErrEvaluatorInvalidArguments-252]
	_ = x[ErrIntegerOverflow-253]
	_ = x[ErrLikeInvalidInputs-254]
	_ = x[ErrCastFailed-255]
	_ = x[ErrInvalidCast-256]
	_ = x[ErrEvaluatorInvalidTimestampFormatPattern-257]
	_ = x[ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing-258]
	_ = x[ErrEvaluatorTimestampFormatPatternDuplicateFields-259]
	_ = x[ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch-260]
	_ = x[ErrEvaluatorUnterminatedTimestampFormatPatternToken-261]
	_ = x[ErrEvaluatorInvalidTimestampFormatPatternToken-262]
	_ = x[ErrEvaluatorInvalidTimestampFormatPatternSymbol-263]
	_ = x[ErrEvaluatorBindingDoesNotExist-264]
	_ = x[ErrMissingHeaders-265]
	_ = x[ErrInvalidColumnIndex-266]
	_ = x[ErrAdminConfigNotificationTargetsFailed-267]
	_ = x[ErrAdminProfilerNotEnabled-268]
	_ = x[ErrInvalidDecompressedSize-269]
	_ = x[ErrAddUserInvalidArgument-270]
	_ = x[ErrAdminAccountNotEligible-271]
	_ = x[ErrAccountNotEligible-272]
	_ = x[ErrAdminServiceAccountNotFound-273]
	_ = x[ErrPostPolicyConditionInvalidFormat-274]
}

const _APIErrorCode_name = "NoneAccessDeniedBadDigestEntityTooSmallEntityTooLargePolicyTooLargeIncompleteBodyInternalErrorInvalidAccessKeyIDInvalidBucketNameInvalidDigestInvalidRangeInvalidRangePartNumberInvalidCopyPartRangeInvalidCopyPartRangeSourceInvalidMaxKeysInvalidEncodingMethodInvalidMaxUploadsInvalidMaxPartsInvalidPartNumberMarkerInvalidPartNumberInvalidRequestBodyInvalidCopySourceInvalidMetadataDirectiveInvalidCopyDestInvalidPolicyDocumentInvalidObjectStateMalformedXMLMissingContentLengthMissingContentMD5MissingRequestBodyErrorMissingSecurityHeaderNoSuchBucketNoSuchBucketPolicyNoSuchBucketLifecycleNoSuchLifecycleConfigurationNoSuchBucketSSEConfigNoSuchCORSConfigurationNoSuchWebsiteConfigurationReplicationConfigurationNotFoundErrorRemoteDestinationNotFoundErrorReplicationDestinationMissingLockRemoteTargetNotFoundErrorReplicationRemoteConnectionErrorBucketRemoteIdenticalToSourceBucketRemoteAlreadyExistsBucketRemoteLabelInUseBucketRemoteArnTypeInvalidBucketRemoteArnInvalidBucketRemoteRemoveDisallowedRemoteTargetNotVersionedErrorReplicationSourceNotVersionedErrorReplicationNeedsVersioningErrorReplicationBucketNeedsVersioningErrorObjectRestoreAlreadyInProgressNoSuchKeyNoSuchUploadInvalidVersionIDNoSuchVersionNotImplementedPreconditionFailedRequestTimeTooSkewedSignatureDoesNotMatchMethodNotAllowedInvalidPartInvalidPartOrderAuthorizationHeaderMalformedMalformedPOSTRequestPOSTFileRequiredSignatureVersionNotSupportedBucketNotEmptyAllAccessDisabledMalformedPolicyMissingFieldsMissingCredTagCredMalformedInvalidRegionInvalidServiceS3InvalidServiceSTSInvalidRequestVersionMissingSignTagMissingSignHeadersTagMalformedDateMalformedPresignedDateMalformedCredentialDateMalformedCredentialRegionMalformedExpiresNegativeExpiresAuthHeaderEmptyExpiredPresignRequestRequestNotReadyYetUnsignedHeadersMissingDateHeaderInvalidQuerySignatureAlgoInvalidQueryParamsBucketAlreadyOwnedByYouInvalidDurationBucketAlreadyExistsMetadataTooLargeUnsupportedMetadataMaximumExpiresSlowDownInvalidPrefixMarkerBadRequestKeyTooLongErrorInvalidBucketObjectLockConfigurationObjectLockConfigurationNotFoundObjectLockConfigurationNotAllowedNoSuchObjectLockConfigurationObjectLockedInvalidRetentionDatePastObjectLockRetainDateUnknownWORMModeDirectiveBucketTaggingNotFoundObjectLockInvalidHeadersInvalidTagDirectiveInvalidEncryptionMethodInsecureSSECustomerRequestSSEMultipartEncryptedSSEEncryptedObjectInvalidEncryptionParametersInvalidSSECustomerAlgorithmInvalidSSECustomerKeyMissingSSECustomerKeyMissingSSECustomerKeyMD5SSECustomerKeyMD5MismatchInvalidSSECustomerParametersIncompatibleEncryptionMethodKMSNotConfiguredKMSAuthFailureNoAccessKeyInvalidTokenEventNotificationARNNotificationRegionNotificationOverlappingFilterNotificationFilterNameInvalidFilterNamePrefixFilterNameSuffixFilterValueInvalidOverlappingConfigsUnsupportedNotificationContentSHA256MismatchReadQuorumWriteQuorumParentIsObjectStorageFullRequestBodyParseObjectExistsAsDirectoryInvalidObjectNameInvalidObjectNamePrefixSlashInvalidResourceNameServerNotInitializedOperationTimedOutClientDisconnectedOperationMaxedOutInvalidRequestInvalidStorageClassBackendDownMalformedJSONAdminNoSuchUserAdminNoSuchGroupAdminGroupNotEmptyAdminNoSuchPolicyAdminInvalidArgumentAdminInvalidAccessKeyAdminInvalidSecretKeyAdminConfigNoQuorumAdminConfigTooLargeAdminConfigBadJSONAdminConfigDuplicateKeysAdminCredentialsMismatchInsecureClientRequestObjectTamperedAdminBucketQuotaExceededAdminNoSuchQuotaConfigurationHealNotImplementedHealNoSuchProcessHealInvalidClientTokenHealMissingBucketHealAlreadyRunningHealOverlappingPathsIncorrectContinuationTokenEmptyRequestBodyUnsupportedFunctionInvalidExpressionTypeBusyUnauthorizedAccessExpressionTooLongIllegalSQLFunctionArgumentInvalidKeyPathInvalidCompressionFormatInvalidFileHeaderInfoInvalidJSONTypeInvalidQuoteFieldsInvalidRequestParameterInvalidDataTypeInvalidTextEncodingInvalidDataSourceInvalidTableAliasMissingRequiredParameterObjectSerializationConflictUnsupportedSQLOperationUnsupportedSQLStructureUnsupportedSyntaxUnsupportedRangeHeaderLexerInvalidCharLexerInvalidOperatorLexerInvalidLiteralLexerInvalidIONLiteralParseExpectedDatePartParseExpectedKeywordParseExpectedTokenTypeParseExpected2TokenTypesParseExpectedNumberParseExpectedRightParenBuiltinFunctionCallParseExpectedTypeNameParseExpectedWhenClauseParseUnsupportedTokenParseUnsupportedLiteralsGroupByParseExpectedMemberParseUnsupportedSelectParseUnsupportedCaseParseUnsupportedCaseClauseParseUnsupportedAliasParseUnsupportedSyntaxParseUnknownOperatorParseMissingIdentAfterAtParseUnexpectedOperatorParseUnexpectedTermParseUnexpectedTokenParseUnexpectedKeywordParseExpectedExpressionParseExpectedLeftParenAfterCastParseExpectedLeftParenValueConstructorParseExpectedLeftParenBuiltinFunctionCallParseExpectedArgumentDelimiterParseCastArityParseInvalidTypeParamParseEmptySelectParseSelectMissingFromParseExpectedIdentForGroupNameParseExpectedIdentForAliasParseUnsupportedCallWithStarParseNonUnaryAgregateFunctionCallParseMalformedJoinParseExpectedIdentForAtParseAsteriskIsNotAloneInSelectListParseCannotMixSqbAndWildcardInSelectListParseInvalidContextForWildcardInSelectListIncorrectSQLFunctionArgumentTypeValueParseFailureEvaluatorInvalidArgumentsIntegerOverflowLikeInvalidInputsCastFailedInvalidCastEvaluatorInvalidTimestampFormatPatternEvaluatorInvalidTimestampFormatPatternSymbolForParsingEvaluatorTimestampFormatPatternDuplicateFieldsEvaluatorTimestampFormatPatternHourClockAmPmMismatchEvaluatorUnterminatedTimestampFormatPatternTokenEvaluatorInvalidTimestampFormatPatternTokenEvaluatorInvalidTimestampFormatPatternSymbolEvaluatorBindingDoesNotExistMissingHeadersInvalidColumnIndexAdminConfigNotificationTargetsFailedAdminProfilerNotEnabledInvalidDecompressedSizeAddUserInvalidArgumentAdminAccountNotEligibleAccountNotEligibleAdminServiceAccountNotFoundPostPolicyConditionInvalidFormat"

var _APIErrorCode_index = [...]uint16{0, 4, 16, 25, 39, 53, 67, 81, 94, 112, 129, 142, 154, 176, 196, 222, 236, 257, 274, 289, 312, 329, 347, 364, 388, 403, 424, 442, 454, 474, 491, 514, 535, 547, 565, 586, 614, 635, 658, 684, 721, 751, 784, 809, 841, 870, 895, 917, 943, 965, 993, 1022, 1056, 1087, 1124, 1154, 1163, 1175, 1191, 1204, 1218, 1236, 1256, 1277, 1293, 1304, 1320, 1348, 1368, 1384, 1412, 1426, 1443, 1458, 1471, 1485, 1498, 1511, 1527, 1544, 1565, 1579, 1600, 1613, 1635, 1658, 1683, 1699, 1714, 1729, 1750, 1768, 1783, 1800, 1825, 1843, 1866, 1881, 1900, 1916, 1935, 1949, 1957, 1976, 1986, 2001, 2037, 2068, 2101, 2130, 2142, 2162, 2186, 2210, 2231, 2255, 2274, 2297, 2323, 2344, 2362, 2389, 2416, 2437, 2458, 2482, 2507, 2535, 2563, 2579, 2593, 2604, 2616, 2633, 2648, 2666, 2695, 2712, 2728, 2744, 2762, 2780, 2803, 2824, 2834, 2845, 2859, 2870, 2886, 2909, 2926, 2954, 2973, 2993, 3010, 3028, 3045, 3059, 3078, 3089, 3102, 3117, 3133, 3151, 3168, 3188, 3209, 3230, 3249, 3268, 3286, 3310, 3334, 3355, 3369, 3393, 3422, 3440, 3457, 3479, 3496, 3514, 3534, 3560, 3576, 3595, 3616, 3620, 3638, 3655, 3681, 3695, 3719, 3740, 3755, 3773, 3796, 3811, 3830, 3847, 3864, 3888, 3915, 3938, 3961, 3978, 4000, 4016, 4036, 4055, 4077, 4098, 4118, 4140, 4164, 4183, 4225, 4246, 4269, 4290, 4321, 4340, 4362, 4382, 4408, 4429, 4451, 4471, 4495, 4518, 4537, 4557, 4579, 4602, 4633, 4671, 4712, 4742, 4756, 4777, 4793, 4815, 4845, 4871, 4899, 4932, 4950, 4973, 5008, 5048, 5090, 5122, 5139, 5164, 5179, 5196, 5206, 5217, 5255, 5309, 5355, 5407, 5455, 5498, 5542, 5570, 5584, 5602, 5638, 5661, 5684, 5706, 5729, 5747, 5774, 5806}

func (i APIErrorCode) String() string {
	if i < 0 || i >= APIErrorCode(len(_APIErrorCode_index)-1) {
		return "APIErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _APIErrorCode_name[_APIErrorCode_index[i]:_APIErrorCode_index[i+1]]
}
