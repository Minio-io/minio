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

package csv

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/klauspost/compress/zip"
	"github.com/minio/minio/pkg/s3select/sql"
)

func TestRead(t *testing.T) {
	cases := []struct {
		content         string
		recordDelimiter string
		fieldDelimiter  string
	}{
		{"1,2,3\na,b,c\n", "\n", ","},
		{"1,2,3\ta,b,c\t", "\t", ","},
		{"1,2,3\r\na,b,c\r\n", "\r\n", ","},
	}

	for i, c := range cases {
		var err error
		var record sql.Record
		var result bytes.Buffer

		r, _ := NewReader(ioutil.NopCloser(strings.NewReader(c.content)), &ReaderArgs{
			FileHeaderInfo:             none,
			RecordDelimiter:            c.recordDelimiter,
			FieldDelimiter:             c.fieldDelimiter,
			QuoteCharacter:             defaultQuoteCharacter,
			QuoteEscapeCharacter:       defaultQuoteEscapeCharacter,
			CommentCharacter:           defaultCommentCharacter,
			AllowQuotedRecordDelimiter: true,
			unmarshaled:                true,
		})

		for {
			record, err = r.Read(record)
			if err != nil {
				break
			}
			record.WriteCSV(&result, []rune(c.fieldDelimiter)[0])
			result.Truncate(result.Len() - 1)
			result.WriteString(c.recordDelimiter)
		}
		r.Close()
		if err != io.EOF {
			t.Fatalf("Case %d failed with %s", i, err)
		}

		if result.String() != c.content {
			t.Errorf("Case %d failed: expected %v result %v", i, c.content, result.String())
		}
	}
}

type tester interface {
	Fatal(...interface{})
}

func openTestFile(t tester, file string) []byte {
	f, err := ioutil.ReadFile(filepath.Join("testdata/testdata.zip"))
	if err != nil {
		t.Fatal(err)
	}
	z, err := zip.NewReader(bytes.NewReader(f), int64(len(f)))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range z.File {
		if f.Name == file {
			rc, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			defer rc.Close()
			b, err := ioutil.ReadAll(rc)
			if err != nil {
				t.Fatal(err)
			}
			return b
		}
	}
	t.Fatal(file, "not found in testdata/testdata.zip")
	return nil
}

