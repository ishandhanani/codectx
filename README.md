# CodeCTX  

CodeCTX is one of many utility tools that let you take a codebase and convert it into a txt file which can be fed into an LLM for context. Here are the current features 

Install via 
```
git clone https://github.com/ishandhanani/codectx
cd codectx
go build -o codectx 
go install
```

Run using
```
./codectx --path="/path/to/your/codebase" --output="your_output_filename"
```

Flags
- `--path`: The directory path that contains your codebase.
- `--filetype`: Specify file extensions to include, separated by commas (e.g., .py,.js,.html). Leave empty to include all files.
- `--output`: The base name for your output file (without extension), e.g., combined_code.
- `--verbose`: Enable verbose mode to display detailed logs during execution.

Example
Combining python and javascript files from a specific directory into an output file named `project_context.txt`
```
./codectx --path="/path/to/project" --filetype=".py,.js" --output="project_context" --verbose
```
