# Mogi
A Simple Immediate-Mode UI Library in GO
<!-- add image -->

![Mogi Logo](mogi.png)
## TODOs
0. padding, margin->  border, border radius -> gap
1. Z-Index in layout while calcuting pos and size
2. gcc renderer_demo.c  renderer.c -Iinclude -Iexternal\glfw -Iexternal\glad -Llib\Release -lrenderer -lglfw3 -lopengl32 -lgdi32 (in nob)
3. text measure implementation
4. nob options: debug, release, build with optimization details
5. fix fps
6. always correct asset path
7. prevent neg. path
block
inline
is border prt of div size?
box-sizing: content-box, border-box
TODO: (Optionally, but recommended) Update the C Renderer struct to cache current_width and current_height, and have get_window_size return these cached values, which are updated by the callback.
1. system to track jitters or frame drops
2. software renderer + SIMD
3. Performance Implication: AddChildren(children...)
4. fps timeseris graph
5. absolute -> relative -> absolute example
5. SetSizeMode(px, percentage)
6. Stack Optimization
7. Nice API i think, SetLayout(common.NewGridLayout().
			SetColumns(2).
			SetRows(2).
			SetColumnGap(10).
			SetRowGap(10).
			SetColumnStretch(1, 1).
			SetRowStretch(1, 1).
			SetColumnStretch(2, 1).
			SetRowStretch(2, 1).
			SetColumnStretch(3, 1).
8. internal component should be private
9. compoent ctx in some methods, for example i need text component size to manually set pos
10. easy way to achieve this: // Standard C code like loops etc work inside components
11. verify if draw_font -> font correct size
12. UUID ID
13. NewWindow
14. easy to make custom component
15. tailwind colors
16. cache absolute position?
17. modern opengl
18. border,padding,border - t,b,r,l
19. renderer interface
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

DevTools
1. Onhover -> border, pos, size

