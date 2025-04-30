# Mogi
A Simple Immediate-Mode UI Library in GO
<!-- add image -->

![Mogi Logo](mogi.png)
## TODOs

1. Z-Index in layout while calcuting pos and size
2. gcc renderer_demo.c  renderer.c -Iinclude -Iexternal\glfw -Iexternal\glad -Llib\Release -lrenderer -lglfw3 -lopengl32 -lgdi32 (in nob)
3. text measure implementation
block
inline
is border prt of div size?
box-sizing: content-box, border-box
TODO: (Optionally, but recommended) Update the C Renderer struct to cache current_width and current_height, and have get_window_size return these cached values, which are updated by the callback.
1. system to track jitters or frame drops
2. software renderer + SIMD
3. Performance Implication: AddChildren(children...)
4. fps timeseris graph
5. SetSizeMode(px, percentage)
// TODO: Benchmarking
// TODO: SIMD implementations
// TODO: bezier curves
// TODO: olivec_ring
// TODO: fuzzer
// TODO: Stencil
🎨 Common uses in 2D UI:
Clipping UI elements (e.g., only render inside a rounded rectangle).

Layered effects (e.g., shadows behind specific shapes).

Masking parts of a texture or image.

Complex shapes without modifying geometry.

