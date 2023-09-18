package types

var PrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAPQykV9lViQ1A20+
KVf+PE+BpSDJXi3Uve4p0dlSWEae+sDsYoCRi6ObXcKSaH7f8mM1oEiIbD5c7SFF
Xt2A7SJ/Ir4mYPYUbKutvPL4duyktA/uEvat7GZbSb3NwaO8wZ9kT4kAi00BxdM1
DKubM0E/DZjtSEA2oPz2xY5TrwzTAgMBAAECgYAPOJuxG4rsBNXq2EYRcwplVkpp
qcOSDcGs97RZ3HUeKcitf85//xJ6JzQH7cJPrjvYjT4pZz9//6DUQxOvsNqW/bO3
vMIqrKPXqBsIqCzxPaLJfiSyqgH6NNCb5qhNU2+/DTNmFhdZ8IVnm0H7gRXxEBbG
HMD4bqJSR3XGE0PO0QJBAPz+cNT8gZbs0PXkYUArSNaW2yt/Rzih3MW4x5+eqODG
3+gISbJG85A0Z65hnQuQiO2EPZSeVWFRNGm9JeozAOkCQQD3GV7FxgCg7/MJZZZ4
Nf5DDSeqsDYmmizpLWZG8MidnrhoQkdL67wHYumCruUB+jIUwi8x3eucLHGkx+N9
e6pbAkEAzoMDp1fWkFQO3ijmGXM7qa7KiN8ES/4UMHF8wZbJU3IDI2xge93ewz+D
wpx7jQ0WOItRmRcFqsKWfhmf8WRgwQJAbZYB0wJqOvXPumYkYnHHruMBqZB2o44S
xuMMjf+xaT4AGLT0O7ZzcG8sknmQNN1KIqywE5SRLnUDfYns2TTkKwJBANqWceIe
zLJ2tOMVj854btO7xhvZgvj9FOgLw+who+LnFXJArNNpu2CzaU13IhZodmiWNGli
VslR62tWjF8IGnk=
-----END PRIVATE KEY-----
`)

// 公钥: 根据私钥生成
// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var PublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD0MpFfZVYkNQNtPilX/jxPgaUg
yV4t1L3uKdHZUlhGnvrA7GKAkYujm13Ckmh+3/JjNaBIiGw+XO0hRV7dgO0ifyK+
JmD2FGyrrbzy+HbspLQP7hL2rexmW0m9zcGjvMGfZE+JAItNAcXTNQyrmzNBPw2Y
7UhANqD89sWOU68M0wIDAQAB
-----END PUBLIC KEY-----
`)

var AesKey = "cure-d111y=1ziukr07k*!r$q=zcgto%" //AES密钥
var Base_binance_url = "https://api.binance.com"

var ApiSecret = "ZWfi7cEzDnG6jUvk3aK9tW53WZweiHLttC4jN2ZQqo5uIklPSSLmfzuDCflUJJFM"
var ApiKey = "CaVXeFXYYkpygoacCuQihH5QIpTP0Y00J3UdC3eEj3wrnq0s9KDBa7V6msNBI16l"
