/*
 * Minio Cloud Storage (C) 2018 Minio, Inc.
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

import React from "react"
import { connect } from "react-redux"
import humanize from "humanize"
import Moment from "moment"
import ObjectItem from "./ObjectItem"
import ObjectActions from "./ObjectActions"
import * as actionsObjects from "./actions"

export const ObjectContainer = ({ object, downloadObject }) => {
  const actionButtons = <ObjectActions object={object} />
  const props = {
    name: object.name,
    contentType: object.contentType,
    size: humanize.filesize(object.size),
    lastModified: Moment(object.lastModified).format("lll"),
    actionButtons: actionButtons
  }
  return <ObjectItem {...props} onClick={() => downloadObject(object.name)} />
}

const mapDispatchToProps = dispatch => {
  return {
    downloadObject: object => dispatch(actionsObjects.downloadObject(object))
  }
}

export default connect(state => state, mapDispatchToProps)(ObjectContainer)