func TestReadExtended(t *testing.T) {
	cases := []struct {
		file            string
		recordDelimiter string
		fieldDelimiter  string
		wantColumns     []string
		wantTenFields   string
		totalFields     int
	}{
		{
			file:            "nyc-taxi-data-100k.csv",
			recordDelimiter: "\n",
			fieldDelimiter:  ",",
			wantColumns:     []string{"trip_id", "vendor_id", "pickup_datetime", "dropoff_datetime", "store_and_fwd_flag", "rate_code_id", "pickup_longitude", "pickup_latitude", "dropoff_longitude", "dropoff_latitude", "passenger_count", "trip_distance", "fare_amount", "extra", "mta_tax", "tip_amount", "tolls_amount", "ehail_fee", "improvement_surcharge", "total_amount", "payment_type", "trip_type", "pickup", "dropoff", "cab_type", "precipitation", "snow_depth", "snowfall", "max_temp", "min_temp", "wind", "pickup_nyct2010_gid", "pickup_ctlabel", "pickup_borocode", "pickup_boroname", "pickup_ct2010", "pickup_boroct2010", "pickup_cdeligibil", "pickup_ntacode", "pickup_ntaname", "pickup_puma", "dropoff_nyct2010_gid", "dropoff_ctlabel", "dropoff_borocode", "dropoff_boroname", "dropoff_ct2010", "dropoff_boroct2010", "dropoff_cdeligibil", "dropoff_ntacode", "dropoff_ntaname", "dropoff_puma"},
			wantTenFields: `3389224,2,2014-03-26 00:26:15,2014-03-26 00:28:38,N,1,-73.950431823730469,40.792251586914063,-73.938949584960937,40.794425964355469,1,0.84,4.5,0.5,0.5,1,0,,,6.5,1,1,75,74,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1828,180,1,Manhattan,018000,1018000,E,MN34,East Harlem North,3804
3389225,2,2014-03-31 09:42:15,2014-03-31 10:01:17,N,1,-73.950340270996094,40.792228698730469,-73.941970825195313,40.842235565185547,1,4.47,17.5,0,0.5,0,0,,,18,2,1,75,244,green,0.16,0.0,0.0,56,36,8.28,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,911,251,1,Manhattan,025100,1025100,E,MN36,Washington Heights South,3801
3389226,2,2014-03-26 17:13:28,2014-03-26 17:19:07,N,1,-73.949493408203125,40.793506622314453,-73.943374633789063,40.786155700683594,1,0.82,5.5,1,0.5,0,0,,,7,1,1,75,75,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1387,164,1,Manhattan,016400,1016400,E,MN33,East Harlem South,3804
3389227,2,2014-03-14 21:07:19,2014-03-14 21:11:41,N,1,-73.950538635253906,40.792228698730469,-73.940811157226563,40.809253692626953,1,1.40,6,0.5,0.5,0,0,,,7,2,1,75,42,green,0.00,0.0,0.0,46,22,5.59,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1184,208,1,Manhattan,020800,1020800,E,MN03,Central Harlem North-Polo Grounds,3803
3389228,1,2014-03-28 13:52:56,2014-03-28 14:29:01,N,1,-73.950569152832031,40.792312622070313,-73.868507385253906,40.688491821289063,2,16.10,46,0,0.5,0,5.33,,,51.83,2,,75,63,green,0.04,0.0,0.0,62,37,5.37,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1544,1182.02,3,Brooklyn,118202,3118202,E,BK83,Cypress Hills-City Line,4008
3389229,2,2014-03-07 09:46:32,2014-03-07 09:55:01,N,1,-73.952301025390625,40.789798736572266,-73.935806274414062,40.794448852539063,1,1.67,8,0,0.5,2,0,,,10.5,1,1,75,74,green,0.00,3.9,0.0,37,26,7.83,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1553,178,1,Manhattan,017800,1017800,E,MN34,East Harlem North,3804
3389230,2,2014-03-17 18:23:05,2014-03-17 18:28:38,N,1,-73.952346801757813,40.789844512939453,-73.946319580078125,40.783851623535156,5,0.95,5.5,1,0.5,0.65,0,,,7.65,1,1,75,263,green,0.00,0.0,0.0,35,23,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,32,156.01,1,Manhattan,015601,1015601,I,MN32,Yorkville,3805
3389231,1,2014-03-19 19:09:36,2014-03-19 19:12:20,N,1,-73.952377319335938,40.789779663085938,-73.947494506835938,40.796474456787109,1,0.50,4,1,0.5,1,0,,,6.5,1,,75,75,green,0.92,0.0,0.0,46,32,7.16,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1401,174.02,1,Manhattan,017402,1017402,E,MN33,East Harlem South,3804
3389232,2,2014-03-20 19:06:28,2014-03-20 19:21:35,N,1,-73.952583312988281,40.789516448974609,-73.985870361328125,40.776973724365234,2,3.04,13,1,0.5,2.8,0,,,17.3,1,1,75,143,green,0.00,0.0,0.0,54,40,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1742,155,1,Manhattan,015500,1015500,I,MN14,Lincoln Square,3806
3389233,2,2014-03-29 09:38:12,2014-03-29 09:44:16,N,1,-73.952728271484375,40.789501190185547,-73.950935363769531,40.775600433349609,1,1.10,6.5,0,0.5,1.3,0,,,8.3,1,1,75,263,green,1.81,0.0,0.0,59,43,10.74,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,2048,138,1,Manhattan,013800,1013800,I,MN32,Yorkville,3805
`,
			totalFields: 308,
		}, {
			file:            "nyc-taxi-data-tabs-100k.csv",
			recordDelimiter: "\n",
			fieldDelimiter:  "\t",
			wantColumns:     []string{"trip_id", "vendor_id", "pickup_datetime", "dropoff_datetime", "store_and_fwd_flag", "rate_code_id", "pickup_longitude", "pickup_latitude", "dropoff_longitude", "dropoff_latitude", "passenger_count", "trip_distance", "fare_amount", "extra", "mta_tax", "tip_amount", "tolls_amount", "ehail_fee", "improvement_surcharge", "total_amount", "payment_type", "trip_type", "pickup", "dropoff", "cab_type", "precipitation", "snow_depth", "snowfall", "max_temp", "min_temp", "wind", "pickup_nyct2010_gid", "pickup_ctlabel", "pickup_borocode", "pickup_boroname", "pickup_ct2010", "pickup_boroct2010", "pickup_cdeligibil", "pickup_ntacode", "pickup_ntaname", "pickup_puma", "dropoff_nyct2010_gid", "dropoff_ctlabel", "dropoff_borocode", "dropoff_boroname", "dropoff_ct2010", "dropoff_boroct2010", "dropoff_cdeligibil", "dropoff_ntacode", "dropoff_ntaname", "dropoff_puma"},
			wantTenFields: `3389224,2,2014-03-26 00:26:15,2014-03-26 00:28:38,N,1,-73.950431823730469,40.792251586914063,-73.938949584960937,40.794425964355469,1,0.84,4.5,0.5,0.5,1,0,,,6.5,1,1,75,74,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1828,180,1,Manhattan,018000,1018000,E,MN34,East Harlem North,3804
3389225,2,2014-03-31 09:42:15,2014-03-31 10:01:17,N,1,-73.950340270996094,40.792228698730469,-73.941970825195313,40.842235565185547,1,4.47,17.5,0,0.5,0,0,,,18,2,1,75,244,green,0.16,0.0,0.0,56,36,8.28,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,911,251,1,Manhattan,025100,1025100,E,MN36,Washington Heights South,3801
3389226,2,2014-03-26 17:13:28,2014-03-26 17:19:07,N,1,-73.949493408203125,40.793506622314453,-73.943374633789063,40.786155700683594,1,0.82,5.5,1,0.5,0,0,,,7,1,1,75,75,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1387,164,1,Manhattan,016400,1016400,E,MN33,East Harlem South,3804
3389227,2,2014-03-14 21:07:19,2014-03-14 21:11:41,N,1,-73.950538635253906,40.792228698730469,-73.940811157226563,40.809253692626953,1,1.40,6,0.5,0.5,0,0,,,7,2,1,75,42,green,0.00,0.0,0.0,46,22,5.59,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1184,208,1,Manhattan,020800,1020800,E,MN03,Central Harlem North-Polo Grounds,3803
3389228,1,2014-03-28 13:52:56,2014-03-28 14:29:01,N,1,-73.950569152832031,40.792312622070313,-73.868507385253906,40.688491821289063,2,16.10,46,0,0.5,0,5.33,,,51.83,2,,75,63,green,0.04,0.0,0.0,62,37,5.37,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1544,1182.02,3,Brooklyn,118202,3118202,E,BK83,Cypress Hills-City Line,4008
3389229,2,2014-03-07 09:46:32,2014-03-07 09:55:01,N,1,-73.952301025390625,40.789798736572266,-73.935806274414062,40.794448852539063,1,1.67,8,0,0.5,2,0,,,10.5,1,1,75,74,green,0.00,3.9,0.0,37,26,7.83,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1553,178,1,Manhattan,017800,1017800,E,MN34,East Harlem North,3804
3389230,2,2014-03-17 18:23:05,2014-03-17 18:28:38,N,1,-73.952346801757813,40.789844512939453,-73.946319580078125,40.783851623535156,5,0.95,5.5,1,0.5,0.65,0,,,7.65,1,1,75,263,green,0.00,0.0,0.0,35,23,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,32,156.01,1,Manhattan,015601,1015601,I,MN32,Yorkville,3805
3389231,1,2014-03-19 19:09:36,2014-03-19 19:12:20,N,1,-73.952377319335938,40.789779663085938,-73.947494506835938,40.796474456787109,1,0.50,4,1,0.5,1,0,,,6.5,1,,75,75,green,0.92,0.0,0.0,46,32,7.16,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1401,174.02,1,Manhattan,017402,1017402,E,MN33,East Harlem South,3804
3389232,2,2014-03-20 19:06:28,2014-03-20 19:21:35,N,1,-73.952583312988281,40.789516448974609,-73.985870361328125,40.776973724365234,2,3.04,13,1,0.5,2.8,0,,,17.3,1,1,75,143,green,0.00,0.0,0.0,54,40,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1742,155,1,Manhattan,015500,1015500,I,MN14,Lincoln Square,3806
3389233,2,2014-03-29 09:38:12,2014-03-29 09:44:16,N,1,-73.952728271484375,40.789501190185547,-73.950935363769531,40.775600433349609,1,1.10,6.5,0,0.5,1.3,0,,,8.3,1,1,75,263,green,1.81,0.0,0.0,59,43,10.74,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,2048,138,1,Manhattan,013800,1013800,I,MN32,Yorkville,3805
`,
			totalFields: 308,
		}, {
			file:            "nyc-taxi-data-100k-single-delim.csv",
			recordDelimiter: "^",
			fieldDelimiter:  ",",
			wantColumns:     []string{"trip_id", "vendor_id", "pickup_datetime", "dropoff_datetime", "store_and_fwd_flag", "rate_code_id", "pickup_longitude", "pickup_latitude", "dropoff_longitude", "dropoff_latitude", "passenger_count", "trip_distance", "fare_amount", "extra", "mta_tax", "tip_amount", "tolls_amount", "ehail_fee", "improvement_surcharge", "total_amount", "payment_type", "trip_type", "pickup", "dropoff", "cab_type", "precipitation", "snow_depth", "snowfall", "max_temp", "min_temp", "wind", "pickup_nyct2010_gid", "pickup_ctlabel", "pickup_borocode", "pickup_boroname", "pickup_ct2010", "pickup_boroct2010", "pickup_cdeligibil", "pickup_ntacode", "pickup_ntaname", "pickup_puma", "dropoff_nyct2010_gid", "dropoff_ctlabel", "dropoff_borocode", "dropoff_boroname", "dropoff_ct2010", "dropoff_boroct2010", "dropoff_cdeligibil", "dropoff_ntacode", "dropoff_ntaname", "dropoff_puma"},
			wantTenFields: `3389224,2,2014-03-26 00:26:15,2014-03-26 00:28:38,N,1,-73.950431823730469,40.792251586914063,-73.938949584960937,40.794425964355469,1,0.84,4.5,0.5,0.5,1,0,,,6.5,1,1,75,74,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1828,180,1,Manhattan,018000,1018000,E,MN34,East Harlem North,3804
3389225,2,2014-03-31 09:42:15,2014-03-31 10:01:17,N,1,-73.950340270996094,40.792228698730469,-73.941970825195313,40.842235565185547,1,4.47,17.5,0,0.5,0,0,,,18,2,1,75,244,green,0.16,0.0,0.0,56,36,8.28,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,911,251,1,Manhattan,025100,1025100,E,MN36,Washington Heights South,3801
3389226,2,2014-03-26 17:13:28,2014-03-26 17:19:07,N,1,-73.949493408203125,40.793506622314453,-73.943374633789063,40.786155700683594,1,0.82,5.5,1,0.5,0,0,,,7,1,1,75,75,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1387,164,1,Manhattan,016400,1016400,E,MN33,East Harlem South,3804
3389227,2,2014-03-14 21:07:19,2014-03-14 21:11:41,N,1,-73.950538635253906,40.792228698730469,-73.940811157226563,40.809253692626953,1,1.40,6,0.5,0.5,0,0,,,7,2,1,75,42,green,0.00,0.0,0.0,46,22,5.59,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1184,208,1,Manhattan,020800,1020800,E,MN03,Central Harlem North-Polo Grounds,3803
3389228,1,2014-03-28 13:52:56,2014-03-28 14:29:01,N,1,-73.950569152832031,40.792312622070313,-73.868507385253906,40.688491821289063,2,16.10,46,0,0.5,0,5.33,,,51.83,2,,75,63,green,0.04,0.0,0.0,62,37,5.37,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1544,1182.02,3,Brooklyn,118202,3118202,E,BK83,Cypress Hills-City Line,4008
3389229,2,2014-03-07 09:46:32,2014-03-07 09:55:01,N,1,-73.952301025390625,40.789798736572266,-73.935806274414062,40.794448852539063,1,1.67,8,0,0.5,2,0,,,10.5,1,1,75,74,green,0.00,3.9,0.0,37,26,7.83,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1553,178,1,Manhattan,017800,1017800,E,MN34,East Harlem North,3804
3389230,2,2014-03-17 18:23:05,2014-03-17 18:28:38,N,1,-73.952346801757813,40.789844512939453,-73.946319580078125,40.783851623535156,5,0.95,5.5,1,0.5,0.65,0,,,7.65,1,1,75,263,green,0.00,0.0,0.0,35,23,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,32,156.01,1,Manhattan,015601,1015601,I,MN32,Yorkville,3805
3389231,1,2014-03-19 19:09:36,2014-03-19 19:12:20,N,1,-73.952377319335938,40.789779663085938,-73.947494506835938,40.796474456787109,1,0.50,4,1,0.5,1,0,,,6.5,1,,75,75,green,0.92,0.0,0.0,46,32,7.16,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1401,174.02,1,Manhattan,017402,1017402,E,MN33,East Harlem South,3804
3389232,2,2014-03-20 19:06:28,2014-03-20 19:21:35,N,1,-73.952583312988281,40.789516448974609,-73.985870361328125,40.776973724365234,2,3.04,13,1,0.5,2.8,0,,,17.3,1,1,75,143,green,0.00,0.0,0.0,54,40,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1742,155,1,Manhattan,015500,1015500,I,MN14,Lincoln Square,3806
3389233,2,2014-03-29 09:38:12,2014-03-29 09:44:16,N,1,-73.952728271484375,40.789501190185547,-73.950935363769531,40.775600433349609,1,1.10,6.5,0,0.5,1.3,0,,,8.3,1,1,75,263,green,1.81,0.0,0.0,59,43,10.74,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,2048,138,1,Manhattan,013800,1013800,I,MN32,Yorkville,3805
`,
			totalFields: 308,
		}, {
			file:            "nyc-taxi-data-100k-multi-delim.csv",
			recordDelimiter: "^Y",
			fieldDelimiter:  ",",
			wantColumns:     []string{"trip_id", "vendor_id", "pickup_datetime", "dropoff_datetime", "store_and_fwd_flag", "rate_code_id", "pickup_longitude", "pickup_latitude", "dropoff_longitude", "dropoff_latitude", "passenger_count", "trip_distance", "fare_amount", "extra", "mta_tax", "tip_amount", "tolls_amount", "ehail_fee", "improvement_surcharge", "total_amount", "payment_type", "trip_type", "pickup", "dropoff", "cab_type", "precipitation", "snow_depth", "snowfall", "max_temp", "min_temp", "wind", "pickup_nyct2010_gid", "pickup_ctlabel", "pickup_borocode", "pickup_boroname", "pickup_ct2010", "pickup_boroct2010", "pickup_cdeligibil", "pickup_ntacode", "pickup_ntaname", "pickup_puma", "dropoff_nyct2010_gid", "dropoff_ctlabel", "dropoff_borocode", "dropoff_boroname", "dropoff_ct2010", "dropoff_boroct2010", "dropoff_cdeligibil", "dropoff_ntacode", "dropoff_ntaname", "dropoff_puma"},
			wantTenFields: `3389224,2,2014-03-26 00:26:15,2014-03-26 00:28:38,N,1,-73.950431823730469,40.792251586914063,-73.938949584960937,40.794425964355469,1,0.84,4.5,0.5,0.5,1,0,,,6.5,1,1,75,74,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1828,180,1,Manhattan,018000,1018000,E,MN34,East Harlem North,3804
3389225,2,2014-03-31 09:42:15,2014-03-31 10:01:17,N,1,-73.950340270996094,40.792228698730469,-73.941970825195313,40.842235565185547,1,4.47,17.5,0,0.5,0,0,,,18,2,1,75,244,green,0.16,0.0,0.0,56,36,8.28,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,911,251,1,Manhattan,025100,1025100,E,MN36,Washington Heights South,3801
3389226,2,2014-03-26 17:13:28,2014-03-26 17:19:07,N,1,-73.949493408203125,40.793506622314453,-73.943374633789063,40.786155700683594,1,0.82,5.5,1,0.5,0,0,,,7,1,1,75,75,green,0.00,0.0,0.0,36,24,11.86,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1387,164,1,Manhattan,016400,1016400,E,MN33,East Harlem South,3804
3389227,2,2014-03-14 21:07:19,2014-03-14 21:11:41,N,1,-73.950538635253906,40.792228698730469,-73.940811157226563,40.809253692626953,1,1.40,6,0.5,0.5,0,0,,,7,2,1,75,42,green,0.00,0.0,0.0,46,22,5.59,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1184,208,1,Manhattan,020800,1020800,E,MN03,Central Harlem North-Polo Grounds,3803
3389228,1,2014-03-28 13:52:56,2014-03-28 14:29:01,N,1,-73.950569152832031,40.792312622070313,-73.868507385253906,40.688491821289063,2,16.10,46,0,0.5,0,5.33,,,51.83,2,,75,63,green,0.04,0.0,0.0,62,37,5.37,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1544,1182.02,3,Brooklyn,118202,3118202,E,BK83,Cypress Hills-City Line,4008
3389229,2,2014-03-07 09:46:32,2014-03-07 09:55:01,N,1,-73.952301025390625,40.789798736572266,-73.935806274414062,40.794448852539063,1,1.67,8,0,0.5,2,0,,,10.5,1,1,75,74,green,0.00,3.9,0.0,37,26,7.83,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1553,178,1,Manhattan,017800,1017800,E,MN34,East Harlem North,3804
3389230,2,2014-03-17 18:23:05,2014-03-17 18:28:38,N,1,-73.952346801757813,40.789844512939453,-73.946319580078125,40.783851623535156,5,0.95,5.5,1,0.5,0.65,0,,,7.65,1,1,75,263,green,0.00,0.0,0.0,35,23,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,32,156.01,1,Manhattan,015601,1015601,I,MN32,Yorkville,3805
3389231,1,2014-03-19 19:09:36,2014-03-19 19:12:20,N,1,-73.952377319335938,40.789779663085938,-73.947494506835938,40.796474456787109,1,0.50,4,1,0.5,1,0,,,6.5,1,,75,75,green,0.92,0.0,0.0,46,32,7.16,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1401,174.02,1,Manhattan,017402,1017402,E,MN33,East Harlem South,3804
3389232,2,2014-03-20 19:06:28,2014-03-20 19:21:35,N,1,-73.952583312988281,40.789516448974609,-73.985870361328125,40.776973724365234,2,3.04,13,1,0.5,2.8,0,,,17.3,1,1,75,143,green,0.00,0.0,0.0,54,40,8.05,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,1742,155,1,Manhattan,015500,1015500,I,MN14,Lincoln Square,3806
3389233,2,2014-03-29 09:38:12,2014-03-29 09:44:16,N,1,-73.952728271484375,40.789501190185547,-73.950935363769531,40.775600433349609,1,1.10,6.5,0,0.5,1.3,0,,,8.3,1,1,75,263,green,1.81,0.0,0.0,59,43,10.74,1267,168,1,Manhattan,016800,1016800,E,MN33,East Harlem South,3804,2048,138,1,Manhattan,013800,1013800,I,MN32,Yorkville,3805
`,
			totalFields: 308,
		},
	}

	for i, c := range cases {
		t.Run(c.file, func(t *testing.T) {

			var err error
			var record sql.Record
			var result bytes.Buffer

			r, _ := NewReader(ioutil.NopCloser(bytes.NewReader(openTestFile(t, c.file))), &ReaderArgs{
				FileHeaderInfo:             use,
				RecordDelimiter:            c.recordDelimiter,
				FieldDelimiter:             c.fieldDelimiter,
				QuoteCharacter:             defaultQuoteCharacter,
				QuoteEscapeCharacter:       defaultQuoteEscapeCharacter,
				CommentCharacter:           defaultCommentCharacter,
				AllowQuotedRecordDelimiter: true,
				unmarshaled:                true,
			})
			fields := 0
			for {
				record, err = r.Read(record)
				if err != nil {
					break
				}
				if fields < 10 {
					// Write with fixed delimiters, newlines.
					err := record.WriteCSV(&result, ',')
					if err != nil {
						t.Error(err)
					}
				}
				fields++
			}
			r.Close()
			if err != io.EOF {
				t.Fatalf("Case %d failed with %s", i, err)
			}
			if !reflect.DeepEqual(r.columnNames, c.wantColumns) {
				t.Errorf("Case %d failed: expected %#v, got result %#v", i, c.wantColumns, r.columnNames)
			}
			if result.String() != c.wantTenFields {
				t.Errorf("Case %d failed: expected %v, got result %v", i, c.wantTenFields, result.String())
			}
			if fields != c.totalFields {
				t.Errorf("Case %d failed: expected %v results %v", i, c.totalFields, fields)
			}
		})
	}
}

