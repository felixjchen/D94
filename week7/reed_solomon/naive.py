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
print(f"Detects and corrects {t//2} errors")

def interpolate(xs, ys, degree):
  '''Returns the coefficient for a polynomial with degree'''

  assert len(xs) >= degree + 1 and len(ys) >= degree + 1, f"Cannot interpolate {xs} {ys}, missing data"
  # We only need n+1 points for a degree poly
  xs = xs[:degree+1]
  ys = ys[:degree+1]

  V = np.vander(xs, degree + 1)

  x = np.linalg.solve(V,  ys)
  x = np.rint(x)
  return x

def verify(coefficients, xs, ys):
  ''' Verifies a polynomial with coefficients maps xs to ys'''
  V = np.vander(xs, k)
  return np.allclose(V @ coefficients, ys)


# Evaluate the polynomial p(x) at [1, n]
# p(x) is a polynomial of degree k-1, with coefficients from message
def encode(message):
  V  = np.vander(evaluation_points, k)
  coefficients = np.array(message)
  return V @ coefficients

# Interpolate
def decode(codeword):
  ''' 
  # Cases:
  # 1. No errors
  # 2. 0 < # errors <= t//2 
  # 3. t//2 < # errors <= t 
  # '''

  # Cases:
  # 1. No errors
  # interpolate a degree k-1 polynomial (with k points), verify all points

  # Get first k-1 deg poly 
  coefficients = interpolate(evaluation_points, codeword, k-1)
  if verify(coefficients, evaluation_points, codeword):
    print("No errors!")
    return list(coefficients)

  # 2. 0 < # errors <= t//2 

  # O(n^(t//2)) possible errors
  l = list(itertools.product(range(n), repeat=t//2))
  # Remove duplicates
  l = [tuple(set(i)) for i in l]
  l = list(set(l))
  l = sorted(l, key= lambda x: len(x))

  for errors in l:
    dp("guessing errors:", set(errors))

    xs = np.array(list(set(evaluation_points) - set(errors)))
    ys = codeword[xs]
    coefficients = interpolate(xs, ys, k-1)

    if verify(coefficients, xs, ys):
      print(f"Found errors {errors}")
      return coefficients

  return "Cannot recover"

if __name__ == "__main__":
  message = [1,0,3]
  print(f"Message  {message}")
  assert len(message) == k 

  codeword = encode(message)
  print(f"Codeword {codeword}")
  assert len(codeword) == n 

  codeword[1] = -44
  codeword[3] = -44
  red_print("With error", codeword)

  ecc_message = decode(codeword)
  print(f"Recovered message  {ecc_message}")