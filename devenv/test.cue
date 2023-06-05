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
packages: {
	pkgTest: {
    import: "fhub.dev/test"
    launch: "start"
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
  }
}
functions: {
  test: {
    package: "pkgTest"
    launch: "test"
    input: {
      arg0: string
      arg1: string
    }
    output: {
      ok: bool
    }
  }
}