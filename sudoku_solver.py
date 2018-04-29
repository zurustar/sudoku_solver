#!/usr/bin/env python

import sys

dic = []


# -----------------------------------------------------------------------------
def _update1(q):
  updated = False
  for s in range(9*9):
    if len(q[s]) != 1:
      continue
    sx, sy, sg = dic[s][0], dic[s][1], dic[s][2]
    for d in range(9*9):
      if s != d:
        dx, dy, dg = dic[d][0], dic[d][1], dic[d][2]
        if sx == dx:
          if q[s][0] in q[d]:
            q[d].remove(q[s][0])
            updated = True
            print(sx, sy, "が", q[s][0], "なので" , dx, dy, "から消す")
        elif sy == dy:
          if q[s][0] in q[d]:
            q[d].remove(q[s][0])
            updated = True
            print(sx, sy, "が", q[s][0], "なので" , dx, dy, "から消す")
        elif sg == dg:
          if q[s][0] in q[d]:
            q[d].remove(q[s][0])
            updated = True
            print(sx, sy, "が", q[s][0], "なので" , dx, dy, "から消す")
  return updated

# -----------------------------------------------------------------------------
def _update2(q):
  updated = False
  for s in range(9*9):
    if len(q[s]) == 1:
      continue
    sx, sy, sg = dic[s][0], dic[s][1], dic[s][2]
    for v in q[s]:
      found = False
      for d in range(9*9):
        if s != d:
          dx, dy, dg = dic[d][0], dic[d][1], dic[d][2]
          if sx == dx:
            if v in q[d]:
              found = True
      if not found:
        print(sx, sy, "の", v, "が見当たらないので確定")
        q[s] = [v]
        updated = True
    for v in q[s]:
      found = False
      for d in range(9*9):
        if s != d:
          dx, dy, dg = dic[d][0], dic[d][1], dic[d][2]
          if sy == dy:
            if v in q[d]:
              found = True
      if not found:
        print(sx, sy, "の", v, "が見当たらないので確定")
        q[s] = [v]
        updated = True
    for v in q[s]:
      found = False
      for d in range(9*9):
        if s != d:
          dx, dy, dg = dic[d][0], dic[d][1], dic[d][2]
          if sg == dg:
            if v in q[d]:
              found = True
      if not found:
        print(sx, sy, "の", v, "が見当たらないので確定")
        q[s] = [v]
        updated = True
  return updated

# -----------------------------------------------------------------------------
def update(q):
  print("." * 20, "update", "." * 20)
  show(q)
  while True:
    result = _update1(q)
    if result == False:
      result = _update2(q)
      if result == False:
        break
  show(q)
  return q

# -----------------------------------------------------------------------------
def solve(q):
  print("- " * 10, "solve", " -" * 10)
  update(q)

# -----------------------------------------------------------------------------
def show(q):
  print("#"*9)
  for y in range(9):
    s = ""
    for x in range(9):
      if len(q[y*9+x]) == 1:
        s += str(q[y*9+x][0])
      else:
        s += "-"
    print(s)
  print(q)

# -----------------------------------------------------------------------------
def init_dic():
  g = [0,0,0,1,1,1,2,2,2]
  for y in range(9):
    for x in range(9):
      dic.append([x, y, g[y]*3+g[x]])

# -----------------------------------------------------------------------------
def main():
  init_dic()
  for filename in sys.argv[1:]:
    fp = open(filename, "r")
    buf = fp.read()
    fp.close()
    q = []
    for c in buf:
      if "0" <= c and c <= "9":
        value = int(c)
        if value == 0:
          q.append([1,2,3,4,5,6,7,8,9])
        else:
          q.append([value])
    if len(q) == 9 * 9:
      solve(update(q))

# -----------------------------------------------------------------------------
if __name__ == '__main__':
  main()
