## How to create encrypted LVol

Simplyblock logical volumes supports encryption at rest by leveraging [crypro bdev](https://spdk.io/doc/bdev.html) module using SPDK Software Accel framework.

The framework internally uses AES_XTS cipher which is used to sequence arbitrary length of block data. This type of cipher is standardly used in backup software.

AES_XTS cipher requires 2 keys of length either 16 bytes of 32 bytes. The hex value of the key can be generated using the command.

for example to generate a 32 byte key
```
openssl rand -hex 32
```
encode the generated 32 byte key
```
echo -n '7b3695268e2a6611a25ac4b1ee15f27f9bf6ea9783dada66a4a730ebf0492bfd' | base64

echo -n '78505636c8133d9be42e347f82785b81a879cd8133046f8fc0b36f17b078ad0c' | base64
```
After the keys are generated, an encrypted pvc can created by passing the generated keys `crypto_key1` and `crypto_key2` secret data
```
apiVersion: v1
Kind: Secret
metadata:
  name: simplyblock-pvc-keys
  namespace: default
data:
  crypto_key1: N2IzNjk1MjY4ZTJhNjYxMWEyNWFjNGIxZWUxNWYyN2Y5YmY2ZWE5NzgzZGFkYTY2YTRhNzMwZWJmMDQ5MmJmZA==
  crypto_key2: Nzg1MDU2MzZjODEzM2Q5YmU0MmUzNDdmODI3ODViODFhODc5Y2Q4MTMzMDQ2ZjhmYzBiMzZmMTdiMDc4YWQwYw==
```
Then set the encryption parameter of the storageclass to `True`
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: spdkcsi-sc
provisioner: csi.simplyblock.io
parameters:
  ...
  encryption: "True"
```
Finally set pvc annotation which will dynamically retrieve the secret data `crypto_key1` and `crypto_key2` that is required for encrypted pvc
```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: spdkcsi-pvc
  annotations:
    simplybk/secret-name: simplyblock-pvc-keys
    simplybk/secret-namespace: default
spec:
  ...
  storageClassName: spdkcsi-sc
```
