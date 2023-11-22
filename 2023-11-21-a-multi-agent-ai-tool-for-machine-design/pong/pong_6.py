# Here's a complete pong game implementation based on the provided context. I've incorporated player control for paddles, collision detection with the paddles, the ball bouncing off the top and bottom of the screen, and scoring:

import pygame
from pygame.locals import *

# Constants
WIDTH, HEIGHT = 600, 400
BALL_RADIUS = 20
PADDLE_WIDTH, PADDLE_HEIGHT = 15, 80
FPS = 60
PADDLE_SPEED = 4
BALL_SPEED = 2
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)

# Initialize Pygame
pygame.init()
screen = pygame.display.set_mode((WIDTH, HEIGHT))
clock = pygame.time.Clock()

# Game objects
ball = pygame.Rect(WIDTH // 2, HEIGHT // 2, BALL_RADIUS, BALL_RADIUS)
paddle1 = pygame.Rect(PADDLE_WIDTH, HEIGHT // 2, PADDLE_WIDTH, PADDLE_HEIGHT)
paddle2 = pygame.Rect(WIDTH - (2 * PADDLE_WIDTH), HEIGHT // 2, PADDLE_WIDTH, PADDLE_HEIGHT)

# Initial speeds
ball_dx = BALL_SPEED
ball_dy = BALL_SPEED
paddle1_dy = 0
paddle2_dy = 0

# Scoring
player1_score = 0
player2_score = 0

running = True
while running:
    # Event handling
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

    # Movement update
    ball.x += ball_dx
    ball.y += ball_dy
    paddle1.y += paddle1_dy
    paddle2.y += paddle2_dy

    # Collision handling
    if ball.colliderect(paddle1) or ball.colliderect(paddle2):
        ball_dx *= -1
    elif ball.x <= 0:
        player2_score += 1
        ball.center = (WIDTH // 2, HEIGHT // 2)
        ball_dy = -ball_dy
    elif ball.x >= WIDTH:
        player1_score += 1
        ball.center = (WIDTH // 2, HEIGHT // 2)
        ball_dy = -ball_dy
    if paddle1.y <= 0 or paddle1.y + PADDLE_HEIGHT >= HEIGHT:
        paddle1_dy = 0
    if paddle2.y <= 0 or paddle2.y + PADDLE_HEIGHT >= HEIGHT:
        paddle2_dy = 0

    # Draw everything
    screen.fill(BLACK)
    for entity in [ball, paddle1, paddle2]:
        pygame.draw.rect(screen, WHITE, entity)

    pygame.display.flip()
    clock.tick(FPS)

pygame.quit()
