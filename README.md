# Go-Burokkuchen

This is a very simple implementation of a blokchain in the go programming language. This project follows a guide from articles created by [Jeiwan](https://github.com/Jeiwan).

---

## Table of Contents

- [Go-Burokkuchen](#go-burokkuchen)
	- [Table of Contents](#table-of-contents)
	- [Chapter 1: Basic Prototype](#chapter-1-basic-prototype)
		- [Block](#block)
		- [Blockchain](#blockchain)
	- [Chapter 2: Proof of Work](#chapter-2-proof-of-work)
		- [Proof of Work](#proof-of-work)
		- [Hashing](#hashing)
		- [Hashcash](#hashcash)

---

## Chapter 1: [Basic Prototype](https://jeiwan.net/posts/building-blockchain-in-go-part-1/)

### Block

Implemented a simplified version of a block with a data structure as follows:

```golang
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}
```

The fields in the block structure are:

-   `Timestamp` is the current timestamp (when the block is created).
-   `Data` is the valuable information contained in the block.
-   `PrevBlockHash` stores the hash of the previous block.
-   `Hash` is the hash of the block

In Bitcoin specifications `Timestamp`, `PrevBlockHash`, and `Hash` are [**block headers**](https://developer.bitcoin.org/reference/block_chain.html#), which form a separate data structure, and transactions (`Data` in our case) is also a separate data structure.

### Blockchain

In its essence a blockchain is a database with an ordered, back-linked list structure. This means that every blocks are stored in the insertion order and each block is linked to the previous one. The implementation of this structure is as follows:

```golang
type Blockchain struct {
	blocks []*Block
}
```

The blockchain will be used as an array of blocks, with each block having a connection to the previous one. The actual blockchain is much more complex though.

---

## Chapter 2: [Proof of Work](https://jeiwan.net/posts/building-blockchain-in-go-part-2/)

### Proof of Work

A key idea of blockchain is that one has to perform some hard work to put data in it. It is this hard work that makes blockchain secure and consistent. Also, a reward is paid for this hard work (this is how people get coins for mining).

Proof of Work algorithms must meet a requirement: **doing the work is hard**, but **verifying the proof is easy**. A proof is usually handed to someone else, so for them, it shouldn’t take much time to verify it.

### Hashing

Hashing is a process of obtaining a hash for specified data. A hash is a unique representation of the data it was calculated on. A hash function is a function that takes data of arbitrary size and produces a fixed size hash. Here are some key features of hashing:

1. Original data cannot be restored from a hash. Thus, **hashing is not encryption**.
2. Certain data can **have only one hash** and the **hash is unique**.
3. Changing even one byte in the input data will result in a completely different hash.

In blockchain, hashing is used to **guarantee the consistency of a block**. The input data for a hashing algorithm contains the hash of the previous block, thus making it impossible (or, at least, quite difficult) to modify a block in the chain: one has to recalculate its hash and hashes of all the blocks after it.

### Hashcash

Bitcoin uses [_Hashcash_](https://en.wikipedia.org/wiki/Hashcash), a Proof of Work algorithm that was initially developed to prevent email spam. It can be split into the following steps:

1. Take some publicly known data (in case of email, it’s receiver’s email address; in case of Bitcoin, it’s block headers).
2. Add a counter to it. The counter starts at 0.
3. Get a hash of the data + counter combination.
4. Check that the hash meets certain requirements. 1. If it does, you’re done. 2. If it doesn’t, increase the counter and repeat the steps 3 and 4.

Thus, this is a brute force algorithm: you change the counter, calculate a new hash, check it, increment the counter, calculate a hash, check it again, and so on. That’s why it’s computationally expensive.

---
