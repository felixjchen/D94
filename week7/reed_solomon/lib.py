def colored(r, g, b, text):
    return "\033[38;2;{};{};{}m{}\033[38;2;255;255;255m".format(r, g, b, text)

def red_print(*text):
  print(*(colored(240,0,0, i) for i in text))

def pad_left(toPad, char, size):
    while len(toPad) < size:
        toPad = char + toPad
    return toPad


def print_block(binary_block, parity_indexes):
    out = "[ "
    for i, b in enumerate(binary_block):
        if i in parity_indexes:
            out += colored(0, 0, 150, b)
        else:
            out += colored(0, 200, 0, b)

        out += ' '
    out += "]"
    print(out)
