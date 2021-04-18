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

// Package mimedb is a database of file extension to mime content-type.
// Definitions are imported from NodeJS mime-db project under MIT license.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const progTempl = `// DO NOT EDIT THIS FILE. IT IS AUTO-GENERATED BY "gen-db.go". //
/*
 * mimedb: Mime Database, (C) 2016 MinIO, Inc.
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

// Package mimedb is a database of file extension to mime content-type.
// Definitions are imported from NodeJS mime-db project under MIT license.
package mimedb

// DB - Mime is a collection of mime types with extension as key and content-type as value.
var DB = map[string]struct {
	ContentType  string
	Compressible bool
}{
{{range $extension, $entry := . }}	"{{$extension}}": {
		ContentType:  "{{$entry.ContentType}}",
		Compressible: {{$entry.Compressible}},
	},
{{end}}}
`

type mimeEntry struct {
	ContentType  string `json:"contentType"`
	Compressible bool   `json:"compresible"`
}

type mimeDB map[string]mimeEntry

//  JSON data from gobindata and parse them into extDB.
func convertDB(jsonFile string) (mimeDB, error) {
	// Structure of JSON data from mime-db project.
	type dbEntry struct {
		Source       string   `json:"source"`
		Compressible bool     `json:"compresible"`
		Extensions   []string `json:"extensions"`
	}

	// Access embedded "db.json" inside go-bindata.
	jsonDB, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	// Convert db.json into go's typed structure.
	db := make(map[string]dbEntry)
	if err := json.Unmarshal(jsonDB, &db); err != nil {
		return nil, err
	}

	mDB := make(mimeDB)

	// Generate a new database from mime-db.
	for key, val := range db {
		if len(val.Extensions) > 0 {
			/* Denormalize - each extension has its own
			unique content-type now. Looks will be fast. */
			for _, ext := range val.Extensions {
				/* Single extension type may map to
				multiple content-types. In that case,
				simply prefer the longest content-type
				to maintain some level of
				consistency. Only guarantee is,
				whatever content type is assigned, it
				is appropriate and valid type. */
				if strings.Compare(mDB[ext].ContentType, key) < 0 {
					mDB[ext] = mimeEntry{
						ContentType:  key,
						Compressible: val.Compressible,
					}
				}
			}
		}
	}
	return mDB, nil
}

func main() {
	// Take input json file from command-line".
	if len(os.Args) != 2 {
		fmt.Print("Syntax:\n\tgen-db /path/to/db.json\n")
		os.Exit(1)
	}

	// Load and convert db.json into new database with extension
	// as key.
	mDB, err := convertDB(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Generate db embedded go program.
	tmpl := template.New("mimedb")
	mimeTmpl, err := tmpl.Parse(progTempl)
	if err != nil {
		panic(err)
	}

	err = mimeTmpl.Execute(os.Stdout, mDB)
	if err != nil {
		panic(err)
	}
}
