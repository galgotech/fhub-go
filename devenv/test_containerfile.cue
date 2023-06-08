name: "test"
specVersion: "1.0"
version: "v1"
import: [
  "fhub/internaltest.cue"
]
build: {
  container: {
    containerFile: "Containerfile"
    source: "/app"
  }
}
serving: {
  http: {
    url: "https://fhub.dev/test"
  }
}
functions: {
  test: {
    input: {
      arg0: string
      arg1: string
    }
    output: {
      ok: bool
    }
  }
}