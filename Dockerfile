FROM golang:1.13

COPY SudokuSolver.go /usr/bin/
ENTRYPOINT ["go", "run", "/usr/bin/SudokuSolver.go"]
CMD [""]


