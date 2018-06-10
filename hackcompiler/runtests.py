import subprocess as sp

# Specify where to save coverage output.
coverage_output = "tests/coverage.txt"
coverage_html = "tests/coverage.html"

# Run the test commands.
sp.run(["go", "test", "-v", "-cover", "-coverprofile="+coverage_output], shell=True)
sp.run(["go", "tool", "cover", "-html="+coverage_output, "-o", coverage_html], shell=True)

