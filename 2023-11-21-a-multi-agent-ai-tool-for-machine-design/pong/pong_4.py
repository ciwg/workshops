# Your current program is trying to use the `pygame.math.cos()` and `pygame.math.sin()` functions which do not exist in Pygame module. Instead, Python's built-in `math` module should be used for these trigonometric functions. Here is the corrected program:

import pygame
import math
import sys

# General setup
pygame.init()
clock = pygame.time.Clock()

# Main Window
screen_width = 640
screen_height = 480
screen = pygame.display.set_mode((screen_width,screen_height))
pygame.display.set_caption('Pong Game')

# Game rectangles
ball = pygame.Rect(screen_width/2 - 15/2, screen_height/2 - 15/2, 15, 15)
player = pygame.Rect(screen_width - 20, screen_height/2 - 70/2, 10,70)
opponent = pygame.Rect(10, screen_height/2 - 70/2, 10,70)

bg_color = pygame.Color('grey12')
light_grey = (200,200,200)

ball_speed = 3
ball_dx = ball_speed * math.cos(45)
ball_dy = ball_speed * math.sin(45)
player_speed = 2

def ball_animation():
    global ball_speed, ball_dx, ball_dy
    ball.x += ball_dx
    ball.y += ball_dy

    if ball.left <= 0 or ball.right >= screen_width:
        ball_dx *= -1
    if ball.top <= 0 or ball.bottom >= screen_height:
        ball_dy *= -1

    if ball.colliderect(player) or ball.colliderect(opponent):
        ball_dx *= -1

def player_animation():
    player.y += player_speed

    if player.top <= 0:
        player.top = 0
    if player.bottom >= screen_height:
        player.bottom = screen_height

def opponent_ai():
    if opponent.top < ball.y:
        opponent.top += player_speed
    if opponent.bottom > ball.y:
        opponent.bottom -= player_speed
    if opponent.top <= 0:
        opponent.top = 0
    if opponent.bottom >= screen_height:
        opponent.bottom = screen_height

# Game loop
while True:
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            sys.exit()

    ball_animation()
    player_animation()
    opponent_ai()

    screen.fill(bg_color)
    pygame.draw.rect(screen, light_grey, player)
    pygame.draw.rect(screen, light_grey, opponent)
    pygame.draw.ellipse(screen, light_grey, ball)
    pygame.draw.aaline(screen, light_grey, (screen_width/2,0),(screen_width/2,screen_height))

    pygame.display.flip()
    clock.tick(60)

# Please note that math.sin() and math.cos() functions take argument in radians, not in degrees. If you want to go with degrees instead of radians, you have to convert the degrees into radians like so: `ball_speed * math.cos(math.radians(45))`, and `ball_speed * math.sin(math.radians(45))`.
