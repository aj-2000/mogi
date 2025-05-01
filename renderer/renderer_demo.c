#include "include/renderer.h" // Adjust path if needed
#include <stdio.h>            // For snprintf, fprintf
#include <stdlib.h>           // For exit, EXIT_FAILURE
#include <math.h>             // For sinf, cosf

// Define M_PI if needed (some compilers require _USE_MATH_DEFINES before #include <math.h>)
#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

// Helper to get mouse position (requires access to GLFWwindow*)
// NOTE: This slightly breaks the abstraction of the void* renderer,
// but is necessary for mouse interaction without adding a dedicated
// input system to the renderer library itself.
void get_mouse_position(void* renderer_ptr, double* xpos, double* ypos) {
    Renderer* ctx = (Renderer*)renderer_ptr; // Cast to access internal struct
    if (ctx && ctx->window) {
        glfwGetCursorPos(ctx->window, xpos, ypos);
    } else {
        *xpos = 0.0;
        *ypos = 0.0;
    }
}

// Helper for simple key state polling
int is_key_pressed(void* renderer_ptr, int key) {
     Renderer* ctx = (Renderer*)renderer_ptr; // Cast to access internal struct
     if (ctx && ctx->window) {
         return glfwGetKey(ctx->window, key) == GLFW_PRESS;
     }
     return 0;
}


