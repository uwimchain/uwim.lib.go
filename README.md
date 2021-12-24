# uwim.lib.go

<h2>Installation:</h2> 

```
go get github.com/uwimchain/uwim.lib.go
```

```go
import "github.com/uwimchain/uwim.lib.go"
```
<h2>Generate mnemonic phrase</h2>

```go
mnemonic, err := uwim_lib_go.GenerateMnemonic()
```

To generate public , secret keys or addresses from a mnemonic phrase, you can use a ready-made mnemonic phrase.

<h2>Generation of a Seed string from a mnemonic phrase</h2>
  
```go
seed, err := uwim_lib_go.SeedFromMnemonic(mnemonic)
```
<h2>Generating a secret key from a Seed string or mnemonic phrase</h2>

```go
secret_key, err := uwim_lib_go.SecretKeyFromSeed(seed);<br><br>
secret_key, err := uwim_lib_go.SecretKeyFromMnemonic(mnemonic)
```
<h2>Generation of a public key from a secret key or mnemonic phrase</h2>

```go
public_key, err := uwim_lib_go.PublicKeyFromSecretKey(secret_key)
public_key, err := uwim_lib_go.PublicKeyFromMnemonic(mnemonic)
```
<h2>Generating a user's address from a public key or mnemonic phrase</h2>

To generate an address, you can use a public key or a mnemonic phrase. You must also specify one of the three available prefixes. If you specify any other prefix, the function will return an error<br><br>

<h3>Generation of an address with the "uw" prefix - the address of the user's wallet</h3>

```go
uw_address, err := uwim_lib_go.AddressFromPublicKey(public_key, "uw")
uw_address, err := uwim_lib_go.AddressFromMnemonic(mnemonic, "uw")
```
<h3>Generation of an address with the "sc" prefix - the address of the smart contract</h3>

```go
sc_address, err := uwim_lib_go.AddressFromPublicKey(public_key, "sc")
sc_address, err := uwim_lib_go.AddressFromMnemonic(mnemonic, "sc")
```
<h3>Generation of an address with the "nd" prefix – the address of the node.</h3>

```go
nd_address, err := uwim_lib_go.AddressFromPublicKey(public_key, "nd")
nd_address, err := uwim_lib_go.AddressFromMnemonic(mnemonic, "nd")
```
<h2>Receiving a RAW transaction line for sending to the blockchain API </h2>

In order to generate a RAW transaction line, you need to specify such data as: 

<ul>
 Mnemonic phrase (the sender of the transaction);<br>
  Sender's address (must be generated from a mnemonic phrase or match it);<br>
  Recipient's address;<br>
  The number of coins you want to transfer (for some transaction types or transaction subtypes, the number of coins may be zero);<br>
  The designation of the token whose coins you want to transfer (for example: "uwm");<br>
  Transaction subtype (for example: "default_transaction");<br>
  Data of the comment to the transaction in JSON format (for each type or subtype of the transaction, its own comment data is indicated or not indicated at all);<br>
  Transaction type (Number 1 or 3);
</ul>
  
```go
transaction_raw, err := uwim_lib_go.GetRawTransaction(
    mnemonic,
    sender_address,
    recipient_addres,
    amount,
    token_label,
    transaction_comment_title,
    transaction_comment_data,
    transaction_type
)
```
