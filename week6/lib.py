def colored(r, g, b, text):
    return "\033[38;2;{};{};{}m{}\033[38;2;255;255;255m".format(r, g, b, text)


def pad_left(toPad, char, size):
    while len(toPad) < size:
        toPad = char + toPad
    return toPad


def flip_bit(codeword, flip):
    new_value = "1" if codeword[flip] == "0" else "0"
    codeword = list(codeword)
    codeword[flip] = new_value
    codeword = "".join(codeword)
    return codeword


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
