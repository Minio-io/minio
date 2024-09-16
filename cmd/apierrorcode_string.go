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
	_ = x[ErrAccessKeyDisabled-9]
	_ = x[ErrLDAPAccessKeyDisabled-10]
	_ = x[ErrInvalidArgument-11]
	_ = x[ErrInvalidBucketName-12]
	_ = x[ErrInvalidDigest-13]
	_ = x[ErrInvalidRange-14]
	_ = x[ErrInvalidRangePartNumber-15]
	_ = x[ErrInvalidCopyPartRange-16]
	_ = x[ErrInvalidCopyPartRangeSource-17]
	_ = x[ErrInvalidMaxKeys-18]
	_ = x[ErrInvalidEncodingMethod-19]
	_ = x[ErrInvalidMaxUploads-20]
	_ = x[ErrInvalidMaxParts-21]
	_ = x[ErrInvalidPartNumberMarker-22]
	_ = x[ErrInvalidPartNumber-23]
	_ = x[ErrInvalidRequestBody-24]
	_ = x[ErrInvalidCopySource-25]
	_ = x[ErrInvalidMetadataDirective-26]
	_ = x[ErrInvalidCopyDest-27]
	_ = x[ErrInvalidPolicyDocument-28]
	_ = x[ErrInvalidObjectState-29]
	_ = x[ErrMalformedXML-30]
	_ = x[ErrMissingContentLength-31]
	_ = x[ErrMissingContentMD5-32]
	_ = x[ErrMissingRequestBodyError-33]
	_ = x[ErrMissingSecurityHeader-34]
	_ = x[ErrNoSuchBucket-35]
	_ = x[ErrNoSuchBucketPolicy-36]
	_ = x[ErrNoSuchBucketLifecycle-37]
	_ = x[ErrNoSuchLifecycleConfiguration-38]
	_ = x[ErrInvalidLifecycleWithObjectLock-39]
	_ = x[ErrNoSuchBucketSSEConfig-40]
	_ = x[ErrNoSuchCORSConfiguration-41]
	_ = x[ErrNoSuchWebsiteConfiguration-42]
	_ = x[ErrReplicationConfigurationNotFoundError-43]
	_ = x[ErrRemoteDestinationNotFoundError-44]
	_ = x[ErrReplicationDestinationMissingLock-45]
	_ = x[ErrRemoteTargetNotFoundError-46]
	_ = x[ErrReplicationRemoteConnectionError-47]
	_ = x[ErrReplicationBandwidthLimitError-48]
	_ = x[ErrBucketRemoteIdenticalToSource-49]
	_ = x[ErrBucketRemoteAlreadyExists-50]
	_ = x[ErrBucketRemoteLabelInUse-51]
	_ = x[ErrBucketRemoteArnTypeInvalid-52]
	_ = x[ErrBucketRemoteArnInvalid-53]
	_ = x[ErrBucketRemoteRemoveDisallowed-54]
	_ = x[ErrRemoteTargetNotVersionedError-55]
	_ = x[ErrReplicationSourceNotVersionedError-56]
	_ = x[ErrReplicationNeedsVersioningError-57]
	_ = x[ErrReplicationBucketNeedsVersioningError-58]
	_ = x[ErrReplicationDenyEditError-59]
	_ = x[ErrRemoteTargetDenyAddError-60]
	_ = x[ErrReplicationNoExistingObjects-61]
	_ = x[ErrReplicationValidationError-62]
	_ = x[ErrReplicationPermissionCheckError-63]
	_ = x[ErrObjectRestoreAlreadyInProgress-64]
	_ = x[ErrNoSuchKey-65]
	_ = x[ErrNoSuchUpload-66]
	_ = x[ErrInvalidVersionID-67]
	_ = x[ErrNoSuchVersion-68]
	_ = x[ErrNotImplemented-69]
	_ = x[ErrPreconditionFailed-70]
	_ = x[ErrRequestTimeTooSkewed-71]
	_ = x[ErrSignatureDoesNotMatch-72]
	_ = x[ErrMethodNotAllowed-73]
	_ = x[ErrInvalidPart-74]
	_ = x[ErrInvalidPartOrder-75]
	_ = x[ErrMissingPart-76]
	_ = x[ErrAuthorizationHeaderMalformed-77]
	_ = x[ErrMalformedPOSTRequest-78]
	_ = x[ErrPOSTFileRequired-79]
	_ = x[ErrSignatureVersionNotSupported-80]
	_ = x[ErrBucketNotEmpty-81]
	_ = x[ErrAllAccessDisabled-82]
	_ = x[ErrPolicyInvalidVersion-83]
	_ = x[ErrMissingFields-84]
	_ = x[ErrMissingCredTag-85]
	_ = x[ErrCredMalformed-86]
	_ = x[ErrInvalidRegion-87]
	_ = x[ErrInvalidServiceS3-88]
	_ = x[ErrInvalidServiceSTS-89]
	_ = x[ErrInvalidRequestVersion-90]
	_ = x[ErrMissingSignTag-91]
	_ = x[ErrMissingSignHeadersTag-92]
	_ = x[ErrMalformedDate-93]
	_ = x[ErrMalformedPresignedDate-94]
	_ = x[ErrMalformedCredentialDate-95]
	_ = x[ErrMalformedExpires-96]
	_ = x[ErrNegativeExpires-97]
	_ = x[ErrAuthHeaderEmpty-98]
	_ = x[ErrExpiredPresignRequest-99]
	_ = x[ErrRequestNotReadyYet-100]
	_ = x[ErrUnsignedHeaders-101]
	_ = x[ErrMissingDateHeader-102]
	_ = x[ErrInvalidQuerySignatureAlgo-103]
	_ = x[ErrInvalidQueryParams-104]
	_ = x[ErrBucketAlreadyOwnedByYou-105]
	_ = x[ErrInvalidDuration-106]
	_ = x[ErrBucketAlreadyExists-107]
	_ = x[ErrMetadataTooLarge-108]
	_ = x[ErrUnsupportedMetadata-109]
	_ = x[ErrUnsupportedHostHeader-110]
	_ = x[ErrMaximumExpires-111]
	_ = x[ErrSlowDownRead-112]
	_ = x[ErrSlowDownWrite-113]
	_ = x[ErrMaxVersionsExceeded-114]
	_ = x[ErrInvalidPrefixMarker-115]
	_ = x[ErrBadRequest-116]
	_ = x[ErrKeyTooLongError-117]
	_ = x[ErrInvalidBucketObjectLockConfiguration-118]
	_ = x[ErrObjectLockConfigurationNotFound-119]
	_ = x[ErrObjectLockConfigurationNotAllowed-120]
	_ = x[ErrNoSuchObjectLockConfiguration-121]
	_ = x[ErrObjectLocked-122]
	_ = x[ErrInvalidRetentionDate-123]
	_ = x[ErrPastObjectLockRetainDate-124]
	_ = x[ErrUnknownWORMModeDirective-125]
	_ = x[ErrBucketTaggingNotFound-126]
	_ = x[ErrObjectLockInvalidHeaders-127]
	_ = x[ErrInvalidTagDirective-128]
	_ = x[ErrPolicyAlreadyAttached-129]
	_ = x[ErrPolicyNotAttached-130]
	_ = x[ErrExcessData-131]
	_ = x[ErrInvalidEncryptionMethod-132]
	_ = x[ErrInvalidEncryptionKeyID-133]
	_ = x[ErrInsecureSSECustomerRequest-134]
	_ = x[ErrSSEMultipartEncrypted-135]
	_ = x[ErrSSEEncryptedObject-136]
	_ = x[ErrInvalidEncryptionParameters-137]
	_ = x[ErrInvalidEncryptionParametersSSEC-138]
	_ = x[ErrInvalidSSECustomerAlgorithm-139]
	_ = x[ErrInvalidSSECustomerKey-140]
	_ = x[ErrMissingSSECustomerKey-141]
	_ = x[ErrMissingSSECustomerKeyMD5-142]
	_ = x[ErrSSECustomerKeyMD5Mismatch-143]
	_ = x[ErrInvalidSSECustomerParameters-144]
	_ = x[ErrIncompatibleEncryptionMethod-145]
	_ = x[ErrKMSNotConfigured-146]
	_ = x[ErrKMSKeyNotFoundException-147]
	_ = x[ErrKMSDefaultKeyAlreadyConfigured-148]
	_ = x[ErrNoAccessKey-149]
	_ = x[ErrInvalidToken-150]
	_ = x[ErrEventNotification-151]
	_ = x[ErrARNNotification-152]
	_ = x[ErrRegionNotification-153]
	_ = x[ErrOverlappingFilterNotification-154]
	_ = x[ErrFilterNameInvalid-155]
	_ = x[ErrFilterNamePrefix-156]
	_ = x[ErrFilterNameSuffix-157]
	_ = x[ErrFilterValueInvalid-158]
	_ = x[ErrOverlappingConfigs-159]
	_ = x[ErrUnsupportedNotification-160]
	_ = x[ErrContentSHA256Mismatch-161]
	_ = x[ErrContentChecksumMismatch-162]
	_ = x[ErrStorageFull-163]
	_ = x[ErrRequestBodyParse-164]
	_ = x[ErrObjectExistsAsDirectory-165]
	_ = x[ErrInvalidObjectName-166]
	_ = x[ErrInvalidObjectNamePrefixSlash-167]
	_ = x[ErrInvalidResourceName-168]
	_ = x[ErrInvalidLifecycleQueryParameter-169]
	_ = x[ErrServerNotInitialized-170]
	_ = x[ErrBucketMetadataNotInitialized-171]
	_ = x[ErrRequestTimedout-172]
	_ = x[ErrClientDisconnected-173]
	_ = x[ErrTooManyRequests-174]
	_ = x[ErrInvalidRequest-175]
	_ = x[ErrTransitionStorageClassNotFoundError-176]
	_ = x[ErrInvalidStorageClass-177]
	_ = x[ErrBackendDown-178]
	_ = x[ErrMalformedJSON-179]
	_ = x[ErrAdminNoSuchUser-180]
	_ = x[ErrAdminNoSuchUserLDAPWarn-181]
	_ = x[ErrAdminLDAPExpectedLoginName-182]
	_ = x[ErrAdminNoSuchGroup-183]
	_ = x[ErrAdminGroupNotEmpty-184]
	_ = x[ErrAdminGroupDisabled-185]
	_ = x[ErrAdminInvalidGroupName-186]
	_ = x[ErrAdminNoSuchJob-187]
	_ = x[ErrAdminNoSuchPolicy-188]
	_ = x[ErrAdminPolicyChangeAlreadyApplied-189]
	_ = x[ErrAdminInvalidArgument-190]
	_ = x[ErrAdminInvalidAccessKey-191]
	_ = x[ErrAdminInvalidSecretKey-192]
	_ = x[ErrAdminConfigNoQuorum-193]
	_ = x[ErrAdminConfigTooLarge-194]
	_ = x[ErrAdminConfigBadJSON-195]
	_ = x[ErrAdminNoSuchConfigTarget-196]
	_ = x[ErrAdminConfigEnvOverridden-197]
	_ = x[ErrAdminConfigDuplicateKeys-198]
	_ = x[ErrAdminConfigInvalidIDPType-199]
	_ = x[ErrAdminConfigLDAPNonDefaultConfigName-200]
	_ = x[ErrAdminConfigLDAPValidation-201]
	_ = x[ErrAdminConfigIDPCfgNameAlreadyExists-202]
	_ = x[ErrAdminConfigIDPCfgNameDoesNotExist-203]
	_ = x[ErrInsecureClientRequest-204]
	_ = x[ErrObjectTampered-205]
	_ = x[ErrAdminLDAPNotEnabled-206]
	_ = x[ErrSiteReplicationInvalidRequest-207]
	_ = x[ErrSiteReplicationPeerResp-208]
	_ = x[ErrSiteReplicationBackendIssue-209]
	_ = x[ErrSiteReplicationServiceAccountError-210]
	_ = x[ErrSiteReplicationBucketConfigError-211]
	_ = x[ErrSiteReplicationBucketMetaError-212]
	_ = x[ErrSiteReplicationIAMError-213]
	_ = x[ErrSiteReplicationConfigMissing-214]
	_ = x[ErrSiteReplicationIAMConfigMismatch-215]
	_ = x[ErrAdminRebalanceAlreadyStarted-216]
	_ = x[ErrAdminRebalanceNotStarted-217]
	_ = x[ErrAdminBucketQuotaExceeded-218]
	_ = x[ErrAdminNoSuchQuotaConfiguration-219]
	_ = x[ErrHealNotImplemented-220]
	_ = x[ErrHealNoSuchProcess-221]
	_ = x[ErrHealInvalidClientToken-222]
	_ = x[ErrHealMissingBucket-223]
	_ = x[ErrHealAlreadyRunning-224]
	_ = x[ErrHealOverlappingPaths-225]
	_ = x[ErrIncorrectContinuationToken-226]
	_ = x[ErrEmptyRequestBody-227]
	_ = x[ErrUnsupportedFunction-228]
	_ = x[ErrInvalidExpressionType-229]
	_ = x[ErrBusy-230]
	_ = x[ErrUnauthorizedAccess-231]
	_ = x[ErrExpressionTooLong-232]
	_ = x[ErrIllegalSQLFunctionArgument-233]
	_ = x[ErrInvalidKeyPath-234]
	_ = x[ErrInvalidCompressionFormat-235]
	_ = x[ErrInvalidFileHeaderInfo-236]
	_ = x[ErrInvalidJSONType-237]
	_ = x[ErrInvalidQuoteFields-238]
	_ = x[ErrInvalidRequestParameter-239]
	_ = x[ErrInvalidDataType-240]
	_ = x[ErrInvalidTextEncoding-241]
	_ = x[ErrInvalidDataSource-242]
	_ = x[ErrInvalidTableAlias-243]
	_ = x[ErrMissingRequiredParameter-244]
	_ = x[ErrObjectSerializationConflict-245]
	_ = x[ErrUnsupportedSQLOperation-246]
	_ = x[ErrUnsupportedSQLStructure-247]
	_ = x[ErrUnsupportedSyntax-248]
	_ = x[ErrUnsupportedRangeHeader-249]
	_ = x[ErrLexerInvalidChar-250]
	_ = x[ErrLexerInvalidOperator-251]
	_ = x[ErrLexerInvalidLiteral-252]
	_ = x[ErrLexerInvalidIONLiteral-253]
	_ = x[ErrParseExpectedDatePart-254]
	_ = x[ErrParseExpectedKeyword-255]
	_ = x[ErrParseExpectedTokenType-256]
	_ = x[ErrParseExpected2TokenTypes-257]
	_ = x[ErrParseExpectedNumber-258]
	_ = x[ErrParseExpectedRightParenBuiltinFunctionCall-259]
	_ = x[ErrParseExpectedTypeName-260]
	_ = x[ErrParseExpectedWhenClause-261]
	_ = x[ErrParseUnsupportedToken-262]
	_ = x[ErrParseUnsupportedLiteralsGroupBy-263]
	_ = x[ErrParseExpectedMember-264]
	_ = x[ErrParseUnsupportedSelect-265]
	_ = x[ErrParseUnsupportedCase-266]
	_ = x[ErrParseUnsupportedCaseClause-267]
	_ = x[ErrParseUnsupportedAlias-268]
	_ = x[ErrParseUnsupportedSyntax-269]
	_ = x[ErrParseUnknownOperator-270]
	_ = x[ErrParseMissingIdentAfterAt-271]
	_ = x[ErrParseUnexpectedOperator-272]
	_ = x[ErrParseUnexpectedTerm-273]
	_ = x[ErrParseUnexpectedToken-274]
	_ = x[ErrParseUnexpectedKeyword-275]
	_ = x[ErrParseExpectedExpression-276]
	_ = x[ErrParseExpectedLeftParenAfterCast-277]
	_ = x[ErrParseExpectedLeftParenValueConstructor-278]
	_ = x[ErrParseExpectedLeftParenBuiltinFunctionCall-279]
	_ = x[ErrParseExpectedArgumentDelimiter-280]
	_ = x[ErrParseCastArity-281]
	_ = x[ErrParseInvalidTypeParam-282]
	_ = x[ErrParseEmptySelect-283]
	_ = x[ErrParseSelectMissingFrom-284]
	_ = x[ErrParseExpectedIdentForGroupName-285]
	_ = x[ErrParseExpectedIdentForAlias-286]
	_ = x[ErrParseUnsupportedCallWithStar-287]
	_ = x[ErrParseNonUnaryAggregateFunctionCall-288]
	_ = x[ErrParseMalformedJoin-289]
	_ = x[ErrParseExpectedIdentForAt-290]
	_ = x[ErrParseAsteriskIsNotAloneInSelectList-291]
	_ = x[ErrParseCannotMixSqbAndWildcardInSelectList-292]
	_ = x[ErrParseInvalidContextForWildcardInSelectList-293]
	_ = x[ErrIncorrectSQLFunctionArgumentType-294]
	_ = x[ErrValueParseFailure-295]
	_ = x[ErrEvaluatorInvalidArguments-296]
	_ = x[ErrIntegerOverflow-297]
	_ = x[ErrLikeInvalidInputs-298]
	_ = x[ErrCastFailed-299]
	_ = x[ErrInvalidCast-300]
	_ = x[ErrEvaluatorInvalidTimestampFormatPattern-301]
	_ = x[ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing-302]
	_ = x[ErrEvaluatorTimestampFormatPatternDuplicateFields-303]
	_ = x[ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch-304]
	_ = x[ErrEvaluatorUnterminatedTimestampFormatPatternToken-305]
	_ = x[ErrEvaluatorInvalidTimestampFormatPatternToken-306]
	_ = x[ErrEvaluatorInvalidTimestampFormatPatternSymbol-307]
	_ = x[ErrEvaluatorBindingDoesNotExist-308]
	_ = x[ErrMissingHeaders-309]
	_ = x[ErrInvalidColumnIndex-310]
	_ = x[ErrAdminConfigNotificationTargetsFailed-311]
	_ = x[ErrAdminProfilerNotEnabled-312]
	_ = x[ErrInvalidDecompressedSize-313]
	_ = x[ErrAddUserInvalidArgument-314]
	_ = x[ErrAddUserValidUTF-315]
	_ = x[ErrAdminResourceInvalidArgument-316]
	_ = x[ErrAdminAccountNotEligible-317]
	_ = x[ErrAccountNotEligible-318]
	_ = x[ErrAdminServiceAccountNotFound-319]
	_ = x[ErrPostPolicyConditionInvalidFormat-320]
	_ = x[ErrInvalidChecksum-321]
	_ = x[ErrLambdaARNInvalid-322]
	_ = x[ErrLambdaARNNotFound-323]
	_ = x[ErrInvalidAttributeName-324]
	_ = x[ErrAdminNoAccessKey-325]
	_ = x[ErrAdminNoSecretKey-326]
	_ = x[ErrIAMNotInitialized-327]
	_ = x[apiErrCodeEnd-328]
}

