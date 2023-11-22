# Here is a Python game meeting your requirements.

import pygame, random

# Constants for the game
WIDTH, HEIGHT = 600, 400
BALL_RADIUS = 20
PADDLE_WIDTH, PADDLE_HEIGHT = 15, 80
BALL_SPEED = 2
PADDLE_SPEED = 4
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)
FPS = 60

pygame.init()
clock = pygame.time.Clock()
screen = pygame.display.set_mode((WIDTH, HEIGHT))

def draw_paddle(x, y):
    pygame.draw.rect(screen, WHITE, pygame.Rect(x, y, PADDLE_WIDTH, PADDLE_HEIGHT))

def draw_ball(x, y):
    pygame.draw.circle(screen, WHITE, (x, y), BALL_RADIUS)

def game():
    x1, y1 = 0, HEIGHT // 2
    x2, y2 = WIDTH - PADDLE_WIDTH, HEIGHT // 2
    ball_x, ball_y = WIDTH // 2, HEIGHT // 2
    ball_dx = ball_dy = BALL_SPEED
    score1, score2 = 0, 0

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

        keys = pygame.key.get_pressed()
        if keys[pygame.K_w]:
            y1 = max(y1 - PADDLE_SPEED, 0)
        if keys[pygame.K_s]:
            y1 = min(y1 + PADDLE_SPEED, HEIGHT - PADDLE_HEIGHT)
        if keys[pygame.K_UP]:
            y2 = max(y2 - PADDLE_SPEED, 0)
        if keys[pygame.K_DOWN]:
            y2 = min(y2 + PADDLE_SPEED, HEIGHT - PADDLE_HEIGHT)

        ball_x += ball_dx
        ball_y += ball_dy

        if ball_y - BALL_RADIUS < 0 or ball_y + BALL_RADIUS > HEIGHT:
            ball_dy *= -1
        if ball_x - BALL_RADIUS < 0:
            if y1 < ball_y < y1 + PADDLE_HEIGHT:
                ball_dx *= -1
            else:
                ball_x, ball_y = WIDTH // 2, HEIGHT // 2
                score2 += 1
        elif ball_x + BALL_RADIUS > WIDTH:
            if y2 < ball_y < y2 + PADDLE_HEIGHT:
                ball_dx *= -1
            else:
                ball_x, ball_y = WIDTH // 2, HEIGHT // 2
                score1 += 1

        screen.fill(BLACK)
        draw_paddle(x1, y1)
        draw_paddle(x2, y2)
        draw_ball(ball_x, ball_y)
        pygame.display.flip()
        clock.tick(FPS)
    
    pygame.quit()
    
if __name__ == "__main__":
    game()

# This program creates a pong game where two players can play against each other. They can control their paddles using the 'w', 's' keys and the 'up', 'down' arrow keys. The game has a scoring system and the speed of the game objects is regulated according to the specified constants. The game runs on a main game loop which ends when the game window is closed.
