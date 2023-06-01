packages: {
	pkgTestInternal: {
    import: "fhub.dev/test_internal"
    launch: "startInternal"
    serving: {
      http: {
        url: "https://fhub.dev/testinternal"
      }
    }
  }
}
functions: {
  test_internal: {
    package: "pkgTestInternal"
    launch: "testInternal"
    input: {
      arg0: string
      arg1: string
    }
    output: {
      ok: bool
    }
  }
}