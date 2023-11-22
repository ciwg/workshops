import pygame
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
ball_dx = ball_speed * pygame.math.cos(45)
ball_dy = ball_speed * pygame.math.sin(45)
player_speed = 2

def ball_movement():
    global ball_dx, ball_dy

    ball.x += ball_dx
    ball.y += ball_dy

    if ball.left <= 0 or ball.right >= screen_width:
        ball_dx *= -1
    if ball.top <= 0 or ball.bottom >= screen_height:
        ball_dy *= -1

    if ball.colliderect(player) or ball.colliderect(opponent):
        ball_dx *= -1

def player_movement():
    player.y += player_speed

    if player.top <= 0:
        player.top = 0
    if player.bottom >= screen_height:
        player.bottom = screen_height

def opponent_movement():
    if opponent.top < ball.y:
        opponent.top += player_speed
    if opponent.bottom > ball.y:
        opponent.bottom -= player_speed
    if opponent.top <= 0:
        opponent.top = 0
    if opponent.bottom >= screen_height:
        opponent.bottom = screen_height

while True:
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            sys.exit()

    ball_movement()
    player_movement()
    opponent_movement()

    screen.fill(bg_color)
    pygame.draw.rect(screen, light_grey, player)
    pygame.draw.rect(screen, light_grey, opponent)
    pygame.draw.ellipse(screen, light_grey, ball)
    pygame.draw.aaline(screen, light_grey, (screen_width/2,0),(screen_width/2,screen_height))

    pygame.display.flip()
    clock.tick(60)

'''
This code creates a simple pong game where a ball bounces off the paddles. The `ball_movement` function checks if the ball is colliding with the paddles or the window borders and reverses direction in case of collision. The `player_movement` and `opponent_movement` functions regulate the movement of player and opponent paddles, making sure they don't move outside of the window.
'''
