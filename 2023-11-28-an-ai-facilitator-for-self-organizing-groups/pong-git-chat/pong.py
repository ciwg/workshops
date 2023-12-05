
import pygame
import sys
from pygame.locals import *

# General constants
WIDTH = 600
HEIGHT = 400

# Colors constants
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)

# Paddle constants
PADDLE_WIDTH = 15
PADDLE_HEIGHT = 80
PADDLE_SPEED = 2

# Ball constants
BALL_DIAMETER = 15
BALL_SPEED = 2

# FPS constants
FPS = 60

pygame.init()

# Set up some stuff
window = pygame.display.set_mode((WIDTH, HEIGHT))
clock = pygame.time.Clock()

# Score variables
score1 = 0
score2 = 0

pygame.font.init()
myfont = pygame.font.SysFont('Comic Sans MS', 30)

# Initialize ball and paddle objects
ball = pygame.Rect(WIDTH // 2, HEIGHT // 2, BALL_DIAMETER, BALL_DIAMETER)
paddle1 = pygame.Rect(0, HEIGHT // 2, PADDLE_WIDTH, PADDLE_HEIGHT)
paddle2 = pygame.Rect(WIDTH - PADDLE_WIDTH, HEIGHT // 2, PADDLE_WIDTH, PADDLE_HEIGHT)

# Set initial velocities
ball_dx = ball_dy = BALL_SPEED
paddle1_dy = paddle2_dy = 0

# Main game loop
while True:
    for event in pygame.event.get():
        if event.type == QUIT:
            pygame.quit()
            sys.exit()
        elif event.type == KEYDOWN:
            if event.key == K_w:
                paddle1_dy = -PADDLE_SPEED
            elif event.key == K_s:
                paddle1_dy = PADDLE_SPEED
            elif event.key == K_UP:
                paddle2_dy = -PADDLE_SPEED
            elif event.key == K_DOWN:
                paddle2_dy = PADDLE_SPEED
        elif event.type == KEYUP:
            if event.key in (K_w, K_s):
                paddle1_dy = 0
            elif event.key in (K_UP, K_DOWN):
                paddle2_dy = 0

    # Update ball and paddle positions
    ball.move_ip(ball_dx, ball_dy)
    paddle1.move_ip(0, paddle1_dy)
    paddle2.move_ip(0, paddle2_dy)

    # Collision detection
    if ball.colliderect(paddle1) or ball.colliderect(paddle2):
        ball_dx *= -1
    elif ball.top <= 0 or ball.bottom >= HEIGHT:
        ball_dy *= -1
    
    # Score when the ball goes past a paddle
    if ball.left <= 0:  # Player 2 scores
        score2 += 1
        ball.center = (WIDTH // 2, HEIGHT // 2)
        ball_dx *= -1
    elif ball.right >= WIDTH:  # Player 1 scores
        score1 += 1
        ball.center = (WIDTH // 2, HEIGHT // 2)
        ball_dx *= -1

    # Update the display
    window.fill(BLACK)
    pygame.draw.rect(window, WHITE, paddle1)
    pygame.draw.rect(window, WHITE, paddle2)
    pygame.draw.ellipse(window, WHITE, ball)
    pygame.draw.aaline(window, WHITE, (WIDTH // 2, 0), (WIDTH // 2, HEIGHT))
    scoretext = myfont.render(f'{score1} score:score score {score2}', False, WHITE)
    window.blit(scoretext, (WIDTH // 2 - scoretext.get_width() // 2, 10))  # render score

    pygame.display.update()
    clock.tick(FPS) # Implement FPS cap
