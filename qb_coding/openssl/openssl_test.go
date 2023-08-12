package openssl

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	plaintext := "Hello World!"
	passphrase := "z4yH36a6zerhfE5427ZV"
	opensslEncrypted := "U2FsdGVkX19ZM5qQJGe/d5A/4pccgH+arBGTp+QnWPU="

	o := NewOpenSSL()

	enc, err := o.EncryptBytes(passphrase, []byte(plaintext), BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}
	opensslEncrypted = string(enc)

	dec, err := o.DecryptBytes(passphrase, []byte(opensslEncrypted), BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}
	fmt.Printf("Decrypted text: %s\n", string(dec))

}
