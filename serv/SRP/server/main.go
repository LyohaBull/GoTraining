package main

// Simple test program to test the SRP library
// Author: Sudhi Herle
// April 2014

import (
	"crypto"
	"crypto/sha256"
	"fmt"

	"github.com/opencoff/go-srp"
)

func lookup() (string, string) {
	c := sha256.New()
	fmt.Println(c.Size())
	s, _ := srp.NewWithHash(crypto.SHA256, 1024)
	//s, _ := srp.New(1024)
	v, _ := s.Verifier([]byte("alex"), []byte("dsfsf4545454t"))
	id, verif := v.Encode()
	return id, verif
}
func main() {
	//ih := "b8fd39a36525c3a6aa40660f6093b2ae5024d23dcd1d9da9421ad53989fb07f8"

	//vh := "128:eeaf0ab9adb38dd69c33f80afa8fc5e86072618775ff3c0b9ea2314c9c256576d674df7496ea81d3383b4813d692c6e0e0d5d8e250b98be48e495c1d6089dad15dc7d7b46154d6b6ce8ef4ad69b15d4982559b297bcf1885c529f566660e57ec68edbc3c05726cc02fd4cbf4976eaa9afd5138fe8376435b9fc61d2fc0eb06e3:2:17:b8fd39a36525c3a6aa40660f6093b2ae5024d23dcd1d9da9421ad53989fb07f8:6b3678dc2f743205e0c16971e00d440cc32041d1af47f801b5b6db9746fc7e758ec5c0872b721d0b33f14d8c16c853c0e34102d2735444893210775d1a0a50d153f6e06653de79a075dfb6d9acb3e36bf86bb174b0484f5dd1af5dd7332a970554f369b10e0d9c00a957dd673cd2cf269210a7eeac623137f6c1405f63ac372e:9b5eb0d02f858161ae50b3cf10b930a593a5f0a0c58dc6e0f869393c20d5f8475e6ceab7c5e2a5fae2b13a0f8c67c85063770fc038e4becbcb9dba7a7c4a392afdb731fb5096ec712cbe4ac7dfe64d823f8d5c659bd42980fde484f0f80f6b50b675c496f75d316ebe053c09dd086647431baba2864a448a9b5d129e0cc5a31e"
	// Store ih, vh in durable storage
	_, verifier := lookup()
	fmt.Println(verifier)
	//creds := "b8fd39a36525c3a6aa40660f6093b2ae5024d23dcd1d9da9421ad53989fb07f8:9753c09dd65ba66786df164b63a26a5d3c3c1fa088445e58fb71bd56c000d579f62053ed5b9c5cf2c5570ad80ca5c7051bec826a8b284f0139a5763588d320011283b2e77691e3065b9f79dfaaa71d37ba4289cab72c4b6c23f9a06ddaca86343434f33e01388af7083fbb068a71511d7fb46fa3c1c3083d7370215f5dca39ab"
	creds := ""
	fmt.Scanf("%s\n", &creds)
	// Begin the server by parsing the client public key and identity.
	_, A, err := srp.ServerBegin(creds)
	if err != nil {
		panic(err)
	}

	// Now, pretend to lookup the user db using "I" as the key and
	// fetch salt, verifier etc.
	s, v, err := srp.MakeSRPVerifier(verifier)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server Begin; <v, A>:\n   %s\n   %x\n", verifier, A.Bytes())
	srv, err := s.NewServer(v, A)
	if err != nil {
		panic(err)
	}

	// Generate the credentials to send to client
	creds = srv.Credentials()

	// Send the server public key and salt to server
	fmt.Printf("Server Begin; <s, B> --> client:\n   %s\n", creds)

	// client processes the server creds and generates
	// a mutual authenticator; the authenticator is sent
	// to the server as proof that the client derived its keys.
	//cauth := "155e16e3746035fe1e0c67c6f8eb1fb26a273864d36d790888e7f52b9b91a59c"

	cauth := ""
	fmt.Scanf("%s\n", &cauth)
	// Receive the proof of authentication from client
	proof, ok := srv.ClientOk(cauth)

	if !ok {
		panic("client auth failed")
	}

	// Send proof to the client
	fmt.Printf("Server Authenticator: M' --> Server\n   %s\n", proof)

	// Verify the server's proof
	// Now, we have successfully authenticated the client to the
	// server and vice versa.

	ks := srv.RawKey()

	fmt.Printf("Server Key: %x\n", ks)

}

// EOF
