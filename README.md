# keyvault
Practice Project: key vault to store your secrets, provide RESTFul API, build with Gin

> Just a practice project, it works in internal network. It's not a production-ready application.

# Features

- client and server use certificate to validate each other
- server perform as a CA, and able to issue certificate to client
- only authorized clients can communicate with keyvault
- the certificate must have OU property, server will use it as the namespace of secret to query database
- server has a master key to encrypt/decrypt the secrets using AES
- the master key is randomly generated and stored in a file that only visible to the user who start the service
- the master key is also encrypted in the file and only visible during runtime
- server has RESTFul API to store the secrets, get secrets, issue certificates
- communication is protected by TLS

The secret still visible as a plain text to authorized client. If the client choose to print it out or log it, there is nothing we can do. The sensitive data should never been seen anywhere.

## Workflow

### Storing a new secret
- client initiate the request to store a new secret, a secret is a string
- server received the string and encode it with base64 and then encrypted by server master key, store them in database
- if the secret already exist under same namespace, new record won't be created

### Getting an existing secret
- parse the namespace from certificate
- query secret name under the specific namespace and return 404 Not Found if no record found
- decrypt the encrypted text with master key of the record
- decode the decrypted text and return to client

## Usage

Creating a new secret (human)

```
POST /keyvault
{
    "name": "K8S_ADMIN_PASSWORD",
    "value": "the!realpassw0rd"
}

Response
{
    "message": "success/failed"
}
```

Getting an existing secret (in the program, certificates deployed)
```
GET /keyvault/?s=K8S_ADMIN_PASSWORD

Response
{
    "name": "K8S_ADMIN_PASSWORD",
    "value": "the!realpassw0rd"
}
```

Issue a CSR (human)
```
POST /keyvault/ca
{
    "csr": "<valid_x509_content>"
}

Response
{
    "signed": "<valid_x509_content>"
}
```