int main() {
    int initial_width = 1024;
    int initial_height = 768;
    const char* font_path = "../JetBrainsMonoNL-Regular.ttf"; // ADJUST PATH AS NEEDED

    // --- Initialization ---
    void* renderer = create_renderer(initial_width, initial_height, "Renderer Demo");
    if (!renderer) {
        fprintf(stderr, "Fatal: Failed to create renderer.\n");
        return EXIT_FAILURE;
    }

    // --- Load Fonts ---
    FontData* font_large = load_font(font_path, 32.0f);
    if (!font_large) {
        fprintf(stderr, "Fatal: Failed to load large font: %s\n", font_path);
        destroy_renderer(renderer);
        return EXIT_FAILURE;
    }

    FontData* font_medium = load_font(font_path, 20.0f);
    if (!font_medium) {
        fprintf(stderr, "Fatal: Failed to load medium font: %s\n", font_path);
        destroy_font(font_large);
        destroy_renderer(renderer);
        return EXIT_FAILURE;
    }

     FontData* font_small = load_font(font_path, 14.0f);
    if (!font_small) {
         fprintf(stderr, "Fatal: Failed to load small font: %s\n", font_path);
         destroy_font(font_large);
         destroy_font(font_medium);
         destroy_renderer(renderer);
         return EXIT_FAILURE;
    }

    // --- Colors ---
    ColorRGBA clear_color = {0.05f, 0.05f, 0.1f, 1.0f}; // Dark blueish
    ColorRGBA white = {1.0f, 1.0f, 1.0f, 1.0f};
    ColorRGBA red = {1.0f, 0.2f, 0.2f, 1.0f};
    ColorRGBA green = {0.2f, 1.0f, 0.3f, 1.0f};
    ColorRGBA blue = {0.3f, 0.4f, 1.0f, 1.0f};
    ColorRGBA yellow = {1.0f, 1.0f, 0.3f, 1.0f};
    ColorRGBA magenta = {1.0f, 0.3f, 1.0f, 1.0f};
    ColorRGBA cyan = {0.3f, 1.0f, 1.0f, 1.0f};
    ColorRGBA orange = {1.0f, 0.6f, 0.1f, 1.0f};
    ColorRGBA gray = {0.5f, 0.5f, 0.5f, 1.0f};
    ColorRGBA black = {0.0f, 0.0f, 0.0f, 1.0f};

    // --- State Variables for Animation ---
    float total_time = 0.0f;
    float moving_rect_x = 50.0f;
    float moving_rect_dir = 150.0f; // Pixels per second
    float growing_circle_radius = 20.0f;
    float growing_circle_rate = 30.0f; // Radius change per second
    float rotating_line_angle = 0.0f;
    float rotation_speed = 90.0f; // Degrees per second
    int vsync_enabled = 1; // Assume VSync is initially on (as set in create_renderer)
    int v_key_pressed_last_frame = 0;

    GLuint image_texture = load_texture("../mogi.png"); // Load an image texture


    // --- Main Loop ---
    while (!window_should_close(renderer)) {
        // --- Input & Time ---
        handle_events(renderer); // Process resize, close events, etc.
        float dt = get_delta_time(renderer);
        total_time += dt;
        Vec2 window_size = get_window_size(renderer); // Get current size

        // --- Update Animations ---
        // Moving Rectangle
        moving_rect_x += moving_rect_dir * dt;
        if (moving_rect_x > window_size.x - 150.0f || moving_rect_x < 50.0f) {
            moving_rect_dir *= -1.0f; // Reverse direction
            // Clamp position to prevent going too far off-screen if dt is large
            if (moving_rect_x > window_size.x - 150.0f) moving_rect_x = window_size.x - 150.0f;
            if (moving_rect_x < 50.0f) moving_rect_x = 50.0f;
        }

        // Growing/Shrinking Circle
        growing_circle_radius += growing_circle_rate * dt;
        if (growing_circle_radius > 80.0f || growing_circle_radius < 10.0f) {
            growing_circle_rate *= -1.0f; // Reverse growth
             if (growing_circle_radius > 80.0f) growing_circle_radius = 80.0f;
             if (growing_circle_radius < 10.0f) growing_circle_radius = 10.0f;
        }

        // Rotating Line
        rotating_line_angle += rotation_speed * dt;
        if (rotating_line_angle >= 360.0f) {
            rotating_line_angle -= 360.0f;
        }

        // Color Cycling (simple lerp between red and blue)
        float cycle = (sinf(total_time * 1.5f) + 1.0f) / 2.0f; // 0.0 to 1.0 cycle
        ColorRGBA cycling_color = {
            red.r * (1.0f - cycle) + blue.r * cycle,
            red.g * (1.0f - cycle) + blue.g * cycle,
            red.b * (1.0f - cycle) + blue.b * cycle,
            1.0f
        };

        // VSync Toggle Input
        int v_key_currently_pressed = is_key_pressed(renderer, GLFW_KEY_V);
        if (v_key_currently_pressed && !v_key_pressed_last_frame) {
            vsync_enabled = !vsync_enabled;
            set_vsync(renderer, vsync_enabled);
            printf("VSync %s\n", vsync_enabled ? "Enabled" : "Disabled");
        }
        v_key_pressed_last_frame = v_key_currently_pressed;


        // --- Rendering ---
        clear_screen(renderer, clear_color);

        // --- Draw Rectangles ---
        Rect static_rect = {{window_size.x - 220.0f, 50.0f}, 200.0f, 80.0f};
        draw_rectangle_filled(renderer, static_rect, gray);

        Rect moving_outline_rect = {{moving_rect_x, 50.0f}, 100.0f, 60.0f};
        draw_rectangle_outline(renderer, moving_outline_rect, cyan);

        Rect cycling_rect = {{50.0f, window_size.y - 100.0f}, 150.0f, 50.0f};
        draw_rectangle_filled_outline(renderer, cycling_rect, cycling_color, white);


        // --- Draw Circles ---
        Circle static_circle = {{100.0f, 180.0f}, 40.0f};
        draw_circle_filled(renderer, static_circle, magenta);

        Circle growing_circle = {{window_size.x - 120.0f, window_size.y - 120.0f}, growing_circle_radius};
        draw_circle_outline(renderer, growing_circle, yellow);

        // Circle following mouse
        double mouse_x, mouse_y;
        get_mouse_position(renderer, &mouse_x, &mouse_y);
        Circle mouse_circle = {{(float)mouse_x, (float)mouse_y}, 25.0f};
        draw_circle_filled_outline(renderer, mouse_circle, blue, white);


        // --- Draw Lines ---
        Line thick_line = {{250.0f, 150.0f}, {450.0f, 250.0f}};
        draw_line_thick(renderer, thick_line, orange, 5.0f); // 5px thick

        // Rotating dashed line
        Vec2 center = {window_size.x / 2.0f, window_size.y / 2.0f};
        float angle_rad = rotating_line_angle * (float)M_PI / 180.0f;
        float line_len = 150.0f;
        Line rotating_line = {
            center,
            {center.x + cosf(angle_rad) * line_len, center.y + sinf(angle_rad) * line_len}
        };
        draw_line_dashed(renderer, rotating_line, green, 10.0f, 5.0f); // 10px dash, 5px gap

        // Dotted line
        Line dotted_line_def = {{50.0f, window_size.y - 150.0f}, {window_size.x - 50.0f, window_size.y - 150.0f}};
        draw_line_dotted(renderer, dotted_line_def, white, 2.0f, 1.8f); // 2px radius dots, 1.8x gap


        // --- Draw Text ---
        // Static text
        draw_text(renderer, font_large, "Renderer Demo", (Vec2){20.0f, 20.0f}, white);

        // Dynamic text (FPS/Delta Time)
        char fps_buffer[100];
        float fps = (dt > 0) ? (1.0f / dt) : 0.0f;
        snprintf(fps_buffer, sizeof(fps_buffer), "FPS: %.1f (%.3f ms) VSync: %s (Press V)",
                 fps, dt * 1000.0f, vsync_enabled ? "ON" : "OFF");
        // Position FPS counter at top-right using text width calculation
        float fps_text_width = calculate_text_width(font_small, fps_buffer);
        Vec2 fps_pos = {window_size.x - fps_text_width - 10.0f, 10.0f};
        draw_text(renderer, font_small, fps_buffer, fps_pos, yellow);

        // Text inside a rectangle (using calculate_text_width for centering)
        const char* centered_text = "Centered Text";
        float text_w = calculate_text_width(font_medium, centered_text);
        float text_h = font_medium->ascent - font_medium->descent; // Approx height using metrics
        Vec2 text_pos_centered = {
            static_rect.position.x + (static_rect.width - text_w) / 2.0f,
            static_rect.position.y + (static_rect.height - text_h) / 2.0f // Basic vertical center
        };
        draw_text(renderer, font_medium, centered_text, text_pos_centered, black); // Black text on gray rect

        // More static text examples
        draw_text(renderer, font_medium, "Medium Font Example", (Vec2){20.0f, 300.0f}, cyan);
        draw_text(renderer, font_small, "Small Font Example - 0123456789 !@#$^&*()", (Vec2){20.0f, 330.0f}, green);

        // Image example
        Rect image_rect = {{window_size.x - 200.0f, window_size.y - 200.0f}, 180.0f, 180.0f};
        draw_texture(renderer, image_texture, image_rect, white); // Draw the image
        
        // --- Present ---
        present_screen(renderer);
    }

    free_texture(image_texture); // Free the texture when done


    // --- Cleanup ---
    printf("Cleaning up...\n");
    destroy_font(font_large);
    destroy_font(font_medium);
    destroy_font(font_small);
    destroy_renderer(renderer);
    printf("Cleanup complete. Exiting.\n");

    return 0;
}
