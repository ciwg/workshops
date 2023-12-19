#!/bin/bash

docker start gpt_4_container

docker exec -i gpt_4_container ./gpt-4

docker stop gpt_4_container
