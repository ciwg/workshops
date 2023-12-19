# System Architecture

This system architecture centers on autonomous 'grok agents', each paired with a different 'brain' (AI, Human, IoT), interacting with a shared version control system (VCS) such as Git. Each brain is coupled with an agent facilitating a smooth exchange of instructions and feedback, thereby contributing to an organized flow of operations.

## agent.dot

The `agent.dot` diagram illustrates the high-level interaction between the filesystem and the grok agents. The agents, each coupled with a unique type of 'brain', consistently fetch and push updates to and from the versioning system, creating a cyclical flow of information exchange.

## agent-flow.dot

The `agent-flow.dot` diagram provides a deeper insight into the operation of individual agents during their collaboration with the VCS and their respective brains. The depicted workflow is as follows:

1. Brain uses `clone` to duplicate the repository and initialize the workspace.
2. Brain switches the current working directory (using `cd`) to the name of the branch cloned.
3. Brain issues the `join` command to create and join a specific branch in the VCS.
4. The agent fetches updates from the VCS regularly.
5. If there are no merge conflicts, the agent merges new commits from other branches and produces a list of updated files on stdout, for the Brain to act upon.
6. If merge conflicts occur, the Brain assists the agent in resolving these before the changes are committed.
7. The agent pushes final commits to the VCS.

In this operation, the Brain has two primary responsibilities:

- Acting on the list of updated files that the agent aggregates after the fetch and merge operations.
- Assisting to resolve conflicts when the agent encounters merge conflicts during commit operations.

# Initialization of an Agent's Workspace

To set up an agent's workspace, the Brain first clones the repository using the `clone` command. The Brain then changes its current directory to match the branch name, just before issuing the `join` command. The `join` command creates a new branch in the VCS and the agent begins operating within the context of this branch. This process simplifies the initialization process for the agent's workspace.

```bash
grok agent clone <clone-source-url>
cd <branch-name>
grok agent join <branch-name>
```

# Running an Autonomous Agent

An AI Brain can operate a fully autonomous grok agent using the `grok agent run` command. The command requires following arguments: 
'-r' for the role name, 
'-i' for input files, and 
'-o' for output files. 

The command accepts multiple comma-separated file names for `-i` and `-o` options. The role name implies the branch name, thus simplifying the operation and eliminating redundancy.

```bash
grok agent run -r <role-name> -i <input-file1,input-file2,...> -o <output-file1,output-file2,...>
```

# Joining a Group and Monitoring the VCS with grok Agents

A Brain can allow a grok agent to join a group, creating its own branch and begin tracking file changes within a Git repository. It uses the `grok agent join` command to accomplish this. The branch name is inferred from the role specified in the `grok agent run` command.

Here are the steps for the `join` command:

1. By issuing the `join` command, the Brain signifies the grok agent to initiate its own branch in the VCS and join the group of grok agents working on the repository.
2. Each agent within the group operates autonomously on its dedicated branch, regularly fetching updates from the VCS.
3. The agent watches for file changes in the repository, taking a similar approach to `inotifywait`. The notifications about these changes come from the git log, not directly from the filesystem.
4. Upon detecting modifications within certain directories or files, the Agent pauses its operation and notifies the Brain.
5. If file changes occur, the Agent fetches any new commits from different branches and shares an updated file list with the Brain. The Brain then processes these updates.

This iterative interaction between the Brain and the agent forms a continuous development loop, with file change notifications setting the course for consequent actions.