func BenchmarkReaderBasic(b *testing.B) {
	const dataDir = "testdata"
	args := ReaderArgs{
		FileHeaderInfo:             use,
		RecordDelimiter:            "\n",
		FieldDelimiter:             ",",
		QuoteCharacter:             defaultQuoteCharacter,
		QuoteEscapeCharacter:       defaultQuoteEscapeCharacter,
		CommentCharacter:           defaultCommentCharacter,
		AllowQuotedRecordDelimiter: true,
		unmarshaled:                true,
	}
	f := openTestFile(b, "nyc-taxi-data-100k.csv")
	r, err := NewReader(ioutil.NopCloser(bytes.NewBuffer(f)), &args)
	if err != nil {
		b.Fatalf("Reading init failed with %s", err)
	}
	defer r.Close()
	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(len(f)))
	var record sql.Record
	for i := 0; i < b.N; i++ {
		r, err = NewReader(ioutil.NopCloser(bytes.NewBuffer(f)), &args)
		if err != nil {
			b.Fatalf("Reading init failed with %s", err)
		}

		for err == nil {
			record, err = r.Read(record)
			if err != nil && err != io.EOF {
				b.Fatalf("Reading failed with %s", err)
			}
		}
		r.Close()
	}
}

func BenchmarkReaderReplace(b *testing.B) {
	const dataDir = "testdata"
	args := ReaderArgs{
		FileHeaderInfo:             use,
		RecordDelimiter:            "^",
		FieldDelimiter:             ",",
		QuoteCharacter:             defaultQuoteCharacter,
		QuoteEscapeCharacter:       defaultQuoteEscapeCharacter,
		CommentCharacter:           defaultCommentCharacter,
		AllowQuotedRecordDelimiter: true,
		unmarshaled:                true,
	}
	f := openTestFile(b, "nyc-taxi-data-100k-single-delim.csv")
	r, err := NewReader(ioutil.NopCloser(bytes.NewBuffer(f)), &args)
	if err != nil {
		b.Fatalf("Reading init failed with %s", err)
	}
	defer r.Close()
	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(len(f)))
	var record sql.Record
	for i := 0; i < b.N; i++ {
		r, err = NewReader(ioutil.NopCloser(bytes.NewBuffer(f)), &args)
		if err != nil {
			b.Fatalf("Reading init failed with %s", err)
		}

		for err == nil {
			record, err = r.Read(record)
			if err != nil && err != io.EOF {
				b.Fatalf("Reading failed with %s", err)
			}
		}
		r.Close()
	}
}

