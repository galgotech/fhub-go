name: "test"
specVersion: "1.0"
version: "v1"
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