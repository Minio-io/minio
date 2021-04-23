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

package parquet

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

//Types supported by Parquet.  These types are intended to be used in combination
//with the encodings to control the on disk storage format.
//For example INT16 is not included as a type since a good encoding of INT32
//would handle this.
type Type int64

const (
	Type_BOOLEAN              Type = 0
	Type_INT32                Type = 1
	Type_INT64                Type = 2
	Type_INT96                Type = 3
	Type_FLOAT                Type = 4
	Type_DOUBLE               Type = 5
	Type_BYTE_ARRAY           Type = 6
	Type_FIXED_LEN_BYTE_ARRAY Type = 7
)

func (p Type) String() string {
	switch p {
	case Type_BOOLEAN:
		return "BOOLEAN"
	case Type_INT32:
		return "INT32"
	case Type_INT64:
		return "INT64"
	case Type_INT96:
		return "INT96"
	case Type_FLOAT:
		return "FLOAT"
	case Type_DOUBLE:
		return "DOUBLE"
	case Type_BYTE_ARRAY:
		return "BYTE_ARRAY"
	case Type_FIXED_LEN_BYTE_ARRAY:
		return "FIXED_LEN_BYTE_ARRAY"
	}
	return "<UNSET>"
}

func TypeFromString(s string) (Type, error) {
	switch s {
	case "BOOLEAN":
		return Type_BOOLEAN, nil
	case "INT32":
		return Type_INT32, nil
	case "INT64":
		return Type_INT64, nil
	case "INT96":
		return Type_INT96, nil
	case "FLOAT":
		return Type_FLOAT, nil
	case "DOUBLE":
		return Type_DOUBLE, nil
	case "BYTE_ARRAY":
		return Type_BYTE_ARRAY, nil
	case "FIXED_LEN_BYTE_ARRAY":
		return Type_FIXED_LEN_BYTE_ARRAY, nil
	}
	return Type(0), fmt.Errorf("not a valid Type string")
}

func TypePtr(v Type) *Type { return &v }

func (p Type) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Type) UnmarshalText(text []byte) error {
	q, err := TypeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *Type) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = Type(v)
	return nil
}

func (p *Type) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

//Common types used by frameworks(e.g. hive, pig) using parquet.  This helps map
//between types in those frameworks to the base types in parquet.  This is only
//metadata and not needed to read or write the data.
type ConvertedType int64

const (
	ConvertedType_UTF8             ConvertedType = 0
	ConvertedType_MAP              ConvertedType = 1
	ConvertedType_MAP_KEY_VALUE    ConvertedType = 2
	ConvertedType_LIST             ConvertedType = 3
	ConvertedType_ENUM             ConvertedType = 4
	ConvertedType_DECIMAL          ConvertedType = 5
	ConvertedType_DATE             ConvertedType = 6
	ConvertedType_TIME_MILLIS      ConvertedType = 7
	ConvertedType_TIME_MICROS      ConvertedType = 8
	ConvertedType_TIMESTAMP_MILLIS ConvertedType = 9
	ConvertedType_TIMESTAMP_MICROS ConvertedType = 10
	ConvertedType_UINT_8           ConvertedType = 11
	ConvertedType_UINT_16          ConvertedType = 12
	ConvertedType_UINT_32          ConvertedType = 13
	ConvertedType_UINT_64          ConvertedType = 14
	ConvertedType_INT_8            ConvertedType = 15
	ConvertedType_INT_16           ConvertedType = 16
	ConvertedType_INT_32           ConvertedType = 17
	ConvertedType_INT_64           ConvertedType = 18
	ConvertedType_JSON             ConvertedType = 19
	ConvertedType_BSON             ConvertedType = 20
	ConvertedType_INTERVAL         ConvertedType = 21
)

func (p ConvertedType) String() string {
	switch p {
	case ConvertedType_UTF8:
		return "UTF8"
	case ConvertedType_MAP:
		return "MAP"
	case ConvertedType_MAP_KEY_VALUE:
		return "MAP_KEY_VALUE"
	case ConvertedType_LIST:
		return "LIST"
	case ConvertedType_ENUM:
		return "ENUM"
	case ConvertedType_DECIMAL:
		return "DECIMAL"
	case ConvertedType_DATE:
		return "DATE"
	case ConvertedType_TIME_MILLIS:
		return "TIME_MILLIS"
	case ConvertedType_TIME_MICROS:
		return "TIME_MICROS"
	case ConvertedType_TIMESTAMP_MILLIS:
		return "TIMESTAMP_MILLIS"
	case ConvertedType_TIMESTAMP_MICROS:
		return "TIMESTAMP_MICROS"
	case ConvertedType_UINT_8:
		return "UINT_8"
	case ConvertedType_UINT_16:
		return "UINT_16"
	case ConvertedType_UINT_32:
		return "UINT_32"
	case ConvertedType_UINT_64:
		return "UINT_64"
	case ConvertedType_INT_8:
		return "INT_8"
	case ConvertedType_INT_16:
		return "INT_16"
	case ConvertedType_INT_32:
		return "INT_32"
	case ConvertedType_INT_64:
		return "INT_64"
	case ConvertedType_JSON:
		return "JSON"
	case ConvertedType_BSON:
		return "BSON"
	case ConvertedType_INTERVAL:
		return "INTERVAL"
	}
	return "<UNSET>"
}

func ConvertedTypeFromString(s string) (ConvertedType, error) {
	switch s {
	case "UTF8":
		return ConvertedType_UTF8, nil
	case "MAP":
		return ConvertedType_MAP, nil
	case "MAP_KEY_VALUE":
		return ConvertedType_MAP_KEY_VALUE, nil
	case "LIST":
		return ConvertedType_LIST, nil
	case "ENUM":
		return ConvertedType_ENUM, nil
	case "DECIMAL":
		return ConvertedType_DECIMAL, nil
	case "DATE":
		return ConvertedType_DATE, nil
	case "TIME_MILLIS":
		return ConvertedType_TIME_MILLIS, nil
	case "TIME_MICROS":
		return ConvertedType_TIME_MICROS, nil
	case "TIMESTAMP_MILLIS":
		return ConvertedType_TIMESTAMP_MILLIS, nil
	case "TIMESTAMP_MICROS":
		return ConvertedType_TIMESTAMP_MICROS, nil
	case "UINT_8":
		return ConvertedType_UINT_8, nil
	case "UINT_16":
		return ConvertedType_UINT_16, nil
	case "UINT_32":
		return ConvertedType_UINT_32, nil
	case "UINT_64":
		return ConvertedType_UINT_64, nil
	case "INT_8":
		return ConvertedType_INT_8, nil
	case "INT_16":
		return ConvertedType_INT_16, nil
	case "INT_32":
		return ConvertedType_INT_32, nil
	case "INT_64":
		return ConvertedType_INT_64, nil
	case "JSON":
		return ConvertedType_JSON, nil
	case "BSON":
		return ConvertedType_BSON, nil
	case "INTERVAL":
		return ConvertedType_INTERVAL, nil
	}
	return ConvertedType(0), fmt.Errorf("not a valid ConvertedType string")
}

func ConvertedTypePtr(v ConvertedType) *ConvertedType { return &v }

func (p ConvertedType) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *ConvertedType) UnmarshalText(text []byte) error {
	q, err := ConvertedTypeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *ConvertedType) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = ConvertedType(v)
	return nil
}

func (p *ConvertedType) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

//Representation of Schemas
type FieldRepetitionType int64

const (
	FieldRepetitionType_REQUIRED FieldRepetitionType = 0
	FieldRepetitionType_OPTIONAL FieldRepetitionType = 1
	FieldRepetitionType_REPEATED FieldRepetitionType = 2
)

func (p FieldRepetitionType) String() string {
	switch p {
	case FieldRepetitionType_REQUIRED:
		return "REQUIRED"
	case FieldRepetitionType_OPTIONAL:
		return "OPTIONAL"
	case FieldRepetitionType_REPEATED:
		return "REPEATED"
	}
	return "<UNSET>"
}

func FieldRepetitionTypeFromString(s string) (FieldRepetitionType, error) {
	switch s {
	case "REQUIRED":
		return FieldRepetitionType_REQUIRED, nil
	case "OPTIONAL":
		return FieldRepetitionType_OPTIONAL, nil
	case "REPEATED":
		return FieldRepetitionType_REPEATED, nil
	}
	return FieldRepetitionType(0), fmt.Errorf("not a valid FieldRepetitionType string")
}

func FieldRepetitionTypePtr(v FieldRepetitionType) *FieldRepetitionType { return &v }

func (p FieldRepetitionType) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *FieldRepetitionType) UnmarshalText(text []byte) error {
	q, err := FieldRepetitionTypeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *FieldRepetitionType) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = FieldRepetitionType(v)
	return nil
}

func (p *FieldRepetitionType) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

//Encodings supported by Parquet.  Not all encodings are valid for all types.  These
//enums are also used to specify the encoding of definition and repetition levels.
//See the accompanying doc for the details of the more complicated encodings.
type Encoding int64

const (
	Encoding_PLAIN                   Encoding = 0
	Encoding_PLAIN_DICTIONARY        Encoding = 2
	Encoding_RLE                     Encoding = 3
	Encoding_BIT_PACKED              Encoding = 4
	Encoding_DELTA_BINARY_PACKED     Encoding = 5
	Encoding_DELTA_LENGTH_BYTE_ARRAY Encoding = 6
	Encoding_DELTA_BYTE_ARRAY        Encoding = 7
	Encoding_RLE_DICTIONARY          Encoding = 8
)

func (p Encoding) String() string {
	switch p {
	case Encoding_PLAIN:
		return "PLAIN"
	case Encoding_PLAIN_DICTIONARY:
		return "PLAIN_DICTIONARY"
	case Encoding_RLE:
		return "RLE"
	case Encoding_BIT_PACKED:
		return "BIT_PACKED"
	case Encoding_DELTA_BINARY_PACKED:
		return "DELTA_BINARY_PACKED"
	case Encoding_DELTA_LENGTH_BYTE_ARRAY:
		return "DELTA_LENGTH_BYTE_ARRAY"
	case Encoding_DELTA_BYTE_ARRAY:
		return "DELTA_BYTE_ARRAY"
	case Encoding_RLE_DICTIONARY:
		return "RLE_DICTIONARY"
	}
	return "<UNSET>"
}

func EncodingFromString(s string) (Encoding, error) {
	switch s {
	case "PLAIN":
		return Encoding_PLAIN, nil
	case "PLAIN_DICTIONARY":
		return Encoding_PLAIN_DICTIONARY, nil
	case "RLE":
		return Encoding_RLE, nil
	case "BIT_PACKED":
		return Encoding_BIT_PACKED, nil
	case "DELTA_BINARY_PACKED":
		return Encoding_DELTA_BINARY_PACKED, nil
	case "DELTA_LENGTH_BYTE_ARRAY":
		return Encoding_DELTA_LENGTH_BYTE_ARRAY, nil
	case "DELTA_BYTE_ARRAY":
		return Encoding_DELTA_BYTE_ARRAY, nil
	case "RLE_DICTIONARY":
		return Encoding_RLE_DICTIONARY, nil
	}
	return Encoding(0), fmt.Errorf("not a valid Encoding string")
}

func EncodingPtr(v Encoding) *Encoding { return &v }

func (p Encoding) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Encoding) UnmarshalText(text []byte) error {
	q, err := EncodingFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *Encoding) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = Encoding(v)
	return nil
}

func (p *Encoding) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

//Supported compression algorithms.
//
//Codecs added in 2.4 can be read by readers based on 2.4 and later.
//Codec support may vary between readers based on the format version and
//libraries available at runtime. Gzip, Snappy, and LZ4 codecs are
//widely available, while Zstd and Brotli require additional libraries.
type CompressionCodec int64

const (
	CompressionCodec_UNCOMPRESSED CompressionCodec = 0
	CompressionCodec_SNAPPY       CompressionCodec = 1
	CompressionCodec_GZIP         CompressionCodec = 2
	CompressionCodec_LZO          CompressionCodec = 3
	CompressionCodec_BROTLI       CompressionCodec = 4
	CompressionCodec_LZ4          CompressionCodec = 5
	CompressionCodec_ZSTD         CompressionCodec = 6
)

func (p CompressionCodec) String() string {
	switch p {
	case CompressionCodec_UNCOMPRESSED:
		return "UNCOMPRESSED"
	case CompressionCodec_SNAPPY:
		return "SNAPPY"
	case CompressionCodec_GZIP:
		return "GZIP"
	case CompressionCodec_LZO:
		return "LZO"
	case CompressionCodec_BROTLI:
		return "BROTLI"
	case CompressionCodec_LZ4:
		return "LZ4"
	case CompressionCodec_ZSTD:
		return "ZSTD"
	}
	return "<UNSET>"
}

func CompressionCodecFromString(s string) (CompressionCodec, error) {
	switch s {
	case "UNCOMPRESSED":
		return CompressionCodec_UNCOMPRESSED, nil
	case "SNAPPY":
		return CompressionCodec_SNAPPY, nil
	case "GZIP":
		return CompressionCodec_GZIP, nil
	case "LZO":
		return CompressionCodec_LZO, nil
	case "BROTLI":
		return CompressionCodec_BROTLI, nil
	case "LZ4":
		return CompressionCodec_LZ4, nil
	case "ZSTD":
		return CompressionCodec_ZSTD, nil
	}
	return CompressionCodec(0), fmt.Errorf("not a valid CompressionCodec string")
}

func CompressionCodecPtr(v CompressionCodec) *CompressionCodec { return &v }

func (p CompressionCodec) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *CompressionCodec) UnmarshalText(text []byte) error {
	q, err := CompressionCodecFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *CompressionCodec) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = CompressionCodec(v)
	return nil
}

func (p *CompressionCodec) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

type PageType int64

const (
	PageType_DATA_PAGE       PageType = 0
	PageType_INDEX_PAGE      PageType = 1
	PageType_DICTIONARY_PAGE PageType = 2
	PageType_DATA_PAGE_V2    PageType = 3
)

func (p PageType) String() string {
	switch p {
	case PageType_DATA_PAGE:
		return "DATA_PAGE"
	case PageType_INDEX_PAGE:
		return "INDEX_PAGE"
	case PageType_DICTIONARY_PAGE:
		return "DICTIONARY_PAGE"
	case PageType_DATA_PAGE_V2:
		return "DATA_PAGE_V2"
	}
	return "<UNSET>"
}

func PageTypeFromString(s string) (PageType, error) {
	switch s {
	case "DATA_PAGE":
		return PageType_DATA_PAGE, nil
	case "INDEX_PAGE":
		return PageType_INDEX_PAGE, nil
	case "DICTIONARY_PAGE":
		return PageType_DICTIONARY_PAGE, nil
	case "DATA_PAGE_V2":
		return PageType_DATA_PAGE_V2, nil
	}
	return PageType(0), fmt.Errorf("not a valid PageType string")
}

func PageTypePtr(v PageType) *PageType { return &v }

func (p PageType) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *PageType) UnmarshalText(text []byte) error {
	q, err := PageTypeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *PageType) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = PageType(v)
	return nil
}

func (p *PageType) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

//Enum to annotate whether lists of min/max elements inside ColumnIndex
//are ordered and if so, in which direction.
type BoundaryOrder int64

const (
	BoundaryOrder_UNORDERED  BoundaryOrder = 0
	BoundaryOrder_ASCENDING  BoundaryOrder = 1
	BoundaryOrder_DESCENDING BoundaryOrder = 2
)

func (p BoundaryOrder) String() string {
	switch p {
	case BoundaryOrder_UNORDERED:
		return "UNORDERED"
	case BoundaryOrder_ASCENDING:
		return "ASCENDING"
	case BoundaryOrder_DESCENDING:
		return "DESCENDING"
	}
	return "<UNSET>"
}

func BoundaryOrderFromString(s string) (BoundaryOrder, error) {
	switch s {
	case "UNORDERED":
		return BoundaryOrder_UNORDERED, nil
	case "ASCENDING":
		return BoundaryOrder_ASCENDING, nil
	case "DESCENDING":
		return BoundaryOrder_DESCENDING, nil
	}
	return BoundaryOrder(0), fmt.Errorf("not a valid BoundaryOrder string")
}

func BoundaryOrderPtr(v BoundaryOrder) *BoundaryOrder { return &v }

func (p BoundaryOrder) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *BoundaryOrder) UnmarshalText(text []byte) error {
	q, err := BoundaryOrderFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *BoundaryOrder) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = BoundaryOrder(v)
	return nil
}

