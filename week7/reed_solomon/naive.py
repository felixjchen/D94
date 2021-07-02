from lib import *
import numpy as np
# https://www.cs.cmu.edu/afs/cs/project/pscico-guyb/realworld/www/slidesF15/ecc4.pdf

verbose = True 
def dp(*args):
  args = (colored(0,100,250, i) for i in args)
  if verbose:
    print(*args)


n = 5
k = 3
d = n-k+1

t = n-k 

print(f"A [{n}, {k}, {d}] reed solomon code")
print(f"Detects {t} errors")
print(f"Detects and corrects {t//2} errors")


# Evaluate the polynomial p(x) at [1, n]
# p(x) is a polynomial of degree k-1, with coefficients from message
def encode(message):

  def p(x):
    res = 0
    for i, coeff in enumerate(message[::-1]):
      res += coeff * (x**i)
    return res 

  xs = range(1, n+1)
  ys = [p(i) for i in xs]
  dp(ys)
  return ys

# Interpolate
def decode(codeword):
  ''' 
  # Cases:
  # 1. No errors
  # 2. 0 < # errors <= t//2 
  # 3. t//2 < # errors <= t 
  # '''

  codeword = np.array(codeword)

  # Cases:
  # 1. No errors
  # interpolate a degree k-1 polynomial (with k points), verify all points

  # Get first k-1 deg poly 
  ys = codeword[:k]
  xs = range(1,k+1)
  V = np.vander(np.array(xs), k)
  dp(V, ys)
  coefficients = np.linalg.solve(V, ys)
  coefficients = np.rint(coefficients)
  dp(coefficients)

  # Check coefficients
  xs = range(1,n+1)
  V = np.vander(np.array(xs), k)
  if np.allclose(V @ coefficients, codeword):
    print("No errors!")
    return list(coefficients)

  # 2. 0 < # errors <= t//2 

  # 3. t//2 < # errors <= t 
  pass

if __name__ == "__main__":
  message = [1,0,0]
  assert len(message) == k 

  codeword = encode(message)
  assert len(codeword) == n 

  codeword[1] = -44
  red_print("With error", codeword)

  ecc_message = decode(codeword)

  print(ecc_message)