func BenchmarkReaderReplaceTwo(b *testing.B) {
	const dataDir = "testdata"
	args := ReaderArgs{
		FileHeaderInfo:             use,
		RecordDelimiter:            "^Y",
		FieldDelimiter:             ",",
		QuoteCharacter:             defaultQuoteCharacter,
		QuoteEscapeCharacter:       defaultQuoteEscapeCharacter,
		CommentCharacter:           defaultCommentCharacter,
		AllowQuotedRecordDelimiter: true,
		unmarshaled:                true,
	}
	f := openTestFile(b, "nyc-taxi-data-100k-multi-delim.csv")
	r, err := NewReader(ioutil.NopCloser(bytes.NewBuffer(f)), &args)
	if err != nil {
		b.Fatalf("Reading init failed with %s", err)
	}
	defer r.Close()
	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(len(f)))
	var record sql.Record
	for i := 0; i < b.N; i++ {
		r, err = NewReader(ioutil.NopCloser(bytes.NewBuffer(f)), &args)
		if err != nil {
			b.Fatalf("Reading init failed with %s", err)
		}

		for err == nil {
			record, err = r.Read(record)
			if err != nil && err != io.EOF {
				b.Fatalf("Reading failed with %s", err)
			}
		}
		r.Close()
	}
}
