// Copyright 2023 The fhub-runtime-go Authors
// This file is part of fhub-runtime-go.
//
// This file is part of fhub-runtime-go.
// fhub-runtime-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// fhub-runtime-go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with fhub-runtime-go. If not, see <https://www.gnu.org/licenses/>.

name: string
specVersion: "1.0"
version: string
env: {
  [string]: string
}
packages: {
  [string]: {
    import: string
    launch?: string
    build: {
      container: {
        image?: string
        context?: string
        dockerfile?: string | *"Dockerfile"
        target: string
      }
    }
    serving: {
      http?: {
        url: string
      }
      grpc?: {}
    }
  }
}
functions: {
  [string]: {
    package: string
    launch: string
    inputs: {
      [string]: number | string | bool
    }
    outputs: {
      [string]: number | string | bool
    }
  }
}