package cmd

import (
	"errors"
	"net/http"
	"testing"

	"reflect"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func ossErrResponse(code string) error {
	return traceError(oss.ServiceError{
		Code: code,
	})
}

func TestOSSToObjectError(t *testing.T) {
	testCases := []struct {
		inputErr       error
		expectedErr    error
		bucket, object string
	}{
		{
			inputErr:    ossErrResponse("BucketAlreadyExists"),
			expectedErr: BucketAlreadyOwnedByYou{},
		},
		{
			inputErr:    ossErrResponse("BucketNotEmpty"),
			expectedErr: BucketNotEmpty{},
		},
		{
			inputErr:    ossErrResponse("InvalidBucketName"),
			expectedErr: BucketNameInvalid{},
		},
		{
			inputErr:    ossErrResponse("NoSuchBucket"),
			expectedErr: BucketNotFound{},
		},
		// with empty object, NoSuchKey is interpreted as BucketNotFound
		{
			inputErr:    ossErrResponse("NoSuchKey"),
			expectedErr: BucketNotFound{},
		},
		{
			inputErr:    ossErrResponse("NoSuchUpload"),
			expectedErr: InvalidUploadID{},
		},
		{
			inputErr:    ossErrResponse("InvalidObjectName"),
			expectedErr: ObjectNameInvalid{},
		},
		{
			inputErr:    ossErrResponse("AccessDenied"),
			expectedErr: PrefixAccessDenied{},
		},
		{
			inputErr:    ossErrResponse("NoSuchUpload"),
			expectedErr: InvalidUploadID{},
		},
		{
			inputErr:    ossErrResponse("EntityTooSmall"),
			expectedErr: PartTooSmall{},
		},
		{
			inputErr:    nil,
			expectedErr: nil,
		},
		// Special test case for NoSuchKey with object name
		{
			inputErr:    ossErrResponse("NoSuchKey"),
			expectedErr: ObjectNotFound{Bucket: "bucket", Object: "object"},
			bucket:      "bucket",
			object:      "object",
		},

		// Special test case for error value that is not of
		// type (*Error)
		{
			inputErr:    errors.New("not a *Error"),
			expectedErr: errors.New("not a *Error"),
		},
	}

	for i, tc := range testCases {
		actualErr := ossToObjectError(tc.inputErr, tc.bucket, tc.object)
		if e, ok := actualErr.(*Error); ok && e.e != tc.expectedErr {
			t.Errorf("Test case %d: Expected error '%v' but received error '%v'", i+1, tc.expectedErr, e.e)
		}
	}
}

func TestS3MetaToOSSOptions(t *testing.T) {
	var err error
	var headers map[string]string

	headers = map[string]string{
		"invalid--meta": "value",
	}
	_, err = s3MetaToOSSOptions(headers)
	if err = errorCause(err); err != nil {
		if _, ok := err.(UnsupportedMetadata); !ok {
			t.Fatalf("Test failed with unexpected error %s, expected UnsupportedMetadata", err)
		}
	}

	headers = map[string]string{
		"accept-encoding":          "gzip", // not this
		"content-encoding":         "gzip",
		"X-Amz-Meta-Hdr":           "value",
		"X-Amz-Meta-X_test_key":    "value",
		"X-Amz-Meta-X__test__key":  "value",
		"X-Amz-Meta-X-Test__key":   "value",
		"X-Amz-Meta-X-Amz-Key":     "hu3ZSqtqwn+aL4V2VhAeov4i+bG3KyCtRMSXQFRHXOk=",
		"X-Amz-Meta-X-Amz-Matdesc": "{}",
		"X-Amz-Meta-X-Amz-Iv":      "eWmyryl8kq+EVnnsE7jpOg==",
	}
	opts, err := s3MetaToOSSOptions(headers)
	if err != nil {
		t.Fatalf("Test failed, with %s", err)
	}
	if len(opts) != len(headers)-1 {
		t.Fatalf("Test failed, S3 metdata is not fully transformed. expeted: %d, actual: %d", len(headers)-1, len(opts))
	}
}

func TestOSSHeaderToS3Meta(t *testing.T) {
	meta := map[string]string{
		"x-oss-meta-first_name":       "myname",
		"X-OSS-Meta-x_test_key":       "value",
		"X-Oss-Meta-x_test__key":      "value",
		"X-Oss-Meta-x__test__key":     "value",
		"X-Oss-Meta-x____test____key": "value",
		"X-Oss-Meta-x_amz_key":        "hu3ZSqtqwn+aL4V2VhAeov4i+bG3KyCtRMSXQFRHXOk=",
		"X-Oss-Meta-x_amz_matdesc":    "{}",
		"x-oss-meta-x_amz_iv":         "eWmyryl8kq+EVnnsE7jpOg==",
	}
	header := make(http.Header)
	for k, v := range meta {
		header.Set(k, v)
	}

	expectedMeta := map[string]string{
		"X-Amz-Meta-First-Name":    "myname",
		"X-Amz-Meta-X-Test-Key":    "value",
		"X-Amz-Meta-X-Test_key":    "value",
		"X-Amz-Meta-X_test_key":    "value",
		"X-Amz-Meta-X__test__key":  "value",
		"X-Amz-Meta-X-Amz-Key":     "hu3ZSqtqwn+aL4V2VhAeov4i+bG3KyCtRMSXQFRHXOk=",
		"X-Amz-Meta-X-Amz-Matdesc": "{}",
		"X-Amz-Meta-X-Amz-Iv":      "eWmyryl8kq+EVnnsE7jpOg==",
	}
	actualMeta := ossHeaderToS3Meta(header)
	for k, v := range expectedMeta {
		if v2, ok := actualMeta[k]; !ok {
			t.Errorf("Test failed for key %s: missing key", k)
		} else if v != v2 {
			t.Errorf("Test failed for key %s, expected '%s', got '%s'", k, v, v2)
		}
	}
}

func TestOSSBuildListObjectPartsParams(t *testing.T) {
	expected := map[string]interface{}{
		"uploadId":           "test",
		"part-number-marker": "123",
		"max-parts":          "456",
	}
	actual := ossBuildListObjectPartsParams("test", 123, 456)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Test failed, expected %v, got %v", expected, actual)
	}
}
