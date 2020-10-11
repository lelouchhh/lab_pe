cd db/

go build -o bin

cd ..

touch delete.sh
echo "rm db/bin;rm delete.sh" > delete.sh