func (p *BoundaryOrder) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Statistics per row group and per page
// All fields are optional.
//
// Attributes:
//  - Max: DEPRECATED: min and max value of the column. Use min_value and max_value.
//
// Values are encoded using PLAIN encoding, except that variable-length byte
// arrays do not include a length prefix.
//
// These fields encode min and max values determined by signed comparison
// only. New files should use the correct order for a column's logical type
// and store the values in the min_value and max_value fields.
//
// To support older readers, these may be set when the column order is
// signed.
//  - Min
//  - NullCount: count of null value in the column
//  - DistinctCount: count of distinct values occurring
//  - MaxValue: Min and max values for the column, determined by its ColumnOrder.
//
// Values are encoded using PLAIN encoding, except that variable-length byte
// arrays do not include a length prefix.
//  - MinValue
type Statistics struct {
	Max           []byte `thrift:"max,1" db:"max" json:"max,omitempty"`
	Min           []byte `thrift:"min,2" db:"min" json:"min,omitempty"`
	NullCount     *int64 `thrift:"null_count,3" db:"null_count" json:"null_count,omitempty"`
	DistinctCount *int64 `thrift:"distinct_count,4" db:"distinct_count" json:"distinct_count,omitempty"`
	MaxValue      []byte `thrift:"max_value,5" db:"max_value" json:"max_value,omitempty"`
	MinValue      []byte `thrift:"min_value,6" db:"min_value" json:"min_value,omitempty"`
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

var Statistics_Max_DEFAULT []byte

func (p *Statistics) GetMax() []byte {
	return p.Max
}

var Statistics_Min_DEFAULT []byte

func (p *Statistics) GetMin() []byte {
	return p.Min
}

var Statistics_NullCount_DEFAULT int64

func (p *Statistics) GetNullCount() int64 {
	if !p.IsSetNullCount() {
		return Statistics_NullCount_DEFAULT
	}
	return *p.NullCount
}

var Statistics_DistinctCount_DEFAULT int64

func (p *Statistics) GetDistinctCount() int64 {
	if !p.IsSetDistinctCount() {
		return Statistics_DistinctCount_DEFAULT
	}
	return *p.DistinctCount
}

var Statistics_MaxValue_DEFAULT []byte

func (p *Statistics) GetMaxValue() []byte {
	return p.MaxValue
}

var Statistics_MinValue_DEFAULT []byte

func (p *Statistics) GetMinValue() []byte {
	return p.MinValue
}
func (p *Statistics) IsSetMax() bool {
	return p.Max != nil
}

func (p *Statistics) IsSetMin() bool {
	return p.Min != nil
}

func (p *Statistics) IsSetNullCount() bool {
	return p.NullCount != nil
}

func (p *Statistics) IsSetDistinctCount() bool {
	return p.DistinctCount != nil
}

func (p *Statistics) IsSetMaxValue() bool {
	return p.MaxValue != nil
}

func (p *Statistics) IsSetMinValue() bool {
	return p.MinValue != nil
}

func (p *Statistics) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Statistics) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Max = v
	}
	return nil
}

func (p *Statistics) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Min = v
	}
	return nil
}

func (p *Statistics) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.NullCount = &v
	}
	return nil
}

func (p *Statistics) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.DistinctCount = &v
	}
	return nil
}

func (p *Statistics) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.MaxValue = v
	}
	return nil
}

func (p *Statistics) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.MinValue = v
	}
	return nil
}

func (p *Statistics) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Statistics"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Statistics) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetMax() {
		if err := oprot.WriteFieldBegin("max", thrift.STRING, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:max: ", p), err)
		}
		if err := oprot.WriteBinary(p.Max); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.max (1) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:max: ", p), err)
		}
	}
	return err
}

func (p *Statistics) writeField2(oprot thrift.TProtocol) (err error) {
	if p.IsSetMin() {
		if err := oprot.WriteFieldBegin("min", thrift.STRING, 2); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:min: ", p), err)
		}
		if err := oprot.WriteBinary(p.Min); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.min (2) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 2:min: ", p), err)
		}
	}
	return err
}

func (p *Statistics) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetNullCount() {
		if err := oprot.WriteFieldBegin("null_count", thrift.I64, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:null_count: ", p), err)
		}
		if err := oprot.WriteI64(*p.NullCount); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.null_count (3) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:null_count: ", p), err)
		}
	}
	return err
}

func (p *Statistics) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetDistinctCount() {
		if err := oprot.WriteFieldBegin("distinct_count", thrift.I64, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:distinct_count: ", p), err)
		}
		if err := oprot.WriteI64(*p.DistinctCount); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.distinct_count (4) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:distinct_count: ", p), err)
		}
	}
	return err
}

func (p *Statistics) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetMaxValue() {
		if err := oprot.WriteFieldBegin("max_value", thrift.STRING, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:max_value: ", p), err)
		}
		if err := oprot.WriteBinary(p.MaxValue); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.max_value (5) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:max_value: ", p), err)
		}
	}
	return err
}

func (p *Statistics) writeField6(oprot thrift.TProtocol) (err error) {
	if p.IsSetMinValue() {
		if err := oprot.WriteFieldBegin("min_value", thrift.STRING, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:min_value: ", p), err)
		}
		if err := oprot.WriteBinary(p.MinValue); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.min_value (6) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:min_value: ", p), err)
		}
	}
	return err
}

func (p *Statistics) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Statistics(%+v)", *p)
}

// Empty structs to use as logical type annotations
type StringType struct {
}

func NewStringType() *StringType {
	return &StringType{}
}

func (p *StringType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StringType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("StringType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StringType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StringType(%+v)", *p)
}

type UUIDType struct {
}

func NewUUIDType() *UUIDType {
	return &UUIDType{}
}

func (p *UUIDType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *UUIDType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("UUIDType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *UUIDType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UUIDType(%+v)", *p)
}

type MapType struct {
}

func NewMapType() *MapType {
	return &MapType{}
}

func (p *MapType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *MapType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MapType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *MapType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MapType(%+v)", *p)
}

type ListType struct {
}

func NewListType() *ListType {
	return &ListType{}
}

func (p *ListType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *ListType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ListType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ListType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ListType(%+v)", *p)
}

type EnumType struct {
}

func NewEnumType() *EnumType {
	return &EnumType{}
}

func (p *EnumType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *EnumType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("EnumType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *EnumType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("EnumType(%+v)", *p)
}

type DateType struct {
}

func NewDateType() *DateType {
	return &DateType{}
}

func (p *DateType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *DateType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("DateType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *DateType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DateType(%+v)", *p)
}

// Logical type to annotate a column that is always null.
//
// Sometimes when discovering the schema of existing data, values are always
// null and the physical type can't be determined. This annotation signals
// the case where the physical type was guessed from all null values.
type NullType struct {
}

func NewNullType() *NullType {
	return &NullType{}
}

func (p *NullType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *NullType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("NullType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *NullType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("NullType(%+v)", *p)
}

// Decimal logical type annotation
//
// To maintain forward-compatibility in v1, implementations using this logical
// type must also set scale and precision on the annotated SchemaElement.
//
// Allowed for physical types: INT32, INT64, FIXED, and BINARY
//
// Attributes:
//  - Scale
//  - Precision
type DecimalType struct {
	Scale     int32 `thrift:"scale,1,required" db:"scale" json:"scale"`
	Precision int32 `thrift:"precision,2,required" db:"precision" json:"precision"`
}

func NewDecimalType() *DecimalType {
	return &DecimalType{}
}

func (p *DecimalType) GetScale() int32 {
	return p.Scale
}

func (p *DecimalType) GetPrecision() int32 {
	return p.Precision
}
func (p *DecimalType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetScale bool = false
	var issetPrecision bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetScale = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetPrecision = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetScale {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Scale is not set"))
	}
	if !issetPrecision {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Precision is not set"))
	}
	return nil
}

func (p *DecimalType) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Scale = v
	}
	return nil
}

func (p *DecimalType) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Precision = v
	}
	return nil
}

func (p *DecimalType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("DecimalType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *DecimalType) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("scale", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:scale: ", p), err)
	}
	if err := oprot.WriteI32(p.Scale); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.scale (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:scale: ", p), err)
	}
	return err
}

func (p *DecimalType) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("precision", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:precision: ", p), err)
	}
	if err := oprot.WriteI32(p.Precision); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.precision (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:precision: ", p), err)
	}
	return err
}

func (p *DecimalType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DecimalType(%+v)", *p)
}

// Time units for logical types
type MilliSeconds struct {
}

func NewMilliSeconds() *MilliSeconds {
	return &MilliSeconds{}
}

func (p *MilliSeconds) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *MilliSeconds) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MilliSeconds"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *MilliSeconds) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MilliSeconds(%+v)", *p)
}

type MicroSeconds struct {
}

func NewMicroSeconds() *MicroSeconds {
	return &MicroSeconds{}
}

func (p *MicroSeconds) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *MicroSeconds) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MicroSeconds"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *MicroSeconds) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MicroSeconds(%+v)", *p)
}

type NanoSeconds struct {
}

func NewNanoSeconds() *NanoSeconds {
	return &NanoSeconds{}
}

func (p *NanoSeconds) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *NanoSeconds) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("NanoSeconds"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *NanoSeconds) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("NanoSeconds(%+v)", *p)
}

// Attributes:
//  - MILLIS
//  - MICROS
//  - NANOS
type TimeUnit struct {
	MILLIS *MilliSeconds `thrift:"MILLIS,1" db:"MILLIS" json:"MILLIS,omitempty"`
	MICROS *MicroSeconds `thrift:"MICROS,2" db:"MICROS" json:"MICROS,omitempty"`
	NANOS  *NanoSeconds  `thrift:"NANOS,3" db:"NANOS" json:"NANOS,omitempty"`
}

func NewTimeUnit() *TimeUnit {
	return &TimeUnit{}
}

var TimeUnit_MILLIS_DEFAULT *MilliSeconds

func (p *TimeUnit) GetMILLIS() *MilliSeconds {
	if !p.IsSetMILLIS() {
		return TimeUnit_MILLIS_DEFAULT
	}
	return p.MILLIS
}

var TimeUnit_MICROS_DEFAULT *MicroSeconds

func (p *TimeUnit) GetMICROS() *MicroSeconds {
	if !p.IsSetMICROS() {
		return TimeUnit_MICROS_DEFAULT
	}
	return p.MICROS
}

var TimeUnit_NANOS_DEFAULT *NanoSeconds

func (p *TimeUnit) GetNANOS() *NanoSeconds {
	if !p.IsSetNANOS() {
		return TimeUnit_NANOS_DEFAULT
	}
	return p.NANOS
}
func (p *TimeUnit) CountSetFieldsTimeUnit() int {
	count := 0
	if p.IsSetMILLIS() {
		count++
	}
	if p.IsSetMICROS() {
		count++
	}
	if p.IsSetNANOS() {
		count++
	}
	return count

}

func (p *TimeUnit) IsSetMILLIS() bool {
	return p.MILLIS != nil
}

func (p *TimeUnit) IsSetMICROS() bool {
	return p.MICROS != nil
}

func (p *TimeUnit) IsSetNANOS() bool {
	return p.NANOS != nil
}

func (p *TimeUnit) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TimeUnit) ReadField1(iprot thrift.TProtocol) error {
	p.MILLIS = &MilliSeconds{}
	if err := p.MILLIS.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.MILLIS), err)
	}
	return nil
}

func (p *TimeUnit) ReadField2(iprot thrift.TProtocol) error {
	p.MICROS = &MicroSeconds{}
	if err := p.MICROS.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.MICROS), err)
	}
	return nil
}

func (p *TimeUnit) ReadField3(iprot thrift.TProtocol) error {
	p.NANOS = &NanoSeconds{}
	if err := p.NANOS.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.NANOS), err)
	}
	return nil
}

func (p *TimeUnit) Write(oprot thrift.TProtocol) error {
	if c := p.CountSetFieldsTimeUnit(); c != 1 {
		return fmt.Errorf("%T write union: exactly one field must be set (%d set).", p, c)
	}
	if err := oprot.WriteStructBegin("TimeUnit"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TimeUnit) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetMILLIS() {
		if err := oprot.WriteFieldBegin("MILLIS", thrift.STRUCT, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:MILLIS: ", p), err)
		}
		if err := p.MILLIS.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.MILLIS), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:MILLIS: ", p), err)
		}
	}
	return err
}

func (p *TimeUnit) writeField2(oprot thrift.TProtocol) (err error) {
	if p.IsSetMICROS() {
		if err := oprot.WriteFieldBegin("MICROS", thrift.STRUCT, 2); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:MICROS: ", p), err)
		}
		if err := p.MICROS.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.MICROS), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 2:MICROS: ", p), err)
		}
	}
	return err
}

func (p *TimeUnit) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetNANOS() {
		if err := oprot.WriteFieldBegin("NANOS", thrift.STRUCT, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:NANOS: ", p), err)
		}
		if err := p.NANOS.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.NANOS), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:NANOS: ", p), err)
		}
	}
	return err
}

func (p *TimeUnit) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TimeUnit(%+v)", *p)
}

// Timestamp logical type annotation
//
// Allowed for physical types: INT64
//
// Attributes:
//  - IsAdjustedToUTC
//  - Unit
type TimestampType struct {
	IsAdjustedToUTC bool      `thrift:"isAdjustedToUTC,1,required" db:"isAdjustedToUTC" json:"isAdjustedToUTC"`
	Unit            *TimeUnit `thrift:"unit,2,required" db:"unit" json:"unit"`
}

func NewTimestampType() *TimestampType {
	return &TimestampType{}
}

func (p *TimestampType) GetIsAdjustedToUTC() bool {
	return p.IsAdjustedToUTC
}

var TimestampType_Unit_DEFAULT *TimeUnit

func (p *TimestampType) GetUnit() *TimeUnit {
	if !p.IsSetUnit() {
		return TimestampType_Unit_DEFAULT
	}
	return p.Unit
}
func (p *TimestampType) IsSetUnit() bool {
	return p.Unit != nil
}

func (p *TimestampType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetIsAdjustedToUTC bool = false
	var issetUnit bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetIsAdjustedToUTC = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetUnit = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetIsAdjustedToUTC {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsAdjustedToUTC is not set"))
	}
	if !issetUnit {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Unit is not set"))
	}
	return nil
}

func (p *TimestampType) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.IsAdjustedToUTC = v
	}
	return nil
}

func (p *TimestampType) ReadField2(iprot thrift.TProtocol) error {
	p.Unit = &TimeUnit{}
	if err := p.Unit.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Unit), err)
	}
	return nil
}

func (p *TimestampType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("TimestampType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TimestampType) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("isAdjustedToUTC", thrift.BOOL, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:isAdjustedToUTC: ", p), err)
	}
	if err := oprot.WriteBool(p.IsAdjustedToUTC); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.isAdjustedToUTC (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:isAdjustedToUTC: ", p), err)
	}
	return err
}

func (p *TimestampType) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("unit", thrift.STRUCT, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:unit: ", p), err)
	}
	if err := p.Unit.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Unit), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:unit: ", p), err)
	}
	return err
}

func (p *TimestampType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TimestampType(%+v)", *p)
}

// Time logical type annotation
//
// Allowed for physical types: INT32 (millis), INT64 (micros, nanos)
//
// Attributes:
//  - IsAdjustedToUTC
//  - Unit
type TimeType struct {
	IsAdjustedToUTC bool      `thrift:"isAdjustedToUTC,1,required" db:"isAdjustedToUTC" json:"isAdjustedToUTC"`
	Unit            *TimeUnit `thrift:"unit,2,required" db:"unit" json:"unit"`
}

func NewTimeType() *TimeType {
	return &TimeType{}
}

func (p *TimeType) GetIsAdjustedToUTC() bool {
	return p.IsAdjustedToUTC
}

var TimeType_Unit_DEFAULT *TimeUnit

func (p *TimeType) GetUnit() *TimeUnit {
	if !p.IsSetUnit() {
		return TimeType_Unit_DEFAULT
	}
	return p.Unit
}
func (p *TimeType) IsSetUnit() bool {
	return p.Unit != nil
}

func (p *TimeType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetIsAdjustedToUTC bool = false
	var issetUnit bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetIsAdjustedToUTC = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetUnit = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetIsAdjustedToUTC {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsAdjustedToUTC is not set"))
	}
	if !issetUnit {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Unit is not set"))
	}
	return nil
}

func (p *TimeType) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.IsAdjustedToUTC = v
	}
	return nil
}

func (p *TimeType) ReadField2(iprot thrift.TProtocol) error {
	p.Unit = &TimeUnit{}
	if err := p.Unit.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Unit), err)
	}
	return nil
}

func (p *TimeType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("TimeType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TimeType) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("isAdjustedToUTC", thrift.BOOL, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:isAdjustedToUTC: ", p), err)
	}
	if err := oprot.WriteBool(p.IsAdjustedToUTC); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.isAdjustedToUTC (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:isAdjustedToUTC: ", p), err)
	}
	return err
}

func (p *TimeType) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("unit", thrift.STRUCT, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:unit: ", p), err)
	}
	if err := p.Unit.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Unit), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:unit: ", p), err)
	}
	return err
}

