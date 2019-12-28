# sudoku_solver


goを実行する環境をおもちで無い場合は、カレントディレクトリに問題がかかれたファイル（ここでは1.txtと仮定）を置いておき、以下のような漢字で実行すれば解こうとしてくれる。

docker container run -v \`pwd\`:/tmp/work $(docker image build -q .) /tmp/work/1.txt

