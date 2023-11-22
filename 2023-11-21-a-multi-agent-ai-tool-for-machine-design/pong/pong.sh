#!/bin/bash

# iteration:
echo "i want a pong game in pygame" > requirements.txt
grok msg "you are an expert pygame developer. how would this code be structured?" < requirements.txt > structure.txt
grok msg "you are an export pygame developer. what are the main classes and functions?" < structure.txt > classes.txt
grok msg "you are an expert pygame developer. improve the code." < structure.txt > pong.txt
grok msg "you are a markdown parser.  extract the code from inside the quotes." < pong.txt > pong_1.py
grok msg "you are an expert pygame developer. make the paddles move faster. your answer must be a complete python program." < pong_1.py > pong_2.py
grok msg "you are an expert pygame developer. make the ball bounce off the paddles or leave the screen. your answer must be a complete python program." < pong_2.py > pong_3.py
grok msg "you are an expert pygame developer. fix the following error. your answer must be a complete python program.
    Traceback (most recent call last):
      File 'pong_3.py', line 23, in <module>
        ball_dx = ball_speed * pygame.math.cos(45)
    AttributeError: module 'pygame.math' has no attribute 'cos'
    " < pong_3.py > pong_4.py
grok msg "you are an expert pygame developer. allow players to control the paddles. your answer must be a complete python program." < pong_4.py > pong_5.py
grok add *
echo "fast paddles, bounce off paddles, players control paddles" | grok ctx 4000 > ctx1.txt
grok msg "you are an expert pygame developer. write a complete pong game with fast paddles, bounce off paddles, players control paddles, based on the provided context. your answer must be a complete python program." < ctx1.txt > pong_6.py

# iteration:
grok msg "you are a product manager who can read python code.  write a complete list of requirements for the provided pong game." < pong_6.py > requirements7.txt
grok msg "you are an expert pygame developer. write a complete pong game that meets the provided requirements. your answer must be a complete python program." < requirements7.txt > pong_7.py
grok msg "you are a software tester.  does the provided code meet the following requirements?  $(cat requirements7.txt)" < pong_7.py > testres7.txt
echo "the scores aren't displayed" > humanres7.txt

# iteration:
grok msg "you are an expert pygame developer. your customer says $(cat humanres7.txt). fix the provided code. your answer must be a complete python program." < pong_7.py > pong_8.py









