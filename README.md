## Run
```
Run: go run ./src/.
Mac Build: env GOOS=js GOARCH=wasm go build -o ./dist/main.wasm ./src/.
Windows Build: $Env:GOOS='js';$Env:GOARCH='wasm';go build -ldflags="-s -w" -o ../dist/main.wasm
```

## Credits
- Early GameBoy font: https://www.dafont.com/early-gameboy.font
- Kirokaze Gameboy Palette: https://lospec.com/palette-list/kirokaze-gameboy
- GB Mbius Palette: https://lospec.com/palette-list/gb-mbius
- Ice Cream GB Palette: https://lospec.com/palette-list/ice-cream-gb
- Hollow Palette: https://lospec.com/palette-list/hollow

## Todo
- [x] Let the shapes go to the edges
- [x] Visualize the blocks preventing a shape from getting extracted
- [ ] Green shapes moving up should draw over everything
- [ ] Pulse blockign shapes
- [ ] Add single block "shape"
- [ ] Show upcomming shapes
- [ ] Switch which of the upcomming shapes you are placing
- [x] Block placing if trapped under another block
- [ ] Add new lines
- [ ] Score
- [ ] Reset
- [ ] Holding key joggs shape
- [x] Mouse control that plays nicely with keyboard controls