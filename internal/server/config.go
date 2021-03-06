package server

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"regexp"
)

var SERVER_CERT = []byte(`-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUa/FoCcBpHGnYD/1C6OSRumK+76MwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjAzMjgyMTQ1MzBaFw0yMzAz
MjgyMTQ1MzBaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQDDyyF/a6xc1UCi45D94w1uMOaSJUwaCMfxTPR/HqFg
w3g69/Pa/fgRyom9MY5m6PwmgVwS+4RKcpcQjuAE86OxubUUir8pCt3RmFQrwfvq
GOEYL5GUqSh4AG6nfkWV3NDegIfV3GmEGIkEGrYsIOj8DeXLTpXOAcaM3cKPsfVk
lpjqSNcgJqDJtvPqNRjYVhn65ey4TD94dGO/3t95QbNFibV+s08JmE0RtESdEC28
hvC/gAqokolUX2L9yHZw9wMVs1pP5uleUvNI8OvZ/aeUkwYCse2qvb+BTqHW9YjX
kT3kaI1u+j2xku7SjpDyc5E11w2Z3AoP8uVyp6HiY7M0K5oApH8r2YVtpypfCuFB
owko82L7CPAeR2HUi+9KZgFVp3dABb+qHm1rmleNkGteQixHnvNFANmt+IlMBLQx
skmoYw1kq/INFCI5f8I9NzUqMIyscbLffXnvKAaIkqq2gclIFALrPZx9Y456jI03
15ZZP3aZXe22vM5OX0P9S9ov8lHlexQNVjwd0WijkHRHowwtvsoYFiub0g8Q85tp
Dgc+F5kY2J+gu7+KTO3YaM4m31y+bdwDkK/FjisRM5HwsHT98vupMDqOrgHq1FFf
JtmBw+FlLNMYkvDQ6hyngg+DeFTEb9/9tSm8/hbvG6nIq+1L6R4UHjOyXbsht5w/
RwIDAQABo1MwUTAdBgNVHQ4EFgQUywzdUNpAeMQVH07f2fMQaD3tlvIwHwYDVR0j
BBgwFoAUywzdUNpAeMQVH07f2fMQaD3tlvIwDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEARf/nJvyLfQ51d1/6dDezN9PtZkJsf+dPrabnOaulf9Wc
T9aSn0Mx79U5w5pIkPzbAPqXKYU8EwAohv4vB3ngJpa/VKMEjJ79leXo2P0849rK
OJC8MVu5aKnCSZME20gKV2+wMIdyjslqznx304yB1cm7MM77Yejy2vX0wyj/11sV
A60gdkMJLCk13WfRpX/t8cqlhYaQww4HhRfJezJsF/ShAcxX0P/amvEMv3DgTTlU
cKL4/h4hN1fu3v6BqeHF6FdKkDpmhauQ8C/hCh4vCX6p3lmVmW/Cu0JvptMPVHuR
JrPerUiKuyYxMQKwTkIDs3MtoE/oAAOOcPf1LoTI3F8eUJ+TQbtTfiY8ppFkP/Ju
sHIM664TWD0KABUURSzikGGv4UU7ao1xXiYcjeqRmDZuD6siI8FeRCJaQff1vOaJ
9bgf4eCFxDsgQwHzEagLA2kb3IUvV7wdGDk4kau5I1ELCDj+sWcUG3HRzptZfmxy
ZtYHOcHr5Plv1Aj/c4HHxKBKHcoj8bahVSdozgtLi+GEmpC2Idc+QlBTSESk5nw7
7VxB9kmRJiY7AhHg/5LITk9uqdtJZLwfYCXAoKCkehlOWjL1ssKUtSwXQfVhpQJg
0BnYsEH9bnMP4l4HLlAOadDalHY7wnkAFP+exi50CiPPYr3Fg05F6Nci0XfnJ6o=
-----END CERTIFICATE-----`)

var SERVER_KEY = []byte(`-----BEGIN PRIVATE KEY-----
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
pqYDEuWdZ1VBCIaG5sZtuDpfrO+6eH8=
-----END PRIVATE KEY-----`)

