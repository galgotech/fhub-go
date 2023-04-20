name: "test"
specVersion: "1.0"
version: "v1"
env:  {
  name: "value"
}
packages: {
	pkgTest: {
    import: "fhub.dev/test"
    launch: "start"
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