#include <stdio.h>
#include "include/renderer.h"  // Adjust path to your renderer.h file

int main() {
    // Create the renderer (window)
    void* renderer = create_renderer(800, 600, "Test Renderer");

    if (renderer == NULL) {
        printf("Failed to create renderer\n");
        return -1;
    }

    while (!window_should_close(renderer)) {
        clear_screen(renderer, (ColorRGBA){0.0f, 0.0f, 0.0f, 1.0f}); // Clear screen with black

        draw_rectangle(renderer, (Rect){(Vec2){100, 100}, 200, 150}, (ColorRGBA){1.0f, 0.0f, 0.0f, 1.0f}); // Draw red rectangle

        present_screen(renderer);
        
        handle_events(renderer);
    }

    destroy_renderer(renderer);
    return 0;
}
