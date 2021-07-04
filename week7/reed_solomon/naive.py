from lib import *
import numpy as np
import itertools

# https://www.cs.cmu.edu/afs/cs/project/pscico-guyb/realworld/www/slidesF15/ecc4.pdf

verbose = True 
def dp(*args):
  args = (colored(0,100,250, i) for i in args)
  if verbose:
    print(*args)


n = 7
k = 3
d = n-k+1

t = n-k 

evaluation_points = list(range(n))

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

  return [p(i) for i in evaluation_points]

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
  xs = evaluation_points[:k]
  V = np.vander(np.array(xs), k)
  # dp(V, ys)
  coefficients = np.linalg.solve(V, ys)
  coefficients = np.rint(coefficients)
  # dp(coefficients)

  # Check coefficients works for all points
  V = np.vander(np.array(evaluation_points), k)
  if np.allclose(V @ coefficients, codeword):
    print("No errors!")
    return list(coefficients)

  # 2. 0 < # errors <= t//2 
  for errors in itertools.product(range(0,n), repeat=t//2):
    dp("guessing errors:", set(errors))

    remaining_points = set(evaluation_points) - set(errors)
    dp(f"Checking points {remaining_points}")

  # 3. t//2 < # errors <= t 
  pass

if __name__ == "__main__":
  message = [1,0,0]
  print(f"Message  {message}")
  assert len(message) == k 

  codeword = encode(message)
  print(f"Codeword {codeword}")
  assert len(codeword) == n 

  codeword[1] = -44
  codeword[4] = -44
  red_print("With error", codeword)

  ecc_message = decode(codeword)
  print(f"Recovered message  {ecc_message}")