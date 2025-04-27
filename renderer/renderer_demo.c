#include "include/renderer.h"
#include <stdio.h> // For printf

int main() {
    int screen_width = 800;
    int screen_height = 600;

    void* renderer = create_renderer(screen_width, screen_height, "Text Rendering Test");
    if (!renderer) {
        return -1;
    }

    // --- Load a font ---
    // Make sure you have a TTF font file (e.g., "arial.ttf") accessible
    // Adjust the path as needed.
    FontData* my_font = load_font("../JetBrainsMonoNL-Regular.ttf", 32.0f); // Load Arial at 32 pixels height
    if (!my_font) {
        fprintf(stderr, "Failed to load font.\n");
        destroy_renderer(renderer);
        return -1;
    }

    FontData* small_font = load_font("../JetBrainsMonoNL-Regular.ttf", 16.0f); // Load Consolas at 16 pixels
    if (!small_font) {
         fprintf(stderr, "Failed to load small font.\n");
         destroy_font(my_font); // Clean up previously loaded font
         destroy_renderer(renderer);
         return -1;
    }


    ColorRGBA clear_color = {0.1f, 0.1f, 0.15f, 1.0f};
    ColorRGBA rect_color = {0.8f, 0.2f, 0.2f, 1.0f};
    ColorRGBA text_color = {1.0f, 1.0f, 1.0f, 1.0f}; // White text
    ColorRGBA green_text = {0.2f, 0.9f, 0.3f, 1.0f}; // Green text

    Rect my_rect = {{100.0f, 400.0f}, 200.0f, 50.0f};

    // Main loop
    while (!window_should_close(renderer)) {
        handle_events(renderer);

        clear_screen(renderer, clear_color);

        // Draw other shapesfonts/consola.tt
        draw_rectangle(renderer, my_rect, rect_color);

        // --- Draw text ---
        Vec2 pos1 = {50.0f, 100.0f};
        draw_text(renderer, my_font, "Hello, stb_truetype!", pos1, text_color);

        Vec2 pos2 = {50.0f, 150.0f};
        draw_text(renderer, my_font, "Testing 123 @#$^&*", pos2, green_text);

        char buffer[100];
        snprintf(buffer, sizeof(buffer), "Small Font Example (Size: %.0f)", small_font->font_height_pixels);

        Vec2 pos3 = {50.0f, 250.0f};
        draw_text(renderer, small_font, buffer, pos3, text_color);

        Vec2 pos4 = {50.0f, 280.0f};
        draw_text(renderer, small_font, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", pos4, text_color);

        Vec2 pos5 = {50.0f, 300.0f};
        draw_text(renderer, small_font, "abcdefghijklmnopqrstuvwxyz", pos5, text_color);

        Vec2 pos6 = {50.0f, 320.0f};
        draw_text(renderer, small_font, "0123456789 .,:;!?()[]{}", pos6, text_color);


        present_screen(renderer);
        handle_events(renderer);
    }

    // --- Clean up ---
    destroy_font(my_font);
    destroy_font(small_font);
    destroy_renderer(renderer);

    return 0;
}