func (p *TimeType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TimeType(%+v)", *p)
}

// Integer logical type annotation
//
// bitWidth must be 8, 16, 32, or 64.
//
// Allowed for physical types: INT32, INT64
//
// Attributes:
//  - BitWidth
//  - IsSigned
type IntType struct {
	BitWidth int8 `thrift:"bitWidth,1,required" db:"bitWidth" json:"bitWidth"`
	IsSigned bool `thrift:"isSigned,2,required" db:"isSigned" json:"isSigned"`
}

func NewIntType() *IntType {
	return &IntType{}
}

func (p *IntType) GetBitWidth() int8 {
	return p.BitWidth
}

func (p *IntType) GetIsSigned() bool {
	return p.IsSigned
}
func (p *IntType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetBitWidth bool = false
	var issetIsSigned bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetBitWidth = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetIsSigned = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetBitWidth {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field BitWidth is not set"))
	}
	if !issetIsSigned {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsSigned is not set"))
	}
	return nil
}

func (p *IntType) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadByte(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := v
		p.BitWidth = temp
	}
	return nil
}

func (p *IntType) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.IsSigned = v
	}
	return nil
}

func (p *IntType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("IntType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *IntType) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("bitWidth", thrift.BYTE, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:bitWidth: ", p), err)
	}
	if err := oprot.WriteByte(p.BitWidth); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.bitWidth (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:bitWidth: ", p), err)
	}
	return err
}

func (p *IntType) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("isSigned", thrift.BOOL, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:isSigned: ", p), err)
	}
	if err := oprot.WriteBool(p.IsSigned); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.isSigned (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:isSigned: ", p), err)
	}
	return err
}

func (p *IntType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("IntType(%+v)", *p)
}

// Embedded JSON logical type annotation
//
// Allowed for physical types: BINARY
type JsonType struct {
}

func NewJsonType() *JsonType {
	return &JsonType{}
}

func (p *JsonType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *JsonType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("JsonType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *JsonType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("JsonType(%+v)", *p)
}

// Embedded BSON logical type annotation
//
// Allowed for physical types: BINARY
type BsonType struct {
}

func NewBsonType() *BsonType {
	return &BsonType{}
}

func (p *BsonType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *BsonType) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("BsonType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *BsonType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BsonType(%+v)", *p)
}

// LogicalType annotations to replace ConvertedType.
//
// To maintain compatibility, implementations using LogicalType for a
// SchemaElement must also set the corresponding ConvertedType from the
// following table.
//
// Attributes:
//  - STRING
//  - MAP
//  - LIST
//  - ENUM
//  - DECIMAL
//  - DATE
//  - TIME
//  - TIMESTAMP
//  - INTEGER
//  - UNKNOWN
//  - JSON
//  - BSON
//  - UUID
type LogicalType struct {
	STRING    *StringType    `thrift:"STRING,1" db:"STRING" json:"STRING,omitempty"`
	MAP       *MapType       `thrift:"MAP,2" db:"MAP" json:"MAP,omitempty"`
	LIST      *ListType      `thrift:"LIST,3" db:"LIST" json:"LIST,omitempty"`
	ENUM      *EnumType      `thrift:"ENUM,4" db:"ENUM" json:"ENUM,omitempty"`
	DECIMAL   *DecimalType   `thrift:"DECIMAL,5" db:"DECIMAL" json:"DECIMAL,omitempty"`
	DATE      *DateType      `thrift:"DATE,6" db:"DATE" json:"DATE,omitempty"`
	TIME      *TimeType      `thrift:"TIME,7" db:"TIME" json:"TIME,omitempty"`
	TIMESTAMP *TimestampType `thrift:"TIMESTAMP,8" db:"TIMESTAMP" json:"TIMESTAMP,omitempty"`
	// unused field # 9
	INTEGER *IntType  `thrift:"INTEGER,10" db:"INTEGER" json:"INTEGER,omitempty"`
	UNKNOWN *NullType `thrift:"UNKNOWN,11" db:"UNKNOWN" json:"UNKNOWN,omitempty"`
	JSON    *JsonType `thrift:"JSON,12" db:"JSON" json:"JSON,omitempty"`
	BSON    *BsonType `thrift:"BSON,13" db:"BSON" json:"BSON,omitempty"`
	UUID    *UUIDType `thrift:"UUID,14" db:"UUID" json:"UUID,omitempty"`
}

func NewLogicalType() *LogicalType {
	return &LogicalType{}
}

var LogicalType_STRING_DEFAULT *StringType

func (p *LogicalType) GetSTRING() *StringType {
	if !p.IsSetSTRING() {
		return LogicalType_STRING_DEFAULT
	}
	return p.STRING
}

var LogicalType_MAP_DEFAULT *MapType

func (p *LogicalType) GetMAP() *MapType {
	if !p.IsSetMAP() {
		return LogicalType_MAP_DEFAULT
	}
	return p.MAP
}

var LogicalType_LIST_DEFAULT *ListType

func (p *LogicalType) GetLIST() *ListType {
	if !p.IsSetLIST() {
		return LogicalType_LIST_DEFAULT
	}
	return p.LIST
}

var LogicalType_ENUM_DEFAULT *EnumType

func (p *LogicalType) GetENUM() *EnumType {
	if !p.IsSetENUM() {
		return LogicalType_ENUM_DEFAULT
	}
	return p.ENUM
}

var LogicalType_DECIMAL_DEFAULT *DecimalType

func (p *LogicalType) GetDECIMAL() *DecimalType {
	if !p.IsSetDECIMAL() {
		return LogicalType_DECIMAL_DEFAULT
	}
	return p.DECIMAL
}

var LogicalType_DATE_DEFAULT *DateType

func (p *LogicalType) GetDATE() *DateType {
	if !p.IsSetDATE() {
		return LogicalType_DATE_DEFAULT
	}
	return p.DATE
}

var LogicalType_TIME_DEFAULT *TimeType

func (p *LogicalType) GetTIME() *TimeType {
	if !p.IsSetTIME() {
		return LogicalType_TIME_DEFAULT
	}
	return p.TIME
}

var LogicalType_TIMESTAMP_DEFAULT *TimestampType

func (p *LogicalType) GetTIMESTAMP() *TimestampType {
	if !p.IsSetTIMESTAMP() {
		return LogicalType_TIMESTAMP_DEFAULT
	}
	return p.TIMESTAMP
}

var LogicalType_INTEGER_DEFAULT *IntType

func (p *LogicalType) GetINTEGER() *IntType {
	if !p.IsSetINTEGER() {
		return LogicalType_INTEGER_DEFAULT
	}
	return p.INTEGER
}

var LogicalType_UNKNOWN_DEFAULT *NullType

func (p *LogicalType) GetUNKNOWN() *NullType {
	if !p.IsSetUNKNOWN() {
		return LogicalType_UNKNOWN_DEFAULT
	}
	return p.UNKNOWN
}

var LogicalType_JSON_DEFAULT *JsonType

func (p *LogicalType) GetJSON() *JsonType {
	if !p.IsSetJSON() {
		return LogicalType_JSON_DEFAULT
	}
	return p.JSON
}

var LogicalType_BSON_DEFAULT *BsonType

func (p *LogicalType) GetBSON() *BsonType {
	if !p.IsSetBSON() {
		return LogicalType_BSON_DEFAULT
	}
	return p.BSON
}

var LogicalType_UUID_DEFAULT *UUIDType

func (p *LogicalType) GetUUID() *UUIDType {
	if !p.IsSetUUID() {
		return LogicalType_UUID_DEFAULT
	}
	return p.UUID
}
func (p *LogicalType) CountSetFieldsLogicalType() int {
	count := 0
	if p.IsSetSTRING() {
		count++
	}
	if p.IsSetMAP() {
		count++
	}
	if p.IsSetLIST() {
		count++
	}
	if p.IsSetENUM() {
		count++
	}
	if p.IsSetDECIMAL() {
		count++
	}
	if p.IsSetDATE() {
		count++
	}
	if p.IsSetTIME() {
		count++
	}
	if p.IsSetTIMESTAMP() {
		count++
	}
	if p.IsSetINTEGER() {
		count++
	}
	if p.IsSetUNKNOWN() {
		count++
	}
	if p.IsSetJSON() {
		count++
	}
	if p.IsSetBSON() {
		count++
	}
	if p.IsSetUUID() {
		count++
	}
	return count

}

func (p *LogicalType) IsSetSTRING() bool {
	return p.STRING != nil
}

func (p *LogicalType) IsSetMAP() bool {
	return p.MAP != nil
}

func (p *LogicalType) IsSetLIST() bool {
	return p.LIST != nil
}

func (p *LogicalType) IsSetENUM() bool {
	return p.ENUM != nil
}

func (p *LogicalType) IsSetDECIMAL() bool {
	return p.DECIMAL != nil
}

func (p *LogicalType) IsSetDATE() bool {
	return p.DATE != nil
}

func (p *LogicalType) IsSetTIME() bool {
	return p.TIME != nil
}

func (p *LogicalType) IsSetTIMESTAMP() bool {
	return p.TIMESTAMP != nil
}

func (p *LogicalType) IsSetINTEGER() bool {
	return p.INTEGER != nil
}

func (p *LogicalType) IsSetUNKNOWN() bool {
	return p.UNKNOWN != nil
}

func (p *LogicalType) IsSetJSON() bool {
	return p.JSON != nil
}

func (p *LogicalType) IsSetBSON() bool {
	return p.BSON != nil
}

func (p *LogicalType) IsSetUUID() bool {
	return p.UUID != nil
}

func (p *LogicalType) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.ReadField8(iprot); err != nil {
				return err
			}
		case 10:
			if err := p.ReadField10(iprot); err != nil {
				return err
			}
		case 11:
			if err := p.ReadField11(iprot); err != nil {
				return err
			}
		case 12:
			if err := p.ReadField12(iprot); err != nil {
				return err
			}
		case 13:
			if err := p.ReadField13(iprot); err != nil {
				return err
			}
		case 14:
			if err := p.ReadField14(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *LogicalType) ReadField1(iprot thrift.TProtocol) error {
	p.STRING = &StringType{}
	if err := p.STRING.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.STRING), err)
	}
	return nil
}

func (p *LogicalType) ReadField2(iprot thrift.TProtocol) error {
	p.MAP = &MapType{}
	if err := p.MAP.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.MAP), err)
	}
	return nil
}

func (p *LogicalType) ReadField3(iprot thrift.TProtocol) error {
	p.LIST = &ListType{}
	if err := p.LIST.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.LIST), err)
	}
	return nil
}

func (p *LogicalType) ReadField4(iprot thrift.TProtocol) error {
	p.ENUM = &EnumType{}
	if err := p.ENUM.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.ENUM), err)
	}
	return nil
}

func (p *LogicalType) ReadField5(iprot thrift.TProtocol) error {
	p.DECIMAL = &DecimalType{}
	if err := p.DECIMAL.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.DECIMAL), err)
	}
	return nil
}

func (p *LogicalType) ReadField6(iprot thrift.TProtocol) error {
	p.DATE = &DateType{}
	if err := p.DATE.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.DATE), err)
	}
	return nil
}

func (p *LogicalType) ReadField7(iprot thrift.TProtocol) error {
	p.TIME = &TimeType{}
	if err := p.TIME.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.TIME), err)
	}
	return nil
}

func (p *LogicalType) ReadField8(iprot thrift.TProtocol) error {
	p.TIMESTAMP = &TimestampType{}
	if err := p.TIMESTAMP.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.TIMESTAMP), err)
	}
	return nil
}

func (p *LogicalType) ReadField10(iprot thrift.TProtocol) error {
	p.INTEGER = &IntType{}
	if err := p.INTEGER.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.INTEGER), err)
	}
	return nil
}

func (p *LogicalType) ReadField11(iprot thrift.TProtocol) error {
	p.UNKNOWN = &NullType{}
	if err := p.UNKNOWN.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.UNKNOWN), err)
	}
	return nil
}

func (p *LogicalType) ReadField12(iprot thrift.TProtocol) error {
	p.JSON = &JsonType{}
	if err := p.JSON.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.JSON), err)
	}
	return nil
}

func (p *LogicalType) ReadField13(iprot thrift.TProtocol) error {
	p.BSON = &BsonType{}
	if err := p.BSON.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.BSON), err)
	}
	return nil
}

func (p *LogicalType) ReadField14(iprot thrift.TProtocol) error {
	p.UUID = &UUIDType{}
	if err := p.UUID.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.UUID), err)
	}
	return nil
}

func (p *LogicalType) Write(oprot thrift.TProtocol) error {
	if c := p.CountSetFieldsLogicalType(); c != 1 {
		return fmt.Errorf("%T write union: exactly one field must be set (%d set).", p, c)
	}
	if err := oprot.WriteStructBegin("LogicalType"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
		if err := p.writeField8(oprot); err != nil {
			return err
		}
		if err := p.writeField10(oprot); err != nil {
			return err
		}
		if err := p.writeField11(oprot); err != nil {
			return err
		}
		if err := p.writeField12(oprot); err != nil {
			return err
		}
		if err := p.writeField13(oprot); err != nil {
			return err
		}
		if err := p.writeField14(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *LogicalType) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetSTRING() {
		if err := oprot.WriteFieldBegin("STRING", thrift.STRUCT, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:STRING: ", p), err)
		}
		if err := p.STRING.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.STRING), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:STRING: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField2(oprot thrift.TProtocol) (err error) {
	if p.IsSetMAP() {
		if err := oprot.WriteFieldBegin("MAP", thrift.STRUCT, 2); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:MAP: ", p), err)
		}
		if err := p.MAP.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.MAP), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 2:MAP: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetLIST() {
		if err := oprot.WriteFieldBegin("LIST", thrift.STRUCT, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:LIST: ", p), err)
		}
		if err := p.LIST.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.LIST), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:LIST: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetENUM() {
		if err := oprot.WriteFieldBegin("ENUM", thrift.STRUCT, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:ENUM: ", p), err)
		}
		if err := p.ENUM.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.ENUM), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:ENUM: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetDECIMAL() {
		if err := oprot.WriteFieldBegin("DECIMAL", thrift.STRUCT, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:DECIMAL: ", p), err)
		}
		if err := p.DECIMAL.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.DECIMAL), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:DECIMAL: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField6(oprot thrift.TProtocol) (err error) {
	if p.IsSetDATE() {
		if err := oprot.WriteFieldBegin("DATE", thrift.STRUCT, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:DATE: ", p), err)
		}
		if err := p.DATE.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.DATE), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:DATE: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField7(oprot thrift.TProtocol) (err error) {
	if p.IsSetTIME() {
		if err := oprot.WriteFieldBegin("TIME", thrift.STRUCT, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:TIME: ", p), err)
		}
		if err := p.TIME.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.TIME), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:TIME: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField8(oprot thrift.TProtocol) (err error) {
	if p.IsSetTIMESTAMP() {
		if err := oprot.WriteFieldBegin("TIMESTAMP", thrift.STRUCT, 8); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:TIMESTAMP: ", p), err)
		}
		if err := p.TIMESTAMP.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.TIMESTAMP), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 8:TIMESTAMP: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField10(oprot thrift.TProtocol) (err error) {
	if p.IsSetINTEGER() {
		if err := oprot.WriteFieldBegin("INTEGER", thrift.STRUCT, 10); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:INTEGER: ", p), err)
		}
		if err := p.INTEGER.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.INTEGER), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 10:INTEGER: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField11(oprot thrift.TProtocol) (err error) {
	if p.IsSetUNKNOWN() {
		if err := oprot.WriteFieldBegin("UNKNOWN", thrift.STRUCT, 11); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 11:UNKNOWN: ", p), err)
		}
		if err := p.UNKNOWN.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.UNKNOWN), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 11:UNKNOWN: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField12(oprot thrift.TProtocol) (err error) {
	if p.IsSetJSON() {
		if err := oprot.WriteFieldBegin("JSON", thrift.STRUCT, 12); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 12:JSON: ", p), err)
		}
		if err := p.JSON.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.JSON), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 12:JSON: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField13(oprot thrift.TProtocol) (err error) {
	if p.IsSetBSON() {
		if err := oprot.WriteFieldBegin("BSON", thrift.STRUCT, 13); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 13:BSON: ", p), err)
		}
		if err := p.BSON.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.BSON), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 13:BSON: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) writeField14(oprot thrift.TProtocol) (err error) {
	if p.IsSetUUID() {
		if err := oprot.WriteFieldBegin("UUID", thrift.STRUCT, 14); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 14:UUID: ", p), err)
		}
		if err := p.UUID.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.UUID), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 14:UUID: ", p), err)
		}
	}
	return err
}

func (p *LogicalType) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LogicalType(%+v)", *p)
}

