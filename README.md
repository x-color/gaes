# gaes

gaes is tool for AES-256-CBC Encription and Decription.

## Usage

Encription

```bash
$ gaes encrypt plain.txt encrypted.txt

$ gaes e -w file.txt
```

Decryption

```bash
$ gaes decrypt encrypted.txt plain.txt

$ gaes d -w file.txt
```

Peek encrypted file

```bash
$ gaes peek encrypted.txt

$ gaes p encrypted.txt
```

## Install 

Prerequisite Tools

- Git
- Go (at least Go 1.11)

```bash
$ git clone https://github.com/x-color/gaes.git
$ cd gaes
$ GO111MODULE=on go install
```