/**
 * 登录用RSA加解密处理
 */

package base

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"

	"github.com/sirupsen/logrus"
)

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}

/**
 * 生成pkcs8格式公钥私钥
 */
func CreatePkcs8Keys(keyLength int) (privateKey, publicKey string) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: MarshalPKCS8PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))

	return
}

/**
 * 解密
 */
func Decrypt(secretData []byte, publicKey, privateKey string) ([]byte, error) {
	rsaPublicKey := rsaPublicKeyFromString(publicKey)
	rsaPrivateKey := rsaPrivateKeyFromString(privateKey)

	blockLength := rsaPublicKey.N.BitLen() / 8
	if len(secretData) <= blockLength {
		return rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, secretData)
	}

	buffer := bytes.NewBufferString("")

	pages := len(secretData) / blockLength
	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(secretData) {
				continue
			}
			end = len(secretData)
		}

		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, secretData[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 加密
 */
func Encrypt(data []byte, publicKey string) ([]byte, error) {
	rsaPublicKey := rsaPublicKeyFromString(publicKey)

	blockLength := rsaPublicKey.N.BitLen()/8 - 11
	if len(data) <= blockLength {
		return rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(data))
	}

	buffer := bytes.NewBufferString("")

	pages := len(data) / blockLength

	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}

		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, data[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

// 从RSA私钥字符串转rsa.PrivateKey对象
func rsaPrivateKeyFromString(secret string) *rsa.PrivateKey {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		panic(err)
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}

	return privateKey.(*rsa.PrivateKey)
}

// 从RSA公钥字符串转rsa.PublicKey
func rsaPublicKeyFromString(secret string) *rsa.PublicKey {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		panic(err)
	}

	return rsaPublicKey.(*rsa.PublicKey)
}