// Represents a element inside a schema definition.
//  - if it is a group (inner node) then type is undefined and num_children is defined
//  - if it is a primitive type (leaf) then type is defined and num_children is undefined
// the nodes are listed in depth first traversal order.
//
// Attributes:
//  - Type: Data type for this field. Not set if the current element is a non-leaf node
//  - TypeLength: If type is FIXED_LEN_BYTE_ARRAY, this is the byte length of the vales.
// Otherwise, if specified, this is the maximum bit length to store any of the values.
// (e.g. a low cardinality INT col could have this set to 3).  Note that this is
// in the schema, and therefore fixed for the entire file.
//  - RepetitionType: repetition of the field. The root of the schema does not have a repetition_type.
// All other nodes must have one
//  - Name: Name of the field in the schema
//  - NumChildren: Nested fields.  Since thrift does not support nested fields,
// the nesting is flattened to a single list by a depth-first traversal.
// The children count is used to construct the nested relationship.
// This field is not set when the element is a primitive type
//  - ConvertedType: When the schema is the result of a conversion from another model
// Used to record the original type to help with cross conversion.
//  - Scale: Used when this column contains decimal data.
// See the DECIMAL converted type for more details.
//  - Precision
//  - FieldID: When the original schema supports field ids, this will save the
// original field id in the parquet schema
//  - LogicalType: The logical type of this SchemaElement
//
// LogicalType replaces ConvertedType, but ConvertedType is still required
// for some logical types to ensure forward-compatibility in format v1.
type SchemaElement struct {
	Type           *Type                `thrift:"type,1" db:"type" json:"type,omitempty"`
	TypeLength     *int32               `thrift:"type_length,2" db:"type_length" json:"type_length,omitempty"`
	RepetitionType *FieldRepetitionType `thrift:"repetition_type,3" db:"repetition_type" json:"repetition_type,omitempty"`
	Name           string               `thrift:"name,4,required" db:"name" json:"name"`
	NumChildren    *int32               `thrift:"num_children,5" db:"num_children" json:"num_children,omitempty"`
	ConvertedType  *ConvertedType       `thrift:"converted_type,6" db:"converted_type" json:"converted_type,omitempty"`
	Scale          *int32               `thrift:"scale,7" db:"scale" json:"scale,omitempty"`
	Precision      *int32               `thrift:"precision,8" db:"precision" json:"precision,omitempty"`
	FieldID        *int32               `thrift:"field_id,9" db:"field_id" json:"field_id,omitempty"`
	LogicalType    *LogicalType         `thrift:"logicalType,10" db:"logicalType" json:"logicalType,omitempty"`
}

func NewSchemaElement() *SchemaElement {
	return &SchemaElement{}
}

var SchemaElement_Type_DEFAULT Type

func (p *SchemaElement) GetType() Type {
	if !p.IsSetType() {
		return SchemaElement_Type_DEFAULT
	}
	return *p.Type
}

var SchemaElement_TypeLength_DEFAULT int32

func (p *SchemaElement) GetTypeLength() int32 {
	if !p.IsSetTypeLength() {
		return SchemaElement_TypeLength_DEFAULT
	}
	return *p.TypeLength
}

var SchemaElement_RepetitionType_DEFAULT FieldRepetitionType

func (p *SchemaElement) GetRepetitionType() FieldRepetitionType {
	if !p.IsSetRepetitionType() {
		return SchemaElement_RepetitionType_DEFAULT
	}
	return *p.RepetitionType
}

func (p *SchemaElement) GetName() string {
	return p.Name
}

var SchemaElement_NumChildren_DEFAULT int32

func (p *SchemaElement) GetNumChildren() int32 {
	if !p.IsSetNumChildren() {
		return SchemaElement_NumChildren_DEFAULT
	}
	return *p.NumChildren
}

var SchemaElement_ConvertedType_DEFAULT ConvertedType

func (p *SchemaElement) GetConvertedType() ConvertedType {
	if !p.IsSetConvertedType() {
		return SchemaElement_ConvertedType_DEFAULT
	}
	return *p.ConvertedType
}

var SchemaElement_Scale_DEFAULT int32

func (p *SchemaElement) GetScale() int32 {
	if !p.IsSetScale() {
		return SchemaElement_Scale_DEFAULT
	}
	return *p.Scale
}

var SchemaElement_Precision_DEFAULT int32

func (p *SchemaElement) GetPrecision() int32 {
	if !p.IsSetPrecision() {
		return SchemaElement_Precision_DEFAULT
	}
	return *p.Precision
}

var SchemaElement_FieldID_DEFAULT int32

func (p *SchemaElement) GetFieldID() int32 {
	if !p.IsSetFieldID() {
		return SchemaElement_FieldID_DEFAULT
	}
	return *p.FieldID
}

var SchemaElement_LogicalType_DEFAULT *LogicalType

func (p *SchemaElement) GetLogicalType() *LogicalType {
	if !p.IsSetLogicalType() {
		return SchemaElement_LogicalType_DEFAULT
	}
	return p.LogicalType
}
func (p *SchemaElement) IsSetType() bool {
	return p.Type != nil
}

func (p *SchemaElement) IsSetTypeLength() bool {
	return p.TypeLength != nil
}

func (p *SchemaElement) IsSetRepetitionType() bool {
	return p.RepetitionType != nil
}

func (p *SchemaElement) IsSetNumChildren() bool {
	return p.NumChildren != nil
}

func (p *SchemaElement) IsSetConvertedType() bool {
	return p.ConvertedType != nil
}

func (p *SchemaElement) IsSetScale() bool {
	return p.Scale != nil
}

func (p *SchemaElement) IsSetPrecision() bool {
	return p.Precision != nil
}

func (p *SchemaElement) IsSetFieldID() bool {
	return p.FieldID != nil
}

func (p *SchemaElement) IsSetLogicalType() bool {
	return p.LogicalType != nil
}

func (p *SchemaElement) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetName bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
			issetName = true
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.ReadField8(iprot); err != nil {
				return err
			}
		case 9:
			if err := p.ReadField9(iprot); err != nil {
				return err
			}
		case 10:
			if err := p.ReadField10(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetName {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Name is not set"))
	}
	return nil
}

func (p *SchemaElement) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := Type(v)
		p.Type = &temp
	}
	return nil
}

func (p *SchemaElement) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.TypeLength = &v
	}
	return nil
}

func (p *SchemaElement) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		temp := FieldRepetitionType(v)
		p.RepetitionType = &temp
	}
	return nil
}

func (p *SchemaElement) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *SchemaElement) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.NumChildren = &v
	}
	return nil
}

func (p *SchemaElement) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		temp := ConvertedType(v)
		p.ConvertedType = &temp
	}
	return nil
}

func (p *SchemaElement) ReadField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 7: ", err)
	} else {
		p.Scale = &v
	}
	return nil
}

func (p *SchemaElement) ReadField8(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 8: ", err)
	} else {
		p.Precision = &v
	}
	return nil
}

func (p *SchemaElement) ReadField9(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 9: ", err)
	} else {
		p.FieldID = &v
	}
	return nil
}

func (p *SchemaElement) ReadField10(iprot thrift.TProtocol) error {
	p.LogicalType = &LogicalType{}
	if err := p.LogicalType.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.LogicalType), err)
	}
	return nil
}

func (p *SchemaElement) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("SchemaElement"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
		if err := p.writeField8(oprot); err != nil {
			return err
		}
		if err := p.writeField9(oprot); err != nil {
			return err
		}
		if err := p.writeField10(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *SchemaElement) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetType() {
		if err := oprot.WriteFieldBegin("type", thrift.I32, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:type: ", p), err)
		}
		if err := oprot.WriteI32(int32(*p.Type)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.type (1) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:type: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField2(oprot thrift.TProtocol) (err error) {
	if p.IsSetTypeLength() {
		if err := oprot.WriteFieldBegin("type_length", thrift.I32, 2); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:type_length: ", p), err)
		}
		if err := oprot.WriteI32(*p.TypeLength); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.type_length (2) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 2:type_length: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetRepetitionType() {
		if err := oprot.WriteFieldBegin("repetition_type", thrift.I32, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:repetition_type: ", p), err)
		}
		if err := oprot.WriteI32(int32(*p.RepetitionType)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.repetition_type (3) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:repetition_type: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:name: ", p), err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.name (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:name: ", p), err)
	}
	return err
}

func (p *SchemaElement) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetNumChildren() {
		if err := oprot.WriteFieldBegin("num_children", thrift.I32, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:num_children: ", p), err)
		}
		if err := oprot.WriteI32(*p.NumChildren); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.num_children (5) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:num_children: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField6(oprot thrift.TProtocol) (err error) {
	if p.IsSetConvertedType() {
		if err := oprot.WriteFieldBegin("converted_type", thrift.I32, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:converted_type: ", p), err)
		}
		if err := oprot.WriteI32(int32(*p.ConvertedType)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.converted_type (6) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:converted_type: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField7(oprot thrift.TProtocol) (err error) {
	if p.IsSetScale() {
		if err := oprot.WriteFieldBegin("scale", thrift.I32, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:scale: ", p), err)
		}
		if err := oprot.WriteI32(*p.Scale); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.scale (7) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:scale: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField8(oprot thrift.TProtocol) (err error) {
	if p.IsSetPrecision() {
		if err := oprot.WriteFieldBegin("precision", thrift.I32, 8); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:precision: ", p), err)
		}
		if err := oprot.WriteI32(*p.Precision); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.precision (8) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 8:precision: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField9(oprot thrift.TProtocol) (err error) {
	if p.IsSetFieldID() {
		if err := oprot.WriteFieldBegin("field_id", thrift.I32, 9); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 9:field_id: ", p), err)
		}
		if err := oprot.WriteI32(*p.FieldID); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.field_id (9) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 9:field_id: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) writeField10(oprot thrift.TProtocol) (err error) {
	if p.IsSetLogicalType() {
		if err := oprot.WriteFieldBegin("logicalType", thrift.STRUCT, 10); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:logicalType: ", p), err)
		}
		if err := p.LogicalType.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.LogicalType), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 10:logicalType: ", p), err)
		}
	}
	return err
}

func (p *SchemaElement) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SchemaElement(%+v)", *p)
}

// Data page header
//
// Attributes:
//  - NumValues: Number of values, including NULLs, in this data page. *
//  - Encoding: Encoding used for this data page *
//  - DefinitionLevelEncoding: Encoding used for definition levels *
//  - RepetitionLevelEncoding: Encoding used for repetition levels *
//  - Statistics: Optional statistics for the data in this page*
type DataPageHeader struct {
	NumValues               int32       `thrift:"num_values,1,required" db:"num_values" json:"num_values"`
	Encoding                Encoding    `thrift:"encoding,2,required" db:"encoding" json:"encoding"`
	DefinitionLevelEncoding Encoding    `thrift:"definition_level_encoding,3,required" db:"definition_level_encoding" json:"definition_level_encoding"`
	RepetitionLevelEncoding Encoding    `thrift:"repetition_level_encoding,4,required" db:"repetition_level_encoding" json:"repetition_level_encoding"`
	Statistics              *Statistics `thrift:"statistics,5" db:"statistics" json:"statistics,omitempty"`
}

func NewDataPageHeader() *DataPageHeader {
	return &DataPageHeader{}
}

func (p *DataPageHeader) GetNumValues() int32 {
	return p.NumValues
}

func (p *DataPageHeader) GetEncoding() Encoding {
	return p.Encoding
}

func (p *DataPageHeader) GetDefinitionLevelEncoding() Encoding {
	return p.DefinitionLevelEncoding
}

func (p *DataPageHeader) GetRepetitionLevelEncoding() Encoding {
	return p.RepetitionLevelEncoding
}

var DataPageHeader_Statistics_DEFAULT *Statistics

func (p *DataPageHeader) GetStatistics() *Statistics {
	if !p.IsSetStatistics() {
		return DataPageHeader_Statistics_DEFAULT
	}
	return p.Statistics
}
func (p *DataPageHeader) IsSetStatistics() bool {
	return p.Statistics != nil
}

func (p *DataPageHeader) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetNumValues bool = false
	var issetEncoding bool = false
	var issetDefinitionLevelEncoding bool = false
	var issetRepetitionLevelEncoding bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetNumValues = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetEncoding = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetDefinitionLevelEncoding = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
			issetRepetitionLevelEncoding = true
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetNumValues {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumValues is not set"))
	}
	if !issetEncoding {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Encoding is not set"))
	}
	if !issetDefinitionLevelEncoding {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field DefinitionLevelEncoding is not set"))
	}
	if !issetRepetitionLevelEncoding {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RepetitionLevelEncoding is not set"))
	}
	return nil
}

func (p *DataPageHeader) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.NumValues = v
	}
	return nil
}

func (p *DataPageHeader) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		temp := Encoding(v)
		p.Encoding = temp
	}
	return nil
}

func (p *DataPageHeader) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		temp := Encoding(v)
		p.DefinitionLevelEncoding = temp
	}
	return nil
}

func (p *DataPageHeader) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		temp := Encoding(v)
		p.RepetitionLevelEncoding = temp
	}
	return nil
}

func (p *DataPageHeader) ReadField5(iprot thrift.TProtocol) error {
	p.Statistics = &Statistics{}
	if err := p.Statistics.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Statistics), err)
	}
	return nil
}

func (p *DataPageHeader) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("DataPageHeader"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *DataPageHeader) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_values", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:num_values: ", p), err)
	}
	if err := oprot.WriteI32(p.NumValues); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_values (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:num_values: ", p), err)
	}
	return err
}

func (p *DataPageHeader) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("encoding", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:encoding: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Encoding)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.encoding (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:encoding: ", p), err)
	}
	return err
}

func (p *DataPageHeader) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("definition_level_encoding", thrift.I32, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:definition_level_encoding: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.DefinitionLevelEncoding)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.definition_level_encoding (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:definition_level_encoding: ", p), err)
	}
	return err
}

func (p *DataPageHeader) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("repetition_level_encoding", thrift.I32, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:repetition_level_encoding: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.RepetitionLevelEncoding)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.repetition_level_encoding (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:repetition_level_encoding: ", p), err)
	}
	return err
}

func (p *DataPageHeader) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetStatistics() {
		if err := oprot.WriteFieldBegin("statistics", thrift.STRUCT, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:statistics: ", p), err)
		}
		if err := p.Statistics.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Statistics), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:statistics: ", p), err)
		}
	}
	return err
}

func (p *DataPageHeader) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DataPageHeader(%+v)", *p)
}

type IndexPageHeader struct {
}

func NewIndexPageHeader() *IndexPageHeader {
	return &IndexPageHeader{}
}

func (p *IndexPageHeader) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *IndexPageHeader) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("IndexPageHeader"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *IndexPageHeader) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("IndexPageHeader(%+v)", *p)
}

// TODO: *
//
// Attributes:
//  - NumValues: Number of values in the dictionary *
//  - Encoding: Encoding using this dictionary page *
//  - IsSorted: If true, the entries in the dictionary are sorted in ascending order *
type DictionaryPageHeader struct {
	NumValues int32    `thrift:"num_values,1,required" db:"num_values" json:"num_values"`
	Encoding  Encoding `thrift:"encoding,2,required" db:"encoding" json:"encoding"`
	IsSorted  *bool    `thrift:"is_sorted,3" db:"is_sorted" json:"is_sorted,omitempty"`
}

func NewDictionaryPageHeader() *DictionaryPageHeader {
	return &DictionaryPageHeader{}
}

func (p *DictionaryPageHeader) GetNumValues() int32 {
	return p.NumValues
}

func (p *DictionaryPageHeader) GetEncoding() Encoding {
	return p.Encoding
}

var DictionaryPageHeader_IsSorted_DEFAULT bool

func (p *DictionaryPageHeader) GetIsSorted() bool {
	if !p.IsSetIsSorted() {
		return DictionaryPageHeader_IsSorted_DEFAULT
	}
	return *p.IsSorted
}
func (p *DictionaryPageHeader) IsSetIsSorted() bool {
	return p.IsSorted != nil
}

func (p *DictionaryPageHeader) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetNumValues bool = false
	var issetEncoding bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetNumValues = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetEncoding = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetNumValues {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumValues is not set"))
	}
	if !issetEncoding {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Encoding is not set"))
	}
	return nil
}

func (p *DictionaryPageHeader) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.NumValues = v
	}
	return nil
}

func (p *DictionaryPageHeader) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		temp := Encoding(v)
		p.Encoding = temp
	}
	return nil
}

func (p *DictionaryPageHeader) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.IsSorted = &v
	}
	return nil
}

func (p *DictionaryPageHeader) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("DictionaryPageHeader"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *DictionaryPageHeader) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_values", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:num_values: ", p), err)
	}
	if err := oprot.WriteI32(p.NumValues); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_values (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:num_values: ", p), err)
	}
	return err
}

func (p *DictionaryPageHeader) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("encoding", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:encoding: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Encoding)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.encoding (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:encoding: ", p), err)
	}
	return err
}

