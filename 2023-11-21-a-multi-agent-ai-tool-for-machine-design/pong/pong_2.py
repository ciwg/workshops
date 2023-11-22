'''
Your request is to increase the speed of the paddles. The speed of the paddles is determined by the variable `PADDLE_SPEED`. In your existing program, `PADDLE_SPEED` is set to 2. You can increase the paddle speed by increasing this value. Here is the modified program with the speed of the paddles set to 4:
'''

import pygame
from pygame.locals import *

# Initialize Pygame
pygame.init()

# Constant setup
WIDTH, HEIGHT = 600, 400
BALL_RADIUS = 20
PADDLE_WIDTH, PADDLE_HEIGHT = 15, 80
FPS = 60
PADDLE_SPEED = 4 # <-- Changed speed here
BALL_SPEED = 2
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)

# Object setup
ball = pygame.rect.Rect(WIDTH // 2, HEIGHT // 2, BALL_RADIUS, BALL_RADIUS)
paddle1 = pygame.rect.Rect(PADDLE_WIDTH, HEIGHT // 2, PADDLE_WIDTH, PADDLE_HEIGHT)
paddle2 = pygame.rect.Rect(WIDTH - (2 * PADDLE_WIDTH), HEIGHT // 2, PADDLE_WIDTH, PADDLE_HEIGHT)

ball_dx = BALL_SPEED
ball_dy = BALL_SPEED
paddle1_dy = 0
paddle2_dy = 0
clock = pygame.time.Clock()

# Scoring
player1_score = 0
player2_score = 0

screen = pygame.display.set_mode((WIDTH, HEIGHT))

running = True
while running:
    # Process events
    for event in pygame.event.get():
        if event.type == QUIT:
            running = False
        elif event.type == KEYDOWN:
            if event.key == K_UP:
                paddle2_dy = -PADDLE_SPEED
            elif event.key == K_DOWN:
                paddle2_dy = PADDLE_SPEED
            elif event.key == K_w:
                paddle1_dy = -PADDLE_SPEED
            elif event.key == K_s:
                paddle1_dy = PADDLE_SPEED
        elif event.type == KEYUP:
            if event.key == K_UP and paddle2_dy < 0:
                paddle2_dy = 0
            elif event.key == K_DOWN and paddle2_dy > 0:
                paddle2_dy = 0
            elif event.key == K_w and paddle1_dy < 0:
                paddle1_dy = 0
            elif event.key == K_s and paddle1_dy > 0:
                paddle1_dy = 0

    # Update game state
    ball.move_ip(ball_dx, ball_dy)
    paddle1.move_ip(0, paddle1_dy)
    paddle2.move_ip(0, paddle2_dy)

    # Handle collisions
    if ball.left < 0 or ball.right > WIDTH:
        ball_dx *= -1
    if ball.top < 0 or ball.bottom > HEIGHT:
        ball_dy *= -1
    if paddle1.top < 0 or paddle1.bottom > HEIGHT:
        paddle1_dy = 0
    if paddle2.top < 0 or paddle2.bottom > HEIGHT:
        paddle2_dy = 0

    # Drawing
    screen.fill(BLACK)
    pygame.draw.rect(screen, WHITE, ball)
    pygame.draw.rect(screen, WHITE, paddle1)
    pygame.draw.rect(screen, WHITE, paddle2)

    pygame.display.flip()
    clock.tick(FPS)

pygame.quit()
