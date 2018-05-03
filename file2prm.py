#!/usr/bin/env python

import sys

def main(filename):
	fp = open(filename, "r")
	buf = fp.read()
	fp.close()
	s = ""
	for c in buf:
		if "0" <= c and c <= "9":
			s += c
	print(s)


if __name__ == '__main__':
	main(sys.argv[1])
