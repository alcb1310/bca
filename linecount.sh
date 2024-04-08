#! /bin/bash
linecount=`bash -c "find . -type f \( -name '*.go' -and -not -name '*_templ.go' -or -name '*.templ' \) | xargs wc -l | tail -1"`
lines=$(echo $linecount | awk '{print $1}')

filecount=`bash -c "find . -type f \( -name '*.go' -and -not -name '*_templ.go' -or -name '*.templ' \) | wc -l"`

generated=`bash -c "find . -name '*_templ.go' | xargs wc -l | tail -1"`
generatedlines=$(echo $generated | awk '{print $1}')

generatedfile=`bash -c "find . -name '*_templ.go' | wc -l"`

echo "Total lines of code: "$lines
echo "Total lines of generated code: "$generatedlines
echo "Total files: "$filecount
echo "Total generated files: "$generatedfile
echo "Average lines per file: "$((lines/filecount))
echo "Average lines per generated file: "$((generatedlines/generatedfile))