func (p *DictionaryPageHeader) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetIsSorted() {
		if err := oprot.WriteFieldBegin("is_sorted", thrift.BOOL, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:is_sorted: ", p), err)
		}
		if err := oprot.WriteBool(*p.IsSorted); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.is_sorted (3) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:is_sorted: ", p), err)
		}
	}
	return err
}

func (p *DictionaryPageHeader) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DictionaryPageHeader(%+v)", *p)
}

// New page format allowing reading levels without decompressing the data
// Repetition and definition levels are uncompressed
// The remaining section containing the data is compressed if is_compressed is true
//
//
// Attributes:
//  - NumValues: Number of values, including NULLs, in this data page. *
//  - NumNulls: Number of NULL values, in this data page.
// Number of non-null = num_values - num_nulls which is also the number of values in the data section *
//  - NumRows: Number of rows in this data page. which means pages change on record boundaries (r = 0) *
//  - Encoding: Encoding used for data in this page *
//  - DefinitionLevelsByteLength: length of the definition levels
//  - RepetitionLevelsByteLength: length of the repetition levels
//  - IsCompressed: whether the values are compressed.
// Which means the section of the page between
// definition_levels_byte_length + repetition_levels_byte_length + 1 and compressed_page_size (included)
// is compressed with the compression_codec.
// If missing it is considered compressed
//  - Statistics: optional statistics for this column chunk
type DataPageHeaderV2 struct {
	NumValues                  int32       `thrift:"num_values,1,required" db:"num_values" json:"num_values"`
	NumNulls                   int32       `thrift:"num_nulls,2,required" db:"num_nulls" json:"num_nulls"`
	NumRows                    int32       `thrift:"num_rows,3,required" db:"num_rows" json:"num_rows"`
	Encoding                   Encoding    `thrift:"encoding,4,required" db:"encoding" json:"encoding"`
	DefinitionLevelsByteLength int32       `thrift:"definition_levels_byte_length,5,required" db:"definition_levels_byte_length" json:"definition_levels_byte_length"`
	RepetitionLevelsByteLength int32       `thrift:"repetition_levels_byte_length,6,required" db:"repetition_levels_byte_length" json:"repetition_levels_byte_length"`
	IsCompressed               bool        `thrift:"is_compressed,7" db:"is_compressed" json:"is_compressed,omitempty"`
	Statistics                 *Statistics `thrift:"statistics,8" db:"statistics" json:"statistics,omitempty"`
}

func NewDataPageHeaderV2() *DataPageHeaderV2 {
	return &DataPageHeaderV2{
		IsCompressed: true,
	}
}

func (p *DataPageHeaderV2) GetNumValues() int32 {
	return p.NumValues
}

func (p *DataPageHeaderV2) GetNumNulls() int32 {
	return p.NumNulls
}

func (p *DataPageHeaderV2) GetNumRows() int32 {
	return p.NumRows
}

func (p *DataPageHeaderV2) GetEncoding() Encoding {
	return p.Encoding
}

func (p *DataPageHeaderV2) GetDefinitionLevelsByteLength() int32 {
	return p.DefinitionLevelsByteLength
}

func (p *DataPageHeaderV2) GetRepetitionLevelsByteLength() int32 {
	return p.RepetitionLevelsByteLength
}

var DataPageHeaderV2_IsCompressed_DEFAULT bool = true

func (p *DataPageHeaderV2) GetIsCompressed() bool {
	return p.IsCompressed
}

var DataPageHeaderV2_Statistics_DEFAULT *Statistics

func (p *DataPageHeaderV2) GetStatistics() *Statistics {
	if !p.IsSetStatistics() {
		return DataPageHeaderV2_Statistics_DEFAULT
	}
	return p.Statistics
}
func (p *DataPageHeaderV2) IsSetIsCompressed() bool {
	return p.IsCompressed != DataPageHeaderV2_IsCompressed_DEFAULT
}

func (p *DataPageHeaderV2) IsSetStatistics() bool {
	return p.Statistics != nil
}

func (p *DataPageHeaderV2) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetNumValues bool = false
	var issetNumNulls bool = false
	var issetNumRows bool = false
	var issetEncoding bool = false
	var issetDefinitionLevelsByteLength bool = false
	var issetRepetitionLevelsByteLength bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetNumValues = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetNumNulls = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetNumRows = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
			issetEncoding = true
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
			issetDefinitionLevelsByteLength = true
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
			issetRepetitionLevelsByteLength = true
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.ReadField8(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetNumValues {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumValues is not set"))
	}
	if !issetNumNulls {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumNulls is not set"))
	}
	if !issetNumRows {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumRows is not set"))
	}
	if !issetEncoding {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Encoding is not set"))
	}
	if !issetDefinitionLevelsByteLength {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field DefinitionLevelsByteLength is not set"))
	}
	if !issetRepetitionLevelsByteLength {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RepetitionLevelsByteLength is not set"))
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.NumValues = v
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.NumNulls = v
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.NumRows = v
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		temp := Encoding(v)
		p.Encoding = temp
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.DefinitionLevelsByteLength = v
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.RepetitionLevelsByteLength = v
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 7: ", err)
	} else {
		p.IsCompressed = v
	}
	return nil
}

func (p *DataPageHeaderV2) ReadField8(iprot thrift.TProtocol) error {
	p.Statistics = &Statistics{}
	if err := p.Statistics.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Statistics), err)
	}
	return nil
}

func (p *DataPageHeaderV2) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("DataPageHeaderV2"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
		if err := p.writeField8(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *DataPageHeaderV2) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_values", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:num_values: ", p), err)
	}
	if err := oprot.WriteI32(p.NumValues); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_values (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:num_values: ", p), err)
	}
	return err
}

func (p *DataPageHeaderV2) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_nulls", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:num_nulls: ", p), err)
	}
	if err := oprot.WriteI32(p.NumNulls); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_nulls (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:num_nulls: ", p), err)
	}
	return err
}

func (p *DataPageHeaderV2) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_rows", thrift.I32, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:num_rows: ", p), err)
	}
	if err := oprot.WriteI32(p.NumRows); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_rows (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:num_rows: ", p), err)
	}
	return err
}

func (p *DataPageHeaderV2) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("encoding", thrift.I32, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:encoding: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Encoding)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.encoding (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:encoding: ", p), err)
	}
	return err
}

func (p *DataPageHeaderV2) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("definition_levels_byte_length", thrift.I32, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:definition_levels_byte_length: ", p), err)
	}
	if err := oprot.WriteI32(p.DefinitionLevelsByteLength); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.definition_levels_byte_length (5) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:definition_levels_byte_length: ", p), err)
	}
	return err
}

func (p *DataPageHeaderV2) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("repetition_levels_byte_length", thrift.I32, 6); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:repetition_levels_byte_length: ", p), err)
	}
	if err := oprot.WriteI32(p.RepetitionLevelsByteLength); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.repetition_levels_byte_length (6) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 6:repetition_levels_byte_length: ", p), err)
	}
	return err
}

func (p *DataPageHeaderV2) writeField7(oprot thrift.TProtocol) (err error) {
	if p.IsSetIsCompressed() {
		if err := oprot.WriteFieldBegin("is_compressed", thrift.BOOL, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:is_compressed: ", p), err)
		}
		if err := oprot.WriteBool(p.IsCompressed); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.is_compressed (7) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:is_compressed: ", p), err)
		}
	}
	return err
}

func (p *DataPageHeaderV2) writeField8(oprot thrift.TProtocol) (err error) {
	if p.IsSetStatistics() {
		if err := oprot.WriteFieldBegin("statistics", thrift.STRUCT, 8); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:statistics: ", p), err)
		}
		if err := p.Statistics.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Statistics), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 8:statistics: ", p), err)
		}
	}
	return err
}

func (p *DataPageHeaderV2) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DataPageHeaderV2(%+v)", *p)
}

// Attributes:
//  - Type: the type of the page: indicates which of the *_header fields is set *
//  - UncompressedPageSize: Uncompressed page size in bytes (not including this header) *
//  - CompressedPageSize: Compressed page size in bytes (not including this header) *
//  - Crc: 32bit crc for the data below. This allows for disabling checksumming in HDFS
// if only a few pages needs to be read
//
//  - DataPageHeader
//  - IndexPageHeader
//  - DictionaryPageHeader
//  - DataPageHeaderV2
type PageHeader struct {
	Type                 PageType              `thrift:"type,1,required" db:"type" json:"type"`
	UncompressedPageSize int32                 `thrift:"uncompressed_page_size,2,required" db:"uncompressed_page_size" json:"uncompressed_page_size"`
	CompressedPageSize   int32                 `thrift:"compressed_page_size,3,required" db:"compressed_page_size" json:"compressed_page_size"`
	Crc                  *int32                `thrift:"crc,4" db:"crc" json:"crc,omitempty"`
	DataPageHeader       *DataPageHeader       `thrift:"data_page_header,5" db:"data_page_header" json:"data_page_header,omitempty"`
	IndexPageHeader      *IndexPageHeader      `thrift:"index_page_header,6" db:"index_page_header" json:"index_page_header,omitempty"`
	DictionaryPageHeader *DictionaryPageHeader `thrift:"dictionary_page_header,7" db:"dictionary_page_header" json:"dictionary_page_header,omitempty"`
	DataPageHeaderV2     *DataPageHeaderV2     `thrift:"data_page_header_v2,8" db:"data_page_header_v2" json:"data_page_header_v2,omitempty"`
}

func NewPageHeader() *PageHeader {
	return &PageHeader{}
}

func (p *PageHeader) GetType() PageType {
	return p.Type
}

func (p *PageHeader) GetUncompressedPageSize() int32 {
	return p.UncompressedPageSize
}

func (p *PageHeader) GetCompressedPageSize() int32 {
	return p.CompressedPageSize
}

var PageHeader_Crc_DEFAULT int32

func (p *PageHeader) GetCrc() int32 {
	if !p.IsSetCrc() {
		return PageHeader_Crc_DEFAULT
	}
	return *p.Crc
}

var PageHeader_DataPageHeader_DEFAULT *DataPageHeader

func (p *PageHeader) GetDataPageHeader() *DataPageHeader {
	if !p.IsSetDataPageHeader() {
		return PageHeader_DataPageHeader_DEFAULT
	}
	return p.DataPageHeader
}

var PageHeader_IndexPageHeader_DEFAULT *IndexPageHeader

func (p *PageHeader) GetIndexPageHeader() *IndexPageHeader {
	if !p.IsSetIndexPageHeader() {
		return PageHeader_IndexPageHeader_DEFAULT
	}
	return p.IndexPageHeader
}

var PageHeader_DictionaryPageHeader_DEFAULT *DictionaryPageHeader

func (p *PageHeader) GetDictionaryPageHeader() *DictionaryPageHeader {
	if !p.IsSetDictionaryPageHeader() {
		return PageHeader_DictionaryPageHeader_DEFAULT
	}
	return p.DictionaryPageHeader
}

var PageHeader_DataPageHeaderV2_DEFAULT *DataPageHeaderV2

func (p *PageHeader) GetDataPageHeaderV2() *DataPageHeaderV2 {
	if !p.IsSetDataPageHeaderV2() {
		return PageHeader_DataPageHeaderV2_DEFAULT
	}
	return p.DataPageHeaderV2
}
func (p *PageHeader) IsSetCrc() bool {
	return p.Crc != nil
}

func (p *PageHeader) IsSetDataPageHeader() bool {
	return p.DataPageHeader != nil
}

func (p *PageHeader) IsSetIndexPageHeader() bool {
	return p.IndexPageHeader != nil
}

func (p *PageHeader) IsSetDictionaryPageHeader() bool {
	return p.DictionaryPageHeader != nil
}

func (p *PageHeader) IsSetDataPageHeaderV2() bool {
	return p.DataPageHeaderV2 != nil
}

func (p *PageHeader) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetType bool = false
	var issetUncompressedPageSize bool = false
	var issetCompressedPageSize bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetType = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetUncompressedPageSize = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetCompressedPageSize = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.ReadField8(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetType {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Type is not set"))
	}
	if !issetUncompressedPageSize {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field UncompressedPageSize is not set"))
	}
	if !issetCompressedPageSize {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field CompressedPageSize is not set"))
	}
	return nil
}

func (p *PageHeader) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := PageType(v)
		p.Type = temp
	}
	return nil
}

func (p *PageHeader) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.UncompressedPageSize = v
	}
	return nil
}

func (p *PageHeader) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.CompressedPageSize = v
	}
	return nil
}

func (p *PageHeader) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.Crc = &v
	}
	return nil
}

func (p *PageHeader) ReadField5(iprot thrift.TProtocol) error {
	p.DataPageHeader = &DataPageHeader{}
	if err := p.DataPageHeader.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.DataPageHeader), err)
	}
	return nil
}

func (p *PageHeader) ReadField6(iprot thrift.TProtocol) error {
	p.IndexPageHeader = &IndexPageHeader{}
	if err := p.IndexPageHeader.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.IndexPageHeader), err)
	}
	return nil
}

func (p *PageHeader) ReadField7(iprot thrift.TProtocol) error {
	p.DictionaryPageHeader = &DictionaryPageHeader{}
	if err := p.DictionaryPageHeader.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.DictionaryPageHeader), err)
	}
	return nil
}

func (p *PageHeader) ReadField8(iprot thrift.TProtocol) error {
	p.DataPageHeaderV2 = &DataPageHeaderV2{
		IsCompressed: true,
	}
	if err := p.DataPageHeaderV2.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.DataPageHeaderV2), err)
	}
	return nil
}

func (p *PageHeader) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("PageHeader"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
		if err := p.writeField8(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *PageHeader) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("type", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:type: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Type)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.type (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:type: ", p), err)
	}
	return err
}

func (p *PageHeader) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("uncompressed_page_size", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:uncompressed_page_size: ", p), err)
	}
	if err := oprot.WriteI32(p.UncompressedPageSize); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.uncompressed_page_size (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:uncompressed_page_size: ", p), err)
	}
	return err
}

func (p *PageHeader) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("compressed_page_size", thrift.I32, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:compressed_page_size: ", p), err)
	}
	if err := oprot.WriteI32(p.CompressedPageSize); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.compressed_page_size (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:compressed_page_size: ", p), err)
	}
	return err
}

func (p *PageHeader) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetCrc() {
		if err := oprot.WriteFieldBegin("crc", thrift.I32, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:crc: ", p), err)
		}
		if err := oprot.WriteI32(*p.Crc); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.crc (4) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:crc: ", p), err)
		}
	}
	return err
}

func (p *PageHeader) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetDataPageHeader() {
		if err := oprot.WriteFieldBegin("data_page_header", thrift.STRUCT, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:data_page_header: ", p), err)
		}
		if err := p.DataPageHeader.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.DataPageHeader), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:data_page_header: ", p), err)
		}
	}
	return err
}

func (p *PageHeader) writeField6(oprot thrift.TProtocol) (err error) {
	if p.IsSetIndexPageHeader() {
		if err := oprot.WriteFieldBegin("index_page_header", thrift.STRUCT, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:index_page_header: ", p), err)
		}
		if err := p.IndexPageHeader.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.IndexPageHeader), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:index_page_header: ", p), err)
		}
	}
	return err
}

func (p *PageHeader) writeField7(oprot thrift.TProtocol) (err error) {
	if p.IsSetDictionaryPageHeader() {
		if err := oprot.WriteFieldBegin("dictionary_page_header", thrift.STRUCT, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:dictionary_page_header: ", p), err)
		}
		if err := p.DictionaryPageHeader.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.DictionaryPageHeader), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:dictionary_page_header: ", p), err)
		}
	}
	return err
}

func (p *PageHeader) writeField8(oprot thrift.TProtocol) (err error) {
	if p.IsSetDataPageHeaderV2() {
		if err := oprot.WriteFieldBegin("data_page_header_v2", thrift.STRUCT, 8); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:data_page_header_v2: ", p), err)
		}
		if err := p.DataPageHeaderV2.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.DataPageHeaderV2), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 8:data_page_header_v2: ", p), err)
		}
	}
	return err
}

func (p *PageHeader) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PageHeader(%+v)", *p)
}

// Wrapper struct to store key values
//
// Attributes:
//  - Key
//  - Value
type KeyValue struct {
	Key   string  `thrift:"key,1,required" db:"key" json:"key"`
	Value *string `thrift:"value,2" db:"value" json:"value,omitempty"`
}

func NewKeyValue() *KeyValue {
	return &KeyValue{}
}

func (p *KeyValue) GetKey() string {
	return p.Key
}

var KeyValue_Value_DEFAULT string

func (p *KeyValue) GetValue() string {
	if !p.IsSetValue() {
		return KeyValue_Value_DEFAULT
	}
	return *p.Value
}
func (p *KeyValue) IsSetValue() bool {
	return p.Value != nil
}

func (p *KeyValue) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetKey bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetKey = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetKey {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Key is not set"))
	}
	return nil
}

