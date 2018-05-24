#!/usr/bin/env python

import sys
import copy

#peers = []


#------------------------------------------------------------------------------
def init_peers():
  peers = []
  for sy in range(9):
    for sx in range(9):
      s = sy * 9 + sx
      peers.append([])
      for dy in range(9):
        for dx in range(9):
          d = dy * 9 + dx
          if s != d:
            if sx == dx or sy == dy:
              peers[s].append(d)
            elif (sy//3)*3+sx//3 == (dy//3)*3+dx//3:
              peers[s].append(d)
  return peers

#------------------------------------------------------------------------------
def check(board, peers):
  for i in range(81):
    if len(board[i]) == 1:
      for peer in peers[i]:
        if len(board[peer]) == 1:
          if board[i][0] == board[peer][0]:
            raise ValueError
#------------------------------------------------------------------------------
def _update1(board, peers):
  updated = False
  for target_pos in range(81):
    if len(board[target_pos]) == 1:
      target_value = board[target_pos][0]
      for peer_pos in peers[target_pos]:
        if target_value in board[peer_pos]:
          board[peer_pos].remove(target_value)
          if len(board[peer_pos]) == 0:
            raise ValueError
          updated = True
  return board, updated
#------------------------------------------------------------------------------
def _update2(board, peers):
  updated = False
  for target_pos in range(81):
    if len(board[target_pos]) != 1:
      for target_cand_value in board[target_pos]:
        found = False
        for peer_pos in peers[target_pos]:
          for peer_cand_value in board[peer_pos]:
            found = True
            break
        if found == False:
          board[target_pos] = [target_cand_value]
          updated = True
          break
  return board, updated
#------------------------------------------------------------------------------
def update(board, peers):
  updated = True 
  while updated:
    board, updated = _update1(board, peers)
    if not updated:
      board, updated = _update2(board, peers)
  return board

#------------------------------------------------------------------------------
def solve(board, peers, depth):
  #print("depth =", depth)
  result = True
  check(board, peers)
  board = update(board, peers)
  #show(board, depth)
  counter = [[],[],[],[],[],[],[],[],[],[]]
  for i in range(81):
    l = len(board[i])
    if l == 0:
      raise ValueError
    counter[l].append(i)
  if len(counter[1]) == 81:
    show(board)
    sys.exit()
  for i in [2,3,4,5,6,7,8,9]:
    for p in counter[i]:
      for c in board[p]:
        nb = copy.deepcopy(board)
        nb[p] = [c]
        try:
          solve(nb, peers, depth+1)
        except ValueError:
          board[p].remove(c)
          pass
        
#------------------------------------------------------------------------------
def show_detail(board):
  s = ""
  for y in range(9):
    s += "|"
    for x in range(9):
      for v in [1,2,3,4,5,6,7,8,9]:
        if v in board[y*9+x]:
          s += str(v)
        else:
          s += " "
      s += "|"
    s += "\n"
    if y in [2,5]:
      s += "-" * 91 + "\n"
  print(s)
#------------------------------------------------------------------------------
def show(board, depth = 0):
  s = ""
  for y in range(9):
    s += " " * depth
    for x in range(9):
      pos = y * 9 + x
      if len(board[pos]) == 0:
        raise ValueError
      elif len(board[pos]) == 1:
        s += str(board[pos][0])
      else:
        s += "0"
    s += "\n"
  print(s)

#------------------------------------------------------------------------------
def main(path):
  board = []
  with open(path) as fp:
    buf = fp.read()
    for c in buf:
      if c == "0":
        board.append([1,2,3,4,5,6,7,8,9])
      elif c in "123456789":
        board.append([int(c)])
  if len(board) == 81:
    solve(board, init_peers(), 0)
        
#------------------------------------------------------------------------------
if __name__ == '__main__':
  main(sys.argv[1])

