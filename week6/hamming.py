import math
from lib import *

r = 5

assert r >= 2

n = 2**r - 1
k = 2**r - 1 - r
rate = k / n

print(f"[{n}, {k}, 3] code, rate = {rate}")

block = [i for i in range(1, n+1)]

# print(f"block: {block}")

binary_block = [f"{i:b}" for i in block]
binary_block = [pad_left(i, "0", r) for i in binary_block]

# print(f"binary_block: {binary_block}")

parity_indexes = [2**i-1 for i in range(r)]
data_indexes = [i for i in range(0, n) if i not in parity_indexes]

print_block(binary_block, parity_indexes)


def encode(message):
    # Correct message len
    assert len(message) == k
    # Binary message
    assert all([i in ["0", "1"] for i in message])

    # Empty codeword
    codeword = [None for _ in range(n)]

    # 1. Place message into codeword's data indicies
    message = list(message)
    i = 0
    while message:
        codeword[data_indexes[i]] = message.pop(0)
        i += 1

    # 2. Construct parity bits
    for i in range(r):
        # Loop over entire block
        count = 0
        for j in range(n):
            # If this is a bit I'm responsible for
            if binary_block[j][r-i-1] == "1":
                # Sum for parity
                if codeword[j] == "1":
                    count += 1

        codeword[parity_indexes[i]] = "1" if count % 2 == 1 else "0"

    # list -> str
    return "".join(codeword)


def decode(codeword):
    assert len(codeword) == n

    failed_parity_checks = 0
    error_pos = 0

    # Do parity checks
    for i in range(r):
        # Loop over entire block
        count = 0
        for j in range(n):
            # If this is a bit I'm responsible for
            if binary_block[j][r-i-1] == "1":
                # Sum for parity
                if codeword[j] == "1":
                    count += 1

        if count % 2 != 0:
            failed_parity_checks += 1
            error_pos += parity_indexes[i] + 1

    if failed_parity_checks > 0:
        new_value = "1" if codeword[error_pos - 1] == "0" else "0"
        print(colored(255, 0, 0, f"correcting index {error_pos - 1}"))
        print(colored(255, 0, 0, codeword))
        codeword = codeword[:error_pos - 1] + new_value + codeword[error_pos:]
        print(colored(0, 200, 0, codeword))

    # Extract message
    return "".join(bit for i, bit in enumerate(codeword) if i in data_indexes)


if __name__ == "__main__":
    message = "10111011101101110111010101"
    print(f"original message {message}")

    codeword = encode(message)

    print(f"original codeword {codeword}")

    codeword = flip_bit(codeword, 3)

    print(f"noisy codeword {codeword}")

    ecc_message = decode(codeword)

    print(f"corrected message {ecc_message}")

    assert message == ecc_message