func (p *KeyValue) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Key = v
	}
	return nil
}

func (p *KeyValue) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Value = &v
	}
	return nil
}

func (p *KeyValue) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("KeyValue"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *KeyValue) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("key", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:key: ", p), err)
	}
	if err := oprot.WriteString(string(p.Key)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.key (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:key: ", p), err)
	}
	return err
}

func (p *KeyValue) writeField2(oprot thrift.TProtocol) (err error) {
	if p.IsSetValue() {
		if err := oprot.WriteFieldBegin("value", thrift.STRING, 2); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:value: ", p), err)
		}
		if err := oprot.WriteString(string(*p.Value)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.value (2) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 2:value: ", p), err)
		}
	}
	return err
}

func (p *KeyValue) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KeyValue(%+v)", *p)
}

// Wrapper struct to specify sort order
//
// Attributes:
//  - ColumnIdx: The column index (in this row group) *
//  - Descending: If true, indicates this column is sorted in descending order. *
//  - NullsFirst: If true, nulls will come before non-null values, otherwise,
// nulls go at the end.
type SortingColumn struct {
	ColumnIdx  int32 `thrift:"column_idx,1,required" db:"column_idx" json:"column_idx"`
	Descending bool  `thrift:"descending,2,required" db:"descending" json:"descending"`
	NullsFirst bool  `thrift:"nulls_first,3,required" db:"nulls_first" json:"nulls_first"`
}

func NewSortingColumn() *SortingColumn {
	return &SortingColumn{}
}

func (p *SortingColumn) GetColumnIdx() int32 {
	return p.ColumnIdx
}

func (p *SortingColumn) GetDescending() bool {
	return p.Descending
}

func (p *SortingColumn) GetNullsFirst() bool {
	return p.NullsFirst
}
func (p *SortingColumn) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetColumnIdx bool = false
	var issetDescending bool = false
	var issetNullsFirst bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetColumnIdx = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetDescending = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetNullsFirst = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetColumnIdx {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field ColumnIdx is not set"))
	}
	if !issetDescending {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Descending is not set"))
	}
	if !issetNullsFirst {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NullsFirst is not set"))
	}
	return nil
}

func (p *SortingColumn) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.ColumnIdx = v
	}
	return nil
}

func (p *SortingColumn) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Descending = v
	}
	return nil
}

func (p *SortingColumn) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.NullsFirst = v
	}
	return nil
}

func (p *SortingColumn) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("SortingColumn"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *SortingColumn) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("column_idx", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:column_idx: ", p), err)
	}
	if err := oprot.WriteI32(p.ColumnIdx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.column_idx (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:column_idx: ", p), err)
	}
	return err
}

func (p *SortingColumn) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("descending", thrift.BOOL, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:descending: ", p), err)
	}
	if err := oprot.WriteBool(p.Descending); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.descending (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:descending: ", p), err)
	}
	return err
}

func (p *SortingColumn) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("nulls_first", thrift.BOOL, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:nulls_first: ", p), err)
	}
	if err := oprot.WriteBool(p.NullsFirst); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.nulls_first (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:nulls_first: ", p), err)
	}
	return err
}

func (p *SortingColumn) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SortingColumn(%+v)", *p)
}

// statistics of a given page type and encoding
//
// Attributes:
//  - PageType: the page type (data/dic/...) *
//  - Encoding: encoding of the page *
//  - Count: number of pages of this type with this encoding *
type PageEncodingStats struct {
	PageType PageType `thrift:"page_type,1,required" db:"page_type" json:"page_type"`
	Encoding Encoding `thrift:"encoding,2,required" db:"encoding" json:"encoding"`
	Count    int32    `thrift:"count,3,required" db:"count" json:"count"`
}

func NewPageEncodingStats() *PageEncodingStats {
	return &PageEncodingStats{}
}

func (p *PageEncodingStats) GetPageType() PageType {
	return p.PageType
}

func (p *PageEncodingStats) GetEncoding() Encoding {
	return p.Encoding
}

func (p *PageEncodingStats) GetCount() int32 {
	return p.Count
}
func (p *PageEncodingStats) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetPageType bool = false
	var issetEncoding bool = false
	var issetCount bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetPageType = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetEncoding = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetCount = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetPageType {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field PageType is not set"))
	}
	if !issetEncoding {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Encoding is not set"))
	}
	if !issetCount {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Count is not set"))
	}
	return nil
}

func (p *PageEncodingStats) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := PageType(v)
		p.PageType = temp
	}
	return nil
}

func (p *PageEncodingStats) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		temp := Encoding(v)
		p.Encoding = temp
	}
	return nil
}

func (p *PageEncodingStats) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Count = v
	}
	return nil
}

func (p *PageEncodingStats) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("PageEncodingStats"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *PageEncodingStats) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("page_type", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:page_type: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.PageType)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.page_type (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:page_type: ", p), err)
	}
	return err
}

func (p *PageEncodingStats) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("encoding", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:encoding: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Encoding)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.encoding (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:encoding: ", p), err)
	}
	return err
}

func (p *PageEncodingStats) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("count", thrift.I32, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:count: ", p), err)
	}
	if err := oprot.WriteI32(p.Count); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.count (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:count: ", p), err)
	}
	return err
}

func (p *PageEncodingStats) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PageEncodingStats(%+v)", *p)
}

// Description for column metadata
//
// Attributes:
//  - Type: Type of this column *
//  - Encodings: Set of all encodings used for this column. The purpose is to validate
// whether we can decode those pages. *
//  - PathInSchema: Path in schema *
//  - Codec: Compression codec *
//  - NumValues: Number of values in this column *
//  - TotalUncompressedSize: total byte size of all uncompressed pages in this column chunk (including the headers) *
//  - TotalCompressedSize: total byte size of all compressed pages in this column chunk (including the headers) *
//  - KeyValueMetadata: Optional key/value metadata *
//  - DataPageOffset: Byte offset from beginning of file to first data page *
//  - IndexPageOffset: Byte offset from beginning of file to root index page *
//  - DictionaryPageOffset: Byte offset from the beginning of file to first (only) dictionary page *
//  - Statistics: optional statistics for this column chunk
//  - EncodingStats: Set of all encodings used for pages in this column chunk.
// This information can be used to determine if all data pages are
// dictionary encoded for example *
type ColumnMetaData struct {
	Type                  Type                 `thrift:"type,1,required" db:"type" json:"type"`
	Encodings             []Encoding           `thrift:"encodings,2,required" db:"encodings" json:"encodings"`
	PathInSchema          []string             `thrift:"path_in_schema,3,required" db:"path_in_schema" json:"path_in_schema"`
	Codec                 CompressionCodec     `thrift:"codec,4,required" db:"codec" json:"codec"`
	NumValues             int64                `thrift:"num_values,5,required" db:"num_values" json:"num_values"`
	TotalUncompressedSize int64                `thrift:"total_uncompressed_size,6,required" db:"total_uncompressed_size" json:"total_uncompressed_size"`
	TotalCompressedSize   int64                `thrift:"total_compressed_size,7,required" db:"total_compressed_size" json:"total_compressed_size"`
	KeyValueMetadata      []*KeyValue          `thrift:"key_value_metadata,8" db:"key_value_metadata" json:"key_value_metadata,omitempty"`
	DataPageOffset        int64                `thrift:"data_page_offset,9,required" db:"data_page_offset" json:"data_page_offset"`
	IndexPageOffset       *int64               `thrift:"index_page_offset,10" db:"index_page_offset" json:"index_page_offset,omitempty"`
	DictionaryPageOffset  *int64               `thrift:"dictionary_page_offset,11" db:"dictionary_page_offset" json:"dictionary_page_offset,omitempty"`
	Statistics            *Statistics          `thrift:"statistics,12" db:"statistics" json:"statistics,omitempty"`
	EncodingStats         []*PageEncodingStats `thrift:"encoding_stats,13" db:"encoding_stats" json:"encoding_stats,omitempty"`
}

func NewColumnMetaData() *ColumnMetaData {
	return &ColumnMetaData{}
}

func (p *ColumnMetaData) GetType() Type {
	return p.Type
}

func (p *ColumnMetaData) GetEncodings() []Encoding {
	return p.Encodings
}

func (p *ColumnMetaData) GetPathInSchema() []string {
	return p.PathInSchema
}

func (p *ColumnMetaData) GetCodec() CompressionCodec {
	return p.Codec
}

func (p *ColumnMetaData) GetNumValues() int64 {
	return p.NumValues
}

func (p *ColumnMetaData) GetTotalUncompressedSize() int64 {
	return p.TotalUncompressedSize
}

func (p *ColumnMetaData) GetTotalCompressedSize() int64 {
	return p.TotalCompressedSize
}

var ColumnMetaData_KeyValueMetadata_DEFAULT []*KeyValue

func (p *ColumnMetaData) GetKeyValueMetadata() []*KeyValue {
	return p.KeyValueMetadata
}

func (p *ColumnMetaData) GetDataPageOffset() int64 {
	return p.DataPageOffset
}

var ColumnMetaData_IndexPageOffset_DEFAULT int64

func (p *ColumnMetaData) GetIndexPageOffset() int64 {
	if !p.IsSetIndexPageOffset() {
		return ColumnMetaData_IndexPageOffset_DEFAULT
	}
	return *p.IndexPageOffset
}

var ColumnMetaData_DictionaryPageOffset_DEFAULT int64

func (p *ColumnMetaData) GetDictionaryPageOffset() int64 {
	if !p.IsSetDictionaryPageOffset() {
		return ColumnMetaData_DictionaryPageOffset_DEFAULT
	}
	return *p.DictionaryPageOffset
}

var ColumnMetaData_Statistics_DEFAULT *Statistics

func (p *ColumnMetaData) GetStatistics() *Statistics {
	if !p.IsSetStatistics() {
		return ColumnMetaData_Statistics_DEFAULT
	}
	return p.Statistics
}

var ColumnMetaData_EncodingStats_DEFAULT []*PageEncodingStats

func (p *ColumnMetaData) GetEncodingStats() []*PageEncodingStats {
	return p.EncodingStats
}
func (p *ColumnMetaData) IsSetKeyValueMetadata() bool {
	return p.KeyValueMetadata != nil
}

func (p *ColumnMetaData) IsSetIndexPageOffset() bool {
	return p.IndexPageOffset != nil
}

func (p *ColumnMetaData) IsSetDictionaryPageOffset() bool {
	return p.DictionaryPageOffset != nil
}

func (p *ColumnMetaData) IsSetStatistics() bool {
	return p.Statistics != nil
}

func (p *ColumnMetaData) IsSetEncodingStats() bool {
	return p.EncodingStats != nil
}

func (p *ColumnMetaData) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetType bool = false
	var issetEncodings bool = false
	var issetPathInSchema bool = false
	var issetCodec bool = false
	var issetNumValues bool = false
	var issetTotalUncompressedSize bool = false
	var issetTotalCompressedSize bool = false
	var issetDataPageOffset bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetType = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetEncodings = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetPathInSchema = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
			issetCodec = true
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
			issetNumValues = true
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
			issetTotalUncompressedSize = true
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
			issetTotalCompressedSize = true
		case 8:
			if err := p.ReadField8(iprot); err != nil {
				return err
			}
		case 9:
			if err := p.ReadField9(iprot); err != nil {
				return err
			}
			issetDataPageOffset = true
		case 10:
			if err := p.ReadField10(iprot); err != nil {
				return err
			}
		case 11:
			if err := p.ReadField11(iprot); err != nil {
				return err
			}
		case 12:
			if err := p.ReadField12(iprot); err != nil {
				return err
			}
		case 13:
			if err := p.ReadField13(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetType {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Type is not set"))
	}
	if !issetEncodings {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Encodings is not set"))
	}
	if !issetPathInSchema {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field PathInSchema is not set"))
	}
	if !issetCodec {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Codec is not set"))
	}
	if !issetNumValues {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumValues is not set"))
	}
	if !issetTotalUncompressedSize {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TotalUncompressedSize is not set"))
	}
	if !issetTotalCompressedSize {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TotalCompressedSize is not set"))
	}
	if !issetDataPageOffset {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field DataPageOffset is not set"))
	}
	return nil
}

func (p *ColumnMetaData) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := Type(v)
		p.Type = temp
	}
	return nil
}

func (p *ColumnMetaData) ReadField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]Encoding, 0, size)
	p.Encodings = tSlice
	for i := 0; i < size; i++ {
		var _elem0 Encoding
		if v, err := iprot.ReadI32(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			temp := Encoding(v)
			_elem0 = temp
		}
		p.Encodings = append(p.Encodings, _elem0)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnMetaData) ReadField3(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]string, 0, size)
	p.PathInSchema = tSlice
	for i := 0; i < size; i++ {
		var _elem1 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem1 = v
		}
		p.PathInSchema = append(p.PathInSchema, _elem1)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnMetaData) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		temp := CompressionCodec(v)
		p.Codec = temp
	}
	return nil
}

func (p *ColumnMetaData) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.NumValues = v
	}
	return nil
}

func (p *ColumnMetaData) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.TotalUncompressedSize = v
	}
	return nil
}

func (p *ColumnMetaData) ReadField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 7: ", err)
	} else {
		p.TotalCompressedSize = v
	}
	return nil
}

func (p *ColumnMetaData) ReadField8(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*KeyValue, 0, size)
	p.KeyValueMetadata = tSlice
	for i := 0; i < size; i++ {
		_elem2 := &KeyValue{}
		if err := _elem2.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem2), err)
		}
		p.KeyValueMetadata = append(p.KeyValueMetadata, _elem2)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnMetaData) ReadField9(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 9: ", err)
	} else {
		p.DataPageOffset = v
	}
	return nil
}

func (p *ColumnMetaData) ReadField10(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 10: ", err)
	} else {
		p.IndexPageOffset = &v
	}
	return nil
}

func (p *ColumnMetaData) ReadField11(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 11: ", err)
	} else {
		p.DictionaryPageOffset = &v
	}
	return nil
}

func (p *ColumnMetaData) ReadField12(iprot thrift.TProtocol) error {
	p.Statistics = &Statistics{}
	if err := p.Statistics.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Statistics), err)
	}
	return nil
}

func (p *ColumnMetaData) ReadField13(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*PageEncodingStats, 0, size)
	p.EncodingStats = tSlice
	for i := 0; i < size; i++ {
		_elem3 := &PageEncodingStats{}
		if err := _elem3.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem3), err)
		}
		p.EncodingStats = append(p.EncodingStats, _elem3)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnMetaData) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ColumnMetaData"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
		if err := p.writeField8(oprot); err != nil {
			return err
		}
		if err := p.writeField9(oprot); err != nil {
			return err
		}
		if err := p.writeField10(oprot); err != nil {
			return err
		}
		if err := p.writeField11(oprot); err != nil {
			return err
		}
		if err := p.writeField12(oprot); err != nil {
			return err
		}
		if err := p.writeField13(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ColumnMetaData) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("type", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:type: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Type)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.type (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:type: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("encodings", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:encodings: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.I32, len(p.Encodings)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Encodings {
		if err := oprot.WriteI32(int32(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:encodings: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("path_in_schema", thrift.LIST, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:path_in_schema: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.PathInSchema)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.PathInSchema {
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:path_in_schema: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("codec", thrift.I32, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:codec: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Codec)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.codec (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:codec: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_values", thrift.I64, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:num_values: ", p), err)
	}
	if err := oprot.WriteI64(p.NumValues); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_values (5) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:num_values: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("total_uncompressed_size", thrift.I64, 6); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:total_uncompressed_size: ", p), err)
	}
	if err := oprot.WriteI64(p.TotalUncompressedSize); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.total_uncompressed_size (6) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 6:total_uncompressed_size: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField7(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("total_compressed_size", thrift.I64, 7); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:total_compressed_size: ", p), err)
	}
	if err := oprot.WriteI64(p.TotalCompressedSize); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.total_compressed_size (7) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 7:total_compressed_size: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField8(oprot thrift.TProtocol) (err error) {
	if p.IsSetKeyValueMetadata() {
		if err := oprot.WriteFieldBegin("key_value_metadata", thrift.LIST, 8); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:key_value_metadata: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.KeyValueMetadata)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.KeyValueMetadata {
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 8:key_value_metadata: ", p), err)
		}
	}
	return err
}

func (p *ColumnMetaData) writeField9(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("data_page_offset", thrift.I64, 9); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 9:data_page_offset: ", p), err)
	}
	if err := oprot.WriteI64(p.DataPageOffset); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.data_page_offset (9) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 9:data_page_offset: ", p), err)
	}
	return err
}

func (p *ColumnMetaData) writeField10(oprot thrift.TProtocol) (err error) {
	if p.IsSetIndexPageOffset() {
		if err := oprot.WriteFieldBegin("index_page_offset", thrift.I64, 10); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:index_page_offset: ", p), err)
		}
		if err := oprot.WriteI64(*p.IndexPageOffset); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.index_page_offset (10) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 10:index_page_offset: ", p), err)
		}
	}
	return err
}

