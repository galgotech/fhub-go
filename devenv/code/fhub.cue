name: "test"
specVersion: "1.0"
version: "v1"
constants:  {
  "const_name": "test"
}
env: [
  "NAME"
]
import: [
  "fhub/internaltest.cue"
]
build: {
  local: {
    source: "./"
  }
}
serving: {
  http: {
    url: "https://fhub.dev/test"
  }
}
functions: {
  FuncTest: {
    input: {
      a: string
    }
    output: {
      response: string
    }
  }
}