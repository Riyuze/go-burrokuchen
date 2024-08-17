# Go-Burokkuchen

This is a very simple implementation of a blokchain in the go programming language. This project follows a guide from articles created by [Jeiwan](https://github.com/Jeiwan).

---

# Chapter 1: [Basic Prototype](https://jeiwan.net/posts/building-blockchain-in-go-part-1/)

## Block

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

In Bitcoin specifications `Timestamp`, `PrevBlockHash`, and `Hash` are [block headers](https://developer.bitcoin.org/reference/block_chain.html#), which form a separate data structure, and transactions (`Data` in our case) is also a separate data structure.

## Blockchain

In its essence a blockchain is a database with an ordered, back-linked list structure. This means that every blocks are stored in the insertion order and each block is linked to the previous one. The implementation of this structure is as follows:

```golang
type Blockchain struct {
	blocks []*Block
}
```

The blockchain will be used as an array of blocks, with each block having a connection to the previous one. The actual blockchain is much more complex though.

---