func (p *ColumnMetaData) writeField11(oprot thrift.TProtocol) (err error) {
	if p.IsSetDictionaryPageOffset() {
		if err := oprot.WriteFieldBegin("dictionary_page_offset", thrift.I64, 11); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 11:dictionary_page_offset: ", p), err)
		}
		if err := oprot.WriteI64(*p.DictionaryPageOffset); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.dictionary_page_offset (11) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 11:dictionary_page_offset: ", p), err)
		}
	}
	return err
}

func (p *ColumnMetaData) writeField12(oprot thrift.TProtocol) (err error) {
	if p.IsSetStatistics() {
		if err := oprot.WriteFieldBegin("statistics", thrift.STRUCT, 12); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 12:statistics: ", p), err)
		}
		if err := p.Statistics.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Statistics), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 12:statistics: ", p), err)
		}
	}
	return err
}

func (p *ColumnMetaData) writeField13(oprot thrift.TProtocol) (err error) {
	if p.IsSetEncodingStats() {
		if err := oprot.WriteFieldBegin("encoding_stats", thrift.LIST, 13); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 13:encoding_stats: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.EncodingStats)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.EncodingStats {
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 13:encoding_stats: ", p), err)
		}
	}
	return err
}

func (p *ColumnMetaData) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnMetaData(%+v)", *p)
}

// Attributes:
//  - FilePath: File where column data is stored.  If not set, assumed to be same file as
// metadata.  This path is relative to the current file.
//
//  - FileOffset: Byte offset in file_path to the ColumnMetaData *
//  - MetaData: Column metadata for this chunk. This is the same content as what is at
// file_path/file_offset.  Having it here has it replicated in the file
// metadata.
//
//  - OffsetIndexOffset: File offset of ColumnChunk's OffsetIndex *
//  - OffsetIndexLength: Size of ColumnChunk's OffsetIndex, in bytes *
//  - ColumnIndexOffset: File offset of ColumnChunk's ColumnIndex *
//  - ColumnIndexLength: Size of ColumnChunk's ColumnIndex, in bytes *
type ColumnChunk struct {
	FilePath          *string         `thrift:"file_path,1" db:"file_path" json:"file_path,omitempty"`
	FileOffset        int64           `thrift:"file_offset,2,required" db:"file_offset" json:"file_offset"`
	MetaData          *ColumnMetaData `thrift:"meta_data,3" db:"meta_data" json:"meta_data,omitempty"`
	OffsetIndexOffset *int64          `thrift:"offset_index_offset,4" db:"offset_index_offset" json:"offset_index_offset,omitempty"`
	OffsetIndexLength *int32          `thrift:"offset_index_length,5" db:"offset_index_length" json:"offset_index_length,omitempty"`
	ColumnIndexOffset *int64          `thrift:"column_index_offset,6" db:"column_index_offset" json:"column_index_offset,omitempty"`
	ColumnIndexLength *int32          `thrift:"column_index_length,7" db:"column_index_length" json:"column_index_length,omitempty"`
}

func NewColumnChunk() *ColumnChunk {
	return &ColumnChunk{}
}

var ColumnChunk_FilePath_DEFAULT string

func (p *ColumnChunk) GetFilePath() string {
	if !p.IsSetFilePath() {
		return ColumnChunk_FilePath_DEFAULT
	}
	return *p.FilePath
}

func (p *ColumnChunk) GetFileOffset() int64 {
	return p.FileOffset
}

var ColumnChunk_MetaData_DEFAULT *ColumnMetaData

func (p *ColumnChunk) GetMetaData() *ColumnMetaData {
	if !p.IsSetMetaData() {
		return ColumnChunk_MetaData_DEFAULT
	}
	return p.MetaData
}

var ColumnChunk_OffsetIndexOffset_DEFAULT int64

func (p *ColumnChunk) GetOffsetIndexOffset() int64 {
	if !p.IsSetOffsetIndexOffset() {
		return ColumnChunk_OffsetIndexOffset_DEFAULT
	}
	return *p.OffsetIndexOffset
}

var ColumnChunk_OffsetIndexLength_DEFAULT int32

func (p *ColumnChunk) GetOffsetIndexLength() int32 {
	if !p.IsSetOffsetIndexLength() {
		return ColumnChunk_OffsetIndexLength_DEFAULT
	}
	return *p.OffsetIndexLength
}

var ColumnChunk_ColumnIndexOffset_DEFAULT int64

func (p *ColumnChunk) GetColumnIndexOffset() int64 {
	if !p.IsSetColumnIndexOffset() {
		return ColumnChunk_ColumnIndexOffset_DEFAULT
	}
	return *p.ColumnIndexOffset
}

var ColumnChunk_ColumnIndexLength_DEFAULT int32

func (p *ColumnChunk) GetColumnIndexLength() int32 {
	if !p.IsSetColumnIndexLength() {
		return ColumnChunk_ColumnIndexLength_DEFAULT
	}
	return *p.ColumnIndexLength
}
func (p *ColumnChunk) IsSetFilePath() bool {
	return p.FilePath != nil
}

func (p *ColumnChunk) IsSetMetaData() bool {
	return p.MetaData != nil
}

func (p *ColumnChunk) IsSetOffsetIndexOffset() bool {
	return p.OffsetIndexOffset != nil
}

func (p *ColumnChunk) IsSetOffsetIndexLength() bool {
	return p.OffsetIndexLength != nil
}

func (p *ColumnChunk) IsSetColumnIndexOffset() bool {
	return p.ColumnIndexOffset != nil
}

func (p *ColumnChunk) IsSetColumnIndexLength() bool {
	return p.ColumnIndexLength != nil
}

func (p *ColumnChunk) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetFileOffset bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetFileOffset = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetFileOffset {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field FileOffset is not set"))
	}
	return nil
}

func (p *ColumnChunk) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.FilePath = &v
	}
	return nil
}

func (p *ColumnChunk) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.FileOffset = v
	}
	return nil
}

func (p *ColumnChunk) ReadField3(iprot thrift.TProtocol) error {
	p.MetaData = &ColumnMetaData{}
	if err := p.MetaData.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.MetaData), err)
	}
	return nil
}

func (p *ColumnChunk) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.OffsetIndexOffset = &v
	}
	return nil
}

func (p *ColumnChunk) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.OffsetIndexLength = &v
	}
	return nil
}

func (p *ColumnChunk) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.ColumnIndexOffset = &v
	}
	return nil
}

func (p *ColumnChunk) ReadField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 7: ", err)
	} else {
		p.ColumnIndexLength = &v
	}
	return nil
}

func (p *ColumnChunk) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ColumnChunk"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ColumnChunk) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetFilePath() {
		if err := oprot.WriteFieldBegin("file_path", thrift.STRING, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:file_path: ", p), err)
		}
		if err := oprot.WriteString(string(*p.FilePath)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.file_path (1) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:file_path: ", p), err)
		}
	}
	return err
}

func (p *ColumnChunk) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("file_offset", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:file_offset: ", p), err)
	}
	if err := oprot.WriteI64(p.FileOffset); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.file_offset (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:file_offset: ", p), err)
	}
	return err
}

func (p *ColumnChunk) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetMetaData() {
		if err := oprot.WriteFieldBegin("meta_data", thrift.STRUCT, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:meta_data: ", p), err)
		}
		if err := p.MetaData.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.MetaData), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:meta_data: ", p), err)
		}
	}
	return err
}

func (p *ColumnChunk) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetOffsetIndexOffset() {
		if err := oprot.WriteFieldBegin("offset_index_offset", thrift.I64, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:offset_index_offset: ", p), err)
		}
		if err := oprot.WriteI64(*p.OffsetIndexOffset); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.offset_index_offset (4) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:offset_index_offset: ", p), err)
		}
	}
	return err
}

func (p *ColumnChunk) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetOffsetIndexLength() {
		if err := oprot.WriteFieldBegin("offset_index_length", thrift.I32, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:offset_index_length: ", p), err)
		}
		if err := oprot.WriteI32(*p.OffsetIndexLength); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.offset_index_length (5) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:offset_index_length: ", p), err)
		}
	}
	return err
}

func (p *ColumnChunk) writeField6(oprot thrift.TProtocol) (err error) {
	if p.IsSetColumnIndexOffset() {
		if err := oprot.WriteFieldBegin("column_index_offset", thrift.I64, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:column_index_offset: ", p), err)
		}
		if err := oprot.WriteI64(*p.ColumnIndexOffset); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.column_index_offset (6) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:column_index_offset: ", p), err)
		}
	}
	return err
}

func (p *ColumnChunk) writeField7(oprot thrift.TProtocol) (err error) {
	if p.IsSetColumnIndexLength() {
		if err := oprot.WriteFieldBegin("column_index_length", thrift.I32, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:column_index_length: ", p), err)
		}
		if err := oprot.WriteI32(*p.ColumnIndexLength); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.column_index_length (7) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:column_index_length: ", p), err)
		}
	}
	return err
}

func (p *ColumnChunk) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnChunk(%+v)", *p)
}

// Attributes:
//  - Columns: Metadata for each column chunk in this row group.
// This list must have the same order as the SchemaElement list in FileMetaData.
//
//  - TotalByteSize: Total byte size of all the uncompressed column data in this row group *
//  - NumRows: Number of rows in this row group *
//  - SortingColumns: If set, specifies a sort ordering of the rows in this RowGroup.
// The sorting columns can be a subset of all the columns.
type RowGroup struct {
	Columns        []*ColumnChunk   `thrift:"columns,1,required" db:"columns" json:"columns"`
	TotalByteSize  int64            `thrift:"total_byte_size,2,required" db:"total_byte_size" json:"total_byte_size"`
	NumRows        int64            `thrift:"num_rows,3,required" db:"num_rows" json:"num_rows"`
	SortingColumns []*SortingColumn `thrift:"sorting_columns,4" db:"sorting_columns" json:"sorting_columns,omitempty"`
}

func NewRowGroup() *RowGroup {
	return &RowGroup{}
}

func (p *RowGroup) GetColumns() []*ColumnChunk {
	return p.Columns
}

func (p *RowGroup) GetTotalByteSize() int64 {
	return p.TotalByteSize
}

func (p *RowGroup) GetNumRows() int64 {
	return p.NumRows
}

var RowGroup_SortingColumns_DEFAULT []*SortingColumn

func (p *RowGroup) GetSortingColumns() []*SortingColumn {
	return p.SortingColumns
}
func (p *RowGroup) IsSetSortingColumns() bool {
	return p.SortingColumns != nil
}

func (p *RowGroup) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetColumns bool = false
	var issetTotalByteSize bool = false
	var issetNumRows bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetColumns = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetTotalByteSize = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetNumRows = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetColumns {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Columns is not set"))
	}
	if !issetTotalByteSize {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TotalByteSize is not set"))
	}
	if !issetNumRows {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumRows is not set"))
	}
	return nil
}

func (p *RowGroup) ReadField1(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*ColumnChunk, 0, size)
	p.Columns = tSlice
	for i := 0; i < size; i++ {
		_elem4 := &ColumnChunk{}
		if err := _elem4.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem4), err)
		}
		p.Columns = append(p.Columns, _elem4)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *RowGroup) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.TotalByteSize = v
	}
	return nil
}

func (p *RowGroup) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.NumRows = v
	}
	return nil
}

func (p *RowGroup) ReadField4(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*SortingColumn, 0, size)
	p.SortingColumns = tSlice
	for i := 0; i < size; i++ {
		_elem5 := &SortingColumn{}
		if err := _elem5.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem5), err)
		}
		p.SortingColumns = append(p.SortingColumns, _elem5)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *RowGroup) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("RowGroup"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *RowGroup) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("columns", thrift.LIST, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:columns: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Columns)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Columns {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:columns: ", p), err)
	}
	return err
}

func (p *RowGroup) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("total_byte_size", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:total_byte_size: ", p), err)
	}
	if err := oprot.WriteI64(p.TotalByteSize); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.total_byte_size (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:total_byte_size: ", p), err)
	}
	return err
}

func (p *RowGroup) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_rows", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:num_rows: ", p), err)
	}
	if err := oprot.WriteI64(p.NumRows); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_rows (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:num_rows: ", p), err)
	}
	return err
}

func (p *RowGroup) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetSortingColumns() {
		if err := oprot.WriteFieldBegin("sorting_columns", thrift.LIST, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:sorting_columns: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.SortingColumns)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.SortingColumns {
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:sorting_columns: ", p), err)
		}
	}
	return err
}

func (p *RowGroup) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RowGroup(%+v)", *p)
}

// Empty struct to signal the order defined by the physical or logical type
type TypeDefinedOrder struct {
}

func NewTypeDefinedOrder() *TypeDefinedOrder {
	return &TypeDefinedOrder{}
}

func (p *TypeDefinedOrder) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TypeDefinedOrder) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("TypeDefinedOrder"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TypeDefinedOrder) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TypeDefinedOrder(%+v)", *p)
}

// Union to specify the order used for the min_value and max_value fields for a
// column. This union takes the role of an enhanced enum that allows rich
// elements (which will be needed for a collation-based ordering in the future).
//
// Possible values are:
// * TypeDefinedOrder - the column uses the order defined by its logical or
//                      physical type (if there is no logical type).
//
// If the reader does not support the value of this union, min and max stats
// for this column should be ignored.
//
// Attributes:
//  - TYPE_ORDER: The sort orders for logical types are:
//   UTF8 - unsigned byte-wise comparison
//   INT8 - signed comparison
//   INT16 - signed comparison
//   INT32 - signed comparison
//   INT64 - signed comparison
//   UINT8 - unsigned comparison
//   UINT16 - unsigned comparison
//   UINT32 - unsigned comparison
//   UINT64 - unsigned comparison
//   DECIMAL - signed comparison of the represented value
//   DATE - signed comparison
//   TIME_MILLIS - signed comparison
//   TIME_MICROS - signed comparison
//   TIMESTAMP_MILLIS - signed comparison
//   TIMESTAMP_MICROS - signed comparison
//   INTERVAL - unsigned comparison
//   JSON - unsigned byte-wise comparison
//   BSON - unsigned byte-wise comparison
//   ENUM - unsigned byte-wise comparison
//   LIST - undefined
//   MAP - undefined
//
// In the absence of logical types, the sort order is determined by the physical type:
//   BOOLEAN - false, true
//   INT32 - signed comparison
//   INT64 - signed comparison
//   INT96 (only used for legacy timestamps) - undefined
//   FLOAT - signed comparison of the represented value (*)
//   DOUBLE - signed comparison of the represented value (*)
//   BYTE_ARRAY - unsigned byte-wise comparison
//   FIXED_LEN_BYTE_ARRAY - unsigned byte-wise comparison
//
// (*) Because the sorting order is not specified properly for floating
//     point values (relations vs. total ordering) the following
//     compatibility rules should be applied when reading statistics:
//     - If the min is a NaN, it should be ignored.
//     - If the max is a NaN, it should be ignored.
//     - If the min is +0, the row group may contain -0 values as well.
//     - If the max is -0, the row group may contain +0 values as well.
//     - When looking for NaN values, min and max should be ignored.
type ColumnOrder struct {
	TYPE_ORDER *TypeDefinedOrder `thrift:"TYPE_ORDER,1" db:"TYPE_ORDER" json:"TYPE_ORDER,omitempty"`
}

func NewColumnOrder() *ColumnOrder {
	return &ColumnOrder{}
}

var ColumnOrder_TYPE_ORDER_DEFAULT *TypeDefinedOrder

func (p *ColumnOrder) GetTYPE_ORDER() *TypeDefinedOrder {
	if !p.IsSetTYPE_ORDER() {
		return ColumnOrder_TYPE_ORDER_DEFAULT
	}
	return p.TYPE_ORDER
}
func (p *ColumnOrder) CountSetFieldsColumnOrder() int {
	count := 0
	if p.IsSetTYPE_ORDER() {
		count++
	}
	return count

}

func (p *ColumnOrder) IsSetTYPE_ORDER() bool {
	return p.TYPE_ORDER != nil
}

func (p *ColumnOrder) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *ColumnOrder) ReadField1(iprot thrift.TProtocol) error {
	p.TYPE_ORDER = &TypeDefinedOrder{}
	if err := p.TYPE_ORDER.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.TYPE_ORDER), err)
	}
	return nil
}

func (p *ColumnOrder) Write(oprot thrift.TProtocol) error {
	if c := p.CountSetFieldsColumnOrder(); c != 1 {
		return fmt.Errorf("%T write union: exactly one field must be set (%d set).", p, c)
	}
	if err := oprot.WriteStructBegin("ColumnOrder"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ColumnOrder) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetTYPE_ORDER() {
		if err := oprot.WriteFieldBegin("TYPE_ORDER", thrift.STRUCT, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:TYPE_ORDER: ", p), err)
		}
		if err := p.TYPE_ORDER.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.TYPE_ORDER), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:TYPE_ORDER: ", p), err)
		}
	}
	return err
}

