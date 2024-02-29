cd ../src || exit
black -l 120 .
isort .
autoflake -r --in-place --remove-all-unused-imports --remove-unused-variables .
cd - || exit
