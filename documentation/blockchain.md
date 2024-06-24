## Bitcoin Blockchain

### Database structure

Bitcoin core uses 'two buckets' to store data.

1. **`blocks`** stores metadata describing all the blocks in a chain.
2. **`chainstate`** stores state of chain. It is all currently unspent transaction outputs and some metadata. `chainstate` does not store transaction, it stores a set of unspent transaction output.

### Transaction

#### Transaction lifecycle

1. In the beginning, there is the genesis block that contains a coinbase transaction. There are no real inputs in coinbase transactions, so signining is not required. The output of the coinbase transaction contains a hashed public key(RIPEMD16(SHA256(PubKey)) algorithms are used.)

2. When one sends coins, a transaction is created. Inputs of the transaction will reference outputs from previous transactions. Every input will store a public key (not hased) and a signature of the whole transaction.

3. Other nodes in the Bitcoin network that receive the transaction will verify it. Besides other things, they will check that: the hash of the public key in an input matches the hash of the reference output(this ensures that the sender spends only coins belonging to them); the signature is correct (this ensures that the transaction is created by the real owner of the coins).

4. When a miner node is ready to mine a new block, it will put the transaction in a block and start mining it.

5. When the block is mined, every other node in the network receives a message saying the block is mined and adds the block to the blockchain.

6. After a block is added to the blockchain, the transaction is completed, its outputs can be referenced in new transactions.

#### Inputs and Outputs

1. Inputs

   - **Reference Previous Outputs:** Each input in a transaction refers to an unspent output from a previous transaction.
   - **Spend the Value:** The value of each input is derived from the previous output it references.

2. Outputs
   - **Define New Outputs:** Each output specifies a value and a destination address (or script in Bitcoin’s case).

#### Detailed Example:

1. Previous Transaction (TX1):

   - Outputs:
     - Output 0: 0.5 BTC to address_A
     - Output 1: 1.0 BTC to address_B

2. Current Transaction (TX2):

   - Inputs:

     - Input 0: References TX1: Output 0 (0.5 BTC)
     - Input 1: References TX1: Output 1 (1.0 BTC)

   - Outputs:
     - Output 0: 1.3 BTC to address_C
     - Output 1: 0.1 BTC to address_D

Calculation:
Total Input Value:

Input 0: 0.5 BTC (from TX1: Output 0)
Input 1: 1.0 BTC (from TX1: Output 1)
Total: 0.5 BTC + 1.0 BTC = 1.5 BTC
Total Output Value:

Output 0: 1.3 BTC to address_C
Output 1: 0.1 BTC to address_D
Total: 1.3 BTC + 0.1 BTC = 1.4 BTC
Transaction Fee:

Total Input Value (1.5 BTC) - Total Output Value (1.4 BTC) = 0.1 BTC

Validation Rules:
Sum of Inputs: The total value of inputs must equal the total value of outputs plus the transaction fee.
Unspent Outputs: The inputs must reference outputs that have not been spent yet, ensuring the value is available.

#### coinbase transaction

When a miner starts mining a block, it adds a coinbase transaction to it. A coinbase transaction is a special type of transactions, which doesn’t require previously existing outputs. It creates outputs (i.e., “coins”) out of nowhere. This is the reward miners get for mining new blocks.

#### Notice

1. There are outputs that are not linked to inputs.
2. In one transaction, inputs can reference ouptuts from multiple transactions.
3. An input must reference an output.

### Wallet

In bitcoin, wallet is a pair of public and private key. Bitcoin relies on a combination of cryptography algorithms to create these keys, and guarantee that no one else in the world can access your coins without getting physical access to your keys.

A private key is generated using the curve, and a public key is generated from the private key. One thing to notice: in elliptic curve based algorithms, public keys are points on a elliptic curve. Thus, a public key is a combination of X, Y coordinates. In Bitcoin, these coordinates are concatenated and form a public key.

#### steps to convert a public key to base58 wallet address

Here are the steps to convert a public key into a Base58 address:

1. Take the public key and hash it twice with RIPEMD160(SHA256(PubKey)) hashing algorithms.
2. Prepend the version of the address generation algorithm to the hash.
3. Calculate the checksum by hashing the result of step 2 with SHA256(SHA256(payload)). The checksum is the first four bytes of the resulted hash.
4. Append the checksum to the version+PubKeyHash combination.
5. Encode the version+PubKeyHash+checksum combination with Base58.

### Merkle tree

The full Bitcoin database takes more than 140 Gb of disk space. Because of the decentralized nature of Bitcoin, every node in the network must be independent and self-sufficient, i.e. every node must store a full copy of the blockchain. With many people starting using Bitcoin, this rule becomes more difficult to follow: it’s not likely that everyone will run a full node. Also, since nodes are full-fledged participants of the network, they have responsibilities: they must verify transactions and blocks. Also, there’s certain internet traffic required to interact with other nodes and download new blocks.

There was a solution for this problem: Simplified Payment Verification (SPV). SPV is a light Bitcoin node that doesn’t download the whole blockchain and doesn’t verify blocks and transactions. Instead, it finds transactions in blocks (to verify payments) and is linked to a full node to retrieve just necessary data. This mechanism allows having multiple light wallet nodes with running just one full node.

For SPV to be possible, there should be a way to check if a block contains certain transaction without downloading the whole block. And this is where Merkle tree comes into play.

Merkle trees are used by Bitcoin to obtain transactions hash, which is then saved in block headers and is considered by the proof-of-work system. Until now, we just concatenated hashes of each transaction in a block and applied SHA-256 to them. This is also a good way of getting a unique representation of block transactions, but it doesn’t have benefits of Merkle trees.

### Network

Blockchain network is decentralized, which means there’re no servers that do stuff and clients that use servers to get or process data. In blockchain network there are nodes, and each node is a full-fledged member of the network. A node is everything: it’s both a client and a server.

Blockchain network is a P2P (Peer-to-Peer) network, which means that nodes are connected directly to each other. It’s topology is flat, since there are no hierarchy in node roles.

#### Node roles

Blockchain nodes can play different roles in the network. Here they are:

1. Miner
   Such nodes are run on powerful or specialized hardware (like ASIC), and their only goal is to mine new blocks as fast as possible. Miners are only possible in blockchains that use Proof-of-Work, because mining actually means solving PoW puzzles. In Proof-of-Stake blockchains, for example, there’s no mining.

2. Full node
   These nodes validate blocks mined by miners and verify transactions. To do this, they must have the whole copy of blockchain. Also, such nodes perform such routing operations, like helping other nodes to discover each other.
   It’s very crucial for network to have many full nodes, because it’s these nodes that make decisions: they decide if a block or transaction is valid.

3. SPV
   SPV stands for Simplified Payment Verification. These nodes don’t store a full copy of blockchain, but they still able to verify transactions (not all of them, but a subset, for example, those that were sent to specific address). An SPV node depends on a full node to get data from, and there could be many SPV nodes connected to one full node. SPV makes wallet applications possible: one don’t need to download full blockchain, but still can verify their transactions.
