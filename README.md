## Run
```
Run: go run ./src/.
Mac Build: env GOOS=js GOARCH=wasm go build -o ./dist/main.wasm ./src/.
Windows Build: $Env:GOOS='js';$Env:GOARCH='wasm';go build -ldflags="-s -w" -o ../dist/main.wasm
```

## Controls
- Mouse or W-A-S-D to navigate the menu move the shape
- Q-E to rotate the shape
- Mouse click or SPACE to place the shape
- R to reset the current level
- 1,2,3,4 to change the color palette 

## Credits
- Ebitengine: https://ebitengine.org/
- Early GameBoy font: https://www.dafont.com/early-gameboy.font
- Tkachevica-4px font: https://www.dafontfree.net/tkachevica-4px-regular/f55920.htm
- Hollow Palette (key 2): https://lospec.com/palette-list/hollow
- Gold BG Palette (key 4): https://lospec.com/palette-list/gold-gb
- Lospec-gb Palette (key 5): https://lospec.com/palette-list/lospec-gb

## Todo
- [x] Let the shapes go to the edges
- [x] Visualize the blocks preventing a shape from getting extracted
- [x] Shake when blocked
- [x] Make it clear when you hit a dead end "Press R to reset"
- [x] Level hints
- [ ] Green shapes moving up should draw over everything
- [x] Pulse blocking shapes when trying to extract
- [ ] Add single block "shape"
- [x] Show upcoming shapes
- [x] Switch which of the upcoming shapes you are placing
- [x] Block placing if trapped under another block
- [ ] Add new lines
- [ ] Score
- [x] Reset
- [ ] Holding key jogs shape
- [x] Mouse control that plays nicely with keyboard controls
- [x] Hook up play menu item
- [ ] Help screen
- [ ] Credits screen
- [ ] Logo
- [ ] SoundFX