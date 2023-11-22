---
layout: true
class: center, middle, inverse
---

# Demo Meeting Notes

[.footer: RemarkJS Slide Presentation]

---
template: inverse
class: middle

# Outline

1. Pong Game Demo
2. Org Design
3. LLM Servers
4. Vector Embeddings
5. Multi-Agent Algorithm
6. Consensus Formation in Self-Organizing Groups
7. Postmortem 

---

# Pong Game Demo

- Next time, start with the showing of Pong Game
- Run `pong_8.py`
- Show `pong.sh`

---

# Org Design

- Discussion on importance of organization design.
- Importance of consensus formation in self-organizing groups in the realm of multi-agent systems.

---

# LLM Servers

- Talk about GPT-4 API as a chat continuation service.
- Both Grokker and ChatGPT function by maintaining local context.
- Local context maintained in session for ChatGPT or in the Grokker repository.
- Use of `grok add filename` to add a file to local context.

---

# Vector Embeddings

- Overview of embeddings - transform text to a point in multidimensional space.
- Both ChatGPT and Grokker maintain a vector database.
- Refer [A Beginner's Guide to Vector Embeddings](https://www.pinecone.io/learn/vector-database/) for more detail.

---

# Multi-Agent Algorithm

- Definition of an agent as a single sequence of chat messages fitting into a token limit.
- Use of larger local context with grokker's vector database for creating chat message sequence.
- Understanding of "there ain’t no such thing as a free lunch" principle.
- The cooperative, promise-based and self-organizing nature of a better algorithm.

---

# Consensus Formation in Self-Organizing Groups

- Review of Promise Theory by Mark Burgess
- Workshop next week: Consensus Formation in Self-Organizing Groups

---

# Postmortem

- Suggestion of a rolling slide deck to summarise and focus discussion.
- Use of markdown-source slide deck for version control.

---

# References: 

- [remarkjs](https://remarkjs.com/)
- [revealjs](https://revealjs.com/)
