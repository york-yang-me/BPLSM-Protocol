# BPLSM-Protocol
A Secure Multi-party Computation Protocol for Universal Data Privacy
Protection based on Blockchain
## Step 0 Preparation
 - Go SDK 1.13.7
 - Go mod
 - Truffle
 - Ganache v2.1.0

## Step 1 Off-chain test
### 1. BPLSM-bplsm_test.go
Run these test functions as follows:

```go
func BenchmarkBPLSMPedersenG(b *testing.B) {...}

func BenchmarkBPLSMKeyGen(b *testing.B) {...}
```
### 2. BLS-bls_test.go
Run test function as follows:

```go
func BenchmarkBLSAggregateSignatureN(b *testing.B) {...}
```
## Step 2 On-chain test
==You should use Truffle and Ganache to test on-chain verification code.==
Run these commands in your terminal:

```bash
truffle compile
truffle migrate
truffle test
```
Then you will see these right information in your terminal.
This is one test in my personal computers:
![terminal test](https://img-blog.csdnimg.cn/20200917233839716.png#pic_center)
>My computer configuration：
>Win 10
>CPU: 8GB intel i7
