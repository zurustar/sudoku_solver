#
# 起動時にsudoku_solver.jsonがおいてあるディレクトリを
# コンテナ側の/sudoku_solver/config にマウントしてください
#
FROM ubuntu:18.04
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y git golang
RUN git clone https://github.com/zurustar/sudoku_solver.git
RUN go get github.com/ChimeraCoder/anaconda
RUN cd /sudoku_solver && go build ./sudoku_solver.go
RUN cd /sudoku_solver && go build ./twitter_bot.go
CMD /sudoku_solver/sudoku_solver
