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
  test: {
    package: "pkgTest"
    launch: "FuncTest"
    input: {
      arg0: string
      arg1: string | int | null | bool | float | number | bytes | *""
      arg2: {[string]: string | int}
      arg3: [...(string | int)]
      arg4: float
    }
    output: {
      ok: bool
    }
  }
}