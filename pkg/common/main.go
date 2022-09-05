package common

import (
	"bytes"
	"io"
)

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

var SAMPLE_CERT = `-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUAhdObDJnggMc26EDvd11ItBaPewwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0xOTAzMzAyMTE2MjhaFw0yMDAz
MjkyMTE2MjhaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQCsqC6OCmA9bG8qA/qo8W8afBYANo4YKrcohAY9Mg+e
rabrPxS3fVCDUEn4A7EKlHIzRuobxVvUSOK25rGC0SftQmaRhqzhtV5oL84EvKzy
7Q83gCDZbocsf+GqdZxEayFm89WNQCCMr1BIRwVY7w8pElgAPVujRisjE6MV2//z
nQwigxuNGA06fFmT+IjCy7p3owhPanC26yn0KaFSPWqsxvl5/gD21rAJpYKtOf6V
Xl85K4tjXCqJ00oPZa6y5y1eSyEb9UE225LsBFoARgcoNCcwVpRdkiXt8TD1KI7K
LNhqOxeRi3y2i9e0f1vgo6BX6xCUx14tRs5nJyfEp7Nknyg7yK8a7FZ5eazRuRlS
CrySP3BjxEZTVySW9wxDaEY250vxtcFJanccD5/C/BAHwAN6eLczgCm4JXWA6HNq
G0LNouAONj9u5XY6gfJRz3Iyh4bHwMmXH4gN9Ep6MeeXlN13B8zQ9udG15UwNELl
FJ24KhRj+hAX8uEHeBLZN4hUKzwpkBeN9r/aotZsdVFAp9Gl1RrpOKwaShuh1fY3
nT+P0CECB4ugoIkuY/IgIGNwGzdBeAZ+pCxUuq95eGKFg0wdA7aJtV11IdcxdpYU
ljJlFlluc/CSY6G7I626nrYumLK8Lt1GCmAKu7fnZ02gKGJvbk7VEjj7/M7BywHm
XwIDAQABo1MwUTAdBgNVHQ4EFgQUybknD8HpuXXis6Jq6nT2H3R6I90wHwYDVR0j
BBgwFoAUybknD8HpuXXis6Jq6nT2H3R6I90wDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEAUHlDMCMjifD8M+7MTejca9dYOXK3kujgFjKnSig2Fa4q
97n9/9I+r3AAGz9YFBPw1bmHlnE7ieV1NJuCSe+L3z1s2Cq5ulQZV1qiP3NFdjUd
J1P8TBI18aA38huX4RwAUn1sA8B0ApiwRe8f3JxgvmCx3T9JSWaUavrMNUKM8aXf
gv6FrtrxKcE5MfdS72GYqY2+zy8HQQC0aVNhkoZk6GhlVLTG1HrZ22b9zQpn+pRc
adwWZ752ZQMSUrntPYA6BK/4aVl65m47dA0k1EnK8Hz5DSEr5W+NP1m8xsw3yFKo
1NQB0RULPvMiHqb3czMKY6ORMNZOvYqyUlbErjCcAc9jbvgOzmT91GFf6u3mUgmm
z8QvMTFzuAtmhmyKSjtiLmHZvhxQBXr2fX+7Zryz139EGijSCZOspwFHZkJFkg/U
Sc0vud9uKcbklt5GluM7HaxZWPTdze03Z4wbvU1nhKyGdf4bTvM3HeqAK2sURN3Y
1uxEV4k4S++gj2UUcfxtHGRiF3YGEfzt7pmigVo/Cr73ZXzPmr3qhxDyBSvyAARl
O57+cPhsWazbxf8tHDvSJEO+bVtFY4JRoLi2TpNbug75LmBFHojgxqgI/3qogJ2Q
lKUKOFnVNnDJhVjLh5DeNYbyyU/f+xFqTzQjSyuw+4FegKbzfq7oliRnDeT0Wgs=
-----END CERTIFICATE-----`