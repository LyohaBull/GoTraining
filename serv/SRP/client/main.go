package main

// Simple test program to test the SRP library
// Author: Sudhi Herle
// April 2014

import (
	"fmt"

	"github.com/opencoff/go-srp"
)

func main() {
	bits := 1024
	pass := []byte("alice")
	i := []byte("password123")

	s, err := srp.New(bits)
	if err != nil {
		panic(err)
	}

	c, err := s.NewClient(i, pass)
	if err != nil {
		panic(err)
	}

	// client credentials (public key and identity) to send to server
	creds := c.Credentials()

	fmt.Printf("Client Begin; <I, A> --> server:\n   %s\n", creds)

	// Now, pretend to lookup the user db using "I" as the key and
	// fetch salt, verifier etc.
	//serv_creds := "fcb16f98b189d94f3d0aa80aa2341396dec54c503746714a3e1dcc784242050dbad6d59b5e4e62181f20504c10becfddfbc40a78de4d68b99139184d3d191c22320f22fe53f7d9ab76a8558b77e7a848da1c3ea236255e051ceca0d17523ec783a7da7355f90f99907828e8ae5d64752e043ee87e21a7afbf34b1cd7537aa4dc:0bb7c836a04a6254e504d635df5b93af7585685d12a06a68755488a3b8470cbf6fae8c5274ae55d169414105e0c4a89c8b5615b2e6877b0745b139f4a44daf423fd2d704daf61c65094515441bcdd2d1ee7cf639e50f039bc382acfd36539d211d302595c9f310e2d45e1754743afb36b60887dbc0ff9c8304d7ae7479a16377"
	serv_creds := ""
	fmt.Scanf("%s\n", &serv_creds)
	cauth, err := c.Generate(serv_creds)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Client Authenticator: M --> Server\n   %s\n", cauth)

	// Receive the proof of authentication from client

	//proof := "5a6c919de7c891e0c7dc95566da8bb6eafe27147c94248870f5d541756e2381b"
	proof := ""
	fmt.Scanf("%s\n", &proof)
	// Verify the server's proof

	if !c.ServerOk(proof) {
		panic("server auth failed")
	}

	// Now, we have successfully authenticated the client to the
	// server and vice versa.

	kc := c.RawKey()
	fmt.Printf("Client Key: %x\n", kc)

}

// EOF
