Here's the bash script based on your requirements:

```bash
#!/bin/bash

# Ensure 'grok msg' command is installed
if ! command -v grok msg &> /dev/null
then
    echo "grok msg command doesn't exist. Please, install it and rerun the script."
    exit
fi

# For task related to the integration of grok msg AI tool 
# Assume 'grok msg' is the command to run your AI tool, replace it with actual command if it's different
# Calling the grok msg for code evolution tool 
grok msg "evolve code"

# For task related to installing 'grok' command without the user consent step
# Assume 'install-grok' is the command to install 'grok', replace it with actual command if it's different
# Installing 'grok' command without user consent
yes | install-grok

# For debugging
# Since no specific bugs or code are provided, this section only contains pseudo code 
# Execution of the debug process
# ...

# Testing the grok msg AI tool integration and 'grok' command installation
# Testing result will be saved to 'test_result.txt'
# Assume 'test' is the command to run your test, replace it with actual command if it's different
test > test_result.txt 

# As per user stories, presenting the test results
echo "Test results:"
cat test_result.txt

echo "Script executed successfully"
```

This script assumes the grok msg AI and 'grok' command installation can be achieved by just calling them in the terminal, this might not be true in a real situation. Also, the script doesn't handle any errors coming from these commands, so you would need to further optimize this for good exception handling and further automation.
