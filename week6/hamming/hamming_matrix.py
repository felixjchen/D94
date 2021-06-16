import numpy as np
from lib import *

r = 4

assert r >= 2

n = 2**r - 1
k = 2**r - 1 - r
rate = k / n

parity_indexes = [2**i-1 for i in range(r)]
data_indexes = [i for i in range(0, n) if i not in parity_indexes]

block = [i for i in range(1, n+1)]
binary_block = [f"{i:b}" for i in block]
binary_block = [pad_left(i, "0", r) for i in binary_block]

print(f"[{n}, {k}, 3] code, rate = {rate}")
print_block(binary_block, parity_indexes)

table = np.zeros([r, n])

# Construct table for parity bit responsibility, each row is one parity bit
for i in range(r):
    for j in range(n):
        if binary_block[j][r-i-1] == "1":
            table[i][j] = 1

# Pop out parity bit indexed columns.. the rest should be a matrix for calucating parity bits
for col in reversed(parity_indexes):
    table = np.delete(table, col, axis=1)

# Identity matrix copies over data bits, we still need to calculate parity bits
G = np.concatenate((np.eye(k), table.T), axis=1)
H = np.concatenate((table, np.eye(n-k)), axis=1)
print(f"Generator Matrix \n{G}")
print(f"Parity Check Matrix \n{H}")


def encode(message: str):
    # Correct message len
    assert len(message) == k
    # Binary message
    assert all([i in ["0", "1"] for i in message])

    message = np.array([int(i) for i in message]).reshape([1, k])
    codeword = (message @ G) % 2

    # convert np array to string
    return "".join(str(i) for i in codeword.astype(int).flatten())


def decode(codeword: str):
    assert len(codeword) == n
    codeword_a = np.array([int(i) for i in codeword]).reshape([1, n])
    parity_check = ((codeword_a @ H.T) % 2).flatten()

    if all(i == 0 for i in parity_check):
        return codeword[:k]

    for i, col in enumerate(H.T):
        if np.allclose(col, parity_check):
            print(colored(255, 0, 0, f"correcting index {i}"))
            print(colored(255, 0, 0, codeword))
            new_value = "1" if codeword[i] == "0" else "0"
            codeword = codeword[:i] + new_value + codeword[i+1:]
            print(colored(0, 200, 0, codeword))
            return codeword[:k]


if __name__ == "__main__":
    message = "10110101101"
    print(f"original message {message}")

    codeword = encode(message)

    print(f"original codeword {codeword}")

    codeword = flip_bit(codeword, 5)

    print(f"noisy codeword {codeword}")

    ecc_message = decode(codeword)

    print(f"corrected message {ecc_message}")

    assert message == ecc_message