var CA_CERT = []byte(`-----BEGIN CERTIFICATE-----
MIIF9DCCA9ygAwIBAgIJAODqYUwoVjJkMA0GCSqGSIb3DQEBCwUAMIGOMQswCQYD
VQQGEwJJTDEPMA0GA1UECAwGQ2VudGVyMQwwCgYDVQQHDANMb2QxEDAOBgNVBAoM
B0dvUHJveHkxEDAOBgNVBAsMB0dvUHJveHkxGjAYBgNVBAMMEWdvcHJveHkuZ2l0
aHViLmlvMSAwHgYJKoZIhvcNAQkBFhFlbGF6YXJsQGdtYWlsLmNvbTAeFw0xNzA0
MDUyMDAwMTBaFw0zNzAzMzEyMDAwMTBaMIGOMQswCQYDVQQGEwJJTDEPMA0GA1UE
CAwGQ2VudGVyMQwwCgYDVQQHDANMb2QxEDAOBgNVBAoMB0dvUHJveHkxEDAOBgNV
BAsMB0dvUHJveHkxGjAYBgNVBAMMEWdvcHJveHkuZ2l0aHViLmlvMSAwHgYJKoZI
hvcNAQkBFhFlbGF6YXJsQGdtYWlsLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIP
ADCCAgoCggIBAJ4Qy+H6hhoY1s0QRcvIhxrjSHaO/RbaFj3rwqcnpOgFq07gRdI9
3c0TFKQJHpgv6feLRhEvX/YllFYu4J35lM9ZcYY4qlKFuStcX8Jm8fqpgtmAMBzP
sqtqDi8M9RQGKENzU9IFOnCV7SAeh45scMuI3wz8wrjBcH7zquHkvqUSYZz035t9
V6WTrHyTEvT4w+lFOVN2bA/6DAIxrjBiF6DhoJqnha0SZtDfv77XpwGG3EhA/qoh
hiYrDruYK7zJdESQL44LwzMPupVigqalfv+YHfQjbhT951IVurW2NJgRyBE62dLr
lHYdtT9tCTCrd+KJNMJ+jp9hAjdIu1Br/kifU4F4+4ZLMR9Ueji0GkkPKsYdyMnq
j0p0PogyvP1l4qmboPImMYtaoFuYmMYlebgC9LN10bL91K4+jLt0I1YntEzrqgJo
WsJztYDw543NzSy5W+/cq4XRYgtq1b0RWwuUiswezmMoeyHZ8BQJe2xMjAOllASD
fqa8OK3WABHJpy4zUrnUBiMuPITzD/FuDx4C5IwwlC68gHAZblNqpBZCX0nFCtKj
YOcI2So5HbQ2OC8QF+zGVuduHUSok4hSy2BBfZ1pfvziqBeetWJwFvapGB44nIHh
WKNKvqOxLNIy7e+TGRiWOomrAWM18VSR9LZbBxpJK7PLSzWqYJYTRCZHAgMBAAGj
UzBRMB0GA1UdDgQWBBR4uDD9Y6x7iUoHO+32ioOcw1ICZTAfBgNVHSMEGDAWgBR4
uDD9Y6x7iUoHO+32ioOcw1ICZTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEB
CwUAA4ICAQAaCEupzGGqcdh+L7BzhX7zyd7yzAKUoLxFrxaZY34Xyj3lcx1XoK6F
AqsH2JM25GixgadzhNt92JP7vzoWeHZtLfstrPS638Y1zZi6toy4E49viYjFk5J0
C6ZcFC04VYWWx6z0HwJuAS08tZ37JuFXpJGfXJOjZCQyxse0Lg0tuKLMeXDCk2Y3
Ba0noeuNyHRoWXXPyiUoeApkVCU5gIsyiJSWOjhJ5hpJG06rQNfNYexgKrrraEin
o0jmEMtJMx5TtD83hSnLCnFGBBq5lkE7jgXME1KsbIE3lJZzRX1mQwUK8CJDYxye
i6M/dzSvy0SsPvz8fTAlprXRtWWtJQmxgWENp3Dv+0Pmux/l+ilk7KA4sMXGhsfr
bvTOeWl1/uoFTPYiWR/ww7QEPLq23yDFY04Q7Un0qjIk8ExvaY8lCkXMgc8i7sGY
VfvOYb0zm67EfAQl3TW8Ky5fl5CcxpVCD360Bzi6hwjYixa3qEeBggOixFQBFWft
8wrkKTHpOQXjn4sDPtet8imm9UYEtzWrFX6T9MFYkBR0/yye0FIh9+YPiTA6WB86
NCNwK5Yl6HuvF97CIH5CdgO+5C7KifUtqTOL8pQKbNwy0S3sNYvB+njGvRpR7pKV
BUnFpB/Atptqr4CUlTXrc5IPLAqAfmwk5IKcwy3EXUbruf9Dwz69YA==
-----END CERTIFICATE-----`)

var CA_KEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)

type ServerConfigRaw struct {
	CertificatePath    string   `json:"certificatePath"`
	CertificateKeyPath string   `json:"certificateKeyPath"`
	LayersPath         string   `json:"layersPath"`
	UrlPatterns        []string `json:"urlPatterns"`
}

func LoadServerConfigFromFile(path string) (*ServerConfig, error) {
	config := new(ServerConfigRaw)

	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configData, config)
	if err != nil {
		return nil, err
	}

	return LoadServerConfig(config)
}

type ServerConfig struct {
	Hosts       []string
	UrlPatterns []*regexp.Regexp
	CACert      *tls.Certificate
	LayersPath  string
}

func LoadServerConfig(rawConfig *ServerConfigRaw) (*ServerConfig, error) {
	certData, err := ioutil.ReadFile(rawConfig.CertificatePath)
	if err != nil {
		return nil, err
	}

	keyData, err := ioutil.ReadFile(rawConfig.CertificateKeyPath)
	if err != nil {
		return nil, err
	}

	caCert, _ := tls.X509KeyPair(certData, keyData)
	re := regexp.MustCompile("production.cloudflare.docker.com:443/registry-v2/docker/registry/v2/blobs/sha256/.+/(.+)/.*")
	return &ServerConfig{
		Hosts: []string{
			"production.cloudflare.docker.com:443",
		},
		UrlPatterns: []*regexp.Regexp{
			re,
		},
		CACert:     &caCert,
		LayersPath: rawConfig.LayersPath,
	}, nil
}
