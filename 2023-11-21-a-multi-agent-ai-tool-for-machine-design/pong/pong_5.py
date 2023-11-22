import pygame
import sys

# General setup
pygame.init()
clock = pygame.time.Clock()

# Main Window
screen_width = 800
screen_height = 600
screen = pygame.display.set_mode((screen_width, screen_height))

# Game rectangles
ball = pygame.Rect(screen_width/2 - 15/2, screen_height/2 - 15/2, 15, 15)
player = pygame.Rect(screen_width - 20, screen_height/2 - 70/2, 10,70)
opponent = pygame.Rect(10, screen_height/2 - 70/2, 10, 70)

bg_color = pygame.Color('grey12')
light_grey = (200, 200, 200)

player_speed = 2
opponent_speed = 2

def draw_objects():
    screen.fill(bg_color)
    pygame.draw.rect(screen, light_grey, player)
    pygame.draw.rect(screen, light_grey, opponent)
    pygame.draw.ellipse(screen, light_grey, ball)
    pygame.draw.aaline(screen, light_grey, (screen_width/2,0), (screen_width/2,screen_height))

def player_movement():
    keys = pygame.key.get_pressed()
    if keys[pygame.K_UP]:
        player.move_ip(0, -player_speed)
    if keys[pygame.K_DOWN]:
        player.move_ip(0, player_speed)
    if player.top <= 0:
        player.top = 0
    if player.bottom >= screen_height:
        player.bottom = screen_height

def opponent_movement():
    if opponent.top < ball.y:
        opponent.move_ip(0, opponent_speed)
    if opponent.bottom > ball.y:
        opponent.move_ip(0, -opponent_speed)
    if opponent.top <= 0:
        opponent.top = 0
    if opponent.bottom >= screen_height:
        opponent.bottom = screen_height

while True:
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            sys.exit()

    draw_objects()
    player_movement()
    opponent_movement()

    pygame.display.flip()
    clock.tick(60)

# This program allows the user to control the player paddle using the 'Up' and 'Down' arrow keys. The opponent paddle is controlled automatically with a simple AI. If the ball is above the paddle's center, the paddle moves upward. Otherwise, the paddle moves downward.