func (p *ColumnOrder) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnOrder(%+v)", *p)
}

// Attributes:
//  - Offset: Offset of the page in the file *
//  - CompressedPageSize: Size of the page, including header. Sum of compressed_page_size and header
// length
//  - FirstRowIndex: Index within the RowGroup of the first row of the page; this means pages
// change on record boundaries (r = 0).
type PageLocation struct {
	Offset             int64 `thrift:"offset,1,required" db:"offset" json:"offset"`
	CompressedPageSize int32 `thrift:"compressed_page_size,2,required" db:"compressed_page_size" json:"compressed_page_size"`
	FirstRowIndex      int64 `thrift:"first_row_index,3,required" db:"first_row_index" json:"first_row_index"`
}

func NewPageLocation() *PageLocation {
	return &PageLocation{}
}

func (p *PageLocation) GetOffset() int64 {
	return p.Offset
}

func (p *PageLocation) GetCompressedPageSize() int32 {
	return p.CompressedPageSize
}

func (p *PageLocation) GetFirstRowIndex() int64 {
	return p.FirstRowIndex
}
func (p *PageLocation) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetOffset bool = false
	var issetCompressedPageSize bool = false
	var issetFirstRowIndex bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetOffset = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetCompressedPageSize = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetFirstRowIndex = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetOffset {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Offset is not set"))
	}
	if !issetCompressedPageSize {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field CompressedPageSize is not set"))
	}
	if !issetFirstRowIndex {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field FirstRowIndex is not set"))
	}
	return nil
}

func (p *PageLocation) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Offset = v
	}
	return nil
}

func (p *PageLocation) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.CompressedPageSize = v
	}
	return nil
}

func (p *PageLocation) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.FirstRowIndex = v
	}
	return nil
}

func (p *PageLocation) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("PageLocation"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *PageLocation) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("offset", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:offset: ", p), err)
	}
	if err := oprot.WriteI64(p.Offset); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.offset (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:offset: ", p), err)
	}
	return err
}

func (p *PageLocation) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("compressed_page_size", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:compressed_page_size: ", p), err)
	}
	if err := oprot.WriteI32(p.CompressedPageSize); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.compressed_page_size (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:compressed_page_size: ", p), err)
	}
	return err
}

func (p *PageLocation) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("first_row_index", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:first_row_index: ", p), err)
	}
	if err := oprot.WriteI64(p.FirstRowIndex); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.first_row_index (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:first_row_index: ", p), err)
	}
	return err
}

func (p *PageLocation) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PageLocation(%+v)", *p)
}

// Attributes:
//  - PageLocations: PageLocations, ordered by increasing PageLocation.offset. It is required
// that page_locations[i].first_row_index < page_locations[i+1].first_row_index.
type OffsetIndex struct {
	PageLocations []*PageLocation `thrift:"page_locations,1,required" db:"page_locations" json:"page_locations"`
}

func NewOffsetIndex() *OffsetIndex {
	return &OffsetIndex{}
}

func (p *OffsetIndex) GetPageLocations() []*PageLocation {
	return p.PageLocations
}
func (p *OffsetIndex) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetPageLocations bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetPageLocations = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetPageLocations {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field PageLocations is not set"))
	}
	return nil
}

func (p *OffsetIndex) ReadField1(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*PageLocation, 0, size)
	p.PageLocations = tSlice
	for i := 0; i < size; i++ {
		_elem6 := &PageLocation{}
		if err := _elem6.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem6), err)
		}
		p.PageLocations = append(p.PageLocations, _elem6)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *OffsetIndex) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("OffsetIndex"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *OffsetIndex) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("page_locations", thrift.LIST, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:page_locations: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.PageLocations)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.PageLocations {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:page_locations: ", p), err)
	}
	return err
}

func (p *OffsetIndex) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("OffsetIndex(%+v)", *p)
}

// Description for ColumnIndex.
// Each <array-field>[i] refers to the page at OffsetIndex.page_locations[i]
//
// Attributes:
//  - NullPages: A list of Boolean values to determine the validity of the corresponding
// min and max values. If true, a page contains only null values, and writers
// have to set the corresponding entries in min_values and max_values to
// byte[0], so that all lists have the same length. If false, the
// corresponding entries in min_values and max_values must be valid.
//  - MinValues: Two lists containing lower and upper bounds for the values of each page.
// These may be the actual minimum and maximum values found on a page, but
// can also be (more compact) values that do not exist on a page. For
// example, instead of storing ""Blart Versenwald III", a writer may set
// min_values[i]="B", max_values[i]="C". Such more compact values must still
// be valid values within the column's logical type. Readers must make sure
// that list entries are populated before using them by inspecting null_pages.
//  - MaxValues
//  - BoundaryOrder: Stores whether both min_values and max_values are orderd and if so, in
// which direction. This allows readers to perform binary searches in both
// lists. Readers cannot assume that max_values[i] <= min_values[i+1], even
// if the lists are ordered.
//  - NullCounts: A list containing the number of null values for each page *
type ColumnIndex struct {
	NullPages     []bool        `thrift:"null_pages,1,required" db:"null_pages" json:"null_pages"`
	MinValues     [][]byte      `thrift:"min_values,2,required" db:"min_values" json:"min_values"`
	MaxValues     [][]byte      `thrift:"max_values,3,required" db:"max_values" json:"max_values"`
	BoundaryOrder BoundaryOrder `thrift:"boundary_order,4,required" db:"boundary_order" json:"boundary_order"`
	NullCounts    []int64       `thrift:"null_counts,5" db:"null_counts" json:"null_counts,omitempty"`
}

func NewColumnIndex() *ColumnIndex {
	return &ColumnIndex{}
}

func (p *ColumnIndex) GetNullPages() []bool {
	return p.NullPages
}

func (p *ColumnIndex) GetMinValues() [][]byte {
	return p.MinValues
}

func (p *ColumnIndex) GetMaxValues() [][]byte {
	return p.MaxValues
}

func (p *ColumnIndex) GetBoundaryOrder() BoundaryOrder {
	return p.BoundaryOrder
}

var ColumnIndex_NullCounts_DEFAULT []int64

func (p *ColumnIndex) GetNullCounts() []int64 {
	return p.NullCounts
}
func (p *ColumnIndex) IsSetNullCounts() bool {
	return p.NullCounts != nil
}

func (p *ColumnIndex) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetNullPages bool = false
	var issetMinValues bool = false
	var issetMaxValues bool = false
	var issetBoundaryOrder bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetNullPages = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetMinValues = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetMaxValues = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
			issetBoundaryOrder = true
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetNullPages {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NullPages is not set"))
	}
	if !issetMinValues {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field MinValues is not set"))
	}
	if !issetMaxValues {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field MaxValues is not set"))
	}
	if !issetBoundaryOrder {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field BoundaryOrder is not set"))
	}
	return nil
}

func (p *ColumnIndex) ReadField1(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]bool, 0, size)
	p.NullPages = tSlice
	for i := 0; i < size; i++ {
		var _elem7 bool
		if v, err := iprot.ReadBool(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem7 = v
		}
		p.NullPages = append(p.NullPages, _elem7)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnIndex) ReadField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([][]byte, 0, size)
	p.MinValues = tSlice
	for i := 0; i < size; i++ {
		var _elem8 []byte
		if v, err := iprot.ReadBinary(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem8 = v
		}
		p.MinValues = append(p.MinValues, _elem8)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnIndex) ReadField3(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([][]byte, 0, size)
	p.MaxValues = tSlice
	for i := 0; i < size; i++ {
		var _elem9 []byte
		if v, err := iprot.ReadBinary(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem9 = v
		}
		p.MaxValues = append(p.MaxValues, _elem9)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnIndex) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		temp := BoundaryOrder(v)
		p.BoundaryOrder = temp
	}
	return nil
}

func (p *ColumnIndex) ReadField5(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]int64, 0, size)
	p.NullCounts = tSlice
	for i := 0; i < size; i++ {
		var _elem10 int64
		if v, err := iprot.ReadI64(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem10 = v
		}
		p.NullCounts = append(p.NullCounts, _elem10)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ColumnIndex) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ColumnIndex"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ColumnIndex) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("null_pages", thrift.LIST, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:null_pages: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.BOOL, len(p.NullPages)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.NullPages {
		if err := oprot.WriteBool(v); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:null_pages: ", p), err)
	}
	return err
}

func (p *ColumnIndex) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("min_values", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:min_values: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.MinValues)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.MinValues {
		if err := oprot.WriteBinary(v); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:min_values: ", p), err)
	}
	return err
}

func (p *ColumnIndex) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("max_values", thrift.LIST, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:max_values: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.MaxValues)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.MaxValues {
		if err := oprot.WriteBinary(v); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:max_values: ", p), err)
	}
	return err
}

func (p *ColumnIndex) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("boundary_order", thrift.I32, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:boundary_order: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.BoundaryOrder)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.boundary_order (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:boundary_order: ", p), err)
	}
	return err
}

func (p *ColumnIndex) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetNullCounts() {
		if err := oprot.WriteFieldBegin("null_counts", thrift.LIST, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:null_counts: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.I64, len(p.NullCounts)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.NullCounts {
			if err := oprot.WriteI64(v); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:null_counts: ", p), err)
		}
	}
	return err
}

func (p *ColumnIndex) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ColumnIndex(%+v)", *p)
}

// Description for file metadata
//
// Attributes:
//  - Version: Version of this file *
//  - Schema: Parquet schema for this file.  This schema contains metadata for all the columns.
// The schema is represented as a tree with a single root.  The nodes of the tree
// are flattened to a list by doing a depth-first traversal.
// The column metadata contains the path in the schema for that column which can be
// used to map columns to nodes in the schema.
// The first element is the root *
//  - NumRows: Number of rows in this file *
//  - RowGroups: Row groups in this file *
//  - KeyValueMetadata: Optional key/value metadata *
//  - CreatedBy: String for application that wrote this file.  This should be in the format
// <Application> version <App Version> (build <App Build Hash>).
// e.g. impala version 1.0 (build 6cf94d29b2b7115df4de2c06e2ab4326d721eb55)
//
//  - ColumnOrders: Sort order used for the min_value and max_value fields of each column in
// this file. Each sort order corresponds to one column, determined by its
// position in the list, matching the position of the column in the schema.
//
// Without column_orders, the meaning of the min_value and max_value fields is
// undefined. To ensure well-defined behavior, if min_value and max_value are
// written to a Parquet file, column_orders must be written as well.
//
// The obsolete min and max fields are always sorted by signed comparison
// regardless of column_orders.
type FileMetaData struct {
	Version          int32            `thrift:"version,1,required" db:"version" json:"version"`
	Schema           []*SchemaElement `thrift:"schema,2,required" db:"schema" json:"schema"`
	NumRows          int64            `thrift:"num_rows,3,required" db:"num_rows" json:"num_rows"`
	RowGroups        []*RowGroup      `thrift:"row_groups,4,required" db:"row_groups" json:"row_groups"`
	KeyValueMetadata []*KeyValue      `thrift:"key_value_metadata,5" db:"key_value_metadata" json:"key_value_metadata,omitempty"`
	CreatedBy        *string          `thrift:"created_by,6" db:"created_by" json:"created_by,omitempty"`
	ColumnOrders     []*ColumnOrder   `thrift:"column_orders,7" db:"column_orders" json:"column_orders,omitempty"`
}

func NewFileMetaData() *FileMetaData {
	return &FileMetaData{}
}

func (p *FileMetaData) GetVersion() int32 {
	return p.Version
}

func (p *FileMetaData) GetSchema() []*SchemaElement {
	return p.Schema
}

func (p *FileMetaData) GetNumRows() int64 {
	return p.NumRows
}

func (p *FileMetaData) GetRowGroups() []*RowGroup {
	return p.RowGroups
}

var FileMetaData_KeyValueMetadata_DEFAULT []*KeyValue

func (p *FileMetaData) GetKeyValueMetadata() []*KeyValue {
	return p.KeyValueMetadata
}

var FileMetaData_CreatedBy_DEFAULT string

func (p *FileMetaData) GetCreatedBy() string {
	if !p.IsSetCreatedBy() {
		return FileMetaData_CreatedBy_DEFAULT
	}
	return *p.CreatedBy
}

var FileMetaData_ColumnOrders_DEFAULT []*ColumnOrder

func (p *FileMetaData) GetColumnOrders() []*ColumnOrder {
	return p.ColumnOrders
}
func (p *FileMetaData) IsSetKeyValueMetadata() bool {
	return p.KeyValueMetadata != nil
}

func (p *FileMetaData) IsSetCreatedBy() bool {
	return p.CreatedBy != nil
}

func (p *FileMetaData) IsSetColumnOrders() bool {
	return p.ColumnOrders != nil
}

func (p *FileMetaData) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetVersion bool = false
	var issetSchema bool = false
	var issetNumRows bool = false
	var issetRowGroups bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetVersion = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetSchema = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetNumRows = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
			issetRowGroups = true
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetVersion {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Version is not set"))
	}
	if !issetSchema {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Schema is not set"))
	}
	if !issetNumRows {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field NumRows is not set"))
	}
	if !issetRowGroups {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RowGroups is not set"))
	}
	return nil
}

func (p *FileMetaData) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Version = v
	}
	return nil
}

func (p *FileMetaData) ReadField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*SchemaElement, 0, size)
	p.Schema = tSlice
	for i := 0; i < size; i++ {
		_elem11 := &SchemaElement{}
		if err := _elem11.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem11), err)
		}
		p.Schema = append(p.Schema, _elem11)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *FileMetaData) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.NumRows = v
	}
	return nil
}

func (p *FileMetaData) ReadField4(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*RowGroup, 0, size)
	p.RowGroups = tSlice
	for i := 0; i < size; i++ {
		_elem12 := &RowGroup{}
		if err := _elem12.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem12), err)
		}
		p.RowGroups = append(p.RowGroups, _elem12)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *FileMetaData) ReadField5(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*KeyValue, 0, size)
	p.KeyValueMetadata = tSlice
	for i := 0; i < size; i++ {
		_elem13 := &KeyValue{}
		if err := _elem13.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem13), err)
		}
		p.KeyValueMetadata = append(p.KeyValueMetadata, _elem13)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *FileMetaData) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.CreatedBy = &v
	}
	return nil
}

func (p *FileMetaData) ReadField7(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*ColumnOrder, 0, size)
	p.ColumnOrders = tSlice
	for i := 0; i < size; i++ {
		_elem14 := &ColumnOrder{}
		if err := _elem14.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem14), err)
		}
		p.ColumnOrders = append(p.ColumnOrders, _elem14)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *FileMetaData) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("FileMetaData"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
			return err
		}
		if err := p.writeField5(oprot); err != nil {
			return err
		}
		if err := p.writeField6(oprot); err != nil {
			return err
		}
		if err := p.writeField7(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *FileMetaData) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("version", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:version: ", p), err)
	}
	if err := oprot.WriteI32(p.Version); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.version (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:version: ", p), err)
	}
	return err
}

func (p *FileMetaData) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("schema", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:schema: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Schema)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Schema {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:schema: ", p), err)
	}
	return err
}

func (p *FileMetaData) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num_rows", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:num_rows: ", p), err)
	}
	if err := oprot.WriteI64(p.NumRows); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.num_rows (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:num_rows: ", p), err)
	}
	return err
}

func (p *FileMetaData) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("row_groups", thrift.LIST, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:row_groups: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.RowGroups)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.RowGroups {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:row_groups: ", p), err)
	}
	return err
}

func (p *FileMetaData) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetKeyValueMetadata() {
		if err := oprot.WriteFieldBegin("key_value_metadata", thrift.LIST, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:key_value_metadata: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.KeyValueMetadata)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.KeyValueMetadata {
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:key_value_metadata: ", p), err)
		}
	}
	return err
}

func (p *FileMetaData) writeField6(oprot thrift.TProtocol) (err error) {
	if p.IsSetCreatedBy() {
		if err := oprot.WriteFieldBegin("created_by", thrift.STRING, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:created_by: ", p), err)
		}
		if err := oprot.WriteString(string(*p.CreatedBy)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.created_by (6) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:created_by: ", p), err)
		}
	}
	return err
}

func (p *FileMetaData) writeField7(oprot thrift.TProtocol) (err error) {
	if p.IsSetColumnOrders() {
		if err := oprot.WriteFieldBegin("column_orders", thrift.LIST, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:column_orders: ", p), err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.ColumnOrders)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.ColumnOrders {
			if err := v.Write(oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:column_orders: ", p), err)
		}
	}
	return err
}

func (p *FileMetaData) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("FileMetaData(%+v)", *p)
}
