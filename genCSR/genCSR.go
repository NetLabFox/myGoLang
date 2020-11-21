package gencsr

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	b64 "encoding/base64"
	"fmt"
)

/*func main() {
		reader := rand.Reader
		var subject pkix.Name
		C := []string{}
		C = append(C, "TW")
		OU := []string{"CHT"}
		subject.Country = C
		subject.OrganizationalUnit = OU
		subject.CommonName = "張香傑"
		keys, err := rsa.GenerateKey(reader, 2048)
		template := x509.CertificateRequest{
			SignatureAlgorithm: 4,
			Subject:            subject}
		if err != nil {
			fmt.Println("產生金鑰成功")
		}
		csr, err := x509.CreateCertificateRequest(reader, &template, keys)
		sEnc := b64.StdEncoding.EncodeToString([]byte(csr))
		fmt.Println(sEnc)
	GenCSR("TW", "CHT", []string{"資訊處"}, "張香傑")
}*/

//GenCSR 透過相關參數產生CSR
func GenCSR(Country string, Organization string, OrganizationalUnit []string, CommonName string) (csr string, err error) {
	reader := rand.Reader
	var subject pkix.Name
	/*C := []string{}
	O := []string{}
	OU := []string{}
	CN := CommonName*/
	/*C = append(C, Country)
	O = append(O, Organization)
	OU = append(OU, OrganizationalUnit...)*/
	subject.Country = append(subject.Country, Country)
	subject.OrganizationalUnit = append(subject.OrganizationalUnit, OrganizationalUnit...)
	subject.Organization = append(subject.Organization, Organization)
	subject.CommonName = CommonName

	keys, err := rsa.GenerateKey(reader, 2048)
	if err != nil {
		fmt.Println("產生金鑰成功")
	}
	template := x509.CertificateRequest{
		SignatureAlgorithm: 4,
		Subject:            subject}

	bCsr, err := x509.CreateCertificateRequest(reader, &template, keys)
	csr = b64.StdEncoding.EncodeToString([]byte(bCsr))
	fmt.Println(csr)
	return
}
