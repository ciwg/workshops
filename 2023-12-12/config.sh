
agents="facilitator pygame-developer tester user"
# associate outputs with agents
declare -A agent_outputs
agent_outputs[user]="requirements.md"
agent_outputs[facilitator]="tasks.md"
agent_outputs[pygame-developer]="pong.py"
agent_outputs[tester]="test-results.md"