func Test() {
	// private, public := CreatePkcs8Keys(2048)
	// logrus.Info(private)
	// logrus.Info(public)

	private := `MIIEuwIBADALBgkqhkiG9w0BAQEEggSnMIIEowIBAAKCAQEAvrKWFpLztzPQ3UoK
Tt1QNAxfSoiqfqTefq0w92Jj9lYkVC6fQR2NBdQKsZeT6q8jAcEYWVnXHiGT3BRN
+WLEgMU52uSJSfdTulj9Nt3yYOtWzCmM2xGri14jhIa3IdnhfcX2XNoNB/fNgEeo
0o+t2+9yHSe9wW6RtPCBE8U46MGwTnFZtlBXikA66135ubyYtYd2LrPN7nPYEkNL
8TtYMbMTDmyzgUa5lOXypcbc+zd+QOgbIXtCkNvN17Nk9tFOB5yPN3toLNzOkmMh
YREv0dXndJ8klDynrPkkCfDxdIvq9RjPNoqTxhweCeXHDAekIn3lg54VC3CacQrP
luSupwIDAQABAoIBAHdNsRp0W2ctUqlvDd3jFa9KYj92GvxaVxx3a+AJPTK7F8VW
2alaPIT98KbEhvTXFxac4Ifd7fha129jgJjaEsfhG933BnEw+7/ktp4h4uaBtW7L
O+U+O81YWu4pfd7+udT/Ca9zd52ZiYaMznDVFNc5CXJ2D4A5lYzWvlpJE96BYj5V
u6SeSJCQtxv8Y4Ey7n/Er6HkjMwWCoWiaClpPzJeftJUibHxH063HLPvKGXZFGX5
9vF04xyfzuJQiLgbCTXrwwTaxYO+8UnS7esqlge4fkpzhYhFtTnyC3NSEgj91tJH
YeoeKFsbgPVdu3Evk4u0Gxdna0jEGnDqf4CZCcECgYEA3OpmD0edOioJrTqURy5G
qu1pKd+67QhmvtDITV7wKQfpnJEk/p4LEv2WfgEMkIouj14H3EJiuECexKlve0uq
w9wdofcpBCxCJkotmL/KSd6I6eALKa+XwK6YgtyEmrggixaicZfzpmNfaMAShOIZ
piSIQWQNgfIsBvoBZS+qDkcCgYEA3Pui4gWxVyNA7Sc0KBRm8IGAGQGXnYSzGotE
/+8JKk7YdYqaHiEY4+ZCja484nKWK+ZZ+Wg434TJTGzWaN/uNOY3H78yNttYRgnB
/QzkHVlp/IgiqyfOge86elH+7wCAh4wbJK+UDvaCRCifdiS9sTLKQsR/GBOyLwvw
KDTWrKECgYAw4XiFpv3mEckkWFLY0Sd3yKI9TrDIo9RAImg/nmMbYRHSv9bks8mV
gSDcbpT+ImUc+dxZYyL+y+WVdDwjluGJBtpTrSGZN8XHPSCLrNwwrhmzTgyKQ70b
OEaspeh9Z4Jj5DU7VzjlNxW0UtOGLZUpSuoPNfk7KH+PZ6AJaJuDHwKBgQCd3jsP
42c8zBefFInDNEgR+0HrG2MYCev1w5bIjBjtG6Sx3BGcAqMIdMAI/XfLgnbb59VR
Qu6WaANy0LIf/BHtwqWQzYNvAyY96symHeZ9PRplaU/zHB4AX0pUhm1sitxHeYUO
oUxRoDORw7+fpEHL7G/oYP420iNSTuIDpzPR4QKBgHUAONjv7DGVy5QbMKWdpkMk
9v/KcvtEhpfJSiSJYkK7Kbjnb5thAmK/b6LxU/jV8hf4qR9i5qsTLe4VE3LeHxmH
S7ra8PtGKRX4AZ/2i43+XZP3YqjSPPNR0LNLFY3HAmPYhFoWkLsLJl1Wzm6GteAk
3OgA0tYDb7foAq8RProS`

	public := `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvrKWFpLztzPQ3UoKTt1Q
NAxfSoiqfqTefq0w92Jj9lYkVC6fQR2NBdQKsZeT6q8jAcEYWVnXHiGT3BRN+WLE
gMU52uSJSfdTulj9Nt3yYOtWzCmM2xGri14jhIa3IdnhfcX2XNoNB/fNgEeo0o+t
2+9yHSe9wW6RtPCBE8U46MGwTnFZtlBXikA66135ubyYtYd2LrPN7nPYEkNL8TtY
MbMTDmyzgUa5lOXypcbc+zd+QOgbIXtCkNvN17Nk9tFOB5yPN3toLNzOkmMhYREv
0dXndJ8klDynrPkkCfDxdIvq9RjPNoqTxhweCeXHDAekIn3lg54VC3CacQrPluSu
pwIDAQAB`

	plaintext := "123456"

	// 123456
	// dmWmo2dbja/kXgtFoMjk2DgsSZ55LRI1S5mGuPjA7dCmati51BRNDV6jDNmosK73e1uWgDK8RHjB6CnukaKLBEFDjBt8nqNKKeWWrHqnbeXas/zhr3JaUWRCzqGGl4xhFK3YyJgK0SxUhnEgu2M7kULFK9S3g9hlK65PUTczfR8IlSl+gbfxvx7ZloH2pC5I3l+B+mKQYmF2uKCOgJECU18E0vkAjOhnF3bZX4ebjbA5esbQCkbVG9HBHxbLr5h/phsOzYGwzh4K2QxY+TKzR0rtoSX65f6KgvpmHxwXS7PrwAWPS9L8jsEjBIaV3GCQuHzRMgi8LG0WOI69QUvLXg==

	// 666666
	// UURHcWIdFWdZUbuOLkYZvmfwEw4clscb2dP3duQExqg9p1C1wQ37bPLxaJh1oo7S+hwIttptdSElAq2LwxQOLLaiqrQZ+F/c1ArtT+x/GPPEA3bzReiu4cfdLjDPvEfXdOGTK9gAt2vGZ6TQle8YYN4Ze4THnuFxwm4B+BMuzN60peRJI1QzVjx15hqdQLVLnTF0+DBDVNwQHVJgtMh/YAzVMf5q2r/Ru5LjlUe13+OuxwdPd5fRzMZIfp1OtB+4C//T76celyzH/n8QACwCCHOya1nxm8vcYiC3tmiO4hJBckY73VNtDLjgdQRkihFwKkn+/ZzAmB7AHDIdywWz8A==
	secrettext, err := Encrypt([]byte(plaintext), public)
	if err != nil {
		panic(err)
	}

	encodeString := base64.StdEncoding.EncodeToString(secrettext)

	logrus.Info("密文：", encodeString)

	descrpttext, err := Decrypt([]byte(secrettext), public, private)
	if err != nil {
		panic(err)
	}

	logrus.Info(string(descrpttext))
}
