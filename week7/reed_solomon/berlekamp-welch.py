import numpy as np
import galois
from lib import *

n = 7
k = 3
d = n-k+1
print(f"A [{n}, {k}, {d}] reed solomon code")

t = n-k

evaluation_points = list(range(n))

print(f"Detects and corrects {t//2} errors")


def encode(message):
    # Evaluate the polynomial F at [0, n-1]
    # F is a polynomial of degree k-1, with coefficients from message
    V = np.vander(evaluation_points, k)
    coefficients = np.array(message)
    return V @ coefficients % n


def decode(codeword):
    # check errors for with size [1, t//2]
    for s in range(t//2, 0, -1):
        dp(f"checking {s} sized errors")

        # E has degree s
        vander_E = np.vander(evaluation_points, s)
        # Q has degree s+k-2
        vander_q = np.vander(evaluation_points, s+k)
        # => Q / E has degree k-1

        left = vander_E * \
            np.repeat(np.expand_dims(codeword, axis=1), s, axis=1)
        right = - vander_q

        A = np.concatenate((left, right), axis=1)
        b = - codeword * np.power(np.array(evaluation_points), s)

        # Why does this solve only work in GF(7)?
        GFn = galois.GF(n)
        A = GFn(A % n)
        b = GFn(b % n)
        dp(A, b)

        try:
            x = np.linalg.solve(A, b)
            dp(x)

            coefficients_E = np.insert(np.array(x[:s]), 0,  1)
            coefficients_Q = np.array(x[s:])
            message, remainder = np.polydiv(coefficients_Q, coefficients_E)
            message, remainder = message % n, remainder[1] % n
            if remainder == 0:
                return message

        except np.linalg.LinAlgError:
            pass


if __name__ == "__main__":
    message = [3, 2, 1]
    print(f"Message  {message}")

    codeword = encode(message)
    print(f"Codeword {codeword}")

    codeword[1] = 5
    # codeword[4] = 3
    red_print("With error", codeword)

    decoded_message = decode(codeword)

    print(f"Decoded message {decoded_message}")