const _APIErrorCode_name = "NoneAccessDeniedBadDigestEntityTooSmallEntityTooLargePolicyTooLargeIncompleteBodyInternalErrorInvalidAccessKeyIDAccessKeyDisabledLDAPAccessKeyDisabledInvalidArgumentInvalidBucketNameInvalidDigestInvalidRangeInvalidRangePartNumberInvalidCopyPartRangeInvalidCopyPartRangeSourceInvalidMaxKeysInvalidEncodingMethodInvalidMaxUploadsInvalidMaxPartsInvalidPartNumberMarkerInvalidPartNumberInvalidRequestBodyInvalidCopySourceInvalidMetadataDirectiveInvalidCopyDestInvalidPolicyDocumentInvalidObjectStateMalformedXMLMissingContentLengthMissingContentMD5MissingRequestBodyErrorMissingSecurityHeaderNoSuchBucketNoSuchBucketPolicyNoSuchBucketLifecycleNoSuchLifecycleConfigurationInvalidLifecycleWithObjectLockNoSuchBucketSSEConfigNoSuchCORSConfigurationNoSuchWebsiteConfigurationReplicationConfigurationNotFoundErrorRemoteDestinationNotFoundErrorReplicationDestinationMissingLockRemoteTargetNotFoundErrorReplicationRemoteConnectionErrorReplicationBandwidthLimitErrorBucketRemoteIdenticalToSourceBucketRemoteAlreadyExistsBucketRemoteLabelInUseBucketRemoteArnTypeInvalidBucketRemoteArnInvalidBucketRemoteRemoveDisallowedRemoteTargetNotVersionedErrorReplicationSourceNotVersionedErrorReplicationNeedsVersioningErrorReplicationBucketNeedsVersioningErrorReplicationDenyEditErrorRemoteTargetDenyAddErrorReplicationNoExistingObjectsReplicationValidationErrorReplicationPermissionCheckErrorObjectRestoreAlreadyInProgressNoSuchKeyNoSuchUploadInvalidVersionIDNoSuchVersionNotImplementedPreconditionFailedRequestTimeTooSkewedSignatureDoesNotMatchMethodNotAllowedInvalidPartInvalidPartOrderMissingPartAuthorizationHeaderMalformedMalformedPOSTRequestPOSTFileRequiredSignatureVersionNotSupportedBucketNotEmptyAllAccessDisabledPolicyInvalidVersionMissingFieldsMissingCredTagCredMalformedInvalidRegionInvalidServiceS3InvalidServiceSTSInvalidRequestVersionMissingSignTagMissingSignHeadersTagMalformedDateMalformedPresignedDateMalformedCredentialDateMalformedExpiresNegativeExpiresAuthHeaderEmptyExpiredPresignRequestRequestNotReadyYetUnsignedHeadersMissingDateHeaderInvalidQuerySignatureAlgoInvalidQueryParamsBucketAlreadyOwnedByYouInvalidDurationBucketAlreadyExistsMetadataTooLargeUnsupportedMetadataUnsupportedHostHeaderMaximumExpiresSlowDownReadSlowDownWriteMaxVersionsExceededInvalidPrefixMarkerBadRequestKeyTooLongErrorInvalidBucketObjectLockConfigurationObjectLockConfigurationNotFoundObjectLockConfigurationNotAllowedNoSuchObjectLockConfigurationObjectLockedInvalidRetentionDatePastObjectLockRetainDateUnknownWORMModeDirectiveBucketTaggingNotFoundObjectLockInvalidHeadersInvalidTagDirectivePolicyAlreadyAttachedPolicyNotAttachedExcessDataInvalidEncryptionMethodInvalidEncryptionKeyIDInsecureSSECustomerRequestSSEMultipartEncryptedSSEEncryptedObjectInvalidEncryptionParametersInvalidEncryptionParametersSSECInvalidSSECustomerAlgorithmInvalidSSECustomerKeyMissingSSECustomerKeyMissingSSECustomerKeyMD5SSECustomerKeyMD5MismatchInvalidSSECustomerParametersIncompatibleEncryptionMethodKMSNotConfiguredKMSKeyNotFoundExceptionKMSDefaultKeyAlreadyConfiguredNoAccessKeyInvalidTokenEventNotificationARNNotificationRegionNotificationOverlappingFilterNotificationFilterNameInvalidFilterNamePrefixFilterNameSuffixFilterValueInvalidOverlappingConfigsUnsupportedNotificationContentSHA256MismatchContentChecksumMismatchStorageFullRequestBodyParseObjectExistsAsDirectoryInvalidObjectNameInvalidObjectNamePrefixSlashInvalidResourceNameInvalidLifecycleQueryParameterServerNotInitializedBucketMetadataNotInitializedRequestTimedoutClientDisconnectedTooManyRequestsInvalidRequestTransitionStorageClassNotFoundErrorInvalidStorageClassBackendDownMalformedJSONAdminNoSuchUserAdminNoSuchUserLDAPWarnAdminLDAPExpectedLoginNameAdminNoSuchGroupAdminGroupNotEmptyAdminGroupDisabledAdminInvalidGroupNameAdminNoSuchJobAdminNoSuchPolicyAdminPolicyChangeAlreadyAppliedAdminInvalidArgumentAdminInvalidAccessKeyAdminInvalidSecretKeyAdminConfigNoQuorumAdminConfigTooLargeAdminConfigBadJSONAdminNoSuchConfigTargetAdminConfigEnvOverriddenAdminConfigDuplicateKeysAdminConfigInvalidIDPTypeAdminConfigLDAPNonDefaultConfigNameAdminConfigLDAPValidationAdminConfigIDPCfgNameAlreadyExistsAdminConfigIDPCfgNameDoesNotExistInsecureClientRequestObjectTamperedAdminLDAPNotEnabledSiteReplicationInvalidRequestSiteReplicationPeerRespSiteReplicationBackendIssueSiteReplicationServiceAccountErrorSiteReplicationBucketConfigErrorSiteReplicationBucketMetaErrorSiteReplicationIAMErrorSiteReplicationConfigMissingSiteReplicationIAMConfigMismatchAdminRebalanceAlreadyStartedAdminRebalanceNotStartedAdminBucketQuotaExceededAdminNoSuchQuotaConfigurationHealNotImplementedHealNoSuchProcessHealInvalidClientTokenHealMissingBucketHealAlreadyRunningHealOverlappingPathsIncorrectContinuationTokenEmptyRequestBodyUnsupportedFunctionInvalidExpressionTypeBusyUnauthorizedAccessExpressionTooLongIllegalSQLFunctionArgumentInvalidKeyPathInvalidCompressionFormatInvalidFileHeaderInfoInvalidJSONTypeInvalidQuoteFieldsInvalidRequestParameterInvalidDataTypeInvalidTextEncodingInvalidDataSourceInvalidTableAliasMissingRequiredParameterObjectSerializationConflictUnsupportedSQLOperationUnsupportedSQLStructureUnsupportedSyntaxUnsupportedRangeHeaderLexerInvalidCharLexerInvalidOperatorLexerInvalidLiteralLexerInvalidIONLiteralParseExpectedDatePartParseExpectedKeywordParseExpectedTokenTypeParseExpected2TokenTypesParseExpectedNumberParseExpectedRightParenBuiltinFunctionCallParseExpectedTypeNameParseExpectedWhenClauseParseUnsupportedTokenParseUnsupportedLiteralsGroupByParseExpectedMemberParseUnsupportedSelectParseUnsupportedCaseParseUnsupportedCaseClauseParseUnsupportedAliasParseUnsupportedSyntaxParseUnknownOperatorParseMissingIdentAfterAtParseUnexpectedOperatorParseUnexpectedTermParseUnexpectedTokenParseUnexpectedKeywordParseExpectedExpressionParseExpectedLeftParenAfterCastParseExpectedLeftParenValueConstructorParseExpectedLeftParenBuiltinFunctionCallParseExpectedArgumentDelimiterParseCastArityParseInvalidTypeParamParseEmptySelectParseSelectMissingFromParseExpectedIdentForGroupNameParseExpectedIdentForAliasParseUnsupportedCallWithStarParseNonUnaryAggregateFunctionCallParseMalformedJoinParseExpectedIdentForAtParseAsteriskIsNotAloneInSelectListParseCannotMixSqbAndWildcardInSelectListParseInvalidContextForWildcardInSelectListIncorrectSQLFunctionArgumentTypeValueParseFailureEvaluatorInvalidArgumentsIntegerOverflowLikeInvalidInputsCastFailedInvalidCastEvaluatorInvalidTimestampFormatPatternEvaluatorInvalidTimestampFormatPatternSymbolForParsingEvaluatorTimestampFormatPatternDuplicateFieldsEvaluatorTimestampFormatPatternHourClockAmPmMismatchEvaluatorUnterminatedTimestampFormatPatternTokenEvaluatorInvalidTimestampFormatPatternTokenEvaluatorInvalidTimestampFormatPatternSymbolEvaluatorBindingDoesNotExistMissingHeadersInvalidColumnIndexAdminConfigNotificationTargetsFailedAdminProfilerNotEnabledInvalidDecompressedSizeAddUserInvalidArgumentAddUserValidUTFAdminResourceInvalidArgumentAdminAccountNotEligibleAccountNotEligibleAdminServiceAccountNotFoundPostPolicyConditionInvalidFormatInvalidChecksumLambdaARNInvalidLambdaARNNotFoundInvalidAttributeNameAdminNoAccessKeyAdminNoSecretKeyIAMNotInitializedapiErrCodeEnd"

