# Testing Only

## How To Generate

`mkcert featureguards.dev "*.featureguards.dev" 127.0.0.1 ::1`

`ssh-keygen -t rsa -b 4096 -m PEM -f jwt.key`
`openssl rsa -in jwt.key -pubout -outform PEM -out jwt.key.pub`
