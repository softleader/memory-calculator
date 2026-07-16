# Certificate contract fixtures

These test-only fixtures come from `paketo-buildpacks/libjvm` tag `v1.46.0`
(commit `d0895b1355131c76a1ef2d998ea1cfcda19c1cce`):

- `truststore.p12` is the passwordless PKCS#12 truststore required by that
  runtime version's `PasswordLessPKCS12Keystore` implementation.
- `certificate.pem` is a public X.509 certificate loaded into the copied
  truststore by the contract test.

Neither fixture contains a private key or production secret.