var _APIErrorCode_index = [...]uint16{0, 4, 16, 25, 39, 53, 67, 81, 94, 112, 129, 150, 165, 182, 195, 207, 229, 249, 275, 289, 310, 327, 342, 365, 382, 400, 417, 441, 456, 477, 495, 507, 527, 544, 567, 588, 600, 618, 639, 667, 697, 718, 741, 767, 804, 834, 867, 892, 924, 954, 983, 1008, 1030, 1056, 1078, 1106, 1135, 1169, 1200, 1237, 1261, 1285, 1313, 1339, 1370, 1400, 1409, 1421, 1437, 1450, 1464, 1482, 1502, 1523, 1539, 1550, 1566, 1577, 1605, 1625, 1641, 1669, 1683, 1700, 1720, 1733, 1747, 1760, 1773, 1789, 1806, 1827, 1841, 1862, 1875, 1897, 1920, 1936, 1951, 1966, 1987, 2005, 2020, 2037, 2062, 2080, 2103, 2118, 2137, 2153, 2172, 2193, 2207, 2219, 2232, 2251, 2270, 2280, 2295, 2331, 2362, 2395, 2424, 2436, 2456, 2480, 2504, 2525, 2549, 2568, 2589, 2606, 2616, 2639, 2661, 2687, 2708, 2726, 2753, 2784, 2811, 2832, 2853, 2877, 2902, 2930, 2958, 2974, 2997, 3027, 3038, 3050, 3067, 3082, 3100, 3129, 3146, 3162, 3178, 3196, 3214, 3237, 3258, 3281, 3292, 3308, 3331, 3348, 3376, 3395, 3425, 3445, 3473, 3488, 3506, 3521, 3535, 3570, 3589, 3600, 3613, 3628, 3651, 3677, 3693, 3711, 3729, 3750, 3764, 3781, 3812, 3832, 3853, 3874, 3893, 3912, 3930, 3953, 3977, 4001, 4026, 4061, 4086, 4120, 4153, 4174, 4188, 4207, 4236, 4259, 4286, 4320, 4352, 4382, 4405, 4433, 4465, 4493, 4517, 4541, 4570, 4588, 4605, 4627, 4644, 4662, 4682, 4708, 4724, 4743, 4764, 4768, 4786, 4803, 4829, 4843, 4867, 4888, 4903, 4921, 4944, 4959, 4978, 4995, 5012, 5036, 5063, 5086, 5109, 5126, 5148, 5164, 5184, 5203, 5225, 5246, 5266, 5288, 5312, 5331, 5373, 5394, 5417, 5438, 5469, 5488, 5510, 5530, 5556, 5577, 5599, 5619, 5643, 5666, 5685, 5705, 5727, 5750, 5781, 5819, 5860, 5890, 5904, 5925, 5941, 5963, 5993, 6019, 6047, 6081, 6099, 6122, 6157, 6197, 6239, 6271, 6288, 6313, 6328, 6345, 6355, 6366, 6404, 6458, 6504, 6556, 6604, 6647, 6691, 6719, 6733, 6751, 6787, 6810, 6833, 6855, 6870, 6898, 6921, 6939, 6966, 6998, 7013, 7029, 7046, 7066, 7082, 7098, 7115, 7128}

func (i APIErrorCode) String() string {
	if i < 0 || i >= APIErrorCode(len(_APIErrorCode_index)-1) {
		return "APIErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _APIErrorCode_name[_APIErrorCode_index[i]:_APIErrorCode_index[i+1]]
}
