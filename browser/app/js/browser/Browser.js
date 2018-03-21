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
import SideBar from "./SideBar"
import MainContent from "./MainContent"
import AlertContainer from "../alert/AlertContainer"

const Aux = props => props.children

class Browser extends React.Component {
  render() {
    return (
      <Aux>
        <SideBar />
        <MainContent />
        <AlertContainer />
      </Aux>
    )
  }
}

export default connect(state => state)(Browser)
