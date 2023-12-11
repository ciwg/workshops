# System Architecture

This system architecture centers on autonomous 'grok agents', each paired with a different 'brain' (AI, Human, IoT), interacting with a shared version control system (VCS) such as Git. Each brain is coupled with an agent facilitating a smooth exchange of instructions and feedback, thereby contributing to an organized flow of operations.

## agent.dot

The `agent.dot` diagram illustrates the high-level interaction between the filing system and the grok agents. The agents, each coupled with a unique type of 'brain', consistently fetch and push updates to and from the versioning system, creating a cyclical flow of information exchange.

## agent-flow.dot

The `agent-flow.dot` diagram provides a deeper insight into the operation of individual agents during their collaboration with the VCS and their respective brains. The depicted workflow is as follows:

1. Brain uses `clone` to duplicate the repository and initialize the workspace.
2. Brain switches the current working directory (using `cd`) to the name of the branch cloned.
3. Brain issues the `join` command to join a specific branch. Upon executing this operation, the agent pushes any pending updates.
4. The agent fetches updates from the VCS regularly.
5. If there are no merge conflicts, the agent merges new commits from other branches and produces a list of updated files on stdout, for the Brain to act upon.
6. If merge conflicts occur, the Brain assists the agent in resolving these before the changes are committed.
7. The agent pushes final commits to the VCS.

In this operation, the Brain has two primary responsibilities:

- Acting on the list of updated files that the agent aggregates after the fetch and merge operations.
- Assisting to resolve conflicts when the agent encounters merge conflicts during commit operations.

# Initialization of an Agent's Workspace

To set up an agent's workspace, the Brain first clones the repository using the `clone` command. The Brain then changes its current directory to match the branch name, just before issuing the `join` command. This process simplifies the initialization process for the agent's workspace.

```bash
grok agent clone <clone-source-url>
cd <branch-name>
grok agent join <branch-name>
```

# Running an Autonomous Agent

An AI Brain can operate a fully autonomous grok agent using the `grok agent run` command. This command needs the following arguments: 
'-r' for the role name, 
'-i' for input files, and 
'-o' for output files. 

The input and output files arguments can handle multiple comma-separated file names. Using the role name to infer the branch name eliminates redundancy and simplifies the operation.

```bash
grok agent run -r <role-name> -i <input-file1,input-file2,...> -o <output-file1,output-file2,...>
```

# Joining a Group of grok Agents and Monitoring the VCS

A Brain can join a group of grok agents and start monitoring file changes within a Git repository using the `grok agent join` command. Here, the branch name is inferred from the role defined during the `run` command.

Executing the `join` command follows these steps:

1. The Brain integrates with a group of grok agents, each of which operates on its dedicated branch.
2. The Agent pushes any local changes that haven't been synced yet, if any.
3. The Agent fetches updates from the VCS regularly.
4. Appearing similar to `inotifywait`, the Agent watches for file changes - it actually gets notifications about these changes from the git log, not directly from the filesystem.
5. The Agent halts execution and waits for modifications within certain directories or files.
6. Upon spotting file changes, the Agent fetches any new commits from different branches and lists the updated files on stdout for the Brain to address.

This ongoing interaction forms a continuous development loop, driven by the changes pushed to the VCS, setting the course for the next cycle.
