#!/bin/bash

# Define the input file
input_file="input.txt"

# Define the output file
output_file="factorize_output.txt"

# Remove the output file if it exists
if [ -f "$output_file" ]; then
  rm "$output_file"
fi

# Check if the input file exists
if [ ! -f "$input_file" ]; then
  echo "Input file $input_file not found."
  exit 1
fi

# Read each line from the input file and factorize the number
while IFS= read -r number; do
  echo "$number"
  ./factorize "$number"
done < "$input_file"