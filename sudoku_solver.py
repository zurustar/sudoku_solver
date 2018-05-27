#!/usr/bin/env python

import sys
import copy

#--------------------------------------------------------------------
#
# 同じ値があってはいけない場所一覧を作成する。
#
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

#--------------------------------------------------------------------
#
# 同じ値があってはいけない場所に同じ値があったら例外をなげる
#
def check(board, peers):
  for i in range(81):
    # 確定している場所をみつけたら、
    if len(board[i]) == 1:
      # 同じあったらいけない場所をもってきて、
      for peer in peers[i]:
        # そこも値が確定しているかを調べ、
        if len(board[peer]) == 1:
          # 確定している値が同じなら
          if board[i][0] == board[peer][0]:
            # おかしいので例外。
            raise ValueError

#--------------------------------------------------------------------
#
# 値が確定している場所があったら、その場所からみて同じ値があっては
# いけない場所からその確定している値を削除する
#
def _update1(board, peers):
  updated = False # この関数内での処理で値が更新されたか？
  for target_pos in range(81):
    # 候補がひとつ＝確定済み だったら
    if len(board[target_pos]) == 1:
      # 確定した値を抜き出して、
      target_value = board[target_pos][0]
      # 同じ値があってはいけない所
      for peer_pos in peers[target_pos]:
        # その値があったら、、
        if target_value in board[peer_pos]:
          # 消す。
          board[peer_pos].remove(target_value)
          # 値がゼロ個になったら
          if len(board[peer_pos]) == 0:
            # 仮おきした値が間違っているので例外
            raise ValueError
          updated = True
  return board, updated

#--------------------------------------------------------------------
#
# ある場所の候補について、その場所からみて同じ値があってはいけない
# 場所の候補に含まれていなかったら、その候補の値に確定する
#
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

#--------------------------------------------------------------------
#
# _update1と_update2を更新できなくなるまで繰り返し実行する
#
def update(board, peers):
  updated = True 
  while updated:
    board, updated = _update1(board, peers)
    if not updated:
      board, updated = _update2(board, peers)
  return board

#-------------------------------------------------------------------
#
# 現在の状態に矛盾が生じていないかを確認し、
# 更新してみて、解けているかを確認し、解けていなかったら
# 仮おきして再帰的にこの関数を呼び出す。
#
def solve(board, peers, depth):
  # 矛盾が生じていないかを確認
  check(board, peers)
  # ルールに則り解いてみる
  board = update(board, peers)
  show_detail(board, depth)
  # 残りの候補数を調べる
  counter = [[],[],[],[],[],[],[],[],[],[]]
  for i in range(81):
    l = len(board[i])
    if l == 0:
      raise ValueError
    counter[l].append(i)
  # すべて残りの候補がひとつになっていたら解けたという意味。
  if len(counter[1]) == 81:
    show(board)
    sys.exit()
  # 残り候補数が少ないものから順に借り置きして試しに解いてみる
  for i in [2,3,4,5,6,7,8,9]:
    for p in counter[i]:
      for c in board[p]:
        nb = copy.deepcopy(board)
        nb[p] = [c]
        try:
          solve(nb, peers, depth+1)
        except ValueError:
          board[p].remove(c)
        
#--------------------------------------------------------------------
def show_detail(board, depth):
  s = ""
  for y in range(9):
    s += " " * depth
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
      s += " " * depth + "-" * 91 + "\n"
  print(s)
#--------------------------------------------------------------------
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

#--------------------------------------------------------------------
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
        
#--------------------------------------------------------------------
if __name__ == '__main__':
  main(sys.argv[1])

