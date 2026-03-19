Cryptography:
    Generate Keys:
        use: Curve25519
    Key Exchange:
        use: ECDH
        Generate a shared secret key between two parties

        Example:
            Alice:
                private: a
                public: A
            Bob: 
                private: b
                public: B
            
            # Public key exchange
            Alice receives B
            Bob receives A

            # Calculating general secret key and getting the same result+
            Alice: S = ECDH(a, B)
            Bob: S = ECDH(b, A)

            *S - shared secret key
    Key Derivation Function (KDF):
        Convert the shared secret S into a fixed-length symmetric key K for encryption

        use: bcrypt (todo for research)
        parameters:
            salt: same for both parties (pre-agreed or transmitted securely)
            iterations: same for both parties
        
        process:
            K = KDF(S, salt, iterations)
        *K - result after KDF

        Bob computes K using the same S, salt, and iterations, so the result is identical
    Message Encryption:
        method: AES-GCM
        process:
            cipherMessage = AES-GCM(K, plainMessage, nonce)
            # nonce should be random and unique per message
    Message Sending:
        Alice sends:
            {
                cipherMessage: "...",
                publicKey: A,       # optional if Bob already knows it
                nonce: "..."
            }
    Message Receiving:
        Bob receives:
            cipherMessage, Alice's publicKey (A), nonce
            # Recompute shared secret
            S = ECDH(b, A)
            # Use same KDF parameters
            K = KDF(S, salt, iterations)
            # Decrypt
            plainMessage = AES-GCM-Decrypt(K, cipherMessage, nonce)