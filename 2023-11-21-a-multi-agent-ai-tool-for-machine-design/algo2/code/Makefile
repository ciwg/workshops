all: code.sh

tasks.md: stories.md bugs.md
	./manager.sh stories.md bugs.md > tasks.md

code.sh: tasks.md 
	./developer.sh tasks.md code.sh > newcode.sh
	mv newcode.sh code.sh

bugs.md: stories.md code.sh
	./tester.sh stories.md code.sh > bugs.md

