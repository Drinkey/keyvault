## Client - draft

The command-line format

`keyvault [object] [action] [arguments]`

## Usage

Before running the client, make sure the service is running and configuration is correct.

`keyvault certificate init --client [namespace]`
- Function
    - initialize the client, create the certificates for further communication
    - private key and csr are self generated automatically, the csr will be signed by service
- Details
  - this command will create a directory `~/.keyvault_certs/` and put all certificates inside it.
    - `keyvault_client.pkey`
    - `[namespace].csr`
    - `[namespace].crt`
  - certificate parameters are specified by `config.json` (default name)


`keyvault ns create --name [namespace]`
- create a new namespace

`keyvault secret create --ns [namespace] --key [key_name]`
- create a new secret with key [key_name] under namespace [namespace]

`keyvault secret get --ns [namespace] --key [key_name]`
- get a existing secret with key [key_name] under namespace [namespace]
