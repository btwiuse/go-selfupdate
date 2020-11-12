/*
Package update provides functionality to implement secure, self-updating Go programs (or other single-file targets).

Basic Example

This example shows how to update a program remotely from a URL.

	import (
		"fmt"
		"net/http"

		"github.com/creativeprojects/go-selfupdate/update"
	)

	func doUpdate(url string) error {
		// request the new file
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		err := update.Apply(resp.Body, update.Options{})
		if err != nil {
			if rerr := update.RollbackError(err); rerr != nil {
				fmt.Println("Failed to rollback from bad update: %v", rerr)
			}
		}
		return err
	}


Binary Patching

Go binaries can often be large. It can be advantageous to only ship a binary patch to a client
instead of the complete program text of a new version.

This example shows how to update a program with a bsdiff binary patch. Other patch formats
may be applied by implementing the Patcher interface.

	import (
		"encoding/hex"
		"io"

		"github.com/creativeprojects/go-selfupdate/update"
	)

	func updateWithPatch(patch io.Reader) error {
		err := update.Apply(patch, update.Options{
			Patcher: update.NewBSDiffPatcher()
		})
		if err != nil {
			// error handling
		}
		return err
	}

Checksum Verification

Updating executable code on a computer can be a dangerous operation unless you
take the appropriate steps to guarantee the authenticity of the new code. While
checksum verification is important, it should always be combined with signature
verification (next section) to guarantee that the code came from a trusted party.

go-update validates SHA256 checksums by default, but this is pluggable via the Hash
property on the Options struct.

This example shows how to guarantee that the newly-updated binary is verified to
have an appropriate checksum (that was otherwise retrived via a secure channel)
specified as a hex string.

	import (
		"crypto"
		_ "crypto/sha256"
		"encoding/hex"
		"io"

		"github.com/creativeprojects/go-selfupdate/update"
	)

	func updateWithChecksum(binary io.Reader, hexChecksum string) error {
		checksum, err := hex.DecodeString(hexChecksum)
		if err != nil {
			return err
		}
		err = update.Apply(binary, update.Options{
			Hash: crypto.SHA256, 	// this is the default, you don't need to specify it
			Checksum: checksum,
		})
		if err != nil {
			// error handling
		}
		return err
	}

Cryptographic Signature Verification

Cryptographic verification of new code from an update is an extremely important way to guarantee the
security and integrity of your updates.

Verification is performed by validating the signature of a hash of the new file. This
means nothing changes if you apply your update with a patch.

This example shows how to add signature verification to your updates. To make all of this work
an application distributor must first create a public/private key pair and embed the public key
into their application. When they issue a new release, the issuer must sign the new executable file
with the private key and distribute the signature along with the update.

	import (
		"crypto"
		_ "crypto/sha256"
		"encoding/hex"
		"io"

		"github.com/creativeprojects/go-selfupdate/update"
	)

	var publicKey = []byte(`
	-----BEGIN PUBLIC KEY-----
	MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEtrVmBxQvheRArXjg2vG1xIprWGuCyESx
	MMY8pjmjepSy2kuz+nl9aFLqmr+rDNdYvEBqQaZrYMc6k29gjvoQnQ==
	-----END PUBLIC KEY-----
	`)

	func verifiedUpdate(binary io.Reader, hexChecksum, hexSignature string) {
		checksum, err := hex.DecodeString(hexChecksum)
		if err != nil {
			return err
		}
		signature, err := hex.DecodeString(hexSignature)
		if err != nil {
			return err
		}
		opts := update.Options{
			Checksum: checksum,
			Signature: signature,
			Hash: crypto.SHA256, 	                 // this is the default, you don't need to specify it
			Verifier: update.NewECDSAVerifier(),   // this is the default, you don't need to specify it
		}
		err = opts.SetPublicKeyPEM(publicKey)
		if err != nil {
			return err
		}
		err = update.Apply(binary, opts)
		if err != nil {
			// error handling
		}
		return err
	}


*/
package update
